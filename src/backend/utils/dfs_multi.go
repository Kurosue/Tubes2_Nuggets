package utils

import (
	"fmt"
	"maps"
	"strings"
	"sync"
	"time"
)

func DFSMultiple(recipeMap RecipeMap, recipesEl RecipeElement, targetElement Element, maxPaths int, resultChan chan Message) {
    var mutex sync.Mutex
    var wg sync.WaitGroup
    resultSet := make(map[string]bool)
    numWorkers := min(maxPaths * 2, 10)
    wg.Add(numWorkers)
    start := time.Now().UnixNano()
    
    for i := 0; i < numWorkers; i++ {
        go func(workerID int) {
            defer wg.Done()
            seed := workerID
            visited := make(map[string]bool)
            nodesVisited := 0
            recipePath := DFSHelperWithVariation(recipeMap, recipesEl, targetElement.Name, visited, &nodesVisited, 0, seed)
            pathSignature := generatePathSignature(recipePath)
            
            mutex.Lock()
            if len(recipePath) > 0 && len(resultSet) < maxPaths && !resultSet[pathSignature] {
                resultSet[pathSignature] = true
                resultChan <- Message{
                    RecipePath: recipePath,
                    NodesVisited: nodesVisited,
                    Duration: float32(time.Now().UnixNano() - start) / 1000000,
                }
            }
            mutex.Unlock()
        }(i)
    }
    wg.Wait()
}

func DFSHelperWithVariation(recipeMap RecipeMap, recipesEl RecipeElement, targetElement string, visited map[string]bool, nodesVisited *int, currentDepth int, seed int) []RecipePath {
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
    var combos [][2]string
    
    for _, recipe := range recipesEl[targetElement].Recipes {
        ing1Name := recipe[0]
        ing2Name := recipe[1]
        
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
        
        combos = append(combos, recipe)
    }
    
    // If no valid combinations found
    if len(combos) == 0 {
        return []RecipePath{}
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
    
    // Choose only one recipe based on the seed
    chosenIndex := seed % len(combos)
    ing1 := combos[chosenIndex][0]
    ing2 := combos[chosenIndex][1]
    
    var wg sync.WaitGroup
    wg.Add(2)
    ing1Channel := make(chan []RecipePath)
    ing2Channel := make(chan []RecipePath)
    
    go func() {
        defer wg.Done()
        visited1 := make(map[string]bool)
        maps.Copy(visited1, visited)
        subPath1 := DFSHelperWithVariation(recipeMap, recipesEl, ing1, visited1, nodesVisited, currentDepth + 1, seed)
        ing1Channel <- subPath1
    }()
    go func() {
        defer wg.Done()
        visited2 := make(map[string]bool)
        maps.Copy(visited2, visited)
        subPath2 := DFSHelperWithVariation(recipeMap, recipesEl, ing2, visited2, nodesVisited, currentDepth + 1, seed + 1) // Add variety
        ing2Channel <- subPath2
    }()
    subPath1 := <-ing1Channel
    subPath2 := <-ing2Channel
    wg.Wait()
    
    result = append(result, RecipePath{ Ingredient1: ing1, Ingredient2: ing2, Result: targetElement })
    result = append(result, subPath1...)
    result = append(result, subPath2...)
    return result
}

// generatePathSignature creates a unique signature for a recipe path
func generatePathSignature(messages []RecipePath) string {
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
