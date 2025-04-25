package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	// "strconv"
)

type Element struct {
	Name     string   `json:"name"`
	Recipes  []string `json:"recipes"`
	Image    string   `json:"image"`
	PageURL  string   `json:"page_url"`
	Tier	 int      `json:"tier"`
}

const BASE_URL = "https://little-alchemy.fandom.com"
const ELEMENTS_URL = BASE_URL + "/wiki/Elements_(Little_Alchemy_2)"
const HEADERS = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"

func sanitizeFilename(name string) string {
	re := regexp.MustCompile(`[\\/*?:"<>|]`)
	return re.ReplaceAllString(name, "")
}

func fetchImageFromElementPage(name, pageURL string) string {
	res, err := http.Get(pageURL)
	if err != nil {
		log.Println("Error fetching page:", err)
		return ""
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println("Error parsing HTML:", err)
		return ""
	}

	// Cari gambar di halaman elemen
	imgTag := doc.Find("img.pi-image-thumbnail")
	if imgTag.Length() > 0 {
		imgURL, exists := imgTag.Attr("src")
		if exists {
			// Mengunduh gambar
			imgData, err := http.Get(imgURL)
			if err != nil {
				log.Println("Error fetching image:", err)
				return ""
			}
			defer imgData.Body.Close()

			// Simpan gambar	
			imgFilename := sanitizeFilename(name) + ".png"
			file, err := os.Create("images/" + imgFilename)
			if err != nil {
				log.Println("Error saving image:", err)
				return ""
			}
			defer file.Close()
			_, err = io.Copy(file, imgData.Body)
			if err != nil {
				log.Println("Error copying image data:", err)
				return ""
			}
			return imgFilename
		}
	}
	return ""
}

func main() {
	// Setup kolektor
	c := colly.NewCollector(
		colly.UserAgent(HEADERS),
	)

	var elements []Element

	tableCount := 1 
	currentTier:= -1 
	// Scrape halaman utama
	c.OnHTML("table.list-table.col-list.icon-hover", func(e *colly.HTMLElement) {
		if(tableCount == 2) { // skip Time
			tableCount++
		} else {
			currentTier++
			tableCount++
		}
		e.ForEach("tr", func(i int, el *colly.HTMLElement) {
			if i == 0 {
				return
			}
			name := el.ChildText("td:nth-child(1) a")
			// fmt.Println(name)
			if name == "" {
				return
			}

			// Ambil URL elemen
			link := ""
			el.DOM.Find("td:nth-child(2) a").EachWithBreak(func(i int, s *goquery.Selection) bool {
				if goquery.NodeName(s.Parent()) != "span" {
					link, _ = s.Attr("href")
					return false // stop setelah ketemu yang pertama bukan dalam <span>
				}
				return true
			})

			if link == "" {
				return
			}

			elementPageURL := BASE_URL + link
			fmt.Println("Scraping:", elementPageURL)

			// Coba ambil gambar dari halaman elemen
			// imgFilename := fetchImageFromElementPage(name, elementPageURL)
			// // if imgFilename == "" {
			// // 	// Gambar tidak ditemukan dari halaman elemen, coba ambil dari tabel
			// // 	imgTag := el.ChildAttr("td:nth-child(1) img", "src")
			// // 	if imgTag != "" {
			// // 		imgURL := "https:" + imgTag
			// // 		// Mengunduh gambar
			// // 		imgData, err := http.Get(imgURL)
			// // 		if err != nil {
			// // 			log.Println("Error fetching image from table:", err)
			// // 			return
			// // 		}
			// // 		defer imgData.Body.Close()

			// // 		// Simpan gambar
			// // 		imgFilename := sanitizeFilename(name) + ".png"
			// // 		file, err := os.Create("images/" + imgFilename)
			// // 		if err != nil {
			// // 			log.Println("Error saving image:", err)
			// // 			return
			// // 		}
			// // 		defer file.Close()
			// // 		_, err = io.Copy(file, imgData.Body)
			// // 		if err != nil {
			// // 			log.Println("Error copying image data:", err)
			// // 		}
			// // 	}
			// // }

			// Ambil resep elemen
			recipeList := []string{}
			el.DOM.Find("td:nth-child(2) li").Each(func(i int, li *goquery.Selection) {
				recipeParts := []string{}
				li.Find("a").Each(func(j int, aTag *goquery.Selection) {
					if aTag.Parent().Is("span") {
						return
					}
		
					recipeParts = append(recipeParts, aTag.Text())
				})
				if len(recipeParts) >= 2 {
					recipeList = append(recipeList, fmt.Sprintf("%s + %s", recipeParts[0], recipeParts[1]))
				}
			})
		

			fmt.Println(recipeList)

			elements = append(elements, Element{
				Name:    name,
				Recipes: recipeList,
				Image:   "images/" + name + ".png",
				PageURL: elementPageURL,
				Tier: currentTier,
			})

			// Simpan hasil scraping
			fmt.Printf("Scraped: %s (Tier %d)\n", name, currentTier)
		})
	})

	// Mulai mengunjungi halaman utama
	err := c.Visit(ELEMENTS_URL)
	if err != nil {
		log.Fatal("Error visiting page:", err)
	}

	// Simpan ke file JSON
	file, err := os.Create("elements.json")
	if err != nil {
		log.Fatal("Error creating JSON file:", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // untuk pretty-print
	if err := encoder.Encode(elements); err != nil {
		log.Fatal("Error encoding JSON:", err)
	}

	fmt.Println("\nDone!")
}
