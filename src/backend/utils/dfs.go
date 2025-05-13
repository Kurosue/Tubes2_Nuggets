package utils

import (
	// "fmt"
	"maps"
	"sync"
	"time"
)

func DFSHelper(recipeMap RecipeMap, recipesEl RecipeElement, targetElement string, visited map[string]bool, nodesVisited* int, currentDepth int) []RecipePath {
    // Base case
    *nodesVisited++
    target, targetExists := recipesEl[targetElement]
    if !targetExists {
        return []RecipePath{}
    }
    if visited[targetElement] || BaseElement[targetElement] {
        return []RecipePath{{ Ingredient1: "", Ingredient2: "", Result: targetElement }}
    }
    
    visited[targetElement] = true
    result := []RecipePath{}
    
    // Find recipe to create this element
    for _, recipe := range target.Recipes {
        ing1Name := recipe[0]
        ing2Name := recipe[1]
        // fmt.Printf("Combination: %s, Result: %s\n", combination, targetElement)
        // fmt.Printf("Ingredients: %s, %s\n", ing1, ing2)
        
        // Skip yang ga ada di resepnya
        if ing1Name == "Time" || ing2Name == "Time" {
            continue
        }
        ing1, ing1Exists := recipesEl[ing1Name]
        ing2, ing2Exists := recipesEl[ing2Name]
        if !ing1Exists || !ing2Exists {
            continue
        }
        if ing1.Tier >= target.Tier || ing2.Tier >= target.Tier {
            continue
        }

        var wg sync.WaitGroup
        wg.Add(2)
        ing1Channel := make(chan []RecipePath)
        ing2Channel := make(chan []RecipePath)
        go func() {
            defer wg.Done()
            visited1 := make(map[string]bool)
            maps.Copy(visited1, visited)
            subPath1 := DFSHelper(recipeMap, recipesEl, ing1Name, visited1, nodesVisited, currentDepth + 1)
            ing1Channel <- subPath1
            
        }()
        go func() {
            defer wg.Done()
            visited2 := make(map[string]bool)
            maps.Copy(visited2, visited)
            subPath2 := DFSHelper(recipeMap, recipesEl, ing2Name, visited2, nodesVisited, currentDepth + 1)
            ing2Channel <- subPath2
            
        }()
        wg.Wait()
        subPath1 := <-ing1Channel
        subPath2 := <-ing2Channel

        result = append(result, RecipePath{ Ingredient1: ing1Name, Ingredient2: ing2Name, Result: targetElement })
        result = append(result, subPath1...)
        result = append(result, subPath2...)
        return result
    }
    return result
}

func DFS(recipeMap RecipeMap, recipesEl RecipeElement, targetElement Element, resultChan chan Message) {
    visited := make(map[string]bool)
    nodesVisited := 0
    start := time.Now().UnixNano()
    if _, exists := recipesEl[targetElement.Name]; !exists {
        resultChan <- Message{
            RecipePath: []RecipePath{},
            NodesVisited: 0,
            Duration: float32(time.Now().UnixNano() - start) / 1000000,
        }
        return
    }
    recipePath := DFSHelper(recipeMap, recipesEl, targetElement.Name, visited, &nodesVisited, 0)
    resultChan <- Message{
        RecipePath: recipePath,
        NodesVisited: nodesVisited,
        Duration: float32(time.Now().UnixNano() - start) / 1000000,
    }
}
