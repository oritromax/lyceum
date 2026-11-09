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

Asked whether `close` should refuse to open the PR until he pushes (A) or push the branch itself and then open
it (B), the owner chose:

> B, cause that was my initial idea

## What and why

- **Symptom:** `close` commits the closed plan and the `close` event, then, when the branch already has an
  upstream, runs `gh pr create` (or `glab mr create`) straight away. The close commit is not on the remote, so the
  PR opens without it; the only hint is the trailing `Next: push the close commit so the pull request has it`. If
  the owner merges before pushing, main never records the close.
- **Expected:** `close` pushes the unit branch itself after the close commit and before it opens the pull
  request, so the PR always carries the close: `git push -u origin <branch>` when the branch has no upstream,
  `git push` to its upstream otherwise, then `gh pr create` / `glab mr create` as today. With no `gh`/`glab`
  it still pushes, then prints the create step. If the push fails, the unit stays closed (the commit is local),
  `close` opens no PR and prints the exact push and create commands to run by hand. With no origin remote,
  nothing changes. This reverses "Lyceum never pushes" for `close` only; no other command pushes
  [DERIVED: owner chose B].
- **Trigger:** push the branch (e.g. to run CI) before `lyceum close`, then run `lyceum close`.
- **Out of scope:** force pushes or pushing any branch other than the unit's own; detecting PRs already merged without a close (that is what PR #33
  fixed by hand).

## Spec delta

## Steps

## Progress

## Log
