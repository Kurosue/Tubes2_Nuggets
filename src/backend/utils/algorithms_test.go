package utils

import (
	// "fmt"
	"testing"
)


func TestDFS(t *testing.T) {
	var recipes RecipeMap
	var recipesEl RecipeElement
	var err error
	recipes, recipesEl, err = LoadRecipes("../scrap/elements.json")
	if err != nil {
		t.Fatalf("Failed to load recipes: %v", err)
	}
	test := "Gold" // Change this to the element you want to test

	// One recipe Test
	// pathEl, node := DFS(recipes, recipesEl, recipesEl[test])
	// result := DFSResult{Messages: pathEl, NodesVisited: node}
	// t.Logf("Path to %s: ", test)
	// t.Logf("Total nodes: %d", result.NodesVisited)
	// t.Logf("\n%s", VisualizeDFS(result)) // visualization max depth to 10


	// Multiple recipe Test
	resultChan := make(chan Message)
	go func() {
		DFSMultiple(recipes, recipesEl, recipesEl[test], 5, resultChan)
		close(resultChan)
	}()
	i := 0
	for result := range resultChan {
		t.Logf("recipe: %v, duration %v, visited: %v", i+1, result.Duration, result.NodesVisited)
		t.Log(VisualizeMessageTree(result.RecipePath)) // visualization max depth to 10
		i++
	}
}

// How to run, use `go test -v .` in the terminal
