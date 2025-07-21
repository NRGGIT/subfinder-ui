package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/user/subfinder-service/backend/internal/queue"
	"github.com/user/subfinder-service/backend/pkg/models"
)

// Server represents the API server
type Server struct {
	port   string
	router *gin.Engine
	queue  *queue.JobQueue
	logger *log.Logger
	server *http.Server
}

// NewServer creates a new API server
func NewServer(port string, queue *queue.JobQueue, logger *log.Logger) *Server {
	router := gin.Default()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	server := &Server{
		port:   port,
		router: router,
		queue:  queue,
		logger: logger,
	}

	// Set up routes
	server.setupRoutes()

	return server
}

// setupRoutes sets up the API routes
func (s *Server) setupRoutes() {
	// Health check endpoint
	s.router.GET("/health", s.handleHealthCheck)

	// API endpoints
	api := s.router.Group("/subfinder")
	{
		// Submit a new job
		api.POST("", s.handleSubmitJob)

		// Get job status/results
		api.GET("/:id", s.handleGetJob)

		// Get service status
		api.GET("/status", s.handleGetStatus)

		// Get all jobs
		api.GET("/jobs", s.handleGetAllJobs)
	}
}

// Start starts the API server
func (s *Server) Start() error {
	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%s", s.port),
		Handler: s.router,
	}

	s.logger.Printf("Starting API server on port %s", s.port)
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the API server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// handleHealthCheck handles the health check endpoint
func (s *Server) handleHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// handleSubmitJob handles the submit job endpoint
func (s *Server) handleSubmitJob(c *gin.Context) {
	var request models.JobRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Invalid request: %v", err),
		})
		return
	}

	// Validate the domain
	if request.Domain == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Domain is required",
		})
		return
	}

	// Set default configuration values if not provided
	if request.Config.MaxDepth <= 0 {
		request.Config.MaxDepth = 1
	}
	if request.Config.Timeout <= 0 {
		request.Config.Timeout = 60
	}
	if request.Config.RateLimit <= 0 {
		request.Config.RateLimit = 10
	}
	// ExcludeWww is false by default, so no need to set it explicitly

	s.logger.Printf("Received job submission for domain %s", request.Domain)

	// Create a new job
	job := &models.Job{
		ID:        uuid.New().String(),
		Domain:    request.Domain,
		Config:    request.Config,
		Status:    models.JobStatusQueued,
		CreatedAt: time.Now(),
	}

	// Enqueue the job
	if err := s.queue.Enqueue(job); err != nil {
		s.logger.Printf("Failed to enqueue job %s: %v", job.ID, err)
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": fmt.Sprintf("Failed to enqueue job: %v", err),
		})
		return
	}

	s.logger.Printf("Enqueued job %s for domain %s", job.ID, job.Domain)

	// Return the job ID and status
	c.JSON(http.StatusAccepted, models.JobResponse{
		JobID:                   job.ID,
		Status:                  job.Status,
		EstimatedCompletionTime: job.EstimatedCompletionTime,
	})
}

// handleGetJob handles the get job endpoint
func (s *Server) handleGetJob(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Job ID is required",
		})
		return
	}

	s.logger.Printf("Retrieving job %s", id)

	// Get the job from the queue
	job, ok := s.queue.Get(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("Job %s not found", id),
		})
		return
	}

	s.logger.Printf("Job %s status %s", id, job.Status)

	// Return the job
	c.JSON(http.StatusOK, job)
}

// handleGetStatus handles the get status endpoint
func (s *Server) handleGetStatus(c *gin.Context) {
	// Get all jobs
	jobs := s.queue.List()

	s.logger.Printf("Reporting status for %d job(s)", len(jobs))

	// Count jobs by status
	queued := 0
	running := 0
	completed := 0
	failed := 0

	// Create a simplified job list for the response
	jobList := make([]gin.H, 0, len(jobs))
	for _, job := range jobs {
		switch job.Status {
		case models.JobStatusQueued:
			queued++
		case models.JobStatusRunning:
			running++
		case models.JobStatusCompleted:
			completed++
		case models.JobStatusFailed:
			failed++
		}

		// Add job to the list
		jobList = append(jobList, gin.H{
			"job_id":     job.ID,
			"domain":     job.Domain,
			"status":     job.Status,
			"created_at": job.CreatedAt,
		})
	}

	// Return the status and job list
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"jobs": gin.H{
			"total":     len(jobs),
			"queued":    queued,
			"running":   running,
			"completed": completed,
			"failed":    failed,
			"list":      jobList,
		},
		"time": time.Now().Format(time.RFC3339),
	})
}

// handleGetAllJobs handles the get all jobs endpoint
func (s *Server) handleGetAllJobs(c *gin.Context) {
	// Get all jobs
	jobs := s.queue.List()

	s.logger.Printf("Listing %d job(s)", len(jobs))

	// Create a simplified job list for the response
	jobList := make([]gin.H, 0, len(jobs))
	for _, job := range jobs {
		jobList = append(jobList, gin.H{
			"job_id":     job.ID,
			"domain":     job.Domain,
			"status":     job.Status,
			"created_at": job.CreatedAt,
		})
	}

	// Return the job list
	c.JSON(http.StatusOK, gin.H{
		"jobs": jobList,
	})
}
