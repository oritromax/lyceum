---
id: e6a1a0-hooks-fail-open-when-lyceum-is-not-on-pa
type: bug
title: hooks fail open when lyceum is not on PATH after npx init
branch: fix/e6a1a0-hooks-fail-open-when-lyceum-is-not-on-pa
status: implemented
gate: cause
approvals: []
files: [internal/cli/initharness.go, internal/cli/init.go, internal/templates/templates.go, internal/templates/pre-commit, internal/templates/opencode-plugin.js, internal/templates/pi-extension.ts, internal/templates/templates_test.go, cmd/lyceum/testdata/script/ac052_init_hooks_ci.txtar, cmd/lyceum/testdata/script/ac053_init_wiring.txtar, cmd/lyceum/testdata/script/ac092_init_rerun.txtar, cmd/lyceum/testdata/script/ac197_init_and_precommit_without_lyceum.txtar, README.md, docs-site/src/content/docs/start/install.md, docs-site/src/content/docs/index.mdx, .claude/settings.json, .codex/hooks.json]
acs: [AC197]
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

## Root cause

- Every hook `init` writes assumes `lyceum` is on PATH and treats "not found" as no answer, which harnesses read
  as allow:
  - `initHarness.command()` (`internal/cli/initharness.go`) wrote the bare `lyceum hook pre-tool --harness <name>`
    into Claude Code, Codex and Command Code settings. With no `lyceum`, `sh` exits 127; these harnesses block only
    on exit 2 (or an explicit deny answer), so the call goes ahead unchecked.
  - The OpenCode and Pi shims (`internal/templates/opencode-plugin.js`, `pi-extension.ts`) caught the empty stdout
    of a failed `spawnSync` and returned no answer, by design ("as when lyceum is not on PATH: the call goes ahead").
  - The pre-commit block (`internal/templates/pre-commit`) ran `lyceum check --pre-commit || true`; `|| true` hides
    the 127 behind the shell's bare "not found".
- Nothing told the owner: `init` never checks whether `lyceum` will be on PATH once it ends, and run from npx it
  only is while npx runs.
- The repro `TestScripts/ac197_hook_fails_closed_without_lyceum` runs each JSON harness's command from the
  settings file `init` wrote, as the harness does (`sh -c`), with a PATH holding only `sh`, and requires exit 2 with
  the install commands on stderr. On the old code it exits 127 with `lyceum: not found`.

## Fix

- JSON hooks: the command is `command -v lyceum >/dev/null 2>&1 || { echo '<message>' >&2; exit 2; }; lyceum hook
  pre-tool --harness <name>`. It still contains `lyceum hook pre-tool`, so `mergeHook` replaces an old bare entry on
  re-run (ac092) and eject's `unmergeHook`/`wired()` recognise it unchanged.
- Shims: on `spawnSync` error `ENOENT` they deny with the message (OpenCode: deny answer, so `execute.before`
  throws and the shell permission is `deny`; Pi: `{block: true, reason}`).
- Pre-commit: `if command -v lyceum …; then lyceum check --pre-commit || true; else echo '<message>' >&2; fi`, still
  exit 0. Doctor compares the block with the template, so it reports an old block until `init` reruns.
- `init` warns on stderr, naming the install commands, when no `lyceum` is on PATH outside npx's `/_npx/` cache.
- `templates.InstallHint` holds the install text; `TestInstallHint_AC197` keeps the pre-commit block and shims in
  step with it.
- README.md, docs `start/install.md` and the docs home page say a global install is needed for the hooks.

## Spec delta

- `lyceum/spec/enforcement.md`: AC197 (new); "Harness hooks" gains the fail-closed rule and `init`'s warning;
  "`lyceum check`" gains the pre-commit message. Registry range in `lyceum/SPEC.md` gains AC197; `lyceum/LIBRARY.md`
  names it.
- AC197 [DERIVED: plan What and why] Given the hooks `init` wrote and no `lyceum` on PATH, when a harness runs its
  hook, then the Claude Code, Codex and Command Code hooks exit 2 and the OpenCode and Pi shims deny, each with a
  message naming the install commands, and the pre-commit hook says it did not check the commit, names them and
  exits 0; given `init` run where no `lyceum` will be on PATH afterwards (only npx's cache, or none), then it warns
  naming them.
- [AGENT-PROPOSED: Codex and Command Code are taken to block on exit 2 as Claude Code does; not verified live]
- [AGENT-PROPOSED: a `lyceum` on PATH outside npx's cache counts as installed, so `npx lyceum init` with a global
  install does not warn]

## Steps

## Progress

- [x] repro: `TestScripts/ac197_hook_fails_closed_without_lyceum` committed alone; `Cleared repro`.
- [x] diagnose: root cause and AC197 written.
- [x] fix: hooks, shims, pre-commit, init warning, docs; this repo's installed hooks refreshed with the built binary.
- [ ] verify

## Log
