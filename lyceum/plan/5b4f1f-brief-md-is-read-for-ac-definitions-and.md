---
id: 5b4f1f-brief-md-is-read-for-ac-definitions-and
type: bug
title: brief.md is read for AC definitions and counts as an area for the project gate
branch: fix/5b4f1f-brief-md-is-read-for-ac-definitions-and
status: planned
gate: repro
approvals: []
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

## Spec delta

## Steps

## Progress

## Log
