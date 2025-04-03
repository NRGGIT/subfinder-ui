// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  devtools: { enabled: true },
  modules: ['@nuxt/ui'],
  ui: {
    global: true,
    icons: ['lucide']
  },
  app: {
    head: {
      title: 'Subfinder UI',
      meta: [
        { name: 'description', content: 'UI for Subfinder subdomain enumeration tool' }
      ]
    }
  },
  runtimeConfig: {
    public: {
      apiBaseUrl: process.env.NUXT_PUBLIC_API_BASE_URL || 'http://localhost:8080'
    }
  },
  css: ['~/assets/css/main.css']
})
