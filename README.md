# Class Campaign

LINE MINI App for running class money-collecting campaigns. See `GLOSSARY.md` for terms and `docs/architecture.md` for the design.

The web app (`apps/web`) is a Better-T-Stack scaffold (React, TanStack Router, Tailwind, Biome) with shared shadcn/ui primitives in `packages/ui`. The server (`apps/server`) is Go and embeds the built web app.

## Commands

- `bun install`: install workspace dependencies
- `bun run dev:web`: Vite dev server on port 3001, proxies `${BASE_PATH}/api` to `localhost:8080`
- `bun run build`: build the web app into `apps/server/internal/spa/dist`
- `bun run check-types`: generate routes, then `tsc --noEmit` (needs `apps/web/src/api/types.gen.ts` from tygo)
- `bun run check`: Biome format and lint

`BASE_PATH` (default `/mmtc34`) is read by Vite at build time and by the server at run time.

## Adding shared UI primitives

```bash
bunx shadcn@latest add dialog -c packages/ui
```

```tsx
import { Button } from "@mmtc34/ui/components/button";
```
