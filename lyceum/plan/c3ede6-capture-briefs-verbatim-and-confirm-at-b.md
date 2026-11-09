---
id: c3ede6-capture-briefs-verbatim-and-confirm-at-b
type: feat
title: capture briefs verbatim and confirm at both ends
branch: feat/c3ede6-capture-briefs-verbatim-and-confirm-at-b
status: in-progress
gate: build
approvals:
  - gate: spec
    hashes:
      lyceum/SPEC.md: sha256:1fb14b52d122f4b30123a1594291c728deef67725a5c349cae041e93f1e1afdc
      lyceum/plan/c3ede6-capture-briefs-verbatim-and-confirm-at-b.md: sha256:bef3cdb491255ab7491e7f19945c5bb1358bc9e4a5fd33abb98f44f0b246de9b
      lyceum/spec/units.md: sha256:1938d20627b65ce84fd779e37379ffcdb3baf00b83928c05e10953f50ddd6913
    commit: 372108f
    by: oritromax@gmail.com
    at: 2026-11-09T05:28:55Z
areas: [units]
files:
  - internal/templates/skills/init/SKILL.md
  - internal/templates/skills/feat/SKILL.md
  - internal/templates/skills/bug/SKILL.md
  - internal/templates/skills/doc/SKILL.md
  - internal/templates/templates_test.go
  - cmd/lyceum/testdata/script/ac191_project_hashes_brief.txtar
  - .claude/skills/lyceum-init/SKILL.md
  - .claude/skills/lyceum-feat/SKILL.md
  - .claude/skills/lyceum-bug/SKILL.md
  - .claude/skills/lyceum-doc/SKILL.md
  - .agents/skills/lyceum-init/SKILL.md
  - .agents/skills/lyceum-feat/SKILL.md
  - .agents/skills/lyceum-bug/SKILL.md
  - .agents/skills/lyceum-doc/SKILL.md
  - README.md
  - docs-site/src/content/docs/start/quickstart.md
  - docs-site/src/content/docs/guides/existing-projects.md
acs: [AC185, AC186, AC187, AC188, AC189, AC190, AC191, AC192, AC193, AC194, AC195]
---

## What and why

- Interviews run from chat only. A brief sent over several messages is never kept verbatim, and nothing asks
  "is that everything?" at either end. Approved intent is lost and `[DERIVED]` / `[AGENT-PROPOSED]` markers have
  nothing to be checked against.
- Add an opening and a closing checkpoint to `lyceum-init`, `lyceum-feat`, `lyceum-bug` and `lyceum-doc`; keep the
  brief verbatim (`lyceum/spec/brief.md` at greenfield, the plan's `## Brief` in a unit); freeze it under the gate
  that approves what it led to. Spec: `lyceum/spec/units.md`, "Interviews" → "Checkpoints", and "Greenfield".
- Hash decision: register `brief.md` as area `brief` instead of hashing it explicitly. No Go change:
  - `approve project` already hashes `SPEC.md` plus every registered area file (`approval.FilesFor`).
  - `check`'s registry check already fails any `lyceum/spec/*.md` with no registry row, so an unregistered
    `brief.md` would need a Go exemption anyway. Registering is the zero-code path; explicit hashing needs changes
    in `approval` and an exemption in `check`.
  - The ACs cell of the `brief` row stays empty; the registry check only validates ACs a file defines.
- Unit briefs go in the plan, not in `brief.md`. Also no Go change: the plan hash (`PlanHashBasis`) covers every
  section but Progress and Log, so `approve spec` (feat) and `approve cause` (bug) freeze `## Brief`. Writing to
  `brief.md` from a unit would break the project approval and, for a doc unit, need `approve project` again.
- Phase 1 is skill prose plus tests. AC191 pins the existing CLI behaviour for `brief.md` with a testscript.

## Spec delta

- `units`: "Interviews" gains "Checkpoints" (opening, closing, ready is not approval, late additions, provenance,
  freezing). "Greenfield" names `brief.md` as registered area `brief`.
- New ACs: AC185–AC195. Registry row `units` extended.

## Steps

1. `init` skill (AC185, AC186, AC187, AC188, AC189, AC190, AC192): opening checkpoint after the open first
   question; closing loop before step 4's STOP; ready is not approval; late additions; register the `brief` row;
   provenance against `brief.md`. Keep it under the 120-line budget.
2. `feat` and `bug` skills (AC188, AC189, AC192, AC193, AC194): opening checkpoint into `## Brief` before the
   interview; closing loop before `approve spec` / `approve cause`. Both are near the 120-line budget: tight wording.
3. `doc` skill (AC189, AC193, AC195): opening checkpoint; closing loop before `lyceum test`; note there is no
   approval. Under the 60-line budget.
4. Tests in `internal/templates/templates_test.go` (AC185–AC190, AC192–AC195): one `TestSkillCheckpoints_ACnnn`
   per AC, asserting the wording with `oneLine`, like `TestSkillRules_AC164`.
5. Testscript `ac191_project_hashes_brief.txtar` (AC191): greenfield init, register `brief` plus one area,
   `approve project` hashes `lyceum/spec/brief.md`, edit and commit it, `check` names it.
6. Refresh the installed copies: build, run `dist/lyceum init` in this repo so `.claude/skills` and `.agents/skills`
   match the binary (AC120, `doctor`); confirm no other file changes.
7. Prose mirrors (hand-written, not generated): the interview lines in `README.md`, `quickstart.md` and
   `existing-projects.md` name the two checkpoints in one clause.
8. `go test ./...`, `golangci-lint run`, `lyceum test`.

## Progress

- [x] Step 1: `init` skill checkpoints, `brief` row, provenance (71 lines).

## Log

- Owner decisions in the brief: checkpoints 1–7 as listed in "Checkpoints"; hash approach left to the agent
  (chosen: register `brief.md`).
