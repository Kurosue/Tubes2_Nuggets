package utils

import (
    // "strings"
    "fmt"
)

func BFS(start string, recipes RecipeMap) map[string]string {
    // visited := make(map[string]bool)
    // queue := []string{start}
    path := make(map[string]string)
    path[start] = ""

    fmt.Print(recipes[start])
    return path
}