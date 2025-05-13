// API client for Element API

export interface Element {
  name: string;
  recipes: [string, string][];
  image: string;
  page_url: string;
  tier: number;
}

export interface RecipePath {
  ingredient1: string;
  ingredient2: string;
  result: string;
}

export interface AlgorithmResponse {
  recipePath: RecipePath[];
  nodesVisited: number;
  duration: number; // milliseconds
}

const API_BASE_URL = new URL("/api", process.env.NEXT_PUBLIC_BACKEND).href;

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
  const url = new URL("/api/find-recipe", process.env.NEXT_PUBLIC_BACKEND);
  url.searchParams.set("algorithm", algorithm);
  url.searchParams.set("direction", direction);
  url.searchParams.set("target", targetElement);
  url.searchParams.set("count", `${count}`);
  const socket = new WebSocket(url);
  const startTime = performance.now();
  
  socket.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data);
      if(typeof data.duration != "number")
        data.duration = performance.now() - startTime;
      onMessage(data);
    } catch (err) {
      console.error('Error parsing WebSocket message:', err);
    }
  };
  
  socket.onerror = onError;
  
  return socket;
}
