// API client for Element API

export interface TimingInfo {
  algorithm: string;
  duration: number; // milliseconds
}

export interface Element {
  name: string;
  recipes: string[];
  image: string;
  page_url: string;
  tier: number;
  parsed_recipes: string[][];
}

export interface Message {
  Ingredient1: string;
  Ingredient2: string;
  Result: string;
  Depth: number;
}

export interface AlgorithmResponse {
  recipe: Message[];
  nodesVisited: number;
  recipeIndex: number;
  totalRecipes: number;
  timingInfo: TimingInfo;
}

const API_BASE_URL = 'http://localhost:8888/api';

// Fetch all elements
export async function fetchElements(): Promise<Element[]> {
  const response = await fetch(`${API_BASE_URL}/elements`);
  if (!response.ok) {
    throw new Error('Failed to fetch elements');
  }
  const data = await response.json();
  return data.elements;
}

// Create a WebSocket connection to find recipes
export function createRecipeWebSocket(
  algorithm: 'dfs' | 'bfs' | 'bfs-shortest',
  direction: 'target',
  targetElement: string,
  count: number,
  onMessage: (result: AlgorithmResponse) => void,
  onError: (error: Event) => void
): WebSocket {
  const url = `ws://localhost:8888/api/find-recipe?algorithm=${algorithm}&direction=${direction}&target=${encodeURIComponent(targetElement)}&count=${count}`;
  const socket = new WebSocket(url);

  const startTime = performance.now();
  
  socket.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data);
      const endTime = performance.now();
      const duration = endTime - startTime;
      if (!data.timingInfo) {
      data.timingInfo = {
        algorithm: algorithm,
        duration: duration
      };
    } else if (typeof data.timingInfo !== 'object') {
      // If timingInfo is not an object (maybe it's a number from backend)
      data.timingInfo = {
        algorithm: algorithm,
        duration: typeof data.timingInfo === 'number' ? data.timingInfo : (duration)
      };
    }
      onMessage(data);
    } catch (err) {
      console.error('Error parsing WebSocket message:', err);
    }
  };
  
  socket.onerror = onError;
  
  return socket;
}