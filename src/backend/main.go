package main

import (
	"fmt"
	"github.com/Kurosue/Tubes2_Nugget/utils"
)

func main() {
	// Load recipes
	recipes, err := utils.LoadRecipes("./scrap/elements.json")
	if err != nil {
		fmt.Printf("Error loading recipes: %v\n", err)
		return
	}
    // print some recipes for debugging
    for k, v := range recipes {
        fmt.Printf("%s -> %s\n", k, v)
    }
    // Start BFS from "Nugget"
    start := "Picnic"
    path := utils.BFS(start, recipes)
    // Print the path
    for k, v := range path {
        fmt.Printf("%s -> %s\n", k, v)
    }
}
