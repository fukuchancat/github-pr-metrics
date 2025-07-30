package metrics

import (
	"fmt"
	"sort"
	"time"

	"github.com/fukuchancat/github-pr-metrics/internal/api"
	"github.com/fukuchancat/github-pr-metrics/pkg/utils"
)

// Computes statistical summaries across PR collections by time period
type AggregatedMetricsCalculator struct {
	logger *utils.Logger
}

// Initializes calculator with logger dependency
func NewAggregatedMetricsCalculator(logger *utils.Logger) *AggregatedMetricsCalculator {
	return &AggregatedMetricsCalculator{
		logger: logger,
	}
}

// Groups PRs by ISO week and computes averages and medians
func (c *AggregatedMetricsCalculator) CalculateWeeklyAggregatedMetrics(prMetrics []*api.PRMetrics) ([]*api.AggregatedMetrics, error) {
	c.logger.Info("Calculating weekly aggregated metrics")

	// Group PRs by week
	weeklyPRs := make(map[string][]*api.PRMetrics)
	weeklyStartDates := make(map[string]time.Time)
	weeklyEndDates := make(map[string]time.Time)

	for _, pr := range prMetrics {
		// Skip PRs that haven't been merged
		if pr.MergedAt.IsZero() {
			continue
		}

		// Get the week number (ISO week)
		year, week := pr.MergedAt.ISOWeek()
		weekKey := fmt.Sprintf("%d-W%02d", year, week)

		// Calculate the start and end date of the week
		// ISO week starts on Monday
		startOfWeek := getStartOfISOWeek(pr.MergedAt)
		endOfWeek := startOfWeek.AddDate(0, 0, 6) // End of week (Sunday)

		if _, exists := weeklyPRs[weekKey]; !exists {
			weeklyPRs[weekKey] = []*api.PRMetrics{}
			weeklyStartDates[weekKey] = startOfWeek
			weeklyEndDates[weekKey] = endOfWeek
		}

		weeklyPRs[weekKey] = append(weeklyPRs[weekKey], pr)
	}

	// Calculate aggregated metrics for each week
	var weeklyMetrics []*api.AggregatedMetrics

	for weekKey, prs := range weeklyPRs {
		aggregated := c.calculateAggregatedMetrics(weekKey, weeklyStartDates[weekKey], weeklyEndDates[weekKey], prs)
		weeklyMetrics = append(weeklyMetrics, aggregated)
	}

	// Sort by period
	sort.Slice(weeklyMetrics, func(i, j int) bool {
		return weeklyMetrics[i].Period < weeklyMetrics[j].Period
	})

	c.logger.Info("Successfully calculated weekly aggregated metrics for %d weeks", len(weeklyMetrics))
	return weeklyMetrics, nil
}

// Groups PRs by calendar month and computes statistical summaries
func (c *AggregatedMetricsCalculator) CalculateMonthlyAggregatedMetrics(prMetrics []*api.PRMetrics) ([]*api.AggregatedMetrics, error) {
	c.logger.Info("Calculating monthly aggregated metrics")

	// Group PRs by month
	monthlyPRs := make(map[string][]*api.PRMetrics)
	monthlyStartDates := make(map[string]time.Time)
	monthlyEndDates := make(map[string]time.Time)

	for _, pr := range prMetrics {
		// Skip PRs that haven't been merged
		if pr.MergedAt.IsZero() {
			continue
		}

		// Get the month
		year, month, _ := pr.MergedAt.Date()
		monthKey := fmt.Sprintf("%d-%02d", year, month)

		// Calculate the start and end date of the month
		startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, pr.MergedAt.Location())
		endOfMonth := startOfMonth.AddDate(0, 1, -1) // Last day of month

		if _, exists := monthlyPRs[monthKey]; !exists {
			monthlyPRs[monthKey] = []*api.PRMetrics{}
			monthlyStartDates[monthKey] = startOfMonth
			monthlyEndDates[monthKey] = endOfMonth
		}

		monthlyPRs[monthKey] = append(monthlyPRs[monthKey], pr)
	}

	// Calculate aggregated metrics for each month
	var monthlyMetrics []*api.AggregatedMetrics

	for monthKey, prs := range monthlyPRs {
		aggregated := c.calculateAggregatedMetrics(monthKey, monthlyStartDates[monthKey], monthlyEndDates[monthKey], prs)
		monthlyMetrics = append(monthlyMetrics, aggregated)
	}

	// Sort by period
	sort.Slice(monthlyMetrics, func(i, j int) bool {
		return monthlyMetrics[i].Period < monthlyMetrics[j].Period
	})

	c.logger.Info("Successfully calculated monthly aggregated metrics for %d months", len(monthlyMetrics))
	return monthlyMetrics, nil
}

// Computes averages and medians for all metrics within a PR group
func (c *AggregatedMetricsCalculator) calculateAggregatedMetrics(period string, startDate, endDate time.Time, prs []*api.PRMetrics) *api.AggregatedMetrics {
	prCount := len(prs)
	if prCount == 0 {
		return &api.AggregatedMetrics{
			Period:    period,
			StartDate: startDate,
			EndDate:   endDate,
			PRCount:   0,
		}
	}

	// Initialize sums and slices for median calculation
	var (
		sumCommitCount                int
		sumCommentCount               int
		sumReviewCount                int
		sumApprovalCount              int
		sumAdditions                  int
		sumDeletions                  int
		sumChangedFiles               int
		sumCommitCountDuringPR        int
		sumFirstCommitToCreateHours   float64
		sumCreateToLastCommitHours    float64
		sumFirstCommitToMergeHours    float64
		sumLastCommitToMergeHours     float64
		sumCreatedToFirstCommentHours float64
		sumTimeToApprovalHours        float64
		sumTotalPRLifetimeHours       float64
		sumMaxNoCommentPeriodHours    float64
		sumMaxNoCommitPeriodHours     float64
		sumMaxNoActivityPeriodHours   float64

		countFirstCommitToCreate   int
		countCreateToLastCommit    int
		countFirstCommitToMerge    int
		countLastCommitToMerge     int
		countCreatedToFirstComment int
		countTimeToApproval        int
		countTotalPRLifetime       int
		countMaxNoCommentPeriod    int
		countMaxNoCommitPeriod     int
		countMaxNoActivityPeriod   int

		commitCounts               []int
		commentCounts              []int
		reviewCounts               []int
		approvalCounts             []int
		additions                  []int
		deletions                  []int
		changedFiles               []int
		commitCountsDuringPR       []int
		firstCommitToCreateHours   []float64
		createToLastCommitHours    []float64
		firstCommitToMergeHours    []float64
		lastCommitToMergeHours     []float64
		createdToFirstCommentHours []float64
		timeToApprovalHours        []float64
		totalPRLifetimeHours       []float64
		maxNoCommentPeriodHours    []float64
		maxNoCommitPeriodHours     []float64
		maxNoActivityPeriodHours   []float64
	)

	// Calculate sums and collect values for median calculation
	for _, pr := range prs {
		// Sums for averages
		sumCommitCount += pr.CommitCount
		sumCommentCount += pr.CommentCount
		sumReviewCount += pr.ReviewCount
		sumApprovalCount += pr.ApprovalCount
		sumAdditions += pr.Additions
		sumDeletions += pr.Deletions
		sumChangedFiles += pr.ChangedFiles
		sumCommitCountDuringPR += pr.CommitCountDuringPR

		// Values for median calculation
		commitCounts = append(commitCounts, pr.CommitCount)
		commentCounts = append(commentCounts, pr.CommentCount)
		reviewCounts = append(reviewCounts, pr.ReviewCount)
		approvalCounts = append(approvalCounts, pr.ApprovalCount)
		additions = append(additions, pr.Additions)
		deletions = append(deletions, pr.Deletions)
		changedFiles = append(changedFiles, pr.ChangedFiles)
		commitCountsDuringPR = append(commitCountsDuringPR, pr.CommitCountDuringPR)

		// Time metrics
		if pr.FirstCommitToCreateHours > 0 {
			sumFirstCommitToCreateHours += pr.FirstCommitToCreateHours
			countFirstCommitToCreate++
			firstCommitToCreateHours = append(firstCommitToCreateHours, pr.FirstCommitToCreateHours)
		}

		if pr.CreateToLastCommitHours > 0 {
			sumCreateToLastCommitHours += pr.CreateToLastCommitHours
			countCreateToLastCommit++
			createToLastCommitHours = append(createToLastCommitHours, pr.CreateToLastCommitHours)
		}

		if pr.FirstCommitToMergeHours > 0 {
			sumFirstCommitToMergeHours += pr.FirstCommitToMergeHours
			countFirstCommitToMerge++
			firstCommitToMergeHours = append(firstCommitToMergeHours, pr.FirstCommitToMergeHours)
		}

		if pr.LastCommitToMergeHours > 0 {
			sumLastCommitToMergeHours += pr.LastCommitToMergeHours
			countLastCommitToMerge++
			lastCommitToMergeHours = append(lastCommitToMergeHours, pr.LastCommitToMergeHours)
		}

		if pr.CreatedToFirstCommentHours > 0 {
			sumCreatedToFirstCommentHours += pr.CreatedToFirstCommentHours
			countCreatedToFirstComment++
			createdToFirstCommentHours = append(createdToFirstCommentHours, pr.CreatedToFirstCommentHours)
		}

		if pr.TimeToApprovalHours > 0 {
			sumTimeToApprovalHours += pr.TimeToApprovalHours
			countTimeToApproval++
			timeToApprovalHours = append(timeToApprovalHours, pr.TimeToApprovalHours)
		}

		if pr.TotalPRLifetimeHours > 0 {
			sumTotalPRLifetimeHours += pr.TotalPRLifetimeHours
			countTotalPRLifetime++
			totalPRLifetimeHours = append(totalPRLifetimeHours, pr.TotalPRLifetimeHours)
		}

		if pr.MaxNoCommentPeriodHours > 0 {
			sumMaxNoCommentPeriodHours += pr.MaxNoCommentPeriodHours
			countMaxNoCommentPeriod++
			maxNoCommentPeriodHours = append(maxNoCommentPeriodHours, pr.MaxNoCommentPeriodHours)
		}

		if pr.MaxNoCommitPeriodHours > 0 {
			sumMaxNoCommitPeriodHours += pr.MaxNoCommitPeriodHours
			countMaxNoCommitPeriod++
			maxNoCommitPeriodHours = append(maxNoCommitPeriodHours, pr.MaxNoCommitPeriodHours)
		}

		if pr.MaxNoActivityPeriodHours > 0 {
			sumMaxNoActivityPeriodHours += pr.MaxNoActivityPeriodHours
			countMaxNoActivityPeriod++
			maxNoActivityPeriodHours = append(maxNoActivityPeriodHours, pr.MaxNoActivityPeriodHours)
		}
	}

	// Calculate averages and medians
	metrics := &api.AggregatedMetrics{
		Period:                 period,
		StartDate:              startDate,
		EndDate:                endDate,
		PRCount:                prCount,
		AvgCommitCount:         float64(sumCommitCount) / float64(prCount),
		AvgCommentCount:        float64(sumCommentCount) / float64(prCount),
		AvgReviewCount:         float64(sumReviewCount) / float64(prCount),
		AvgApprovalCount:       float64(sumApprovalCount) / float64(prCount),
		AvgAdditions:           float64(sumAdditions) / float64(prCount),
		AvgDeletions:           float64(sumDeletions) / float64(prCount),
		AvgChangedFiles:        float64(sumChangedFiles) / float64(prCount),
		AvgCommitCountDuringPR: float64(sumCommitCountDuringPR) / float64(prCount),

		// Calculate medians for count metrics
		MedianCommitCount:         calculateMedianInt(commitCounts),
		MedianCommentCount:        calculateMedianInt(commentCounts),
		MedianReviewCount:         calculateMedianInt(reviewCounts),
		MedianApprovalCount:       calculateMedianInt(approvalCounts),
		MedianAdditions:           calculateMedianInt(additions),
		MedianDeletions:           calculateMedianInt(deletions),
		MedianChangedFiles:        calculateMedianInt(changedFiles),
		MedianCommitCountDuringPR: calculateMedianInt(commitCountsDuringPR),
	}

	// Calculate averages for time metrics (only if we have valid data)
	if countFirstCommitToCreate > 0 {
		metrics.AvgFirstCommitToCreateHours = sumFirstCommitToCreateHours / float64(countFirstCommitToCreate)
		metrics.MedianFirstCommitToCreateHours = calculateMedianFloat(firstCommitToCreateHours)
	}

	if countCreateToLastCommit > 0 {
		metrics.AvgCreateToLastCommitHours = sumCreateToLastCommitHours / float64(countCreateToLastCommit)
		metrics.MedianCreateToLastCommitHours = calculateMedianFloat(createToLastCommitHours)
	}

	if countFirstCommitToMerge > 0 {
		metrics.AvgFirstCommitToMergeHours = sumFirstCommitToMergeHours / float64(countFirstCommitToMerge)
		metrics.MedianFirstCommitToMergeHours = calculateMedianFloat(firstCommitToMergeHours)
	}

	if countLastCommitToMerge > 0 {
		metrics.AvgLastCommitToMergeHours = sumLastCommitToMergeHours / float64(countLastCommitToMerge)
		metrics.MedianLastCommitToMergeHours = calculateMedianFloat(lastCommitToMergeHours)
	}

	if countCreatedToFirstComment > 0 {
		metrics.AvgCreatedToFirstCommentHours = sumCreatedToFirstCommentHours / float64(countCreatedToFirstComment)
		metrics.MedianCreatedToFirstCommentHours = calculateMedianFloat(createdToFirstCommentHours)
	}

	if countTimeToApproval > 0 {
		metrics.AvgTimeToApprovalHours = sumTimeToApprovalHours / float64(countTimeToApproval)
		metrics.MedianTimeToApprovalHours = calculateMedianFloat(timeToApprovalHours)
	}

	if countTotalPRLifetime > 0 {
		metrics.AvgTotalPRLifetimeHours = sumTotalPRLifetimeHours / float64(countTotalPRLifetime)
		metrics.MedianTotalPRLifetimeHours = calculateMedianFloat(totalPRLifetimeHours)
	}

	if countMaxNoCommentPeriod > 0 {
		metrics.AvgMaxNoCommentPeriodHours = sumMaxNoCommentPeriodHours / float64(countMaxNoCommentPeriod)
		metrics.MedianMaxNoCommentPeriodHours = calculateMedianFloat(maxNoCommentPeriodHours)
	}

	if countMaxNoCommitPeriod > 0 {
		metrics.AvgMaxNoCommitPeriodHours = sumMaxNoCommitPeriodHours / float64(countMaxNoCommitPeriod)
		metrics.MedianMaxNoCommitPeriodHours = calculateMedianFloat(maxNoCommitPeriodHours)
	}

	if countMaxNoActivityPeriod > 0 {
		metrics.AvgMaxNoActivityPeriodHours = sumMaxNoActivityPeriodHours / float64(countMaxNoActivityPeriod)
		metrics.MedianMaxNoActivityPeriodHours = calculateMedianFloat(maxNoActivityPeriodHours)
	}

	return metrics
}
