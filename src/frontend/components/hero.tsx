"use client";

import { useState, useEffect } from "react";
import { motion } from "framer-motion";
import { Button } from "@/components/ui/button";
import Link from "next/link";
import { MotionWrapper } from "@/components/motion-wrapper";

interface HeroProps {
  title?: string;
  subtitle?: string;
  onStartClick: () => void;
}

export default function Hero({
  title = "Visual Algorithm Playground",
  subtitle = "Discover All 720 Elements from 4 Basic Elements!",
  onStartClick
}: HeroProps) {
  
  // Animated icons for the hero
  const icons = ["🔍", "🧪", "⚗️", "🌍"];
  const [currentIcon, setCurrentIcon] = useState(0);
  
  // Cycle through icons
  useEffect(() => {
    const interval = setInterval(() => {
      setCurrentIcon(prev => (prev + 1) % icons.length);
    }, 2000);
    
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="w-full bg-gradient-to-b from-primary-50 to-background py-16 px-4 sm:px-6 text-center">
      <div className="max-w-4xl mx-auto">
        <MotionWrapper animation="fadeInDown" delay={0.2}>
          <motion.div 
            className="inline-block mb-6 text-4xl"
            animate={{ 
              rotate: [0, 10, -10, 0],
              scale: [1, 1.1, 1]
            }}
            transition={{ 
              duration: 2,
              repeat: Infinity,
              repeatDelay: 1
            }}
          >
            {icons[currentIcon]}
          </motion.div>
        </MotionWrapper>
        
        <MotionWrapper animation="fadeInUp" delay={0.4}>
          <h1 className="text-4xl md:text-6xl font-bold hero-text-animation mb-6">
            {title}
          </h1>
        </MotionWrapper>
        
        <MotionWrapper animation="fadeInUp" delay={0.6}>
          <p className="text-xl md:text-2xl text-text-muted mb-8">
            {subtitle}
          </p>
        </MotionWrapper>
        
        <MotionWrapper animation="scale" delay={0.8}>
          <Button 
            onClick={onStartClick}
            size="lg" 
            className="bg-primary hover:bg-primary-dark text-white px-8 py-6 text-lg rounded-xl"
          >
            Start Exploring
            <motion.span 
              className="ml-2"
              animate={{ x: [0, 5, 0] }}
              transition={{ repeat: Infinity, duration: 1.5 }}
            >
              ↓
            </motion.span>
          </Button>
        </MotionWrapper>
      </div>
    </div>
  );
}