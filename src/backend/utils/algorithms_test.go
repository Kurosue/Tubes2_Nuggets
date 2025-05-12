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
	test := "Science" // Change this to the element you want to test

	// One recipe Test
	// pathEl, node := DFS(recipes, recipesEl, recipesEl[test])
	// result := DFSResult{Messages: pathEl, NodesVisited: node}
	// t.Logf("Path to %s: ", test)
	// t.Logf("Total nodes: %d", result.NodesVisited)
	// t.Logf("\n%s", VisualizeDFS(result)) // visualization max depth to 10


	// Multiple recipe Test
	pathElMulti := DFSMultiple(recipes, recipesEl, recipesEl[test], 5)
	for i, el := range pathElMulti.RecipePaths {
		t.Logf("recipe: %v", i+1)
		t.Log(VisualizeMessageTree(el)) // visualization max depth to 10
	}
	t.Logf("Total nodes: %d", pathElMulti.NodesVisited)
}

// How to run, use `go test -v .` in the terminal

