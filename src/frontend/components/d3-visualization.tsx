import React, { useEffect, useRef, useImperativeHandle } from "react";
import * as d3 from "d3";
import { RecipePath } from "@/lib/api/service";

type D3CanvasRefType = {
  handler: {
    refreshData: (messages: RecipePath[]) => void;
  };
};

const D3Canvas = React.forwardRef(function D3Canvas(
  { className }: { className: string },
  ref: React.Ref<D3CanvasRefType>
) {
  const svgRef = useRef<SVGSVGElement>(null);
  const containerRef = useRef<SVGGElement | null>(null);
  const currentMessages = useRef<RecipePath[]>([]);
  const zoomTransformRef = useRef<any>(null);

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
        zoomTransformRef.current = event.transform;
        container.attr("transform", event.transform);
      });
    
    svg.call(zoom as any);
    
    // Initial render with empty data
    refreshData([]);

    svg.append("text")
    .attr("id", "emoji-placeholder")
    .attr("x", "50%")
    .attr("y", "45%")
    .attr("text-anchor", "middle")
    .attr("dominant-baseline", "middle")
    .attr("font-size", "50px")
    .attr("fill", "#ccc")
    .text("🧪");

    // Caption under the emoji
    svg.append("text")
        .attr("id", "caption-placeholder")
        .attr("x", "50%")
        .attr("y", "55%")
        .attr("text-anchor", "middle")
        .attr("dominant-baseline", "middle")
        .attr("font-size", "16px")
        .attr("fill", "#999")
        .text("Enter an element to discover its magical recipe");
    
    // Handle window resize
    const handleResize = () => {
      refreshData(currentMessages.current);
    };
    
    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  }, []);

  // The main function to refresh the tree with new message data
  const refreshData = (recipePath: RecipePath[]) => {
    if (!containerRef.current || !svgRef.current) return;
    
    // Store current messages
    currentMessages.current = recipePath;
    
    if (recipePath.length == 0) return;
    
    const width = svgRef.current.clientWidth;
    const height = svgRef.current.clientHeight;
    const svg = d3.select(svgRef.current);
    
    // Clear content
    d3.select(containerRef.current).selectAll("*").remove();
    d3.select("#emoji-placeholder").remove();
    d3.select("#caption-placeholder").remove();

    // Find target element
    function findRootNode(
      path: RecipePath, 
      visited: Set<string> = new Set()
    ) {
      const next = recipePath.find(m => m.ingredient1 == path.result || m.ingredient2 == path.result);
      if(next == null || visited.has(path.result)) return path;
      visited.add(path.result);
      return findRootNode(next, visited);
    }
    const targetElement = findRootNode(recipePath[0]);
    
    const recipeMap: Record<string, RecipePath> = {};
    recipePath.forEach(msg => {
      if (!recipeMap[msg.result]) {
        recipeMap[msg.result] = msg;
      }
    });
    
    // Helper
    interface TreeNode {
      name: string;
      path: RecipePath;
      depth: number;
      children: TreeNode[];
    }
    
    // Recursively build the tree - similar to your Go code
    function buildTreeNode(
      path: RecipePath, 
      visited: Set<string> = new Set(),
      depth = 0
    ): TreeNode {
      if (visited.has(path.result)) { // Prevent infinite recursion
        return { name: path.result, path: path, depth: depth, children: [] };
      }
      
      // Create node for current message
      const node: TreeNode = {
        name: path.result,
        path: path,
        depth: depth,
        children: []
      };
      
      // Mark as visited
      try {
        visited.add(path.result);
        
        // Add ingredient1 if it exists
        if (path.ingredient1 && recipeMap[path.ingredient1])
          node.children.push(buildTreeNode(recipeMap[path.ingredient1], visited, depth + 1));
        
        // Add ingredient2 if it exists
        if (path.ingredient2 && recipeMap[path.ingredient2])
          node.children.push(buildTreeNode(recipeMap[path.ingredient2], visited, depth + 1));
      } finally {
        visited.delete(path.result);
      }
      
      return node;
    }
    
    // Build the hierarchical tree data
    const treeData = buildTreeNode(targetElement);
    
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
      .attr("stroke", d => d.data.path.ingredient1 == "" ? "#ff8800" : "#999") // Highlight base elements
      .attr("stroke-width", 1.5);
    
    // Add element images
    nodes.append("image")
      .attr("xlink:href", d => new URL(`/images/${d.data.name}.svg`, process.env.NEXT_PUBLIC_BACKEND).href)
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
    nodes.filter(d => d.data.path.ingredient1 == "")
      .append("text")
      .attr("dy", 30)
      .attr("text-anchor", "middle")
      .attr("font-size", "10px")
      .attr("fill", d => d.data.path.ingredient1 == "" ? "#ff8800" : "#333") // Highlight base elements
      .attr("font-weight", d => d.data.path.ingredient1 == "" ? "bold" : "normal")
      .attr("stroke", "black")  // Add white outline
      .attr("stroke-width", "0.3px") // Thin white stroke
      .attr("paint-order", "stroke") 
      .text(d => {
        const recipePath = d.data.path;
        // Show "Base" suffix for base elements
        return recipePath.ingredient1 == "" && recipePath.ingredient2 == "" ? 
          `${recipePath.result} (Base)` : recipePath.result;
      });
    
    // Add recipe formula and depth (for non-base elements)
    nodes.filter(d => d.data.path.ingredient1 != "")
      .append("text")
      .attr("dy", -25)
      .attr("text-anchor", "middle")
      .attr("font-size", "9px")
      .attr("fill", d => d.data.path.ingredient1 == "" ? "#ff8800" : "#333") // Highlight base elements
      .attr("font-weight", d => d.data.path.ingredient1 == "" ? "bold" : "normal")
      .attr("stroke", "black")  // Add white outline
      .attr("stroke-width", "0.3px") // Thin white stroke
      .attr("paint-order", "stroke") 
      .text(d => {
        const recipePath = d.data.path;
        return `${recipePath.result} = ${recipePath.ingredient1} + ${recipePath.ingredient2} (Depth: ${d.data.depth})`;
      });
    
    // Add depth label for all nodes
    nodes.append("text")
      .attr("dy", 45)
      .attr("text-anchor", "middle")
      .attr("font-size", "9px")
      .attr("fill", "#666")
      .attr("fill", d => d.data.path.ingredient1 == "" ? "#ff8800" : "#333") // Highlight base elements
      .attr("font-weight", d => d.data.path.ingredient1 == "" ? "bold" : "normal")
      .attr("stroke", "black")  // Add white outline
      .attr("stroke-width", "0.3px") // Thin white stroke
      .attr("paint-order", "stroke") 
      .text(d => `Depth: ${d.data.depth}`);
    container.attr("transform", zoomTransformRef.current);
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
