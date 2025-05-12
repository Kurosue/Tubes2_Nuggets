package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Kurosue/Tubes2_Nuggets/scrap"
	"github.com/Kurosue/Tubes2_Nuggets/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ElementApp struct {
	Name     string   `json:"name"`
    Recipes  []string `json:"recipes"`
    Image    string   `json:"image"`
    PageURL  string   `json:"page_url"`
    Tier     int      `json:"tier"`
	ParsedRecipes [][]string `json:"parsed_recipes"`
}

type AlgorithmResponse struct {
        Recipe      []utils.Message `json:"recipe"`
        NodesVisited int      `json:"nodesVisited"`
        RecipeIndex int       `json:"recipeIndex"`
        TotalRecipes int      `json:"totalRecipes"`
    }

var cachedElements []ElementApp
var cachedRecipesMap utils.RecipeMap
var cachedRecipesEl utils.RecipeElement
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, restrict origin
	},
}

func main() {
	app := gin.Default()

	app.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"GET", "POST", "PATCH", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))

	router := app
	router.Static("/backend/scrap/images", "../scrap/images")

	// Define the API endpoint
	if err := initData(); err != nil {
		panic("Failed to initialize recipes: " + err.Error())
	}

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the Element API!"})
	})

	router.PATCH("/api/rescrape", scrapRecipes)
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
	recipesMap, recipesEl, err := utils.LoadRecipes("../scrap/elements.json")
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
	cachedRecipesEl = recipesEl
	cachedRecipesMap = recipesMap
	return nil
}

var scrapingRecipes bool = false;

func scrapRecipes(c *gin.Context) {
	if scrapingRecipes {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Still scraping" })
		return
	}
	scrapingRecipes = true
	defer func() {
		scrapingRecipes = false
	}()
	scrap.DoScrap(false)
	if err := initData(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize recipes: " + err.Error() })
		return
	}
}

func getElements(c *gin.Context) {
	// Return the cached elements
	c.JSON(http.StatusOK, gin.H{"elements": cachedElements})
}

func findPath(c *gin.Context) {
	// Get query parameters
	var results []AlgorithmResponse
	algorithm := c.Query("algorithm")
	direction := c.Query("direction")
	count, err := strconv.Atoi(c.Query("count"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid count"})
		return
	}
	targetElement := c.Query("target")

	// Validate parameters
	if algorithm != "dfs" && algorithm != "bfs" && algorithm != "bfs-shortest" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid algorithm"})
		return
	}
	if direction != "source" && direction != "target" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid direction"})
		return
	}
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Count is required"})
		return
	}

	if _, exists := cachedRecipesEl[targetElement]; !exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Target element not found"})
        return
    }

	// Upgrade to WebSocket
	wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "WebSocket upgrade failed"})
		return
	}
	defer wsConn.Close()

	if algorithm == "dfs" && direction == "target" {
		// Perform DFS from target to source
		if count == 1 {
			path , node := utils.DFS(cachedRecipesMap, cachedRecipesEl, cachedRecipesEl[targetElement])
			results = append(results, AlgorithmResponse{
				Recipe:      path,
				NodesVisited: node,
				RecipeIndex: 1,
				TotalRecipes: 1,
			})
		} else {
			multiplePath := utils.DFSMultiple(cachedRecipesMap, cachedRecipesEl, cachedRecipesEl[targetElement], int(count))
			for i, recipe := range multiplePath.RecipePaths {
				results = append(results, AlgorithmResponse{
					Recipe:      recipe,
					NodesVisited: multiplePath.NodesVisited,
					RecipeIndex: i + 1,
					TotalRecipes: multiplePath.PathsFound,
				})
			}
		}
	}

	if algorithm == "bfs" && direction == "target" {
		// Perform BFS from target to source
		path, node := utils.BFSP(targetElement, cachedRecipesMap, cachedRecipesEl)
		path = path[:count]
		for i, recipe := range path {
			results = append(results, AlgorithmResponse{
				Recipe:      recipe,
				NodesVisited: node,
				RecipeIndex: i + 1,
				TotalRecipes: len(path),
			})
		}
	}

	if algorithm == "bfs-shortest" && direction == "target" {
		// Perform BFS from target to source
		path, node := utils.BFSShortestPath(targetElement, cachedRecipesMap, cachedRecipesEl)
		results = append(results, AlgorithmResponse{
			Recipe:      path,
			NodesVisited: node,
			RecipeIndex: 1,
			TotalRecipes: 1,
		})
	}

	// Dummy
	for _, result := range results {
        err := wsConn.WriteJSON(result)
        if err != nil {
            // Log error
            break
        }
    }
}
