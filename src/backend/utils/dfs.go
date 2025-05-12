package utils

import (
    // "fmt"
    "sync"
)

type DFSResult struct {
    Messages []Message
    NodesVisited int
}

func DFSHelper(recipeMap RecipeMap, recipesEl RecipeElement, targetElement string, visited map[string]bool, nodesVisited* int, currentDepth int) []Message {
    // Base case
    *nodesVisited++
    if targetElement == "Air" || targetElement == "Water" || targetElement == "Earth" || targetElement == "Fire" {
        if _, exists := recipesEl[targetElement]; exists {
            return []Message{{Ingredient1: "", Ingredient2: "", Result: targetElement, Depth: currentDepth}}
        }
    }
    
    if visited[targetElement] {
        if _, exists := recipesEl[targetElement]; exists {
            return []Message{{Ingredient1: "", Ingredient2: "", Result: targetElement, Depth: currentDepth}}
        }
    }
    
    visited[targetElement] = true
    
    result := []Message{}
    
    // Find recipe to create this element
    for _, combination := range recipesEl[targetElement].Recipes {
            ing1, ing2 := DecomposeKeyWithPlus(combination)
            // fmt.Printf("Combination: %s, Result: %s\n", combination, targetElement)
            // fmt.Printf("Ingredients: %s, %s\n", ing1, ing2)
            
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

            result = append(result, Message{Ingredient1: ing1, Ingredient2: ing2, Result: targetElement, Depth: currentDepth})
        
            wg.Wait()
            result = append(result, subPath1...)
            result = append(result, subPath2...)
            
            return result
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
