package main

import (
	"fmt"
	"github.com/Kurosue/Tubes2_Nugget/utils"
)

func main() {
	// Load recipes
	recipes, err := utils.LoadRecipe("./scrap/elements.json")
	if err != nil {
		fmt.Printf("Error loading recipes: %v\n", err)
		return
	}
    // Start BFS from "Nugget"
    start := "Steam"
    fmt.Println(recipes[start])

}
