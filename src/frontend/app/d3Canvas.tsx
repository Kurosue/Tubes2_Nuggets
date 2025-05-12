import React, { useEffect, useRef, useImperativeHandle } from "react";
import * as d3 from "d3";
import { Message } from "@/lib/api/service";

type D3CanvasRefType = {
  handler: {
    refreshData: (messages: Message[]) => void;
  };
};

const D3Canvas = React.forwardRef(function D3Canvas(
  { className }: { className: string },
  ref: React.Ref<D3CanvasRefType>
) {
  const svgRef = useRef<SVGSVGElement>(null);
  const containerRef = useRef<SVGGElement | null>(null);
  const currentMessages = useRef<Message[]>([]);

  useEffect(() => {
    if (!svgRef.current) return;

    const svg = d3.select(svgRef.current);
    svg.selectAll("*").remove();
    
    // Create container for zooming
    const container = svg.append("g").attr("class", "container");
    containerRef.current = container.node()!;
    
    // Add zoom behavior
    const zoom = d3.zoom()
      .scaleExtent([0.1, 3])
      .on("zoom", (event) => {
        container.attr("transform", event.transform);
      });
    
    svg.call(zoom as any);
    
    // Initial render with empty data
    refreshData([]);
    
    // Handle window resize
    const handleResize = () => {
      refreshData(currentMessages.current);
    };
    
    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  }, []);

  // The main function to refresh the tree with new message data
  const refreshData = (messages: Message[]) => {
    if (!containerRef.current || !svgRef.current) return;
    
    // Store current messages
    currentMessages.current = messages;
    
    if (messages.length === 0) return;
    
    const width = svgRef.current.clientWidth;
    const height = svgRef.current.clientHeight;
    
    // Clear content
    d3.select(containerRef.current).selectAll("*").remove();
    
    // Find target element
    const targetMsg = messages.find(msg => msg.Depth === 0) || messages[0];
    
    const messageMap: Record<string, Message> = {};
    messages.forEach(msg => {
      if (!messageMap[msg.Result] || msg.Depth < messageMap[msg.Result].Depth) {
        messageMap[msg.Result] = msg;
      }
    });
    
    // Helper
    interface TreeNode {
      name: string;
      message: Message;
      children: TreeNode[];
    }
    
    // Recursively build the tree - similar to your Go code
    function buildTreeNode(
      msg: Message, 
      visited: Set<string> = new Set(),
      depth = 0
    ): TreeNode {
      if (depth > 10 || visited.has(msg.Result)) { // Prevent infinite recursion
        return { name: msg.Result, message: msg, children: [] };
      }
      
      // Create node for current message
      const node: TreeNode = {
        name: msg.Result,
        message: msg,
        children: []
      };
      
      // Mark as visited
      visited.add(msg.Result);
      
      // Add ingredient1 if it exists
      if (msg.Ingredient1 && messageMap[msg.Ingredient1]) {
        // Create copy of visited set for this branch
        const visited1 = new Set(visited);
        node.children.push(buildTreeNode(messageMap[msg.Ingredient1], visited1, depth + 1));
      }
      
      // Add ingredient2 if it exists
      if (msg.Ingredient2 && messageMap[msg.Ingredient2]) {
        // Create copy of visited set for this branch
        const visited2 = new Set(visited);
        node.children.push(buildTreeNode(messageMap[msg.Ingredient2], visited2, depth + 1));
      }
      
      return node;
    }
    
    // Build the hierarchical tree data
    const treeData = buildTreeNode(targetMsg);
    
    // Create d3 hierarchy
    const root = d3.hierarchy(treeData);
    
    // Create tree layout
    const treeLayout = d3.tree<TreeNode>()
      .nodeSize([100, 120]); // [horizontal, vertical] spacing
    
    // Apply the layout
    const hierarchyData = treeLayout(root);
    
    // Create container group and center it
    const container = d3.select(containerRef.current);
    
    // Center the tree horizontally
    const descendants = hierarchyData.descendants();
    const minX = d3.min(descendants, d => d.x) || 0;
    const maxX = d3.max(descendants, d => d.x) || 0;
    const centerX = (maxX + minX) / 2;
    
    // Transform to center the tree
    container.attr("transform", `translate(${width / 2 - centerX}, 50)`);
    
    // Create container for links
    const linkGroup = container.append("g").attr("class", "links");
    
    // Draw links
    linkGroup.selectAll("path")
      .data(hierarchyData.links())
      .enter()
      .append("path")
      .attr("fill", "none")
      .attr("stroke", "#999")
      .attr("stroke-width", 1.5)
      .attr("d", d3.linkVertical<any, any>()
        .x(d => d.x)
        .y(d => d.y)
      );
    
    // Create container for nodes
    const nodeGroup = container.append("g").attr("class", "nodes");
    
    // Create node groups
    const nodes = nodeGroup.selectAll(".node")
      .data(hierarchyData.descendants())
      .enter()
      .append("g")
      .attr("class", "node")
      .attr("transform", d => `translate(${d.x},${d.y})`);
    
    // Add white background circle for better visibility
    nodes.append("circle")
      .attr("r", 20)
      .attr("fill", "white")
      .attr("stroke", d => d.data.message.Ingredient1 === "" ? "#ff8800" : "#999") // Highlight base elements
      .attr("stroke-width", 1.5);
    
    // Add element images
    nodes.append("image")
      .attr("xlink:href", d => `/images/${d.data.name}.svg`)
      .attr("width", 40)
      .attr("height", 40)
      .attr("x", -20)
      .attr("y", -20)
      .on("error", function() {
        // Fallback for failed images
        const d = d3.select(this).datum() as any;
        d3.select(this)
          .attr("xlink:href", `https://placehold.co/40x40/orange/white?text=${d.data.name.charAt(0)}`);
      });
    
    // Add element name labels
    nodes.append("text")
      .attr("dy", 30)
      .attr("text-anchor", "middle")
      .attr("font-size", "10px")
      .text(d => {
        const msg = d.data.message;
        // Show "Base" suffix for base elements
        return msg.Ingredient1 === "" && msg.Ingredient2 === "" ? 
          `${msg.Result} (Base)` : msg.Result;
      });
    
    // Add recipe formula and depth (for non-base elements)
    nodes.filter(d => d.data.message.Ingredient1 !== "")
      .append("text")
      .attr("dy", -25)
      .attr("text-anchor", "middle")
      .attr("font-size", "9px")
      .text(d => {
        const msg = d.data.message;
        return `${msg.Result} = ${msg.Ingredient1} + ${msg.Ingredient2} (Depth: ${msg.Depth})`;
      });
    
    // Add depth label for all nodes
    nodes.append("text")
      .attr("dy", 45)
      .attr("text-anchor", "middle")
      .attr("font-size", "9px")
      .attr("fill", "#666")
      .text(d => `Depth: ${d.data.message.Depth}`);
  };
  
  // Expose method to refresh the tree data
  useImperativeHandle(ref, () => ({
    handler: {
      refreshData
    }
  }), []);
  
  return (
    <svg 
      ref={svgRef} 
      className={className}
      width="100%"
      height="100%"
    ></svg>
  );
});

export default D3Canvas;