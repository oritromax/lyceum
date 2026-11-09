---
id: 1ed791-merging-main-into-a-bug-branch-after-app
type: bug
title: merging main into a bug branch after approve cause breaks its approval
branch: fix/1ed791-merging-main-into-a-bug-branch-after-app
status: implemented
gate: close
approvals:
  - gate: cause
    hashes:
      .agents/skills/lyceum-bug/SKILL.md: sha256:77f9e4e7aa7d1f13d5be98d2460c027dae7409e047615ff482313ffef25c9dfb
      .claude/skills/lyceum-bug/SKILL.md: sha256:77f9e4e7aa7d1f13d5be98d2460c027dae7409e047615ff482313ffef25c9dfb
      internal/check/check.go: sha256:e190ff02379a00c519369dc0c4c62497e8797be045b1e873bea4be6064122ab9
      internal/cli/approve.go: sha256:4fcedd42e032d4dd635d93da1d39089e979fb68dfa14968883dc3d30ff4971b4
      internal/cli/check.go: sha256:b3a3d2a5b117b21f3d615eabe552682c2e1301a0c1d8a6af193e51820a9b344e
      internal/cli/reopen.go: sha256:441c00ed2e6cbd9751e0e8e8df53a8645bfdd0edfa2d70d2bb22d770b2644bf3
      internal/cli/status.go: sha256:2bdf0c4889e3f597849ac4cad5d0a568c167be4e0426001afe7f0c5ee3e313e5
      internal/cli/test.go: sha256:3a7f2d9556596b27a8e8dd32bcaa0adae8b2eee2ee2bf4a9a94a304cfd6dcbf4
      internal/templates/skills/bug/SKILL.md: sha256:77f9e4e7aa7d1f13d5be98d2460c027dae7409e047615ff482313ffef25c9dfb
      lyceum/LIBRARY.md: sha256:b160397596959699454e561ea3d5c6eb81b7e976c5fd6047b57faed46f4fb83b
      lyceum/SPEC.md: sha256:5f92c3f525310e438dba23ae23f3d71d32d943a22a148774e4f9f86c25549617
      lyceum/plan/1ed791-merging-main-into-a-bug-branch-after-app.md: sha256:0a553165e3eb54259e8821bcd088a2a680185c79d2ba79cb24c0bca9ef7e9e9d
      lyceum/spec/enforcement.md: sha256:7ee75a889f63e75fed19e67b0f02734387a09d83fa0545e843b854644e4d0f45
    commit: 5e841f6
    by: oritromax@gmail.com
    at: 2026-11-09T09:17:53Z
acs: [AC201]
files:
  - internal/check/check.go
  - internal/cli/reopen.go
  - internal/cli/approve.go
  - internal/cli/check.go
  - internal/cli/status.go
  - internal/cli/test.go
  - internal/templates/skills/bug/SKILL.md
  - .agents/skills/lyceum-bug/SKILL.md
  - .claude/skills/lyceum-bug/SKILL.md
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

## Root cause

- The only reopen for a bug is `check.Reopened`, which fires at gate `cause` (before approval) on a code commit
  after verify, and sends the bug to `fix`. Nothing handles a bug at `close`: `reopenBug` (called by `status`,
  `test`, `approve` and local `check`) returns early, and the AC17 replay accepts a `gate-reopen` only from `cause`
  to `fix`.
- So at `close`, a changed cause-hashed file falls through to `ApprovalIntegrity`, whose only `Next` is "restore",
  and the old `approve cause` keeps counting (`CauseApproval` sees it), so no step asks for a new approval. With
  `approve cause` refused at `close` and `amend` refused for bugs, there is no way out.
- The premise "code changes at close already reopen to fix" does not hold on main: `Reopened` checks
  `gate == cause` only, so at `close` a code change also lands on `approval-integrity`
  [DERIVED: read from internal/check/check.go]. The fix covers it by composition: drift on a hashed code file
  reopens to `cause`, then the existing rule sends it on to `fix`.
- The repro (`TestScripts/ac201_bug_merge_reopens_cause`) approves cause with `lyceum/SPEC.md` in the hashes,
  commits a change to it, and expects `status` to print `Back to cause`; on main, `status` stays at `close` and
  lists "restore lyceum/SPEC.md" as the pending step.

## Spec delta

- AC201 added to `lyceum/spec/enforcement.md` and its registry range in `lyceum/SPEC.md`; the Reopen check row
  describes the reopen to `cause`. The bug skill's `close` section says what to do after it.

## Steps

## Progress

## Log
