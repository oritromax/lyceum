---
id: 311a4a-close-help-text-says-it-folds-the-spec-d
type: chore
title: close help text says it folds the spec delta
branch: chore/311a4a-close-help-text-says-it-folds-the-spec-d
status: in-progress
gate: build
approvals: []
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
