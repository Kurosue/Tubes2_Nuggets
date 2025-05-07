package utils

import "fmt"

// "strings"

func DFSBotHelper(recipeMap RecipeMap, recipesEl RecipeElement, targetElement string, visited map[string]bool) []Element {
    if visited[targetElement] {
        return []Element{}
    }
    
    visited[targetElement] = true
    
    path := []Element{}
    if el, exists := recipesEl[targetElement]; exists {
        path = append(path, el)
    } else {
        path = append(path, Element{Name: targetElement})
    }
    
    for combination, result := range recipeMap {
        if result == targetElement {
            ing1, ing2 := DecomposeKey(combination)
			fmt.Printf("Combination: %s, Result: %s\n", combination, result)

			if( ing1 == "Time" || ing2 == "Time" ){ // Ini gimana yak , soalnya harus nemu lebih dari 100 items
				continue
			}
            
            if ing1 == "Air" || ing1 == "Water" || ing1 == "Earth" || ing1 == "Fire" {
                path = append(path, Element{Name: ing1})
            } else if !visited[ing1] {
                subPath := DFSBotHelper(recipeMap, recipesEl, ing1, visited)
                path = append(path, subPath...)
            }
            
            if ing2 == "Air" || ing2 == "Water" || ing2 == "Earth" || ing2 == "Fire" {
                path = append(path, Element{Name: ing2})
            } else if !visited[ing2] {
                subPath := DFSBotHelper(recipeMap, recipesEl, ing2, visited)
                path = append(path, subPath...)
            }
            
            break
        }
    }
    
    return path
}

func DFSBot(recipeMap RecipeMap, recipesEl RecipeElement, targetElement Element) []Element {
    visited := make(map[string]bool)
    return DFSBotHelper(recipeMap, recipesEl, targetElement.Name, visited)
}