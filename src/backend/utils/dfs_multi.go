package utils

import (
    "fmt"
    "strings"
    "sync"
)

type MultiDFSResult struct {
    RecipePaths   [][]Message 
    NodesVisited  int       
    PathsFound    int         
}

func DFSMultiple(recipeMap RecipeMap, recipesEl RecipeElement, targetElement Element, maxPaths int) [][]Message {
    resultChan := make(chan []Message, maxPaths*2)
    
    var mutex sync.Mutex
    resultSet := make(map[string]bool)
    
    var wg sync.WaitGroup
    
    numWorkers := min(maxPaths*2, 10)
    wg.Add(numWorkers)
    
    for i := 0; i < numWorkers; i++ {
        go func(workerID int) {
            defer wg.Done()
            
            seed := workerID
            
            visited := make(map[string]bool)
            nodesVisited := 0
            messages := DFSHelperWithVariation(recipeMap, recipesEl, targetElement.Name, visited, &nodesVisited, 0, seed)
            
            pathSignature := generatePathSignature(messages)
            
            mutex.Lock()
            if len(messages) > 0 && !resultSet[pathSignature] {
                resultSet[pathSignature] = true
                resultChan <- messages
            }
            mutex.Unlock()
        }(i)
    }
    
    go func() {
        wg.Wait()
        close(resultChan)
    }()
    
    var allPaths [][]Message
    for path := range resultChan {
        if len(allPaths) < maxPaths {
            allPaths = append(allPaths, path)
        } else {
            break
        }
    }
    
    return allPaths
}

func DFSHelperWithVariation(recipeMap RecipeMap, recipesEl RecipeElement, targetElement string, 
                          visited map[string]bool, nodesVisited *int, currentDepth int, seed int) []Message {
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
    
    // Get all possible combinations that produce this element
    type recipePair struct {
        ing1, ing2 string
    }
    var combos []recipePair
    
    for _, combination := range recipesEl[targetElement].Recipes {
            ing1, ing2 := DecomposeKey(combination)
            
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
            
            combos = append(combos, recipePair{ing1, ing2})
    }
    
    if len(combos) > 1 {
        switch seed % 3 {
        case 0:
            // Original order
            break
        case 1:
            // Reverse order
            for i := 0; i < len(combos)/2; i++ {
                j := len(combos) - i - 1
                combos[i], combos[j] = combos[j], combos[i]
            }
        case 2:
            // Variety
            if len(combos) > 2 {
                mid := len(combos) / 2
                // Swap first and middle elements
                combos[0], combos[mid] = combos[mid], combos[0]
            }
        }
    }
    
    for _, combo := range combos {
        ing1 := combo.ing1
        ing2 := combo.ing2
        
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
            subPath1 := DFSHelperWithVariation(recipeMap, recipesEl, ing1, visited1, nodesVisited, currentDepth+1, seed)
            ing1Channel <- subPath1
        }()
        
        go func() {
            defer wg.Done()
            subPath2 := DFSHelperWithVariation(recipeMap, recipesEl, ing2, visited2, nodesVisited, currentDepth+1, seed+1) // Add variety
            ing2Channel <- subPath2
        }()
        
        subPath1 := <-ing1Channel
        subPath2 := <-ing2Channel
        
        wg.Wait()
        
        // Create a message for this combination
        result = append(result, Message{Ingredient1: ing1, Ingredient2: ing2, Result: targetElement, Depth: currentDepth})
        result = append(result, subPath1...)
        result = append(result, subPath2...)
        
        return result
    }
    
    return result
}

// generatePathSignature creates a unique signature for a recipe path
func generatePathSignature(messages []Message) string {
    var elements []string
    
    for _, msg := range messages {
        elements = append(elements, fmt.Sprintf("%s=%s+%s", 
            msg.Result, msg.Ingredient1, msg.Ingredient2))
    }
    
    return strings.Join(elements, "|")
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}