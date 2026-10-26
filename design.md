# Design — MMTC Events

Locked design system for the web app. Every page restyle reads this file first. Extend it here
instead of overriding locally.

## Genre
modern-minimal. Theme: custom, "jersey navy, white paper, electric blue". Light mode only; no
`.dark` class is ever set.

## Macrostructure family
App pages only: Workbench. Member shell is a single `max-w-md` column with a white header and a
bottom tab bar. Admin shell is a full-width navy header with a scrolling tab row over a
`max-w-3xl` column. Pages carry no enrichment; function carries the page.

## Theme
Tokens live in `packages/ui/src/styles/globals.css` under the shadcn variable names. Pages use
Tailwind token utilities (`bg-primary`, `text-secondary`), never raw colour values.

- `--background` paper, near white, hue 265
- `--foreground` navy-tinted ink
- `--primary` oklch(0.30 0.09 265), jersey navy: header, primary buttons, headings
- `--secondary` oklch(0.52 0.21 262), bright blue: links, selected states, status, focus ring
- `--accent` pale blue tint with navy text: hover and highlighted rows
- `--muted` / `--muted-foreground` quiet surfaces and secondary text
- `--destructive` the only warm colour. No gold or yellow anywhere.
- `--primary-muted-foreground`, `--primary-raised` text and surfaces drawn on the navy header

Accent stays under about 5% of a viewport. Blue is a highlighter, not a fill.

## Typography
LINE Seed Sans TH (SIL OFL), self-hosted woff2 in `apps/web/public/fonts`, `font-display: swap`.
Regular covers 400 to 500, Bold 600 to 700, ExtraBold 800 to 900. Page titles `text-2xl
font-extrabold`, section titles `text-base font-semibold`, body `text-sm`. Thai leading is set in
the theme (`--text-*--line-height`); never use `leading-none` or `leading-tight` on Thai text.

## Shape and size
Radius follows `--radius`: controls `rounded-md`, cards and dialogs `rounded-xl`, switches and
badges `rounded-full`. Buttons and inputs are 40px tall (`lg` 48px). Touch targets stay at least
36px, bottom tabs 56px.

## Motion
Easings `--ease-out`, `--ease-in`, `--ease-in-out`; durations `--dur-micro`, `--dur-short`.
Colour and opacity transitions only. `prefers-reduced-motion` drops transitions to 150ms and
removes slide and zoom.

## Admin Mode
`lib/admin-mode.ts` stores the last mode per device. Admins signed in with LINE get the
`โหมดสมาชิก / โหมดแอดมิน` switch in both shell headers. Members and the Super Admin never see it.

## Responsive floor
No horizontal scroll at 320, 375, 414 and 768px. Clickable labels never wrap. Use `min-w-0` on
flex and grid children that hold text, and let wide rows scroll inside their own container.
