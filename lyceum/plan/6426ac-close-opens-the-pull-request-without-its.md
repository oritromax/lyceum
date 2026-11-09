---
id: 6426ac-close-opens-the-pull-request-without-its
type: bug
title: close opens the pull request without its close commit
branch: fix/6426ac-close-opens-the-pull-request-without-its
status: planned
gate: repro
approvals: []
---

## Brief

Owner, 2026-11-09, filing the parked follow-ups one unit at a time:

> Lets open them up as unit, one after another, so they don't cause cascading issues. Open one, fix it, PR close
> and then move to the next until all of these are cleared.

The follow-up as listed to him:

> `lyceum close` doesn't push the close commit before opening the PR.

Observed the same day: ten units (5b4f1f, c3ede6, 6b21e1, d20234, e0b7ea, e7033e, ffa56b, 0c520e, 0e888d, a516e0)
had their PRs merged while `lyceum close`'s commit (`status: closed` + the `close` event) sat only in the local
branch, so main kept listing them at `close` until PR #33 carried the commits over by hand.

## What and why

- **Symptom:** `close` commits the closed plan and the `close` event, then, when the branch already has an
  upstream, runs `gh pr create` (or `glab mr create`) straight away. The close commit is not on the remote, so the
  PR opens without it; the only hint is the trailing `Next: push the close commit so the pull request has it`. If
  the owner merges before pushing, main never records the close.
- **Expected:** `close` never opens a pull request that lacks the close commit. When the branch has an upstream
  and the local branch is ahead of it, `close` does not run `gh`/`glab`; it prints the push and then the exact
  create command, the same shape as the no-upstream case. Lyceum still never pushes
  [DERIVED: openPR's doc comment and units.md "Close" step 5 — the CLI prints the push, it does not run it].
- **Trigger:** push the branch (e.g. to run CI) before `lyceum close`, then run `lyceum close`.
- **Out of scope:** making `close` push; detecting PRs already merged without a close (that is what PR #33
  fixed by hand).

## Spec delta

## Steps

## Progress

## Log
