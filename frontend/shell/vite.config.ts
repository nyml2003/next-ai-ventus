import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import { resolve } from 'path';

export default defineConfig({
  plugins: [react()],
  root: '../',
  resolve: {
    alias: {
      '@ventus/types': resolve(__dirname, '../packages/types/src'),
      '@ventus/utils': resolve(__dirname, '../packages/utils/src'),
      '@ventus/ui': resolve(__dirname, '../packages/ui/src'),
      '@ventus/markdown': resolve(__dirname, '../packages/markdown/src'),
      '@ventus/request': resolve(__dirname, '../packages/request/src'),
      '@ventus/store': resolve(__dirname, '../packages/store/src'),
      '@ventus/orchestration': resolve(__dirname, '../packages/orchestration/src'),
    },
    dedupe: ['react', 'react-dom'],
  },
  optimizeDeps: {
    include: ['react', 'react-dom', 'react/jsx-runtime', 'react/jsx-dev-runtime'],
  },
  build: {
    outDir: './dist',
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
  server: {
    port: 3000,
    proxy: {
      '/api': 'http://localhost:8080',
      '/uploads': 'http://localhost:8080',
    },
    // 重写规则：/pages/post/{slug} -> /pages/post/index.html
    fs: {
      strict: false,
    },
  },
});
