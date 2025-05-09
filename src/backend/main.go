package main

import (
    "fmt"
    "strings"
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
    start := "Pyramid"
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