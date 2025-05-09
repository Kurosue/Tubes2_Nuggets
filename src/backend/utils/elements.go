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
    // Nanti dimasukin pas parsing json nya
    ParsedRecipes []Recipe
}

type Recipe struct {
    First  string
    Second string
}

type ElementMap map[string]*Element

func LoadRecipe(path string) (ElementMap, error) {
    var elems []Element
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    if err := json.Unmarshal(data, &elems); err != nil {
        return nil, err
    }

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
    }
    
    return elementMap, nil
}