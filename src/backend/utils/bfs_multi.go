package utils

import (
	"sync"
	"time"
)

func BFSNRecipes(target string, Elements map[string]Element, maxRecipes int, resultChan chan Message) {
	resultCount := 0
	var resultMutex sync.Mutex
	queue := []*Node{{ Result: target }}
	var queueMutex sync.Mutex
	
	var wg sync.WaitGroup
	maxWorkers := 10
	semaphore := make(chan struct{}, maxWorkers)
	
	visitedNode := 0
	done := make(chan struct{})
	start := time.Now().UnixNano()
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
		visitedNode++
		resultMutex.Unlock()
		
		if BaseElement[current.Result] {
			resultMutex.Lock()
			resultCount++
			resultChan <- Message{
				RecipePath: FlattenTreeToRecipePaths(current),
				NodesVisited: visitedNode,
				Duration: float32(time.Now().UnixNano() - start) / 1000000,
			}
			if resultCount >= maxRecipes {
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
		
		for _, recipe := range elem.Recipes {
			ing1 := recipe[0]
			ing2 := recipe[1]
			
			e1, ok1 := Elements[ing1]
			e2, ok2 := Elements[ing2]
			if !ok1 || !ok2 || e1.Tier >= elem.Tier || e2.Tier >= elem.Tier {
				continue
			}
			
			var leftMutex, rightMutex sync.Mutex
			var left, right *Node
			
			var ingWg sync.WaitGroup
			ingWg.Add(2)
			go func() {
				defer ingWg.Done()
				subVisited := 0
				result := BFSShortestNodeImpl(ing1, Elements, &subVisited)
				visitedNode += subVisited
				leftMutex.Lock()
				left = result
				leftMutex.Unlock()
			}()
			go func() {
				defer ingWg.Done()
				subVisited := 0
				result := BFSShortestNodeImpl(ing2, Elements, &subVisited)
				visitedNode += subVisited
				rightMutex.Lock()
				right = result
				rightMutex.Unlock()
			}()
			ingWg.Wait()
			
			newNode := &Node{
				Result:      current.Result,
				Ingredient1: left,
				Ingredient2: right,
			}
			resultMutex.Lock()
			resultCount++
			resultChan <- Message{
				RecipePath: FlattenTreeToRecipePaths(newNode),
				NodesVisited: visitedNode,
				Duration: float32(time.Now().UnixNano() - start) / 1000000,
			}
			if resultCount >= maxRecipes {
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
}
