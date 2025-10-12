import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import * as path from 'node:path'

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src')
    }
  },
  server: {
    host: true,
    port: 5173,
    proxy: {
      // Proxy API requests to the Go backend during development
      '/health': process.env.VITE_API_PROXY_TARGET || 'http://localhost:8080',
      '/api': process.env.VITE_API_PROXY_TARGET || 'http://localhost:8080'
    }
  }
})
