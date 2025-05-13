package utils
import (
	"strings"
)

func BFSNRecipes(target string, Elements map[string]Element, maxRecipes int, visitedNode *int) []*Node {
	var resultRoots []*Node
	queue := []*Node{{Result: target, Depth: 0}}

	for len(queue) > 0 && len(resultRoots) < maxRecipes {
		current := queue[0]
		queue = queue[1:]
		*visitedNode++

		if BaseElement[current.Result] {
			resultRoots = append(resultRoots, current)
			continue
		}

		elem, ok := Elements[current.Result]
		if !ok {
			continue
		}

		tier := elem.Tier

		for _, recipe := range elem.Recipes {
			parts := strings.Split(recipe, "+")
			if len(parts) != 2 {
				continue
			}

			ing1 := strings.TrimSpace(parts[0][:len(parts[0])-1])
			ing2 := strings.TrimSpace(parts[1][1:])

			e1, ok1 := Elements[ing1]
			e2, ok2 := Elements[ing2]
			if !ok1 || !ok2 {
				continue
			}
			if e1.Tier >= tier || e2.Tier >= tier {
				continue
			}

			// Rekursif bangun pohon dari ingredient
			subVisited := 0
			left := BFSShortestNode(ing1, Elements, &subVisited)
			right := BFSShortestNode(ing2, Elements, &subVisited)

			newNode := &Node{
				Result:      current.Result,
				Depth:       current.Depth,
				Ingredient1: left,
				Ingredient2: right,
			}
			resultRoots = append(resultRoots, newNode)
			if len(resultRoots) >= maxRecipes {
				break
			}
		}
	}

	return resultRoots
}

func FlattenMultipleTrees(roots []*Node) [][]Message {
	var all [][]Message
	for _, root := range roots {
		all = append(all, FlattenTreeToMessages(root))
	}
	return all
}
