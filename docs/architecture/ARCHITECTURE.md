---
created: 2026-07-23
last_updated: 2026-07-23
---
# traderepublic-portfolio-downloader architecture

## Overview

This repository owns a local Trade Republic portfolio-downloader fork with a
released legacy Go module and an in-progress v2 Go module. Both modules
authenticate and retrieve brokerage data only when an operator explicitly
runs a binary. Routine architecture and test gates remain non-live.

The repository does not own Trade Republic's API, account/session lifecycle,
Docker itself, external release infrastructure, or the operator's portfolio
data. Architecture conformance does not make brokerage authentication,
downloads, response writes, document writes, generation, reset, or release
commands safe to run automatically.

The project was already physically layered before STD-028 adoption. The target
keeps those cohesive package boundaries, deletes one non-executable commented
prototype, and makes both real Go modules part of the normal check. Folder
depth is not an objective.

## Governed roots and resources

| Surface | Paths | Policy |
|---|---|---|
| Legacy Go source | `cmd/**`, `internal/**` | v1 delivery, transport, portfolio, persistence, and composition owners |
| Root test support | `tests/**` | Architecture regression and legacy fixtures |
| v2 Go source | `v2/cmd/**`, `v2/internal/**`, `v2/pkg/**` | v2 delivery, workflows, adapters, protocol, and composition owners |
| Maintenance scripts | `scripts/**`, `entrypoint.sh` | Exact explicit maintenance/container entry points |
| Non-code resources | `.continue/**`, `.github/**`, `.vscode/**`, `assets/**`, `docs/**`, `v2/tests/**` | Agent, automation, developer, fixture, and knowledge owners |
| Exact controls | The `resource_files` list in `appcheck.toml` | Build, package, gate, symbol, and v2 nested-module files |

The source roots deliberately stop at `v2/cmd`, `v2/internal`, and `v2/pkg`
rather than treating the nested module's build metadata as executable source.
The root limit is ten direct executable files per source/script root; no
current root reaches it. Script and composition roots are limited to 200
lines.

## Components

| Owner | Owned paths | Owns | Does not own | Public seam | Allowed dependencies |
|---|---|---|---|---|---|
| `v1_core` | Direct files `internal/{const,counter,datetime,datetime_test,time,time_test}.go` | Shared legacy constants, counters, date/time values | IO, brokerage transport, delivery | Exported symbols in the root `internal` package | None |
| `v1_persistence` | `internal/{database,filesystem,reader,writer}/**` | SQLite/CSV/JSON/filesystem ports and adapters | Brokerage protocol or portfolio policy | Exported repository, reader, writer, and CSV/JSON types/functions | `v1_core` |
| `v1_broker_transport` | `internal/console/**`, `internal/traderepublc/api/**` | Console authentication, REST/WebSocket access, generated REST client | Portfolio transaction interpretation | Exported client/auth/timeline packages | Core and persistence |
| `v1_portfolio_domain` | `internal/traderepublc/portfolio/**` | Activity, document, instrument, and transaction processing | Process lifecycle and external release | Exported portfolio package handlers/builders/models | Core, persistence, broker transport |
| `v1_delivery` | `cmd/portfoliodownloader/app.go` | Legacy application orchestration | Reusable transport/domain/IO behavior | `portfoliodownloader.NewApp` and `App.Run` | All v1 inner owners |
| `v1_composition_roots` | `cmd/example-generator/main.go`, `cmd/portfoliodownloader/{dev,public}/**` | Process setup, flags, Wire assembly, fixture-generator dispatch | Reusable business behavior | Exact `main.go` paths | v1 owners and test support |
| `v1_test_support` | `assets/**`, `tests/**` | Deterministic fixtures and architecture regression proof | Production behavior | Go test packages and fixture values | v1 core, persistence, transport, portfolio |
| `v2_foundation` | `v2/internal/const.go`, `v2/internal/bus/**`, `v2/internal/console/**` | v2 constants, event bus, operator input primitives | Brokerage protocol and transaction decisions | Exported v2 internal package types | None |
| `v2_protocol` | `v2/pkg/traderepublic/**` | Schemas, generated protocol types, public publisher/WebSocket seam | Operator workflow or output policy | `pkg/traderepublic` exported API | None |
| `v2_broker_transport` | `v2/internal/traderepublic/**` | v2 auth and API adapters | Portfolio interpretation and output | Exported internal auth/API constructors and interfaces | Foundation and protocol |
| `v2_workflows` | `v2/internal/{instrument,message,timelinedetails,timelinetransactions,transaction}/**` | v2 message mapping and portfolio transaction workflows | Process lifecycle and filesystem output | Exported handlers, mappers, builders, and models | Foundation, protocol, broker transport |
| `v2_output_adapters` | `v2/internal/{file,writer}/**` | v2 CSV/output adapters | Brokerage calls and transaction policy | Exported output handlers/interfaces | Foundation and workflows |
| `v2_delivery` | `v2/cmd/dev/app/**`, `v2/cmd/portfolio-downloader/{app,args}.go` | v2 application assembly and command arguments | Reusable inner behavior | Command app/argument types | All v2 inner owners |
| `v2_composition_roots` | `v2/cmd/{dev,portfolio-downloader}/main.go` | Process lifecycle and dependency composition | Reusable behavior | Exact `main.go` paths | v2 delivery and inner owners |
| `v2_test_support` | `v2/tests/**` | Frozen JSON protocol fixtures | Production behavior | Fixture paths | None |
| `maintenance_entrypoints` | `entrypoint.sh`, `scripts/**` | Container build dispatch and generated-client maintenance | Runtime portfolio behavior | Exact executable paths | None |
| `automation_contract` | `.github/**`, `docker-compose.yml` | CI/release and container metadata | Runtime implementation | Workflow and Compose identities | None |
| `project_knowledge` | Root/v2 READMEs, REBUILD, LICENSE, SYMBOLS, `docs/**` | Human architecture/setup/navigation truth | Runtime behavior | Document paths | None |
| `agent_control` | `.claudeignore`, `.continue/**`, `AGENTS.md` | Repository-local agent rules | Runtime behavior | `AGENTS.md` | None |
| `developer_control` | `.vscode/**` | Editor launch/settings metadata | Production behavior | VS Code metadata paths | None |
| `build_control` | Exact build/gate/package files in `appcheck.toml` | Go modules, Make gates, lint/generator config, architecture/hook control | Product behavior | `make architecture-check`, `make check` | None |

## Allowed dependency direction

The two product generations are independent:

```text
v1_composition_roots -> v1_delivery -> v1_portfolio_domain
                                \-> v1_broker_transport
v1_portfolio_domain -> v1_broker_transport -> v1_persistence -> v1_core

v2_composition_roots -> v2_delivery -> v2_output_adapters -> v2_workflows
                                      \-> v2_broker_transport
v2_workflows -> v2_broker_transport -> v2_protocol
v2_workflows -> v2_foundation
```

No v1 owner may import v2, and no v2 owner may import v1. Generated protocol
code is inner data/transport authority, never a delivery dependency. A new
edge requires this document, `appcheck.toml`, the focused architecture test,
and the architecture gate to change before the import lands.

## Named public seams and compatibility

The following identities are compatibility contracts:

- Root module `github.com/dhojayev/traderepublic-portfolio-downloader`.
- Nested module `github.com/dhojayev/traderepublic-portfolio-downloader/v2`.
- Legacy public binary source `cmd/portfoliodownloader/public/main.go`.
- Legacy development binary source `cmd/portfoliodownloader/dev/main.go`.
- v2 operator binary source `v2/cmd/portfolio-downloader/main.go`.
- v2 development source `v2/cmd/dev/main.go`.
- Container launcher `entrypoint.sh`.
- Generated-client maintenance path `scripts/generate-rest-client.sh`.
- Auth/output identities `.session`, `.refresh`, `transactions.csv`,
  `responses/**`, `documents/**`, and v2 debug/output equivalents.

No CodeProjects source consumer imports this repository. The legacy public
binary remains external-release compatibility scope, so its implementation is
not moved during adoption.

The former `v2/cmd/websocket-downloader` tree is not a compatibility seam. Its
Go file contained only comments, no executable `main`, and referenced removed
packages. Its adjacent README named nonexistent Make targets. Both are
retired, and the v2 root README now names only real entry points.

## Composition-root and maintenance policy

Only the six exact files in `architecture.composition_roots` may compose
Go processes or maintenance dispatch; `entrypoint.sh` is the separately
declared container entry point. Go `main.go` files may parse flags, build
dependencies, call a named app/operation seam, map errors to exit status, and
nothing more. Reusable transport, portfolio, persistence, mapping, or output
behavior belongs in an inner owner.

`scripts/generate-rest-client.sh` and `entrypoint.sh` remain executable because
real maintenance/container consumers use those paths. Generation, container
assembly, and cross-platform release builds are not routine checks. A new
script requires a real external consumer, an exact declared composition root,
an entry in `script-surfaces.toml` when applicable, and a named implementation
owner.

## Existing debt and temporary exceptions

The architecture baseline is empty. No flat-root, oversized-script,
oversized-composition-root, ownership, dependency-edge, cycle, or
configuration finding is retained.

Large `*_gen.go` protocol files are generated, owned source artifacts rather
than hand-written scripts or composition roots. Modify their schema/generator
inputs and regenerate them; do not hand-split generated output to satisfy a
folder metric.

The root generator script currently references an absent `openapi-rest.yaml`.
That pre-existing maintenance defect is not an architecture baseline item and
the adoption does not claim the generator passes. Repair requires identifying
the authoritative API contract before the script is run.

Several v2 tests are explicitly skipped because the rewrite remains
incomplete. The normal gate reports those skips honestly; they are not proof
of live brokerage behavior. The v2 model-builder test was corrected to assert
the intentional pointer-valued optional fields without changing production
code.

The shared `regenerate-symbols` command indexes Python, TypeScript, and Swift,
not Go. Running it here produces an empty index, so `SYMBOLS.md` and
`SYMBOLS-full.md` remain manually maintained Go navigation artifacts. Their
architecture paths are governed, but a generator green status is not claimed.

## Evolution rules

- Keep v1 and v2 dependency graphs independent until an explicit retirement or
  migration decision names the compatibility boundary.
- Put Trade Republic wire/schema mechanics in the matching transport/protocol
  owner; put portfolio interpretation in the matching domain/workflow owner.
- Put filesystem/CSV/SQLite behavior in persistence/output owners.
- Keep `cmd/**/main.go`, `entrypoint.sh`, and `scripts/**` composition-only.
- Delete dead commented programs instead of preserving them as apparent
  executable authority.
- Never add brokerage, Docker, generation, reset, or output mutation to
  `make check`.
- When a nested Go module is added, add its deterministic test lane to the root
  gate in the same change.
- Do not baseline adopter-created debt, invalid declarations, cycles, or
  ownership gaps.

## Agent pre-write checklist

Before adding or moving production behavior, state:

1. **Generation** — is the change v1, v2, or shared build/knowledge control?
2. **Owning module** — which component row changes with the behavior?
3. **Public seam** — what exact exported type/function, binary path, or file
   identity does the consumer use?
4. **Allowed edge** — is every first-party import already inward and
   generation-local?
5. **Composition root** — if touching `main.go`, `entrypoint.sh`, or
   `scripts/**`, is the change wiring-only?
6. **Operational boundary** — can the change authenticate, call Trade
   Republic, start Docker, generate code, reset state, or write portfolio
   output, and if so is it excluded from routine validation?
7. **Compatibility** — were released binary paths, Go module identities,
   auth/output paths, generated inputs, and external consumers checked?

Run:

```bash
make architecture-check
make check
```

For a stable review candidate, also run both generated architecture hook stages
from a record-inclusive disposable index and record the exact candidate
identity in `docs/architecture/adoption/verification.json`.
