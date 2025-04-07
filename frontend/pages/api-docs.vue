<template>
  <div>
    <UCard class="mb-6 max-w-4xl mx-auto">
      <template #header>
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold">API Documentation</h2>
        </div>
      </template>
      
      <p class="mb-4">
        This page documents the available API endpoints for the Subfinder service. 
        You can use these endpoints to programmatically interact with the service.
      </p>
      
      <div class="mb-8">
        <h3 class="text-md font-semibold mb-2">Base URL</h3>
        <UBadge color="gray" class="font-mono">{{ baseUrl }}</UBadge>
      </div>
      
      <!-- API Endpoints -->
      <div v-for="(endpoint, index) in endpoints" :key="index" class="mb-8 border-b border-gray-200 pb-6">
        <div class="flex items-start mb-2">
          <UBadge :color="getMethodColor(endpoint.method)" class="mr-2">{{ endpoint.method }}</UBadge>
          <code class="font-mono bg-gray-100 px-2 py-1 rounded">{{ endpoint.path }}</code>
        </div>
        
        <p class="mb-4">{{ endpoint.description }}</p>
        
        <!-- Request Parameters -->
        <div v-if="endpoint.request" class="mb-4">
          <h4 class="text-sm font-semibold mb-2">Request</h4>
          
          <!-- Path Parameters -->
          <div v-if="endpoint.request.pathParams && endpoint.request.pathParams.length > 0" class="mb-4">
            <h5 class="text-xs font-semibold mb-1">Path Parameters</h5>
            <UTable :columns="paramColumns" :rows="endpoint.request.pathParams" class="w-full text-sm">
              <template #cell-required="{ row }">
                <UBadge v-if="row.required" color="red" size="xs">Required</UBadge>
                <UBadge v-else color="gray" size="xs">Optional</UBadge>
              </template>
            </UTable>
          </div>
          
          <!-- Request Body -->
          <div v-if="endpoint.request.body" class="mb-4">
            <h5 class="text-xs font-semibold mb-1">Request Body</h5>
            <UCard class="bg-gray-50">
              <pre class="text-xs overflow-auto p-2"><code>{{ formatJson(endpoint.request.body) }}</code></pre>
            </UCard>
            
            <!-- Body Parameters -->
            <div v-if="endpoint.request.bodyParams && endpoint.request.bodyParams.length > 0" class="mt-2">
              <UTable :columns="paramColumns" :rows="endpoint.request.bodyParams" class="w-full text-sm">
                <template #cell-required="{ row }">
                  <UBadge v-if="row.required" color="red" size="xs">Required</UBadge>
                  <UBadge v-else color="gray" size="xs">Optional</UBadge>
                </template>
              </UTable>
            </div>
          </div>
        </div>
        
        <!-- Response -->
        <div v-if="endpoint.response" class="mb-4">
          <h4 class="text-sm font-semibold mb-2">Response</h4>
          
          <!-- Status Codes -->
          <div v-if="endpoint.response.statusCodes && endpoint.response.statusCodes.length > 0" class="mb-2">
            <h5 class="text-xs font-semibold mb-1">Status Codes</h5>
            <UTable :columns="statusColumns" :rows="endpoint.response.statusCodes" class="w-full text-sm">
              <template #cell-code="{ row }">
                <UBadge :color="getStatusColor(row.code)" size="xs">{{ row.code }}</UBadge>
              </template>
            </UTable>
          </div>
          
          <!-- Response Body -->
          <div v-if="endpoint.response.body" class="mb-4">
            <h5 class="text-xs font-semibold mb-1">Response Body</h5>
            <UCard class="bg-gray-50">
              <pre class="text-xs overflow-auto p-2"><code>{{ formatJson(endpoint.response.body) }}</code></pre>
            </UCard>
            
            <!-- Response Fields -->
            <div v-if="endpoint.response.fields && endpoint.response.fields.length > 0" class="mt-2">
              <h5 class="text-xs font-semibold mb-1">Response Fields</h5>
              <UTable :columns="paramColumns" :rows="endpoint.response.fields" class="w-full text-sm">
                <template #cell-required="{ row }">
                  <UBadge v-if="row.required" color="red" size="xs">Required</UBadge>
                  <UBadge v-else color="gray" size="xs">Optional</UBadge>
                </template>
              </UTable>
            </div>
          </div>
        </div>
        
        <!-- Example -->
        <div v-if="endpoint.example" class="mb-4">
          <h4 class="text-sm font-semibold mb-2">Example</h4>
          
          <UCard class="bg-gray-50 mb-2">
            <div class="text-xs font-semibold mb-1">Request</div>
            <pre class="text-xs overflow-auto p-2"><code>{{ endpoint.example.request }}</code></pre>
          </UCard>
          
          <UCard class="bg-gray-50">
            <div class="text-xs font-semibold mb-1">Response</div>
            <pre class="text-xs overflow-auto p-2"><code>{{ endpoint.example.response }}</code></pre>
          </UCard>
        </div>
      </div>
    </UCard>
  </div>
</template>

<script setup lang="ts">
const config = useRuntimeConfig()
const baseUrl = config.public.apiBaseUrl

// Table columns
const paramColumns = [
  { key: 'name', label: 'Parameter' },
  { key: 'type', label: 'Type' },
  { key: 'description', label: 'Description' },
  { key: 'required', label: 'Required' }
]

const statusColumns = [
  { key: 'code', label: 'Code' },
  { key: 'description', label: 'Description' }
]

// Helper functions
function getMethodColor(method) {
  const colors = {
    GET: 'blue',
    POST: 'green',
    PUT: 'orange',
    DELETE: 'red',
    PATCH: 'purple'
  }
  
  return colors[method] || 'gray'
}

function getStatusColor(code) {
  if (code >= 200 && code < 300) return 'green'
  if (code >= 300 && code < 400) return 'blue'
  if (code >= 400 && code < 500) return 'orange'
  if (code >= 500) return 'red'
  return 'gray'
}

function formatJson(obj) {
  return JSON.stringify(obj, null, 2)
}

// API Endpoints
const endpoints = [
  {
    method: 'POST',
    path: '/subfinder',
    description: 'Submit a new subdomain enumeration job. This endpoint accepts a domain and configuration options, and returns a job ID that can be used to track the job status.',
    request: {
      body: {
        domain: 'example.com',
        config: {
          max_depth: 1,
          include_ips: false,
          sources: ['virustotal', 'crtsh'],
          timeout: 60,
          rate_limit: 10,
          include_wildcards: false,
          exclude_unresolvable: false,
          exclude_www: false
        }
      },
      bodyParams: [
        { name: 'domain', type: 'string', description: 'Domain to search for subdomains', required: true },
        { name: 'config', type: 'object', description: 'Configuration options for subfinder', required: false },
        { name: 'config.max_depth', type: 'number', description: 'Maximum depth level for subdomains (1-5)', required: false },
        { name: 'config.include_ips', type: 'boolean', description: 'Whether to include IP addresses in the results', required: false },
        { name: 'config.sources', type: 'array', description: 'List of sources to use (e.g., "virustotal", "crtsh", etc.)', required: false },
        { name: 'config.timeout', type: 'number', description: 'Timeout in seconds for the subfinder operation (10-300)', required: false },
        { name: 'config.rate_limit', type: 'number', description: 'Rate limit for requests per second (1-50)', required: false },
        { name: 'config.include_wildcards', type: 'boolean', description: 'Whether to include wildcard subdomains in the results', required: false },
        { name: 'config.exclude_unresolvable', type: 'boolean', description: 'Whether to exclude subdomains that don\'t resolve', required: false },
        { name: 'config.exclude_www', type: 'boolean', description: 'Whether to exclude subdomains with www prefix', required: false }
      ]
    },
    response: {
      statusCodes: [
        { code: 202, description: 'Job accepted and queued' },
        { code: 400, description: 'Invalid request parameters' },
        { code: 503, description: 'Service unavailable, failed to enqueue job' }
      ],
      body: {
        job_id: '123e4567-e89b-12d3-a456-426614174000',
        status: 'queued',
        estimated_completion_time: '2025-04-04T10:15:30Z'
      },
      fields: [
        { name: 'job_id', type: 'string', description: 'Unique identifier for the job', required: true },
        { name: 'status', type: 'string', description: 'Current status of the job (queued, running, completed, failed)', required: true },
        { name: 'estimated_completion_time', type: 'string', description: 'Estimated time when the job will be completed (ISO 8601 format)', required: false }
      ]
    },
    example: {
      request: `curl -X POST ${config.public.apiBaseUrl}/subfinder \\
  -H "Content-Type: application/json" \\
  -d '{
    "domain": "example.com",
    "config": {
      "max_depth": 2,
      "include_ips": true,
      "timeout": 120,
      "rate_limit": 20
    }
  }'`,
      response: `{
  "job_id": "123e4567-e89b-12d3-a456-426614174000",
  "status": "queued",
  "estimated_completion_time": "2025-04-04T10:15:30Z"
}`
    }
  },
  {
    method: 'GET',
    path: '/subfinder/:id',
    description: 'Get the status and results of a specific job. This endpoint returns detailed information about the job, including any subdomains found if the job is completed.',
    request: {
      pathParams: [
        { name: 'id', type: 'string', description: 'Job ID', required: true }
      ]
    },
    response: {
      statusCodes: [
        { code: 200, description: 'Job found' },
        { code: 400, description: 'Invalid job ID' },
        { code: 404, description: 'Job not found' }
      ],
      body: {
        job_id: '123e4567-e89b-12d3-a456-426614174000',
        domain: 'example.com',
        config: {
          max_depth: 2,
          include_ips: true,
          sources: [],
          timeout: 120,
          rate_limit: 20,
          include_wildcards: false,
          exclude_unresolvable: true,
          exclude_www: false
        },
        status: 'completed',
        created_at: '2025-04-04T10:00:00Z',
        started_at: '2025-04-04T10:00:05Z',
        completed_at: '2025-04-04T10:02:30Z',
        subdomains: [
          { subdomain: 'api.example.com', ip: '192.168.1.10', source: 'virustotal' },
          { subdomain: 'blog.example.com', ip: '192.168.1.11', source: 'crtsh' },
          { subdomain: 'mail.example.com', ip: '192.168.1.12', source: 'dnsdumpster' },
          { subdomain: 'dev.example.com', ip: '192.168.1.13', source: 'virustotal' }
        ],
        stats: {
          total_found: 4,
          execution_time: '2m 25s',
          sources_used: ['virustotal', 'crtsh', 'dnsdumpster']
        }
      },
      fields: [
        { name: 'job_id', type: 'string', description: 'Unique identifier for the job', required: true },
        { name: 'domain', type: 'string', description: 'Domain that was searched', required: true },
        { name: 'config', type: 'object', description: 'Configuration options used for the job', required: true },
        { name: 'status', type: 'string', description: 'Current status of the job (queued, running, completed, failed)', required: true },
        { name: 'created_at', type: 'string', description: 'Time when the job was created (ISO 8601 format)', required: true },
        { name: 'started_at', type: 'string', description: 'Time when the job was started (ISO 8601 format)', required: false },
        { name: 'completed_at', type: 'string', description: 'Time when the job was completed (ISO 8601 format)', required: false },
        { name: 'error', type: 'string', description: 'Error message if the job failed', required: false },
        { name: 'subdomains', type: 'array', description: 'List of subdomain objects found', required: false },
        { name: 'subdomains[].subdomain', type: 'string', description: 'The found subdomain name', required: true },
        { name: 'subdomains[].ip', type: 'string', description: 'IP address of the subdomain (only present if config.include_ips is true)', required: false },
        { name: 'subdomains[].source', type: 'string', description: 'Source where the subdomain was found', required: true },
        { name: 'stats', type: 'object', description: 'Statistics about the job', required: false },
        { name: 'stats.total_found', type: 'number', description: 'Total number of subdomains found', required: false },
        { name: 'stats.execution_time', type: 'string', description: 'Time taken to execute the job', required: false },
        { name: 'stats.sources_used', type: 'array', description: 'Sources used to find subdomains', required: false }
      ]
    },
    example: {
      request: `curl ${config.public.apiBaseUrl}/subfinder/123e4567-e89b-12d3-a456-426614174000`,
      response: `{
  "job_id": "123e4567-e89b-12d3-a456-426614174000",
  "domain": "example.com",
  "config": {
    "max_depth": 2,
    "include_ips": true,
    "sources": [],
    "timeout": 120,
    "rate_limit": 20,
    "include_wildcards": false,
    "exclude_unresolvable": true,
    "exclude_www": false
  },
  "status": "completed",
  "created_at": "2025-04-04T10:00:00Z",
  "started_at": "2025-04-04T10:00:05Z",
  "completed_at": "2025-04-04T10:02:30Z",
  "subdomains": [
    { "subdomain": "api.example.com", "ip": "192.168.1.10", "source": "virustotal" },
    { "subdomain": "blog.example.com", "ip": "192.168.1.11", "source": "crtsh" },
    { "subdomain": "mail.example.com", "ip": "192.168.1.12", "source": "dnsdumpster" },
    { "subdomain": "dev.example.com", "ip": "192.168.1.13", "source": "virustotal" }
  ],
  "stats": {
    "total_found": 4,
    "execution_time": "2m 25s",
    "sources_used": ["virustotal", "crtsh", "dnsdumpster"]
  }
}`
    }
  },
  {
    method: 'GET',
    path: '/subfinder/status',
    description: 'Get the current status of the service, including statistics about all jobs. This endpoint provides an overview of the service\'s current state and workload.',
    response: {
      statusCodes: [
        { code: 200, description: 'Status retrieved successfully' }
      ],
      body: {
        status: 'ok',
        jobs: {
          total: 10,
          queued: 2,
          running: 1,
          completed: 6,
          failed: 1,
          list: [
            {
              job_id: '123e4567-e89b-12d3-a456-426614174000',
              domain: 'example.com',
              status: 'completed',
              created_at: '2025-04-04T10:00:00Z'
            },
            {
              job_id: '223e4567-e89b-12d3-a456-426614174001',
              domain: 'test.com',
              status: 'running',
              created_at: '2025-04-04T10:05:00Z'
            }
          ]
        },
        time: '2025-04-04T10:10:00Z'
      },
      fields: [
        { name: 'status', type: 'string', description: 'Status of the service (ok)', required: true },
        { name: 'jobs', type: 'object', description: 'Information about all jobs', required: true },
        { name: 'jobs.total', type: 'number', description: 'Total number of jobs', required: true },
        { name: 'jobs.queued', type: 'number', description: 'Number of queued jobs', required: true },
        { name: 'jobs.running', type: 'number', description: 'Number of running jobs', required: true },
        { name: 'jobs.completed', type: 'number', description: 'Number of completed jobs', required: true },
        { name: 'jobs.failed', type: 'number', description: 'Number of failed jobs', required: true },
        { name: 'jobs.list', type: 'array', description: 'List of all jobs', required: true },
        { name: 'time', type: 'string', description: 'Current time (ISO 8601 format)', required: true }
      ]
    },
    example: {
      request: `curl ${config.public.apiBaseUrl}/subfinder/status`,
      response: `{
  "status": "ok",
  "jobs": {
    "total": 10,
    "queued": 2,
    "running": 1,
    "completed": 6,
    "failed": 1,
    "list": [
      {
        "job_id": "123e4567-e89b-12d3-a456-426614174000",
        "domain": "example.com",
        "status": "completed",
        "created_at": "2025-04-04T10:00:00Z"
      },
      {
        "job_id": "223e4567-e89b-12d3-a456-426614174001",
        "domain": "test.com",
        "status": "running",
        "created_at": "2025-04-04T10:05:00Z"
      }
    ]
  },
  "time": "2025-04-04T10:10:00Z"
}`
    }
  },
  {
    method: 'GET',
    path: '/subfinder/jobs',
    description: 'Get a list of all jobs. This endpoint returns a simplified list of all jobs in the system.',
    response: {
      statusCodes: [
        { code: 200, description: 'Jobs retrieved successfully' }
      ],
      body: {
        jobs: [
          {
            job_id: '123e4567-e89b-12d3-a456-426614174000',
            domain: 'example.com',
            status: 'completed',
            created_at: '2025-04-04T10:00:00Z'
          },
          {
            job_id: '223e4567-e89b-12d3-a456-426614174001',
            domain: 'test.com',
            status: 'running',
            created_at: '2025-04-04T10:05:00Z'
          }
        ]
      },
      fields: [
        { name: 'jobs', type: 'array', description: 'List of all jobs', required: true },
        { name: 'jobs[].job_id', type: 'string', description: 'Unique identifier for the job', required: true },
        { name: 'jobs[].domain', type: 'string', description: 'Domain that was searched', required: true },
        { name: 'jobs[].status', type: 'string', description: 'Current status of the job', required: true },
        { name: 'jobs[].created_at', type: 'string', description: 'Time when the job was created (ISO 8601 format)', required: true }
      ]
    },
    example: {
      request: `curl ${config.public.apiBaseUrl}/subfinder/jobs`,
      response: `{
  "jobs": [
    {
      "job_id": "123e4567-e89b-12d3-a456-426614174000",
      "domain": "example.com",
      "status": "completed",
      "created_at": "2025-04-04T10:00:00Z"
    },
    {
      "job_id": "223e4567-e89b-12d3-a456-426614174001",
      "domain": "test.com",
      "status": "running",
      "created_at": "2025-04-04T10:05:00Z"
    }
  ]
}`
    }
  },
  {
    method: 'GET',
    path: '/health',
    description: 'Check the health of the service. This endpoint can be used for monitoring and health checks.',
    response: {
      statusCodes: [
        { code: 200, description: 'Service is healthy' }
      ],
      body: {
        status: 'ok',
        time: '2025-04-04T10:10:00Z'
      },
      fields: [
        { name: 'status', type: 'string', description: 'Status of the service (ok)', required: true },
        { name: 'time', type: 'string', description: 'Current time (ISO 8601 format)', required: true }
      ]
    },
    example: {
      request: `curl ${config.public.apiBaseUrl}/health`,
      response: `{
  "status": "ok",
  "time": "2025-04-04T10:10:00Z"
}`
    }
  }
]
</script>
