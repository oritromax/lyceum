---
id: 0b1405-an-existing-testing-md-without-frontmatt
type: bug
title: an existing TESTING.md without frontmatter leaves brownfield init at a dead end
branch: fix/0b1405-an-existing-testing-md-without-frontmatt
status: implemented
gate: close
approvals:
  - gate: cause
    hashes:
      .agents/skills/lyceum-init/SKILL.md: sha256:b882d08dfed154e3cd23daca5bbf274ba75de0a5298b8c590b500281ded14fac
      .claude/skills/lyceum-init/SKILL.md: sha256:b882d08dfed154e3cd23daca5bbf274ba75de0a5298b8c590b500281ded14fac
      README.md: sha256:0d4386eb052d2da32025e6a63215526275f81b5c99d97a2f8665d8160b44fe66
      docs-site/src/content/docs/guides/existing-projects.md: sha256:57acc6ff6e13d58fe0608df28d6a6e896660cb831bd28f715521173c8992241a
      internal/cli/approve.go: sha256:08049d7517d688086b8bdce05dafda9737f6048391092004e6c0420c56526839
      internal/cli/init.go: sha256:350746467713d191f44b8d36de0d90fd3bf72d3d99b658daff447e02381e33ff
      internal/cli/status.go: sha256:da678a55883d72d42ce19b3bb48b7fae88104ddadad88944a3a5065ee37e41e7
      internal/templates/skills/init/SKILL.md: sha256:b882d08dfed154e3cd23daca5bbf274ba75de0a5298b8c590b500281ded14fac
      internal/testspec/testspec.go: sha256:cc5f1c0a461e73fde06ed88be532762a1fa311b9a1e7bb8d16f0ef670ad8df61
      lyceum/LIBRARY.md: sha256:d8858322d1006b93a1cacf5975cf440b7af29b27aca094f2eb13001c86ed34b9
      lyceum/SPEC.md: sha256:e683e54d4c1f262612c6507784aadad885f60a2f1d6b08ee9e330b9cc6fe8f8d
      lyceum/plan/0b1405-an-existing-testing-md-without-frontmatt.md: sha256:6f3eee8663a775473d28f6f260a86b79e3f94eba6aae8bb1441a4dfaedfbe271
      lyceum/spec/distribution.md: sha256:557893345742d7705e6bbbe4e87d9db858bf7898ab86b372839591624210a7ff
    commit: 6ae2945
    by: oritromax@gmail.com
    at: 2026-11-09T10:14:28Z
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
