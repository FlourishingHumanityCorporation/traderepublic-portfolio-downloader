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
- Keep changes scoped to this repo unless the user asks for a workspace-wide change.
- Check `.meta/projects.json` before changing ports, dependency relationships, project names, or automation ownership.
- Prefer existing local conventions, scripts, and docs over new machinery.
- Never print secrets from `.env`, local credentials, browser profiles, keychains, or private data stores.

## Local Orientation

- `README.md` is the first local product/usage orientation surface.

## Commands

- `make check` - run local Go tests with `go test -v ./...`.
- `make test` - same test surface as `make check`.
- Do not run the downloader binary, Docker Compose, `make reset`, or response/document writing paths during registry or audit work unless the user explicitly approves live brokerage operations.

## Important Files

- `README.md`
- `Makefile`
- `go.mod`
- `docker-compose.yml`
- `appcheck.toml`
- `REBUILD.md`

## Verification

- Run the narrowest relevant local check after edits.
- For user-visible behavior, prefer the real UI, CLI, appcheck, logs, or documented proof surface over a shallow green check.
- If verification cannot run, report the exact blocker and residual risk.
