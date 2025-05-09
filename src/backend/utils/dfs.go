package utils

import (
    // "fmt"
    "sync"
)

func DFSHelper(recipeMap RecipeMap, recipesEl RecipeElement, targetElement string, visited map[string]bool) []Element {
    // Base case
    if targetElement == "Air" || targetElement == "Water" || targetElement == "Earth" || targetElement == "Fire" {
        if el, exists := recipesEl[targetElement]; exists {
            return []Element{el}
        } else {
            return []Element{{Name: targetElement}}
        }
    }
    
    if visited[targetElement] {
        if el, exists := recipesEl[targetElement]; exists {
            return []Element{el}
        } else {
            return []Element{{Name: targetElement}}
        }
    }
    
    visited[targetElement] = true
    
    // Start with the current element
    result := []Element{}
    if el, exists := recipesEl[targetElement]; exists {
        result = append(result, el)
    } else {
        result = append(result, Element{Name: targetElement})
    }
    
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

            if recipesEl[ing1].Tier >= recipesEl[targetElement].Tier || recipesEl[ing2].Tier >= recipesEl[targetElement].Tier {
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
            ing1Channel := make(chan []Element)
            ing2Channel := make(chan []Element)
            
            go func() {
                defer wg.Done()
                subPath1 := DFSHelper(recipeMap, recipesEl, ing1, visited1)
                ing1Channel <- subPath1
                
            }()
            
            go func() {
                defer wg.Done()
                subPath2 := DFSHelper(recipeMap, recipesEl, ing2, visited2)
                ing2Channel <- subPath2
                
            }()
            subPath1 := <-ing1Channel
            subPath2 := <-ing2Channel

            wg.Wait()
            result = append(result, subPath1...)
            result = append(result, subPath2...)
            
            return result
        }
    }
    
    return result
}


func DFS(recipeMap RecipeMap, recipesEl RecipeElement, targetElement Element) []Element {
    visited := make(map[string]bool)
    return DFSHelper(recipeMap, recipesEl, targetElement.Name, visited)
}