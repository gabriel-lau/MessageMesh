import type { Config } from 'tailwindcss';
import flowbitePlugin from 'flowbite/plugin';

export default {
  content: ['./src/**/*.{html,js,svelte,ts}', './node_modules/flowbite-svelte/**/*.{html,js,svelte,ts}'],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        // flowbite-svelte
        primary: {
          50: '#f3faf5',
          100: '#e3f5ea',
          200: '#c7ebd6',
          300: '#9cd9b5',
          400: '#68c08d',
          500: '#44a36c',
          600: '#338657',
          700: '#2b6a46',
          800: '#26553b',
          900: '#214632',
          950: '#0d2618'
        }
      }
    }
  },
  plugins: [require('@tailwindcss/typography'), flowbitePlugin]
} as Config;
