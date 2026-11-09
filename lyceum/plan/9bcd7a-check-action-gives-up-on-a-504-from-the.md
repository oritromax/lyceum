---
id: 9bcd7a-check-action-gives-up-on-a-504-from-the
type: bug
title: check action gives up on a 504 from the release download
branch: fix/9bcd7a-check-action-gives-up-on-a-504-from-the
status: in-progress
gate: diagnose
approvals: []
files: [check-action/run.sh, check-action/action.yml, cmd/lyceum/script_test.go, cmd/lyceum/testdata/script/ac199_check_action_download_fallback.txtar, docs/ci.md, docs-site/src/content/docs/guides/ci.md]
acs: [AC199]
---

## Brief

Owner, 2026-11-09, filing the parked follow-ups one unit at a time:

> Lets open them up as unit, one after another, so they don't cause cascading issues. Open one, fix it, PR close
> and then move to the next until all of these are cleared.

The follow-up as listed to him:

> The check action needs `--retry-all-errors`, a longer backoff, and a fallback to the API asset endpoint for the
> binary download. I'd do this one first, because it's the one users actually see.

Observed the same day: jai run 42581536020 and Lyceum PRs #32, #33, #34 failed `check` with
`curl: (22) The requested URL returned error: 504` four times, then
`lyceum check-action: cannot download https://github.com/oritromax/lyceum/releases/download/v1.0.0/lyceum_1.0.0_linux_amd64.tar.gz`.
From outside CI the public download URL (archive and `checksums.txt`) returned 504 every time while
`api.github.com/repos/oritromax/lyceum/releases/assets/<id>` with `Accept: application/octet-stream` returned 200.
The asset was `uploaded`; githubstatus.com said all systems operational.

## What and why

- **Symptom:** `check-action/run.sh` downloads with `curl -fsSL --retry 3`. curl does count a 504 as transient,
  and the log shows four attempts, but the default backoff (1 s, 2 s, 4 s) gives up about 7 seconds in, then
  the job fails before `lyceum check --ci` runs. A CDN problem on GitHub's release-download path that lasts minutes
  turns every consumer's CI red with no Lyceum finding behind it. [DERIVED: the owner's follow-up named
  `--retry-all-errors`; curl's man page lists 504 as transient, so the gap is the window, not the error class]
- **Expected:** a transient failure on the download path is retried with a longer backoff, and if the public
  release URL still fails, the action fetches the same asset through the GitHub REST API (with the job's token when
  one is given), then verifies the checksum exactly as now. Only when both paths fail does the job fail, naming both.
- **Trigger:** any 5xx on `github.com/<repo>/releases/download/...` lasting longer than the retry window.
- **Out of scope:** changing which release is picked, caching the binary across runs, and the checksum rule.

## Root cause

- `check-action/run.sh` fetched each file with one call, `curl -fsSL --retry 3 -o … "$base_url/$tag/$f"`, then
  `fail`ed. Two gaps:
  - The retry window is curl's default backoff for `--retry 3` (1 s, 2 s, 4 s): four tries over about 7 seconds.
    An outage on GitHub's release-download CDN that lasts minutes outlives it.
  - There is one source. The same asset is also served by the REST API
    (`api.github.com/repos/<repo>/releases/assets/<id>`, `Accept: application/octet-stream`), which stayed up in the
    incident, but the script never tries it.
- The repro `TestScripts/ac199_check_action_download_fallback` serves the release from a fake GitHub
  (`fakegithub`, an in-process HTTP server) whose public download path answers 504 every time while its REST API
  works. On the old code `run.sh` printed four `curl: (22) … 504` lines and
  `cannot download …/down/v1.0.0/lyceum_1.0.0_linux_amd64.tar.gz` after 7 s, so the first `exec` failed.

## Fix

- `run.sh`: each download retries 5 times, 10 s apart, within 120 s (`--retry-delay`, `--retry-max-time`), and
  with `--retry-all-errors` when `curl --help all` lists it (older curl on a runner just goes without).
- If the public URL still fails, `api_download` fetches the release by tag from the REST API once, finds the asset
  id by name with `tr`/`awk` (the script has never used `jq`), and downloads it with
  `Accept: application/octet-stream`, under the same retries. `Authorization: Bearer …` goes only to the API, and
  only when a token is set.
- Both paths failing: `cannot download <file> from <public URL> or <release API URL>`.
- Checksum verification and everything after it is unchanged.
- `action.yml`: optional `token` input, default `github.token`, passed as `LYCEUM_ACTION_TOKEN`.
- Test overrides, like the existing `LYCEUM_ACTION_BASE_URL`: `LYCEUM_ACTION_API_URL`, `LYCEUM_ACTION_RETRIES`,
  `LYCEUM_ACTION_RETRY_DELAY`.
- Docs: `docs/ci.md` and docs-site `guides/ci.md` step 2 say how the download retries and falls back, and list the
  `token` input.

## Spec delta

- `lyceum/spec/enforcement.md`: AC199 (new). Registry range in `lyceum/SPEC.md` gains AC199; `lyceum/LIBRARY.md`
  names it.
- AC199 Given the public release download answering 5xx, when the Action downloads the archive or `checksums.txt`,
  then it retries over a longer window, then fetches the same asset through the GitHub REST API, sending the token
  when one is set; it verifies the checksum as before and fails only when both paths fail, naming both URLs.
  [DERIVED: token, checksum and both-paths failure from the plan's What and why]
- [AGENT-PROPOSED: 5 retries 10 s apart within 120 s per URL, plus `--retry-all-errors` where curl has it]
- [AGENT-PROPOSED: optional `token` input defaulting to `github.token`; the token is never sent to the public URL]

## Steps

## Progress

- [x] repro: `TestScripts/ac199_check_action_download_fallback` committed alone; `Cleared repro`.
- [x] diagnose: root cause and AC199 written.
- [x] fix: `run.sh` retries and REST API fallback, `token` input, docs.
- [ ] verify

## Log
