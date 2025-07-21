package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"

	"strings"
	"syscall"
	"time"

	"github.com/user/subfinder-service/backend/internal/api"
	"github.com/user/subfinder-service/backend/internal/queue"
	"github.com/user/subfinder-service/backend/internal/worker"
)

func main() {
	// Set up logger
	logger := log.New(os.Stdout, "[SUBFINDER-SERVICE] ", log.LstdFlags)
	logger.Println("Starting subfinder service...")

	if err := checkSubfinder(logger); err != nil {
		logger.Fatalf("subfinder not available: %v", err)
	}

	// Create job queue
	jobQueue := queue.NewJobQueue()

	// Create worker pool
	workerCount := getEnvInt("WORKER_COUNT", runtime.NumCPU())
	logger.Printf("Using %d worker(s)", workerCount)
	workerPool := worker.NewWorkerPool(workerCount, jobQueue, logger)

	// Start worker pool
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	workerPool.Start(ctx)

	// Create and start API server
	port := getEnv("PORT", "8080")
	server := api.NewServer(port, jobQueue, logger)
	go func() {
		if err := server.Start(); err != nil {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	logger.Println("Shutting down server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	// Wait for worker pool to finish
	cancel()
	workerPool.Wait()

	logger.Println("Server exited properly")
}

// getEnv returns the value of an environment variable or a default value if not set
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvInt returns the integer value of an environment variable or a default value if not set
func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := parseInt(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// parseInt converts a string to an integer
func parseInt(value string) (int, error) {
	var result int
	_, err := fmt.Sscanf(value, "%d", &result)
	return result, err
}

// checkSubfinder verifies that the subfinder binary is available and logs its version.
func checkSubfinder(logger *log.Logger) error {
	path, err := exec.LookPath("subfinder")
	if err != nil {
		return err
	}
	cmd := exec.Command(path, "-version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Printf("failed to execute subfinder -version: %v", err)
	}
	logger.Printf("Using subfinder at %s version: %s", path, strings.TrimSpace(string(out)))
	return nil
}
