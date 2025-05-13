package utils

import (
	"maps"
	"sync"
	"time"
)

type QueuePathItem struct {
	element      string
	tier         int
	recipePath   []RecipePath
	pendingElems map[string]bool
}

func BFSP(startElement string, recipeMap RecipeMap, elements RecipeElement, maxResult int, resultChan chan Message) {
	var visitedMu sync.RWMutex
    var nodeVisitedMu sync.Mutex
	var resultCountMu sync.Mutex

	visited := make(map[string]bool)
	visited[startElement] = true

	initialPending := make(map[string]bool)
	initialPending[startElement] = true

	queue := []QueuePathItem{{
		element:      startElement,
		tier:         elements[startElement].Tier,
		recipePath:   []RecipePath{},
		pendingElems: initialPending,
	}}

    nodeVisited := 0
	resultCount := 0
    start := time.Now().UnixNano()

	const maxWorkers = 4
	sem := make(chan struct{}, maxWorkers)

	for len(queue) > 0 && resultCount < maxResult {
		currentLevel := queue
		queue = nil

		var queueMu sync.Mutex
		var wg sync.WaitGroup

		for _, current := range currentLevel {

            nodeVisitedMu.Lock()
            nodeVisited++
            nodeVisitedMu.Unlock()

			if !current.pendingElems[current.element] {
				continue
			}

			wg.Add(1)
			sem <- struct{}{}

			go func(item QueuePathItem) {
				defer func() {
					<-sem
					wg.Done()
				}()

				pendingCopy := make(map[string]bool)
				maps.Copy(pendingCopy, item.pendingElems)

				delete(pendingCopy, item.element)

				if BaseElement[item.element] {
					newPath := append([]RecipePath{}, item.recipePath...)
					newPath = append(newPath, RecipePath{
						Ingredient1: "",
						Ingredient2: "",
						Result: item.element,
					})

					if len(pendingCopy) == 0 {
						resultCountMu.Lock()
						shouldResult := resultCount < maxResult
						resultCount++
						resultCountMu.Unlock()
						if(shouldResult) {
							resultChan <- Message{
								RecipePath: newPath,
								NodesVisited: nodeVisited,
								Duration: float32(time.Now().UnixNano() - start) / 1000000,
							}
						}
					} else {
						for nextElem := range pendingCopy {
							queueMu.Lock()
							queue = append(queue, QueuePathItem{
								element:      nextElem,
								tier:         elements[nextElem].Tier,
								recipePath:   newPath,
								pendingElems: pendingCopy,
							})
							queueMu.Unlock()
							break
						}
					}
					return
				}

				for _, recipe := range elements[item.element].Recipes {
					ing1, ok1 := elements[recipe[0]]
					ing2, ok2 := elements[recipe[1]]
					if !ok1 || !ok2 {
						continue
					}

					if ing1.Tier < item.tier && ing2.Tier < item.tier {
						msg := RecipePath{
							Ingredient1: ing1.Name,
							Ingredient2: ing2.Name,
							Result:      item.element,
						}

						newPath := append([]RecipePath{}, item.recipePath...)
						newPath = append(newPath, msg)

						newPending := make(map[string]bool)
						maps.Copy(newPending, pendingCopy)

						visitedMu.Lock()
						if !BaseElement[ing1.Name] {
							newPending[ing1.Name] = true
							visited[ing1.Name] = true
						}
						if !BaseElement[ing2.Name] {
							newPending[ing2.Name] = true
							visited[ing2.Name] = true
						}
						visitedMu.Unlock()

						if BaseElement[ing1.Name] {
							baseMsg := RecipePath{
								Ingredient1: "",
								Ingredient2: "",
								Result: ing1.Name,
							}
							newPath = append(newPath, baseMsg)
						}
						if BaseElement[ing2.Name] {
							baseMsg := RecipePath{
								Ingredient1: "",
								Ingredient2: "",
								Result: ing2.Name,
							}
							newPath = append(newPath, baseMsg)
						}

						if len(newPending) == 0 {
							resultCountMu.Lock()
							shouldResult := resultCount < maxResult
							resultCount++
							resultCountMu.Unlock()
							if(shouldResult) {
								resultChan <- Message{
									RecipePath: newPath,
									NodesVisited: nodeVisited,
									Duration: float32(time.Now().UnixNano() - start) / 1000000,
								}
							}
						} else {
							for nextElem := range newPending {
								queueMu.Lock()
								queue = append(queue, QueuePathItem{
									element:      nextElem,
									tier:         elements[nextElem].Tier,
									recipePath:   newPath,
									pendingElems: newPending,
								})
								queueMu.Unlock()
								break
							}
						}
					}
				}
			}(current)
		}
		wg.Wait()
	}
}
