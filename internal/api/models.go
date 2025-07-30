package api

import (
	"time"
)

// Contains comprehensive analytics data for a single pull request
type PRMetrics struct {
	Number                     int
	Title                      string
	Author                     string
	Milestone                  string
	CreatedAt                  time.Time
	MergedAt                   time.Time
	State                      string
	CommitCount                int
	FirstCommitAt              time.Time
	LastCommitAt               time.Time
	FirstCommitToCreateHours   float64
	CreateToLastCommitHours    float64
	CommitCountDuringPR        int
	FirstCommitToMergeHours    float64
	LastCommitToMergeHours     float64
	CommentCount               int
	FirstCommentAt             time.Time
	CreatedToFirstCommentHours float64
	ReviewCount                int
	Additions                  int
	Deletions                  int
	ChangedFiles               int
	ApprovalCount              int
	TimeToApprovalHours        float64
	TotalPRLifetimeHours       float64
	MaxNoCommentPeriodHours    float64
	MaxNoCommitPeriodHours     float64
	MaxNoActivityPeriodHours   float64
}

// Contains statistical summaries of PR metrics over a time period
type AggregatedMetrics struct {
	Period                           string // YYYY-WW for week, YYYY-MM for month
	StartDate                        time.Time
	EndDate                          time.Time
	PRCount                          int
	AvgCommitCount                   float64
	AvgCommentCount                  float64
	AvgReviewCount                   float64
	AvgApprovalCount                 float64
	AvgAdditions                     float64
	AvgDeletions                     float64
	AvgChangedFiles                  float64
	AvgFirstCommitToCreateHours      float64
	AvgCreateToLastCommitHours       float64
	AvgCommitCountDuringPR           float64
	AvgFirstCommitToMergeHours       float64
	AvgLastCommitToMergeHours        float64
	AvgCreatedToFirstCommentHours    float64
	AvgTimeToApprovalHours           float64
	AvgTotalPRLifetimeHours          float64
	AvgMaxNoCommentPeriodHours       float64
	AvgMaxNoCommitPeriodHours        float64
	AvgMaxNoActivityPeriodHours      float64
	MedianCommitCount                float64
	MedianCommentCount               float64
	MedianReviewCount                float64
	MedianApprovalCount              float64
	MedianAdditions                  float64
	MedianDeletions                  float64
	MedianChangedFiles               float64
	MedianFirstCommitToCreateHours   float64
	MedianCreateToLastCommitHours    float64
	MedianCommitCountDuringPR        float64
	MedianFirstCommitToMergeHours    float64
	MedianLastCommitToMergeHours     float64
	MedianCreatedToFirstCommentHours float64
	MedianTimeToApprovalHours        float64
	MedianTotalPRLifetimeHours       float64
	MedianMaxNoCommentPeriodHours    float64
	MedianMaxNoCommitPeriodHours     float64
	MedianMaxNoActivityPeriodHours   float64
}
