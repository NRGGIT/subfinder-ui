<template>
  <UApp>
    <div class="min-h-screen">
      <header class="brand-header py-4">
        <div class="page-container flex items-center justify-between">
          <div class="flex items-center">
            <NuxtLink to="/" class="text-xl font-bold text-white">Subfinder UI</NuxtLink>
          </div>
          <nav class="flex items-center space-x-4">
            <NuxtLink 
              to="/" 
              class="text-white hover:text-gray-200 px-3 py-2 rounded-md transition-colors"
              :class="{ 'bg-acv-color-nav-active': $route.path === '/' }"
            >
              <span class="flex items-center">
                <UIcon name="i-lucide-home" class="mr-2" />
                Home
              </span>
            </NuxtLink>
            <NuxtLink 
              to="/status" 
              class="text-white hover:text-gray-200 px-3 py-2 rounded-md transition-colors"
              :class="{ 'bg-acv-color-nav-active': $route.path === '/status' }"
            >
              <span class="flex items-center">
                <UIcon name="i-lucide-activity" class="mr-2" />
                Status
              </span>
            </NuxtLink>
            <NuxtLink 
              to="/api-docs" 
              class="text-white hover:text-gray-200 px-3 py-2 rounded-md transition-colors"
              :class="{ 'bg-acv-color-nav-active': $route.path === '/api-docs' }"
            >
              <span class="flex items-center">
                <UIcon name="i-lucide-code" class="mr-2" />
                API Docs
              </span>
            </NuxtLink>
          </nav>
        </div>
      </header>
      <main class="page-container py-6">
        <NuxtPage />
      </main>
    </div>
  </UApp>
</template>

<script setup>
// App-wide setup code can go here
</script>
