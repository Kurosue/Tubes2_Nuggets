package main

import (
    "net/http"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
    "github.com/Kurosue/Tubes2_Nugget/utils"
	"strings"
)

type ElementApp struct {
	Name     string   `json:"name"`
    Recipes  []string `json:"recipes"`
    Image    string   `json:"image"`
    PageURL  string   `json:"page_url"`
    Tier     int      `json:"tier"`
	ParsedRecipes [][]string `json:"parsed_recipes"`
}

var cachedElements []ElementApp
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, restrict origin
	},
}

func main() {
	app := gin.Default()
	router := app

	// Define the API endpoint
	if err := initData(); err != nil {
		panic("Failed to initialize recipes: " + err.Error())
	}

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the Element API!"})
	})

	router.PATCH("/api/rescrape", loadRecipes) // ini mau scrap ulang apa gimans ye
	router.GET("/api/elements", getElements)
	router.GET("/api/elements/:name", func(c *gin.Context) {
		name := c.Param("name")
		for _, element := range cachedElements {
			if element.Name == name {
				c.JSON(http.StatusOK, element)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Element not found"})
	})

	router.GET("/api/elements/:name/image", func(c *gin.Context) {
		name := c.Param("name")
		for _, element := range cachedElements {
			if element.Name == name {
				c.JSON(http.StatusOK, gin.H{"image": element.Image})
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "Element not found"})
	})

	router.GET("/api/find-recipe", findPath)
	
	// Start the server
	if err := app.Run(":8888"); err != nil {
		panic(err)
	}
}

func initData() error {
	// Load recipes
	var elementsData []ElementApp
	_, recipesEl, err := utils.LoadRecipes("../scrap/elements.json")
	if err != nil {
		return err
	}

	// Parse the recipes
	for _, recipe := range recipesEl {
		parsedRecipes := make([][]string, 0)
		for _, r := range recipe.Recipes {
			parts := strings.Split(r, "+")
			if len(parts) != 2 {
				continue
			}
			a := strings.TrimSpace(parts[0])
			b := strings.TrimSpace(parts[1])
			parsedRecipes = append(parsedRecipes, []string{a, b})
		}
		temp := ElementApp{
			Name: recipe.Name,
			Recipes: recipe.Recipes,
			Image: recipe.Image,
			PageURL: recipe.PageURL,
			Tier: recipe.Tier,
			ParsedRecipes: parsedRecipes,
		}
		elementsData = append(elementsData, temp)
	}

	cachedElements = elementsData
	return nil
}

func loadRecipes(c *gin.Context) {
	// Load recipes
	var elementsData []ElementApp
	_, recipesEl, err := utils.LoadRecipes("../scrap/elements.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load recipes"})
		return
	}

	// Parse the recipes
	for _, recipe := range recipesEl {
		parsedRecipes := make([][]string, 0)
		for _, r := range recipe.Recipes {
			parts := strings.Split(r, "+")
			if len(parts) != 2 {
				continue
			}
			a := strings.TrimSpace(parts[0])
			b := strings.TrimSpace(parts[1])
			parsedRecipes = append(parsedRecipes, []string{a, b})
		}
		temp := ElementApp{
			Name: recipe.Name,
			Recipes: recipe.Recipes,
			Image: recipe.Image,
			PageURL: recipe.PageURL,
			Tier: recipe.Tier,
			ParsedRecipes: parsedRecipes,
		}
		elementsData = append(elementsData, temp)
	}

	// Return the loaded recipes as JSON
	c.JSON(http.StatusOK, gin.H{"elementsData": elementsData})
}

func getElements(c *gin.Context) {
	// Return the cached elements
	c.JSON(http.StatusOK, gin.H{"elements": cachedElements})
}

func findPath(c *gin.Context) {
	// Get query parameters
	algorithm := c.Query("algorithm")
	direction := c.Query("direction")
	count := c.Query("count")

	// Validate parameters
	if algorithm != "dfs" && algorithm != "bfs" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid algorithm"})
		return
	}
	if direction != "source" && direction != "target" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid direction"})
		return
	}
	if count == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Count is required"})
		return
	}

	// Upgrade to WebSocket
	wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "WebSocket upgrade failed"})
		return
	}
	defer wsConn.Close()

	// Dummy
	messages := []map[string]string{
		{"sourceA": "water", "sourceB": "earth", "target": "mud"},
		{"sourceA": "air", "sourceB": "fire", "target": "energy"},
	}

	for _, msg := range messages {
		err := wsConn.WriteJSON(msg)
		if err != nil {
			break
		}
	}
}