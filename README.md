# MMTC Events

LINE MINI App where one school class runs money-collecting Campaigns (jersey orders, donations, parties). Anyone who signs in with LINE and fills in a profile is a Member at once; there is no approval step. Members order like in an online shop, pay by transfer using the Campaign's QR image and bank account, and upload a transfer Slip (one or more images) per Order. Admins review Slips, upload a Slip on a Member's behalf and add or remove its images (on the admin Order page, where they can review it at once), record cash and Adjustments, and export factory, shipping and pickup lists.

Terms are defined in `GLOSSARY.md`. The design is in `docs/architecture.md`, and the decisions behind it in `docs/adr/`.

## Layout

- `apps/server` is the Go server (Gin, SQLite, sqlc, goose). It serves the API under `${BASE_PATH}/api` and embeds the built web app for every other path under `${BASE_PATH}`. With `BASE_PATH=/` the API is at `/api` and the web app owns the whole host.
- `apps/web` is the React SPA (TanStack Router and Query, Tailwind, Biome), scaffolded with Better-T-Stack. `packages/ui` holds the shared shadcn/ui primitives.

## Toolchain

`mise install` installs Go, Bun, sqlc and tygo at the versions pinned in `mise.toml`. Then run `bun install` at the repo root.

## Develop

```bash
mise run seed     # create ./apps/server/data/app.db with the football 2026 Campaign; its images are embedded in the binary and copied into data/files/images when absent
mise run dev      # Go server on :8080 with DEV_LOGIN on, Vite on :3001
```

Open `http://localhost:3001/mmtc34/`. With `DEV_LOGIN` on, the login page takes any LINE userId instead of LIFF. The dev server also enables the Super Admin as `admin` / `dev-super-admin-pw`; sign in at `/mmtc34/admin-login` and promote Members to Admin.

Admins who signed in with LINE switch between Admin Mode and Member Mode with the toggle in the header. The choice is remembered on the device (`localStorage` key `mode`) and the server never sees it. The Super Admin has no Orders of its own and is always in admin.

Each Campaign picks its bank from a fixed list of Thai banks and can carry an uploaded QR image (for example from K PLUS). The app does not generate a payment QR: bank-account Thai QR (EMVCo tag 29 sub-tag 04) was tried and failed in real banking apps, so Members scan the owner's uploaded QR and type the amount shown on their Order.

| Task | What it runs |
| --- | --- |
| `mise run gen` | `sqlc generate` and `tygo generate` (writes `apps/web/src/api/types.gen.ts`) |
| `mise run test` | `go test ./...` in `apps/server`, the HTTP-level test suite |
| `mise run check` | `go vet`, TypeScript typecheck and Biome |
| `mise run build` | web build, then a static `bin/mmtc34` binary |
| `mise run docker` | `docker build -t mmtc34 .` |

## Configuration

| Variable | Required | Meaning |
| --- | --- | --- |
| `BASE_PATH` | no, default `/mmtc34` | URL prefix, or `/` (or empty) to serve at the root of a dedicated subdomain. The default applies only when the variable is unset. Baked into the web build (Docker build arg) and read by the server; the server refuses to start if they differ. |
| `LIFF_ID` | yes, unless `DEV_LOGIN` | LIFF app ID served to the browser. |
| `LINE_CHANNEL_ID` | yes, unless `DEV_LOGIN` | Channel ID used to verify LIFF ID tokens with LINE. |
| `SESSION_SECRET` | yes | At least 32 bytes; signs the session cookie. |
| `SUPER_ADMIN_USERNAME` | no, default `admin` | Username of the Super Admin account. |
| `SUPER_ADMIN_PASSWORD` | yes, to bootstrap | At least 12 characters. Enables the Super Admin login at `${BASE_PATH}/admin-login`; unset disables it. |
| `DEV_LOGIN` | no, default off | Enables login by raw userId and non-Secure cookies. Never set it in production. |
| `DATA_DIR` | no, default `data` (`/data` in Docker) | Holds `app.db` and uploaded images and Slips. |
| `PORT` | no, default `8080` | Listen port. |

## Deploy

`compose.yaml` runs one container with one volume at `/data`, published on `127.0.0.1:${HOST_PORT:-8080}`. Put the variables above (plus `HOST_PORT` if 8080 is taken) in a `.env` next to it. Pending migrations run at startup. Then:

```bash
mise run up        # docker compose up -d --build
mise run up:seed   # once, only to load the example football 2026 Campaign
mise run logs      # follow the app log
mise run tunnel    # Cloudflare quick tunnel to HOST_PORT; prints a new https://*.trycloudflare.com URL
mise run down
```

For a LINE test, set the MINI App channel's Endpoint URL to `<tunnel URL>${BASE_PATH}/` and open `https://miniapp.line.me/<LIFF_ID>` in LINE. A quick tunnel gets a new URL on every start, so update the Endpoint URL each time.

The gateway must forward `${BASE_PATH}` without stripping it. Backups are the host's job: snapshot the whole `/data` volume, or use `sqlite3 app.db .backup`, never a plain file copy while the app runs.

### Dedicated subdomain

To serve the app at the root of its own hostname, such as `https://mmtc-events.codesook.dev`, put `BASE_PATH=/` in `.env`. `compose.yaml` passes it as both the build arg and the runtime variable, so `mise run up` rebuilds the image for the root and the server starts with the matching value. Then:

- Set the MINI App channel's Endpoint URL to `https://mmtc-events.codesook.dev/`.
- Sign in as Super Admin at `https://mmtc-events.codesook.dev/admin-login`.
- The session cookie is scoped to `Path=/`, and the gateway forwards every path to the container.

Changing `BASE_PATH` needs a rebuild (`mise run up`), because the web build bakes the value in.

### First Admin

1. Open `https://<host>${BASE_PATH}/admin-login` in a normal browser and sign in with `SUPER_ADMIN_USERNAME` / `SUPER_ADMIN_PASSWORD`.
2. Open the app in LINE with your own account and fill in the profile. You become a Member at once.
3. As the Super Admin, go to สมาชิก, open your Member page → แก้ไข and set the role to แอดมิน.

Promoted Admins can promote others too. The Super Admin cannot be demoted or deleted; to rotate its password, change `SUPER_ADMIN_PASSWORD` and `SESSION_SECRET` (the latter ends every open session).

## Adding shared UI primitives

```bash
bunx shadcn@latest add dialog -c packages/ui
```

```tsx
import { Button } from "@mmtc34/ui/components/button";
```
