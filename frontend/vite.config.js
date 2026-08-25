import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    host: '0.0.0.0',  // 🔥 Слушаем ВСЕ интерфейсы (IPv4 + IPv6)
    port: 5173,
    strictPort: true,
  },
})