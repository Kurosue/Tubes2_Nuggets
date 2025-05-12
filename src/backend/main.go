package main

import (
	"fmt"

	"github.com/Kurosue/Tubes2_Nuggets/utils"
)

func main() {
    // Load recipes
    rese, elmt, err := utils.LoadRecipes("./scrap/elements.json")
    if err != nil {
        fmt.Printf("Error loading recipes: %v\n", err)
        return
    }

    start := "Pyramid"
    res, _ := utils.BFSP(start, rese, elmt)
    fmt.Println("Result:")
    for _, path := range res {
        fmt.Println(path)
    }
    
}


