import dotenv from 'dotenv'
import { resolve } from 'path'

// Load .env from the root of the project
dotenv.config({ path: resolve(process.cwd(), '../.env') })

const backendUrl = process.env.PUBLIC_URL || 'http://localhost:8080'
const frontendUrl = process.env.FRONTEND_URL || 'http://localhost:3000'

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  runtimeConfig: {
    public: {
      backendUrl,
      frontendUrl
    }
  },
  modules: [
    '@nuxt/eslint',
    '@nuxt/ui',
    '@pinia/nuxt'
  ],

  devtools: {
    enabled: true
  },

  css: ['~/assets/css/main.css'],

  routeRules: {
    '/': { prerender: true },
    '/api/**': { proxy: `${backendUrl}/api/**` }
  },

  compatibilityDate: '2025-01-15',

  eslint: {
    config: {
      stylistic: {
        commaDangle: 'never',
        braceStyle: '1tbs'
      }
    }
  }
})
