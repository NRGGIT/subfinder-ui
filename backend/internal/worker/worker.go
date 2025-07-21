package worker

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/user/subfinder-service/backend/internal/queue"
	"github.com/user/subfinder-service/backend/internal/subfinder"
	"github.com/user/subfinder-service/backend/pkg/models"
)

// WorkerPool represents a pool of workers that process jobs from a queue
type WorkerPool struct {
	count     int
	queue     *queue.JobQueue
	logger    *log.Logger
	wg        sync.WaitGroup
	subfinder *subfinder.Client
}

// NewWorkerPool creates a new worker pool with the specified number of workers
func NewWorkerPool(count int, queue *queue.JobQueue, logger *log.Logger) *WorkerPool {
	return &WorkerPool{
		count:     count,
		queue:     queue,
		logger:    logger,
		subfinder: subfinder.NewClient(logger),
	}
}

// Start starts the worker pool
func (p *WorkerPool) Start(ctx context.Context) {
	p.logger.Printf("Starting worker pool with %d workers", p.count)

	for i := 0; i < p.count; i++ {
		p.wg.Add(1)
		go p.worker(ctx, i)
	}
}

// Wait waits for all workers to finish
func (p *WorkerPool) Wait() {
	p.wg.Wait()
}

// worker processes jobs from the queue
func (p *WorkerPool) worker(ctx context.Context, id int) {
	defer func() {
		if r := recover(); r != nil {
			p.logger.Printf("Worker %d recovered from panic: %v", id, r)
		}
		p.wg.Done()
	}()

	p.logger.Printf("Worker %d started", id)

	for {
		select {
		case <-ctx.Done():
			p.logger.Printf("Worker %d stopped", id)
			return
		default:
			// Try to get a job from the queue
			jobID, ok := p.queue.Dequeue()
			if !ok {
				// No jobs in the queue, wait a bit and try again
				time.Sleep(100 * time.Millisecond)
				continue
			}
			p.logger.Printf("Worker %d dequeued job %s", id, jobID)

			// Get the job from the queue
			job, ok := p.queue.Get(jobID)
			if !ok {
				p.logger.Printf("Worker %d: Job %s not found", id, jobID)
				continue
			}

			// Process the job
			p.processJob(ctx, job)
		}
	}
}

// processJob processes a job
func (p *WorkerPool) processJob(ctx context.Context, job *models.Job) {
	p.logger.Printf("Processing job %s for domain %s with config %+v", job.ID, job.Domain, job.Config)

	// Update job status to running
	now := time.Now()
	job.Status = models.JobStatusRunning
	job.StartedAt = &now

	// Estimate completion time based on domain complexity
	// This is a simple estimation, could be improved with historical data
	estimatedDuration := 30 * time.Second // Default estimation
	estimatedCompletionTime := now.Add(estimatedDuration)
	job.EstimatedCompletionTime = &estimatedCompletionTime
	p.logger.Printf("Job %s estimated completion at %s", job.ID, estimatedCompletionTime.Format(time.RFC3339))

	p.queue.Update(job)

	// Create a context with timeout from the job configuration
	jobCtx := ctx
	if job.Config.Timeout > 0 {
		var cancel context.CancelFunc
		jobCtx, cancel = context.WithTimeout(ctx, time.Duration(job.Config.Timeout)*time.Second)
		defer cancel()
		p.logger.Printf("Job %s timeout set to %ds", job.ID, job.Config.Timeout)
	}

	// Run subfinder
	startTime := time.Now()
	subdomains, sourcesUsed, err := p.subfinder.FindSubdomains(jobCtx, job.Domain, job.Config)
	executionTime := time.Since(startTime)

	// Update job with results
	now = time.Now()
	job.CompletedAt = &now

	if err != nil {
		job.Status = models.JobStatusFailed
		job.Error = err.Error()
		p.logger.Printf("Job %s failed after %s: %v", job.ID, executionTime.String(), err)
	} else {
		job.Status = models.JobStatusCompleted
		job.Subdomains = subdomains
		job.Stats = &models.JobStats{
			TotalFound:    len(subdomains),
			ExecutionTime: executionTime.String(),
			SourcesUsed:   sourcesUsed,
		}
		p.logger.Printf("Job %s completed in %s: found %d subdomains", job.ID, executionTime.String(), len(subdomains))
	}

	p.queue.Update(job)
}
