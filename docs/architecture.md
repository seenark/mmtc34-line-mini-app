# Architecture: MMTC Events LINE MINI App, milestone 1

One Go package of feature files over one pure `domain` package; a table-driven Order state
machine with one `apply` entry point; a schema and route table frozen by the skeleton so every
wave ticket only **adds files**; navigation derived from route `staticData` so the web side has
no shared menu list either.

Go module path: `github.com/seenark/mmtc34-line-mini-app/apps/server`. The code is the source of
truth; this document records the design and why. The architect-arena sketch it was drawn from
(`arena/a/sketch/`) was deleted once milestone 1 landed; the "Synthesis decision" section lists
where the build deviated from that sketch.

Sections 1–9 and the Rationale describe milestone 1 and keep their ticket tables as history. Where
they differ from the code, section 10 is current: open membership (no Applicants, no approval, no
`/waiting` route), Admin Mode, Delivery Method `none`, a bank registry, an uploaded QR image for
payment, and Option size charts.

---

## 1. Usage (caller's view) — written first

### 1a. An HTTP scenario test (ticket .5) — `internal/server/payments_test.go`

```go
func TestSlipRejectReuploadConfirmUnderTransfer(t *testing.T) {
	h := New(t)                                   // temp SQLite + migrations + temp DATA_DIR, fake LINE, fake clock
	admin, alice := h.Admin(), h.Member("alice")  // real login + onboarding over HTTP
	j := h.OpenJersey()                           // 260-baht jersey, open, closes start+24h

	order := alice.Checkout(j.Campaign.ID, Pickup(Shirt(j, "L", 2))) // 520
	order = alice.UploadSlip(order.ID)
	Equal(t, order.Status, api.OrderPendingReview)

	reject := fmt.Sprintf("/api/admin/orders/%d/reject", order.ID)
	Equal(t, admin.POST(reject, api.RejectInput{}).Status(400).Error(), "กรุณาระบุเหตุผล")
	admin.POST(reject, api.RejectInput{Reason: "สลิปไม่ชัด"}).Status(200)
	Equal(t, alice.Order(order.ID).Actions.UploadSlip, true)

	order = alice.UploadSlip(order.ID)
	confirm := fmt.Sprintf("/api/admin/orders/%d/confirm", order.ID)
	order = As[api.Order](admin.POST(confirm, api.AmountInput{Amount: 500}).Status(200))
	Equal(t, order.Balance, int64(20))           // under-transfer stays owed

	h.Member("bob").GET("/api/slips/" + order.Payments[1].Slip).Status(404)
}

func TestSlipRefusedAfterCloseTime(t *testing.T) {
	h := New(t); alice := h.Member("alice"); j := h.OpenJersey()
	order := alice.Checkout(j.Campaign.ID, Pickup(Shirt(j, "M", 1)))
	h.Advance(25 * time.Hour)                                  // cross closeAt, no job runs
	alice.POST(fmt.Sprintf("/api/orders/%d/slip", order.ID), JPEG()).Status(409)
	Equal(t, alice.Order(order.ID).Status, api.OrderCancelled) // computed, not stored
	Equal(t, alice.Order(order.ID).Balance, int64(0))
}
```

Harness (`harness_test.go`, ticket .1): `New(t) *H`, `h.Advance(d)`, `h.Anonymous() *Client`,
`c.GET/POST/PUT/DELETE(path, body) *Res` (body = struct → JSON, `Upload` → multipart),
`r.Status(code) *Res` (fails with method, path, body), `r.Error()`, `r.Text()`, `As[T](r)`,
`Equal(t, got, want)`. Paths are written `/api/...`; the harness prefixes `/mmtc34`, so every
test also exercises BASE_PATH. It runs over `httptest.NewTLSServer` because session cookies are
`Secure`. Fixtures are methods in per-ticket files (`fixtures_auth_test.go` .2 → `h.Admin()`,
`h.Member(nick)`, `h.Visitor(nick)`; `fixtures_campaign_test.go` .3 → `h.OpenJersey()`,
`h.Close(id)`; `fixtures_order_test.go` .4 → `Shirt`, `Pickup`, `Ship`, `c.Checkout`, `c.Order`;
`fixtures_payment_test.go` .5 → `c.UploadSlip`). Fixtures go through HTTP only, never SQL.

### 1b. A lifecycle handler (ticket .6) — `internal/server/ledger.go`

```go
func (s *Server) recordCash(c *gin.Context) {
	in, ok := bind[api.AmountInput](c)
	if !ok { return }
	if in.Amount <= 0 { fail(c, &domain.Error{Kind: domain.Invalid, Msg: "จำนวนเงินต้องมากกว่า 0"}); return }
	me := actorOf(c)
	s.apply(c, domain.EvRecordCash, func(q *db.Queries, o domain.Order) error {
		return q.InsertCashPayment(c, db.InsertCashPaymentParams{
			OrderID: o.ID, Amount: in.Amount, ReviewedBy: me.MemberID, Now: rfc3339(s.now())})
	})
}
```

`apply` loads the Order, asks `o.Decide(actor, ev)`, compare-and-sets the status, runs the
effect, reloads, and writes `200 api.Order` — all in one IMMEDIATE transaction. No handler
checks status or role by hand.

### 1c. A web route adding an Admin screen with a badge (ticket .5) — `apps/web/src/routes/admin/payments.tsx`

```tsx
export const Route = createFileRoute("/admin/payments")({
  staticData: { nav: { slot: "admin", label: "ตรวจสลิป", order: 20, badge: pendingSlipCount } },
  loader: ({ context }) => context.queryClient.ensureQueryData(pendingSlips),
  component: SlipQueue,
});
```

Creating this file is the whole integration: the router plugin regenerates `routeTree.gen.ts`,
and `<NavSlot slot="admin" />` in `admin.tsx` lists it. No menu file is edited.

---

## 2. Module map

### apps/server

```
main.go                         env → db.Open → server.New → http.Server; `seed` subcommand
sqlc.yaml, tygo.yaml
internal/
  api/        leaf package: DTOs + enums (tygo source). config.go members.go campaigns.go orders.go
  domain/     PURE rules, imports only api + time
    errors.go     Error{Kind, Msg(Thai)}, Actor{MemberID, Admin}
    campaign.go   Campaign.Status(now)  ← effective Campaign status, the only one
                  VisibleTo, CanCheckout, ValidateCampaign, ValidateOption, ParseSizes
    order.go      Order, Order.Status() ← effective Order status, the only one
                  rules table, Order.Decide(actor, ev), Order.Actions(actor)
    money.go      Order.Money() ← Total/Paid/Balance, the only one
    checkout.go   ValidateOrder → ValidOrder (unexported fields: only constructible by validation)
    bank.go       Banks() registry (3-digit Bank of Thailand code → Thai name), BankName(code)
    member.go     ValidateProfile, SuperAdminLineUserID (reserved "superadmin")
  db/         sqlc output + open.go (Open: pragmas + goose) + seed.go (Seed: embedded seed/*.sql)
    migrations/00001_init.sql     complete schema (sketch has it in full)
    migrations/20261026120000_events_football_2026.sql   open membership, size_chart, bank_code, qr_source, delivery 'none'
    migrations/20261026180000_drop_qr_source.sql         drops campaigns.qr_source (no generated QR)
    migrations/20261102090000_campaign_photos.sql        campaign_photos (ordered image names per Campaign, sort_order 0 = Cover); moves cover_image into it and drops the column
    queries/<feature>.sql         one per ticket → <feature>.sql.go
    seed/10_football_2026.sql (.3), seed/20_football_2026_sponsors.sql (.10), seed/images/*.jpg (embedded; copied to DATA_DIR/files/images when absent)
  line/       verify.go: Client.Verify(ctx, idToken) → LINE userId
  spa/        embed.go: //go:embed all:dist  (Vite outDir)
  server/     Gin shell, one file per feature, all in package server
    server.go     Config, ConfigFromEnv, Verifier, Deps, New, tx, mountConfig
    routes.go     access-level groups + complete mount list (frozen after .2)
    respond.go    fail (only error→HTTP mapping), bind[T], idParam
    spa.go        NoRoute: SPA files / index fallback / BASE_PATH checks
    session.go    cookie sign/parse, requireSession/Member/Admin, actorOf
    order_core.go loadOrders (only Order read path), toAPIOrder (only mapping), apply, checkout
    auth.go members.go campaigns.go uploads.go orders.go sponsors.go payments.go
    admin_orders.go unlinked.go profile.go ledger.go exports.go roster.go
    harness_test.go fixtures_<feature>_test.go <feature>_test.go
```

Dependency direction: `server → domain → api`, `server → db`, `server → line` (via `Verifier`
interface owned by server). `domain` never sees `db` or Gin; `db` never sees `domain`.

Call chain for any request: `server/<feature>.go` → `domain/*.go` → `db/queries/<feature>.sql`.
Lifecycle requests add `server/order_core.go` (`apply`), which is the same for all of them.

### apps/web

```
vite.config.ts            base = `${BASE_PATH}/`, outDir ../server/internal/spa/dist, dev proxy
public/.keep              copied into dist so the committed dist/.keep survives emptyOutDir
src/main.tsx              QueryClient + router({ basepath: import.meta.env.BASE_URL, context })
src/api/client.ts         api<T>(method, path, body) with credentials; ApiError(status, Thai msg)
src/api/types.gen.ts      tygo output (git-ignored)
src/api/<feature>.ts      queryOptions + mutations per ticket
src/lib/nav.tsx           NavEntry staticData augmentation + <NavSlot slot params actions>
src/lib/liff.ts (.2)      init → idToken → POST /auth/line
src/lib/admin-mode.ts     Admin Mode: localStorage `mode`, launchMode, canSwitchMode, useRecordMode
src/lib/image.ts (.3)     shrinkImage(file): ≤1600px JPEG blob (createImageBitmap + canvas)
src/lib/cart.ts (.4)      localStorage `cart:<memberId>:<campaignId>`; useCart(campaignId); cleared on logout and account deletion
src/components/...        profile-form (.2), campaign/* (.3), order/* (.4), mode-switch (Admin Mode toggle)
src/routes/
  __root.tsx              shell, error boundary showing ApiError message
  login.tsx               LIFF login, or dev form when config.devLogin
  onboarding.tsx          (no waiting route: a new Member goes straight to _app)
  _app.tsx                gate: Member; bottom <NavSlot slot="member">
  _app/campaigns.index.tsx, _app/campaigns.$campaignId.tsx (header + memberCampaign tabs)
  _app/campaigns.$campaignId.{index,sponsors,roster,checkout}.tsx
  _app/orders.index.tsx, _app/orders.$orderId.tsx (the Member's own Order page, for Admins too outside Admin Mode; uploads the Slip inline through `components/payment/slip-upload.tsx`)
  _app/profile.tsx
  admin.tsx               gate: Admin; sidebar <NavSlot slot="admin">
  admin/members.*, admin/campaigns.*, admin/orders.$orderId.tsx (the admin Order page), admin/orders.$orderId.{edit,cash,adjust}.tsx, admin/payments.tsx
```

Web never computes money or lifecycle: it shows `order.total/paid/balance` and renders buttons from
`order.actions` (the server's `Order.Actions(viewer)`). Balance label: `balance > 0` "ค้างจ่าย N บาท",
`= 0` "จ่ายครบแล้ว", `< 0` "จ่ายเกิน N บาท" — a sign check, not arithmetic.

---

## 3. Order state machine (`domain/order.go`, `rules`)

Rows match on the **stored** status. "Owner needs open" means the Member owner additionally needs
the Campaign effectively open; Admins are never bound by close.

| Event | Who | From (stored) | To (stored) | Owner needs open | Side effect in `apply` effect |
|---|---|---|---|---|---|
| view | owner, admin | any | – | no | – (others get 404) |
| upload_slip | owner | awaiting_payment | pending_review | yes | save file, insert pending Payment |
| confirm_slip | admin | pending_review | paid | – | pending Payment → confirmed, amount, reviewer, time |
| reject_slip | admin | pending_review | awaiting_payment | – | pending Payment → rejected, reason |
| record_cash | admin | awaiting_payment, paid | paid | – | insert confirmed Payment, slip NULL |
| cancel | owner | awaiting_payment | cancelled | yes | – |
| cancel | admin | any (idempotent) | cancelled | – | reject a pending Payment ("ยกเลิกคำสั่งซื้อ") |
| edit | admin | any | – | – | replace Items (kept IDs keep unit price), delivery, ship fields |
| adjust | admin | any | – | – | insert Adjustment |

Not an event: **checkout** creates the Order (`Campaign.CanCheckout`: Member needs effectively open;
Admin on behalf needs non-draft). Errors from `Decide`: other Member's Order → 404; party not in any
row for the event → 403; party allowed but status/close forbids → 409 with a Thai reason.

Persistence: `SetOrderStatus ... WHERE id=@id AND status=@from` (`:execrows`); 0 rows → 409.
`payments_one_pending` (partial unique index) makes a double-tapped upload unable to create two
pending Payments. Together: every transition applies at most once.

**Effective status** — `Order.Status()`: stored `awaiting_payment` + `Campaign.Status(now)==closed`
+ `PlacedByAdmin == false` → `cancelled`; otherwise stored. **Effective Campaign status** —
`Campaign.Status(now)`: stored `open` with `CloseAt <= now` → `closed`. `now` is the injected clock,
captured once per request by `loadOrders` (`Order.Now`) and `toAPICampaign`.

**Balance** — `Order.Money()`: `Total = Σ qty×unit_price`; `Paid = Σ confirmed amounts`;
`Balance = (Status()==cancelled ? 0 : Total) − Paid + Σ adjustments`. Every consumer (api.Order,
exports, Slip queue, confirm prefill via `order.total`) reads this. SQL never sums money.

Single read path: `loadOrders(scope)` → 4 queries (`ListOrders` with optional order/campaign/member
filter, then Items/Payments/Adjustments `IN (sqlc.slice(ids))`) → `[]domain.Order`. Status filters
(admin list, queue, roster, exports) run in Go on `o.Status()` so effective status is never
re-implemented in SQL.

---

## 4. Auth and session

1. Web loads `GET /api/config` → `{liffId, devLogin, adminLogin}`.
2. Normal: `liff.init({liffId})` (SDK ≥ 2.31.1); not logged in → `liff.login()`; `liff.getIDToken()`
   → `POST /api/auth/line {idToken}`. Dev (`devLogin`): form → `POST /api/auth/dev {lineUserId}`
   (route returns 404 when `DEV_LOGIN` is off).
3. Server: `Verifier.Verify` (`line.Client`: form `id_token`, `client_id=LINE_CHANNEL_ID` to
   `https://api.line.me/oauth2/v2.1/verify`; `sub` is the userId). Failure → 401.
4. Super Admin (ADR 0005): `POST /api/auth/admin {username, password}` checks
   `SUPER_ADMIN_USERNAME`/`SUPER_ADMIN_PASSWORD` in constant time, at most 5 failures per minute
   per process (then 429). Success upserts the Member row whose `line_user_id` is the reserved
   `superadmin` (admin) and sets the same cookie. LINE and dev login refuse that id.
   Nobody can change the Super Admin's role or delete it. Other Admins are Members it promotes.
5. Cookie `mmtc34_session` = `b64url(userId).exp.b64url(HMAC-SHA256(SESSION_SECRET, userId.exp))`,
   HttpOnly, SameSite=Lax, Path=`BASE_PATH`, Max-Age 30 d, Secure unless `DEV_LOGIN`.
6. Every request under `self/member/admin`: `requireSession` parses the cookie and loads the Member
   **by LINE userId**; `requireMember`/`requireAdmin` gate on the fresh row. Role changes, deletion
   and linking an Order take effect on the next request.
7. `GET /api/me` → `{lineUserId, member|null}`. Web routing: 401 → `/login`; `member == null` →
   `/onboarding`; otherwise `_app`.
8. `POST /api/me` creates the Member at once (open membership, ADR 0006); no approval step. The
   `members.status` column and its `applicant` value survive in the schema but nothing writes
   `applicant`; the 2026-10-26 migration turns existing Applicants into Members.

Access groups (routes.go): `public` (config, auth) · `self` (me, onboard, edit/delete self) ·
`member` (has a Member, not deleted) · `admin` (`/api/admin/*`).

---

## 5. Uploads

- Web: `shrinkImage` → JPEG ≤ 1600 px (~100–300 KB) → `FormData` field `file`.
- Server `saveUpload(c, kind)`: `http.MaxBytesReader` 5 MB (413), sniff first 512 bytes with
  `http.DetectContentType` → `image/jpeg|png|webp` only (400), write temp then rename to
  `DATA_DIR/files/<kind>/<32 hex>.<ext>`. No files table.
- The **directory decides access**: `images/` (Campaign Photos, Option images, QR, logos) via
  `GET /api/images/:name` for Members; `slips/` via `GET /api/slips/:name` for the
  Payment's Order owner or Admins, else 404. Names must match `^[0-9a-f]{32}\.(jpg|png|webp)$`
  (no traversal). Cache: `private, max-age=31536000, immutable`.
- Campaign Photos: the client uploads each picture through the single-file image endpoint, then
  saves the Campaign with the ordered name list. Saving replaces the whole `campaign_photos` set
  in one transaction (array index = `sort_order`; unknown name or duplicate → 400, nothing saved).
  Lists load every row's photos with one `IN (sqlc.slice(ids))` query. A photo removed from a
  Campaign leaves its file in `images/`: there is no delete or sweep.
- Slip upload saves the file inside the `apply` effect, i.e. after the permission/state check.
  Shortcut: a file whose transaction then fails is orphaned; acceptable at this scale.

---

## 6. BASE_PATH

One env var, read at build time by Vite and at run time by Go:

| Consumer | How |
|---|---|
| Gin | `engine.Group(BASE_PATH).Group("/api")`; NoRoute → `spaHandler` |
| Vite | `base: \`${BASE_PATH}/\`` → asset URLs, `import.meta.env.BASE_URL` |
| Router | `createRouter({ basepath: import.meta.env.BASE_URL })` |
| API client | `` `${import.meta.env.BASE_URL}api` `` |
| Cookie | `Path=BASE_PATH`, or `/` at the root |
| Dev | mise `[env] BASE_PATH`; Vite proxies `${base}api` → `:8080` |

`spaHandler` (prefix mode): outside prefix → 404; exactly `/mmtc34` → 301 `/mmtc34/`. In both modes an
unmatched `<base>/api/*` → JSON 404; existing dist file → file (`/assets/*` immutable); else `index.html`
(no-cache). `BASE_PATH=/` (or empty) serves at the root for a dedicated subdomain.
`checkSPABase` makes startup fail if the embedded `index.html` references a different
`<base>/assets/` — a mismatched image fails loudly instead of serving blank pages.

---

## 7. Codegen, build, Docker, mise (`sketch/mise.toml`, `sketch/Dockerfile`, `sketch/compose.yaml`)

Generated artefacts and whether they are committed:

| Artefact | Tool | Committed | Why |
|---|---|---|---|
| `internal/db/*.sql.go`, `models.go`, `db.go` | sqlc | **yes** | `go build/test`, IDE work with no tool step; one output file per query file, so parallel tickets never touch the same generated file; `models.go` depends only on the frozen schema |
| `apps/web/src/api/types.gen.ts` | tygo | no | one file per Go package → every DTO-adding ticket would conflict on it |
| `apps/web/src/routeTree.gen.ts` | TanStack plugin | no | same; regenerated by `bun run dev/build`; `typecheck` script runs `tsr generate && tsc --noEmit` |
| `internal/spa/dist/*` | Vite | only `.keep` | embed needs the directory to exist |

`.gitignore` (ticket .1): replace `dist/` with `apps/web/dist/`, `apps/server/internal/spa/dist/*`,
`!apps/server/internal/spa/dist/.keep`; add `apps/web/src/api/types.gen.ts`,
`apps/web/src/routeTree.gen.ts`. Biome ignores `*.gen.ts`.

Order: `gen:sqlc` ∥ `gen:types` → `web:build` (Vite plugin generates routes, writes into
`internal/spa/dist`) → `go build` (embeds). mise tasks: `gen:sqlc`, `gen:types`, `gen`, `dev:server`
(`DEV_LOGIN=1 go run .`), `dev:web`, `dev` (both), `seed`, `test` (`go test ./...`), `check`
(vet + biome + typecheck), `web:build`, `build` (`bin/mmtc34`), `docker`.
`go.mod` drops the tygo tool entry; `mise.toml` drops `node`.

Dockerfile: `types` (golang, `go install tygo`, `tygo generate`) → `web` (oven/bun, `ARG BASE_PATH`,
`bun run build`) → `build` (`CGO_ENABLED=0 go build`) → `distroless/static:nonroot` with
`/data` owned by 65532, `VOLUME /data`. Compose: one service, one named volume at `/data`
(`app.db`, WAL files, `files/images`, `files/slips`).

---

## 8. Endpoints

| Group | Method + path | Ticket |
|---|---|---|
| public | GET /config | .1 |
| public | POST /auth/line, /auth/dev, /auth/admin (Super Admin, ADR 0005), /auth/logout | .2 |
| self | GET /me, POST /me | .2 |
| self | PUT /me, DELETE /me | .11 |
| member | GET /campaigns, /campaigns/:campaignId | .3 |
| member | GET /images/:name | .3 |
| member | GET /campaigns/:campaignId/sponsors | .10 |
| member | POST /campaigns/:campaignId/orders; GET /orders, /orders/:orderId; POST /orders/:orderId/cancel | .4 |
| member | POST /orders/:orderId/slip; GET /slips/:name | .5 |
| member | GET /campaigns/:campaignId/roster | .9 |
| member | GET /me/addresses | elx |
| admin | GET /members, /members/:memberId | .2 |
| admin | GET/POST /campaigns, PUT /campaigns/:campaignId; POST /campaigns/:campaignId/options; PUT/DELETE /options/:optionId; POST /images | .3 |
| admin | POST /campaigns/:campaignId/sponsors; PUT/DELETE /sponsors/:sponsorId | .10 |
| admin | GET /payments/pending; POST /orders/:orderId/confirm, /orders/:orderId/reject | .5 |
| admin | GET /campaigns/:campaignId/orders?status=; POST /orders/:orderId/cash, /orders/:orderId/adjustments | .6 |
| admin | POST /campaigns/:campaignId/orders (on behalf); PUT /orders/:orderId; POST /members | .7 |
| admin | GET, POST /orders/:orderId/link (move an Unlinked Member's Order, or all, to a Member with LINE; ADR 0006) | elx |
| admin | GET /campaigns/:campaignId/exports/{factory,shipping,pickup}.csv?includePending=1 | .8 |
| admin | PUT /members/:memberId, /members/:memberId/role | .11 |
| admin | GET /members/:memberId/addresses | elx |

Cancel is one endpoint for Member and Admin (the state machine decides by actor), not two.
Confirm/reject address the Order, not the Payment: there is at most one pending Payment per Order.

DTO changes vs. the prototype (`sketch/apps/server/internal/api/`): `CampaignDetail` drops `Sponsors`
(own endpoint, so .10 never edits .3); `OrderInput` drops `MemberID` and `AdminOrderInput{memberId, order}`
is added; `ItemInput.ID` lets an Admin edit keep captured prices; `Order` gains `placedBy` and
`actions: OrderActions`.

---

## 9. Per-ticket file ownership

Rule: after .1–.3 (serial), a wave ticket **only creates files**, except the shared lines listed.
.1 creates every `server/<feature>.go` with an empty `func (s *Server) mountX(r routes) {}` and
the complete mount list in `routes.go`; the owning ticket fills its file. The schema is complete in
`00001_init.sql`; a later change adds a timestamp-versioned migration, never an edit.

| Ticket | Wave | Server (apps/server) | Web (apps/web/src unless noted) | Shared lines |
|---|---|---|---|---|
| .1 skeleton | serial | `main.go`, `go.mod/sum` (gin, modernc sqlite, goose), `sqlc.yaml`, `tygo.yaml`, `internal/api/*`, `internal/db/{open.go,migrations/00001_init.sql,db.go,models.go}`, `internal/spa/{embed.go,dist/.keep}`, `internal/domain/errors.go`, `internal/server/{server.go,routes.go,respond.go,spa.go}` + empty mount files for all features, `harness_test.go`, `config_test.go`; root `mise.toml`, `.gitignore`, `Dockerfile`, `compose.yaml` | Better-T-Stack scaffold, `vite.config.ts`, `public/.keep`, `main.tsx`, `api/client.ts`, `lib/nav.tsx`, `routes/__root.tsx`, `routes/index.tsx` | – |
| .2 login | serial | `line/verify.go`, `server/{session.go,auth.go,members.go}`, `domain/member.go`, `db/queries/members.sql`, `fixtures_auth_test.go`, `auth_test.go` | `lib/liff.ts`, `api/me.ts`, `components/profile-form.tsx`, `routes/{login,onboarding,_app,admin}.tsx`, `routes/admin/{index,members.index,members.$memberId}.tsx` | `routes.go`: middleware on 3 group lines; `package.json`: `@line/liff` |
| .3 campaigns | serial | `domain/campaign.go`, `server/{campaigns.go,uploads.go}`, `db/queries/{campaigns,options}.sql`, `db/seed.go`, `db/seed/10_football_2026.sql`, `main.go` seed branch, `fixtures_campaign_test.go`, `campaigns_test.go` | `lib/image.ts`, `api/campaigns.ts`, `components/campaign/*`, `routes/_app/campaigns.{index,$campaignId,$campaignId.index}.tsx`, `routes/admin/campaigns.{index,new,$campaignId,$campaignId.index}.tsx` | – |
| .4 checkout | **1** | `domain/{order.go,money.go,checkout.go}`, `server/{order_core.go,orders.go}`, `db/queries/orders.sql`, `fixtures_order_test.go`, `orders_test.go` | `lib/cart.ts`, `api/orders.ts`, `components/order/*`, `routes/_app/campaigns.$campaignId.checkout.tsx`, `routes/_app/orders.{index,$orderId}.tsx` | `components/campaign/option-card.tsx` (.3): add-to-cart (only editor in wave 1) |
| .10 sponsors | **1** | `server/sponsors.go`, `db/queries/sponsors.sql`, `db/seed/20_football_2026_sponsors.sql`, `sponsors_test.go` | `api/sponsors.ts`, `routes/_app/campaigns.$campaignId.sponsors.tsx` (memberCampaign tab), `routes/admin/campaigns.$campaignId.sponsors.tsx` (adminCampaign tab) | – |
| .5 slips | **2** | `server/payments.go`, `db/queries/payments.sql`, `fixtures_payment_test.go`, `payments_test.go` | `api/payments.ts`, `components/payment/slip-upload.tsx` (inline on `routes/_app/orders.$orderId.tsx`, when uploadSlip), `routes/admin/payments.tsx` (admin nav + badge) | – |
| .7 on behalf + link | **2** | `server/{admin_orders.go,unlinked.go}`, `db/queries/{admin_orders,unlinked}.sql`, `admin_orders_test.go`, `unlinked_test.go` | `routes/admin/campaigns.$campaignId.new-order.tsx` (adminCampaign tab, reuses `components/order/checkout-form.tsx`), `routes/admin/orders.$orderId.edit.tsx` (adminOrder, when edit), `routes/admin/members.new.tsx` | – |
| .11 profile | **2** | `server/profile.go`, `db/queries/profile.sql`, `profile_test.go` | `routes/_app/profile.tsx` (member nav), `routes/admin/members.$memberId.edit.tsx` (adminMember) | – |
| .6 ledger | **3** | `server/ledger.go`, `db/queries/ledger.sql`, `fixtures_ledger_test.go`, `ledger_test.go` | `api/ledger.ts`, `routes/admin/campaigns.$campaignId.orders.tsx` (adminCampaign tab), `routes/admin/orders.$orderId.{cash,adjust}.tsx` (adminOrder) | – |
| .8 exports | **3** | `server/exports.go`, `exports_test.go` | `routes/admin/campaigns.$campaignId.exports.tsx` (adminCampaign tab) | – |
| .9 roster | **3** | `server/roster.go`, `roster_test.go` | `routes/_app/campaigns.$campaignId.roster.tsx` (memberCampaign tab) | – |

Generated files per wave ticket are its own `db/<file>.sql.go`; `models.go`/`db.go` do not change.
The order page (`orders.$orderId.tsx`, .4) already renders Items, Payments with reject reasons,
Adjustments with notes and Money from the start, so .5/.6 need no edit there; the admin buttons arrive
through `<NavSlot slot="adminOrder" actions={order.actions}>` and the Member's Slip upload renders inline.
Merge risk that remains: if two wave tickets both find a missing DTO field in the same
`api/*.go` file. Rule: a ticket adds new types in a new `api/<feature>.go`; changing an existing
struct is escalated to the parent instead of done in a wave.

---

## 10. After milestone 1

### Admin Mode (web only)

Outside Admin Mode an Admin is an ordinary Member and the Orders they place are their own
(GLOSSARY). The server knows nothing about the mode: an Admin's session always authorises
`/api/admin/*`, and `Order.Actions(actor)` still reports admin actions. Only the web chooses which
shell to show.

- `lib/admin-mode.ts` stores `mode` (`"member"` or `"admin"`) in `localStorage`; a blocked or
  unreadable store reads as `member`. `useRecordMode` writes it when a shell is shown, not when the
  router preloads one.
- `canSwitchMode(member)` is true for an Admin who is not the Super Admin. Only they get the
  `ModeSwitch` in both headers (โหมดสมาชิก / โหมดแอดมิน).
- `launchMode` honors the stored mode on the first call per page load only, so back from `/admin`
  to `/` does not bounce the Admin into admin again. `_app` redirects `/` to `/admin` when it says
  admin.
- The Super Admin has no Orders of their own: `_app` always redirects it to `/admin`, it never gets
  a switch, and the admin header shows ออกจากระบบ instead. It is always in admin.
- `/admin/orders/$orderId` is the admin Order page; `/orders/$orderId` is the Member's page, which
  an Admin also gets for their own Orders outside Admin Mode. The admin page is the one that shows
  admin actions.

### Delivery Method `none`

`api.DeliveryNone` is the third Delivery Method, for an Order that holds only donations (every Item's
Option has no fixed price). `domain.ValidateEdit` (`checkout.go`) makes the two agree both ways:
donation-only requires `none` ("คำสั่งซื้อที่มีแต่สมทบทุนไม่ต้องเลือกวิธีรับสินค้า"), and `none` is rejected
for any Order with a priced Item. Ship fields are discarded for `none` and `pickup`. The shipping and
pickup exports skip `none` Orders because they filter on `ship` and `pickup`. `orders.delivery` allows
`'none'` since the 2026-10-26 migration.

### Address choices

`GET /api/me/addresses` (Member) and `GET /api/admin/members/:memberId/addresses` (Admin, for an
on-behalf Order) return `[]api.AddressChoice{name, address, phone}` from one function,
`writeAddressChoices` in `profile.go`. The profile triple comes first, then the distinct shipping
snapshots of the Member's earlier Orders, newest first. Triples compare on trimmed fields, the first
occurrence wins, and a triple with any empty field is dropped. The web offers them in
`components/order/address-picker`, used by checkout and the on-behalf Order form.

### Order linking (`server/unlinked.go`)

An Unlinked Member has no LINE account (ADR 0002). Once the person has one, an Admin moves their
Orders:

- `GET /api/admin/orders/:orderId/link` → `{unlinkedMember, otherOrders}` so the dialog can offer
  "move all".
- `POST /api/admin/orders/:orderId/link {memberId, allOrders}` → `{moved, unlinkedRemoved}`. The
  target must be a live Member with a LINE account and not the Super Admin; the Order's current owner
  must be unlinked and not erased (409 otherwise).

Only `orders.member_id` changes, so Payments, Slips, Adjustments, the placed-by-Admin flag and
ship snapshots travel with the Order. When the Unlinked Member has no Orders left, its Sponsor
Positions are moved to the target first (the FK is `ON DELETE SET NULL`), then the row is deleted.
All of it runs in one transaction.

### Bank registry and payment QR

- `domain/bank.go` is the registry: 16 banks keyed by the 3-digit Bank of Thailand code, in picker
  order. `GET /api/config` returns it as `banks`. `campaigns.bank_code` replaces the free-text
  `bank_name`; `Campaign.bankName` is derived from the code, and `""` means no bank.
- The app never generates a payment QR. A Campaign shows its uploaded `qr_image` (for example from
  K PLUS) when set, plus the bank name, account name and number with a copy button. The Order
  page also shows the Balance as text, and the Member types that amount in their banking app.
- Decision record: a bank-account Thai QR (EMVCo tag 29, sub-tag 04 = bank code + account number,
  static and dynamic amount) was built and scanned in real banking apps. Every variant failed, and
  the class account has no PromptPay. The generator, `Order.paymentQr` and `campaigns.qr_source`
  were removed (migration `20261026180000_drop_qr_source.sql`).

### Option size chart

`options.size_chart` holds the stored image name of a size chart (`''` = none), set in the Option
editor and shown beside the size picker. Like every referenced image it must exist in
`files/images/` (400 otherwise).

### Seed

`go run . seed` (`mise run seed`; `mise run up:seed` in Docker) creates the football 2026 Campaign
with its two Options (the shirt with `size_chart`, and the open-amount donation). `db/seed.go`
embeds `seed/*.sql` and `seed/images/*.jpg`; `Seed` copies each image into `DATA_DIR/files/images`
only when it is absent, then runs the SQL files in name order inside one transaction. Each SQL file
is guarded by `WHERE NOT EXISTS`, so a second run changes nothing.

### Migration `20261026120000_events_football_2026.sql`

One migration, run with `-- +goose NO TRANSACTION` and a single `StatementBegin`/`StatementEnd`
block: `PRAGMA foreign_keys = OFF`, its own `BEGIN`/`COMMIT`, then `PRAGMA foreign_keys = ON`.
It sets every Applicant to approved, adds `options.size_chart`, `campaigns.bank_code` (guessed
from the old `bank_name`, then that column is dropped) and `campaigns.qr_source` (later dropped by
`20261026180000_drop_qr_source.sql`), and rebuilds `orders` so `delivery` accepts `'none'`. SQLite cannot widen a `CHECK`
in place, and dropping `orders` with foreign keys on would cascade-delete its Items. goose's
connection has to keep that pragma for the whole script, which is why it is one statement on one
connection.

---

## Rationale

### Problem

Eleven tickets, eight of them built concurrently in three waves by agents that see only the
files they open. The domain is small but has three rules that are easy to scatter: who may move
an Order between statuses, what an Order owes, and the two read-time "effective" statuses. The
constraints (Gin, sqlc, goose, tygo, one HTTP test seam, injectable verifier + clock, BASE_PATH
not stripped, one binary) are fixed. The non-obvious part is shaping files so that the copy-the-
nearest-example agent puts each rule in the one right place and parallel branches merge without
touching each other.

### Shape

- **Data first.** The full schema exists on day one with invariants as constraints: FK RESTRICT
  for "Option in use cannot be deleted", `CHECK`s tying Payment status to amount/reason/reviewer,
  a partial unique index for one pending Payment per Order, a `CHECK` that a deleted Member holds
  no contact data. No `CURRENT_TIMESTAMP` defaults: the injected clock is the only time source.
- **`domain` is pure and owns every rule once**: `Campaign.Status`, `Order.Status`,
  `Order.Money`, `rules`/`Decide`/`Actions`, `ValidateOrder`. `ValidOrder` has unexported fields,
  so inserting an unvalidated Order does not compile (per encode-lessons-in-structure).
- **Deep `apply`**: one function hides load → authorise → state check → compare-and-set → effect →
  reload → respond. Handlers for upload, confirm, reject, cash, cancel, edit, adjust are each a
  few lines and cannot forget a check (deep module, per the interface-depth discipline).
- **One read path** (`loadOrders` + `toAPIOrder`): status filters run in Go on effective status,
  so roster, exports, queue and admin list cannot disagree with the Order page.
- **No service/repository layer**: handlers call sqlc and domain directly; a request is three
  files (laziness-protocol, minimize-reader-load).
- **Parallelism by construction**: frozen schema, pre-allocated mount functions, per-feature query
  files, git-ignored single-file generators, and route-derived nav. The only cross-ticket edit in
  the waves is `option-card.tsx` by .4.
- **Validation at boundaries** (boundary-discipline): `bind[T]` for JSON shape, `domain.Validate*`
  for rules, `saveUpload` for files, `requireX` for access. Inside, types are trusted.
- **Deliberately not done**: no audit log, no scheduler, no files table, no cart on the server,
  no repository interfaces or mocks, no SQL money sums, no per-ticket migrations.

### Synthesis decision

Base: candidate A (anthropic seat). Cross-judge (`pstack-judge`) scored A 26/30, B 18/30 and picked A;
the orchestrator's independent read agreed. B (openai seat) is discarded.

Grafts from B, now binding:

- Cart key is `cart:<memberId>:<campaignId>` and is cleared on logout and account deletion.
- Upload limit: the request body is capped at 5 MB plus multipart overhead (6 MB); the file part
  itself must be at most 5 MB (413 otherwise).
- Image and Slip responses send `X-Content-Type-Options: nosniff`.
- Saving a Campaign, Option or Sponsor Position checks that each referenced image name exists in
  `files/images/` (400 otherwise).
- The harness has one test that builds the server with a non-default base path (`/x`) and hits
  `GET /x/api/config`, so nothing server-side hard-codes `/mmtc34`.
- Export CSV cells that start with `=`, `+`, `-` or `@` are prefixed with `'` (formula injection).

Fixes to A's sketch, now binding:

- The seed `//go:embed seed/*.sql` lives in `db/seed.go`, created by .3, not in `db/open.go`.
- Read-only Order views (`GET /orders/:id`) do not open an IMMEDIATE write transaction.
- The harness `Res` records method and path so `Status()` failures print them.
- `BASE_PATH` has one source per artefact: the Docker build arg sets the image's runtime `ENV`;
  compose passes `${BASE_PATH:-/mmtc34}` as both build arg and runtime env. tygo is pinned only in `mise.toml`.
- The Group table is named `class_groups`, avoiding SQLite's `GROUPS` keyword in sqlc's parser.
- `api.Order.placedBy` is renamed `placedByAdmin` (bool).
- Dropped as unrequested scope: Admin Slip upload, and the ban on an Admin demoting themself.
- SQLite connections set `busy_timeout` (5000 ms) alongside WAL, foreign keys and `_txlock=immediate`.

Rules ratified by the orchestrator (recorded in GLOSSARY and ADR 0001):

- An Order placed by an Admin on someone's behalf never reads cancelled because the Campaign closed.
  Admins own those Orders and cancel them explicitly.
- A cancelled Order (stored or effective) owes nothing for its Items: Balance = −Paid + Adjustments.

Rejected from B: runtime BASE_PATH via `<base href>` rewriting (build-time is simpler and the LIFF
endpoint URL fixes the prefix anyway), an Origin check (sibling apps on the shared hostname are
same-origin, so it cannot separate them; SameSite=Lax covers cross-site), checked int64 arithmetic,
a single DB connection, integrator-edited route lists, and hand-written nav links.

### Tradeoffs accepted

- We accept loading a whole Campaign's Orders into memory for every filtered view in exchange for
  effective status and Balance existing in exactly one place. Fine for one class; marked as a shortcut.
- We accept .1 writing empty mount functions for features it does not build in exchange for zero
  shared lines in `routes.go` across waves.
- We accept a full schema up front, written before features exist, in exchange for no migration
  numbering collisions and a stable `models.go`. Gaps found later become new timestamped migrations.
- We accept git-ignoring `types.gen.ts` and `routeTree.gen.ts` (a fresh clone needs `mise run gen`)
  in exchange for no generated-file conflicts. sqlc output stays committed because it splits per file.
- We accept BASE_PATH being baked into the web build (image built per prefix) in exchange for no
  runtime HTML rewriting; startup fails on mismatch instead of serving a blank page.
- We accept one extra request for Sponsor Positions in exchange for .10 never editing .3.
- We accept `domain` importing `api` enums in exchange for not duplicating every status type (a
  hand-synced list between domain and DTOs).

### Alternatives considered

- **Package per feature (`orders/`, `payments/`, … each with handler+store).** Better isolation on
  paper, but Orders are read by six features; each package either imports `orders` internals or
  reimplements loading and Balance. It exposes more surface (exported store types) and hides less.
- **Service + repository layers with interfaces.** Lets unit tests mock storage, which the single
  HTTP seam forbids anyway; adds pass-through methods and puts every request across 4–5 files.
- **State as stored-and-swept** (a goroutine cancelling unpaid Orders at close). Rejected by the
  constraints and also wrong: reopening a Campaign would need to un-cancel.
- **Effective status in SQL (`CASE` with a `@now` param).** Needs the same rule in every query that
  filters or lists; drifts. Go-side filtering keeps one copy.
- **Each ticket appends its own migration and route lines.** Simple, but waves of 2–3 tickets
  appending to the same list conflict every time, and goose sequence numbers collide.
- **Hand-kept admin/member menus.** Every screen ticket would edit the same array; replaced by
  `staticData.nav` + `<NavSlot>`.
- **Runtime BASE_PATH injection into built assets** (placeholder replaced at startup). Vite bakes
  base into JS chunk URLs and preload helpers; string-rewriting bundles is fragile.

### Open questions and risks

Built as the orchestrator's default; the owner has not confirmed these yet:

- An Order an Admin places on behalf (`placedByAdmin`) never reads cancelled when the Campaign
  closes, including one placed while the Campaign was still open. Admins cancel those explicitly.
- A cancelled Order (stored or effective) owes nothing for its Items, so a closed-unpaid Order shows
  Balance 0 and a cancelled paid Order shows a negative Balance (refund due).
- Cash is refused while a Slip is pending review, and refused when it exceeds the Balance.
- Admin edit keeps captured prices only for Items sent back with their `id`; a re-added line takes
  the current price.

Risks still open:

- A repeated partial cash request or a repeated Adjustment is recorded twice. The web disables the
  button while a request is in flight; the server has no replay protection.
- The real LIFF login and LINE ID-token verification have only run against a fake verifier and
  dev login. Check them on a real LINE channel before launch.
- Verify in deployment that Cloudflare Tunnel and the gateway pass multipart bodies up to 6 MB and
  keep the `Path=/mmtc34` cookie.
- Other apps on the same hostname are same-origin and can call this API with the user's cookie.
  Only a separate hostname removes that.
