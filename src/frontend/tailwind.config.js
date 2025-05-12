/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
    "./pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
    "./lib/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      colors: {
        background: {
          DEFAULT: '#f8f7fe',
          card: '#ffffff',
          muted: '#f3f3f7',
          dark: '#282a36',
        },
        primary: {
          DEFAULT: '#8844dd',
          light: '#b06efa',
          dark: '#6633aa',
          50: '#f3e8ff',
          100: '#e5d1fc',
          200: '#d0b2f9',
          300: '#b991f5',
          400: '#a06bee',
          500: '#8844dd',
          600: '#7737c5',
          700: '#632ba7',
          800: '#502287',
          900: '#3c1a66',
        },
        secondary: {
          DEFAULT: '#44ddbb',
          light: '#7aefe0',
          dark: '#2eaa8e',
        },
        text: {
          DEFAULT: '#333344',
          muted: '#6e6e8e',
          light: '#9494a4',
          inverse: '#ffffff',
        },
        element: {
          water: '#55acee',
          fire: '#ff5555', 
          earth: '#bd93f9',
          air: '#f8f8f2',
          combined: '#50fa7b',
        },
        algorithm: {
          bfs: '#8be9fd',
          dfs: '#ff79c6',
          bfsShort: '#50fa7b',
        },
      },
      boxShadow: {
        'card': '0 4px 6px -1px rgba(136, 68, 221, 0.1), 0 2px 4px -1px rgba(136, 68, 221, 0.06)',
        'card-hover': '0 10px 15px -3px rgba(136, 68, 221, 0.2), 0 4px 6px -2px rgba(136, 68, 221, 0.1)',
      },
    },
  },
  plugins: [],
}