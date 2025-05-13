package utils

import (
	"strings"
	"sync"
)

func BFSNRecipes(target string, Elements map[string]Element, maxRecipes int, visitedNode *int) []*Node {
	var resultRoots []*Node
	var resultMutex sync.Mutex
	queue := []*Node{{Result: target, Depth: 0}}
	var queueMutex sync.Mutex
	
	var wg sync.WaitGroup
	
	maxWorkers := 10
	semaphore := make(chan struct{}, maxWorkers)
	
	done := make(chan struct{})
	isDone := false
	
	processNode := func(current *Node) {
		defer wg.Done()
		defer func() { <-semaphore }()
		
		select {
		case <-done:
			return
		default:
		}
		
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
			
			var leftMutex, rightMutex sync.Mutex
			var left, right *Node
			
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
	
	for len(queue) > 0 {
		select {
		case <-done:
			break
		default:
		}
		
		queueMutex.Lock()
		if len(queue) == 0 {
			queueMutex.Unlock()
			break
		}
		current := queue[0]
		queue = queue[1:]
		queueMutex.Unlock()
		
		wg.Add(1)
		
		semaphore <- struct{}{}
		
		go processNode(current)
	}
	
	wg.Wait()
	
	return resultRoots
}

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