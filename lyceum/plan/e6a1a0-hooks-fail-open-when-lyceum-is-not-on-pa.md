---
id: e6a1a0-hooks-fail-open-when-lyceum-is-not-on-pa
type: bug
title: hooks fail open when lyceum is not on PATH after npx init
branch: fix/e6a1a0-hooks-fail-open-when-lyceum-is-not-on-pa
status: closed
gate: close
approvals:
  - gate: cause
    hashes:
      .claude/settings.json: sha256:17e8897d469bec327d839c27ce5169e2ac6db8d5dceac3e0eec5a3f8a313073c
      .codex/hooks.json: sha256:af1c7dd970fb062b57c149a19ad2be437853cedcea52712d3a52745603b76015
      README.md: sha256:9cf4b25a1e013a4d3572ce63820ff198f555ed5aacb2826d2f96ff7a4214c4e3
      docs-site/src/content/docs/index.mdx: sha256:e343178f666aea5dd21b3a4f8d74721d1320d5f0326345681c19fae3848fff81
      docs-site/src/content/docs/start/install.md: sha256:18c76bf8e6b7825e61625d5679ebd92a841d42bdde5910e93b6254bd9add6312
      internal/cli/init.go: sha256:8a217aec028f5a2ecc91aa87bb03f3eee9cd976d625c1def83b25c50fd5f2c4e
      internal/cli/initharness.go: sha256:5b3864403002daf10a7e5e18ead3da272eb01dd4f2742fd2a12381f75a04f992
      internal/templates/opencode-plugin.js: sha256:987acf7d9ebff3b2f013e39f1491ef0b7cedbf6bd266a13bbd1190fe9df916db
      internal/templates/pi-extension.ts: sha256:db1a91a65779c90cab3b5ef9b7a73e0d7c37bb7c333d57e4232f44735fb0f9d2
      internal/templates/pre-commit: sha256:0e1020eb2291d8f7c174996dfd8d9fd2c394554f8e641de49a5dda9d27f1c9ff
      internal/templates/templates.go: sha256:97e47659bf0ccdc2358898fd57e64460961a8bfb29ff44fe613b2f10a3cf96f4
      lyceum/LIBRARY.md: sha256:3ebb50da074fded97217b34862ccef0b93b88e0be05a1b79f1ef841317efdccb
      lyceum/SPEC.md: sha256:df20bf6399bcabfd06a6c6fa6e2ce2041d49ffa5a7966025efb87775159f8daa
      lyceum/plan/e6a1a0-hooks-fail-open-when-lyceum-is-not-on-pa.md: sha256:f50763d3d1c8b41210cbc4ffa59a5328277f3d49e2614eb704cd4a2593089f2e
      lyceum/spec/enforcement.md: sha256:2a2a96ccf1df5b479f2e275d35cdf1e82ffe129457857e53bcaedcba41d3e738
    commit: 0d96974
    by: oritromax@gmail.com
    at: 2026-11-09T07:28:32Z
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
