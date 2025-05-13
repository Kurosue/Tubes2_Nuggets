"use client";

import React, { useCallback, useEffect, useRef, useState } from "react";
import { Separator } from "@/components/ui/separator";
import { Button } from "@/components/ui/button";
import { useIsomorphicLayoutEffect } from "@/lib/utils";
import { motion, AnimatePresence } from "framer-motion";

// Import toast for notifications
import { toast } from "sonner";

// Import components
import Hero from "@/components/hero";
import Features from "@/components/features";
import AlgorithmControls, { TimingInfo } from "@/components/algorithm-controls";
import D3Visualization from "@/components/d3-visualization";

// Import API services
import { 
  fetchElements, 
  createRecipeWebSocket, 
  Element as ElementData, 
  RecipePath, 
  AlgorithmResponse
} from "@/lib/api/service";
import { Slider } from "@/components/ui/slider";

// Splash texts for random selection
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

// D3Canvas reference type
type D3CanvasRefType = {
  handler: {
    refreshData: (messages: RecipePath[]) => void;
  };
};

export default function Page() {
  // State for page sections
  const [showHero, setShowHero] = useState<boolean>(true);
  const [splashText, setSplashText] = useState<string | null>(null);
  
  // Refs
  const playgroundRef = useRef<HTMLDivElement>(null);
  const d3CanvasRef = useRef<D3CanvasRefType>(null);
  
  // State for API data and interaction
  const [elements, setElements] = useState<ElementData[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [targetElement, setTargetElement] = useState("");
  const [filteredElements, setFilteredElements] = useState<ElementData[]>([]);
  const [selectedAlgorithm, setSelectedAlgorithm] = useState<'dfs' | 'bfs' | 'bfs-shortest'>('dfs');
  const [totalRecipes, setTotalRecipes] = useState(5);
  const [isProcessing, setIsProcessing] = useState(false);
  const [currentRecipeIndex, setCurrentRecipeIndex] = useState(0);
  const [timingResults, setTimingResults] = useState<TimingInfo[]>([]);
  const [visualizationStep, setVisualizationStep] = useState(Infinity);
  const resultsRef = useRef<AlgorithmResponse[]>([]);

  // Select a random splash text on initial render
  useIsomorphicLayoutEffect(() => {
    setSplashText(splashTexts[Math.floor(Math.random() * splashTexts.length)]);
  }, []);
  
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

  // Scroll to playground when hero button is clicked
  const handleStartExploring = () => {
    setShowHero(false);
    setTimeout(() => {
      if (playgroundRef.current) {
        playgroundRef.current.scrollIntoView({ behavior: 'smooth' });
      }
    }, 100);
    setShowHero(true);
  };

  // Handle searching for recipes
  const handleSearch = useCallback(() => {
    if (!targetElement) {
      setError("Please enter a target element");
      return;
    }
    
    // Check if element exists
    const elementExists = elements.some(e => 
      e.name.toLowerCase() == targetElement.toLowerCase()
    );
    
    if (!elementExists) {
      setError(`Element "${targetElement}" not found. Please enter a valid element name.`);
      return;
    }
    
    resultsRef.current = [];
    // Trigger rerender automatically below:
    setVisualizationStep(Infinity);
    setIsProcessing(true);
    setCurrentRecipeIndex(0);
    setError(null);

    const socket = createRecipeWebSocket(
        selectedAlgorithm,
        'target',
        targetElement,
        totalRecipes,
        (result) => {
          result.recipePath = (() => {
            const m = new Map<string, RecipePath>();
            for(const p of result.recipePath)
              m.set(p.result, p);
            return [...m.values()];
          })();
          resultsRef.current = [...resultsRef.current, result];
          // Trigger rerender automatically below:
          setTimingResults(prev => [
              ...prev,
              {
                  algorithm: selectedAlgorithm,
                  duration: result.duration,
              }
          ]);
          // toast.success(`Found recipe for ${targetElement}`, {
          //     style: {
          //         background: '#4CAF50', // green background
          //         color: '#fff'
          //     }
          // });
        },
        (err) => {
            console.error("WebSocket error:", err);
            setError("Error connecting to server. Please try again.");
            setIsProcessing(false);
            toast.error("Connection error. Please try again.", {
                style: {
                    background: '#f44336', // red background
                    color: '#fff'
                }
            });
        }
    );

    socket.onclose = () => {
      setIsProcessing(false);

      const totalRecipesFound = resultsRef.current.length;

      if (totalRecipesFound > 0) {
        toast.success(`Found ${totalRecipesFound} recipes for ${targetElement}`, {
          style: {
            background: '#4CAF50',
            color: '#fff'
          }
        });
      } else {
        toast.error(`No recipes found for ${targetElement}`, {
          style: {
            background: '#f44336', // red background
            color: '#fff'
          }
        });
      }
    };

    // Safety timeout (30 seconds)
    const timeout = setTimeout(() => {
      if (socket.readyState == WebSocket.OPEN) {
        socket.close();
        if (resultsRef.current.length == 0) {
          setError("Request timed out. Please try again.");
          toast.error("Request timed out", {
            style: {
              background: '#f44336', // red background
              color: '#fff'
            }
          });
        }
      }
    }, 30000);
    
    return () => {
      clearTimeout(timeout);
      if (socket.readyState == WebSocket.OPEN) {
        socket.close();
      }
    };
  }, [targetElement, selectedAlgorithm, totalRecipes, elements]);

  // Update D3Canvas with new recipe data
  useEffect(() => {
    if (!d3CanvasRef.current || resultsRef.current.length == 0) return;
    
    const currentResult = resultsRef.current[currentRecipeIndex];
    if (!currentResult) return;
    
    // Pass the recipe directly to D3Canvas for rendering
    d3CanvasRef.current.handler.refreshData(currentResult.recipePath.slice(0, visualizationStep));
    
  }, [resultsRef.current, visualizationStep, currentRecipeIndex]);

  return (
    <div className="min-h-screen flex flex-col bg-background">
      {/* Header */}
      <header className="border-b border-primary/10">
        <div className="container md:h-16 p-4 flex items-center justify-between mx-auto">
          <motion.h1 
            className="text-2xl font-bold text-primary"
            initial={{ opacity: 0, x: -20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.5 }}
          >
            Nuggets
          </motion.h1>
          <motion.span 
            className="text-xs px-2 py-1 bg-primary-50 text-primary rounded-full"
            initial={{ opacity: 0, x: 20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.5, delay: 0.2 }}
          >
            Beta
          </motion.span>
        </div>
      </header>
      
      {/* Main content */}
      <main className="flex-grow">
        <AnimatePresence>
          {showHero && (
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0, height: 0 }}
              transition={{ duration: 0.5 }}
            >
              <Hero 
                title="Little Alchemy 2 Recipe Explorer" 
                subtitle={splashText || "Discover All 720 Elements from 4 Basic Elements!"}
                onStartClick={handleStartExploring}
              />
              <Features />
            </motion.div>
          )}
        </AnimatePresence>
        
        {/* Algorithm Playground */}
        <div 
          ref={playgroundRef}
          className="py-8 px-4"
        >
          <motion.div 
            className="container max-w-6xl mx-auto"
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5, delay: showHero ? 0 : 0.3 }}
          >
            <h2 className="text-3xl font-bold text-center mb-2 text-primary">Algorithm Playground</h2>
            <p className="text-center text-text-muted mb-8">Explore element combinations using different algorithms</p>
            
            <div className="grid grid-cols-1 md:grid-cols-[300px_1fr] gap-6">
              {/* Left panel - Controls */}
              <AlgorithmControls
                targetElement={targetElement}
                setTargetElement={setTargetElement}
                selectedAlgorithm={selectedAlgorithm}
                setSelectedAlgorithm={setSelectedAlgorithm}
                totalRecipes={totalRecipes}
                setTotalRecipes={setTotalRecipes}
                error={error}
                isProcessing={isProcessing}
                handleSearch={handleSearch}
                filteredElements={filteredElements}
                setFilteredElements={setFilteredElements}
                currentRecipeIndex={currentRecipeIndex}
                setCurrentRecipeIndex={v => { setCurrentRecipeIndex(v); setVisualizationStep(Infinity); }}
                resultsLength={resultsRef.current.length}
                timingResults={timingResults}
                setTimingResults={setTimingResults}
              />
              
              {/* Right panel - Visualization */}
              <div className="bg-background-card border-2 border-primary/10 rounded-xl overflow-hidden shadow-card">
                {resultsRef.current.length > 0 && (
                  <div className="p-3 bg-background-muted flex justify-between items-center border-b border-primary/10">
                    <div className="text-sm font-medium text-primary">
                      Recipe {currentRecipeIndex + 1} for {targetElement}
                    </div>
                    <div className="text-sm text-text-muted">
                      Total nodes searched: {resultsRef.current[currentRecipeIndex]?.nodesVisited || 0}
                    </div>
                  </div>
                )}
                <div className="h-[600px]">
                  <D3Visualization ref={d3CanvasRef} className="w-full h-full" />
                </div>
                {resultsRef.current.length > 0 && (
                  <div className="p-3 bg-background-muted flex justify-between items-center border-b border-primary/10">
                    <Slider
                      min={1}
                      max={resultsRef.current[currentRecipeIndex]?.recipePath.length || 1}
                      step={1}
                      value={[Math.min(visualizationStep, resultsRef.current[currentRecipeIndex]?.recipePath.length || 1)]}
                      onValueChange={(value) => setVisualizationStep(value[0])}
                      className="py-4"
                    />
                  </div>
                )}
              </div>
            </div>
          </motion.div>
        </div>
      </main>
      
      {/* Footer */}
      <footer className="mt-auto py-6 border-t border-primary/10">
        <div className="container text-center text-text-muted mx-auto">
          <p className="mb-2">© 2025 Nuggets - Visual Algorithm Playground</p>
          <div className="flex flex-col items-center justify-center">
            <p className="font-medium text-primary mb-1 mx">Contributors</p>
            <div className="flex flex-wrap justify-center gap-2 my-1">
              <a 
                href="https://github.com/Kurosue" 
                target="_blank" 
                rel="noopener noreferrer"
                className="px-2 py-1 bg-primary/10 rounded-full text-xs hover:bg-primary/20 transition-colors"
              >
                Muhammad Aditya Rahmadeni - 13523028
              </a>
              <a
                href="https://github.com/NadhifRadityo"
                target="_blank"
                rel="noopener noreferrer"
                className="px-2 py-1 bg-primary/10 rounded-full text-xs hover:bg-primary/20 transition-colors"
              >
                Nadhif Radityo Nugroho - 13523045
              </a>
              <a
                href="https://github.com/ryonlunar"
                target="_blank"
                rel="noopener noreferrer"
                className="px-2 py-1 bg-primary/10 rounded-full text-xs hover:bg-primary/20 transition-colors"
              >
                Adhimas Aryo Bimo - 1323052
              </a>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
}
