---
id: d264a8-ac-ids-are-not-checked-for-uniqueness-ac
type: bug
title: AC IDs are not checked for uniqueness across area files
branch: fix/d264a8-ac-ids-are-not-checked-for-uniqueness-ac
status: implemented
gate: close
approvals:
  - gate: cause
    hashes:
      internal/check/check.go: sha256:2f15740cce11b4b6f1428244f3cfe121acc21f60f457fe95a26f69d9d1eec651
      lyceum/LIBRARY.md: sha256:85d2a0930842a68efd64a0b3fc184a1de045ea15a99703152d0ba8daa9c519b5
      lyceum/SPEC.md: sha256:e80cabc6fb0f8be5c6a6f0cbc3225ae45ad73c95e7bd25e689e701fe493cbf65
      lyceum/plan/d264a8-ac-ids-are-not-checked-for-uniqueness-ac.md: sha256:983336d99bed0d71eee849062bdd9e328676983b3dd9663c815494f04d0b07b4
      lyceum/spec/enforcement.md: sha256:2437f83bc5915bac42b0f8b74908c97cc6e03ec7c1dc726e8da9efc1ac73f14a
    commit: 3a9bf64
    by: oritromax@gmail.com
    at: 2026-11-09T07:34:57Z
files: [internal/check/check.go]
acs: [AC198]
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

## Root cause

- `check.Registry` (`internal/check/check.go`) checked each area file only against its own registry row: an AC its
  file defines must fall in that row's ACs cell. No reader compared definitions across files, so the same ID in two
  files passed once both rows listed it. `acmap.SpecACs` dedupes IDs into a set, so coverage and `unclaimed-acs`
  hid the duplicate too.
- The omission finding's `Next:` ("add them to the row of <file>") was written per file with no knowledge of other
  files, so for an ID another area already defines it told the agent to register the duplicate.
- The repro `TestScripts/ac198_duplicate_ac_ids` defines AC1 in `billing.md` (row) and `orders.md` (bullet). With both
  rows listing AC1, `check --ci` exited 0 on the old code; it must exit 2 naming AC1 and both files. With only one
  row listing it, no hint may say to add AC1 to the other row.

## Spec delta

- `enforcement.md` **AC198** (new). Registry range in `SPEC.md` now ends AC184, AC198; `LIBRARY.md` has a bug entry
  naming it.
- `enforcement.md` check table, Registry row: also fails when two area files define the same AC ID.
- [DERIVED: naming every defining file, and local `check` warning as AC147, come from `## What and why` and the
  existing registry check, not the brief]
- [DERIVED: the hint not steering into a duplicate: the brief shows the hint producing it]
- AC198 is the owner's assigned number, not the next free one (AC196 and AC197 are taken on unmerged branches).

## Steps

## Progress

## Log

- Fix: `Registry` first reads every area file and maps each AC ID to the files defining it. An ID in more than one
  file gets one `registry` failure naming the ID and every file, and is left out of the per-row omission finding,
  so no hint says to register it. Built against `acmap.Definitions(data)` as on main.
- Repro cleared with `lyceum test --repro TestScripts/ac198_duplicate_ac_ids`: on the old code `check --ci` printed
  `exit=0`.
