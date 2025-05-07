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
	test := "Brick"
	pathEl := DFSBot(recipes, recipesEl, recipesEl[test])
	t.Logf("Path to %s: ", test)
	// t.Logf("Path: %s", pathEl)
	for _, step := range pathEl {
		t.Logf("Step: %s", step.Name)
	}

	t.Logf("Total steps: %d", len(pathEl))
	// t.Log(recipes["Gold"])
}

// How to run, use `go test -v .` in the terminal

