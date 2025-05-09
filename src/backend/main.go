package main

import (
	"fmt"
	"log"
	"github.com/Kurosue/Tubes2_Nugget/utils"
)

func main() {
    recipes, recipesEl , err := utils.LoadRecipes("./scrap/elements.json")
    if err != nil {
        log.Fatal(err)
    }

    // Example: what *all* elements can you make from Air + Water?
    all := utils.BFS(recipes, "Fire", "Mist")
    fmt.Println("BFS can reach:", all)

    // Example: find a *path* to “Mist” starting from Air + Water
    if path,node:= utils.DFS(recipes, recipesEl, recipesEl["Gold"]); path != nil {
        fmt.Println("DFS path to Mist:")
        fmt.Print(utils.VisualizeDFS(utils.DFSResult{Messages: path, NodesVisited: node}))
    } else {
        fmt.Println("Mist is not reachable from Air + Water.")
    }
}
