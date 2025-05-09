package utils

import (
    "strings"
)

var BaseElement = map[string]bool{
    "Fire":   true,
    "Water":  true,
    "Earth":  true,
    "Air":    true,
}

type QueueItem struct {
    element string
    tier    int
    depth   int
}

func BFS(start string, recipeMap RecipeMap, elements RecipeElement) (res []Message) {
    // Init queue buat track elemen yang mau diolah
    queue := []QueueItem{{element: start, tier: elements[start].Tier, depth: 0}}
    
    // Struct variabel untuk node yang udah dijelajahi
    visited := make(map[string]struct{
        tier int
        processed bool
    })
    
    visited[start] = struct{tier int; processed bool}{elements[start].Tier, false}
    
    // LOOP BFS
    for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:]
        
        // Skip kalo udah pernah dijelajahi
        if v, exists := visited[current.element]; exists && v.processed {
            continue
        }
        
        visitedInfo := visited[current.element]
        visitedInfo.processed = true
        visited[current.element] = visitedInfo
        
        // Iterasi ke semua resep atau anak 
        for _, recipe := range elements[current.element].Recipes {
            parts := strings.Split(recipe, "+")
            
            first := strings.TrimSpace(parts[0])
            second := strings.TrimSpace(parts[1])
            
            // Kalo kedua elemen udah base elemen, langsung masukkan ke hasil trus continue
            if BaseElement[first] && BaseElement[second] {
                res = append(res, Message{
                    Ingredient1: first,
                    Ingredient2: second,
                    Depth:       current.depth,
                })
                continue
            }
            
            // Skil kalo elemen nya gaada di data json
            _, cek1 := elements[first]
            _, cek2 := elements[second]
            if !cek1 || !cek2 {
                continue
            }
            
            firstTier := elements[first].Tier
            secondTier := elements[second].Tier
            
            // Make sure kalau child atau resep dari sebuah elemen tuh lebih rendah tiernya
            // SESUAI SPEK BTW
            if firstTier < current.tier && secondTier < current.tier {
                // nambah ke list hasil
                res = append(res, Message{
                    Ingredient1: first,
                    Ingredient2: second,
                    Depth:       current.depth,
                })
                
                // Cek elemen pertama dari resep kalau udah ada dan kalau base elemen
                if _, ok := visited[first]; !ok {
                    visited[first] = struct{tier int; processed bool}{firstTier, false}
                    if !BaseElement[first] {
                        queue = append(queue, QueueItem{
                            element: first,
                            tier:    firstTier,
                            depth:   current.depth + 1,
                        })
                    }
                }
                
                // sama tapi elemen kedua dari resep
                if _, ok := visited[second]; !ok {
                    visited[second] = struct{tier int; processed bool}{secondTier, false}
                    if !BaseElement[second] {
                        queue = append(queue, QueueItem{
                            element: second,
                            tier:    secondTier,
                            depth:   current.depth + 1,
                        })
                    }
                }
            }
        }
    }
    
    return res
}