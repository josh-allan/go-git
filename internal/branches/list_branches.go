package branches

import (
	"fmt"
	"strings"

	"github.com/josh-allan/go-git/pkg/git"
)

func listRecentBranches(repo *git.Repo) ([]string, error) {
	output, err := repo.Git(
		"for-each-ref",
		"refs/heads/",
		"--sort=committerdate",
		"--format=%(committerdate:short) %(refname:short)",
	)
	if err != nil {
		return nil, fmt.Errorf("listing recent branches: %w", err)
	}

	return splitLines(output), nil
}

func listBranchesFiltered(repo *git.Repo, flag, target string) ([]string, error) {
	output, err := repo.Git("branch", flag, target)
	if err != nil {
		return nil, fmt.Errorf("listing branches with %s: %w", flag, err)
	}

	return splitLines(output), nil
}

func splitLines(output string) []string {
	if output == "" {
		return nil
	}

	lines := strings.Split(output, "\n")
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(strings.Trim(line, "'"))
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
