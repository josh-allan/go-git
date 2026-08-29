package main

import (
	"context"
	"fmt"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/x/term"
	"github.com/josh-allan/go-git/internal/differ"
	"github.com/josh-allan/go-git/internal/who"
	"github.com/josh-allan/go-git/pkg/git"
	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	cmd := &cobra.Command{
		Use:   "git-who [path]",
		Short: "File ownership and contributor breakdown",
		Long:  "Shows who last touched each file, or a contributor breakdown for a specific file.",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runWho,
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

func runWho(cmd *cobra.Command, args []string) error {
	repo, err := git.LoadRepo()
	if err != nil {
		return err
	}

	width, height, _ := term.GetSize(os.Stdout.Fd())
	if width <= 0 {
		width = 80
	}
	if height <= 0 {
		height = 24
	}

	if len(args) == 0 {
		return showSummary(repo, width, height)
	}

	path := args[0]
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("cannot access %s: %w", path, err)
	}

	if info.IsDir() {
		return showSummary(repo, width, height, path)
	}
	return showContributors(repo, path, width, height)
}

func showSummary(repo *git.Repo, width, height int, paths ...string) error {
	files, err := who.RunSummary(repo, paths...)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		fmt.Println("No tracked files.")
		return nil
	}
	return differ.Page(who.RenderSummary(files, width), width, height)
}

func showContributors(repo *git.Repo, file string, width, height int) error {
	contribs, err := who.RunContributors(repo, file)
	if err != nil {
		return err
	}
	if len(contribs) == 0 {
		fmt.Println("No blame data.")
		return nil
	}
	return differ.Page(who.RenderContributors(contribs, file, width), width, height)
}
