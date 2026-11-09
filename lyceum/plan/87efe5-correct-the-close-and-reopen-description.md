---
id: 87efe5-correct-the-close-and-reopen-description
type: doc
title: correct the close and reopen descriptions and log decisions D3 to D5
branch: docs/87efe5-correct-the-close-and-reopen-description
status: in-progress
gate: draft
approvals: []
---

## Brief

> lets do the doc one, then we will tag the new release

Item list the owner approved, as given to him:

> - The README and quickstart claim `close` "folds the spec change into the area files". That's wrong: `close` copies nothing.
> - The reopen rule is described the old way in `architecture.md` (the `gate-reopen` line) and `units-and-gates.md:60`. The build already fixed `test-runner.md:112`.
> - `decisions.md` has no rows for three of your choices: B (`close` pushes), A (reopen to `cause`) and the AC204 shape.

## What and why

Docs and spec prose still describe `close` and the bug reopen rules as they were before units 6426ac, 1ed791
and 94d01a. `close` copies no spec text (it was written and approved at the `spec` gate); it checks, closes,
commits, pushes the unit branch (AC200) and opens the PR. A bug reopens from `cause` to `fix`, and at `close`
to `cause` (AC201) or `fix` (AC204). The three owner choices behind those units have no decision rows.
Correct the prose before the release is tagged. Non-code files only.

## Spec delta

Prose only, no criterion changes: `commands.md` `close` row purpose; `architecture.md` `gate-reopen` paragraph.

## Steps

1. README.md and quickstart.md: rewrite the `close` sentence (no folding; checks, closes, pushes, opens PR).
2. lyceum/spec/commands.md: fix the `lyceum close` purpose. internal/cli/root.go carries the same string but is
   code; left as a follow-up (reference/cli.md is generated from `--help`).
3. lyceum/spec/architecture.md: `gate-reopen` paragraph lists every reopen (cause→fix, close→cause AC201,
   close→fix AC204, chore upgrade).
4. docs-site concepts/units-and-gates.md: bug section covers the close-gate reopens.
5. lyceum/decisions.md: append D3 (close pushes), D4 (reopen to cause), D5 (AC204 shape).
6. Grep README, docs-site docs and spec for other false fold/push claims about `close`: none beyond the above.

## Progress

## Log
