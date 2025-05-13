package utils

import (
	"sort"
	"time"
)

type Node struct {
    Step int
    Ingredient1 *Node
    Ingredient2 *Node
    Result      string
}

func BFSShortestNodeImpl(target string, Elements map[string]Element, visitedNode *int) (*Node) {
    step := 0
    root := &Node{ Result: target, Step: step }
    nodeMap := map[string]*Node{ target: root }
    queue := []*Node{ root }

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

        for _, recipe := range elem.Recipes {
            ing1 := recipe[0]
			ing2 := recipe[1]

            // Cek ingredient valid dan tier lebih kecil dari result
            e1, ok1 := Elements[ing1]
            e2, ok2 := Elements[ing2]
            if !ok1 || !ok2 {
                continue
            }

            if e1.Tier >= elem.Tier || e2.Tier >= elem.Tier {
                continue
            }

            // Buat node ingredient 1
            if _, exists := nodeMap[ing1]; !exists {
                step++
                nodeMap[ing1] = &Node{ Result: ing1, Step: step }
                queue = append(queue, nodeMap[ing1])
            }

            // Buat node ingredient 2
            if _, exists := nodeMap[ing2]; !exists {
                step++
                nodeMap[ing2] = &Node{ Result: ing2, Step: step }
                queue = append(queue, nodeMap[ing2])
            }

            current.Ingredient1 = nodeMap[ing1]
            current.Ingredient2 = nodeMap[ing2]
            break // resep valid pertama
        }
    }
    return root;
}
func BFSShortestNode(target string, Elements map[string]Element, resultChan chan Message) {
    visitedNode := 0
    start := time.Now().UnixNano()
    root := BFSShortestNodeImpl(target, Elements, &visitedNode)
    resultChan <- Message{
        RecipePath: FlattenTreeToRecipePaths(root),
        NodesVisited: visitedNode,
        Duration: float32(time.Now().UnixNano() - start) / 1000000,
    }
}

func FlattenTreeToRecipePaths(root *Node) []RecipePath {
    var nodes []*Node
	var recipePaths []RecipePath

	var iter func(node *Node)
	iter = func(node *Node) {
		if node == nil {
			return
		}
		iter(node.Ingredient1)
		iter(node.Ingredient2)
        nodes = append(nodes, node)
	}
	iter(root)

    sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Step < nodes[j].Step
	})

    for _, node := range nodes {
		if node.Ingredient1 == nil && node.Ingredient2 == nil {
			// Ini BaseElement (daun)
			recipePaths = append(recipePaths, RecipePath{
                Ingredient1: "",
                Ingredient2: "",
				Result: node.Result,
			})
		} else {
			// Ini kombinasi normal
			recipePaths = append(recipePaths, RecipePath{
				Ingredient1: node.Ingredient1.Result,
				Ingredient2: node.Ingredient2.Result,
				Result:      node.Result,
			})
		}
    }
	return recipePaths
}
