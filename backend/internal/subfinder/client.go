package subfinder

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/user/subfinder-service/pkg/models"
)

// Client represents a client for the subfinder library
type Client struct {
	logger *log.Logger
}

// NewClient creates a new subfinder client
func NewClient(logger *log.Logger) *Client {
	return &Client{
		logger: logger,
	}
}

// FindSubdomains finds subdomains for the specified domain using subfinder
func (c *Client) FindSubdomains(ctx context.Context, domain string, config models.SubfinderConfig) ([]string, []string, error) {
	c.logger.Printf("Finding subdomains for domain %s", domain)

	// Build the command
	args := []string{"-d", domain}

	// Add configuration options
	// Note: subfinder doesn't have a -max-depth flag, so we'll ignore this parameter
	// if config.MaxDepth > 0 {
	// 	args = append(args, "-max-depth", fmt.Sprintf("%d", config.MaxDepth))
	// }


	// Note: -oI flag must be used with -active flag
	if config.IncludeIPs {
		args = append(args, "-oI", "-active")
		// Since we're adding -active, set ExcludeUnresolvable to true
		config.ExcludeUnresolvable = true
	}

	if len(config.Sources) > 0 {
		args = append(args, "-sources", strings.Join(config.Sources, ","))
	}

	if config.RateLimit > 0 {
		args = append(args, "-rate-limit", fmt.Sprintf("%d", config.RateLimit))
	}

	if config.IncludeWildcards {
		args = append(args, "-all")
	}

	if config.ExcludeUnresolvable {
		args = append(args, "-active")
	}

	// Set timeout if specified
	if config.Timeout > 0 {
		args = append(args, "-timeout", fmt.Sprintf("%d", config.Timeout))
	}

	// Add silent mode to get clean output
	args = append(args, "-silent")

	// Create the command
	cmd := exec.CommandContext(ctx, "subfinder", args...)

	// Log the command being executed
	c.logger.Printf("Executing command: subfinder %s", strings.Join(args, " "))
	
	// Run the command and get the output
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.logger.Printf("Command failed with error: %v, output: %s", err, string(output))
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, nil, fmt.Errorf("subfinder failed: %s", string(exitErr.Stderr))
		}
		return nil, nil, fmt.Errorf("failed to run subfinder: %v, output: %s", err, string(output))
	}

	// Parse the output
	subdomains := parseSubfinderOutput(string(output))

	// Apply depth filtering if maxDepth is set
	if config.MaxDepth > 0 {
		c.logger.Printf("Filtering subdomains by depth: %d", config.MaxDepth)
		subdomains = filterSubdomainsByDepth(subdomains, domain, config.MaxDepth)
	}

	// Apply www filtering if excludeWww is set
	if config.ExcludeWww {
		c.logger.Printf("Filtering out www subdomains")
		subdomains = filterWwwSubdomains(subdomains, true)
	}

	// For now, we don't have a way to get the sources used from the CLI output
	// In a real implementation, we would use the subfinder library directly
	sourcesUsed := []string{"all"}

	c.logger.Printf("Found %d subdomains after filtering", len(subdomains))
	return subdomains, sourcesUsed, nil
}

// parseSubfinderOutput parses the output of subfinder
func parseSubfinderOutput(output string) []string {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	var subdomains []string

	for _, line := range lines {
		if line != "" {
			subdomains = append(subdomains, line)
		}
	}

	return subdomains
}

// countDomainLevels counts the number of levels in a domain
// e.g., "example.com" has 2 levels, "sub.example.com" has 3 levels
func countDomainLevels(domain string) int {
	return len(strings.Split(domain, "."))
}

// filterSubdomainsByDepth filters subdomains based on the max depth
// baseDomain is the original domain (e.g., "example.com")
// maxDepth is the maximum depth level (e.g., 1 would include only direct subdomains)
func filterSubdomainsByDepth(subdomains []string, baseDomain string, maxDepth int) []string {
	if maxDepth <= 0 {
		return subdomains // No filtering if maxDepth is not set
	}
	
	baseLevels := countDomainLevels(baseDomain)
	maxLevels := baseLevels + maxDepth
	
	var filtered []string
	for _, subdomain := range subdomains {
		if countDomainLevels(subdomain) <= maxLevels {
			filtered = append(filtered, subdomain)
		}
	}
	
	return filtered
}

// filterWwwSubdomains filters out www subdomains based on the configuration
func filterWwwSubdomains(subdomains []string, excludeWww bool) []string {
	if !excludeWww {
		return subdomains // No filtering if excludeWww is false
	}
	
	var filtered []string
	for _, subdomain := range subdomains {
		if !strings.HasPrefix(subdomain, "www.") {
			filtered = append(filtered, subdomain)
		}
	}
	
	return filtered
}

// Note: In a production environment, we would use the subfinder library directly
// instead of executing the CLI. This implementation is a simplified version that
// uses the CLI for demonstration purposes.
//
// To use the library directly, we would need to import:
// "github.com/projectdiscovery/subfinder/v2/pkg/runner"
// And then use the runner.New() and runner.RunEnumeration() functions.
