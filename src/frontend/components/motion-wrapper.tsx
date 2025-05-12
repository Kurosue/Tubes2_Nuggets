"use client";

import { motion } from "framer-motion";
import { ReactNode } from "react";

interface MotionWrapperProps {
  children: ReactNode;
  delay?: number;
  duration?: number;
  className?: string;
  animation?: "fadeIn" | "fadeInUp" | "fadeInDown" | "slideIn" | "scale" | "bounce";
}

export function MotionWrapper({ 
  children, 
  delay = 0, 
  duration = 0.5, 
  className = "", 
  animation = "fadeIn" 
}: MotionWrapperProps) {
  // Animation variants
  const animations = {
    fadeIn: {
      hidden: { opacity: 0 },
      visible: { opacity: 1, transition: { duration } }
    },
    fadeInUp: {
      hidden: { opacity: 0, y: 20 },
      visible: { opacity: 1, y: 0, transition: { duration } }
    },
    fadeInDown: {
      hidden: { opacity: 0, y: -20 },
      visible: { opacity: 1, y: 0, transition: { duration } }
    },
    slideIn: {
      hidden: { opacity: 0, x: -20 },
      visible: { opacity: 1, x: 0, transition: { duration } }
    },
    scale: {
      hidden: { opacity: 0, scale: 0.8 },
      visible: { opacity: 1, scale: 1, transition: { duration } }
    },
    bounce: {
      hidden: { opacity: 0, scale: 0.8 },
      visible: { 
        opacity: 1, 
        scale: 1, 
        transition: { 
          type: "spring",
          stiffness: 300,
          damping: 10,
          duration 
        } 
      }
    }
  };

  return (
    <motion.div
      className={className}
      initial="hidden"
      animate="visible"
      variants={animations[animation]}
      transition={{ delay }}
    >
      {children}
    </motion.div>
  );
}