import tailwindcss from "@tailwindcss/vite";
import { tanstackRouter } from "@tanstack/router-plugin/vite";
import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";

// BASE_PATH (mise [env] / Docker ARG) is the single source; the router basepath and API prefix
// read it back as import.meta.env.BASE_URL, the server reads the same env var at runtime.
const base = `${process.env.BASE_PATH ?? "/mmtc34"}/`;

export default defineConfig({
  base,
  plugins: [tanstackRouter({ target: "react", autoCodeSplitting: true }), react(), tailwindcss()],
  build: { outDir: "../server/internal/spa/dist", emptyOutDir: true },
  server: { proxy: { [`${base}api`]: "http://localhost:8080" } },
});
