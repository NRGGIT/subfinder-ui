<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <UBreadcrumb :items="breadcrumbItems" />
      <UButton
        to="/status"
        color="gray"
        variant="ghost"
        icon="i-lucide-arrow-left"
        label="Back to Jobs"
      />
    </div>
    
    <div v-if="isLoading">
      <USkeleton class="h-32 mb-4" />
      <USkeleton class="h-64" />
    </div>
    
    <template v-else-if="job">
      <!-- Job Header -->
      <UCard class="mb-6">
        <template #header>
          <div class="flex items-center justify-between">
            <div>
              <h2 class="text-lg font-semibold">Job Details</h2>
              <p class="text-sm text-gray-500">ID: {{ job.job_id }}</p>
            </div>
            <UBadge :color="getStatusColor(job.status)" class="status-badge text-sm">
              {{ job.status }}
            </UBadge>
          </div>
        </template>
        
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <h3 class="font-medium mb-2">Domain</h3>
            <p class="text-lg font-mono">{{ job.domain }}</p>
          </div>
          
          <div>
            <h3 class="font-medium mb-2">Timing</h3>
            <div class="space-y-1">
              <div class="flex items-center text-sm">
                <UIcon name="i-lucide-calendar" class="mr-2 text-gray-500" />
                <span class="text-gray-500">Created:</span>
                <span class="ml-2">{{ formatDate(job.created_at) }}</span>
              </div>
              
              <div v-if="job.started_at" class="flex items-center text-sm">
                <UIcon name="i-lucide-play" class="mr-2 text-gray-500" />
                <span class="text-gray-500">Started:</span>
                <span class="ml-2">{{ formatDate(job.started_at) }}</span>
              </div>
              
              <div v-if="job.completed_at" class="flex items-center text-sm">
                <UIcon name="i-lucide-check" class="mr-2 text-gray-500" />
                <span class="text-gray-500">Completed:</span>
                <span class="ml-2">{{ formatDate(job.completed_at) }}</span>
              </div>
              
              <div v-if="job.stats?.execution_time" class="flex items-center text-sm">
                <UIcon name="i-lucide-clock" class="mr-2 text-gray-500" />
                <span class="text-gray-500">Execution Time:</span>
                <span class="ml-2">{{ job.stats.execution_time }}</span>
              </div>
            </div>
          </div>
        </div>
        
        <template v-if="job.status === 'running'">
          <div class="mt-6">
            <p class="text-sm text-gray-500 mb-2">Job in progress...</p>
            <UProgress indeterminate color="primary" />
          </div>
        </template>
        
        <template v-if="job.error">
          <div class="mt-6 p-4 bg-red-50 text-red-700 rounded-md">
            <div class="flex items-start">
              <UIcon name="i-lucide-alert-circle" class="mr-2 text-red-500 mt-0.5" />
              <div>
                <h4 class="font-medium">Error</h4>
                <p>{{ job.error }}</p>
              </div>
            </div>
          </div>
        </template>
      </UCard>
      
      <!-- Configuration -->
      <UCard class="mb-6">
        <template #header>
          <h3 class="text-lg font-semibold">Configuration</h3>
        </template>
        
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <h4 class="text-sm font-medium text-gray-500 mb-1">Max Depth</h4>
            <p>{{ job.config.max_depth }}</p>
          </div>
          
          <div>
            <h4 class="text-sm font-medium text-gray-500 mb-1">Timeout</h4>
            <p>{{ job.config.timeout }} seconds</p>
          </div>
          
          <div>
            <h4 class="text-sm font-medium text-gray-500 mb-1">Rate Limit</h4>
            <p>{{ job.config.rate_limit }} req/sec</p>
          </div>
          
          <div>
            <h4 class="text-sm font-medium text-gray-500 mb-1">Include IPs</h4>
            <UBadge :color="job.config.include_ips ? 'green' : 'gray'" size="sm">
              {{ job.config.include_ips ? 'Yes' : 'No' }}
            </UBadge>
          </div>
          
          <div>
            <h4 class="text-sm font-medium text-gray-500 mb-1">Include Wildcards</h4>
            <UBadge :color="job.config.include_wildcards ? 'green' : 'gray'" size="sm">
              {{ job.config.include_wildcards ? 'Yes' : 'No' }}
            </UBadge>
          </div>
          
          <div>
            <h4 class="text-sm font-medium text-gray-500 mb-1">Exclude Unresolvable</h4>
            <UBadge :color="job.config.exclude_unresolvable ? 'green' : 'gray'" size="sm">
              {{ job.config.exclude_unresolvable ? 'Yes' : 'No' }}
            </UBadge>
          </div>
          
          <div>
            <h4 class="text-sm font-medium text-gray-500 mb-1">Exclude WWW Subdomains</h4>
            <UBadge :color="job.config.exclude_www ? 'green' : 'gray'" size="sm">
              {{ job.config.exclude_www ? 'Yes' : 'No' }}
            </UBadge>
          </div>
          
          <div class="md:col-span-3">
            <h4 class="text-sm font-medium text-gray-500 mb-1">Sources</h4>
            <div class="flex flex-wrap gap-2">
              <template v-if="job.config.sources && job.config.sources.length">
                <UBadge
                  v-for="source in job.config.sources"
                  :key="source"
                  color="blue"
                  variant="subtle"
                  size="sm"
                >
                  {{ source }}
                </UBadge>
              </template>
              <template v-else>
                <p class="text-sm text-gray-500">All available sources</p>
              </template>
            </div>
          </div>
        </div>
      </UCard>
      
      <!-- Results -->
      <UCard v-if="job.status === 'completed' && job.subdomains">
        <template #header>
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold">Results</h3>
            <div class="flex items-center gap-2">
              <UBadge color="blue" size="sm">
                {{ job.subdomains.length }} subdomains found
              </UBadge>
              <UButton
                v-if="job.subdomains.length"
                color="gray"
                variant="ghost"
                icon="i-lucide-download"
                size="sm"
                @click="downloadResults"
              >
                Download
              </UButton>
            </div>
          </div>
        </template>
        
        <div v-if="job.subdomains.length === 0" class="py-8 text-center">
          <UIcon name="i-lucide-search-x" class="mx-auto mb-2 text-4xl text-gray-400" />
          <p class="text-gray-500">No subdomains found for this domain</p>
        </div>
        
        <div v-else>
          <UInput
            v-model="searchQuery"
            icon="i-lucide-search"
            placeholder="Search subdomains..."
            class="mb-4"
          />
          
          <div class="border rounded-md overflow-hidden">
            <div class="max-h-96 overflow-y-auto">
              <table class="w-full">
                <thead class="bg-gray-50 sticky top-0">
                  <tr>
                    <th class="px-4 py-2 text-left text-sm font-medium text-gray-500">#</th>
                    <th class="px-4 py-2 text-left text-sm font-medium text-gray-500">Subdomain</th>
                    <th v-if="job?.config?.include_ips" class="px-4 py-2 text-left text-sm font-medium text-gray-500">IP Address</th>
                    <th class="px-4 py-2 text-left text-sm font-medium text-gray-500">Source</th>
                  </tr>
                </thead>
                <tbody class="divide-y">
                  <tr
                    v-for="(result, index) in filteredSubdomains"
                    :key="result.subdomain + '-' + index"
                    class="hover:bg-gray-50"
                  >
                    <td class="px-4 py-2 text-sm text-gray-500">{{ index + 1 }}</td>
                    <td class="px-4 py-2 font-mono">{{ result.subdomain }}</td>
                    <td v-if="job?.config?.include_ips" class="px-4 py-2 font-mono">{{ result.ip || 'N/A' }}</td>
                    <td class="px-4 py-2 text-sm">{{ result.source }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </UCard>
      
      <!-- Stats -->
      <UCard v-if="job.stats && job.status === 'completed'" class="mt-6">
        <template #header>
          <h3 class="text-lg font-semibold">Statistics</h3>
        </template>
        
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <h4 class="text-sm font-medium text-gray-500 mb-1">Total Found</h4>
            <p class="text-2xl font-semibold">{{ job.stats.total_found }}</p>
          </div>
          
          <div>
            <h4 class="text-sm font-medium text-gray-500 mb-1">Execution Time</h4>
            <p class="text-2xl font-semibold">{{ job.stats.execution_time }}</p>
          </div>
          
          <div>
            <h4 class="text-sm font-medium text-gray-500 mb-1">Sources Used</h4>
            <div class="flex flex-wrap gap-2">
              <UBadge
                v-for="source in job.stats.sources_used"
                :key="source"
                color="blue"
                variant="subtle"
                size="sm"
              >
                {{ source }}
              </UBadge>
            </div>
          </div>
        </div>
      </UCard>
    </template>
    
    <div v-else class="text-center py-12">
      <UIcon name="i-lucide-file-x" class="mx-auto mb-4 text-6xl text-gray-300" />
      <h3 class="text-lg font-medium text-gray-900 mb-1">Job Not Found</h3>
      <p class="text-gray-500 mb-6">The job you're looking for doesn't exist or has been removed.</p>
      <UButton to="/" color="primary">Go to Home</UButton>
    </div>
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const api = useApi()
const toast = useToast()

const jobId = computed(() => route.params.id as string)
const isLoading = ref(true)
const job = ref(null)
const searchQuery = ref('')

// Computed properties
const breadcrumbItems = computed(() => [
  { label: 'Home', to: '/' },
  { label: 'Jobs', to: '/status' },
  { label: `Job ${jobId.value}`, to: route.fullPath }
])

const filteredSubdomains = computed(() => {
  if (!job.value?.subdomains) {
    return []
  }
  if (!searchQuery.value) {
    return job.value.subdomains
  }
  
  const query = searchQuery.value.toLowerCase()
  return job.value.subdomains.filter(result =>
    result.subdomain.toLowerCase().includes(query) ||
    (result.ip && result.ip.toLowerCase().includes(query)) ||
    result.source.toLowerCase().includes(query)
  )
})

// Methods
function getStatusColor(status) {
  const statusColors = {
    queued: 'blue',
    running: 'orange',
    completed: 'green',
    failed: 'red'
  }
  
  return statusColors[status] || 'gray'
}

function formatDate(dateString) {
  if (!dateString) return ''
  
  const date = new Date(dateString)
  return date.toLocaleString()
}

function downloadResults() {
  if (!job.value?.subdomains) return
  
  // Format as CSV
  const headers = ['Subdomain']
  if (job.value.config.include_ips) {
    headers.push('IP Address')
  }
  headers.push('Source')
  
  const csvContent = [
    headers.join(','),
    ...job.value.subdomains.map(result => {
      const row = [result.subdomain]
      if (job.value.config.include_ips) {
        row.push(result.ip || '') // Add IP if included, else empty string
      }
      row.push(result.source)
      // Escape commas within fields if necessary (basic example)
      return row.map(field => `"${field.replace(/"/g, '""')}"`).join(',')
    })
  ].join('\n')
  
  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)
  
  const a = document.createElement('a')
  a.href = url
  a.download = `subdomains-${job.value.domain}-${job.value.job_id}.csv` // Change extension to .csv
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
  
  toast.add({
    title: 'Download Started',
    description: 'Your results are being downloaded',
    color: 'green',
    icon: 'i-lucide-download'
  })
}

// Fetch job data
async function fetchJob() {
  isLoading.value = true
  
  try {
    const { data, error } = await api.getJob(jobId.value)
    
    if (error.value) {
      throw new Error(error.value.message || 'Failed to fetch job')
    }
    
    job.value = data.value
  } catch (err) {
    console.error('Error fetching job:', err)
    toast.add({
      title: 'Error',
      description: err.message || 'Failed to load job details',
      color: 'red',
      icon: 'i-lucide-alert-circle'
    })
    job.value = null
  } finally {
    isLoading.value = false
  }
}

// Auto-refresh for running jobs
let refreshInterval

function setupRefresh() {
  if (job.value?.status === 'running' || job.value?.status === 'queued') {
    refreshInterval = setInterval(fetchJob, 5000)
  } else if (refreshInterval) {
    clearInterval(refreshInterval)
  }
}

onMounted(() => {
  fetchJob()
})

watch(() => job.value?.status, setupRefresh, { immediate: true })

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
})
</script>

<style scoped>
.status-badge {
  text-transform: capitalize;
}
</style>
