import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import { resolve } from 'path';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@navi/web-ui': resolve(__dirname, './node_modules/@navi/web-ui'),
      '@src': resolve(__dirname, './src'),
      '@utils': resolve(__dirname, './src/utils'),
      '@components': resolve(__dirname, './src/components'),
      '@redux': resolve(__dirname, './src/redux'),
      '@icons': resolve(__dirname, './src/icons'),
      '@store': resolve(__dirname, './src/store'),
    },
  },
});
