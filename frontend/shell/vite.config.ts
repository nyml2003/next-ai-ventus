import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import { resolve } from 'path';

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': resolve(__dirname, '../'),
    },
  },
  build: {
    rollupOptions: {
      input: {
        home: resolve(__dirname, '../pages/home/index.html'),
        post: resolve(__dirname, '../pages/post/index.html'),
        login: resolve(__dirname, '../pages/login/index.html'),
        adminPosts: resolve(__dirname, '../pages/admin-posts/index.html'),
        adminEditor: resolve(__dirname, '../pages/admin-editor/index.html'),
      },
    },
  },
});
