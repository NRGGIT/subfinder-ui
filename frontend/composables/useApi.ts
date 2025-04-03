import { UseFetchOptions } from 'nuxt/app'

export function useApi() {
  const config = useRuntimeConfig()
  const baseUrl = config.public.apiBaseUrl

  /**
   * Wrapper around useFetch for API calls
   */
  async function apiFetch<T>(endpoint: string, options: UseFetchOptions<T> = {}) {
    const url = `${baseUrl}${endpoint}`
    
    return useFetch<T>(url, {
      ...options,
      onResponseError(context) {
        // Handle API errors
        const { response } = context
        const toast = useToast()
        
        toast.add({
          title: 'API Error',
          description: response._data?.error || 'An error occurred while communicating with the API',
          color: 'red',
          icon: 'i-lucide-alert-circle',
          timeout: 5000
        })
        
        if (options.onResponseError) {
          options.onResponseError(context)
        }
      }
    })
  }

  /**
   * Submit a new job
   */
  async function submitJob(domain: string, config: any) {
    return apiFetch('/subfinder', {
      method: 'POST',
      body: {
        domain,
        config
      }
    })
  }

  /**
   * Get job details
   */
  async function getJob(jobId: string) {
    return apiFetch(`/subfinder/${jobId}`)
  }

  /**
   * Get service status
   */
  async function getServiceStatus() {
    return apiFetch('/subfinder/status')
  }

  /**
   * Get all jobs
   */
  async function getAllJobs() {
    return apiFetch('/subfinder/jobs')
  }

  /**
   * Get health status
   */
  async function getHealthStatus() {
    return apiFetch('/health')
  }

  return {
    submitJob,
    getJob,
    getServiceStatus,
    getAllJobs,
    getHealthStatus
  }
}
