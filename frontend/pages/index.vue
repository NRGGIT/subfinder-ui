<template>
  <div>
    <UCard class="mb-6 max-w-3xl mx-auto">
      <template #header>
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold">Submit a New Subdomain Enumeration Job</h2>
        </div>
      </template>
      
      <UForm :state="formState" @submit="onSubmit">
        <!-- Domain Input -->
        <UFormGroup label="Domain" name="domain">
          <UInput v-model="formState.domain" placeholder="example.com" />
        </UFormGroup>
        
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mt-4">
          <!-- Max Depth -->
          <UFormGroup label="Max Depth" name="maxDepth">
            <UInput 
              v-model.number="formState.config.maxDepth" 
              type="number" 
              :min="1" 
              :max="5" 
              placeholder="1-5"
            />
            <template #hint>
              <span class="text-xs">Maximum depth level (1-5)</span>
            </template>
          </UFormGroup>
          
          <!-- Timeout -->
          <UFormGroup label="Timeout (seconds)" name="timeout">
            <UInput 
              v-model.number="formState.config.timeout" 
              type="number" 
              :min="10" 
              :max="300" 
              placeholder="10-300"
            />
            <template #hint>
              <span class="text-xs">Timeout in seconds (10-300)</span>
            </template>
          </UFormGroup>
        </div>
        
        <!-- Rate Limit -->
        <UFormGroup label="Rate Limit" name="rateLimit" class="mt-4">
          <UInput 
            v-model.number="formState.config.rateLimit" 
            type="number" 
            :min="1" 
            :max="50" 
            placeholder="1-50"
          />
          <template #hint>
            <span class="text-xs">Rate limit for requests per second (1-50)</span>
          </template>
        </UFormGroup>
        
        <!-- Sources -->
        <UFormGroup label="Sources" name="sources" class="mt-4">
          <USelectMenu
            v-model="formState.config.sources"
            :options="sourceOptions"
            multiple
            placeholder="Select sources (optional)"
          />
          <template #hint>
            <span class="text-xs">Leave empty to use all available sources</span>
          </template>
        </UFormGroup>
        
        <!-- Boolean Options -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mt-4">
          <UFormGroup name="includeIPs">
            <UCheckbox v-model="formState.config.includeIPs" label="Include IP Addresses" />
          </UFormGroup>
          
          <UFormGroup name="includeWildcards">
            <UCheckbox v-model="formState.config.includeWildcards" label="Include Wildcards" />
          </UFormGroup>
          
          <UFormGroup name="excludeUnresolvable">
            <UCheckbox v-model="formState.config.excludeUnresolvable" label="Exclude Unresolvable" />
          </UFormGroup>
        </div>
        
        <div class="flex justify-end mt-6">
          <UButton type="submit" color="primary" :loading="isSubmitting">
            Submit Job
          </UButton>
        </div>
      </UForm>
    </UCard>
    
    <UCard v-if="recentJobs.length > 0" class="max-w-3xl mx-auto">
      <template #header>
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold">Recent Jobs</h2>
          <UButton to="/status" color="gray" variant="ghost" size="sm">
            View All Jobs
          </UButton>
        </div>
      </template>
      
      <UTable :columns="jobColumns" :rows="recentJobs" class="w-full" hover>
        <template #cell-job_id="{ row }">
          <NuxtLink :to="`/jobs/${row.job_id}`" class="font-mono text-primary hover:underline">
            {{ row.job_id }}
          </NuxtLink>
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
          <UButton
            :to="`/jobs/${row.job_id}`"
            color="primary"
            variant="ghost"
            icon="i-lucide-eye"
            size="sm"
          />
        </template>
      </UTable>
    </UCard>
  </div>
</template>

<script setup lang="ts">
const router = useRouter()
const api = useApi()
const toast = useToast()

// Form state
const formState = reactive({
  domain: '',
  config: {
    maxDepth: 1,
    includeIPs: false,
    sources: [],
    timeout: 60,
    rateLimit: 10,
    includeWildcards: false,
    excludeUnresolvable: false
  }
})

// Source options
const sourceOptions = [
  { value: 'virustotal', label: 'VirusTotal' },
  { value: 'crtsh', label: 'Certificate Search' },
  { value: 'dnsdumpster', label: 'DNS Dumpster' },
  { value: 'threatcrowd', label: 'ThreatCrowd' },
  { value: 'hackertarget', label: 'HackerTarget' },
  { value: 'certspotter', label: 'CertSpotter' },
  { value: 'sublist3r', label: 'Sublist3r' },
  { value: 'threatminer', label: 'ThreatMiner' },
  { value: 'waybackarchive', label: 'Wayback Archive' }
]

// Job submission
const isSubmitting = ref(false)

async function onSubmit() {
  if (!formState.domain) {
    toast.add({
      title: 'Validation Error',
      description: 'Domain is required',
      color: 'red',
      icon: 'i-lucide-alert-circle'
    })
    return
  }
  
  isSubmitting.value = true
  
  try {
    // Transform form state to API format
    const payload = {
      domain: formState.domain,
      config: {
        max_depth: formState.config.maxDepth,
        include_ips: formState.config.includeIPs,
        sources: formState.config.sources.length > 0 ? formState.config.sources : undefined,
        timeout: formState.config.timeout,
        rate_limit: formState.config.rateLimit,
        include_wildcards: formState.config.includeWildcards,
        exclude_unresolvable: formState.config.excludeUnresolvable
      }
    }
    
    const { data, error } = await api.submitJob(formState.domain, payload.config)
    
    if (error.value) {
      throw new Error(error.value.message || 'Failed to submit job')
    }
    
    toast.add({
      title: 'Job Submitted',
      description: `Job ID: ${data.value.job_id}`,
      color: 'green',
      icon: 'i-lucide-check-circle'
    })
    
    // Refresh recent jobs
    await fetchRecentJobs()
    
    // Navigate to job details page
    router.push(`/jobs/${data.value.job_id}`)
  } catch (err) {
    console.error('Job submission error:', err)
    toast.add({
      title: 'Submission Error',
      description: err.message || 'An error occurred while submitting the job',
      color: 'red',
      icon: 'i-lucide-alert-circle'
    })
  } finally {
    isSubmitting.value = false
  }
}

// Recent jobs
const recentJobs = ref([])
const jobColumns = [
  {
    key: 'job_id',
    label: 'Job ID'
  },
  {
    key: 'domain',
    label: 'Domain'
  },
  {
    key: 'status',
    label: 'Status'
  },
  {
    key: 'created_at',
    label: 'Created At'
  },
  {
    key: 'actions',
    label: 'Actions'
  }
]

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

// Fetch recent jobs
async function fetchRecentJobs() {
  try {
    const { data: statusData } = await api.getServiceStatus()
    
    if (statusData.value && statusData.value.jobs && statusData.value.jobs.list) {
      // Get the 5 most recent jobs
      recentJobs.value = statusData.value.jobs.list
        .sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
        .slice(0, 5)
    } else {
      // Fallback to the dedicated jobs endpoint if the list is not included in the status
      const { data: jobsData } = await api.getAllJobs()
      if (jobsData.value && jobsData.value.jobs) {
        // Get the 5 most recent jobs
        recentJobs.value = jobsData.value.jobs
          .sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
          .slice(0, 5)
      } else {
        recentJobs.value = []
      }
    }
  } catch (err) {
    console.error('Error loading recent jobs:', err)
    recentJobs.value = []
  }
}

// Load recent jobs on page load
onMounted(() => {
  fetchRecentJobs()
})
</script>

<style scoped>
.status-badge {
  text-transform: capitalize;
}
</style>
