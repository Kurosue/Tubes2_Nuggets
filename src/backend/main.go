package main

import (
	"fmt"

	"github.com/Kurosue/Tubes2_Nuggets/utils"
)

func main() {
    // Load recipes
    _, elmt, err := utils.LoadRecipes("./scrap/elements.json")
    if err != nil {
        fmt.Printf("Error loading recipes: %v\n", err)
        return
    }

    start := "Pyramid"
    resultChan := make(chan utils.Message)
    go func() {
        utils.BFSNRecipes(start, elmt, 5, resultChan)
        close(resultChan)
    }()
    fmt.Println("Result:")
    for result := range resultChan {
        fmt.Println(utils.VisualizeMessageTree(result.RecipePath)) // visualization max depth to 10
    }
}
