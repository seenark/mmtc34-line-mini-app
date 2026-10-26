# Self-host on the homelab: Go + Gin + SQLite in Docker

The original handoff targeted Cloudflare Workers + Hono + D1 + R2 on free tiers. We instead run the app on the owner's homelab: a Go (Gin) service with SQLite, everything in Docker, sitting behind the owner's existing internal Go (Gin) gateway. This puts the app on infrastructure and a language the owner already runs, and removes the Workers constraints (10 ms CPU, D1 daily caps, presigned R2 uploads) that were shaping the design.

## Consequences

- No Hono RPC end-to-end types; the frontend gets API types another way.
- Slips and images live on the homelab's disk next to the SQLite file. Backing them up is the infrastructure's job, not the app's.
- A dedicated subdomain does not need the path prefix: `BASE_PATH=/` serves the API at `/api` and the web app at the root, with the session cookie scoped to `/`. The prefix stays the default for the shared hostname.
- Public ingress is Cloudflare Tunnel. On the shared gateway hostname, apps are routed by path, so the app works under a path prefix; on its own subdomain it runs at the root (previous bullet).
