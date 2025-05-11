package utils

import (
	"strings"
	"sync"
)

var BaseElement = map[string]bool{
	"Fire":  true,
	"Water": true,
	"Earth": true,
	"Air":   true,
}

type Path struct {
	Messages []Message
}

type QueuePathItem struct {
	element      string
	tier         int
	depth        int
	path         Path
	pendingElems map[string]bool
}

func BFSP(start string, recipeMap RecipeMap, elements RecipeElement) (res [][]Message) {
	var m sync.Mutex
	var visitedMu sync.RWMutex

	visited := make(map[string]bool)

	initialPending := make(map[string]bool)
	initialPending[start] = true

	queue := []QueuePathItem{{
		element:      start,
		tier:         elements[start].Tier,
		depth:        0,
		path:         Path{},
		pendingElems: initialPending,
	}}

	visited[start] = true

	const maxWorkers = 4
	sem := make(chan struct{}, maxWorkers)

	for len(queue) > 0 {
		currentLevel := queue
		queue = nil

		var queueMu sync.Mutex
		var wg sync.WaitGroup

		for _, current := range currentLevel {
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
				for k, v := range item.pendingElems {
					pendingCopy[k] = v
				}

				delete(pendingCopy, item.element)

				if BaseElement[item.element] {
					baseElementMsg := Message{
						Result: item.element,
						Depth:  item.depth,
					}

					newPath := append([]Message{}, item.path.Messages...)
					newPath = append(newPath, baseElementMsg)

					if len(pendingCopy) == 0 {
						m.Lock()
						res = append(res, newPath)
						m.Unlock()
					} else {
						for nextElem := range pendingCopy {
							queueMu.Lock()
							queue = append(queue, QueuePathItem{
								element:      nextElem,
								tier:         elements[nextElem].Tier,
								depth:        item.depth,
								path:         Path{Messages: newPath},
								pendingElems: pendingCopy,
							})
							queueMu.Unlock()
							break
						}
					}
					return
				}

				for _, recipe := range elements[item.element].Recipes {
					parts := strings.Split(recipe, "+")
					if len(parts) != 2 {
						continue
					}

					first := strings.TrimSpace(parts[0])
					second := strings.TrimSpace(parts[1])

					_, ok1 := elements[first]
					_, ok2 := elements[second]
					if !ok1 || !ok2 {
						continue
					}

					firstTier := elements[first].Tier
					secondTier := elements[second].Tier

					if firstTier < item.tier && secondTier < item.tier {
						msg := Message{
							Ingredient1: first,
							Ingredient2: second,
							Result:      item.element,
							Depth:       item.depth,
						}

						newPath := append([]Message{}, item.path.Messages...)
						newPath = append(newPath, msg)

						newPending := make(map[string]bool)
						for k, v := range pendingCopy {
							newPending[k] = v
						}

						if !BaseElement[first] {
							newPending[first] = true
							visitedMu.Lock()
							visited[first] = true
							visitedMu.Unlock()
						}

						if !BaseElement[second] {
							newPending[second] = true
							visitedMu.Lock()
							visited[second] = true
							visitedMu.Unlock()
						}

						if BaseElement[first] && BaseElement[second] {
							baseMsg1 := Message{
								Result: first,
								Depth:  item.depth + 1,
							}
							baseMsg2 := Message{
								Result: second,
								Depth:  item.depth + 1,
							}

							finalPath := append([]Message{}, newPath...)
							finalPath = append(finalPath, baseMsg1, baseMsg2)

							if len(newPending) == 0 {
								m.Lock()
								res = append(res, finalPath)
								m.Unlock()
							} else {
								for nextElem := range newPending {
									queueMu.Lock()
									queue = append(queue, QueuePathItem{
										element:      nextElem,
										tier:         elements[nextElem].Tier,
										depth:        item.depth + 1,
										path:         Path{Messages: finalPath},
										pendingElems: newPending,
									})
									queueMu.Unlock()
									break
								}
							}
							continue
						}

						if BaseElement[first] {
							baseMsg := Message{
								Result: first,
								Depth:  item.depth + 1,
							}
							newPath = append(newPath, baseMsg)
						}

						if BaseElement[second] {
							baseMsg := Message{
								Result: second,
								Depth:  item.depth + 1,
							}
							newPath = append(newPath, baseMsg)
						}

						if len(newPending) > 0 {
							for nextElem := range newPending {
								queueMu.Lock()
								queue = append(queue, QueuePathItem{
									element:      nextElem,
									tier:         elements[nextElem].Tier,
									depth:        item.depth + 1,
									path:         Path{Messages: newPath},
									pendingElems: newPending,
								})
								queueMu.Unlock()
								break
							}
						} else {
							m.Lock()
							res = append(res, newPath)
							m.Unlock()
						}
					}
				}
			}(current)
		}

		wg.Wait()
	}

	return res
}
