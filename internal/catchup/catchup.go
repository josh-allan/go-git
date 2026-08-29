package catchup

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/josh-allan/go-git/pkg/git"
)

type Summary struct {
	Branch      string
	Since       string
	Commits     int
	AuthorCount int
	Span        time.Duration
	Files       []FileChurn
	Dirs        []DirChurn
	Authors     []AuthorSummary
	Added       []string
	Deleted     []string
}

type DirChurn struct {
	Path      string
	Additions int
	Deletions int
	FileCount int
	TopFiles  []FileChurn
}

type FileChurn struct {
	Path      string
	Additions int
	Deletions int
	TopAuthor string
}

type AuthorSummary struct {
	Name    string
	Commits int
	TopPath string
}

func Run(repo *git.Repo, branch, since string) (*Summary, error) {
	logArgs := []string{"log", "--format=%H%x00%an%x00%aI%x00%s", "--numstat"}
	if since != "" {
		logArgs = append(logArgs, "--since="+since)
	}
	logArgs = append(logArgs, branch)

	output, err := repo.Git(logArgs[0], logArgs[1:]...)
	if err != nil {
		return nil, fmt.Errorf("running git log: %w", err)
	}
	if strings.TrimSpace(output) == "" {
		return nil, nil
	}

	s, err := parseLog(output, branch, since)
	if err != nil {
		return nil, err
	}

	s.Added, s.Deleted = listAddedDeleted(repo, branch, since)

	return s, nil
}

type parsedCommit struct {
	author string
	time   time.Time
	files  map[string][2]int
}

func parseLog(output, branch, since string) (*Summary, error) {
	lines := strings.Split(output, "\n")

	var commits []parsedCommit
	var cur *parsedCommit

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "\x00", 4)
		if len(parts) == 4 {
			t, _ := time.Parse(time.RFC3339, parts[2])
			commits = append(commits, parsedCommit{author: parts[1], time: t, files: make(map[string][2]int)})
			cur = &commits[len(commits)-1]
			continue
		}

		if cur == nil {
			continue
		}

		fields := strings.SplitN(line, "\t", 3)
		if len(fields) != 3 {
			continue
		}
		add, _ := strconv.Atoi(fields[0])
		del, _ := strconv.Atoi(fields[1])
		cur.files[fields[2]] = [2]int{add, del}
	}

	if len(commits) == 0 {
		return nil, nil
	}

	type fileAgg struct {
		add, del int
		authors  map[string]int
	}
	type authorAgg struct {
		commits int
		paths   map[string]int
	}

	fileMap := make(map[string]*fileAgg)
	authorMap := make(map[string]*authorAgg)
	earliest, latest := commits[0].time, commits[0].time

	for _, c := range commits {
		if c.time.Before(earliest) {
			earliest = c.time
		}
		if c.time.After(latest) {
			latest = c.time
		}

		a := authorMap[c.author]
		if a == nil {
			a = &authorAgg{paths: make(map[string]int)}
			authorMap[c.author] = a
		}
		a.commits++

		for path, stat := range c.files {
			f := fileMap[path]
			if f == nil {
				f = &fileAgg{authors: make(map[string]int)}
				fileMap[path] = f
			}
			f.add += stat[0]
			f.del += stat[1]
			f.authors[c.author]++
			a.paths[path]++
		}
	}

	files := make([]FileChurn, 0, len(fileMap))
	for path, f := range fileMap {
		files = append(files, FileChurn{
			Path:      path,
			Additions: f.add,
			Deletions: f.del,
			TopAuthor: maxKey(f.authors),
		})
	}
	sort.Slice(files, func(i, j int) bool {
		return (files[i].Additions + files[i].Deletions) > (files[j].Additions + files[j].Deletions)
	})

	authors := make([]AuthorSummary, 0, len(authorMap))
	for name, a := range authorMap {
		authors = append(authors, AuthorSummary{
			Name:    name,
			Commits: a.commits,
			TopPath: maxKey(a.paths),
		})
	}
	sort.Slice(authors, func(i, j int) bool {
		return authors[i].Commits > authors[j].Commits
	})

	dirs := groupByDir(files)

	return &Summary{
		Branch:      branch,
		Since:       since,
		Commits:     len(commits),
		AuthorCount: len(authorMap),
		Span:        latest.Sub(earliest),
		Files:       files,
		Dirs:        dirs,
		Authors:     authors,
	}, nil
}

func groupByDir(files []FileChurn) []DirChurn {
	dirMap := make(map[string]*DirChurn)

	for _, f := range files {
		dir := "."
		if idx := strings.LastIndex(f.Path, "/"); idx >= 0 {
			dir = f.Path[:idx]
		}

		d := dirMap[dir]
		if d == nil {
			d = &DirChurn{Path: dir}
			dirMap[dir] = d
		}
		d.Additions += f.Additions
		d.Deletions += f.Deletions
		d.FileCount++
		if len(d.TopFiles) < 3 {
			d.TopFiles = append(d.TopFiles, f)
		}
	}

	dirs := make([]DirChurn, 0, len(dirMap))
	for _, d := range dirMap {
		dirs = append(dirs, *d)
	}
	sort.Slice(dirs, func(i, j int) bool {
		return (dirs[i].Additions + dirs[i].Deletions) > (dirs[j].Additions + dirs[j].Deletions)
	})

	return dirs
}

func listAddedDeleted(repo *git.Repo, branch, since string) (added, deleted []string) {
	args := []string{"log", "--diff-filter=AD", "--name-status", "--format="}
	if since != "" {
		args = append(args, "--since="+since)
	}
	args = append(args, branch)

	output, err := repo.Git(args[0], args[1:]...)
	if err != nil {
		return nil, nil
	}

	addSeen := make(map[string]bool)
	delSeen := make(map[string]bool)
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if len(line) < 3 {
			continue
		}
		status, name := line[0], strings.TrimSpace(line[1:])
		switch status {
		case 'A':
			if !addSeen[name] {
				addSeen[name] = true
				added = append(added, name)
			}
		case 'D':
			if !delSeen[name] {
				delSeen[name] = true
				deleted = append(deleted, name)
			}
		}
	}
	return added, deleted
}

func maxKey(m map[string]int) string {
	best, bestN := "", 0
	for k, n := range m {
		if n > bestN {
			best, bestN = k, n
		}
	}
	return best
}
