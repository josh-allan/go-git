package root

import (
	"github.com/josh-allan/go-git/internal/branches"
	"github.com/josh-allan/go-git/pkg/git"
	"github.com/spf13/cobra"
)

func mergedCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "merged",
		Short: "Select merged branches to delete",
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := git.LoadRepo()
			if err != nil {
				return err
			}

			target, _ := cmd.Flags().GetString("target")
			return branches.Merged(repo, target)
		},
	}

	cmd.Flags().StringP("target", "t", "origin/main", "target branch to check against")

	return cmd
}
