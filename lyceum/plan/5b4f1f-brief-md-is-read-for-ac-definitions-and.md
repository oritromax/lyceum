---
id: 5b4f1f-brief-md-is-read-for-ac-definitions-and
type: bug
title: brief.md is read for AC definitions and counts as an area for the project gate
branch: fix/5b4f1f-brief-md-is-read-for-ac-definitions-and
status: implemented
gate: close
approvals:
  - gate: cause
    hashes:
      internal/acmap/acmap.go: sha256:87d60525fe8dffa70dc3b74ed36de6aff3c004dbe7f08500785eb93a077c31ec
      internal/approval/approval.go: sha256:d42361433e795553ae2e0581938c6139af11ea0d6a96d090a3f2c172de537920
      internal/check/check.go: sha256:eb7758ee56230d0cba5e89cbaab8363d5f903a1fb4240006d95708276e79ee9a
      internal/cli/approve.go: sha256:6048ce27ea168185c1f6e7bfb21f325a1f441ab70467d5d3d1db12e13912e74a
      lyceum/LIBRARY.md: sha256:b3c2bdb8ca7688ac5721818938c6a03a61a740031d72bd7d102ef04a78221927
      lyceum/SPEC.md: sha256:70f697d337a3104d43fd8a9bc4393a28687e9bb466b4cd2761dcbc5a4840f705
      lyceum/plan/5b4f1f-brief-md-is-read-for-ac-definitions-and.md: sha256:a12ce20bb8a6567f1f7a4cae0ad096a74f276ffd447cfff5680247a4581825be
      lyceum/spec/units.md: sha256:5ece0282654d8328c8ecb62e912f49ef55b48c83c964bd838ad698b825de8f76
    commit: 9ae7ced
    by: oritromax@gmail.com
    at: 2026-11-09T06:31:49Z
files: [internal/acmap/acmap.go, internal/acmap/acmap_test.go, internal/approval/approval.go, internal/check/check.go, internal/cli/approve.go]
acs: [AC196]
---

## Brief

Owner, 2026-11-09, approving the follow-up as filed ("yes"):

> bug: ACs aren't read from `brief.md`, and `brief` doesn't count as an area for the project gate. It's one small
> change in `acmap`/`projectPending` that closes both brief holes.

From the validation run that led to it (binary built from main at `6381e21`, throwaway repos):

> A greenfield repo with only `brief` registered: `check --ci` passes, `status` shows `project: pending`, so a
> project with zero ACs can be approved. A brief that quotes `| AC1 | … |` fails `registry`, and the `Next:` line
> says to add AC1 to brief's row.

## What and why

- **Symptom 1:** `lyceum/spec/brief.md` is the owner's words kept verbatim (AC185, AC190), but `acmap.Definitions`
  reads it like any area file. A brief that quotes `| ACn |` rows or `- **ACn**` bullets therefore "defines" ACs,
  `check` fails `registry` on them, and the `Next:` hint steers the agent into registering them on brief's row.
- **Symptom 2:** `projectPending` counts `brief` as an area. A greenfield repo with only `brief` registered shows
  `project: pending` and can be approved with no AC anywhere, which defeats AC38's "at least one area".
- **Expected:** the brief defines no ACs, whatever it quotes; the project gate is pending only when at least one
  area other than `brief` is registered. `approve project` still hashes `brief.md` (AC191 unchanged).
- **Trigger:** greenfield init; register only `brief`; `lyceum status` (symptom 2). Then add `| AC1 | … |` to
  `brief.md`; `lyceum check --ci` (symptom 1).

## Root cause

- Nothing in the code knew that `brief.md` is not an AC-bearing area. AC190 registers it as area `brief` so
  `approve project` hashes it (AC191). From then on every reader treated it like any area file.
- `acmap.Definitions(data)` took only the text, so it could not tell the brief from an area. Every AC reader goes
  through it: `acmap.SpecACs`/`OpenSpecACs` (coverage, `unclaimed-acs`), `approval.ACFile` (approve cause's AC file
  lookup) and `check.Registry`. A quoted `| ACn |` row or `- **ACn**` bullet in the brief became a definition, so
  `registry` failed on brief's empty ACs cell and its `Next:` said to add them to brief's row.
- `projectPending` (`internal/cli/approve.go`) and `approval.FilesFor` for `project` only checked that the registry is
  non-empty. With only `brief` registered, `project` was pending and could be approved with no AC at all,
  against AC38's "at least one area".
- The repro `TestScripts/ac196_brief_defines_no_acs` registers only `brief` (quoting AC1 and AC2) and asserts that
  `project` is not pending, `approve project` refuses, and `check --ci` passes without naming AC1/AC2. On the old code
  it failed at the first assertion: `status` printed `project: pending`.

## Spec delta

- `units.md` "Checkpoints": the brief is not an AC-bearing area, and doesn't count for the project gate.
- `units.md` **AC196** (new). Registry range in `SPEC.md` is now AC185–AC196; `LIBRARY.md` has a bug entry naming it.
- [DERIVED: `approve project` refusing with only `brief` registered goes beyond the brief, which names only the
  pending state; without the refusal AC38 could still be defeated by approving anyway]

## Steps

## Progress

## Log

- Fix: `acmap.IsBrief` is the one rule. `Definitions` now takes the file name and returns nothing for the brief, so
  all four AC readers ignore it with no caller special-casing. `projectPending` and `FilesFor` (project) need an
  area other than the brief, through the same `IsBrief`. `approve project` still hashes `brief.md` (AC191).
- First `lyceum test --repro ac196_brief_defines_no_acs` ran no test (`-run` matches top-level names), so it did not
  clear. Re-ran as `TestScripts/ac196_brief_defines_no_acs` with the test unchanged: it failed and cleared repro.
