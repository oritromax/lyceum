#!/usr/bin/env python3
"""Static Open Graph card for Lyceum: 1200x630, light. Writes og.svg next to this file.
Needs Python with fontTools and DejaVu Sans. Then: node render.mjs (writes ../../public/og.png)."""
import re, os
from fontTools.ttLib import TTFont
from fontTools.pens.svgPathPen import SVGPathPen
from fontTools.pens.transformPen import TransformPen

OUT = os.path.dirname(os.path.abspath(__file__))
BG, INK, MUTED, ACC = "#f8fafc", "#1e293b", "#475569", "#d97706"
lock = open(os.path.join(OUT, "../../../docs/assets/lockup.svg")).read()
body = re.search(r"<svg[^>]*>(.*)</svg>", lock, re.S).group(1)
vbx, vby, vbw, vbh = map(float, re.search(r'viewBox="([^"]+)"', lock).group(1).split())

W, H = 1200, 630
S = 4.3                                   # lockup scale
lw, lh = vbw * S, vbh * S
lx, ly = 100, 130
lockup = f'<g transform="translate({lx} {ly}) scale({S}) translate({-vbx} {-vby})">{body}</g>'

font = TTFont("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf")
gs, cmap, upm, adv = font.getGlyphSet(), font.getBestCmap(), font["head"].unitsPerEm, font["hmtx"]

def text(s, size, x, y, fill):
    k, parts = size / upm, []
    for ch in s:
        g = cmap[ord(ch)]
        pen = SVGPathPen(gs)
        gs[g].draw(TransformPen(pen, (k, 0, 0, -k, x, y)))
        parts.append(pen.getCommands())
        x += adv[g][0] * k
    return f'<path d="{" ".join(parts)}" fill="{fill}"/>', x

t1, _ = text("The Anti-vibe code toolkit", 56, 100, 528, INK)
bar = f'<rect x="100" y="{ly + lh + 42}" width="64" height="6" rx="3" fill="{ACC}"/>'
svg = (f'<svg xmlns="http://www.w3.org/2000/svg" width="{W}" height="{H}" viewBox="0 0 {W} {H}">'
       f'<rect width="{W}" height="{H}" fill="{BG}"/>{lockup}{bar}{t1}</svg>\n')
open(f"{OUT}/og.svg", "w").write(svg)
print("lockup box", round(lw), "x", round(lh), "bar y", round(ly + lh + 42))
