package metrics

import (
	"sort"
	"time"

	"github.com/fukuchancat/github-pr-metrics/internal/api"
	"github.com/fukuchancat/github-pr-metrics/pkg/utils"
	"github.com/google/go-github/v74/github"
)

// Aggregates GitHub API data to compute comprehensive PR analytics
type PRMetricsCalculator struct {
	client *api.Client
	logger *utils.Logger
}

// Initializes calculator with API client and logger dependencies
func NewPRMetricsCalculator(client *api.Client, logger *utils.Logger) *PRMetricsCalculator {
	return &PRMetricsCalculator{
		client: client,
		logger: logger,
	}
}

// Aggregates commits, comments, reviews, and timing data into comprehensive metrics
func (c *PRMetricsCalculator) CalculatePRMetrics(owner, repo string, pr *github.PullRequest) (*api.PRMetrics, error) {
	c.logger.Debug("Calculating metrics for PR #%d: %s", pr.GetNumber(), pr.GetTitle())

	metrics := api.PRMetrics{
		Number:    pr.GetNumber(),
		Title:     pr.GetTitle(),
		Author:    pr.User.GetLogin(),
		CreatedAt: pr.GetCreatedAt().Time,
		MergedAt:  pr.GetMergedAt().Time,
		State:     pr.GetState(),
	}

	// Get milestone information
	if pr.Milestone != nil {
		metrics.Milestone = pr.Milestone.GetTitle()
	}

	// Get PR details for additions, deletions, and changed files
	additions, deletions, changedFiles, err := c.calculatePRDetails(owner, repo, pr.GetNumber())
	if err != nil {
		return nil, err
	}
	metrics.Additions = additions
	metrics.Deletions = deletions
	metrics.ChangedFiles = changedFiles

	// Get commits and calculate commit-related metrics
	commits, err := c.client.GetPRCommits(owner, repo, pr.GetNumber())
	if err != nil {
		return nil, err
	}
	commitMetrics := c.calculateCommitMetrics(commits, metrics.CreatedAt)
	metrics.CommitCount = commitMetrics.CommitCount
	metrics.FirstCommitAt = commitMetrics.FirstCommitAt
	metrics.LastCommitAt = commitMetrics.LastCommitAt
	metrics.CommitCountDuringPR = commitMetrics.CommitCountDuringPR

	// Get comments and calculate comment-related metrics
	comments, err := c.client.GetPRComments(owner, repo, pr.GetNumber())
	if err != nil {
		c.logger.Warn("Failed to get comments for PR #%d: %v", pr.GetNumber(), err)
		// Continue with empty comments data
	} else {
		commentMetrics := c.calculateCommentMetrics(comments)
		metrics.CommentCount = commentMetrics.CommentCount
		metrics.FirstCommentAt = commentMetrics.FirstCommentAt
	}

	// Calculate review-related metrics
	reviewMetrics, err := c.calculateReviewMetrics(owner, repo, pr.GetNumber())
	if err != nil {
		// Continue with empty reviews data if there's an error
		c.logger.Warn("Failed to get reviews for PR #%d: %v", pr.GetNumber(), err)
	} else {
		metrics.ReviewCount = reviewMetrics.ReviewCount
		metrics.ApprovalCount = reviewMetrics.ApprovalCount

		// Calculate time to first approval
		if !reviewMetrics.FirstApprovalAt.IsZero() {
			metrics.TimeToApprovalHours = reviewMetrics.FirstApprovalAt.Sub(metrics.CreatedAt).Hours()
		}
	}

	// Calculate time-related metrics
	timeMetrics := c.calculateTimeMetrics(
		metrics.CreatedAt,
		metrics.MergedAt,
		metrics.FirstCommitAt,
		metrics.LastCommitAt,
		metrics.FirstCommentAt,
	)
	metrics.FirstCommitToCreateHours = timeMetrics.FirstCommitToCreateHours
	metrics.CreateToLastCommitHours = timeMetrics.CreateToLastCommitHours
	metrics.FirstCommitToMergeHours = timeMetrics.FirstCommitToMergeHours
	metrics.LastCommitToMergeHours = timeMetrics.LastCommitToMergeHours
	metrics.TotalPRLifetimeHours = timeMetrics.TotalPRLifetimeHours
	metrics.CreatedToFirstCommentHours = timeMetrics.CreatedToFirstCommentHours

	// Calculate waiting periods
	if len(commits) > 0 && len(comments) > 0 {
		waitingPeriods := c.calculateWaitingPeriods(commits, comments)
		metrics.MaxNoActivityPeriodHours = waitingPeriods.MaxNoActivityPeriodHours
		metrics.MaxNoCommentPeriodHours = waitingPeriods.MaxNoCommentPeriodHours
		metrics.MaxNoCommitPeriodHours = waitingPeriods.MaxNoCommitPeriodHours
	}

	c.logger.Debug("Calculated metrics for PR #%d: %d commits, %d comments, %d reviews, %d approvals",
		pr.GetNumber(), metrics.CommitCount, metrics.CommentCount, metrics.ReviewCount, metrics.ApprovalCount)

	return &metrics, nil
}

// Fetches additions, deletions, and changed files count from GitHub API
func (c *PRMetricsCalculator) calculatePRDetails(owner, repo string, number int) (int, int, int, error) {
	prDetails, err := c.client.GetPRDetails(owner, repo, number)
	if err != nil {
		return 0, 0, 0, err
	}

	return prDetails.GetAdditions(), prDetails.GetDeletions(), prDetails.GetChangedFiles(), nil
}

// CommitMetricsResult contains timing and frequency data for commits
type CommitMetricsResult struct {
	CommitCount         int
	FirstCommitAt       time.Time
	LastCommitAt        time.Time
	CommitCountDuringPR int
}

// Processes commit timestamps to derive timing and frequency metrics
func (c *PRMetricsCalculator) calculateCommitMetrics(commits []*github.RepositoryCommit, createdAt time.Time) CommitMetricsResult {
	result := CommitMetricsResult{
		CommitCount: len(commits),
	}

	if len(commits) > 0 {
		firstCommit := commits[0]
		lastCommit := commits[len(commits)-1]

		if firstCommit.Commit != nil && firstCommit.Commit.Author != nil && firstCommit.Commit.Author.Date != nil {
			result.FirstCommitAt = firstCommit.Commit.Author.GetDate().Time
		}

		if lastCommit.Commit != nil && lastCommit.Commit.Author != nil && lastCommit.Commit.Author.Date != nil {
			result.LastCommitAt = lastCommit.Commit.Author.GetDate().Time
		}

		// Count commits made during PR (after PR creation)
		commitsDuringPR := 0
		for _, commit := range commits {
			if commit.Commit != nil && commit.Commit.Author != nil && commit.Commit.Author.Date != nil {
				commitTime := commit.Commit.Author.GetDate().Time
				if !commitTime.Before(createdAt) {
					commitsDuringPR++
				}
			}
		}
		result.CommitCountDuringPR = commitsDuringPR
	}

	return result
}

// CommentMetricsResult contains comment count and timing data
type CommentMetricsResult struct {
	CommentCount   int
	FirstCommentAt time.Time
}

// Extracts comment count and first comment timing
func (c *PRMetricsCalculator) calculateCommentMetrics(comments []*github.PullRequestComment) CommentMetricsResult {
	result := CommentMetricsResult{
		CommentCount: len(comments),
	}

	if len(comments) > 0 {
		result.FirstCommentAt = comments[0].GetCreatedAt().Time
	}

	return result
}

// ReviewMetricsResult contains review counts and approval timing data
type ReviewMetricsResult struct {
	ReviewCount     int
	ApprovalCount   int
	FirstApprovalAt time.Time
}

// Processes review states to count approvals and track approval timing
func (c *PRMetricsCalculator) calculateReviewMetrics(owner, repo string, number int) (ReviewMetricsResult, error) {
	result := ReviewMetricsResult{}

	reviews, err := c.client.GetPRReviews(owner, repo, number)
	if err != nil {
		return result, err
	}

	result.ReviewCount = len(reviews)

	// Calculate review-related metrics
	approvalCount := 0
	var firstApprovalAt time.Time

	for _, review := range reviews {
		if review.GetState() == "APPROVED" {
			approvalCount++

			// Record the time of the first approval
			if firstApprovalAt.IsZero() || review.GetSubmittedAt().Before(firstApprovalAt) {
				firstApprovalAt = review.GetSubmittedAt().Time
			}
		}
	}

	result.ApprovalCount = approvalCount
	result.FirstApprovalAt = firstApprovalAt

	return result, nil
}

// TimeMetricsResult contains durations between key PR lifecycle events
type TimeMetricsResult struct {
	FirstCommitToCreateHours   float64
	CreateToLastCommitHours    float64
	FirstCommitToMergeHours    float64
	LastCommitToMergeHours     float64
	TotalPRLifetimeHours       float64
	CreatedToFirstCommentHours float64
}

// Computes duration between key PR lifecycle events
func (c *PRMetricsCalculator) calculateTimeMetrics(createdAt, mergedAt, firstCommitAt, lastCommitAt, firstCommentAt time.Time) TimeMetricsResult {
	result := TimeMetricsResult{}

	// Calculate first commit to PR creation time
	if !firstCommitAt.IsZero() {
		result.FirstCommitToCreateHours = createdAt.Sub(firstCommitAt).Hours()
	}

	// Calculate PR creation to last commit time
	if !lastCommitAt.IsZero() {
		result.CreateToLastCommitHours = lastCommitAt.Sub(createdAt).Hours()
	}

	// Calculate merge-related time metrics
	if !mergedAt.IsZero() {
		if !firstCommitAt.IsZero() {
			result.FirstCommitToMergeHours = mergedAt.Sub(firstCommitAt).Hours()
		}

		if !lastCommitAt.IsZero() {
			result.LastCommitToMergeHours = mergedAt.Sub(lastCommitAt).Hours()
		}

		// Calculate total PR lifetime
		result.TotalPRLifetimeHours = mergedAt.Sub(createdAt).Hours()
	}

	// Calculate time from PR creation to first comment
	if !firstCommentAt.IsZero() {
		result.CreatedToFirstCommentHours = firstCommentAt.Sub(createdAt).Hours()
	}

	return result
}

// WaitingPeriodsResult contains maximum inactivity periods between events
type WaitingPeriodsResult struct {
	MaxNoActivityPeriodHours float64
	MaxNoCommentPeriodHours  float64
	MaxNoCommitPeriodHours   float64
}

// Identifies maximum gaps between commits, comments, and all activities
func (c *PRMetricsCalculator) calculateWaitingPeriods(commits []*github.RepositoryCommit, comments []*github.PullRequestComment) WaitingPeriodsResult {
	result := WaitingPeriodsResult{}

	// Store commit and comment times in a sorted slice
	var allEvents []time.Time

	// Add commit times
	for _, commit := range commits {
		if commit.Commit != nil && commit.Commit.Author != nil && commit.Commit.Author.Date != nil {
			allEvents = append(allEvents, commit.Commit.Author.GetDate().Time)
		}
	}

	// Add comment times
	for _, comment := range comments {
		allEvents = append(allEvents, comment.GetCreatedAt().Time)
	}

	// Sort by time
	sort.Slice(allEvents, func(i, j int) bool {
		return allEvents[i].Before(allEvents[j])
	})

	// Calculate maximum waiting periods
	var maxNoActivityPeriod float64
	var maxNoCommentPeriod float64
	var maxNoCommitPeriod float64

	// Extract commit times only
	var commitTimes []time.Time
	for _, commit := range commits {
		if commit.Commit != nil && commit.Commit.Author != nil && commit.Commit.Author.Date != nil {
			commitTimes = append(commitTimes, commit.Commit.Author.GetDate().Time)
		}
	}
	sort.Slice(commitTimes, func(i, j int) bool {
		return commitTimes[i].Before(commitTimes[j])
	})

	// Extract comment times only
	var commentTimes []time.Time
	for _, comment := range comments {
		commentTimes = append(commentTimes, comment.GetCreatedAt().Time)
	}
	sort.Slice(commentTimes, func(i, j int) bool {
		return commentTimes[i].Before(commentTimes[j])
	})

	// Calculate maximum interval between all activities
	for i := 0; i < len(allEvents)-1; i++ {
		gap := allEvents[i+1].Sub(allEvents[i]).Hours()
		if gap > maxNoActivityPeriod {
			maxNoActivityPeriod = gap
		}
	}

	// Calculate maximum interval between comments
	for i := 0; i < len(commentTimes)-1; i++ {
		gap := commentTimes[i+1].Sub(commentTimes[i]).Hours()
		if gap > maxNoCommentPeriod {
			maxNoCommentPeriod = gap
		}
	}

	// Calculate maximum interval between commits
	for i := 0; i < len(commitTimes)-1; i++ {
		gap := commitTimes[i+1].Sub(commitTimes[i]).Hours()
		if gap > maxNoCommitPeriod {
			maxNoCommitPeriod = gap
		}
	}

	result.MaxNoActivityPeriodHours = maxNoActivityPeriod
	result.MaxNoCommentPeriodHours = maxNoCommentPeriod
	result.MaxNoCommitPeriodHours = maxNoCommitPeriod

	return result
}

// Processes multiple PRs with error handling and progress logging
func (c *PRMetricsCalculator) CalculateAllPRMetrics(owner, repo string, prs []*github.PullRequest) ([]*api.PRMetrics, error) {
	c.logger.Info("Calculating metrics for %d pull requests", len(prs))

	var allMetrics []*api.PRMetrics

	for i, pr := range prs {
		c.logger.Debug("Processing PR #%d (%d/%d)", pr.GetNumber(), i+1, len(prs))

		metrics, err := c.CalculatePRMetrics(owner, repo, pr)
		if err != nil {
			c.logger.Error("Failed to calculate metrics for PR #%d: %v", pr.GetNumber(), err)
			continue
		}

		allMetrics = append(allMetrics, metrics)
	}

	c.logger.Info("Successfully calculated metrics for %d/%d pull requests", len(allMetrics), len(prs))
	return allMetrics, nil
}
