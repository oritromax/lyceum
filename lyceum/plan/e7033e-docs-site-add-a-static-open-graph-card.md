---
id: e7033e-docs-site-add-a-static-open-graph-card
type: chore
title: 'Docs site: add a static Open Graph card'
branch: chore/e7033e-docs-site-add-a-static-open-graph-card
status: implemented
gate: close
files: [docs-site/public/og.png, docs-site/scripts/og/build_og.py, docs-site/scripts/og/render.mjs, docs-site/scripts/og/og.svg, docs-site/astro.config.mjs, docs-site/README.md, .github/workflows/docs-site.yml]
approvals:
  - gate: plan
    hashes:
      lyceum/plan/e7033e-docs-site-add-a-static-open-graph-card.md: sha256:6611e604d81e417d493021913363d7b3fb919669aa8f799350d2169fc88f4522
    commit: 260f5bd
    by: oritromax@gmail.com
    at: 2026-11-08T10:09:57Z
---

## What and why

Closes #19. Shared links to the docs show text only because no page sets `og:image`. Add one static 1200x630
card, light theme: the lockup, an amber bar and the slogan "The Anti-vibe code toolkit". Decided with the
owner: one card for every page (not per-page generated), light variant, slogan as the text.

Starlight already emits `og:title`, `og:description` and `og:url` per page, so only the image tags are added.

## Spec delta

None. No behaviour change in the CLI.

## Steps

1. `docs-site/public/og.png`: the committed card (1200x630, about 40 KB). Served at `/og.png`.
2. `docs-site/scripts/og/`: the source so the card can be regenerated. `build_og.py` reads
   `docs/assets/lockup.svg` and writes `og.svg` with the slogan as outlined glyphs (DejaVu Sans, so no font
   dependency at render time); `render.mjs` rasterises it to `../../public/og.png` with `sharp`, which Astro
   already installs. Regenerating is a manual step that needs Python with `fontTools`; the PNG is committed, so
   the site build needs neither.
3. `astro.config.mjs`: add to the Starlight `head` the tags `og:image` (absolute URL on the site origin),
   `og:image:width`, `og:image:height`, `og:image:alt`, `twitter:card` = `summary_large_image` and
   `twitter:image`.
4. `.github/workflows/docs-site.yml`: the existing output asserts also require `dist/og.png`, and that the
   built home page contains an `og:image` tag.
5. `docs-site/README.md`: how to regenerate the card.
6. Verify: build exits 0; every page in `dist/` has exactly one `og:image` and one `twitter:card`; the image
   URL maps to a file in `dist/`; `file` reports 1200x630; the card regenerates byte-for-byte from the
   committed source; `actionlint`. `lyceum test`, then stop at `close`. Whether a platform renders the preview
   can only be seen live, after the deploy.

## Out of scope (follow-ups)

Per-page cards, a dark variant, using the slogan on the landing hero or in the README, and changing the
site description.

## Progress

1. Done: `docs-site/public/og.png` committed (1200x630).
2. Done: `scripts/og/` holds `build_og.py`, `og.svg` and `render.mjs`; the lockup is read from `docs/assets/lockup.svg`. Regenerating from the repo source reproduces `og.svg` and `og.png` byte for byte.
3. Done: `og:image` (absolute URL), width, height, alt, `twitter:card` and `twitter:image` in the Starlight `head`.
4. Done: workflow asserts `dist/og.png` and an `og:image` tag on the built home page; `actionlint` clean.
5. Done: README section on regenerating the card.
6. Done except the live check: build exits 0 (14 pages); each of the 13 content pages has exactly one `og:image` and one `twitter:card`; `dist/og.png` is 1200x630 and identical to the source. Whether a platform renders the preview can only be seen after the deploy.

## Log
