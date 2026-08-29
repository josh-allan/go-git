package who

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/josh-allan/go-git/pkg/git"
)

type FileSummary struct {
	Path        string
	Author      string
	AuthorEmail string
	AuthorTime  time.Time
	Summary     string
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

	remaining := make(map[string]struct{}, len(files))
	for _, f := range files {
		remaining[f] = struct{}{}
	}

	type commitHeader struct {
		author      string
		authorEmail string
		authorTime  time.Time
		summary     string
	}

	seen := make(map[string]*commitHeader, len(files))

	// Walk the log once with --name-only. Each commit appears as:
	//   author\x00email\x00timestamp\x00subject
	//   <blank line>
	//   file1
	//   file2
	//   <blank line>
	logOutput, err := repo.Git("log", "--format=%an%x00%aE%x00%at%x00%s", "--name-only", "--diff-filter=ACDMRT")
	if err != nil {
		return nil, fmt.Errorf("running git log: %w", err)
	}

	lines := strings.Split(logOutput, "\n")
	var cur *commitHeader

	for _, line := range lines {
		if len(remaining) == 0 {
			break
		}

		line = strings.TrimRight(line, "\r")

		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "\x00", 4)
		if len(parts) == 4 {
			ts, _ := strconv.ParseInt(parts[2], 10, 64)
			cur = &commitHeader{
				author:      parts[0],
				authorEmail: parts[1],
				authorTime:  time.Unix(ts, 0),
				summary:     parts[3],
			}
			continue
		}

		if cur == nil {
			continue
		}

		if _, tracked := remaining[line]; tracked {
			if _, already := seen[line]; !already {
				seen[line] = cur
				delete(remaining, line)
			}
		}
	}

	result := make([]FileSummary, 0, len(seen))
	for _, f := range files {
		c, ok := seen[f]
		if !ok {
			continue
		}
		result = append(result, FileSummary{
			Path:        f,
			Author:      c.author,
			AuthorEmail: c.authorEmail,
			AuthorTime:  c.authorTime,
			Summary:     c.summary,
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
