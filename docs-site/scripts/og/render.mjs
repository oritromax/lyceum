// Rasterises og.svg (next to this file) to ../../public/og.png with sharp, which Astro already installs.
import { fileURLToPath } from 'node:url';
import sharp from 'sharp';

const here = (p) => fileURLToPath(new URL(p, import.meta.url));
await sharp(here('./og.svg'), { density: 96 }).resize(1200, 630).png({ compressionLevel: 9 }).toFile(here('../../public/og.png'));
const m = await sharp(here('../../public/og.png')).metadata();
console.log(`og.png ${m.width}x${m.height} ${m.format}`);
