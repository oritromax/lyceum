---
id: 0e888d-ci-install-golangci-lint-with-go-install
type: chore
title: 'CI: install golangci-lint with go install instead of the release CDN'
branch: chore/0e888d-ci-install-golangci-lint-with-go-install
status: in-progress
gate: build
files: [.github/workflows/ci.yml]
approvals:
  - gate: plan
    hashes:
      lyceum/plan/0e888d-ci-install-golangci-lint-with-go-install.md: sha256:3b010b92561a28b3fbd8fefb1c53485efa3303c124ea6344ac0d10827fb2a513
    commit: 13d1d22
    by: oritromax@gmail.com
    at: 2026-11-08T12:15:31Z
---

## What and why

The `Lint` step in `ci` keeps failing on `main` and on PRs because `golangci-lint-action` downloads the
v2.15.0 binary from `release-assets.githubusercontent.com` and gets HTTP 502/504 (four failed runs in a row on
both OSes, no lint finding in any of them). The same step passed on PR #26, so it is the CDN path and not our
config. Switch the action to `install-mode: goinstall`, which builds the pinned version with the Go toolchain
through the module proxy and skips that CDN.

Checked before planning: `goinstall` is a documented input of `golangci-lint-action@v9`, and
`go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.15.0` builds and runs here with Go 1.27.3,
the version `go.mod` pins and `setup-go` installs in the job.

## Spec delta

None. CI only.

## Steps

1. `.github/workflows/ci.yml`: add `install-mode: goinstall` to the Lint step's `with:`, keeping
   `version: v2.15.0`. Add a one-line comment saying why, so nobody reverts it to the default.
2. Verify: `actionlint`; `lyceum test`; `lyceum check --ci`. The real proof is the PR run: both `Lint` jobs
   (ubuntu and macOS) must pass. Build time per job goes up; I will report the measured `Lint` duration from
   that run next to the previous passing run.
3. Stop at `close`.

## Out of scope

Retry wrappers, a third-party action, changing the lint config or version, and the PR #27 bot ("LBCron").

## Fallback

If the Go proxy path also proves flaky, `install-mode: none` with a vendored or cached binary is the next
step, as a separate chore.

## Progress

1. Done: `install-mode: goinstall` and the explaining comment on the Lint step, `version: v2.15.0` kept.
2. Done: `actionlint` clean. Not yet proven: the PR run, where both Lint jobs must pass.

## Log
