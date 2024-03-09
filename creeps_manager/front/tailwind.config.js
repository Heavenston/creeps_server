/** @type {import('tailwindcss').Config} */
import * as colors from "tailwindcss/colors.js";

export default {
  content: ["./src/**/*.ts", "../**/*.templ"],
  theme: {
    extend: {
      colors: {
        "dark": {
          one: "#101010",
          two: "#191919",
          three: "#262626",
        },
      },
    },
  },
  plugins: [],
}

