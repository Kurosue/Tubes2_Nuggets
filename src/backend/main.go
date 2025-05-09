package main

import (
    "fmt"
    "github.com/Kurosue/Tubes2_Nugget/utils"
)

func main() {
    // Load recipes
    _, _, err := utils.LoadRecipes("./scrap/elements.json")
    if err != nil {
        fmt.Printf("Error loading recipes: %v\n", err)
        return
    }
}


