---
id: df39a2-an-ac-id-defined-twice-in-one-area-file
type: bug
title: an AC ID defined twice in one area file is not flagged
branch: fix/df39a2-an-ac-id-defined-twice-in-one-area-file
status: in-progress
gate: diagnose
approvals: []
---

## Brief

Owner, 2026-11-09, filing the parked follow-ups one unit at a time:

> Lets open them up as unit, one after another, so they don't cause cascading issues. Open one, fix it, PR close
> and then move to the next until all of these are cleared.

The follow-up as listed to him (left over from d264a8):

> Duplicate AC IDs inside a single area file aren't caught.

Probe on main at `d6362ab`: an area file `lyceum/spec/orders.md` with two `| AC1 | … |` rows, registered with
`AC1`, gives `lyceum check --ci` exit 0 and only `warn: unclaimed-acs: … AC1`.

## What and why

- **Symptom:** `check.Registry` records each defining file once per ID (`!slices.Contains(in, f)`), so a second
  definition of an ID in the same file is dropped silently. Two different criteria then share one ID, and a
  test claiming it covers only one of them.
- **Expected:** `check` fails (`--ci` exit 2, reported locally) when one area file defines an AC ID more than
  once, naming the ID and the file, with a `Next:` to give the repeat the next free ID. It applies to both
  definition forms `acmap.Definitions` reads (`| ACn |` rows and `- **ACn**` bullets), including one of each.
  The cross-file finding from AC198 is unchanged. [DERIVED: same severity and shape as AC198's finding]
- **Trigger:** define `| AC1 |` twice in one area file, register AC1, run `check --ci`.
- **Out of scope:** `brief.md` (defines no ACs, AC196); IDs inside code fences.

## Spec delta

## Steps

## Progress

## Log
