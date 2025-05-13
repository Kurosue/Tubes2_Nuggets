package utils

import (
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
	var messages []Message

	var iter func(node *Node)
	iter = func(node *Node) {
		if node == nil {
			return
		}
		iter(node.Ingredient1)
		iter(node.Ingredient2)

		if node.Ingredient1 == nil && node.Ingredient2 == nil {
			// Ini BaseElement (daun)
			messages = append(messages, Message{
				Result: node.Result,
				Depth:  node.Depth,
			})
		} else {
			// Ini kombinasi normal
			messages = append(messages, Message{
				Ingredient1: node.Ingredient1.Result,
				Ingredient2: node.Ingredient2.Result,
				Result:      node.Result,
				Depth:       node.Depth,
			})
		}
	}

	iter(root)

	// Reverse agar target berada di index pertama
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages
}
