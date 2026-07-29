---
created: 2026-07-23
last_updated: 2026-07-23
candidate_tree: 70317206f6b1ace6616ca39dac99367e957dde96
status: PENDING_INDEPENDENT_REVIEW
---
# STD-028 contract review — traderepublic-portfolio-downloader

## Decision

**PENDING — no independent reviewer has approved this candidate.**

The implementation owner completed the following self-audit to prepare a
bounded review. It is evidence for a reviewer, not independent approval.

## Review scope

- Candidate tree excluding adoption records:
  `70317206f6b1ace6616ca39dac99367e957dde96`.
- Base and current canonical `fhc/main`:
  `29fdf803aab8bc886af75a0073e2542d31e27b41`.
- Architecture scope and Git root:
  `tools/traderepublic-portfolio-downloader`.
- Stable diff: 13 tracked paths changed, including two deletions, plus five
  non-record candidate files.

## Owner self-audit

- The architecture marker is required in `appcheck.toml`, and the baseline is
  an exact resource with one owner.
- Every scanned file has exactly one owner; the final record-inclusive census
  is expected to be 259/259 with zero outside governance.
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

- The primary checkout has neither an installed pre-commit hook nor an
  installed composite pre-push hook.
- The candidate is detached, uncommitted, and unpublished.
- The shared control-plane candidate is dirty and behind upstream.
- The root REST generator is already broken by an absent
  `openapi-rest.yaml`; it was not run or claimed green.
- Twenty-three v2 test events remain explicitly skipped.

The fleet ledger must remain `IN_REVIEW` until an independent reviewer records
an approval for this exact candidate identity and installed-hook proof is
conclusive.
