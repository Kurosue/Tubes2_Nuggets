package utils

import (
	"sort"
	"strings"
)

// BaseElement adalah map yang menyimpan elemen dasar
var BaseElement = map[string]bool{
    "Air":   true,
    "Water": true,
    "Earth": true,
    "Fire":  true,
}

type Node struct {
    Ingredient1 *Node
    Ingredient2 *Node
    Result      string
    Depth       int
}

func BFSShortestNode(target string, Elements map[string]Element, visitedNode *int) (*Node) {
    root := &Node{Result: target, Depth: 0}
    nodeMap := map[string]*Node{target: root}
    queue := []*Node{root}

    for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:]
		*visitedNode++

        // Jika base element, tidak perlu proses resep
        if BaseElement[current.Result] {
            continue
        }

        elem, ok := Elements[current.Result]
        if !ok {
            continue
        }

        tier := elem.Tier

        for _, recipe := range elem.Recipes {
            parts := strings.Split(recipe, "+")
            ing1 := strings.TrimSpace(parts[0][:len(parts[0])-1])
            ing2 := strings.TrimSpace(parts[1][1:])

            // Cek ingredient valid dan tier lebih kecil dari result
            e1, ok1 := Elements[ing1]
            e2, ok2 := Elements[ing2]
            if !ok1 || !ok2 {
                continue
            }

            if e1.Tier >= tier || e2.Tier >= tier {
                continue
            }

            // Buat node ingredient 1
            if _, exists := nodeMap[ing1]; !exists {
                nodeMap[ing1] = &Node{Result: ing1, Depth: current.Depth + 1}
                queue = append(queue, nodeMap[ing1])
            }

            // Buat node ingredient 2
            if _, exists := nodeMap[ing2]; !exists {
                nodeMap[ing2] = &Node{Result: ing2, Depth: current.Depth + 1}
                queue = append(queue, nodeMap[ing2])
            }

            current.Ingredient1 = nodeMap[ing1]
            current.Ingredient2 = nodeMap[ing2]
            break // resep valid pertama
        }
    }

    return root
}


func FlattenTreeToMessages(root *Node) []Message {
    var result []Message
    visited := make(map[string]bool)

    var iter func(*Node)
    iter = func(n *Node) {
        if n == nil || visited[n.Result] {
            return
        }

        iter(n.Ingredient1)
        iter(n.Ingredient2)

        visited[n.Result] = true

        var ing1, ing2 string
        if n.Ingredient1 != nil {
            ing1 = n.Ingredient1.Result
        }
        if n.Ingredient2 != nil {
            ing2 = n.Ingredient2.Result
        }

        result = append(result, Message{
            Ingredient1: ing1,
            Ingredient2: ing2,
            Result:      n.Result,
            Depth:       n.Depth,
        })
    }

    iter(root)

    sort.Slice(result, func(i, j int) bool {
        return result[i].Depth < result[j].Depth
    })

    return result
}
