---
id: 94d01a-code-committed-after-approve-cause-escap
type: bug
title: code committed after approve cause escapes verify and re-approval
branch: fix/94d01a-code-committed-after-approve-cause-escap
status: in-progress
gate: diagnose
approvals: []
---

## Brief

Owner, 2026-11-09, after the parked follow-ups were cleared:

> Lets handle the other three things you mentioned … 2 - the probe for code changes at close

Probe on main at `72b183f` (testscript built from AC201's setup: a bug through `approve cause` to `close`, then
one commit):

| Commit after approve cause | Result on main |
|---|---|
| edits a code file the cause approval hashed (`src/login.sh`) | back to `cause`, then `fix` (AC201 + reopen) |
| adds a new code file (`src/other.sh`) | stays at `close`; `warn: scope` only; `check --ci` exit 0 |
| edits a code file the approval did not hash (`t.sh`) | stays at `close`; `warn: scope` only; `check --ci` exit 0 |

Asked whether to open it as a bug with the shape "fire the after-verify reopen at `close` for any non-test change
in the branch's own diff; code that arrived only from main does not count", the owner said:

> yes

## What and why

- **Symptom:** `check.Reopened` (the "a commit after verify cleared … changes files that are not tests" rule)
  only runs while the bug is at `cause`. At `close`, only `CauseDrift` (AC201) runs, and it looks only at files
  the cause approval hashed. A commit after `approve cause` that adds a code file or edits an unhashed one
  leaves the bug at `close`: the suite never ran on that code, the owner never approved it, and CI stays green.
- **Expected:** at `close` (cause approved, not closed), a commit after the latest verify clear that changes a
  file that is neither a test nor under `lyceum/`, and that is in the branch's own diff from its merge base with
  the base branch, sends the bug back to `fix` with a `gate-reopen` event naming the files (verify must clear
  again, then the owner approves cause again). Code that reached the branch only by merging the base branch
  (not in the branch's own diff) does not reopen it. Hashed-file drift keeps AC201's behaviour. A closed unit is
  never reopened. [DERIVED: owner's "yes" to the proposed shape]
- **Trigger:** approve cause, then commit a new code file, then `lyceum status` / `check --ci`.
- **Out of scope:** changing the existing rule at `cause` (it still counts code merged from main).

## Spec delta

## Steps

## Progress

## Log
