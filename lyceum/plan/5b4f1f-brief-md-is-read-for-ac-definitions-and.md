---
id: 5b4f1f-brief-md-is-read-for-ac-definitions-and
type: bug
title: brief.md is read for AC definitions and counts as an area for the project gate
branch: fix/5b4f1f-brief-md-is-read-for-ac-definitions-and
status: implemented
gate: cause
approvals: []
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
