---
id: df39a2-an-ac-id-defined-twice-in-one-area-file
type: bug
title: an AC ID defined twice in one area file is not flagged
branch: fix/df39a2-an-ac-id-defined-twice-in-one-area-file
status: implemented
gate: close
approvals:
  - gate: cause
    hashes:
      internal/check/check.go: sha256:0fab670d4da9c57b2301282afac160879dc03dca08eccd10b95ba32221ee734f
      lyceum/LIBRARY.md: sha256:ae42b58bdfa5b88c79784c48c267d9b97db05c6b0d57416515b065a21375e4d5
      lyceum/SPEC.md: sha256:aff5a0a1f4c6607f9f44115c6fa541d5412cbf45ca92e642dd39711b5a5571ea
      lyceum/plan/df39a2-an-ac-id-defined-twice-in-one-area-file.md: sha256:6479f37037e786d5bba854999b8c8fd44268d1f112e441e27186f0b83b339643
      lyceum/spec/enforcement.md: sha256:1b42b1b5a52f89dd51de33aaacdcb7b3559c8aab01028b6cbd169ecb24b01aac
    commit: aee0ae4
    by: oritromax@gmail.com
    at: 2026-11-09T09:51:31Z
acs: [AC202]
files:
  - internal/check/check.go
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

## Root cause

`check.Registry` (`internal/check/check.go`) builds `definedIn`, the area files defining each AC ID, from
`acmap.Definitions`. It appends a file only when it is not already listed (`!slices.Contains(in, f)`), so a
second definition in the same file is dropped, and the only duplicate finding is the cross-file one (`dups`,
AC198), which needs two distinct files. Nothing else counts definitions per file, so a repeat in one file
leaves no finding. The repro `ac202_duplicate_ac_id_in_one_file` defines AC1 twice in `orders.md` (rows,
bullets, a row and a bullet) and expects `check --ci` exit 2 with the finding; on the current code it exits 0.

## Spec delta

- AC202 (`lyceum/spec/enforcement.md`, registry range in `lyceum/SPEC.md`): an AC ID defined more than once in
  one area file fails `registry`, naming the ID and the file, `Next:` gives the repeat the next free AC ID.
  AC198 unchanged.

## Steps

## Progress

- repro cleared: `TestScripts/ac202_duplicate_ac_id_in_one_file` fails on the current code (`exit=0`).
- fix: `Registry` counts definitions per file and adds one `registry` Fail per ID a file repeats, after the
  AC198 findings; `definedIn` and the AC198 path are untouched.

## Log
