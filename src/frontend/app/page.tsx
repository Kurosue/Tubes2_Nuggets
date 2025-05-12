"use client";
import React, { useCallback, useEffect, useRef, useState } from "react";
import { Separator } from "@/components/ui/separator";
import { Button } from "@/components/ui/button";
import { useIsomorphicLayoutEffect } from "@/lib/utils";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { Slider } from "@/components/ui/slider";
import { Card, CardContent } from "@/components/ui/card";

// Let's replace the deprecated toast with sonner
import { toast } from "sonner";

// Import D3Canvas component
import D3Canvas from "./d3Canvas";

// Import API services
import { fetchElements, createRecipeWebSocket, Element as ElementData, Message, AlgorithmResponse , TimingInfo} from "@/lib/api/service";

function getElementImageUrl(elementName: string): string {
  return `/api/element-image/${elementName}`;
}

const splashTexts = [
  "🌍🔍 Temukan Semua 720 Elemen dari 4 Unsur Dasar!",
  "🧪💡 BFS dan DFS Siap Menguak Resep Alkimia!",
  "🌊🔥💨🌱 Dari Dasar Menuju Keajaiban—Gabungkan dan Temukan!",
  "🧠⚙️ Strategi Algoritma Bertemu Dunia Alkimia!",
  "🧬✨ Temukan Resep Tersembunyi dengan Sekali Klik!",
  "🌀🚀 Pilih DFS atau BFS—Raih Elemen Impianmu!",
  "🌳🧩 Visualisasi Pohon Resep yang Seru dan Informatif!",
  "⏱️📊 Ukur Kecepatanmu—Cek Waktu dan Node yang Dilalui!",
  "🔁🎯 Temukan Banyak Resep dalam Sekejap dengan Multithreading!",
  "🕹️🧙‍♂️ Jadi Alkemis Digital Terbaik di Dunia Little Alchemy 2!",
  "INI TUGASS BESAAR WOEEEEE"
];

// Updated D3CanvasRefType to match our new D3Canvas implementation
type D3CanvasRefType = {
  handler: {
    refreshData: (messages: Message[]) => void;
  };
};

export default function Page() {
  const [splashText, setSplashText] = useState<string | null>(null);
  
  // State for API data and interaction
  const [elements, setElements] = useState<ElementData[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [targetElement, setTargetElement] = useState("");
  const [filteredElements, setFilteredElements] = useState<ElementData[]>([]);
  const [selectedAlgorithm, setSelectedAlgorithm] = useState<'dfs' | 'bfs' | 'bfs-shortest'>('dfs');
  const [totalRecipes, setTotalRecipes] = useState(5);
  const [isProcessing, setIsProcessing] = useState(false);
  const [results, setResults] = useState<AlgorithmResponse[]>([]);
  const [currentRecipeIndex, setCurrentRecipeIndex] = useState(0);
  const [timingResults, setTimingResults] = useState<TimingInfo[]>([]);


  useIsomorphicLayoutEffect(() => setSplashText(splashTexts[Math.floor(Math.random() * splashTexts.length)]), []);
  
  // Updated D3Canvas ref
  const d3CanvasRef = useRef<D3CanvasRefType>(null);
  
  // Fetch elements from backend on load
  useEffect(() => {
    const getElements = async () => {
      try {
        setLoading(true);
        setError(null);
        const elementData = await fetchElements();
        setElements(elementData);
        setFilteredElements(elementData);
        setLoading(false);
      } catch (error) {
        console.error("Failed to fetch elements:", error);
        setError("Failed to load elements. Please try again later.");
        setLoading(false);
      }
    };
    
    getElements();
  }, []);
  
  // Filter elements as user types
  useEffect(() => {
    if (targetElement.trim()) {
      const filtered = elements.filter(element => 
        element.name.toLowerCase().includes(targetElement.toLowerCase())
      );
      setFilteredElements(filtered.slice(0, 10)); // Limit to 10 results
    } else {
      setFilteredElements([]);
    }
  }, [targetElement, elements]);

  // Handle searching for recipes
  const handleSearch = useCallback(() => {
    if (!targetElement) {
      setError("Please enter a target element");
      return;
    }
    
    // Check if element exists
    const elementExists = elements.some(e => 
      e.name.toLowerCase() === targetElement.toLowerCase()
    );
    
    if (!elementExists) {
      setError(`Element "${targetElement}" not found. Please enter a valid element name.`);
      return;
    }
    
    setIsProcessing(true);
    setResults([]);
    setCurrentRecipeIndex(0);

    const socket = createRecipeWebSocket(
        selectedAlgorithm,
        'target',
        targetElement,
        totalRecipes,
        (result) => {
            setResults(prevResults => {
                // Use prevResults inside this callback where it's available
                const newResults = [...prevResults, result];
				
                toast.success(`Recipe ${newResults.length} found!`);

				if (result.timingInfo) {
					setTimingResults(prev => [
						...prev,
						{
							algorithm: result.timingInfo.algorithm,
							duration: result.timingInfo.duration,
						}
					]);
				}
		
                return newResults;
            });
        },
        (err) => {
            console.error("WebSocket error:", err);
            setError("Error connecting to server. Please try again.");
            setIsProcessing(false);
            toast.error("Connection error. Please try again.");
        }
    );
    
    socket.onclose = () => {
      setIsProcessing(false);
      if (results.length > 0) {
        toast.success(`Found ${results.length} recipes for ${targetElement}`);
      } else {
        toast.error(`No recipes found for ${targetElement}`);
      }
    };
    
    // Safety timeout (30 seconds)
    const timeout = setTimeout(() => {
      if (socket.readyState === WebSocket.OPEN) {
        socket.close();
        if (results.length === 0) {
          setError("Request timed out. Please try again.");
          toast.error("Request timed out");
        }
      }
    }, 30000);
    
    return () => {
      clearTimeout(timeout);
      if (socket.readyState === WebSocket.OPEN) {
        socket.close();
      }
    };
  }, [targetElement, selectedAlgorithm, totalRecipes, elements, results.length]);

  // Updated useEffect to refresh D3Canvas with new recipe data
  useEffect(() => {
    if (!d3CanvasRef.current || results.length === 0) return;
    
    const currentResult = results[currentRecipeIndex];
    if (!currentResult) return;
    
    // Pass the recipe directly to D3Canvas for rendering
    d3CanvasRef.current.handler.refreshData(currentResult.recipe);
    
  }, [results, currentRecipeIndex]);

  return (
    <div className="h-full flex flex-col">
      <div className="container md:h-16 p-4 flex flex-col items-start justify-between space-y-2 sm:flex-row sm:items-center sm:space-y-0">
        <h2 className="text-lg font-semibold">Nuggets</h2>
        <div className="w-full ml-auto flex space-x-2 sm:justify-end">{splashText}</div>
      </div>
      <Separator className="h-[2px]" />
      <div className="container h-full px-4 py-6">
        <div className="h-full grid items-stretch gap-6 md:grid-cols-[250px_1fr]">
          {/* Left Panel - Search Controls */}
          <div className="flex flex-col space-y-4 order-1">
            <Card>
              <CardContent className="pt-6">
                <div className="space-y-6">
                  {error && (
                    <div className="bg-destructive/10 text-destructive p-3 rounded-md">
                      {error}
                    </div>
                  )}
                  
                  <div className="space-y-2">
                    <Label htmlFor="targetElement">Target Element</Label>
                    <div className="relative">
                      <Input 
                        id="targetElement"
                        value={targetElement}
                        onChange={(e) => {
                          setTargetElement(e.target.value);
                          setError(null);
                        }}
                        placeholder="Type an element name..."
                        className="w-full"
                      />
                      
                      {/* Show suggestions */}
                      {filteredElements.length > 0 && targetElement && (
                        <div className="absolute z-10 w-full mt-1 bg-background border rounded-md shadow-lg max-h-60 overflow-auto">
                          {filteredElements.map(element => (
                            <div 
                              key={element.name}
                              className="px-4 py-2 hover:bg-accent cursor-pointer flex items-center gap-2"
                              onClick={() => {
                                setTargetElement(element.name);
                                setFilteredElements([]);
                              }}
                            >
                              <img 
                                src={element.image || getElementImageUrl(element.name)}
                                alt={element.name}
                                className="w-6 h-6"
                                onError={(e) => {
                                  (e.target as HTMLImageElement).src = "https://placehold.co/24x24/orange/white?text=" + element.name.charAt(0);
                                }}
                              />
                              <span>{element.name}</span>
                            </div>
                          ))}
                        </div>
                      )}
                    </div>
                  </div>
                  
                  <div className="space-y-2">
                    <Label>Algorithm</Label>
                    <RadioGroup 
                      value={selectedAlgorithm} 
                      onValueChange={(value: 'dfs' | 'bfs' | 'bfs-shortest') => setSelectedAlgorithm(value)}
                      className="flex flex-col space-y-1"
                    >
                      <div className="flex items-center space-x-2">
                        <RadioGroupItem value="bfs" id="bfs" />
                        <Label htmlFor="bfs">BFS</Label>
                      </div>
                      <div className="flex items-center space-x-2">
                        <RadioGroupItem value="dfs" id="dfs" />
                        <Label htmlFor="dfs">DFS</Label>
                      </div>
                      <div className="flex items-center space-x-2">
                        <RadioGroupItem value="bfs-shortest" id="bfs-shortest" />
                        <Label htmlFor="bfs-shortest">BFS-Shortest</Label>
                      </div>
                    </RadioGroup>
                  </div>
                  
                  <div className="space-y-2">
                    <div className="flex justify-between">
                      <Label htmlFor="totalRecipes">Total Recipes: {totalRecipes}</Label>
                    </div>
                    <Slider
                      id="totalRecipes"
                      min={1}
                      max={10}
                      step={1}
                      value={[totalRecipes]}
                      onValueChange={(value) => setTotalRecipes(value[0])}
                      className="py-4"
                    />
                  </div>
                  
                  <Button 
                    onClick={handleSearch} 
                    className="w-full" 
                    disabled={isProcessing || loading || !targetElement}
                  >
                    {isProcessing ? "Processing..." : "Find Recipes"}
                  </Button>
                  
                  {results.length > 0 && (
                    <div className="flex justify-between items-center space-x-2 mt-4">
                      <Button 
                        onClick={() => setCurrentRecipeIndex(prev => Math.max(0, prev - 1))}
                        disabled={currentRecipeIndex === 0}
                        variant="outline"
                        size="sm"
                        className="flex-1"
                      >
                        Previous
                      </Button>
                      	<span className="text-center font-medium flex-1 flex justify-center">
						            {currentRecipeIndex + 1} of {results.length}
						            </span>
                      <Button 
                        onClick={() => setCurrentRecipeIndex(prev => Math.min(results.length - 1, prev + 1))}
                        disabled={currentRecipeIndex === results.length - 1}
                        variant="outline"
                        size="sm"
                        className="flex-1"
                      >
                        Next
                      </Button>
                    </div>
                  )}
				  {/* Add Timing Display here */}
					{timingResults.length > 0 && (
					<div className="mt-4 p-4 border rounded-md bg-accent/30">
						<h3 className="text-lg font-semibold mb-2">Algorithm Performance</h3>
						<div className="space-y-2">
						{(() => {
							// Get the latest result for each algorithm
							const latestByAlgorithm: Record<string, TimingInfo> = {};
							
							// Process array from end to beginning to get the latest entries
							[...timingResults].reverse().forEach(timing => {
							if (!latestByAlgorithm[timing.algorithm]) {
								latestByAlgorithm[timing.algorithm] = timing;
							}
							});
							
							// Return the latest entries
							return Object.values(latestByAlgorithm).map((timing, idx) => (
							<div key={idx} className="flex justify-between items-center p-2 bg-white/50 rounded">
								<span className="font-medium">{timing.algorithm.toUpperCase()}</span>
								<span className="text-blue-600 font-bold">{timing.duration.toFixed(2)} ms</span>
							</div>
							));
						})()}
						</div>
					</div>
					)}
					{/* Add a button to clear timing data */}
					<Button 
					onClick={() => setTimingResults([])} 
					variant="outline"
					size="sm"
					className="mt-3 w-full"
					>
					Clear Timing Data
					</Button>
                </div>
              </CardContent>
            </Card>
          </div>
          
          {/* D3 Canvas */}
          <div className="order-2 bg-background border-2 border-input rounded-md overflow-hidden">
            {results.length > 0 && (
              <div className="p-2 bg-muted flex justify-between items-center border-b">
                <div className="text-sm font-medium">
                  Recipe {currentRecipeIndex + 1} for {targetElement}
                </div>
                <div className="text-sm text-muted-foreground">
                  Total nodes searched: {results[currentRecipeIndex]?.nodesVisited || 0}
                </div>
              </div>
            )}
            <div className="flex-grow">
              <D3Canvas 
                ref={d3CanvasRef}
                className="w-full h-[600px]" 
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}