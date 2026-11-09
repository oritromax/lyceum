---
id: 6426ac-close-opens-the-pull-request-without-its
type: bug
title: close opens the pull request without its close commit
branch: fix/6426ac-close-opens-the-pull-request-without-its
status: in-progress
gate: diagnose
approvals: []
files: [internal/cli/pr.go, internal/gitx/gitx.go, cmd/lyceum/testdata/script/ac034_close_opens_pr.txtar, cmd/lyceum/testdata/script/ac126_close_prints_push_merge.txtar, internal/templates/skills/bug/SKILL.md, internal/templates/skills/chore/SKILL.md, internal/templates/skills/doc/SKILL.md, internal/templates/skills/feat/SKILL.md, internal/templates/skills/refactor/SKILL.md, internal/templates/skills/test/SKILL.md, .agents/skills/lyceum-bug/SKILL.md, .agents/skills/lyceum-chore/SKILL.md, .agents/skills/lyceum-doc/SKILL.md, .agents/skills/lyceum-feat/SKILL.md, .agents/skills/lyceum-refactor/SKILL.md, .agents/skills/lyceum-test/SKILL.md, README.md, docs-site/src/content/docs/concepts/units-and-gates.md, docs-site/src/content/docs/start/quickstart.md]
acs: [AC200]
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

## Root cause

- `openPR` (`internal/cli/pr.go`), close's last step, was written to the rule "Lyceum never pushes". `runClose`
  commits the closed plan and the `close` event locally, then `openPR` either prints the push and stops (no
  upstream) or, when the branch has an upstream, runs `gh pr create` / `glab mr create` at once. Nothing between
  the close commit and the create step puts that commit on the remote, so the PR's head is the pre-close commit;
  the only trace is the `Next: push the close commit…` line, which is easy to miss before merging.
- The repro, `TestScripts/ac200_close_pushes_before_pr`, uses a local bare origin at a path naming `github.com`
  and a fake `gh` that records origin's commit on the `--head` branch when it runs. On the old code `gh` sees the
  branch's last pre-close commit (`guide`), not `lyceum: close …`, so the first assertion fails.

## Spec delta

- `lyceum/spec/units.md`: new **AC200** (close pushes the unit branch, then opens the PR; push failure and no-`gh`
  cases); "Close" step 5 rewritten to say `close` pushes [DERIVED: owner chose B].
- **AC126** text edited minimally: "prints the push and merge commands" was now false when the push succeeds; it
  now says it prints the merge commands, and the push only when it did not push. AC34 is unchanged (still true).
- `lyceum/SPEC.md`: units range gains AC200. `lyceum/LIBRARY.md`: a line naming AC200.

## Steps

- Repro: `ac200_close_pushes_before_pr.txtar`, committed alone; `lyceum test --repro` cleared.
- Fix: `gitx.PushArgs` / `gitx.Push` (2-minute `PushTimeout`, own `runFor`); `openPR` pushes after the close
  commit and before the create step; push failure keeps the unit closed and prints push + create.
- AC34 / AC126 scripts: origins switched to local bare repos (they would now push to the network) and the
  old "does not push" assertions replaced; these encoded the reversed behaviour.
- Docs: skills (templates + `.agents` copies), README, docs-site quickstart and units-and-gates.

## Progress

- [x] repro
- [x] diagnose
- [x] fix
- [ ] verify

## Log

- 2026-11-09: `.claude/skills/lyceum-*` copies not updated: this session had no write permission under
  `.claude/`. They now differ from the templates by the close sentence; the owner syncs them.
