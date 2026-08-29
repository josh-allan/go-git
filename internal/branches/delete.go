package branches

import (
	"fmt"

	"github.com/josh-allan/go-git/pkg/git"
)

func deleteBranches(repo *git.Repo, branches []string, flag string) {
	for _, branch := range branches {
		_, err := repo.Git("branch", flag, branch)
		if err != nil {
			fmt.Printf("failed to delete branch %s: %v\n", branch, err)
		} else {
			fmt.Printf("deleted branch: %s\n", branch)
		}
	}
}
