package utils

import (
	"encoding/json"
	"io/ioutil"
	"strings"
)

// Struct element sesuai json
type Element struct {
    Name     string   `json:"name"`
    Recipes  []string `json:"recipes"`
    Image    string   `json:"image"`
    PageURL  string   `json:"page_url"`
    Tier     int      `json:"tier"`
}

func PairKey(a, b string) string {
    if a < b {
        return a + "|" + b
    }
    return b + "|" + a
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

func LoadRecipes(path string) (RecipeMap, RecipeElement, error) {
    var elems []Element
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, nil, err
    }
    if err := json.Unmarshal(data, &elems); err != nil {
        return nil, nil, err
    }

    recipes := make(RecipeMap)
    recipesElement := make(RecipeElement)
    for _, e := range elems {
        for _, r := range e.Recipes {
            parts := strings.Split(r, "+")
            if len(parts) != 2 {
                continue
            }
            a := strings.TrimSpace(parts[0])
            b := strings.TrimSpace(parts[1])
            recipes[PairKey(a, b)] = e.Name
        }
        recipesElement[e.Name] = e
    }
    return recipes, recipesElement, nil
}

type RecipeMap map[string]string
type RecipeElement map[string]Element