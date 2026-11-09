---
id: 9bcd7a-check-action-gives-up-on-a-504-from-the
type: bug
title: check action gives up on a 504 from the release download
branch: fix/9bcd7a-check-action-gives-up-on-a-504-from-the
status: planned
gate: repro
approvals: []
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

## Spec delta

## Steps

## Progress

## Log
