package utils

import (
    "encoding/json"
    "os"
    "strings"
)

type Element struct {
    Name     string   `json:"name"`
    Recipes  []string `json:"recipes"`
    Image    string   `json:"image"`
    PageURL  string   `json:"page_url"`
    Tier     int      `json:"tier"`
    ParsedRecipes []Recipe
}

type Recipe struct {
    First  string
    Second string
}

func DecomposeKey(key string) (string, string) {
    parts := strings.Split(key, "|")
    if len(parts) != 2 {
        return "", ""
    }
    return parts[0], parts[1]
}

func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func LoadRecipes(path string) (RecipeMap, RecipeElement, error)  {
//   ElementMap sama recipesElement tuh sama gaksih?
    type ElementMap map[string]*Element
    var elems []Element
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, nil, err
    }
    if err := json.Unmarshal(data, &elems); err != nil {
        return nil, nil, err
    }

    recipes := make(RecipeMap)
    recipesElement := make(RecipeElement)
    // Create map of elements
    elementMap := make(ElementMap)
    
    // First pass: Add all elements to the map
    for i := range elems {
        elementMap[elems[i].Name] = &elems[i]
    }
    
    // Second pass: Parse recipes for each element
    for _, e := range elems {
        for _, rec := range e.Recipes {
            parts := strings.Split(rec, " + ")
            if len(parts) == 2 {
                recipe := Recipe{
                    First:  strings.TrimSpace(parts[0]),
                    Second: strings.TrimSpace(parts[1]),
                }
                elementMap[e.Name].ParsedRecipes = append(elementMap[e.Name].ParsedRecipes, recipe)
            }
        }
        recipesElement[e.Name] = e
    }
    return recipes, recipesElement, nil
}

type RecipeMap map[string]string
type RecipeElement map[string]Element
