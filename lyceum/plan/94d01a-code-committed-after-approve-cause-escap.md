---
id: 94d01a-code-committed-after-approve-cause-escap
type: bug
title: code committed after approve cause escapes verify and re-approval
branch: fix/94d01a-code-committed-after-approve-cause-escap
status: closed
gate: close
approvals:
  - gate: cause
    hashes:
      .agents/skills/lyceum-bug/SKILL.md: sha256:9fb9ed110dce2deefb706606d8247de03dbd9ad59c237b50e42ef4c1e336f197
      .claude/skills/lyceum-bug/SKILL.md: sha256:9fb9ed110dce2deefb706606d8247de03dbd9ad59c237b50e42ef4c1e336f197
      internal/check/check.go: sha256:259da18056f035f4e32cd3dba6cdf165dd89538881337323f81f2f7ddca6d85c
      internal/cli/reopen.go: sha256:2618bf239af1e2dbf2f0c61d5d6f1fb4c1872c445f4ce6f0d92212de2555e11d
      internal/templates/skills/bug/SKILL.md: sha256:9fb9ed110dce2deefb706606d8247de03dbd9ad59c237b50e42ef4c1e336f197
      lyceum/LIBRARY.md: sha256:a1b5f59dc9f63bd3aca8705b887b6f68572c60eba605e488080f84f73906067b
      lyceum/SPEC.md: sha256:fc9bd6ae92f0fdc4e4819f0a7eeda9c0df1d48968cf077373265e8bca399ed92
      lyceum/plan/94d01a-code-committed-after-approve-cause-escap.md: sha256:dae0e856f97a2e8b48885f5291a76dd5472dd7769ead2d7810d7bdb3039b05fa
      lyceum/spec/enforcement.md: sha256:d2882370f8a6ddf2bd69bad46846c6021412c685624c51ef8a0e2de988ca7652
      lyceum/spec/test-runner.md: sha256:f8e333ff14e484e04610ba80128c7f71a5e5b27e7be153b1027a199d11efe212
    commit: 76e3de4
    by: oritromax@gmail.com
    at: 2026-11-09T16:13:34Z
acs: [AC204]
files:
  - internal/check/check.go
  - internal/cli/reopen.go
  - internal/templates/skills/bug/SKILL.md
  - .agents/skills/lyceum-bug/SKILL.md
  - .claude/skills/lyceum-bug/SKILL.md
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

## Root cause

- `check.Reopened` (`internal/check/check.go`) returns early unless the bug's gate is `cause`. Once `approve
  cause` moves the bug to `close`, the only reopen left is `check.CauseDrift` (AC201), which compares only the
  files the cause approval hashed. A new code file or an edit to an unhashed code file is in neither, so
  `reopenBug` (`internal/cli/reopen.go`) moves nothing and `check --ci`'s `reopen` check finds nothing; only the
  `scope` warning sees the file.
- Lifting the gate guard alone would over-fire at `close`: `Reopened` lists files from `git log since..HEAD`,
  which includes commits a merge of the base branch brings in. So at `close` a file must also be in the branch's
  own committed diff (`gitx.BranchCommitted` against the base `check` resolves), the same diff `scope` uses.
- The repro `TestScripts/ac204_bug_code_at_close_reopens_fix` first fails on `check --ci` passing after a new
  code file is committed at `close`, then asserts the move to `fix` for a new file and for an unhashed edit,
  and that a code file merged in from main leaves the bug at `close`.

## Spec delta

- `enforcement.md`: AC204 added; the Reopen rule row says an open bug at `close` goes back to `fix` for code in
  the branch's own diff, not for code merged from the base branch. Registry range in `lyceum/SPEC.md` updated.
- `test-runner.md`: the after-verify reopen line says it applies at `close` too (AC204).
- Bug skill (`internal/templates/skills/bug/SKILL.md` and installed copies): the `close` paragraph names it.

## Steps

## Progress

- repro cleared (`TestScripts/ac204_bug_code_at_close_reopens_fix`).
- Fix: `Reopened` also runs for an open bug at `close`, counting only files in the branch's own committed diff;
  `reopenBug` resolves the base as `check` does and passes that diff; `check` passes `ctx.Changed`.

## Log

- `.claude/skills/lyceum-bug/SKILL.md` was not writable in the session; it needs the same edit as the template.
