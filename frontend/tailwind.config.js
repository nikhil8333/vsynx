/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        trust: {
          legitimate: '#10b981',
          suspicious: '#f59e0b',
          malicious: '#ef4444',
          unknown: '#6b7280',
        },
      },
    },
  },
  plugins: [],
}
