---
id: 1ed791-merging-main-into-a-bug-branch-after-app
type: bug
title: merging main into a bug branch after approve cause breaks its approval
branch: fix/1ed791-merging-main-into-a-bug-branch-after-app
status: in-progress
gate: diagnose
approvals: []
---

## Brief

Owner, 2026-11-09, filing the parked follow-ups one unit at a time:

> Lets open them up as unit, one after another, so they don't cause cascading issues. Open one, fix it, PR close
> and then move to the next until all of these are cleared.

The follow-up as listed to him:

> There's no clean way to bring main into a bug branch after `approve cause`. It produces `approval-integrity`
> failures and `scope` false alarms.

Observed the same day on e6a1a0 (cause approved, gate `close`): after #32 merged, main was merged into the branch
to resolve conflicts in `lyceum/SPEC.md`, `lyceum/LIBRARY.md` and `lyceum/spec/enforcement.md`. `check --ci` then
failed `approval-integrity` on all three ("the cause approval no longer matches …: it changed after it was
approved"), and PR #34 was merged through a red `check`.

Asked how the owner should approve merged content (A: drift on a cause-hashed file at `close` sends the bug
back to `cause` with a reopen event; B: allow `approve cause` at `close`; C: ignore drift that came from main),
the owner chose:

> A

## What and why

- **Symptom:** a bug at `close` whose branch takes main in (to resolve a conflict or catch up) fails
  `approval-integrity` for every cause-hashed file main also changed. The only `Next:` is to restore the file,
  which would undo main's change. `approve cause` refuses at `close` ("unit is at gate close"), and `amend`
  refuses for bugs, so the owner has no way to approve the merged content.
- **Expected:** when a bug is at `close` (cause approved, not closed) and a file its cause approval hashed no
  longer matches, the CLI reopens it to `cause` with a `gate-reopen` event naming the files, the same mechanism
  that reopens to `fix` on a code change after verify. Code changes keep reopening to `fix` (verify must re-run
  first). Once reopened, `approval-integrity` does not fail on that drift (it is a pending re-approval, not a
  bypass); `check` points to `lyceum approve cause`, and re-approving hashes the merged content
  [DERIVED: owner chose A]. A closed unit is never reopened.
- **Trigger:** approve cause, then merge a main that edited a file the cause approval hashed (area file,
  `SPEC.md`, `LIBRARY.md`, or a code file), then run `check --ci`.
- **Not this bug:** the `scope` warning on `internal/check/check.go` seen on #34 came from a revert/reapply of
  another unit's merge on that branch, not from merging main; a branch's diff is taken from the merge base, which
  moves with each merge of main. [DERIVED: BranchChanges diffs base...HEAD]
- **Out of scope:** a merge of main that changes code still reopens the bug to `fix` (verify must re-run on the
  merged code); that is intended.

## Spec delta

## Steps

## Progress

## Log
