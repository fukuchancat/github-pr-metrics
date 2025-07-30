package main

import (
	"flag"
	"os"
	"strings"
	"time"

	"github.com/fukuchancat/github-pr-metrics/internal/api"
	"github.com/fukuchancat/github-pr-metrics/internal/metrics"
	"github.com/fukuchancat/github-pr-metrics/internal/output"
	"github.com/fukuchancat/github-pr-metrics/pkg/utils"
)

func main() {
	// Parse command line arguments
	githubURL := flag.String("url", "https://api.github.com", "GitHub API URL")
	token := flag.String("token", "", "GitHub Personal Access Token")
	repo := flag.String("repo", "", "Repository name in format 'owner/repo'")
	startDate := flag.String("start-date", "", "Start date for PR filtering (format: YYYY-MM-DD)")
	endDate := flag.String("end-date", "", "End date for PR filtering (format: YYYY-MM-DD)")
	outputDir := flag.String("output-dir", "output", "Output directory for CSV files")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	help := flag.Bool("help", false, "Show help message")

	// Define short options
	flag.StringVar(githubURL, "u", "https://api.github.com", "GitHub API URL (shorthand)")
	flag.StringVar(token, "t", "", "GitHub Personal Access Token (shorthand)")
	flag.StringVar(repo, "r", "", "Repository name in format 'owner/repo' (shorthand)")
	flag.StringVar(startDate, "s", "", "Start date for PR filtering (shorthand)")
	flag.StringVar(endDate, "e", "", "End date for PR filtering (shorthand)")
	flag.StringVar(outputDir, "o", "output", "Output directory for CSV files (shorthand)")
	flag.BoolVar(verbose, "v", false, "Enable verbose logging (shorthand)")
	flag.BoolVar(help, "h", false, "Show help message (shorthand)")

	flag.Parse()

	// Create logger
	logger := utils.NewLogger(*verbose)

	// Show help message if requested
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	// Validate required arguments
	if *token == "" {
		logger.Fatal("GitHub Personal Access Token is required")
	}

	if *repo == "" {
		logger.Fatal("Repository name is required")
	}

	// Parse repository owner and name
	parts := strings.Split(*repo, "/")
	if len(parts) != 2 {
		logger.Fatal("Repository name must be in format 'owner/repo'")
	}
	owner, repoName := parts[0], parts[1]

	// Parse dates
	var start, end time.Time
	var err error

	if *startDate != "" {
		start, err = time.Parse("2006-01-02", *startDate)
		if err != nil {
			logger.Fatal("Invalid start date format: %v", err)
		}
	} else {
		// Default to 7 days ago
		start = time.Now().AddDate(0, 0, -7)
	}

	if *endDate != "" {
		end, err = time.Parse("2006-01-02", *endDate)
		if err != nil {
			logger.Fatal("Invalid end date format: %v", err)
		}
	} else {
		// Default to today
		end = time.Now()
	}

	logger.Info("Fetching PR metrics for %s/%s from %s to %s", owner, repoName, start.Format("2006-01-02"), end.Format("2006-01-02"))

	// Create GitHub API client
	client, err := api.NewClient(*githubURL, *token, logger)
	if err != nil {
		logger.Fatal("Failed to create GitHub API client: %v", err)
	}

	// Get pull requests
	logger.Debug("Fetching pull requests...")
	prs, err := client.GetPullRequests(owner, repoName, start, end)
	if err != nil {
		logger.Fatal("Failed to fetch pull requests: %v", err)
	}

	logger.Info("Found %d pull requests", len(prs))

	// Calculate metrics for each pull request
	calculator := metrics.NewCalculator(client, logger)
	prMetrics, err := calculator.CalculateAllPRMetrics(owner, repoName, prs)
	if err != nil {
		logger.Fatal("Failed to calculate PR metrics: %v", err)
	}

	// Calculate weekly and monthly aggregated metrics
	logger.Debug("Calculating weekly aggregated metrics...")
	weeklyMetrics, err := calculator.CalculateWeeklyAggregatedMetrics(prMetrics)
	if err != nil {
		logger.Fatal("Failed to calculate weekly metrics: %v", err)
	}
	logger.Info("Calculated metrics for %d weeks", len(weeklyMetrics))

	logger.Debug("Calculating monthly aggregated metrics...")
	monthlyMetrics, err := calculator.CalculateMonthlyAggregatedMetrics(prMetrics)
	if err != nil {
		logger.Fatal("Failed to calculate monthly metrics: %v", err)
	}
	logger.Info("Calculated metrics for %d months", len(monthlyMetrics))

	// Write metrics to CSV files in the output directory
	csvWriter := output.NewCSVWriter(logger)
	err = csvWriter.WriteToDirectory(*outputDir, prMetrics, weeklyMetrics, monthlyMetrics)
	if err != nil {
		logger.Fatal("Failed to write CSV files: %v", err)
	}

	logger.Info("Successfully wrote metrics for %d pull requests to directory: %s", len(prMetrics), *outputDir)
}
