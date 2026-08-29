package branches

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/josh-allan/go-git/pkg/git"
	"github.com/josh-allan/go-git/pkg/prompt"
)

func Merged(repo *git.Repo, target string) error {
	branches, err := listBranchesFiltered(repo, "--merged", target)
	if err != nil {
		return fmt.Errorf("listing merged branches: %w", err)
	}

	if len(branches) == 0 {
		fmt.Println("No merged branches found.")
		return nil
	}

	var selected []string

	err = huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select branches to delete").
				Options(prompt.AsOptions(branches)...).
				Value(&selected),
		),
	).Run()
	if err != nil {
		return fmt.Errorf("running form: %w", err)
	}

	deleteBranches(repo, selected, "-d")
	return nil
}
