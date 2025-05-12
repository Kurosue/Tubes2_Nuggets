"use client";

import { motion } from "framer-motion";
import { MotionWrapper } from "@/components/motion-wrapper";

interface FeatureCardProps {
  title: string;
  description: string;
  icon: string;
  delay: number;
}

function FeatureCard({ title, description, icon, delay }: FeatureCardProps) {
  return (
    <MotionWrapper animation="fadeInUp" delay={delay}>
      <div className="feature-card">
        <div className="icon text-primary">{icon}</div>
        <h3 className="text-xl font-bold mb-2">{title}</h3>
        <p className="text-text-muted">{description}</p>
      </div>
    </MotionWrapper>
  );
}

export default function Features() {
  const features = [
    {
      title: "Multiple Algorithms",
      description: "Compare BFS, DFS, and BFS-Shortest to find the most efficient path to your target element.",
      icon: "🧠",
    },
    {
      title: "Visual Tree Representation",
      description: "See the recipe tree unfold step by step with our interactive D3.js visualization.",
      icon: "🌳",
    },
    {
      title: "Performance Metrics",
      description: "Track algorithm performance with detailed timing and node traversal information.",
      icon: "⚡",
    },
    {
      title: "Multithreading Support",
      description: "Utilize multiple threads to speed up the search process and enhance performance.",
      icon: "🚀",
    }
  ];

  return (
    <section className="py-16 px-4">
      <div className="max-w-6xl mx-auto">
        <MotionWrapper animation="fadeInUp" delay={0.2}>
          <h2 className="text-3xl md:text-4xl font-bold text-center text-primary mb-4">Features</h2>
        </MotionWrapper>
        
        <MotionWrapper animation="fadeInUp" delay={0.3}>
          <p className="text-center text-text-muted mb-12 max-w-3xl mx-auto">
            Explore the power of algorithms through our interactive visualization platform
          </p>
        </MotionWrapper>
        
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 text-text-muted">
          {features.map((feature, index) => (
            <FeatureCard
              key={feature.title}
              title={feature.title}
              description={feature.description}
              icon={feature.icon}
              delay={0.4 + index * 0.1}
            />
          ))}
        </div>
      </div>
    </section>
  );
}