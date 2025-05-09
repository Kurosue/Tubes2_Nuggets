package utils

import (
    "fmt"
    "strings"
)

// TreeNode represents a node in the recipe tree
type TreeNode struct {
    Element  Element
    Children []*TreeNode
    IsBase   bool
}

// BuildRecipeTree constructs a tree from the DFS path
func BuildRecipeTree(elements []Element, recipeMap RecipeMap) *TreeNode {
    if len(elements) == 0 {
        return nil
    }

    // Create root node with target element
    root := &TreeNode{
        Element:  elements[0],
        Children: []*TreeNode{},
        IsBase:   false,
    }

    // Track processed elements
    processed := make(map[string]bool)
    processed[root.Element.Name] = true

    // Create element lookup
    elementMap := make(map[string]Element)
    for _, el := range elements {
        elementMap[el.Name] = el
    }

    // Process the remaining elements from the path
    buildTreeHelper(root, elements, processed, elementMap, recipeMap)

    return root
}

// buildTreeHelper recursively builds the tree
func buildTreeHelper(node *TreeNode, elements []Element, processed map[string]bool, 
                     elementMap map[string]Element, recipeMap RecipeMap) {
    // Check if current element is a base element
    if node.Element.Name == "Air" || node.Element.Name == "Water" || 
       node.Element.Name == "Earth" || node.Element.Name == "Fire" {
        node.IsBase = true
        return
    }

    // Find ingredients for this element
    for key, result := range recipeMap {
        if result == node.Element.Name {
            ing1, ing2 := DecomposeKey(key)
            
            // Create child nodes for ingredients
            if el1, exists := elementMap[ing1]; exists {
                child1 := &TreeNode{
                    Element:  el1,
                    Children: []*TreeNode{},
                    IsBase:   ing1 == "Air" || ing1 == "Water" || ing1 == "Earth" || ing1 == "Fire",
                }
                node.Children = append(node.Children, child1)
                
                // Recursively process this child if not already processed
                if !processed[ing1] {
                    processed[ing1] = true
                    buildTreeHelper(child1, elements, processed, elementMap, recipeMap)
                }
            }
            
            if el2, exists := elementMap[ing2]; exists {
                child2 := &TreeNode{
                    Element:  el2,
                    Children: []*TreeNode{},
                    IsBase:   ing2 == "Air" || ing2 == "Water" || ing2 == "Earth" || ing2 == "Fire",
                }
                node.Children = append(node.Children, child2)
                
                // Recursively process this child if not already processed
                if !processed[ing2] {
                    processed[ing2] = true
                    buildTreeHelper(child2, elements, processed, elementMap, recipeMap)
                }
            }
            
            // Only use the first recipe found
            return
        }
    }
}

// DrawRecipeTree creates a visual tree representation of the recipe path
func DrawRecipeTree(root *TreeNode) string {
    var sb strings.Builder
    sb.WriteString("\n=== Recipe Tree ===\n\n")
    drawTreeNode(&sb, root, "", "")
    sb.WriteString("\n=== End of Tree ===\n")
    return sb.String()
}

// Helper function to draw tree recursively
func drawTreeNode(sb *strings.Builder, node *TreeNode, prefix string, childrenPrefix string) {
    if node == nil {
        return
    }

    // Print the current node
    sb.WriteString(prefix)
    if node.IsBase {
        sb.WriteString(fmt.Sprintf("%s (Base)\n", node.Element.Name))
    } else {
        sb.WriteString(fmt.Sprintf("%s\n", node.Element.Name))
    }

    // Print the children
    for i, child := range node.Children {
        isLast := i == len(node.Children)-1
        
        if isLast {
            drawTreeNode(sb, child, childrenPrefix + "└── ", childrenPrefix + "    ")
        } else {
            drawTreeNode(sb, child, childrenPrefix + "├── ", childrenPrefix + "│   ")
        }
    }
}

// Function to use in tests or applications
func CreateRecipeTree(elements []Element, recipeMap RecipeMap) string {
    // Build and draw tree
    root := BuildRecipeTree(elements, recipeMap)
    return DrawRecipeTree(root)
}