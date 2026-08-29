package main

import (
	"context"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/josh-allan/go-git/internal/branches"
	"github.com/josh-allan/go-git/pkg/git"
	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	cmd := &cobra.Command{
		Use:   "git-merged",
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

	if err := fang.Execute(
		context.Background(),
		cmd,
		fang.WithVersion(version),
		fang.WithNotifySignal(os.Interrupt),
	); err != nil {
		os.Exit(1)
	}
}
