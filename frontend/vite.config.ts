import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    proxy: {
      '/users': 'http://localhost:9000',
      '/posts': 'http://localhost:9000',
      '/storage': 'http://localhost:9000',
    },
  },
})
