import { fileURLToPath, URL } from 'node:url'

import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vite'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  server: {
    // 如果使用docker-compose开发模式，设置为false
    open: false,
    port:5173,
    proxy: {
      // 把key的路径代理到target位置
      // detail: https://cli.vuejs.org/config/#devserver-proxy
      ["/api"]: { // 需要代理的路径   例如 '/api'
        target: "http://127.0.0.1:23000", // 代理到 目标路径
        changeOrigin: true,
        rewrite: path => path.replace(new RegExp('^' + "/api"), ''),
      }
    },
  },
})
