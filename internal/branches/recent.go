package branches

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/josh-allan/go-git/pkg/git"
	"github.com/josh-allan/go-git/pkg/prompt"
)

func Recent(repo *git.Repo) error {
	branches, err := listRecentBranches(repo)
	if err != nil {
		return fmt.Errorf("listing recent branches: %w", err)
	}

	if len(branches) == 0 {
		fmt.Println("No branches found.")
		return nil
	}

	var selected string

	err = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a branch to checkout").
				Options(prompt.AsOptions(branches)...).
				Value(&selected),
		),
	).Run()
	if err != nil {
		return fmt.Errorf("running form: %w", err)
	}

	parts := strings.Fields(selected)
	if len(parts) < 2 {
		return fmt.Errorf("unexpected branch format: %q", selected)
	}

	_, err = repo.Git("checkout", parts[1])
	if err != nil {
		return fmt.Errorf("checking out branch: %w", err)
	}

	return nil
}
