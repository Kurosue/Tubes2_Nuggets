package utils

import (
    "strings"
)

func BFSShortestPath(start string, recipeMap RecipeMap, elements RecipeElement) (resultPath []Message) {
    
    type PathStep struct {
        element string
        tier    int
        path    []Message
    }
    
    visited := make(map[string]bool)
    queue := []PathStep{{
        element: start,
        tier:    elements[start].Tier,
        path:    []Message{},
    }}
    
    visited[start] = true
    
    for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:]
        
        // Check if we reached a base element
        if BaseElement[current.element] {
            // We found a path to a base element, which is our goal
            return current.path
        }
        
        for _, recipe := range elements[current.element].Recipes {
            parts := strings.Split(recipe, "+")
            if len(parts) != 2 {
                continue
            }
            
            first := strings.TrimSpace(parts[0])
            second := strings.TrimSpace(parts[1])
            
            // Skip if elements don't exist
            _, ok1 := elements[first]
            _, ok2 := elements[second]
            if !ok1 || !ok2 {
                continue
            }
            
            // If both ingredients are base elements, we found the shortest solution
            if BaseElement[first] && BaseElement[second] {
                // Create the final recipe message
                finalRecipe := Message{
                    Ingredient1: first, 
                    Ingredient2: second,
                    Depth:       len(current.path),
                }
                
                // Return the complete path including the final recipe
                return append(current.path, finalRecipe)
            }
            
            // Otherwise, continue BFS to both ingredients
            firstTier := elements[first].Tier
            secondTier := elements[second].Tier
            
            // Respect tier constraints
            if firstTier < current.tier && secondTier < current.tier {
                // Create a recipe message for this step
                newRecipe := Message{
                    Ingredient1: first,
                    Ingredient2: second,
                    Depth:       len(current.path),
                }
                
                // Create new paths for both ingredients
                newPath := append(append([]Message{}, current.path...), newRecipe)
                
                if !visited[first] {
                    visited[first] = true
                    queue = append(queue, PathStep{
                        element: first,
                        tier:    firstTier,
                        path:    newPath,
                    })
                }
                
                if !visited[second] {
                    visited[second] = true
                    queue = append(queue, PathStep{
                        element: second,
                        tier:    secondTier,
                        path:    newPath,
                    })
                }
            }
        }
    }
    
    // If we get here, no path was found
    return resultPath
}