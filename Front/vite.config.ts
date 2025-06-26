import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { configDefaults } from 'vitest/config'

export default defineConfig({
  plugins: [vue()],
  test: {
    environment: 'jsdom',
    globals: true, // optionnel mais utile avec Vitest
    setupFiles: './src/test/setup.ts', // si tu veux faire un fichier de config test
    exclude: [...configDefaults.exclude, 'e2e'],
  },
})
