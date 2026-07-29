---
created: 2026-07-23
last_updated: 2026-07-28
candidate_revision: 7b8a88e6f09c52adda7f9e70aefd2be66bbd6ff9
status: PENDING_INDEPENDENT_REVIEW
---
# STD-028 contract review — traderepublic-portfolio-downloader

## Decision

**PENDING — no independent reviewer has approved this candidate.**

The implementation owner completed the following self-audit to prepare a
bounded review. It is evidence for a reviewer, not independent approval.

## Review scope

- Exact tested revision:
  `7b8a88e6f09c52adda7f9e70aefd2be66bbd6ff9`.
- Original adoption base:
  `29fdf803aab8bc886af75a0073e2542d31e27b41`.
- Reconciled canonical `fhc/main`:
  `044f6235d61d96f075218a4a9e6ea63c4f572606`.
- Architecture scope and Git root:
  `tools/traderepublic-portfolio-downloader`.
- The candidate includes the original adoption, the two newer canonical
  guidance/hook commits, and the one-time whitespace/line-ending normalization
  required for generated hooks to pass on all tracked files.

## Owner self-audit

- The architecture marker is required in `appcheck.toml`, and the baseline is
  an exact resource with one owner.
- Every scanned file has exactly one owner; the record-inclusive census is
  260/260 with zero outside governance.
- v1 and v2 are independent Go modules with no declared cross-generation edge.
- Six exact composition roots and one exact container entry point are named.
- The root native gate binds Appcheck to the invoking checkout and runs both Go
  modules; the former root-only coverage gap is closed.
- The dead websocket prototype had no executable `main`, referenced removed
  packages/targets, and had no repository or workspace source consumer.
- Released command paths, Go module identities, auth/output paths, and live
  operation boundaries are retained.
- The architecture baseline is empty; no adopter-created debt is hidden.

## Required reviewer checks

1. Confirm the owner split follows behavior/change cohesion rather than folder
   quotas.
2. Confirm every allowed edge is inward and every important Go import is
   represented honestly despite Appcheck not providing Go AST dependency
   proof.
3. Confirm deletion of `v2/cmd/websocket-downloader/**` is compatible with the
   public fork and release expectations.
4. Confirm the root and v2 test lanes are sufficient non-live compatibility
   proof for this adoption tranche.
5. Confirm the pointer assertion repair changes only the test contract.
6. Confirm generated-hook and installed-hook evidence remain distinct.

## Blocking evidence

- The exact implementation revision `7b8a88e` and its initial evidence commit
  `0526d3b` are published to canonical `fhc/main`; the implementation is an
  ancestor of the published evidence.
- Independent review is still absent for exact revision `7b8a88e`.
- The root REST generator is already broken by an absent
  `openapi-rest.yaml`; it was not run or claimed green.
- Twenty-three v2 test events remain explicitly skipped.

The generated pre-commit and composite pre-push stages both pass, and their
dispatchers are installed. This publication note records integration evidence
only; it is not independent approval. The fleet ledger must remain `IN_REVIEW`
until an independent reviewer records approval for this exact candidate
identity.
