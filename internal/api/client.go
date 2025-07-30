package api

import (
	"context"
	"net/url"
	"strings"
	"time"

	"github.com/fukuchancat/github-pr-metrics/pkg/utils"
	"github.com/google/go-github/v74/github"
)

// Wraps GitHub API with authentication and enterprise server support
type Client struct {
	client *github.Client
	ctx    context.Context
	logger *utils.Logger
}

// Configures GitHub API client with authentication and custom base URL support
func NewClient(apiURL, token string, logger *utils.Logger) (*Client, error) {
	ctx := context.Background()

	// Create a new client with auth token
	client := github.NewClient(nil).WithAuthToken(token)

	// Set custom API URL for GitHub Enterprise
	if apiURL != "https://api.github.com" {
		// Ensure the URL has a trailing slash
		if !strings.HasSuffix(apiURL, "/") {
			apiURL += "/"
		}

		baseURL, err := url.Parse(apiURL)
		if err != nil {
			return nil, err
		}
		client.BaseURL = baseURL
		logger.Debug("Using GitHub Enterprise API URL: %s", baseURL.String())
	}

	return &Client{
		client: client,
		ctx:    ctx,
		logger: logger,
	}, nil
}

// Fetches all PRs created within date range using paginated API calls
func (c *Client) GetPullRequests(owner, repo string, startDate, endDate time.Time) ([]*github.PullRequest, error) {
	c.logger.Debug("Fetching pull requests for %s/%s from %s to %s", owner, repo, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	opts := &github.PullRequestListOptions{
		State: "all",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	var allPRs []*github.PullRequest

	for {
		prs, resp, err := c.client.PullRequests.List(c.ctx, owner, repo, opts)
		if err != nil {
			return nil, err
		}

		// Filter PRs by date
		for _, pr := range prs {
			if pr.CreatedAt != nil {
				createdAt := pr.CreatedAt.Time
				if (createdAt.After(startDate) || createdAt.Equal(startDate)) &&
					(createdAt.Before(endDate) || createdAt.Equal(endDate)) {
					allPRs = append(allPRs, pr)
				}
			}
		}

		c.logger.Debug("Fetched page %d of pull requests (%d total so far)", opts.Page, len(allPRs))

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	c.logger.Debug("Fetched %d pull requests in total", len(allPRs))
	return allPRs, nil
}

// Fetches additions, deletions, and changed files count for a specific PR
func (c *Client) GetPRDetails(owner, repo string, number int) (*github.PullRequest, error) {
	c.logger.Debug("Fetching details for PR #%d", number)
	pr, _, err := c.client.PullRequests.Get(c.ctx, owner, repo, number)
	if err != nil {
		return nil, err
	}

	return pr, nil
}

// Fetches all commits associated with a PR using paginated requests
func (c *Client) GetPRCommits(owner, repo string, number int) ([]*github.RepositoryCommit, error) {
	c.logger.Debug("Fetching commits for PR #%d", number)
	opts := &github.ListOptions{
		PerPage: 100,
	}

	var allCommits []*github.RepositoryCommit

	for {
		commits, resp, err := c.client.PullRequests.ListCommits(c.ctx, owner, repo, number, opts)
		if err != nil {
			return nil, err
		}

		allCommits = append(allCommits, commits...)

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	c.logger.Debug("Fetched %d commits for PR #%d", len(allCommits), number)
	return allCommits, nil
}

// Fetches all review comments for a PR using paginated requests
func (c *Client) GetPRComments(owner, repo string, number int) ([]*github.PullRequestComment, error) {
	c.logger.Debug("Fetching comments for PR #%d", number)
	opts := &github.PullRequestListCommentsOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	var allComments []*github.PullRequestComment

	for {
		comments, resp, err := c.client.PullRequests.ListComments(c.ctx, owner, repo, number, opts)
		if err != nil {
			return nil, err
		}

		allComments = append(allComments, comments...)

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	c.logger.Debug("Fetched %d comments for PR #%d", len(allComments), number)
	return allComments, nil
}

// Fetches all code reviews for a PR using paginated requests
func (c *Client) GetPRReviews(owner, repo string, number int) ([]*github.PullRequestReview, error) {
	c.logger.Debug("Fetching reviews for PR #%d", number)
	opts := &github.ListOptions{
		PerPage: 100,
	}

	var allReviews []*github.PullRequestReview

	for {
		reviews, resp, err := c.client.PullRequests.ListReviews(c.ctx, owner, repo, number, opts)
		if err != nil {
			return nil, err
		}

		allReviews = append(allReviews, reviews...)

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	c.logger.Debug("Fetched %d reviews for PR #%d", len(allReviews), number)
	return allReviews, nil
}
