package metrics

import (
	"github.com/fukuchancat/github-pr-metrics/internal/api"
	"github.com/fukuchancat/github-pr-metrics/pkg/utils"
	"github.com/google/go-github/v74/github"
)

// Orchestrates individual PR and aggregated metrics computation
type Calculator struct {
	prCalculator         *PRMetricsCalculator
	aggregatedCalculator *AggregatedMetricsCalculator
	logger               *utils.Logger
}

// Initializes both individual and aggregated metrics calculators
func NewCalculator(client *api.Client, logger *utils.Logger) *Calculator {
	return &Calculator{
		prCalculator:         NewPRMetricsCalculator(client, logger),
		aggregatedCalculator: NewAggregatedMetricsCalculator(logger),
		logger:               logger,
	}
}

// Delegates PR metrics calculation to the PR calculator
func (c *Calculator) CalculatePRMetrics(owner, repo string, pr *github.PullRequest) (*api.PRMetrics, error) {
	return c.prCalculator.CalculatePRMetrics(owner, repo, pr)
}

// Delegates batch PR metrics calculation to the PR calculator
func (c *Calculator) CalculateAllPRMetrics(owner, repo string, prs []*github.PullRequest) ([]*api.PRMetrics, error) {
	return c.prCalculator.CalculateAllPRMetrics(owner, repo, prs)
}

// Delegates weekly metrics aggregation to the aggregated calculator
func (c *Calculator) CalculateWeeklyAggregatedMetrics(prMetrics []*api.PRMetrics) ([]*api.AggregatedMetrics, error) {
	return c.aggregatedCalculator.CalculateWeeklyAggregatedMetrics(prMetrics)
}

// Delegates monthly metrics aggregation to the aggregated calculator
func (c *Calculator) CalculateMonthlyAggregatedMetrics(prMetrics []*api.PRMetrics) ([]*api.AggregatedMetrics, error) {
	return c.aggregatedCalculator.CalculateMonthlyAggregatedMetrics(prMetrics)
}
