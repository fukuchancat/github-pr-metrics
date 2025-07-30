# GitHub PR Metrics

A Go tool that collects various metrics about pull requests in GitHub repositories and outputs them in CSV format.

## Features

- Collects various metrics about pull requests in GitHub repositories
- Detailed metrics for each PR (commit count, comment count, review count, approval count, etc.)
- Tracks the entire PR lifecycle (from first commit to creation, review, and merge)
- Automatically generates weekly and monthly aggregated metrics
- Outputs results in CSV format for data analysis

## How to use

### Generating a Personal Access Token

Before running the tool, you need to generate a [Fine-grained Personal Access Token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens) with the following permissions:

- **Metadata**: Read-only
- **Pull requests**: Read-only

### Running the Tool

```bash
go install github.com/fukuchancat/github-pr-metrics/cmd/github-pr-metrics@latest
github-pr-metrics --url https://api.github.com --token YOUR_PERSONAL_ACCESS_TOKEN --repo owner/repo --start-date 2022-01-01 --end-date 2022-12-31 --output-dir output --verbose
```

## Example Output

This tool outputs three types of CSV files:

### PR Metrics (pr_metrics.csv)

```csv
PR Number,Title,Author,Milestone,Created At,Merged At,State,Commit Count,First Commit At,Last Commit At,First Commit to Create (Hours),Create to Last Commit (Hours),Commit Count During PR,First Commit to Merge (Hours),Last Commit to Merge (Hours),Comment Count,First Comment At,Created to First Comment (Hours),Review Count,Approval Count,Time to Approval (Hours),Total PR Lifetime (Hours),Max No Comment Period (Hours),Max No Commit Period (Hours),Max No Activity Period (Hours),Additions,Deletions,Changed Files
123,Add user authentication feature,alice,Sprint 42,2023-01-15T10:30:00Z,2023-01-18T15:45:00Z,closed,5,2023-01-14T08:20:00Z,2023-01-17T16:30:00Z,26.17,30.00,4,79.42,23.25,8,2023-01-15T14:20:00Z,3.83,3,2,28.50,77.25,12.33,18.75,12.33,245,37,8
124,Fix navigation bar responsiveness,bob,Sprint 42,2023-01-16T09:15:00Z,2023-01-17T11:30:00Z,closed,2,2023-01-16T08:45:00Z,2023-01-16T14:20:00Z,0.50,5.08,2,26.75,21.17,3,2023-01-16T10:05:00Z,0.83,2,1,14.25,26.25,8.50,5.08,8.50,56,12,3
125,Update documentation for API v2,carol,Sprint 42,2023-01-17T13:45:00Z,,open,1,2023-01-17T13:30:00Z,2023-01-17T13:30:00Z,0.25,0.00,0,0.00,0.00,0,,0.00,0,0,0.00,0.00,0.00,0.00,0.00,128,35,4
```

### Weekly Aggregated Metrics (weekly_metrics.csv)

```csv
Period,Start Date,End Date,PR Count,Avg Commit Count,Median Commit Count,Avg Comment Count,Median Comment Count,Avg Review Count,Median Review Count,Avg Approval Count,Median Approval Count,Avg Additions,Median Additions,Avg Deletions,Median Deletions,Avg Changed Files,Median Changed Files,Avg First Commit to Create (Hours),Median First Commit to Create (Hours),Avg Create to Last Commit (Hours),Median Create to Last Commit (Hours),Avg Commit Count During PR,Median Commit Count During PR,Avg First Commit to Merge (Hours),Median First Commit to Merge (Hours),Avg Last Commit to Merge (Hours),Median Last Commit to Merge (Hours),Avg Created to First Comment (Hours),Median Created to First Comment (Hours),Avg Time to Approval (Hours),Median Time to Approval (Hours),Avg Total PR Lifetime (Hours),Median Total PR Lifetime (Hours),Avg Max No Comment Period (Hours),Median Max No Comment Period (Hours),Avg Max No Commit Period (Hours),Median Max No Commit Period (Hours),Avg Max No Activity Period (Hours),Median Max No Activity Period (Hours)
2025-W30,2025-07-21T00:00:00Z,2025-07-27T00:00:00Z,8,4.62,4.50,11.12,9.00,15.88,14.50,1.62,2.00,189.75,40.50,36.00,10.50,6.88,5.00,8.40,0.21,79.71,31.38,3.50,3.00,123.20,124.10,56.21,18.18,13.74,0.02,86.99,31.94,115.99,111.97,49.03,13.89,76.12,28.86,54.02,20.92
2025-W31,2025-07-28T00:00:00Z,2025-08-03T00:00:00Z,12,14.83,5.00,37.25,16.50,42.17,21.50,1.58,1.00,1036.83,141.00,65.08,21.00,21.17,5.50,90.60,0.89,264.47,59.99,11.33,4.00,374.23,107.44,63.24,18.57,122.82,0.02,247.31,92.50,283.62,102.61,75.36,59.04,85.69,76.49,88.62,70.92
```

### Monthly Aggregated Metrics (monthly_metrics.csv)

```csv
Period,Start Date,End Date,PR Count,Avg Commit Count,Median Commit Count,Avg Comment Count,Median Comment Count,Avg Review Count,Median Review Count,Avg Approval Count,Median Approval Count,Avg Additions,Median Additions,Avg Deletions,Median Deletions,Avg Changed Files,Median Changed Files,Avg First Commit to Create (Hours),Median First Commit to Create (Hours),Avg Create to Last Commit (Hours),Median Create to Last Commit (Hours),Avg Commit Count During PR,Median Commit Count During PR,Avg First Commit to Merge (Hours),Median First Commit to Merge (Hours),Avg Last Commit to Merge (Hours),Median Last Commit to Merge (Hours),Avg Created to First Comment (Hours),Median Created to First Comment (Hours),Avg Time to Approval (Hours),Median Time to Approval (Hours),Avg Total PR Lifetime (Hours),Median Total PR Lifetime (Hours),Avg Max No Comment Period (Hours),Median Max No Comment Period (Hours),Avg Max No Commit Period (Hours),Median Max No Commit Period (Hours),Avg Max No Activity Period (Hours),Median Max No Activity Period (Hours)
2025-07,2025-07-01T00:00:00Z,2025-07-31T00:00:00Z,103,13.64,3.00,15.42,8.00,19.62,13.00,1.52,1.00,732.20,83.00,492.82,25.00,24.40,5.00,29.33,0.15,97.38,23.45,6.01,2.00,118.38,41.88,25.42,4.89,34.87,0.03,78.56,20.42,95.36,29.96,35.93,14.82,59.20,21.22,53.41,19.03
```
