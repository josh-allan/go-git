package root

import (
	"github.com/josh-allan/go-git/internal/branches"
	"github.com/josh-allan/go-git/pkg/git"
	"github.com/spf13/cobra"
)

func squashedCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "squashed",
		Short: "Select squashed branches to delete",
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := git.LoadRepo()
			if err != nil {
				return err
			}

			target, _ := cmd.Flags().GetString("target")
			force, _ := cmd.Flags().GetBool("force")
			return branches.Squashed(repo, target, force)
		},
	}

	cmd.Flags().StringP("target", "t", "origin/main", "target branch to check against")
	cmd.Flags().BoolP("force", "f", false, "force delete branches")

	return cmd
}
