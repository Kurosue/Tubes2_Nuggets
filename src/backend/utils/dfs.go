package utils

type Step struct {
    X, Y, Result string
}

func dfsHelper(recipes RecipeMap, have map[string]bool, history []Step, target string) ([]Step, bool) {
    if have[target] {
        return history, true
    }

    for x := range have {
        for y := range have {
            if x == y { continue }
            key := PairKey(x, y)
            if child, ok := recipes[key]; ok && !have[child] {
                have[child] = true
                history = append(history, Step{X: x, Y: y, Result: child})
                if path, found := dfsHelper(recipes, have, history, target); found {
                    return path, true
				}
                delete(have, child)
                history = history[:len(history)-1]
            }
        }
    }
    return nil, false
}

func DFS(recipes RecipeMap, a, b, target string) []Step {
    have := map[string]bool{a: true, b: true}
    path, _ := dfsHelper(recipes, have, nil, target)
    return path
}
