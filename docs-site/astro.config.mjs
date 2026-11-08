// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

const SITE = 'https://lyceum-docs.pages.dev';

export default defineConfig({
  site: SITE,
  output: 'static',
  integrations: [
    starlight({
      title: 'Lyceum',
      description: 'Spec-driven development with a coding agent, gated by approvals and recorded test runs.',
      favicon: '/favicon.svg',
      logo: {
        light: './src/assets/mark.svg',
        dark: './src/assets/mark-dark.svg',
        alt: 'Lyceum',
      },
      head: [
        { tag: 'link', attrs: { rel: 'icon', href: '/favicon-32.png', sizes: '32x32', type: 'image/png' } },
        { tag: 'link', attrs: { rel: 'apple-touch-icon', href: '/apple-touch-icon.png' } },
        { tag: 'meta', attrs: { property: 'og:image', content: `${SITE}/og.png` } },
        { tag: 'meta', attrs: { property: 'og:image:width', content: '1200' } },
        { tag: 'meta', attrs: { property: 'og:image:height', content: '630' } },
        { tag: 'meta', attrs: { property: 'og:image:alt', content: 'Lyceum: The Anti-vibe code toolkit' } },
        { tag: 'meta', attrs: { name: 'twitter:card', content: 'summary_large_image' } },
        { tag: 'meta', attrs: { name: 'twitter:image', content: `${SITE}/og.png` } },
      ],
      social: [{ icon: 'github', label: 'GitHub', href: 'https://github.com/oritromax/lyceum' }],
      editLink: {
        baseUrl: 'https://github.com/oritromax/lyceum/edit/main/docs-site/',
      },
      sidebar: [
        {
          label: 'Start',
          items: ['start/what-is-lyceum', 'start/install', 'start/quickstart'],
        },
        {
          label: 'Concepts',
          items: ['concepts/units-and-gates', 'concepts/approvals', 'concepts/acceptance-criteria'],
        },
        {
          label: 'Guides',
          items: ['guides/agents', 'guides/ci', 'guides/existing-projects'],
        },
        {
          label: 'Reference',
          items: ['reference/commands', 'reference/cli', 'reference/config'],
        },
      ],
    }),
  ],
});
