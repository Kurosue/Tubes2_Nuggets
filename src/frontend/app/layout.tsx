import "./globals.scss";
import { GeistSans } from "geist/font/sans";
import { Toaster } from "@/components/ui/sonner";
import { icons } from "lucide-react";

export const metadata = {
  title: "Nuggets - Visual Algorithm Playground",
  description: "Explore and visualize 720 elements using algorithms like DFS and BFS",
  icons: {
    icon: "/logo-circle.png",  
  },
};

export default function Layout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" className={GeistSans.className}>
      <body>
        {children}
        <Toaster position="top-right" />
      </body>
    </html>
  );
}