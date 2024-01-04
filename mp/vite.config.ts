import { defineConfig } from 'vite';
import uni from '@dcloudio/vite-plugin-uni';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [uni()],
  resolve: {
    alias: {
      'decode-named-character-reference': 'node_modules/decode-named-character-reference/index.js',
    },
  },
});
