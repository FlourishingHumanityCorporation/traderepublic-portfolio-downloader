# traderepublic-portfolio-downloader Symbol Reference

Workspace navigation index for the Trade Republic downloader fork.

---

## Architecture

**Entry Points:** `cmd/portfoliodownloader/app.go` (application wiring) ·
`cmd/example-generator/main.go` (example fixture generation)

**Core Areas:** console auth/input, Trade Republic API clients, timeline
transaction/activity processors, portfolio document/instrument/transaction
mapping, CSV/JSON filesystem writers, SQLite repository helpers, and fake
response fixtures.

---

## 1. Commands

- `cmd/portfoliodownloader/` - downloader command entrypoint and public/dev
  variants.
- `cmd/example-generator/` - fixture/example generation helper.

## 2. Internal Packages

- `internal/console/` - terminal auth prompts and password input helpers.
- `internal/database/` - SQLite repository setup.
- `internal/filesystem/` - CSV and JSON file read/write helpers.
- `internal/reader/` - request/response reader abstractions.
- `internal/traderepublc/api/` - Trade Republic auth, headers, REST, timeline,
  and WebSocket clients.
- `internal/traderepublc/portfolio/` - activity, document, instrument, and
  transaction mapping.
- `internal/writer/` - writer abstractions.

## 3. Tests And Fixtures

- `internal/**/*_test.go` - local unit tests used by `make check`.
- `tests/fakes/` - static fake Trade Republic responses for deterministic tests.

---

_Full symbol index: [SYMBOLS-full.md](SYMBOLS-full.md)_
