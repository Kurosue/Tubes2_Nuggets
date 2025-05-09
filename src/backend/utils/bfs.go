package utils

import "fmt"

type Tree struct {
    Value    string
    Tier     int
    Children [][2]*Tree
}

var BaseElement = map[string]bool{
    "Fire":   true,
    "Water":  true,
    "Earth":  true,
    "Air":    true,
}

func BFS_Tree(start string, recipes ElementMap) *Tree {
    // Init variable yang diperluin
    visited := make(map[string]*Tree)
    root := &Tree{Value: start, Tier: recipes[start].Tier}
    queue := []*Tree{root}
    
    visited[start] = root

    for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:]

        // Iterasi untuk setiap child (resep dari elemen current)
        for _, pair := range recipes[current.Value].ParsedRecipes {
            fmt.Printf("Pair: %s + %s\n", pair.First, pair.Second)
            // Cek kalau data ada di elements hasil scrap kalau ga atau NULL continue
            if recipes[pair.First] == nil || recipes[pair.Second] == nil {
                continue
            } 

            var firstNode, secondNode *Tree

            firstTier := recipes[pair.First].Tier
            secondTier := recipes[pair.Second].Tier
            
            // If both are base elements, add them to the tree but don't traverse further
            if BaseElement[pair.First] && BaseElement[pair.Second] {
                firstNode = &Tree{Value: pair.First, Tier: firstTier}
                secondNode = &Tree{Value: pair.Second, Tier: secondTier}
                
                current.Children = append(current.Children, [2]*Tree{firstNode, secondNode})
                break
            }
            
            // Continue with normal traversal if not both base elements
            // Tier child tidak boleh lebih dari current (sesuai spek cihuy)
            if(firstTier < current.Tier && secondTier < current.Tier) {
                // Cek kalo dia udah pernah divisit
                if node, ok := visited[pair.First]; ok {
                    // Kalau udah, ambil node nya
                    firstNode = node
                } else {
                    // Kalo belum ada, init tree baru isinya elemen child trus masukin ke `visited dan queue
                    firstNode = &Tree{Value: pair.First, Tier: firstTier}
                    visited[pair.First] = firstNode
                    // Cek kalau node ini tuh base element, kalalu ga, masuk queue
                    if _, ok := BaseElement[pair.First]; !ok {
                        queue = append(queue, firstNode)
                    }
                }

                // Cek kalo dia udah pernah divisit
                if node, ok := visited[pair.Second]; ok {
                    // Kalau udah, ambil node nya 
                    secondNode = node
                } else {
                    // Kalo belum ada, init tree baru isinya elemen child trus masukin ke visited dan queue
                    secondNode = &Tree{Value: pair.Second, Tier: secondTier}
                    visited[pair.Second] = secondNode
                    if _, ok := BaseElement[pair.Second]; !ok {
                        queue = append(queue, secondNode)
                    }
                }

                current.Children = append(current.Children, [2]*Tree{firstNode, secondNode})
            }
        }      
    }
    return root
}