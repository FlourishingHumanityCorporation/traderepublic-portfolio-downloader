---
created: 2026-07-23
last_updated: 2026-07-29
candidate_revision: 7b8a88e6f09c52adda7f9e70aefd2be66bbd6ff9
status: APPROVED
---
# STD-028 contract review — traderepublic-portfolio-downloader

## Decision

**APPROVED — the exact candidate satisfies the zero-debt STD-028
architecture contract.**

An independent read-only reviewer approved exact implementation revision
`7b8a88e6f09c52adda7f9e70aefd2be66bbd6ff9`. It is published as an ancestor
of authoritative `fhc/main`; later project revisions change only adoption
records.

## Approved scope

- Architecture scope and Git root:
  `tools/traderepublic-portfolio-downloader`.
- Twenty-one behavior-derived owners govern all 260 scanned Go, script, test,
  documentation, automation, and control entries.
- The root and v2 Go modules remain independent generations with no
  cross-generation edge.
- Six exact composition roots and the container `entrypoint.sh` are declared.
- First-party import census reconciles every actual edge with the declared
  acyclic graph; the contract does not claim Go AST enforcement that Appcheck
  does not provide.
- Root and v2 module identities, released commands, session/auth/output paths,
  and container identity remain declared compatibility seams.
- The architecture baseline is empty with no ownership, dependency, cycle,
  root, script, composition, stale, or historical finding.

## Independent findings

- Every governed path resolves to one owner without overlap or gap.
- The non-obvious example-generator dependency on v1 test support is declared
  explicitly rather than hidden.
- The candidate contains exactly the five Go `main.go` files and one script
  listed as composition roots.
- The deleted websocket prototype contained no executable `main`, consisted of
  commented code, referenced removed packages, and had no repository or
  workspace source consumer.
- Persisted identities and public release seams are unchanged.
- Exact tree identity, artifact digests, architecture results, hook proof, and
  publication ancestry bind the evidence to `7b8a88e`.

## Accepted contract residuals

- Go edge conformance is protected by import census and review, not a
  language-aware machine edge checker.
- The human dependency-direction diagram omits the two declared test-support
  edges even though the module table, contract, and inventory carry them.
- `v2/cmd/portfolio-downloader` splits one Go package between delivery and
  composition ownership; import-level tooling cannot distinguish those files.
- The focused architecture test does not independently assert that the
  baseline is empty or that `check` depends on both test lanes; Appcheck and
  the native gate carry those proofs.

None of these residuals is hidden architecture debt.

## Evidence boundary

- Exact-revision and record-inclusive architecture runs govern 260 files with
  an empty baseline and no findings.
- `make check` passes architecture, changed-file policy, and deterministic
  tests for both Go modules.
- Both modules pass `go vet`.
- Static Appcheck passes with no warning or failure.
- Generated pre-commit and pre-push stages and the installed composite
  dispatcher pass.

No downloader, brokerage, Docker, reset, output writer, generator, or release
surface was exercised.
