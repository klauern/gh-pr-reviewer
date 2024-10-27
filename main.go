package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	fmt.Println("gh-pr-reviewer extension")

	var prNumber int
	var err error

	repo, err := GetCurrentRepository()
	if err != nil {
		fmt.Printf("Error getting current repository: %v\n", err)
		return
	}

	switch len(os.Args) {
	case 1:
		// No arguments provided, try to infer PR from current branch
		prNumber, err = GetCurrentBranchPR(repo.Owner, repo.Name)
		if err != nil {
			fmt.Printf("Error inferring PR number: %v\n", err)
			return
		}
	case 2:
		// One argument provided, assume it's the PR number
		prNumber, err = strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Printf("Invalid PR number: %s\n", os.Args[1])
			return
		}
	default:
		fmt.Println("Usage: gh pr-reviewer [pr_number]")
		return
	}

	diff, err := GetPullRequestDiff(repo.Owner, repo.Name, prNumber)
	if err != nil {
		fmt.Printf("Error fetching pull request diff: %v\n", err)
		return
	}

	fmt.Println("Pull Request Diff:")
	fmt.Println(diff)

	// TODO: Pass the diff to your LLM for review
}

// For more examples of using go-gh, see:
// https://github.com/cli/go-gh/blob/trunk/example_gh_test.go
