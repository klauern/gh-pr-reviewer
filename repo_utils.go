package main

import (
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/repository"
)

func GetCurrentRepository() (repository.Repository, error) {
	return repository.Current()
}

func GetCurrentBranchPR(owner, repo string) (int, error) {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return 0, fmt.Errorf("failed to create GitHub client: %w", err)
	}

	branch, err := getCurrentBranch()
	if err != nil {
		return 0, fmt.Errorf("failed to get current branch: %w", err)
	}

	var response struct {
		Items []struct {
			Number int `json:"number"`
		} `json:"items"`
	}
	err = client.Get(fmt.Sprintf("search/issues?q=repo:%s/%s+head:%s+type:pr+state:open", owner, repo, branch), &response)
	if err != nil {
		return 0, fmt.Errorf("failed to search for PRs: %w", err)
	}

	if len(response.Items) == 0 {
		return 0, fmt.Errorf("no open PR found for the current branch")
	}

	return response.Items[0].Number, nil
}

func getCurrentBranch() (string, error) {
	output, _, err := git("branch", "--show-current")
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	branchBytes, err := io.ReadAll(output)
	if err != nil {
		return "", fmt.Errorf("failed to read branch name: %w", err)
	}
	return strings.TrimSpace(string(branchBytes)), nil
}

func git(args ...string) (stdout, stderr io.ReadCloser, err error) {
	cmd := exec.Command("git", args...)
	stdout, err = cmd.StdoutPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	stderr, err = cmd.StderrPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}
	err = cmd.Start()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start gh command: %w", err)
	}
	return stdout, stderr, nil
}
