package utils

import (
    "fmt"
    "strings"
)

// VisualizeMessages creates a visual representation of the DFS message path
func VisualizeMessages(messages []Message) string {
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
    
    var sb strings.Builder
    sb.WriteString("\n=== Recipe Tree ===\n\n")
    
    // Draw the tree starting with the target
    drawMessageTree(&sb, messages, target, "", "")
    
    sb.WriteString("\n=== End of Tree ===\n")
    return sb.String()
}

// Helper function to draw the message tree recursively
func drawMessageTree(sb *strings.Builder, messages []Message, currentMsg Message, prefix string, childrenPrefix string) {
    // Print the current node
    sb.WriteString(prefix)
    
    if currentMsg.ingredient1 == "" && currentMsg.ingredient2 == "" {
        // Base element
        sb.WriteString(fmt.Sprintf("%s (Base)\n", currentMsg.result))
    } else {
        // Combination
        sb.WriteString(fmt.Sprintf("%s = %s + %s (Depth: %d)\n", 
            currentMsg.result, currentMsg.ingredient1, currentMsg.ingredient2, currentMsg.depth))
        
        // Find child messages for ingredient1 and ingredient2
        var ing1Messages []Message
        var ing2Messages []Message
        
        for _, msg := range messages {
            if msg.result == currentMsg.ingredient1 {
                ing1Messages = append(ing1Messages, msg)
            } else if msg.result == currentMsg.ingredient2 {
                ing2Messages = append(ing2Messages, msg)
            }
        }
        
        // Draw ingredient1 branch if found
        if len(ing1Messages) > 0 {
            // Sort by depth to get the one with the lowest depth
            var ing1Msg Message = ing1Messages[0]
            for _, msg := range ing1Messages {
                if msg.depth < ing1Msg.depth {
                    ing1Msg = msg
                }
            }
            
            // Draw the branch
            sb.WriteString(childrenPrefix + "├── ")
            drawMessageTree(sb, messages, ing1Msg, "", childrenPrefix + "│   ")
        }
        
        // Draw ingredient2 branch if found
        if len(ing2Messages) > 0 {
            // Sort by depth to get the one with the lowest depth
            var ing2Msg Message = ing2Messages[0]
            for _, msg := range ing2Messages {
                if msg.depth < ing2Msg.depth {
                    ing2Msg = msg
                }
            }
            
            // Draw the branch
            sb.WriteString(childrenPrefix + "└── ")
            drawMessageTree(sb, messages, ing2Msg, "", childrenPrefix + "    ")
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