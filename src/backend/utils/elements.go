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
}

func PairKey(a, b string) string {
    if a < b {
        return a + "|" + b
    }
    return b + "|" + a
}

func LoadRecipes(path string) (RecipeMap, error) {
    var elems []Element
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }
    if err := json.Unmarshal(data, &elems); err != nil {
        return nil, err
    }

    recipes := make(RecipeMap)
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
    }
    return recipes, nil
}

type RecipeMap map[string]string