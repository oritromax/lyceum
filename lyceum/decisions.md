# Lyceum — decision log

Append-only. Rows are never edited or renumbered; a later decision that changes an earlier one says so
(`Revises Dn`, `Resolves Dn`). The spec (`lyceum/SPEC.md` and `lyceum/spec/`) is the current design;
this file is why it says what it says. The log starts at v1: the decisions before it are folded into the spec.

## Decision register

| # | Decision | Source |
|---|---|---|
| D1 | Lyceum v1: the first public version. Earlier design decisions are folded into the spec | owner, 2026-11-06 |
| D2 | `vN` follows the newest verified prerelease until a stable `vN.x.y` exists; from then on only stable releases move it, and `check-action` accepts a prerelease behind `vN` only when no stable release is there. Revises the release rule in the README and the release workflow that `v1` never moves to a prerelease | owner, 2026-11-06 |
| D3 | `lyceum close` pushes the unit branch itself (with `-u` when it has no upstream) before opening the PR; if the push fails the unit stays closed and `close` prints the commands. Revises "Lyceum never pushes" for `close` only (unit 6426ac, AC200). Chosen over refusing to open the PR: "B, cause that was my initial idea" | owner, 2026-11-09 |
| D4 | A bug at `close` whose cause-hashed file changes (e.g. a merge of main) goes back to `cause` with a `gate-reopen` event, rather than failing approval-integrity with no way forward (unit 1ed791, AC201). Option A, chosen over allowing `approve cause` at `close` or ignoring drift from main | owner, 2026-11-09 |
| D5 | At `close`, non-test code in the branch's own diff after the `verify` clear sends the bug back to `fix`; code that arrived only by merging the base branch does not (unit 94d01a, AC204) | owner, 2026-11-09 |
