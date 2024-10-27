package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/cli/go-gh/v2/pkg/api"
)

func GetPullRequestDiff(owner, repo string, prNumber int) (string, error) {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return "", fmt.Errorf("failed to create GitHub client: %w", err)
	}

	url := fmt.Sprintf("repos/%s/%s/pulls/%d", owner, repo, prNumber)
	resp, err := client.Request("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get pull request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	diffURL := resp.Header.Get("Location")
	if diffURL == "" {
		return "", fmt.Errorf("diff URL not found in response headers")
	}

	diffResp, err := http.Get(diffURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch diff: %w", err)
	}
	defer diffResp.Body.Close()

	if diffResp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code when fetching diff: %d", diffResp.StatusCode)
	}

	diffContent, err := io.ReadAll(diffResp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read diff content: %w", err)
	}

	return string(diffContent), nil
}
