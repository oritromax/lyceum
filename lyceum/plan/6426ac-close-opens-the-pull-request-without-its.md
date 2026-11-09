---
id: 6426ac-close-opens-the-pull-request-without-its
type: bug
title: close opens the pull request without its close commit
branch: fix/6426ac-close-opens-the-pull-request-without-its
status: implemented
gate: close
approvals:
  - gate: cause
    hashes:
      .agents/skills/lyceum-bug/SKILL.md: sha256:6b8490f24941601123bee0beb16eecb6431ad73fab32db56433afe5e8d39144b
      .agents/skills/lyceum-chore/SKILL.md: sha256:58228a5846c2ea1605d85e1b4870c636d3512e930c229930d13274aa14ade4a9
      .agents/skills/lyceum-doc/SKILL.md: sha256:ccd67dbe2fc7f6d27b1d7c7de072cb40cd41f41234788c74476978b106039b50
      .agents/skills/lyceum-feat/SKILL.md: sha256:5ee83ac21d335c3c38b4123076228a5d78fdb91380dcef70cacdb012598aaacc
      .agents/skills/lyceum-refactor/SKILL.md: sha256:0236a09de50b9ab29de63fa29f098f4a1530bc31a9ff014f63fe12171adc9a51
      .agents/skills/lyceum-test/SKILL.md: sha256:02ad04b0ac94d735db3561a1c460b830210b6a8a2e1a3cad212d4bb985f3415b
      .claude/skills/lyceum-bug/SKILL.md: sha256:6b8490f24941601123bee0beb16eecb6431ad73fab32db56433afe5e8d39144b
      .claude/skills/lyceum-chore/SKILL.md: sha256:58228a5846c2ea1605d85e1b4870c636d3512e930c229930d13274aa14ade4a9
      .claude/skills/lyceum-doc/SKILL.md: sha256:ccd67dbe2fc7f6d27b1d7c7de072cb40cd41f41234788c74476978b106039b50
      .claude/skills/lyceum-feat/SKILL.md: sha256:5ee83ac21d335c3c38b4123076228a5d78fdb91380dcef70cacdb012598aaacc
      .claude/skills/lyceum-refactor/SKILL.md: sha256:0236a09de50b9ab29de63fa29f098f4a1530bc31a9ff014f63fe12171adc9a51
      .claude/skills/lyceum-test/SKILL.md: sha256:02ad04b0ac94d735db3561a1c460b830210b6a8a2e1a3cad212d4bb985f3415b
      README.md: sha256:ce0801936da2ac10a3c760662970f1c916998294e3f3bba33d4415cc28d90869
      docs-site/src/content/docs/concepts/units-and-gates.md: sha256:deb6b70406dee1d17888e5dd0f06977330078ce087369f954f3c5a1228449112
      docs-site/src/content/docs/start/quickstart.md: sha256:4a42e1c3c36c2d7f6b0aad3dcc5519098d7f025078fea40baab0bd0932b3a338
      internal/cli/pr.go: sha256:6d8c549c597d03f576c8f059b5e6b02aebfd71e66425138124df3d0fe0e0d0ef
      internal/gitx/gitx.go: sha256:0f07cff14ba5b14e6a1394141bc9cbeb97f8a8f26cc7a42b9ba3d64831a9cc2d
      internal/templates/skills/bug/SKILL.md: sha256:6b8490f24941601123bee0beb16eecb6431ad73fab32db56433afe5e8d39144b
      internal/templates/skills/chore/SKILL.md: sha256:58228a5846c2ea1605d85e1b4870c636d3512e930c229930d13274aa14ade4a9
      internal/templates/skills/doc/SKILL.md: sha256:ccd67dbe2fc7f6d27b1d7c7de072cb40cd41f41234788c74476978b106039b50
      internal/templates/skills/feat/SKILL.md: sha256:5ee83ac21d335c3c38b4123076228a5d78fdb91380dcef70cacdb012598aaacc
      internal/templates/skills/refactor/SKILL.md: sha256:0236a09de50b9ab29de63fa29f098f4a1530bc31a9ff014f63fe12171adc9a51
      internal/templates/skills/test/SKILL.md: sha256:02ad04b0ac94d735db3561a1c460b830210b6a8a2e1a3cad212d4bb985f3415b
      lyceum/LIBRARY.md: sha256:9c526672cbaef3c922797cc6845f991c9a34b18f251903fdf4ee5718716f9bea
      lyceum/SPEC.md: sha256:1066c8b105164918a80dcbdc194167107ea594c7b84e989ef800c2bd9438da5e
      lyceum/plan/6426ac-close-opens-the-pull-request-without-its.md: sha256:dcd58cd3d085693cea8e66163efef23f48aa95b4f8e5b3fbba1acf0b7fdf40d5
      lyceum/spec/units.md: sha256:99360cfb6dd61dfbda222993c52c0a291642e151e43dcf14b12b8110da43053d
    commit: 25b99a0
    by: oritromax@gmail.com
    at: 2026-11-09T08:40:11Z
files: [internal/cli/pr.go, internal/gitx/gitx.go, cmd/lyceum/testdata/script/ac034_close_opens_pr.txtar, cmd/lyceum/testdata/script/ac126_close_prints_push_merge.txtar, internal/templates/skills/bug/SKILL.md, internal/templates/skills/chore/SKILL.md, internal/templates/skills/doc/SKILL.md, internal/templates/skills/feat/SKILL.md, internal/templates/skills/refactor/SKILL.md, internal/templates/skills/test/SKILL.md, .agents/skills/lyceum-bug/SKILL.md, .agents/skills/lyceum-chore/SKILL.md, .agents/skills/lyceum-doc/SKILL.md, .agents/skills/lyceum-feat/SKILL.md, .agents/skills/lyceum-refactor/SKILL.md, .agents/skills/lyceum-test/SKILL.md, .claude/skills/lyceum-bug/SKILL.md, .claude/skills/lyceum-chore/SKILL.md, .claude/skills/lyceum-doc/SKILL.md, .claude/skills/lyceum-feat/SKILL.md, .claude/skills/lyceum-refactor/SKILL.md, .claude/skills/lyceum-test/SKILL.md, README.md, docs-site/src/content/docs/concepts/units-and-gates.md, docs-site/src/content/docs/start/quickstart.md]
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
