package main

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/x/term"
	"github.com/josh-allan/go-git/internal/differ"
	"github.com/josh-allan/go-git/pkg/git"
	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	cmd := &cobra.Command{
		Use:   "git-differ [paths...]",
		Short: "Side-by-side diff viewer with styled output",
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := git.LoadRepo()
			if err != nil {
				return err
			}

			cached, _ := cmd.Flags().GetBool("cached")
			stat, _ := cmd.Flags().GetBool("stat")

			diffArgs := []string{"diff", "--no-color"}
			if cached {
				diffArgs = append(diffArgs, "--cached")
			}
			diffArgs = append(diffArgs, args...)

			output, err := repo.Git(diffArgs[0], diffArgs[1:]...)
			if err != nil {
				return fmt.Errorf("running git diff: %w", err)
			}

			if output == "" {
				fmt.Println("No changes.")
				return nil
			}

			files := differ.Parse(output)

			width, height, _ := term.GetSize(os.Stdout.Fd())
			if width <= 0 {
				width = 80
			}
			if height <= 0 {
				height = 24
			}

			if stat {
				fmt.Print(differ.RenderStat(files, width))
				return nil
			}

			rendered := differ.Render(files, width)
			return differ.Page(rendered, width, height)
		},
	}

	cmd.Flags().Bool("cached", false, "show staged changes")
	cmd.Flags().Bool("stat", false, "show diffstat summary")

	if err := fang.Execute(
		context.Background(),
		cmd,
		fang.WithVersion(version),
		fang.WithNotifySignal(os.Interrupt),
	); err != nil {
		os.Exit(1)
	}
}
