import tailwindcss from "@tailwindcss/vite";
import { tanstackRouter } from "@tanstack/router-plugin/vite";
import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";

// BASE_PATH is the single source: the router basepath and API prefix read it back as
// import.meta.env.BASE_URL, and the server reads the same env var at runtime. "/" (or empty)
// builds for the root of a dedicated subdomain.
const base = `${(process.env.BASE_PATH ?? "/mmtc34").replace(/\/+$/, "")}/`;

export default defineConfig({
  base,
  server: {
    port: 3001,
    proxy: { [`${base}api`]: "http://localhost:8080" },
  },
  resolve: {
    tsconfigPaths: true,
  },
  plugins: [
    tailwindcss(),
    tanstackRouter({
      target: "react",
      autoCodeSplitting: true,
    }),
    react(),
  ],
  build: { outDir: "../server/internal/spa/dist", emptyOutDir: true },
});
