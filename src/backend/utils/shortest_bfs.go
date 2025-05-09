package utils

import (
    "container/list"
)

// Status iterasi sekarang (untuk mempermudah akses sama monirtoring)
type State struct {
    Visited map[string]bool // Elemen yang udah dikunjungi di state sekarang
    Path []Message // Path untuk sampai kondisi sekarang
}

func BFSShortestPath(target string, recipeMap RecipeMap, elements RecipeElement) []Message {

    // Cek kalau elemen yang dicari base elemen jadi ga perlu dicari
    if BaseElement[target] {
        return []Message{}
    }

    // Inisialisasi Var yang dibutuhin
    queue := list.New()
    initialState := State{
        Visited: make(map[string]bool),
        Path:    []Message{},
    }
    
    // Init state untuk cek elmeen yang udah dibentuk sama path untuk sampai curr state
    for elem := range BaseElement {
        initialState.Visited[elem] = true
    }
    
    queue.PushBack(initialState)
    visited := make(map[string]bool)
    
    // BFS loop
    for queue.Len() > 0 {
        // Dequeue the front state
        current := queue.Front().Value.(State)
        queue.Remove(queue.Front())
        
        // Kalau sampai target, selesai
        if current.Visited[target] {
            return current.Path
        }
        
        stateKey := stateToString(current.Visited)
        if visited[stateKey] {
            continue
        }
        visited[stateKey] = true
        
        for combinationStr, result := range recipeMap {
            if current.Visited[result] {
                continue
            }
            
            elem1, elem2 := DecomposeKey(combinationStr)
            
            if current.Visited[elem1] && current.Visited[elem2] {
                newState := State{
                    Visited: make(map[string]bool),
                    Path:    make([]Message, len(current.Path)),
                }
                
                for elem := range current.Visited {
                    newState.Visited[elem] = true
                }
                
                newState.Visited[result] = true
                
                copy(newState.Path, current.Path)
                
                newMsg := Message{
                    Ingredient1: elem1,
                    Ingredient2: elem2,
                    Result:      result,
                    Depth:       len(current.Path),
                }
                newState.Path = append(newState.Path, newMsg)
                
                if result == target {
                    return newState.Path
                }
                
                queue.PushBack(newState)
            }
        }
    }
    
    return []Message{}
}

// Helper function 
func stateToString(Visited map[string]bool) string {
    result := ""
    keys := make([]string, 0, len(Visited))
    for k := range Visited {
        if Visited[k] {
            keys = append(keys, k)
        }
    }
    
    for i := 0; i < len(keys); i++ {
        for j := i + 1; j < len(keys); j++ {
            if keys[i] > keys[j] {
                keys[i], keys[j] = keys[j], keys[i]
            }
        }
    }
    
    for _, k := range keys {
        result += k + ","
    }
    return result
}