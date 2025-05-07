package main

import (
	"fmt"
	"log"
	"github.com/Kurosue/Tubes2_Nugget/utils"
)

func main() {
    recipes, _, err := utils.LoadRecipes("./scrap/elements.json")
    if err != nil {
        log.Fatal(err)
    }

    // Example: what *all* elements can you make from Air + Water?
    all := utils.BFS(recipes, "Fire", "Mist")
    fmt.Println("BFS can reach:", all)

    // Example: find a *path* to “Mist” starting from Air + Water
    if path := utils.DFS(recipes, "Air", "Water", "Mist"); path != nil {
        fmt.Println("DFS path to Mist:")
        for _, step := range path {
            fmt.Printf("  %s + %s -> %s\n", step.X, step.Y, step.Result)
        }
    } else {
        fmt.Println("Mist is not reachable from Air + Water.")
    }
}
