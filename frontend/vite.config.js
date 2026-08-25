import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    host: '0.0.0.0',
    port: 5173,
    strictPort: true,
    
    // 🔥 РАЗРЕШАЕМ все хосты для dev-режима (включая ngrok, loca.lt, и т.д.)
    allowedHosts: true,
    
    // 🔥 Настраиваем HMR для работы через туннель
    hmr: {
      // Vite сам определит правильный протокол и хост
      protocol: 'wss',
      clientPort: 443,
    },
  },
})