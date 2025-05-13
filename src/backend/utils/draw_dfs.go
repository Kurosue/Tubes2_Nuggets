package utils

import (
	"fmt"
	"strings"
)

type TreeNode struct {
	Name     string       `json:"name"`
	Path     RecipePath   `json:"path"`
	Depth    int          `json:"depth"`
	Children []*TreeNode  `json:"children"`
}

func buildTreeNode(recipeMap map[string]RecipePath, path RecipePath, visited map[string]bool, depth int) *TreeNode {
	if _, seen := visited[path.Result]; seen {
		// Prevent infinite recursion
		return &TreeNode{
			Name:     path.Result,
			Path:     path,
			Depth:    depth,
			Children: []*TreeNode{},
		}
	}

	// Mark as visited
	visited[path.Result] = true
	defer delete(visited, path.Result)

	// Create the node
	node := &TreeNode{
		Name:     path.Result,
		Path:     path,
		Depth:    depth,
		Children: []*TreeNode{},
	}

	// Add child nodes recursively
	if path.Ingredient1 != "" {
		if childPath, ok := recipeMap[path.Ingredient1]; ok {
			node.Children = append(node.Children, buildTreeNode(recipeMap, childPath, visited, depth+1))
		}
	}
	if path.Ingredient2 != "" {
		if childPath, ok := recipeMap[path.Ingredient2]; ok {
			node.Children = append(node.Children, buildTreeNode(recipeMap, childPath, visited, depth+1))
		}
	}

	return node
}

func findRootNode(recipePaths []RecipePath, path RecipePath, visited map[string]bool) RecipePath {
	if _, seen := visited[path.Result]; seen {
		return path
	}
	visited[path.Result] = true

	// Look for a parent path where this path is an ingredient
	for _, candidate := range recipePaths {
		if candidate.Ingredient1 == path.Result || candidate.Ingredient2 == path.Result {
			return findRootNode(recipePaths, candidate, visited)
		}
	}

	// If no parent found, this is the root
	return path
}

// VisualizeMessages creates a visual representation of the DFS message path
func VisualizeMessages(recipePaths []RecipePath) string {
    if len(recipePaths) == 0 {
        return "No recipe path found."
    }

    var sb strings.Builder
    sb.WriteString("\n=== Recipe Path ===\n\n")
    
    recipeMap := make(map[string]RecipePath)
    for _, msg := range recipePaths {
        if _, exists := recipeMap[msg.Result]; !exists {
            recipeMap[msg.Result] = msg
        }
    }

    visited1 := make(map[string]bool)
    visited2 := make(map[string]bool)
    treeNodes := []*TreeNode{ buildTreeNode(recipeMap, findRootNode(recipePaths, recipePaths[0], visited1), visited2, 0) }
    
    depth := 0
    for len(treeNodes) > 0 {
        currentTreeNodes := treeNodes
        treeNodes = []*TreeNode {}
        sb.WriteString(fmt.Sprintf("Depth %d:\n", depth))
        for _, treeNode := range currentTreeNodes {
            treeNodes = append(treeNodes, treeNode.Children...)
            recipePath := treeNode.Path
            if recipePath.Ingredient1 == "" && recipePath.Ingredient2 == "" {
                // Base element or target
                sb.WriteString(fmt.Sprintf("  • %s (Base)\n", recipePath.Result))
            } else {
                // Combination
                sb.WriteString(fmt.Sprintf("  • %s = %s + %s\n", 
                    recipePath.Result, recipePath.Ingredient1, recipePath.Ingredient2))
            }
        }
        depth++
    }
    
    sb.WriteString("=== End of Path ===\n")
    return sb.String()
}

// VisualizeMessageTree creates a tree visualization of the messages
func VisualizeMessageTree(recipePaths []RecipePath) string {
    if len(recipePaths) == 0 {
        return "No recipe path found."
    }
    
    // Find the target element (depth 0)
    visited1 := make(map[string]bool)
    target := findRootNode(recipePaths, recipePaths[0], visited1)
    
    // Pre-process messages into a more efficient structure
    // Map from result -> message for quick lookup
    recipeMap := make(map[string]RecipePath)
    for _, msg := range recipePaths {
        if _, exists := recipeMap[msg.Result]; !exists {
            recipeMap[msg.Result] = msg
        }
    }
    
    var sb strings.Builder
    sb.WriteString("\n=== Recipe Tree ===\n\n")
    
    // Draw the tree starting with the target
    visited := make(map[string]bool) // Prevent infinite recursion
    drawMessageTree(&sb, recipeMap, target, "", "", visited, 0, 10) // Max depth 10 to prevent excessive rendering
    
    sb.WriteString("\n=== End of Tree ===\n")
    return sb.String()
}

// Helper function to draw the message tree recursively with cycle detection and depth limiting
func drawMessageTree(sb *strings.Builder, messageMap map[string]RecipePath, currentMsg RecipePath, 
                    prefix string, childrenPrefix string, visited map[string]bool, 
                    currentDepth int, maxDepth int) {
    
    // Check for max depth or cycles
    if currentDepth > maxDepth || visited[currentMsg.Result] {
        return
    }
    
    // Mark as visited for cycle detection
    visited[currentMsg.Result] = true
    
    // Print the current node
    sb.WriteString(prefix)
    
    if currentMsg.Ingredient1 == "" && currentMsg.Ingredient2 == "" {
        // Base element
        sb.WriteString(fmt.Sprintf("%s (Base)\n", currentMsg.Result))
    } else {
        // Combination
        sb.WriteString(fmt.Sprintf("%s = %s + %s (Depth: %d)\n", 
            currentMsg.Result, currentMsg.Ingredient1, currentMsg.Ingredient2, currentDepth))
        
        // Draw branches for ingredients if they exist in our map
        ing1Msg, ing1Exists := messageMap[currentMsg.Ingredient1]
        ing2Msg, ing2Exists := messageMap[currentMsg.Ingredient2]
        
        // Create a new visited map for each branch to allow shared ingredients in different branches
        visited1 := make(map[string]bool)
        for k, v := range visited {
            visited1[k] = v
        }
        
        // Draw ingredient1 branch if found
        if ing1Exists {
            branchPrefix := childrenPrefix + "├── "
            if !ing2Exists {
                branchPrefix = childrenPrefix + "└── " // Last branch
            }
            
            sb.WriteString(branchPrefix)
            drawMessageTree(sb, messageMap, ing1Msg, "", childrenPrefix + "│   ", visited1, currentDepth+1, maxDepth)
        }
        
        // Draw ingredient2 branch if found
        if ing2Exists {
            visited2 := make(map[string]bool)
            for k, v := range visited {
                visited2[k] = v
            }
            
            sb.WriteString(childrenPrefix + "└── ")
            drawMessageTree(sb, messageMap, ing2Msg, "", childrenPrefix + "    ", visited2, currentDepth+1, maxDepth)
        }
    }
}

// VisualizeDFS creates a visualization of DFS results
func VisualizeDFS(result Message) string {
    var sb strings.Builder
    
    // Add tree visualization
    sb.WriteString(VisualizeMessageTree(result.RecipePath))
    
    // Add statistics
    sb.WriteString(fmt.Sprintf("\nTotal messages: %d\n", len(result.RecipePath)))
    sb.WriteString(fmt.Sprintf("Total nodes visited: %d\n", result.NodesVisited))
    
    return sb.String()
}
