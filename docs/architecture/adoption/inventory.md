---
created: 2026-07-23
last_updated: 2026-07-23
---
# STD-028 pre-move inventory — traderepublic-portfolio-downloader

Inventory date: 2026-07-23
Registry key: `traderepublic-portfolio-downloader`
Architecture scope and effective Git root:
`/Users/paulrohde/CodeProjects/tools/traderepublic-portfolio-downloader`
Isolated candidate:
`/Users/paulrohde/CodeProjects/.worktrees/traderepublic-portfolio-downloader-std028-20260723`
Base revision: `29fdf803aab8bc886af75a0073e2542d31e27b41`
Base tree: `2328d6294c46052278d03e02c4d7efef3404aba6`
Canonical upstream: `fhc/main`, `0 ahead / 0 behind` at selection
Candidate branch: detached HEAD

The public `origin` remote points at the external upstream fork authority.
The local `main` branch tracks `fhc/main`; it does not track `origin/main`.
Candidate and publication state therefore use `fhc/main`, while the external
remote remains consumer/provenance context only.

The primary checkout was clean and was not edited. No matching downloader
process or adoption worktree existed. No downloader binary, brokerage
authentication, brokerage API/WebSocket request, document or response writer,
Docker command, reset target, code generator, or release build was run.

The starting root `make check` passed, but its `go test -v ./...` command
covered only the root Go module. The nested `v2/go.mod` module was outside that
pattern. A separate non-live `make -C v2 test` exposed one pre-existing test
contract failure: five model fields are intentionally pointers, while the test
compared them directly with scalar values. The same lane also contains
explicit skipped/incomplete v2 cases. The adoption must make the normal gate
compile and test both modules without claiming those skipped cases as live
acceptance.

The project is substantial because it contains two independently changing Go
modules, multiple operator/development entry points, brokerage transports,
portfolio workflows, persistence/output adapters, generated protocol code,
and maintenance/release entry points.

## Behavior inventory

| Capability | Pre-move paths | State or policy owned | Entry points | Candidate owner |
|---|---|---|---|---|
| Legacy runtime primitives | Direct files `internal/{const,counter,datetime,datetime_test,time,time_test}.go` | Shared constants, counters, date/time values | Internal imports | `v1_core` |
| Legacy persistence and IO | `internal/database`, `filesystem`, `reader`, `writer` | SQLite repositories, CSV/JSON/filesystem readers and writers | Internal imports | `v1_persistence` |
| Legacy brokerage transport | `internal/console`, `internal/traderepublc/api` | Console auth, REST/WebSocket protocol access, generated REST client | Root downloader | `v1_broker_transport` |
| Legacy portfolio behavior | `internal/traderepublc/portfolio` | Activity, document, instrument, and transaction processing | Root downloader | `v1_portfolio_domain` |
| Legacy delivery and composition | `cmd/portfoliodownloader/**` | App orchestration, flags, Wire assembly, process exit | Public/dev binaries | `v1_delivery`, `v1_composition_roots` |
| Legacy fixture generation | `cmd/example-generator`, `tests/fakes`, `assets` | Deterministic example CSV inputs | Developer generator | `v1_composition_roots`, `v1_test_support` |
| v2 protocol authority | `v2/pkg/traderepublic/**` | Schemas, generated types, publisher/WebSocket client | v2 internal consumers | `v2_protocol` |
| v2 brokerage transport | `v2/internal/traderepublic/**` | API and auth adapters | v2 operator binary | `v2_broker_transport` |
| v2 portfolio workflows | `v2/internal/{instrument,message,timelinedetails,timelinetransactions,transaction}` | Event mapping and portfolio transaction behavior | v2 operator/dev binaries | `v2_workflows` |
| v2 output adapters | `v2/internal/{file,writer}` | CSV/output delivery | v2 operator binary | `v2_output_adapters` |
| v2 runtime primitives | `v2/internal/const.go`, `v2/internal/bus`, `v2/internal/console` | Constants, event bus, operator input | v2 consumers | `v2_foundation` |
| v2 delivery and composition | `v2/cmd/dev`, `v2/cmd/portfolio-downloader` | Dependency assembly and process lifecycle | v2 dev/operator binaries | `v2_delivery`, `v2_composition_roots` |
| Maintenance/release | `scripts/generate-rest-client.sh`, `entrypoint.sh`, Makefiles, Docker/GitHub metadata | Generated-client maintenance and release/container assembly | Explicit maintenance paths | `maintenance_entrypoints`, `automation_contract`, `build_control` |

## Entry points and compatibility decisions

| Identity | Consumer | Decision |
|---|---|---|
| `cmd/portfoliodownloader/public/main.go` | Released legacy public binary | Preserve exact source path and command behavior |
| `cmd/portfoliodownloader/dev/main.go` | Local response-writing development binary | Preserve exact source path; keep live writes outside routine checks |
| `cmd/example-generator/main.go` | Maintainer fixture generation | Preserve as an explicit maintenance composition root |
| `v2/cmd/portfolio-downloader/main.go` | In-progress v2 operator binary | Preserve exact source path |
| `v2/cmd/dev/main.go` | In-progress v2 fixture/development workflow | Preserve exact source path |
| `v2/cmd/websocket-downloader/**` | No executable consumer | Retire: the Go file was 543 lines of comments, had no active `main`, and its README named packages and Make targets that do not exist |
| `scripts/generate-rest-client.sh` | `script-surfaces.toml`, maintainer | Preserve exact path and declare it; do not run during adoption |
| `entrypoint.sh` | Container image/runtime | Preserve executable path and mode |

Workspace-wide source search found no CodeProjects import or subprocess
consumer of this repository or its command paths. The external GitHub fork
and release users remain compatibility context for the legacy public command;
no public command implementation is moved in this adoption.

## Persisted and external identities

| Identity | Path/value | Owner/reader |
|---|---|---|
| Root Go module | `github.com/dhojayev/traderepublic-portfolio-downloader` | Legacy packages and external builds |
| v2 Go module | `github.com/dhojayev/traderepublic-portfolio-downloader/v2` | v2 packages and external builds |
| Auth state | `.session`, `.refresh` | Legacy and v2 auth adapters; gitignored |
| Root output | `transactions.csv` | Legacy portfolio writer; gitignored by operator convention |
| Documents | `documents/transactions`, `documents/activity` | Legacy/v2 document workflows; gitignored |
| Raw responses | `responses/**`, `v2/responses/**`, `v2/debug/**` | Development/debug adapters; gitignored |
| Example input | `assets/transactions.csv` | Tracked fixture/example generator |
| Generated v1 REST client | `internal/traderepublc/api/restclient/openapi_gen.go` | v1 transport |
| Generated v2 protocol | `v2/pkg/traderepublic/*_gen.go`, schemas and OAPI input | v2 protocol |

No secret value or local auth-state file was read. `v2/.env.example` is a
tracked template; real `.env` and auth/output paths stay outside candidate
identity.

## Resource and build census

| Class | Paths | Governance |
|---|---|---|
| Automation/release | `.github/**`, `docker-compose.yml`, `entrypoint.sh` | `automation_contract`, `maintenance_entrypoints` |
| Developer control | `.continue/**`, `.vscode/**`, `.editorconfig` | `agent_control`, `developer_control`, `build_control` |
| Build/package | Root and `v2` Makefiles, Go modules/sums, golangci configs, `.gitignore` files | `build_control` |
| Generated-input/protocol | `scripts/**`, `script-surfaces.toml`, `v2/pkg/traderepublic/{schemas,oapi}/**` | Maintenance root and `v2_protocol` |
| v2 protocol fixtures | `v2/tests/**` | `v2_test_support` |
| Knowledge | `README.md`, `v2/README.md`, `REBUILD.md`, `LICENSE`, `SYMBOLS*.md`, `docs/**` | `project_knowledge` |
| Fixture | `assets/**`, `tests/**` | `v1_test_support` |

The base tree had no tracked symlinks. The repository contains no active
`CLAUDE.md`. The generated symbol indexes are retained and must be refreshed
after stable paths are known.

## Dependency inventory

The two Go modules have no first-party imports between them.

The legacy graph is:

```text
v1_composition_roots -> v1_delivery -> v1_portfolio_domain
                                \-> v1_broker_transport
v1_portfolio_domain -> v1_broker_transport -> v1_persistence -> v1_core
v1_test_support -> v1_portfolio_domain / transport / persistence / core
```

The v2 graph is:

```text
v2_composition_roots -> v2_delivery -> v2_output_adapters -> v2_workflows
                                      \-> v2_broker_transport
v2_workflows -> v2_broker_transport -> v2_protocol
v2_workflows -> v2_foundation
v2_broker_transport -> v2_foundation
```

Go makes unexported cross-package symbol access impossible. A source/import
census found no first-party cross-generation dependency and no dynamic module
identity string that would be invalidated by the bounded deletion.

## Initial debt classification

- The architecture declaration, human architecture truth, local Appcheck
  binding, generated hook, and fleet ledger entry were absent: adoption work.
- The root native gate omitted the nested v2 module: verification-system debt
  to repair, not baseline.
- The v2 model test compared five optional pointer fields as scalars:
  pre-existing test-contract drift to repair without production changes.
- `v2/cmd/websocket-downloader` was a non-executable commented prototype with
  stale operational documentation: safe deletion after consumer search.
- Generated protocol files exceed normal hand-written file sizes but are not
  scripts or composition roots and remain owned generated artifacts.
- `scripts/generate-rest-client.sh` refers to an absent root
  `openapi-rest.yaml`; generation remains an explicit unverified maintenance
  surface and is not invoked or falsely claimed green.
- The shared symbol generator does not index Go. Its trial run produced empty
  output, which was rejected; the useful Go navigation indexes were restored,
  updated manually, and kept under architecture ownership.

No architecture baseline entry is expected after the stale prototype is
deleted and every remaining path has one truthful owner.
