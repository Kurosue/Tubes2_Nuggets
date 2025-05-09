package main

import (
    "fmt"
    "github.com/Kurosue/Tubes2_Nugget/utils"
)

func main() {
    // Load recipes
    rese, elmt, err := utils.LoadRecipes("./scrap/elements.json")
    if err != nil {
        fmt.Printf("Error loading recipes: %v\n", err)
        return
    }

    start := "Pyramid"
    res := utils.BFSShortestPath(start, rese, elmt)
    if len(res) == 0 {
        fmt.Printf("No path found for %s\n", start)
    } else {
        fmt.Printf("Shortest path to %s:\n", start)
        for _, step := range res {
            fmt.Printf("%s + %s -> %s\n", step.Ingredient1, step.Ingredient2, step.Result)
        }
    }
}


