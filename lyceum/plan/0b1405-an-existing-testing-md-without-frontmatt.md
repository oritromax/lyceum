---
id: 0b1405-an-existing-testing-md-without-frontmatt
type: bug
title: an existing TESTING.md without frontmatter leaves brownfield init at a dead end
branch: fix/0b1405-an-existing-testing-md-without-frontmatt
status: in-progress
gate: fix
approvals: []
acs: [AC203]
files:
  - internal/testspec/testspec.go
  - internal/cli/init.go
  - internal/cli/approve.go
  - internal/cli/status.go
  - internal/templates/skills/init/SKILL.md
  - .agents/skills/lyceum-init/SKILL.md
  - .claude/skills/lyceum-init/SKILL.md
  - README.md
  - docs-site/src/content/docs/guides/existing-projects.md
---

## Brief

Owner, 2026-11-09, reporting it from the jai repo:

> Another bug to handle, ✓  ≈~/work/office/jai  (main *) ⦔ lyceum approve testing
>   Can't approve testing: TESTING.md: no frontmatter: the file must start with a "---" line
>   Next: fix TESTING.md, then run "lyceum approve testing" yourself in a terminal
>
> where testing.md already existed. We have to let the model inspect ( not edit ) files that lyceum needs and
> suggest if they needs editing before lyceum can work properly. Don't open a unit for this yet, but we will get
> back to it.

Then, filing the parked follow-ups one at a time:

> Lets open them up as unit, one after another, so they don't cause cascading issues. Open one, fix it, PR close
> and then move to the next until all of these are cleared.

The proposed shape (CLI reports and never edits; `init`, `status`, `doctor` name the gap; the error says the file
exists but is not in Lyceum's format; the `lyceum-init` skill reads the kept file and proposes the exact
frontmatter for the owner to apply; no auto-fix), filed as a bug, was approved with:

> go ahead

## What and why

- **Symptom (probe on main at `8cc9ad8`):** a repo with its own `TESTING.md` (prose, no frontmatter) runs
  `lyceum init`. `init` keeps the file and prints nothing about it (every created file gets a `Created` line).
  `check` says `No findings`. `status` lists `testing for TESTING.md: run "lyceum approve testing"`, which refuses
  with `no frontmatter: the file must start with a "---" line`. Only `doctor` names it, as `fix TESTING.md`, without
  saying what to add. The owner is sent to an approval that cannot succeed.
- **Expected:**
  1. `init`, when it keeps an existing file Lyceum parses and that file is not in Lyceum's format, prints a line
     naming the file and what is missing (for `TESTING.md`: the frontmatter keys, at the top, content kept).
     [DERIVED: owner's "inspect … and suggest"]
  2. `status` (and `doctor`) point to that fix instead of `approve testing` while the file cannot be approved.
  3. The refusal and the `doctor` message say the file exists but is not in Lyceum's format: add frontmatter at
     the top and keep the content; they name the required keys.
  4. The `lyceum-init` skill (template and installed copies) has a step: for each kept file the CLI reports,
     read it, derive values from the file and the repo (e.g. the suite command its prose names), and propose the
     exact frontmatter block to the owner. The agent does not edit the file. [DERIVED: owner's "inspect ( not
     edit )"]
  5. The CLI never writes or prepends to an existing file.
- **Trigger:** commit a prose `TESTING.md`, run `lyceum init`, then `lyceum status`.
- **Out of scope:** files Lyceum only appends a section to (`AGENTS.md`, `CLAUDE.md`); auto-fixing.

## Root cause

- `init` (`internal/cli/init.go`) never looks at an existing `TESTING.md`: it is not in `initFiles`, so a kept
  file gets no line at all, readable or not.
- `status` (`internal/cli/status.go`) treats any existing, unapproved `TESTING.md` as a pending `testing` gate and
  always gives `gateNext(Judgment, testing)` = `run "lyceum approve testing"`, without checking that the file
  reads. `approve testing` (`internal/cli/approve.go`) does check, via `testspec.Read`, and refuses: the dead end.
- The refusal text comes from `frontmatter()` in `internal/testspec/testspec.go`: `no frontmatter: the file must
  start with a "---" line`. It names neither that the file is the owner's own nor the keys parsing requires
  (`lyceum` = 1 and a non-blank `suite`; `single`, `test_files`, `subsets`, `environments`, `timeout` are
  optional). `doctor` and `approve` print that error as is.
- The repro `TestScripts/ac203_init_kept_testing_not_lyceum_format` runs the probe: it first fails on init
  printing no `Kept TESTING.md` line, and asserts status, doctor and approve after it.
- Survey of other kept files that init writes when missing and later parses: `lyceum/config.yaml` fails every
  command loudly with its parse error (no approval route, so no dead end); `RULES.md` is not parsed; a prose
  `lyceum/SPEC.md` reads as no areas registered. Only `TESTING.md` has the dead end, so only it is covered.
  [DERIVED: survey of initFiles and their readers]

## Spec delta

- `distribution.md`: AC203 added; `lyceum init` step 2 says an unreadable kept TESTING.md is named by `init` and
  the skill proposes its frontmatter. Registry range in `lyceum/SPEC.md` updated.

## Steps

## Progress

- Repro committed alone; `lyceum test --repro TestScripts/ac203_init_kept_testing_not_lyceum_format` cleared
  repro (failed at `ac203_…txtar:11`, no `Kept TESTING.md` line).
- Fix: the no-frontmatter error says the file exists but is not in Lyceum's format and names `lyceum: 1` and
  `suite:` (keeps its `TESTING.md: no frontmatter` prefix); `init` prints `Kept <error>` when it cannot read an
  existing TESTING.md; `status` points to fixing the frontmatter (via `doctor`) instead of `approve testing`
  while it cannot be read; `lyceum-init` skill gains "Files `init` kept" (propose, never edit); README and the
  existing-projects guide get a line each. The CLI writes nothing to the kept file.

## Log
