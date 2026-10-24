# Self-host on the homelab: Go + Gin + SQLite in Docker

The original handoff targeted Cloudflare Workers + Hono + D1 + R2 on free tiers. We instead run the app on the owner's homelab: a Go (Gin) service with SQLite, everything in Docker, sitting behind the owner's existing internal Go (Gin) gateway. This puts the app on infrastructure and a language the owner already runs, and removes the Workers constraints (10 ms CPU, D1 daily caps, presigned R2 uploads) that were shaping the design.

## Consequences

- No Hono RPC end-to-end types; the frontend gets API types another way.
- Slips and images live on the homelab's disk next to the SQLite file. Backing them up is the infrastructure's job, not the app's.
- Public ingress is Cloudflare Tunnel into the gateway on one shared hostname, routed by path, with no subdomain per app. The app must therefore work under a path prefix.
