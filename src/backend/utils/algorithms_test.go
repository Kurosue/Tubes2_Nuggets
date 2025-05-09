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

	// for key, value := range recipes {
	// 	t.Logf("Key: %s, Value: %s", key, value)
	// }

	// for key, value := range recipesEl {
	// 	t.Logf("Key: %s, Value: %v", key, value)
	// }
	test := "Gold"
	pathEl, node := DFS(recipes, recipesEl, recipesEl[test])
	result := DFSResult{Messages: pathEl, NodesVisited: node}
	t.Logf("Path to %s: ", test)
	t.Logf("Nodes visited: %d", node)
	// t.Logf("Path: %v", pathEl)
	// for i, el := range pathEl {
	// 	t.Logf("Step %d: %s", i, el.result)
	// 	t.Logf("Ingredients: %s + %s", el.ingredient1, el.ingredient2)
	// 	t.Logf("Depth: %d", el.depth)
	// }

	t.Logf("\n%s", VisualizeDFS(result))




	// t.Logf("Total steps: %d", len(pathEl))
	// t.Log(recipes["Gold"])
	// pathStr := CreateRecipeTree(pathEl,recipes, recipesEl)
	// t.Log(pathStr)

	// pathElMulti := DFSMulti(recipes, recipesEl, recipesEl[test], 5)
	// for i, el := range pathElMulti {
	// 	t.Logf("recipe: %v", i+1)
	// 	for j, el2 := range el {
	// 		t.Logf("Step %d: %s", j, el2.Name)
	// 	}
	// }
}

// How to run, use `go test -v .` in the terminal

