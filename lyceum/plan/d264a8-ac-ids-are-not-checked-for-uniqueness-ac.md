---
id: d264a8-ac-ids-are-not-checked-for-uniqueness-ac
type: bug
title: AC IDs are not checked for uniqueness across area files
branch: fix/d264a8-ac-ids-are-not-checked-for-uniqueness-ac
status: in-progress
gate: diagnose
approvals: []
---

## Brief

Owner, 2026-11-09, approving the follow-up as filed ("yes"):

> bug: AC IDs aren't checked for uniqueness across area files. This is the real one, and it doesn't depend on the brief.

From the validation run that led to it (binary built from main at `6381e21`, throwaway repo):

> I made `brief.md` and `orders.md` both define `| AC1 | … |`. `check --ci` fails `registry` and says "add them to
> the row of lyceum/spec/brief.md". After following that, AC1 is defined in both files and `check` passes with only
> the unclaimed warning. The spec requires IDs to be "unique across all areas", but nothing enforces it.

## What and why

- **Symptom:** two area files under `lyceum/spec/` can each define the same AC ID (as a `| ACn |` row or a
  `- **ACn**` bullet), and `lyceum check --ci` reports no failing finding once each file's registry row lists it.
- **Expected:** `check` fails, naming the duplicated ID and each file that defines it, locally and in CI.
- **Trigger:** register two areas, define `| AC1 | … |` in both, add AC1 to both registry rows, run
  `lyceum check --ci`.
- **Why it matters:** `approve cause` and release coverage locate an AC by its file. A duplicate makes it ambiguous
  which definition a test or approval covers. The registry `Next:` hint even steers the agent into creating one.

## Spec delta

## Steps

## Progress

## Log
