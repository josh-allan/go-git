package who

import (
	"fmt"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
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

	workers := runtime.NumCPU()
	if workers > 16 {
		workers = 16
	}

	type indexedSummary struct {
		index int
		fs    FileSummary
	}

	ch := make(chan int, len(files))
	for i := range files {
		ch <- i
	}
	close(ch)

	var mu sync.Mutex
	var collected []indexedSummary
	var wg sync.WaitGroup

	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range ch {
				out, err := repo.Git("log", "-1", "--format=%aN%x00%at%x00%s", "--", files[i])
				if err != nil || strings.TrimSpace(out) == "" {
					continue
				}
				parts := strings.SplitN(strings.TrimSpace(out), "\x00", 3)
				if len(parts) < 3 {
					continue
				}
				ts, _ := strconv.ParseInt(parts[1], 10, 64)
				item := indexedSummary{
					index: i,
					fs: FileSummary{
						Path:       files[i],
						Author:     parts[0],
						AuthorTime: time.Unix(ts, 0),
						Summary:    parts[2],
					},
				}
				mu.Lock()
				collected = append(collected, item)
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	slices.SortFunc(collected, func(a, b indexedSummary) int {
		return a.index - b.index
	})
	result := make([]FileSummary, len(collected))
	for i, s := range collected {
		result[i] = s.fs
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
