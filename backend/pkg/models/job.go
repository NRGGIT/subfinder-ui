package models

import (
	"time"
)

// JobStatus represents the current status of a job
type JobStatus string

const (
	JobStatusQueued    JobStatus = "queued"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
)

// SubfinderConfig represents the configuration options for subfinder
type SubfinderConfig struct {
	// Maximum depth level for subdomains (e.g., 2 would include a.example.com and a.b.example.com)
	MaxDepth int `json:"max_depth"`
	
	// Whether to include IP addresses in the results
	IncludeIPs bool `json:"include_ips"`
	
	// List of sources to use (e.g., "virustotal", "crtsh", etc.)
	// If empty, all available sources will be used
	Sources []string `json:"sources"`
	
	// Timeout in seconds for the subfinder operation
	Timeout int `json:"timeout"`
	
	// Rate limit for requests (requests per second)
	RateLimit int `json:"rate_limit"`
	
	// Whether to include wildcard subdomains in the results
	IncludeWildcards bool `json:"include_wildcards"`
	
	// Whether to exclude subdomains that don't resolve
	ExcludeUnresolvable bool `json:"exclude_unresolvable"`
	
	// Whether to exclude subdomains with www prefix
	ExcludeWww bool `json:"exclude_www"`
}

// Job represents a subfinder job
type Job struct {
	// Unique identifier for the job
	ID string `json:"job_id"`
	
	// Domain to search for subdomains
	Domain string `json:"domain"`
	
	// Configuration options for subfinder
	Config SubfinderConfig `json:"config"`
	
	// Current status of the job
	Status JobStatus `json:"status"`
	
	// Time when the job was created
	CreatedAt time.Time `json:"created_at"`
	
	// Time when the job was started
	StartedAt *time.Time `json:"started_at,omitempty"`
	
	// Time when the job was completed
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	
	// Estimated time when the job will be completed
	EstimatedCompletionTime *time.Time `json:"estimated_completion_time,omitempty"`
	
	// Error message if the job failed
	Error string `json:"error,omitempty"`
	
	// List of subdomains found
	Subdomains []string `json:"subdomains,omitempty"`
	
	// Statistics about the job
	Stats *JobStats `json:"stats,omitempty"`
}

// JobStats represents statistics about a job
type JobStats struct {
	// Total number of subdomains found
	TotalFound int `json:"total_found"`
	
	// Time taken to execute the job
	ExecutionTime string `json:"execution_time"`
	
	// Sources used to find subdomains
	SourcesUsed []string `json:"sources_used"`
}

// JobRequest represents a request to create a new job
type JobRequest struct {
	// Domain to search for subdomains
	Domain string `json:"domain" binding:"required"`
	
	// Configuration options for subfinder
	Config SubfinderConfig `json:"config"`
}

// JobResponse represents a response to a job request
type JobResponse struct {
	// Unique identifier for the job
	JobID string `json:"job_id"`
	
	// Current status of the job
	Status JobStatus `json:"status"`
	
	// Estimated time when the job will be completed
	EstimatedCompletionTime *time.Time `json:"estimated_completion_time,omitempty"`
}
