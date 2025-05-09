package utils

import (
	"sync"
)

// DFSMultiHelper finds multiple ways to create the target element
func DFSMultiHelper(
	recipeMap RecipeMap, 
	recipesEl RecipeElement, 
	targetElement string, 
	visited map[string]bool, 
	resultChan chan []Element,
	depth int,
	maxDepth int,
) {
	// Limit recursion depth to avoid excessive branching
	if depth > maxDepth {
		return
	}

	// Base case - primitive elements
	if targetElement == "Air" || targetElement == "Water" || targetElement == "Earth" || targetElement == "Fire" {
		if el, exists := recipesEl[targetElement]; exists {
			resultChan <- []Element{el}
		} else {
			resultChan <- []Element{{Name: targetElement}}
		}
		return
	}
	
	// Check if we've already visited this element to prevent cycles
	if visited[targetElement] {
		if el, exists := recipesEl[targetElement]; exists {
			resultChan <- []Element{el}
		} else {
			resultChan <- []Element{{Name: targetElement}}
		}
		return
	}
	
	// Make a local copy of visited map for this branch
	localVisited := make(map[string]bool)
	for k, v := range visited {
		localVisited[k] = v
	}
	localVisited[targetElement] = true
	
	// Current element
	var currentElement Element
	if el, exists := recipesEl[targetElement]; exists {
		currentElement = el
	} else {
		currentElement = Element{Name: targetElement}
	}
	
	// Track found recipes
	recipesFound := 0
	
	// Find different recipes to create this element
	for combination, product := range recipeMap {
		if product == targetElement {
			ing1, ing2 := DecomposeKey(combination)
			
			// Skip Time-based recipes and elements that the ingredients are not in the recipes
			if ing1 == "Time" || ing2 == "Time" {
				continue
			}
			if _, exists := recipesEl[ing1]; !exists {
				continue
			}
			if _, exists := recipesEl[ing2]; !exists {
				continue
			}

			// Skip recipes where ingredients have higher or equal tier than the target
			targetTier := recipesEl[targetElement].Tier
			if ing1Tier, ok := recipesEl[ing1]; ok && ing1Tier.Tier >= targetTier {
				continue
			}
			if ing2Tier, ok := recipesEl[ing2]; ok && ing2Tier.Tier >= targetTier {
				continue
			}

			// Create separate visited maps for this recipe
			branchVisited := make(map[string]bool)
			for k, v := range localVisited {
				branchVisited[k] = v
			}
			
			// Process ingredients 
			ing1ResultChan := make(chan []Element, 10)
			ing2ResultChan := make(chan []Element, 10)
			
			var wg sync.WaitGroup
			wg.Add(2)
			
			// Process first ingredient
			go func() {
				defer wg.Done()
				DFSMultiHelper(recipeMap, recipesEl, ing1, branchVisited, ing1ResultChan, depth+1, maxDepth)
				close(ing1ResultChan)
			}()
			
			// Process second ingredient
			go func() {
				defer wg.Done()
				DFSMultiHelper(recipeMap, recipesEl, ing2, branchVisited, ing2ResultChan, depth+1, maxDepth)
				close(ing2ResultChan)
			}()
			
			// Wait for all processing to complete
			wg.Wait()
			
			// Collect results
			var ing1Results [][]Element
			for result := range ing1ResultChan {
				ing1Results = append(ing1Results, result)
			}
			
			var ing2Results [][]Element
			for result := range ing2ResultChan {
				ing2Results = append(ing2Results, result)
			}
			
			// Generate all combinations of ingredient results
			for _, ing1Result := range ing1Results {
				for _, ing2Result := range ing2Results {
					// Build full path
					result := []Element{currentElement}
					result = append(result, ing1Result...)
					result = append(result, ing2Result...)
					
					// Send result
					resultChan <- result
					recipesFound++
				}
			}
		}
	}
	
	// If no recipes found, just return the element itself
	if recipesFound == 0 {
		resultChan <- []Element{currentElement}
	}
}

// DFSMulti finds multiple ways to create the target element
func DFSMulti(recipeMap RecipeMap, recipesEl RecipeElement, targetElement Element, maxRecipes int) [][]Element {
	// Set up result collection
	resultChan := make(chan []Element, maxRecipes*10) // Buffer to prevent blocking
	visited := make(map[string]bool)
	
	// Start a goroutine to populate the results
	go func() {
		DFSMultiHelper(recipeMap, recipesEl, targetElement.Name, visited, resultChan, 0, 50)
		close(resultChan)
	}()
	
	// Collect unique results
	results := [][]Element{}
	seenPaths := make(map[string]bool)
	
	// Process results as they come in
	for result := range resultChan {
		// Create a simple hash of this path to avoid duplicates
		pathHash := ""
		for _, el := range result {
			pathHash += el.Name + ","
		}
		
		// Only add unique paths
		if !seenPaths[pathHash] {
			seenPaths[pathHash] = true
			results = append(results, result)
			
			// Stop once we have enough results
			if len(results) >= maxRecipes {
				break
			}
		}
	}
	
	return results
}