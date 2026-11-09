---
id: 0b1405-an-existing-testing-md-without-frontmatt
type: bug
title: an existing TESTING.md without frontmatter leaves brownfield init at a dead end
branch: fix/0b1405-an-existing-testing-md-without-frontmatt
status: in-progress
gate: diagnose
approvals: []
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

## Spec delta

## Steps

## Progress

## Log
