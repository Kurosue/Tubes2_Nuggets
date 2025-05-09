package utils

import (
    "sync"
    "strings"
)

func BFSP(start string, recipeMap RecipeMap, elements RecipeElement) (res []Message) {
    var resMu sync.Mutex
    var visitedMu sync.RWMutex

    visited := make(map[string]struct{
        tier int
        processed bool
    })
    
    // INIT Variabel
    queue := []QueueItem{{element: start, tier: elements[start].Tier, depth: 0}}
    visited[start] = struct{tier int; processed bool}{elements[start].Tier, false}
    
    // Pool thread cuma 4 sesuai core :D
    const maxWorkers = 4
    sem := make(chan struct{}, maxWorkers)
    
    for len(queue) > 0 {
        currentLevel := queue
        queue = nil
        
        var queueMu sync.Mutex
        var wg sync.WaitGroup
        
        for _, current := range currentLevel {
            // Skip kalau udah diproses
            visitedMu.RLock()
            visitedInfo, exists := visited[current.element]
            alreadyProcessed := exists && visitedInfo.processed
            visitedMu.RUnlock()
            
            if alreadyProcessed {
                continue
            }
            
            // Mark as processed (with write lock)
            visitedMu.Lock()
            visitedInfo = visited[current.element]
            visitedInfo.processed = true
            visited[current.element] = visitedInfo
            visitedMu.Unlock()
            
            wg.Add(1)
            sem <- struct{}{} // Acquire semaphore
            
            go func(item QueueItem) {
                defer func() {
                    <-sem // Release semaphore
                    wg.Done()
                }()
                
                var localres []Message
                var nextItems []QueueItem
                
                // Process recipes for this element
                for _, recipe := range elements[item.element].Recipes {
                    parts := strings.Split(recipe, "+")
                    if len(parts) != 2 {
                        continue
                    }
                    
                    first := strings.TrimSpace(parts[0])
                    second := strings.TrimSpace(parts[1])
                    
                    //  Kalau kedua elemen adalah base element
                    if BaseElement[first] && BaseElement[second] {
                        localres = append(localres, Message{
                            Ingredient1: first,
                            Ingredient2: second,
                            Depth:       item.depth,
                        })
                        continue
                    }
                    
                    // Skip kalalu elemen tidak ada di data json
                    _, ok1 := elements[first]
                    _, ok2 := elements[second]
                    if !ok1 || !ok2 {
                        continue
                    }
                    
                    firstTier := elements[first].Tier
                    secondTier := elements[second].Tier
                    
                    // cek tier child harus lebih rendah dari paraent
                    if firstTier < item.tier && secondTier < item.tier {
                        localres = append(localres, Message{
                            Ingredient1: first,
                            Ingredient2: second,
                            Depth:       item.depth,
                        })
                        
                        visitedMu.RLock()
                        _, firstVisited := visited[first]
                        visitedMu.RUnlock()
                        
                        if !firstVisited {
                            visitedMu.Lock()
                            visited[first] = struct{tier int; processed bool}{firstTier, false}
                            visitedMu.Unlock()
                            
                            if !BaseElement[first] {
                                nextItems = append(nextItems, QueueItem{
                                    element: first,
                                    tier:    firstTier,
                                    depth:   item.depth + 1,
                                })
                            }
                        }
                        
                        // Add second ingredient to queue if needed
                        visitedMu.RLock()
                        _, secondVisited := visited[second]
                        visitedMu.RUnlock()
                        
                        if !secondVisited {
                            visitedMu.Lock()
                            visited[second] = struct{tier int; processed bool}{secondTier, false}
                            visitedMu.Unlock()
                            
                            if !BaseElement[second] {
                                nextItems = append(nextItems, QueueItem{
                                    element: second,
                                    tier:    secondTier,
                                    depth:   item.depth + 1,
                                })
                            }
                        }
                    }
                }
                
                if len(localres) > 0 {
                    resMu.Lock()
                    res = append(res, localres...)
                    resMu.Unlock()
                }

                if len(nextItems) > 0 {
                    queueMu.Lock()
                    queue = append(queue, nextItems...)
                    queueMu.Unlock()
                }
            }(current)
        }
        
        wg.Wait()
    }
    
    return res
}