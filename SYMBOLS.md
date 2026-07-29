# traderepublic-portfolio-downloader Symbol Reference

Workspace navigation index for the Trade Republic downloader fork.

---

## Architecture

**Legacy entry points:** `cmd/portfoliodownloader/public/main.go` ·
`cmd/portfoliodownloader/dev/main.go` · `cmd/example-generator/main.go`

**v2 entry points:** `v2/cmd/portfolio-downloader/main.go` ·
`v2/cmd/dev/main.go`

**Operational entry points:** `entrypoint.sh` ·
`scripts/generate-rest-client.sh`

The root and `v2` Go modules are independent generations. Their behavioral
owners, public seams, allowed dependencies, composition roots, and non-live
validation boundary are defined in
[`docs/architecture/ARCHITECTURE.md`](docs/architecture/ARCHITECTURE.md).

## Legacy module

- `cmd/portfoliodownloader/` - application orchestration and public/dev
  composition.
- `internal/traderepublc/api/` - authentication, REST, timeline, and WebSocket
  transport.
- `internal/traderepublc/portfolio/` - portfolio activity, document,
  instrument, and transaction behavior.
- `internal/database/`, `internal/filesystem/`, `internal/reader/`,
  `internal/writer/` - persistence and output adapters.

## v2 module

- `v2/cmd/portfolio-downloader/`, `v2/cmd/dev/` - v2 delivery and composition.
- `v2/internal/traderepublic/` - v2 authentication and brokerage adapters.
- `v2/internal/{instrument,message,timelinedetails,timelinetransactions,transaction}/`
  - v2 portfolio workflows.
- `v2/internal/{file,writer}/` - v2 output adapters.
- `v2/pkg/traderepublic/` - v2 public protocol schemas and generated types.

## Tests and safe gate

- `tests/architecture/` - required architecture-contract regression.
- `tests/fakes/` and `v2/tests/` - deterministic fixtures.
- `make check` - architecture, changed-file policy, and both Go module test
  lanes; it does not run a downloader or other live/maintenance operation.

---

_Full symbol index: [SYMBOLS-full.md](SYMBOLS-full.md)_
