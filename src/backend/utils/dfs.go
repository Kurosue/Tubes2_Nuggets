package utils

import (
    // "fmt"
    "sync"
)

type Message struct {
    ingredient1 string
    ingredient2 string
    result      string
    depth       int
}

type DFSResult struct {
    Messages []Message
    NodesVisited int
}

func DFSHelper(recipeMap RecipeMap, recipesEl RecipeElement, targetElement string, visited map[string]bool, nodesVisited* int, currentDepth int) []Message {
    // Base case
    *nodesVisited++
    if targetElement == "Air" || targetElement == "Water" || targetElement == "Earth" || targetElement == "Fire" {
        if _, exists := recipesEl[targetElement]; exists {
            return []Message{{ingredient1: "", ingredient2: "", result: targetElement, depth: currentDepth}}
        }
    }
    
    if visited[targetElement] {
        if _, exists := recipesEl[targetElement]; exists {
            return []Message{{ingredient1: "", ingredient2: "", result: targetElement, depth: currentDepth}}
        }
    }
    
    visited[targetElement] = true
    
    result := []Message{}
    
    // Find recipe to create this element
    for combination, product := range recipeMap {
        if product == targetElement {
            ing1, ing2 := DecomposeKey(combination)
            // fmt.Printf("Combination: %s, Result: %s\n", combination, product)
            
            // Skip yang ga ada di resepnya
            if ing1 == "Time" || ing2 == "Time" {
                continue
            }
            if _, exists := recipesEl[ing1]; !exists {
                continue
            }
            if _, exists := recipesEl[ing2]; !exists {
                continue
            }

            if recipesEl[ing1].Tier > recipesEl[targetElement].Tier || recipesEl[ing2].Tier > recipesEl[targetElement].Tier {
                continue
            }

            visited1 := make(map[string]bool)
            for k, v := range visited {
                visited1[k] = v
            }
            
            visited2 := make(map[string]bool)
            for k, v := range visited {
                visited2[k] = v
            }

            var wg sync.WaitGroup
            wg.Add(2)
            ing1Channel := make(chan []Message)
            ing2Channel := make(chan []Message)
            
            go func() {
                defer wg.Done()
                subPath1 := DFSHelper(recipeMap, recipesEl, ing1, visited1, nodesVisited, currentDepth+1)
                ing1Channel <- subPath1
                
            }()
            
            go func() {
                defer wg.Done()
                subPath2 := DFSHelper(recipeMap, recipesEl, ing2, visited2, nodesVisited, currentDepth+1)
                ing2Channel <- subPath2
                
            }()
            subPath1 := <-ing1Channel
            subPath2 := <-ing2Channel

            result = append(result, Message{ingredient1: ing1, ingredient2: ing2, result: targetElement, depth: currentDepth})

            wg.Wait()
            result = append(result, subPath1...)
            result = append(result, subPath2...)
            
            return result
        }
    }
    
    return result
}

func DFS(recipeMap RecipeMap, recipesEl RecipeElement, targetElement Element) ([]Message, int) {
    visited := make(map[string]bool)
    nodesVisited := 0
    if _, exists := recipesEl[targetElement.Name]; !exists {
        return []Message{}, 0
    }
    messages := DFSHelper(recipeMap, recipesEl, targetElement.Name, visited, &nodesVisited, 0)
    return messages, nodesVisited
}