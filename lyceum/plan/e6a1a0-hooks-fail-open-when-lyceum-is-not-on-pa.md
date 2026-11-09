---
id: e6a1a0-hooks-fail-open-when-lyceum-is-not-on-pa
type: bug
title: hooks fail open when lyceum is not on PATH after npx init
branch: fix/e6a1a0-hooks-fail-open-when-lyceum-is-not-on-pa
status: in-progress
gate: diagnose
approvals: []
---

## Brief

Owner, 2026-11-09:

> i noticed another problem, if we initialize via npx lyceum init, the binary isn't there for further cheks, so
> claude code can't firectly find it. Not sure if i am missing something here.
>
> ❯ /lyceum-understand Understand what this project does
>   PreToolUse:Bash hook error
>   Failed with non-blocking status code: /bin/sh: 1: lyceum: not found
>
> I ran it on nas, on the jai repo and i initialized it via npx lyceum init

> open a bug unit for this, i will handle the binary installation

Confirmed on the NAS repo `~/work/office/jai`: `lyceum` is not on PATH; the only copy is the npx cache
(`~/.npm/_npx/<hash>/node_modules/.bin/lyceum`). `init` wrote bare `lyceum` into `.claude/settings.json`
(`lyceum hook pre-tool --harness claude`) and `.git/hooks/pre-commit` (`lyceum check --pre-commit || true`).

## What and why

- **Symptom:** after `npx lyceum init` with no global install, every hook `init` wrote calls a `lyceum` that is not
  on PATH. The shell exits 127. Claude Code blocks a tool call only on exit 2, so it reports a non-blocking error
  and the call goes ahead: gated edits, Bash and `lyceum approve` checks are all unenforced. The pre-commit hook
  fails the same way and its `|| true` hides it. Only the CI Action (which installs its own binary) still enforces.
- **Expected:** a project `init` wired never runs unenforced silently. The harness hooks fail closed with a message
  naming the install command when `lyceum` cannot be found; `init` warns, naming the install command, when it ran
  from a binary that will not be on PATH afterwards (npx); the docs say `npx lyceum init` needs a global install
  for the hooks to work.
- **Trigger:** a PATH without `lyceum`; run the hook command `init` writes for Claude Code with a PreToolUse
  payload for an edit to a gated file; observe exit 127, not a deny.
- **Out of scope:** writing the npx cache's absolute path into hooks (npx prunes it, and versions drift).

## Spec delta

## Steps

## Progress

## Log
