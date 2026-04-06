// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  app: {
    head: {
      htmlAttrs: { lang: 'ru' }
    }
  },
  devtools: { enabled: true },
  modules: ['@nuxtjs/tailwindcss'],
  css: ['@/assets/css/tailwind.css'],
  runtimeConfig: {
    public: {
      apiBase: '/api/v1'
    }
  },
  routeRules: {
    '/api/**': {
      proxy: 'http://localhost:8080/api/**'
    }
  }
})
