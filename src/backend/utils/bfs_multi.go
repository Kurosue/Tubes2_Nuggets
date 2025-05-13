package utils

import (
	"strings"
	"sync"
)

// BFSNRecipesMultithreaded is a multithreaded version of BFSNRecipes
// It finds multiple recipe paths to create the target element
func BFSNRecipes(target string, Elements map[string]Element, maxRecipes int, visitedNode *int) []*Node {
	var resultRoots []*Node
	var resultMutex sync.Mutex
	queue := []*Node{{Result: target, Depth: 0}}
	var queueMutex sync.Mutex
	
	// Use a wait group to track all goroutines
	var wg sync.WaitGroup
	
	// Create a semaphore to limit the number of concurrent goroutines
	// This prevents creating too many goroutines which could overwhelm the system
	maxWorkers := 10 // Adjust this based on your system capabilities
	semaphore := make(chan struct{}, maxWorkers)
	
	// Flag to signal that we've reached max recipes
	done := make(chan struct{})
	isDone := false
	
	// Process function that will be run concurrently
	processNode := func(current *Node) {
		defer wg.Done()
		defer func() { <-semaphore }() // Release the semaphore when done
		
		// Check if we're done before processing
		select {
		case <-done:
			return
		default:
			// Continue processing
		}
		
		// Track visited nodes
		resultMutex.Lock()
		*visitedNode++
		resultMutex.Unlock()
		
		if BaseElement[current.Result] {
			resultMutex.Lock()
			resultRoots = append(resultRoots, current)
			if len(resultRoots) >= maxRecipes {
				if !isDone {
					close(done)
					isDone = true
				}
			}
			resultMutex.Unlock()
			return
		}
		
		elem, ok := Elements[current.Result]
		if !ok {
			return
		}
		
		tier := elem.Tier
		
		for _, recipe := range elem.Recipes {
			parts := strings.Split(recipe, "+")
			if len(parts) != 2 {
				continue
			}
			
			ing1 := strings.TrimSpace(parts[0][:len(parts[0])-1])
			ing2 := strings.TrimSpace(parts[1][1:])
			
			e1, ok1 := Elements[ing1]
			e2, ok2 := Elements[ing2]
			if !ok1 || !ok2 || e1.Tier >= tier || e2.Tier >= tier {
				continue
			}
			
			// Process ingredient search in parallel
			var leftMutex, rightMutex sync.Mutex
			var left, right *Node
			
			// Use wait group for ingredient searches
			var ingWg sync.WaitGroup
			ingWg.Add(2)
			
			go func() {
				defer ingWg.Done()
				subVisited := 0
				result := BFSShortestNode(ing1, Elements, &subVisited)
				leftMutex.Lock()
				left = result
				leftMutex.Unlock()
			}()
			
			go func() {
				defer ingWg.Done()
				subVisited := 0
				result := BFSShortestNode(ing2, Elements, &subVisited)
				rightMutex.Lock()
				right = result
				rightMutex.Unlock()
			}()
			
			// Wait for both ingredient searches to complete
			ingWg.Wait()
			
			newNode := &Node{
				Result:      current.Result,
				Depth:       current.Depth,
				Ingredient1: left,
				Ingredient2: right,
			}
			
			resultMutex.Lock()
			resultRoots = append(resultRoots, newNode)
			if len(resultRoots) >= maxRecipes {
				if !isDone {
					close(done)
					isDone = true
				}
				resultMutex.Unlock()
				return
			}
			resultMutex.Unlock()
		}
	}
	
	// Main processing loop
	for len(queue) > 0 {
		// Check if we're done
		select {
		case <-done:
			break
		default:
			// Continue processing
		}
		
		// Get next node from queue
		queueMutex.Lock()
		if len(queue) == 0 {
			queueMutex.Unlock()
			break
		}
		current := queue[0]
		queue = queue[1:]
		queueMutex.Unlock()
		
		// Add this task to our wait group
		wg.Add(1)
		
		// Acquire semaphore (blocks if maxWorkers goroutines are already running)
		semaphore <- struct{}{}
		
		// Process this node in a goroutine
		go processNode(current)
	}
	
	// Wait for all goroutines to complete
	wg.Wait()
	
	return resultRoots
}

// FlattenMultipleTrees flattens multiple recipe trees in parallel
func FlattenMultipleTrees(roots []*Node) [][]Message {
	result := make([][]Message, len(roots))
	var wg sync.WaitGroup
	
	for i, root := range roots {
		wg.Add(1)
		go func(index int, node *Node) {
			defer wg.Done()
			result[index] = FlattenTreeToMessages(node)
		}(i, root)
	}
	
	wg.Wait()
	return result
}