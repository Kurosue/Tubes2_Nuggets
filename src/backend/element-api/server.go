package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/Kurosue/Tubes2_Nuggets/scrap"
	"github.com/Kurosue/Tubes2_Nuggets/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var cachedRecipesMap utils.RecipeMap
var cachedRecipesEl utils.RecipeElement
var cachedElements []utils.Element
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, restrict origin
	},
}

func main() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"GET", "POST", "PATCH", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))

	router.Use(func(c *gin.Context) {
		if c.Request.URL.Path[:7] == "/images" {
			c.Header("Cache-Control", "no-cache, must-revalidate, proxy-revalidate")
		}
	})
	router.Static("/images", "../scrap/images")

	// Define the API endpoint
	if err := initData(); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to initialize recipes: %v", err.Error())
	}

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{ "message": "Welcome to the Element API!" })
	})

	router.PATCH("/api/rescrape", scrapRecipes)
	router.GET("/api/elements", func (c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{ "elements": cachedElements })
	})
	router.GET("/api/elements/:name", func(c *gin.Context) {
		name := c.Param("name")
		for _, element := range cachedRecipesEl {
			if element.Name == name {
				c.JSON(http.StatusOK, element)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{ "error": "Element not found" })
	})

	router.GET("/api/elements/:name/image", func(c *gin.Context) {
		name := c.Param("name")
		for _, element := range cachedRecipesEl {
			if element.Name == name {
				c.JSON(http.StatusOK, gin.H{ "image": element.Image })
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{ "error": "Element not found" })
	})

	router.GET("/api/find-recipe", findPath)

	// Start the server
	if err := router.Run(":8888"); err != nil {
		panic(err)
	}
}

func initData() error {
	// Load recipes
	recipesMap, recipesEl, err := utils.LoadRecipes("../scrap/elements.json")
	if err != nil {
		return err
	}
	cachedRecipesEl = recipesEl
	cachedRecipesMap = recipesMap
	cachedElements = nil
	for _, recipe := range recipesEl {
		cachedElements = append(cachedElements, recipe)
	}
	return nil
}

type SafeBuffer struct {
    mu        sync.RWMutex
	closed    bool
    buffer    bytes.Buffer
	closeChan chan bool
	clientsMu sync.Mutex
	clients   map[*gin.Context]bool
}
func NewSafeBuffer() *SafeBuffer {
	return &SafeBuffer{
		closed: false,
		buffer: bytes.Buffer{},
		closeChan: make(chan bool),
		clients: make(map[*gin.Context]bool),
	}
}
func (b *SafeBuffer) Write(p []byte) (n int, err error) {
    b.mu.Lock()
    defer b.mu.Unlock()
	b.clientsMu.Lock()
    defer b.clientsMu.Unlock()
    for c := range b.clients {
        c.Writer.Write(p)
		c.Writer.Flush()
    }
    return b.buffer.Write(p)
}
func (b *SafeBuffer) String() string {
    b.mu.RLock()
    defer b.mu.RUnlock()
    return b.buffer.String()
}
func (b *SafeBuffer) AddClientAndBlock(c *gin.Context) {
	b.mu.RLock()
	_, _ = c.Writer.Write(b.buffer.Bytes())
	b.mu.RUnlock()
	c.Writer.Flush()
	if b.closed {
		return
	}
	b.clientsMu.Lock()
    b.clients[c] = true
    b.clientsMu.Unlock()
	select {
		case <- b.closeChan:
		case <- c.Request.Context().Done():
	}
	b.clientsMu.Lock()
	delete(b.clients, c)
	b.clientsMu.Unlock()
}
func (b *SafeBuffer) Close() {
	if b.closed {
		return
	}
	b.mu.Lock()
    defer b.mu.Unlock()
	b.closed = true
	close(b.closeChan)
}

var scrapBuffer *SafeBuffer = nil

func scrapRecipes(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
	c.Writer.Flush()
	if scrapBuffer != nil {
		scrapBuffer.AddClientAndBlock(c)
		return
	}
	withImage := c.Query("withImage") != ""
	scrapBuffer = NewSafeBuffer()
	go func() {
		iof := io.MultiWriter(os.Stdout, scrapBuffer)
		l := log.New(iof, "[SCRAP] ", log.LstdFlags)
		scrap.DoScrap(withImage, l)
		scrapBuffer.Close()
		scrapBuffer = nil
		if err := initData(); err != nil {
			l.Fatalf("Failed to initialize recipes: %v", err.Error())
		}
	}()
	scrapBuffer.AddClientAndBlock(c)
}

func findPath(c *gin.Context) {
	// Get query parameters
	algorithm := c.Query("algorithm")
	direction := c.Query("direction")
	count, err := strconv.Atoi(c.Query("count"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{ "error": "Invalid count" })
		return
	}
	targetElement := c.Query("target")

	// Validate parameters
	if algorithm != "dfs" && algorithm != "bfs" && algorithm != "bfs-shortest" {
		c.JSON(http.StatusBadRequest, gin.H{ "error": "Invalid algorithm" })
		return
	}
	if direction != "source" && direction != "target" {
		c.JSON(http.StatusBadRequest, gin.H{ "error": "Invalid direction" })
		return
	}
	if count == 0 {
		c.JSON(http.StatusBadRequest, gin.H{ "error": "Count is required" })
		return
	}

	target, targetExists := cachedRecipesEl[targetElement]
	if !targetExists {
        c.JSON(http.StatusBadRequest, gin.H{ "error": "Target element not found" })
        return
    }

	// Upgrade to WebSocket
	wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{ "error": "WebSocket upgrade failed" })
		return
	}
	defer wsConn.Close()

	resultChan := make(chan utils.Message)

	go func() {
		if algorithm == "dfs" && direction == "target" {
			// Perform DFS from target to source
			if count == 1 {
				utils.DFS(cachedRecipesMap, cachedRecipesEl, target, resultChan)
			} else {
				utils.DFSMultiple(cachedRecipesMap, cachedRecipesEl, target, int(count), resultChan)
			}
		}
		if algorithm == "bfs" && direction == "target" {
			// Perform BFS from target to source
			if count == 1 {
				utils.BFSShortestNode(targetElement, cachedRecipesEl, resultChan)
			} else {
				utils.BFSNRecipes(targetElement, cachedRecipesEl, int(count), resultChan)
			}
		}
		close(resultChan)
	}()

	for wsMsg := range resultChan {
		if err := wsConn.WriteJSON(wsMsg); err != nil {
			break;
		}
	}
}
