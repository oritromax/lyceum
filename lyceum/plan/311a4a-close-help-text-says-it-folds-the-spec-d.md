---
id: 311a4a-close-help-text-says-it-folds-the-spec-d
type: chore
title: close help text says it folds the spec delta
branch: chore/311a4a-close-help-text-says-it-folds-the-spec-d
status: in-progress
gate: build
approvals:
  - gate: plan
    hashes:
      lyceum/plan/311a4a-close-help-text-says-it-folds-the-spec-d.md: sha256:87f5a19f80b3c06cb6ca7dc6d87186b3d579a0f95e5fa07fead873b10f11e80d
    commit: 5f23f46
    by: oritromax@gmail.com
    at: 2026-11-10T05:44:36Z
files:
  - internal/cli/root.go
---

## Brief

Owner, before tagging the next release:

> lets do the doc one, then we will tag the new release

Left over from doc unit 87efe5, which could not change code: `lyceum close --help` (`internal/cli/root.go:142`)
still says "Fold the spec delta, close the unit, open a PR", and the generated `reference/cli.md` copies it.

## What and why

`close` copies nothing into the spec (units.md "Close"; commands.md row fixed in 87efe5). The help text now
matches `lyceum/spec/commands.md:15`: "Check and close the unit, push its branch, open a PR".

## Spec delta

## Steps

1. Change the `close` short help in `internal/cli/root.go` to the commands.md wording.

## Progress

## Log
