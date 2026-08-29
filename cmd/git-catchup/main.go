package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/fang"
	"github.com/charmbracelet/x/term"
	"github.com/josh-allan/go-git/internal/catchup"
	"github.com/josh-allan/go-git/internal/differ"
	"github.com/josh-allan/go-git/pkg/git"
	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	cmd := &cobra.Command{
		Use:   "git-catchup [branch]",
		Short: "See what changed on a branch since you last looked",
		Args:  cobra.MaximumNArgs(1),
		RunE:  run,
	}

	cmd.Flags().String("since", "", "time window (e.g. 3d, 1w, 2025-08-20)")
	cmd.Flags().Bool("full", false, "show all files instead of directory summary")

	if err := fang.Execute(
		context.Background(),
		cmd,
		fang.WithVersion(version),
		fang.WithNotifySignal(os.Interrupt),
	); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	repo, err := git.LoadRepo()
	if err != nil {
		return err
	}

	branch := ""
	if len(args) > 0 {
		branch = args[0]
	} else {
		branch = defaultBranch(repo)
	}

	since, _ := cmd.Flags().GetString("since")
	if since == "" {
		since = inferSince(repo, branch)
	}

	summary, err := catchup.Run(repo, branch, since)
	if err != nil {
		return err
	}
	if summary == nil {
		fmt.Println("Nothing to catch up on.")
		return nil
	}

	width, height, _ := term.GetSize(os.Stdout.Fd())
	if width <= 0 {
		width = 80
	}
	if height <= 0 {
		height = 24
	}

	full, _ := cmd.Flags().GetBool("full")
	rendered := catchup.Render(summary, width, full)
	return differ.Page(rendered, width, height)
}

func defaultBranch(repo *git.Repo) string {
	out, err := repo.Git("symbolic-ref", "refs/remotes/origin/HEAD")
	if err == nil {
		out = strings.TrimSpace(out)
		if ref, ok := strings.CutPrefix(out, "refs/remotes/origin/"); ok {
			return ref
		}
	}
	for _, name := range []string{"main", "master", "trunk"} {
		if _, err := repo.Git("rev-parse", "--verify", name); err == nil {
			return name
		}
	}
	return "HEAD"
}

func inferSince(repo *git.Repo, branch string) string {
	email, err := repo.Git("config", "user.email")
	if err != nil {
		return "1w"
	}

	out, err := repo.Git("log", "-1", "--format=%aI", "--author="+email, branch)
	if err != nil || out == "" {
		return "1w"
	}
	return out
}
