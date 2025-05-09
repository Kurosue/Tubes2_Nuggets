package utils

import (
    "encoding/json"
    "os"
)

// Struct element sesuai json
type Element struct {
    Name     string   `json:"name"`
    Recipes  []string `json:"recipes"`
    Image    string   `json:"image"`
    PageURL  string   `json:"page_url"`
}

// RecipeMap now maps an element's name to its recipes.
type RecipeMap map[string][]string

func LoadRecipes(path string) (RecipeMap, error) {
    var elems []Element
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    if err := json.Unmarshal(data, &elems); err != nil {
        return nil, err
    }

    recipes := make(RecipeMap)
    for _, e := range elems {
        recipes[e.Name] = e.Recipes
    }
    return recipes, nil
}