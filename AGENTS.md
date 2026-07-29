# traderepublic-portfolio-downloader

Codex guidance for this repository.

## What This Is

Repository in the CodeProjects workspace.

## Workspace Context

- Workspace path: `tools/traderepublic-portfolio-downloader`
- Repository name: `traderepublic-portfolio-downloader`
- Kind: developer/operator tool
- Registry entry: `traderepublic-portfolio-downloader` in `.meta/projects.json`

## Tech Stack Signals

- Go
- Docker

## Start Here

- Read `/Users/paulrohde/CodeProjects/AGENTS.md` before making cross-project changes.
- Before adding or moving production behavior, read
  `docs/architecture/ARCHITECTURE.md` and name the generation, owning module,
  public seam, allowed dependency edge, and composition root.
- Keep changes scoped to this repo unless the user asks for a workspace-wide change.
- Check `.meta/projects.json` before changing ports, dependency relationships, project names, or automation ownership.
- Prefer existing local conventions, scripts, and docs over new machinery.
- Never print secrets from `.env`, local credentials, browser profiles, keychains, or private data stores.

## Local Orientation

- `README.md` is the first local product/usage orientation surface.
- `docs/architecture/ARCHITECTURE.md` is the module and dependency authority.
- `docs/architecture/adoption/inventory.md` records entry points, consumers,
  resources, persisted identities, and pre-adoption debt.

## Commands

- `make architecture-check` - run the checkout-bound STD-028 module gate.
- `make check` - run architecture, changed-file policy, and both root/v2 Go
  test lanes.
- `make test` - run both root and nested v2 Go test lanes.
- `make -C v2 check` - run the nested v2 Go tests only.
- Do not run a downloader binary, Docker Compose, `make reset`, code
  generation, release builds, or response/document writing paths during
  registry or audit work unless the user explicitly approves the relevant
  live or maintenance operation.

## Important Files

- `README.md`
- `Makefile`
- `go.mod`
- `docker-compose.yml`
- `appcheck.toml`
- `REBUILD.md`
- `docs/architecture/ARCHITECTURE.md`

## Verification

- Run the narrowest relevant local check after edits.
- A new first-party import must follow a declared inward edge. A new command,
  script, or nested module must declare its composition root and join the
  normal gate in the same change.
- For user-visible behavior, prefer the real UI, CLI, appcheck, logs, or documented proof surface over a shallow green check.
- If verification cannot run, report the exact blocker and residual risk.
