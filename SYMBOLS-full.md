# traderepublic-portfolio-downloader Full Symbol Reference

Workspace navigation index for the Trade Republic downloader fork.

---

## Commands and composition roots

| Path | Purpose |
|---|---|
| `cmd/portfoliodownloader/public/main.go` | Released legacy public binary composition. |
| `cmd/portfoliodownloader/dev/main.go` | Legacy development/response-writing composition. |
| `cmd/example-generator/main.go` | Deterministic example-fixture generator. |
| `v2/cmd/portfolio-downloader/main.go` | In-progress v2 operator binary composition. |
| `v2/cmd/dev/main.go` | In-progress v2 development composition. |
| `entrypoint.sh` | Container build/release entry point. |
| `scripts/generate-rest-client.sh` | Explicit legacy REST-client generation path. |

## Legacy owners

| Owner | Paths | Purpose |
|---|---|---|
| `v1_core` | Direct Go files under `internal/` | Constants, counters, and date/time values. |
| `v1_persistence` | `internal/database/`, `filesystem/`, `reader/`, `writer/` | SQLite, CSV, JSON, filesystem, reader, and writer adapters. |
| `v1_broker_transport` | `internal/console/`, `internal/traderepublc/api/` | Operator input, authentication, REST, timeline, and WebSocket transport. |
| `v1_portfolio_domain` | `internal/traderepublc/portfolio/` | Activity, document, instrument, and transaction behavior. |
| `v1_delivery` | `cmd/portfoliodownloader/app.go` | Legacy application orchestration. |

## v2 owners

| Owner | Paths | Purpose |
|---|---|---|
| `v2_foundation` | `v2/internal/const.go`, `bus/`, `console/` | Constants, event bus, and operator input primitives. |
| `v2_protocol` | `v2/pkg/traderepublic/` | Public schemas, generated protocol types, publisher, and WebSocket seam. |
| `v2_broker_transport` | `v2/internal/traderepublic/` | v2 authentication and API adapters. |
| `v2_workflows` | `v2/internal/instrument/`, `message/`, `timelinedetails/`, `timelinetransactions/`, `transaction/` | Message mapping and portfolio workflows. |
| `v2_output_adapters` | `v2/internal/file/`, `writer/` | CSV and filesystem output. |
| `v2_delivery` | v2 command app/argument packages | v2 application assembly. |

## Test and fixture surface

| Path | Purpose |
|---|---|
| `tests/architecture/` | Required architecture declaration and gate regression. |
| Root `*_test.go` files | Deterministic legacy unit tests. |
| v2 `*_test.go` files | Deterministic v2 unit tests, including visible explicit skips for unfinished behavior. |
| `tests/fakes/`, `assets/`, `v2/tests/` | Tracked test and protocol fixtures. |

## Safe validation

| Command | Notes |
|---|---|
| `make architecture-check` | Checkout-bound owner, dependency, composition-root, root-size, and baseline ratchet. |
| `make check` | Runs architecture, changed-file policy, and root plus nested v2 Go tests. It does not run a downloader, authenticate, start Docker, generate code, reset state, or build a release. |
| `make -C v2 check` | Runs the nested v2 module tests only. |

The shared `regenerate-symbols` tool currently indexes Python, TypeScript, and
Swift, not Go. These Go navigation indexes are therefore maintained manually
and verified through architecture ownership plus native Go tooling.
