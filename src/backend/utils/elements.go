package utils

import (
	"encoding/json"
	"io/ioutil"
	"strings"
)

var BaseElement = map[string]bool{
	"Fire":  true,
	"Water": true,
	"Earth": true,
	"Air":   true,
}

// Nanti jadi return value ke Front end
type RecipePath struct {
    Ingredient1 string `json:"ingredient2"`
    Ingredient2 string `json:"ingredient1"`
    Result      string `json:"result"`
}
type Message struct {
    RecipePath []RecipePath `json:"recipePath"`
    NodesVisited int        `json:"nodesVisited"`
    Duration float32        `json:"duration"`
}

// Struct element sesuai json
type Element struct {
    Name     string      `json:"name"`
    Recipes  [][2]string `json:"recipes"`
    Image    string      `json:"image"`
    PageURL  string      `json:"page_url"`
    Tier    int          `json:"tier"`
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

func DecomposeKeyWithPlus(key string) (string, string) {
    parts := strings.Split(key, "+")
    if len(parts) != 2 {
        return "", ""
    }
    return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
}

func StringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func LoadRecipes(path string) (RecipeMap, RecipeElement, error) {
    var elements []Element
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, nil, err
    }
    if err := json.Unmarshal(data, &elements); err != nil {
        return nil, nil, err
    }

    recipes := make(RecipeMap)
    recipesElement := make(RecipeElement)
    for _, element := range elements {
        for _, recipe := range element.Recipes {
            a := strings.TrimSpace(recipe[0])
            b := strings.TrimSpace(recipe[1])
            recipes[PairKey(a, b)] = element.Name
        }
        recipesElement[element.Name] = element
    }
    return recipes, recipesElement, nil
}

type RecipeMap map[string]string
type RecipeElement map[string]Element
