package queue

import (
	"sync"

	"github.com/user/subfinder-service/backend/pkg/models"
)

// JobQueue represents a queue of jobs to be processed
type JobQueue struct {
	jobs     map[string]*models.Job
	queue    chan string
	mutex    sync.RWMutex
	capacity int
}

// NewJobQueue creates a new job queue with the specified capacity
func NewJobQueue() *JobQueue {
	return &JobQueue{
		jobs:     make(map[string]*models.Job),
		queue:    make(chan string, 100), // Buffer size of 100 jobs
		capacity: 100,
	}
}

// Enqueue adds a job to the queue
func (q *JobQueue) Enqueue(job *models.Job) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	// Store the job in the map
	q.jobs[job.ID] = job

	// Add the job ID to the queue
	select {
	case q.queue <- job.ID:
		return nil
	default:
		// Queue is full
		return ErrQueueFull
	}
}

// Dequeue removes a job from the queue and returns it
func (q *JobQueue) Dequeue() (string, bool) {
	select {
	case jobID := <-q.queue:
		return jobID, true
	default:
		// Queue is empty
		return "", false
	}
}

// Get returns a job by ID
func (q *JobQueue) Get(id string) (*models.Job, bool) {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	job, ok := q.jobs[id]
	return job, ok
}

// Update updates a job in the queue
func (q *JobQueue) Update(job *models.Job) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.jobs[job.ID] = job
}

// Size returns the number of jobs in the queue
func (q *JobQueue) Size() int {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	return len(q.jobs)
}

// List returns a list of all jobs
func (q *JobQueue) List() []*models.Job {
	q.mutex.RLock()
	defer q.mutex.RUnlock()

	jobs := make([]*models.Job, 0, len(q.jobs))
	for _, job := range q.jobs {
		jobs = append(jobs, job)
	}

	return jobs
}

// Errors
var (
	ErrQueueFull = NewError("queue is full")
)

// Error represents an error in the queue
type Error struct {
	message string
}

// NewError creates a new error with the specified message
func NewError(message string) *Error {
	return &Error{message: message}
}

// Error returns the error message
func (e *Error) Error() string {
	return e.message
}
