package utils

import "sort"

func BFS(recipes RecipeMap, a, b string) []string {
    have := map[string]bool{a: true, b: true}
    queue := []string{a, b}         

    for i := 0; i < len(queue); i++ {
        x := queue[i]
        for y := range have {
            if x == y {
                continue
            }
            key := PairKey(x, y)
            if child, ok := recipes[key]; ok && !have[child] {
                have[child] = true
                queue = append(queue, child)
            }
        }
    }

    var result []string
    for elem := range have {
        if elem != a && elem != b {
            result = append(result, elem)
        }
    }
    sort.Strings(result)
    return result
}
