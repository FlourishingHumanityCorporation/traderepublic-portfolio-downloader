# traderepublic-portfolio-downloader Full Symbol Reference

Workspace navigation index for the Trade Republic downloader fork.

---

## Commands

| Path | Purpose |
|------|---------|
| `cmd/portfoliodownloader/` | Main downloader command wiring. |
| `cmd/portfoliodownloader/public/` | Public binary entrypoint. |
| `cmd/portfoliodownloader/dev/` | Development binary entrypoint. |
| `cmd/example-generator/` | Example fixture generation helper. |

## Core Internal Packages

| Path | Purpose |
|------|---------|
| `internal/console/` | Auth prompts and terminal input helpers. |
| `internal/database/` | SQLite repository and persistence helpers. |
| `internal/filesystem/` | CSV and JSON read/write helpers. |
| `internal/reader/` | Reader abstractions for API responses. |
| `internal/traderepublc/api/` | Trade Republic API, auth, REST, timeline, and WebSocket clients. |
| `internal/traderepublc/portfolio/` | Portfolio activity, document, instrument, and transaction mapping. |
| `internal/writer/` | Output writer abstractions. |

## Test Surface

| Path | Purpose |
|------|---------|
| `internal/**/*_test.go` | Local deterministic Go tests run by `make check`. |
| `tests/fakes/` | Static fake responses used by package tests. |

## Safe Validation

| Command | Notes |
|---------|-------|
| `make check` | Runs `go test -v ./...`; does not run the downloader, Docker, live brokerage auth/API calls, document downloads, or response writes. |
