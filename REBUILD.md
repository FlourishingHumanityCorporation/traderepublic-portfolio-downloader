# Rebuild Notes

## Safe Routine Check

```bash
make check
```

`make check` is the registered workspace validation surface. It runs the
checkout-bound architecture gate, changed-file policy, and both the root and
nested `v2` Go module tests. It
does not run a downloader binary, start Docker, authenticate to Trade
Republic, download documents, write API responses, run `make reset`, generate
clients, build release binaries, or mutate portfolio output files.

The nested module can be checked independently with:

```bash
make -C v2 check
```

Architecture ownership and dependency direction are documented in
`docs/architecture/ARCHITECTURE.md`.

## Live Operations

The downloader prompts for brokerage credentials and can write session files,
responses, documents, and CSV output when run for real. Treat live brokerage
auth/API/download flows as explicit opt-in operations.

Code generation is also explicit maintenance. The root
`scripts/generate-rest-client.sh` currently references an absent
`openapi-rest.yaml`; identify the authoritative contract before attempting to
repair or run that generator.
