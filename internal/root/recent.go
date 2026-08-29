package root

import (
	"github.com/josh-allan/go-git/internal/branches"
	"github.com/josh-allan/go-git/pkg/git"
	"github.com/spf13/cobra"
)

func recentCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "recent",
		Short: "Select a recent branch to checkout",
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := git.LoadRepo()
			if err != nil {
				return err
			}
			return branches.Recent(repo)
		},
	}
}
