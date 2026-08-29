package who

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/josh-allan/go-git/pkg/git"
)

type FileSummary struct {
	Path       string
	Author     string
	AuthorTime time.Time
	Summary    string
}

func RunSummary(repo *git.Repo, paths ...string) ([]FileSummary, error) {
	args := []string{"ls-files"}
	args = append(args, paths...)
	output, err := repo.Git(args[0], args[1:]...)
	if err != nil {
		return nil, fmt.Errorf("listing files: %w", err)
	}

	files := splitNonEmpty(output)
	if len(files) == 0 {
		return nil, nil
	}

	var result []FileSummary
	for _, f := range files {
		out, err := repo.Git("log", "-1", "--format=%an%x00%at%x00%s", "--", f)
		if err != nil || strings.TrimSpace(out) == "" {
			continue
		}
		parts := strings.SplitN(strings.TrimSpace(out), "\x00", 3)
		if len(parts) < 3 {
			continue
		}
		ts, _ := strconv.ParseInt(parts[1], 10, 64)
		result = append(result, FileSummary{
			Path:       f,
			Author:     parts[0],
			AuthorTime: time.Unix(ts, 0),
			Summary:    parts[2],
		})
	}

	return result, nil
}

func splitNonEmpty(s string) []string {
	var out []string
	for _, line := range strings.Split(s, "\n") {
		if line = strings.TrimSpace(line); line != "" {
			out = append(out, line)
		}
	}
	return out
}
