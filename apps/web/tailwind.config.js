/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{ts,tsx}"],
  theme: {
    extend: {
      colors: {
        ink: "#1d2521",
        mist: "#eef3f1",
        leaf: "#256f5b",
        coral: "#d95f43",
        gold: "#c99726",
      },
      boxShadow: {
        panel: "0 14px 40px rgba(29, 37, 33, 0.12)",
      },
    },
  },
  plugins: [],
};
