import { clsx, type ClassValue } from "clsx"
import { useEffect, useLayoutEffect } from "react";
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

const useIsomorphicLayoutEffect = typeof window !== "undefined" ? useLayoutEffect : useEffect;
export { useIsomorphicLayoutEffect };
