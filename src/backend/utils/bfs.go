package utils

type Tree struct {
    Value    string
    Tier     int
    Children [][2]*Tree
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
            var firstNode, secondNode *Tree

            firstTier := recipes[pair.First].Tier
            secondTier := recipes[pair.Second].Tier

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
                    queue = append(queue, firstNode)
                }

                // Cek kalo dia udah pernah divisit
                if node, ok := visited[pair.Second]; ok {
                    // Kalau udah, ambil node nya 
                    secondNode = node
                } else {
                    // Kalo belum ada, init tree baru isinya elemen child trus masukin ke visited dan queue
                    secondNode = &Tree{Value: pair.Second, Tier: secondTier}
                    visited[pair.Second] = secondNode
                    queue = append(queue, secondNode)
                    }
                current.Children = append(current.Children, [2]*Tree{firstNode, secondNode})
                }
            }      
        }
    return root
}

func (t *Tree) PrintTree(level int) {
    if t == nil {
        return
    }
    for i := 0; i < level; i++ {
        print("  ")
    }
    println(t.Value, "(Tier:", t.Tier, ")")
    for _, children := range t.Children {
        for _, child := range children {
            child.PrintTree(level + 1)
        }
    }
}