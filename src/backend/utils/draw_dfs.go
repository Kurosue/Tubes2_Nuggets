package utils

import (
    "fmt"
    "strings"
)

// VisualizeMessages creates a visual representation of the DFS message path
func VisualizeMessages(messages []Message) string {
    if len(messages) == 0 {
        return "No recipe path found."
    }

    var sb strings.Builder
    sb.WriteString("\n=== Recipe Path ===\n\n")
    
    // Organize messages by depth
    messagesByDepth := make(map[int][]Message)
    maxDepth := 0
    
    for _, msg := range messages {
        messagesByDepth[msg.depth] = append(messagesByDepth[msg.depth], msg)
        if msg.depth > maxDepth {
            maxDepth = msg.depth
        }
    }
    
    // Print messages by depth
    for depth := 0; depth <= maxDepth; depth++ {
        if msgs, exists := messagesByDepth[depth]; exists {
            sb.WriteString(fmt.Sprintf("Depth %d:\n", depth))
            for _, msg := range msgs {
                if msg.ingredient1 == "" && msg.ingredient2 == "" {
                    // Base element or target
                    sb.WriteString(fmt.Sprintf("  • %s (Base)\n", msg.result))
                } else {
                    // Combination
                    sb.WriteString(fmt.Sprintf("  • %s = %s + %s\n", 
                        msg.result, msg.ingredient1, msg.ingredient2))
                }
            }
            sb.WriteString("\n")
        }
    }
    
    sb.WriteString("=== End of Path ===\n")
    return sb.String()
}

// VisualizeMessageTree creates a tree visualization of the messages
func VisualizeMessageTree(messages []Message) string {
    if len(messages) == 0 {
        return "No recipe path found."
    }
    
    // Find the target element (depth 0)
    var target Message
    for _, msg := range messages {
        if msg.depth == 0 {
            target = msg
            break
        }
    }
    
    // Pre-process messages into a more efficient structure
    // Map from result -> message for quick lookup
    messageMap := make(map[string]Message)
    for _, msg := range messages {
        if existing, exists := messageMap[msg.result]; !exists || msg.depth < existing.depth {
            messageMap[msg.result] = msg
        }
    }
    
    var sb strings.Builder
    sb.WriteString("\n=== Recipe Tree ===\n\n")
    
    // Draw the tree starting with the target
    visited := make(map[string]bool) // Prevent infinite recursion
    drawMessageTree(&sb, messageMap, target, "", "", visited, 0, 10) // Max depth 10 to prevent excessive rendering
    
    sb.WriteString("\n=== End of Tree ===\n")
    return sb.String()
}

// Helper function to draw the message tree recursively with cycle detection and depth limiting
func drawMessageTree(sb *strings.Builder, messageMap map[string]Message, currentMsg Message, 
                    prefix string, childrenPrefix string, visited map[string]bool, 
                    currentDepth int, maxDepth int) {
    
    // Check for max depth or cycles
    if currentDepth > maxDepth || visited[currentMsg.result] {
        return
    }
    
    // Mark as visited for cycle detection
    visited[currentMsg.result] = true
    
    // Print the current node
    sb.WriteString(prefix)
    
    if currentMsg.ingredient1 == "" && currentMsg.ingredient2 == "" {
        // Base element
        sb.WriteString(fmt.Sprintf("%s (Base)\n", currentMsg.result))
    } else {
        // Combination
        sb.WriteString(fmt.Sprintf("%s = %s + %s (Depth: %d)\n", 
            currentMsg.result, currentMsg.ingredient1, currentMsg.ingredient2, currentMsg.depth))
        
        // Draw branches for ingredients if they exist in our map
        ing1Msg, ing1Exists := messageMap[currentMsg.ingredient1]
        ing2Msg, ing2Exists := messageMap[currentMsg.ingredient2]
        
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
func VisualizeDFS(result DFSResult) string {
    var sb strings.Builder
    
    // Add tree visualization
    sb.WriteString(VisualizeMessageTree(result.Messages))
    
    // Add statistics
    sb.WriteString(fmt.Sprintf("\nTotal messages: %d\n", len(result.Messages)))
    sb.WriteString(fmt.Sprintf("Total nodes visited: %d\n", result.NodesVisited))
    
    return sb.String()
}