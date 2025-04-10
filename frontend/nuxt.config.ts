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
      htmlAttributes: {
        'data-theme': 'light', // Default to light mode
        style: '@media (prefers-color-scheme: dark) { html { color-scheme: dark; } }'
      },
      link: [
        {
          rel: 'stylesheet',
          href: 'https://fonts.googleapis.com/css2?family=Inter:wght@400;700&display=swap'
        }
      ],
      meta: [{ name: 'description', content: 'UI for Subfinder subdomain enumeration tool' }],
      title: 'Subfinder UI'
    }
  },
  runtimeConfig: {
    public: {
      apiBaseUrl: process.env.NUXT_PUBLIC_API_BASE_URL || 'http://localhost:8080'
    }
  },
  css: ['~/assets/css/main.css']
})
