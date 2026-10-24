# Class Campaign

LINE MINI App where one school class runs money-collecting Campaigns (jersey orders, donations, parties). Members order like in an online shop and upload a transfer Slip per Order; Admins approve people, review Slips, record cash and Adjustments, and export factory, shipping and pickup lists.

Terms are defined in `GLOSSARY.md`. The design is in `docs/architecture.md`, and the decisions behind it in `docs/adr/`.

## Layout

- `apps/server` is the Go server (Gin, SQLite, sqlc, goose). It serves the API under `${BASE_PATH}/api` and embeds the built web app for every other path under `${BASE_PATH}`.
- `apps/web` is the React SPA (TanStack Router and Query, Tailwind, Biome), scaffolded with Better-T-Stack. `packages/ui` holds the shared shadcn/ui primitives.

## Toolchain

`mise install` installs Go, Bun, sqlc and tygo at the versions pinned in `mise.toml`. Then run `bun install` at the repo root.

## Develop

```bash
mise run seed     # create ./apps/server/data/app.db with the jersey 2569 Campaign
mise run dev      # Go server on :8080 with DEV_LOGIN on, Vite on :3001
```

Open `http://localhost:3001/mmtc34/`. With `DEV_LOGIN` on, the login page takes any LINE userId instead of LIFF. The dev server also enables the Super Admin as `admin` / `dev-super-admin-pw`; sign in at `/mmtc34/admin-login`, approve Applicants and promote them to Admin.

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
| `BASE_PATH` | no, default `/mmtc34` | URL prefix. Baked into the web build (Docker build arg) and read by the server; the server refuses to start if they differ. |
| `LIFF_ID` | yes, unless `DEV_LOGIN` | LIFF app ID served to the browser. |
| `LINE_CHANNEL_ID` | yes, unless `DEV_LOGIN` | Channel ID used to verify LIFF ID tokens with LINE. |
| `SESSION_SECRET` | yes | At least 32 bytes; signs the session cookie. |
| `SUPER_ADMIN_USERNAME` | no, default `admin` | Username of the Super Admin account. |
| `SUPER_ADMIN_PASSWORD` | yes, to bootstrap | At least 12 characters. Enables the Super Admin login at `${BASE_PATH}/admin-login`; unset disables it. |
| `DEV_LOGIN` | no, default off | Enables login by raw userId and non-Secure cookies. Never set it in production. |
| `DATA_DIR` | no, default `data` (`/data` in Docker) | Holds `app.db` and uploaded images and Slips. |
| `PORT` | no, default `8080` | Listen port. |

## Deploy

`compose.yaml` runs one container with one volume at `/data`. Put the variables above in a `.env` next to it, then:

```bash
docker compose up -d --build
docker compose run --rm app seed   # once, only to load the example jersey Campaign
```

The gateway must forward `${BASE_PATH}` without stripping it. Backups are the host's job: snapshot the whole `/data` volume, or use `sqlite3 app.db .backup`, never a plain file copy while the app runs.

### First Admin

1. Open `https://<host>${BASE_PATH}/admin-login` in a normal browser and sign in with `SUPER_ADMIN_USERNAME` / `SUPER_ADMIN_PASSWORD`.
2. Open the app in LINE with your own account and fill in the profile. You become an Applicant.
3. As the Super Admin, go to สมาชิก → รออนุมัติ, approve yourself, then open your Member page → แก้ไข and set the role to แอดมิน.

Promoted Admins can promote others too. The Super Admin cannot be demoted or deleted; to rotate its password, change `SUPER_ADMIN_PASSWORD` and `SESSION_SECRET` (the latter ends every open session).

## Adding shared UI primitives

```bash
bunx shadcn@latest add dialog -c packages/ui
```

```tsx
import { Button } from "@mmtc34/ui/components/button";
```
