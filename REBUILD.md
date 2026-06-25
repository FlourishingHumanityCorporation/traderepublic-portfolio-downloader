# Rebuild Notes

## Safe Routine Check

```bash
make check
```

`make check` delegates to `go test -v ./...`. It is the registered workspace
validation surface and does not run the downloader binary, start Docker,
authenticate to Trade Republic, download documents, write API responses, run
`make reset`, or mutate portfolio output files.

## Live Operations

The downloader prompts for brokerage credentials and can write session files,
responses, documents, and CSV output when run for real. Treat live brokerage
auth/API/download flows as explicit opt-in operations.
