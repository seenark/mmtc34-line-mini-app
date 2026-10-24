# Issue tracker: Beads (local)

Issues for this repo live in the local beads database (`.beads/`), managed with the `bd` CLI. Always pass `--json` where supported and run from the repo root.

## Conventions

- A spec or feature is an `epic`; its tickets are children created with `--parent <epic-id>`.
- Blocking edges are native beads dependencies: `bd create ... --deps blocked-by:<id>` or `bd dep add <blocked> <blocker>`.
- Triage state is a beads label (see `triage-labels.md`).
- Comments go through `bd comment`; the body is the issue description.
- Pre-existing spec: GitHub issue seenark/mmtc34-line-mini-app#1 (milestone 1). The beads epic links to it in its description.

## When a skill says "publish to the issue tracker"

`bd create --title ... --type <task|epic> --body-file <file> --labels <triage-label> [--parent <id>] [--deps blocked-by:<id>] --json`

## When a skill says "fetch the relevant ticket"

`bd show <id> --json`

## Frontier

`bd ready --json` lists open tickets whose blockers are all closed.
