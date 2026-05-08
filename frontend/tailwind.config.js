/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{ts,tsx}"],
  theme: {
    extend: {
      boxShadow: {
        pane: "0 1px 2px rgb(15 23 42 / 0.08), 0 18px 60px rgb(15 23 42 / 0.08)"
      },
      colors: {
        ink: "#172033",
        line: "#d9dee7",
        paper: "#f8fafc",
        steel: "#53637a",
        fern: "#276749",
        coral: "#b45309"
      }
    }
  },
  plugins: []
};
