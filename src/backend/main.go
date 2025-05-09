package main

import (
    "fmt"
    "time"
    "strings"
    "github.com/Kurosue/Tubes2_Nugget/utils"
)

func main() {
    // Load recipes
    recipes, elm, err := utils.LoadRecipes("./scrap/elements.json")
    if err != nil {
        fmt.Printf("Error loading recipes: %v\n", err)
        return
    }
    // timer

    start := "Pyramid"
    
    // Start timer
    startTime := time.Now()
    utils.BFS(start, recipes, elm)
    // Calculate elapsed time
    elapsedTime := time.Since(startTime)
    fmt.Printf("BFS processing time: %v\n", elapsedTime)
    fmt.Printf("Tree for %s (BFS):\n", start)
    // printTree(tree, 0)

    // Start timer
    startTime = time.Now()
    tree := utils.BFSP(start, recipes, elm)
    // Calculate elapsed time
    elapsedTime = time.Since(startTime)
    fmt.Printf("BFSP processing time: %v\n", elapsedTime)
    fmt.Printf("Tree for %s (BFSP):\n", start)
    PrintRecipeTree(tree, 0)
    
    // printTree(treeP, 0)

    startTime = time.Now()
    path := utils.BFSShortestPath(start, recipes, elm)
    // Calculate elapsed time
    elapsedTime = time.Since(startTime)
    fmt.Printf("BFSSP processing time: %v\n", elapsedTime)

    // Print the shortest path
    fmt.Printf("Shortest path for %s:\n", start)
    PrintRecipeTree(path, 0)
}

func PrintRecipeTree(messages []utils.Message, depth int) {
    for _, msg := range messages {
        if depth > 0 {
            fmt.Printf("%s%s + %s = %s\n", strings.Repeat(" ", depth*2), msg.Ingredient1, msg.Ingredient2, msg.Result)
        } else {
            fmt.Printf("%s + %s = %s\n", msg.Ingredient1, msg.Ingredient2, msg.Result)
    tree := utils.BFS_Tree(start, recipes)

    // Print the tree
    fmt.Printf("Tree for %s:\n", start)
    printTree(tree, 0)
      
    if path,node:= utils.DFS(recipes, recipesEl, recipesEl["Gold"]); path != nil {
        fmt.Println("DFS path to Mist:")
        fmt.Print(utils.VisualizeDFS(utils.DFSResult{Messages: path, NodesVisited: node}))
    } else {
        fmt.Println("Mist is not reachable from Air + Water.")
    }
 }
func printTree(node *utils.Tree, level int) {
    if node == nil {
        return
    }
    fmt.Printf("%s%s (Tier: %d)\n", strings.Repeat(" ", level*2), node.Value, node.Tier)
    for _, children := range node.Children {
        for _, child := range children {
            printTree(child, level+1)
        }
    }
}

