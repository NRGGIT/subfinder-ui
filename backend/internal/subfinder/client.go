package subfinder

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/user/subfinder-service/backend/pkg/models"
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
func (c *Client) FindSubdomains(ctx context.Context, domain string, config models.SubfinderConfig) ([]models.SubdomainInfo, []string, error) {
	c.logger.Printf("Finding subdomains for domain %s", domain)

	// Ensure the subfinder binary exists
	if _, err := exec.LookPath("subfinder"); err != nil {
		return nil, nil, fmt.Errorf("subfinder binary not found: %v", err)
	}

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
		if ctx.Err() != nil {
			return nil, nil, fmt.Errorf("subfinder canceled: %v", ctx.Err())
		}
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, nil, fmt.Errorf("subfinder failed: %s", string(exitErr.Stderr))
		}
		return nil, nil, fmt.Errorf("failed to run subfinder: %v, output: %s", err, string(output))
	}

	// Parse the output into structured data
	subdomainInfos := parseSubfinderOutput(string(output), config.IncludeIPs)

	// Apply depth filtering if maxDepth is set
	if config.MaxDepth > 0 {
		c.logger.Printf("Filtering subdomains by depth: %d", config.MaxDepth)
		subdomainInfos = filterSubdomainsByDepth(subdomainInfos, domain, config.MaxDepth)
	}

	// Apply www filtering if excludeWww is set
	if config.ExcludeWww {
		c.logger.Printf("Filtering out www subdomains")
		subdomainInfos = filterWwwSubdomains(subdomainInfos, true)
	}

	// For now, we don't have a way to get the sources used from the CLI output
	// In a real implementation, we would use the subfinder library directly
	sourcesUsed := []string{"all"}

	c.logger.Printf("Found %d subdomains after filtering", len(subdomainInfos))
	return subdomainInfos, sourcesUsed, nil
}

// parseSubfinderOutput parses the output of subfinder into SubdomainInfo structs
func parseSubfinderOutput(output string, includeIPs bool) []models.SubdomainInfo {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	var results []models.SubdomainInfo

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, ",")
		info := models.SubdomainInfo{}

		if len(parts) == 0 {
			continue // Skip empty lines after split
		}

		info.Subdomain = parts[0]

		if includeIPs {
			// Expect format: subdomain,ip[,source]
			if len(parts) > 1 {
				info.IP = parts[1]
			}
			if len(parts) > 2 {
				info.Source = parts[2]
			} else {
				info.Source = "unknown" // Default if source is missing but IP is present
			}
		} else {
			// Expect format: subdomain[,source]
			if len(parts) > 1 {
				info.Source = parts[1]
			} else {
				info.Source = "unknown" // Default if source is missing and no IP expected
			}
		}
		results = append(results, info)
	}

	return results
}

// countDomainLevels counts the number of levels in a domain
// e.g., "example.com" has 2 levels, "sub.example.com" has 3 levels
func countDomainLevels(domain string) int {
	return len(strings.Split(domain, "."))
}

// filterSubdomainsByDepth filters subdomains based on the max depth
// baseDomain is the original domain (e.g., "example.com")
// maxDepth is the maximum depth level (e.g., 1 would include only direct subdomains)
func filterSubdomainsByDepth(subdomains []models.SubdomainInfo, baseDomain string, maxDepth int) []models.SubdomainInfo {
	if maxDepth <= 0 {
		return subdomains // No filtering if maxDepth is not set
	}

	baseLevels := countDomainLevels(baseDomain)
	maxLevels := baseLevels + maxDepth

	var filtered []models.SubdomainInfo
	for _, info := range subdomains {
		if countDomainLevels(info.Subdomain) <= maxLevels {
			filtered = append(filtered, info)
		}
	}

	return filtered
}

// filterWwwSubdomains filters out www subdomains based on the configuration
func filterWwwSubdomains(subdomains []models.SubdomainInfo, excludeWww bool) []models.SubdomainInfo {
	if !excludeWww {
		return subdomains // No filtering if excludeWww is false
	}

	var filtered []models.SubdomainInfo
	for _, info := range subdomains {
		if !strings.HasPrefix(info.Subdomain, "www.") {
			filtered = append(filtered, info)
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
