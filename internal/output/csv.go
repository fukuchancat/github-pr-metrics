package output

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/fukuchancat/github-pr-metrics/internal/api"
	"github.com/fukuchancat/github-pr-metrics/pkg/utils"
)

// Handles exporting PR metrics data to CSV format files
type CSVWriter struct {
	logger *utils.Logger
}

// Initializes CSV writer with logger dependency
func NewCSVWriter(logger *utils.Logger) *CSVWriter {
	return &CSVWriter{
		logger: logger,
	}
}

// Exports PR, weekly, and monthly metrics to separate CSV files in target directory
func (w *CSVWriter) WriteToDirectory(dirPath string, prMetrics []*api.PRMetrics, weeklyMetrics, monthlyMetrics []*api.AggregatedMetrics) error {
	w.logger.Info("Writing metrics to directory: %s", dirPath)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Write PR metrics
	prFilePath := filepath.Join(dirPath, "pr_metrics.csv")
	if err := w.writePRMetricsCSV(prFilePath, prMetrics); err != nil {
		return fmt.Errorf("failed to write PR metrics: %v", err)
	}

	// Write weekly metrics
	weeklyFilePath := filepath.Join(dirPath, "weekly_metrics.csv")
	if err := w.writeAggregatedMetricsCSV(weeklyFilePath, weeklyMetrics, "Weekly"); err != nil {
		return fmt.Errorf("failed to write weekly metrics: %v", err)
	}

	// Write monthly metrics
	monthlyFilePath := filepath.Join(dirPath, "monthly_metrics.csv")
	if err := w.writeAggregatedMetricsCSV(monthlyFilePath, monthlyMetrics, "Monthly"); err != nil {
		return fmt.Errorf("failed to write monthly metrics: %v", err)
	}

	w.logger.Info("Successfully wrote metrics to directory: %s", dirPath)
	return nil
}

// Legacy method for exporting only PR metrics to a single CSV file
func (w *CSVWriter) WriteCSV(filename string, prMetrics []*api.PRMetrics) error {
	return w.writePRMetricsCSV(filename, prMetrics)
}

// Formats and exports individual PR metrics data to CSV format
func (w *CSVWriter) writePRMetricsCSV(filename string, prMetrics []*api.PRMetrics) error {
	w.logger.Info("Writing %d PR metrics to CSV file: %s", len(prMetrics), filename)

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			w.logger.Warn("Failed to close file: %v", err)
		}
	}()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"PR Number",
		"Title",
		"Author",
		"Milestone",
		"Created At",
		"Merged At",
		"State",
		"Commit Count",
		"First Commit At",
		"Last Commit At",
		"First Commit to Create (Hours)",
		"Create to Last Commit (Hours)",
		"Commit Count During PR",
		"First Commit to Merge (Hours)",
		"Last Commit to Merge (Hours)",
		"Comment Count",
		"First Comment At",
		"Created to First Comment (Hours)",
		"Review Count",
		"Approval Count",
		"Time to Approval (Hours)",
		"Total PR Lifetime (Hours)",
		"Max No Comment Period (Hours)",
		"Max No Commit Period (Hours)",
		"Max No Activity Period (Hours)",
		"Additions",
		"Deletions",
		"Changed Files",
	}

	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data
	for _, pr := range prMetrics {
		row := []string{
			strconv.Itoa(pr.Number),
			pr.Title,
			pr.Author,
			pr.Milestone,
			formatTime(pr.CreatedAt),
			formatTime(pr.MergedAt),
			pr.State,
			strconv.Itoa(pr.CommitCount),
			formatTime(pr.FirstCommitAt),
			formatTime(pr.LastCommitAt),
			formatFloat(pr.FirstCommitToCreateHours),
			formatFloat(pr.CreateToLastCommitHours),
			strconv.Itoa(pr.CommitCountDuringPR),
			formatFloat(pr.FirstCommitToMergeHours),
			formatFloat(pr.LastCommitToMergeHours),
			strconv.Itoa(pr.CommentCount),
			formatTime(pr.FirstCommentAt),
			formatFloat(pr.CreatedToFirstCommentHours),
			strconv.Itoa(pr.ReviewCount),
			strconv.Itoa(pr.ApprovalCount),
			formatFloat(pr.TimeToApprovalHours),
			formatFloat(pr.TotalPRLifetimeHours),
			formatFloat(pr.MaxNoCommentPeriodHours),
			formatFloat(pr.MaxNoCommitPeriodHours),
			formatFloat(pr.MaxNoActivityPeriodHours),
			strconv.Itoa(pr.Additions),
			strconv.Itoa(pr.Deletions),
			strconv.Itoa(pr.ChangedFiles),
		}

		if err := writer.Write(row); err != nil {
			return err
		}
	}

	w.logger.Info("Successfully wrote %d PR metrics to CSV file", len(prMetrics))
	return nil
}

// Formats and exports statistical metrics summaries to CSV format
func (w *CSVWriter) writeAggregatedMetricsCSV(filename string, metrics []*api.AggregatedMetrics, metricsType string) error {
	w.logger.Info("Writing %d %s metrics to CSV file: %s", len(metrics), metricsType, filename)

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			w.logger.Warn("Failed to close file: %v", err)
		}
	}()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{
		"Period",
		"Start Date",
		"End Date",
		"PR Count",
		"Avg Commit Count",
		"Median Commit Count",
		"Avg Comment Count",
		"Median Comment Count",
		"Avg Review Count",
		"Median Review Count",
		"Avg Approval Count",
		"Median Approval Count",
		"Avg Additions",
		"Median Additions",
		"Avg Deletions",
		"Median Deletions",
		"Avg Changed Files",
		"Median Changed Files",
		"Avg First Commit to Create (Hours)",
		"Median First Commit to Create (Hours)",
		"Avg Create to Last Commit (Hours)",
		"Median Create to Last Commit (Hours)",
		"Avg Commit Count During PR",
		"Median Commit Count During PR",
		"Avg First Commit to Merge (Hours)",
		"Median First Commit to Merge (Hours)",
		"Avg Last Commit to Merge (Hours)",
		"Median Last Commit to Merge (Hours)",
		"Avg Created to First Comment (Hours)",
		"Median Created to First Comment (Hours)",
		"Avg Time to Approval (Hours)",
		"Median Time to Approval (Hours)",
		"Avg Total PR Lifetime (Hours)",
		"Median Total PR Lifetime (Hours)",
		"Avg Max No Comment Period (Hours)",
		"Median Max No Comment Period (Hours)",
		"Avg Max No Commit Period (Hours)",
		"Median Max No Commit Period (Hours)",
		"Avg Max No Activity Period (Hours)",
		"Median Max No Activity Period (Hours)",
	}

	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data
	for _, m := range metrics {
		row := []string{
			m.Period,
			formatTime(m.StartDate),
			formatTime(m.EndDate),
			strconv.Itoa(m.PRCount),
			formatFloat(m.AvgCommitCount),
			formatFloat(m.MedianCommitCount),
			formatFloat(m.AvgCommentCount),
			formatFloat(m.MedianCommentCount),
			formatFloat(m.AvgReviewCount),
			formatFloat(m.MedianReviewCount),
			formatFloat(m.AvgApprovalCount),
			formatFloat(m.MedianApprovalCount),
			formatFloat(m.AvgAdditions),
			formatFloat(m.MedianAdditions),
			formatFloat(m.AvgDeletions),
			formatFloat(m.MedianDeletions),
			formatFloat(m.AvgChangedFiles),
			formatFloat(m.MedianChangedFiles),
			formatFloat(m.AvgFirstCommitToCreateHours),
			formatFloat(m.MedianFirstCommitToCreateHours),
			formatFloat(m.AvgCreateToLastCommitHours),
			formatFloat(m.MedianCreateToLastCommitHours),
			formatFloat(m.AvgCommitCountDuringPR),
			formatFloat(m.MedianCommitCountDuringPR),
			formatFloat(m.AvgFirstCommitToMergeHours),
			formatFloat(m.MedianFirstCommitToMergeHours),
			formatFloat(m.AvgLastCommitToMergeHours),
			formatFloat(m.MedianLastCommitToMergeHours),
			formatFloat(m.AvgCreatedToFirstCommentHours),
			formatFloat(m.MedianCreatedToFirstCommentHours),
			formatFloat(m.AvgTimeToApprovalHours),
			formatFloat(m.MedianTimeToApprovalHours),
			formatFloat(m.AvgTotalPRLifetimeHours),
			formatFloat(m.MedianTotalPRLifetimeHours),
			formatFloat(m.AvgMaxNoCommentPeriodHours),
			formatFloat(m.MedianMaxNoCommentPeriodHours),
			formatFloat(m.AvgMaxNoCommitPeriodHours),
			formatFloat(m.MedianMaxNoCommitPeriodHours),
			formatFloat(m.AvgMaxNoActivityPeriodHours),
			formatFloat(m.MedianMaxNoActivityPeriodHours),
		}

		if err := writer.Write(row); err != nil {
			return err
		}
	}

	w.logger.Info("Successfully wrote %d %s metrics to CSV file", len(metrics), metricsType)
	return nil
}

// Converts time to RFC3339 format or empty string if zero
func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}

// Formats floating point values with 2 decimal places
func formatFloat(f float64) string {
	if f == 0 {
		return "0.00"
	}
	return strconv.FormatFloat(f, 'f', 2, 64)
}
