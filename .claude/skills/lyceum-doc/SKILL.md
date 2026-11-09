---
name: lyceum-doc
description: Use when the owner wants a human-facing document under docs/ in a Lyceum project, written and verified as a doc unit.
---

# lyceum-doc

Run `lyceum status` first. If it names a unit in flight, its printed next action wins over anything here.
Every `lyceum` command ends with a `Next:` line; follow it.

Gates for a doc: `draft` → `verify` (`lyceum test`) → `close`. No owner approval.

1. No unit on this branch for this doc: run `lyceum add doc "<title>"`. With no title, ask for one.
2. Opening checkpoint: before any question, save the owner's brief verbatim to a `## Brief` section of the plan,
   and append each further message to it. A first reply may be the first of several. Ask once "Anything to add,
   or shall I proceed?" and ask numbered questions only after the owner says proceed.
3. `draft`: write the doc under `docs/`. Describe only what exists: never a command, flag or file the project
   doesn't have. If something is unclear, ask (numbered questions, lettered options, recommended first, then
   stop and wait) or mark it `[NEEDS CLARIFICATION: …]`.
4. Before `lyceum test`, run every command the doc shows, exactly as written, and confirm it exits and prints as
   documented. The CLI does not run them; this step is yours. Fix the doc, not the code. Commit.
5. Closing checkpoint: ask "Anything to add, or shall I consider the doc ready?". Append each addition to
   `## Brief`, update the doc, and ask again. A doc has no approval, so "ready" lets you run `lyceum test`.
   Late additions skip the two-round limit: they go straight into the doc and plan, and anything still
   uncertain is left as `[NEEDS CLARIFICATION: …]`.
6. `verify`: with a clean tree, run `lyceum test`. It clears `verify` when the full suite passes and the branch
   changes no code file. In a project with no tests (no file matches `test_files`), it clears `verify` with no
   tests for a docs-only change. Quote the result; never claim it from reasoning.
7. Once `verify` clears, run `lyceum close`. It sets the plan to closed, pushes the branch and opens a PR, or
   prints the commands left to run. A doc has no owner approval, so there is no plan to amend.

After a gate clears, any confirmation to proceed is a structured yes/no, never an open question.

Never change application code or tests in a doc unit, edit CLI-owned plan fields or `lyceum/events.jsonl`, or work
around a hook refusal (stop and tell the owner).
