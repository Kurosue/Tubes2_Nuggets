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

    // fmt.Println(rese)
    // fmt.Println(elmt)
    // start := "Pyramid"
    node := 0
    res := utils.BFSShortestNode("Picnic", elmt, &node)
    messages := utils.FlattenTreeToMessages(res)
    fmt.Println(node)
    fmt.Println(messages)

    fmt.Println("MultiRecipes")
    node = 0
    resm := utils.BFSNRecipes("Picnic", elmt, 3, &node)
    messagesMultiple := utils.FlattenMultipleTrees(resm)
    fmt.Println(node)
    fmt.Println("Result:")
    for _, path := range messagesMultiple {
        fmt.Println()
        fmt.Println(path)
    }
    // ewaw, i := utils.DFS(elmt, elmt["Wizard"])
    // fmt.Println("DFS ewawult:")
    // fmt.Println(ewaw)
    // fmt.Println("DFS Nodes Visited:")
    // fmt.Println(i)
    // resshort, _ := utils.BFSShortestPath(start, rese, elmt)
    // fmt.Println("Result:")
    // fmt.Println(len(res))
    // fmt.Println(len(resshort))
    // fmt.Println(resshort)
    // for _, path := range res {
    //     fmt.Println(path)
    // }
    
}


