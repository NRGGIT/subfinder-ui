<template>
  <div>
    <div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-6">
      <UCard v-for="(stat, index) in statusStats" :key="index" :ui="{ body: { padding: 'p-4' } }">
        <div class="flex items-center">
          <div class="mr-4 p-3 rounded-full w-12 h-12 flex items-center justify-center" :class="stat.bgClass">
            <UIcon :name="stat.icon" class="text-white text-xl" />
          </div>
          <div>
            <div class="text-sm text-gray-500">{{ stat.label }}</div>
            <div class="text-2xl font-semibold">{{ stat.value }}</div>
          </div>
        </div>
      </UCard>
    </div>
    
    <UCard class="max-w-6xl mx-auto">
      <template #header>
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold">All Jobs</h2>
          <div class="flex items-center gap-2">
            <UInput
              v-model="searchQuery"
              icon="i-lucide-search"
              placeholder="Search jobs..."
              class="w-64"
            />
            <UButton
              color="primary"
              variant="soft"
              icon="i-lucide-refresh-cw"
              :loading="isRefreshing"
              @click="refreshData"
            />
          </div>
        </div>
      </template>
      
      <UTable
        :columns="columns"
        :rows="filteredJobs"
        :loading="isLoading"
        :empty-state="{ icon: 'i-lucide-database', label: 'No jobs found' }"
        @select="onRowSelect"
        hover
      >
        <template #empty-state>
          <div class="flex flex-col items-center justify-center py-6">
            <UIcon name="i-lucide-database" class="text-4xl mb-2 text-gray-400" />
            <p class="text-gray-500">No jobs found</p>
            <UButton
              to="/"
              color="primary"
              variant="ghost"
              class="mt-4"
            >
              Submit a new job
            </UButton>
          </div>
        </template>
        
        <template #cell-job_id="{ row }">
          <NuxtLink :to="`/jobs/${row.job_id}`" class="font-mono text-primary hover:underline">
            {{ row.job_id }}
          </NuxtLink>
        </template>
        
        <template #cell-domain="{ row }">
          <div class="font-medium">{{ row.domain }}</div>
        </template>
        
        <template #cell-status="{ row }">
          <UBadge :color="getStatusColor(row.status)" class="status-badge">
            {{ row.status }}
          </UBadge>
        </template>
        
        <template #cell-created_at="{ row }">
          <div class="text-sm">
            {{ formatDate(row.created_at) }}
          </div>
        </template>
        
        <template #cell-actions="{ row }">
          <div class="flex items-center gap-2">
            <UButton
              :to="`/jobs/${row.job_id}`"
              color="primary"
              variant="ghost"
              icon="i-lucide-eye"
              size="sm"
              aria-label="View job details"
            />
          </div>
        </template>
      </UTable>
    </UCard>
  </div>
</template>

<script setup lang="ts">
const router = useRouter()
const api = useApi()

// Status data
const serviceStatus = ref(null)
const isLoading = ref(true)
const isRefreshing = ref(false)
const jobs = ref([])
const searchQuery = ref('')

// Table columns
const columns = [
  {
    key: 'job_id',
    label: 'Job ID',
    sortable: true
  },
  {
    key: 'domain',
    label: 'Domain',
    sortable: true
  },
  {
    key: 'status',
    label: 'Status',
    sortable: true
  },
  {
    key: 'created_at',
    label: 'Created At',
    sortable: true
  },
  {
    key: 'actions',
    label: 'Actions',
    sortable: false
  }
]

// Computed properties
const statusStats = computed(() => {
  if (!serviceStatus.value) {
    return [
      { label: 'Total Jobs', value: 0, icon: 'i-lucide-database', bgClass: 'bg-blue-500' },
      { label: 'Queued', value: 0, icon: 'i-lucide-clock', bgClass: 'bg-blue-500' },
      { label: 'Running', value: 0, icon: 'i-lucide-loader', bgClass: 'bg-orange-500' },
      { label: 'Completed', value: 0, icon: 'i-lucide-check-circle', bgClass: 'bg-green-500' }
    ]
  }
  
  return [
    { 
      label: 'Total Jobs', 
      value: serviceStatus.value.jobs.total, 
      icon: 'i-lucide-database', 
      bgClass: 'bg-blue-500' 
    },
    { 
      label: 'Queued', 
      value: serviceStatus.value.jobs.queued, 
      icon: 'i-lucide-clock', 
      bgClass: 'bg-blue-500' 
    },
    { 
      label: 'Running', 
      value: serviceStatus.value.jobs.running, 
      icon: 'i-lucide-loader', 
      bgClass: 'bg-orange-500' 
    },
    { 
      label: 'Completed', 
      value: serviceStatus.value.jobs.completed, 
      icon: 'i-lucide-check-circle', 
      bgClass: 'bg-green-500' 
    }
  ]
})

const filteredJobs = computed(() => {
  if (!searchQuery.value) {
    return jobs.value
  }
  
  const query = searchQuery.value.toLowerCase()
  return jobs.value.filter(job => 
    job.job_id.toLowerCase().includes(query) || 
    job.domain.toLowerCase().includes(query) ||
    job.status.toLowerCase().includes(query)
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

function onRowSelect(row) {
  router.push(`/jobs/${row.job_id}`)
}

async function fetchData() {
  isLoading.value = true
  
  try {
    // Get service status
    const { data: statusData } = await api.getServiceStatus()
    serviceStatus.value = statusData.value
    
    // Get job list from the status response
    if (statusData.value && statusData.value.jobs && statusData.value.jobs.list) {
      jobs.value = statusData.value.jobs.list
    } else {
      // Fallback to the dedicated jobs endpoint if the list is not included in the status
      const { data: jobsData } = await api.getAllJobs()
      if (jobsData.value && jobsData.value.jobs) {
        jobs.value = jobsData.value.jobs
      } else {
        jobs.value = []
      }
    }
  } catch (err) {
    console.error('Error fetching status data:', err)
    jobs.value = []
  } finally {
    isLoading.value = false
  }
}

async function refreshData() {
  isRefreshing.value = true
  await fetchData()
  isRefreshing.value = false
}

// Auto-refresh every 30 seconds
let refreshInterval

onMounted(() => {
  fetchData()
  refreshInterval = setInterval(refreshData, 30000)
})

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
