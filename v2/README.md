# Trade Republic Portfolio Downloader v2

`v2` is the in-progress replacement for the legacy root Go module. It has its
own `go.mod`, command composition roots, internal workflows, generated Trade
Republic protocol types, and deterministic test lane.

## Current entry points

- `cmd/portfolio-downloader/main.go` composes the operator downloader.
- `cmd/dev/main.go` composes the development fixture workflow.

The former `cmd/websocket-downloader` tree was an entirely commented-out
prototype that referenced packages and Make targets which no longer exist. It
was not executable and has been retired rather than retained as a misleading
second operator path.

## Owners

- `pkg/traderepublic/**` owns generated protocol types, schemas, and the public
  WebSocket publisher/client seam.
- `internal/traderepublic/**` owns authentication and API adapters.
- `internal/{instrument,message,timelinedetails,timelinetransactions,transaction}`
  owns portfolio-processing workflows.
- `internal/{file,writer}` owns output adapters.
- `internal/{bus,console}` plus `internal/const.go` owns shared runtime
  primitives.
- `cmd/**` owns composition and operator/development delivery only.

The complete allowed dependency graph and evolution rules live in
`../docs/architecture/ARCHITECTURE.md`.

## Safe routine check

From the repository root:

```bash
make check
```

Or for this nested module alone:

```bash
make -C v2 check
```

These commands compile and test both Go modules without running either
downloader, authenticating to Trade Republic, starting Docker, or writing
portfolio responses/documents.

## Generated sources

`make -C v2 generate` rewrites checked-in generated protocol code and may
access the network through Go tools. It is a maintenance operation, not part
of the routine gate.

## License

Same as the main project.
