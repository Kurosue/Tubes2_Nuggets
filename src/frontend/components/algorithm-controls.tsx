"use client";

import { useState, useEffect } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { Slider } from "@/components/ui/slider";
import { Card, CardContent } from "@/components/ui/card";
import { Element, TimingInfo } from "@/lib/api/service";

interface AlgorithmControlsProps {
  targetElement: string;
  setTargetElement: (value: string) => void;
  selectedAlgorithm: 'dfs' | 'bfs' | 'bfs-shortest';
  setSelectedAlgorithm: (value: 'dfs' | 'bfs' | 'bfs-shortest') => void;
  totalRecipes: number;
  setTotalRecipes: (value: number) => void;
  error: string | null;
  isProcessing: boolean;
  handleSearch: () => void;
  filteredElements: Element[];
  setFilteredElements: (elements: Element[]) => void;
  currentRecipeIndex: number;
  setCurrentRecipeIndex: React.Dispatch<React.SetStateAction<number>>; // Updated type definition
  resultsLength: number;
  timingResults: TimingInfo[];
  setTimingResults: (value: TimingInfo[]) => void;
}

export default function AlgorithmControls({
  targetElement,
  setTargetElement,
  selectedAlgorithm,
  setSelectedAlgorithm,
  totalRecipes,
  setTotalRecipes,
  error,
  isProcessing,
  handleSearch,
  filteredElements,
  setFilteredElements,
  currentRecipeIndex,
  setCurrentRecipeIndex,
  resultsLength,
  timingResults,
  setTimingResults
}: AlgorithmControlsProps) {
  
  // Element search animation state
  const [isAnimating, setIsAnimating] = useState(false);
  
  // Trigger animation on element input
  useEffect(() => {
    if (targetElement.length > 0) {
      setIsAnimating(true);
      const timer = setTimeout(() => setIsAnimating(false), 200);
      return () => clearTimeout(timer);
    }
  }, [targetElement]);
  
  // Handle previous recipe navigation
  const handlePrevRecipe = () => {
    setCurrentRecipeIndex((prev) => Math.max(0, prev - 1));
  };
  
  // Handle next recipe navigation
  const handleNextRecipe = () => {
    setCurrentRecipeIndex((prev) => Math.min(resultsLength - 1, prev + 1));
  };
  
  return (
    <Card className="shadow-card text-text-muted">
      <CardContent className="pt-6">
        <div className="space-y-6">
          {error && (
            <motion.div 
              initial={{ opacity: 0, y: -10 }}
              animate={{ opacity: 1, y: 0 }}
              className="bg-error/10 text-error p-3 rounded-md"
            >
              {error}
            </motion.div>
          )}
          
          <div className="space-y-2">
            <Label htmlFor="targetElement">Target Element</Label>
            <motion.div 
              className="relative"
              animate={isAnimating ? {
                boxShadow: [
                  "0 0 0 rgba(136, 68, 221, 0)",
                  "0 0 8px rgba(136, 68, 221, 0.6)",
                  "0 0 0 rgba(136, 68, 221, 0)"
                ]
              } : {}}
              transition={{ duration: 0.5 }}
            >
              <Input 
                id="targetElement"
                value={targetElement}
                onChange={(e) => setTargetElement(e.target.value)}
                placeholder="Type an element name..."
                className="w-full border-primary/20 focus:border-primary focus:ring-2 focus:ring-primary/20"
              />
              
              <AnimatePresence>
                {filteredElements.length > 0 && targetElement && (
                  <motion.div 
                    initial={{ opacity: 0, y: -20 }}
                    animate={{ opacity: 1, y: 0 }}
                    exit={{ opacity: 0, y: -20 }}
                    className="absolute z-10 w-full mt-1 bg-background-card border border-primary/20 rounded-md shadow-lg max-h-60 overflow-auto"
                  >
                    {filteredElements.map((element, index) => (
                      <motion.div 
                        key={element.name}
                        initial={{ opacity: 0, y: -5 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ delay: index * 0.03 }}
                        className="px-4 py-2 hover:bg-primary-50 cursor-pointer flex items-center gap-2"
                        onClick={() => {
                          setTargetElement(element.name);
                          setFilteredElements([]);
                        }}
                      >
                        <img 
                          src={new URL(element.image, process.env.NEXT_PUBLIC_BACKEND).href}
                          alt={element.name}
                          className="w-6 h-6"
                          onError={(e) => {
                            (e.target as HTMLImageElement).src = "https://placehold.co/24x24/8844dd/white?text=" + element.name.charAt(0);
                          }}
                        />
                        <span>{element.name}</span>
                      </motion.div>
                    ))}
                  </motion.div>
                )}
              </AnimatePresence>
            </motion.div>
          </div>
          
          <div className="space-y-2">
            <Label>Algorithm</Label>
            <RadioGroup 
              value={selectedAlgorithm} 
              onValueChange={(value: 'dfs' | 'bfs' | 'bfs-shortest') => setSelectedAlgorithm(value)}
              className="flex flex-col space-y-1"
            >
              <motion.div 
                className="flex items-center space-x-2"
                whileHover={{ x: 2 }}
              >
                <RadioGroupItem value="bfs" id="bfs" />
                <Label 
                  htmlFor="bfs" 
                  className={selectedAlgorithm === 'bfs' ? 'text-algorithm-bfs font-medium' : ''}
                >
                  BFS
                </Label>
              </motion.div>
              <motion.div 
                className="flex items-center space-x-2"
                whileHover={{ x: 2 }}
              >
                <RadioGroupItem value="dfs" id="dfs" />
                <Label 
                  htmlFor="dfs"
                  className={selectedAlgorithm === 'dfs' ? 'text-algorithm-dfs font-medium' : ''}
                >
                  DFS
                </Label>
              </motion.div>
              <motion.div 
                className="flex items-center space-x-2"
                whileHover={{ x: 2 }}
              >
                {/* <RadioGroupItem value="bfs-shortest" id="bfs-shortest" />
                <Label 
                  htmlFor="bfs-shortest"
                  className={selectedAlgorithm === 'bfs-shortest' ? 'text-algorithm-bfsShort font-medium' : ''}
                >
                  BFS-Shortest
                </Label> */}
              </motion.div>
            </RadioGroup>
          </div>

          <span className="text-muted-foreground text-xs">For more than one recipe, multithreading will be used</span>
          
          <div className="space-y-2">
            <div className="flex justify-between items-center">
              <Label htmlFor="totalRecipes">Total Recipes</Label>
              <div className="flex items-center space-x-2">
                <span className="text-primary font-bold">{totalRecipes}</span>
                <span className="text-muted-foreground text-sm">/ 20</span>
              </div>
            </div>
            <Slider
              id="totalRecipes"
              min={1}
              max={20}
              step={1}
              value={[totalRecipes]}
              onValueChange={(value) => setTotalRecipes(value[0])}
              className="py-4"
            />
          </div>
          
          <motion.div
            whileHover={{ scale: 1.02 }}
            whileTap={{ scale: 0.98 }}
          >
            <Button 
              onClick={handleSearch} 
              className="w-full bg-primary hover:bg-primary-dark text-white"
              disabled={isProcessing || !targetElement}
            >
              {isProcessing ? (
                <>
                  <svg className="animate-spin -ml-1 mr-3 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  Processing...
                </>
              ) : (
                "Find Recipes"
              )}
            </Button>
          </motion.div>
          
          {resultsLength > 0 && (
            <div className="flex justify-between items-center space-x-2 mt-4">
              <motion.div whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.95 }}>
                <Button 
                  onClick={handlePrevRecipe}
                  disabled={currentRecipeIndex === 0}
                  variant="outline"
                  size="sm"
                  className="flex-1 border-primary/20"
                >
                  Previous
                </Button>
              </motion.div>
              <span className="text-center font-medium flex-1 flex justify-center">
                {currentRecipeIndex + 1} of {resultsLength}
              </span>
              <motion.div whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.95 }}>
                <Button 
                  onClick={handleNextRecipe}
                  disabled={currentRecipeIndex === resultsLength - 1}
                  variant="outline"
                  size="sm"
                  className="flex-1 border-primary/20"
                >
                  Next
                </Button>
              </motion.div>
            </div>
          )}
          
          {/* Timing Display */}
          {timingResults.length > 0 && (
            <motion.div 
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              className="mt-4 p-4 border rounded-md bg-background-muted"
            >
              <h3 className="text-lg font-semibold mb-2 text-primary">Algorithm Performance</h3>
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
                    <motion.div 
                      key={idx} 
                      className="flex justify-between items-center p-2 bg-background-card rounded"
                      initial={{ opacity: 0, x: -10 }}
                      animate={{ opacity: 1, x: 0 }}
                      transition={{ delay: idx * 0.1 }}
                    >
                      <span className="font-medium">
                        {timing.algorithm.toUpperCase()}
                      </span>
                      <span className={`font-bold ${
                        timing.algorithm === 'dfs' 
                          ? 'text-algorithm-dfs' 
                          : timing.algorithm === 'bfs' 
                            ? 'text-algorithm-bfs' 
                            : 'text-algorithm-bfsShort'
                      }`}>
                        {timing.duration.toFixed(2)} ms
                      </span>
                    </motion.div>
                  ));
                })()}
              </div>
            </motion.div>
          )}
          
          {timingResults.length > 0 && (
            <motion.div whileHover={{ scale: 1.02 }} whileTap={{ scale: 0.98 }}>
              <Button 
                onClick={() => setTimingResults([])} 
                variant="outline"
                size="sm"
                className="mt-3 w-full border-primary/20 text-primary hover:bg-primary-50"
              >
                Clear Timing Data
              </Button>
            </motion.div>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
