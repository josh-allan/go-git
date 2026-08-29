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
		Use:   "git-recent",
		Short: "Select a recent branch to checkout",
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := git.LoadRepo()
			if err != nil {
				return err
			}
			return branches.Recent(repo)
		},
	}

	if err := fang.Execute(
		context.Background(),
		cmd,
		fang.WithVersion(version),
		fang.WithNotifySignal(os.Interrupt),
	); err != nil {
		os.Exit(1)
	}
}
