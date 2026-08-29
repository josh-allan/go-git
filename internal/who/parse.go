package who

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/josh-allan/go-git/pkg/git"
)

type BlameLine struct {
	Hash       string
	Author     string
	AuthorTime time.Time
	Summary    string
	LineNo     int
	Content    string
}

func Run(repo *git.Repo, file string) ([]BlameLine, error) {
	output, err := repo.Git("blame", "--porcelain", "--use-mailmap", file)
	if err != nil {
		return nil, fmt.Errorf("running git blame: %w", err)
	}
	return parsePorcelain(output), nil
}

func parsePorcelain(raw string) []BlameLine {
	lines := strings.Split(raw, "\n")

	commits := make(map[string]*commitInfo)
	var result []BlameLine
	var current *commitInfo
	var lineNo int

	for _, line := range lines {
		if len(line) >= 40 && isHexPrefix(line[:40]) {
			parts := strings.Fields(line)
			hash := parts[0]

			if len(parts) >= 3 {
				lineNo, _ = strconv.Atoi(parts[2])
			}

			if ci, ok := commits[hash]; ok {
				current = ci
			} else {
				current = &commitInfo{hash: hash}
				commits[hash] = current
			}
			continue
		}

		if current == nil {
			continue
		}

		if strings.HasPrefix(line, "\t") {
			result = append(result, BlameLine{
				Hash:       current.hash,
				Author:     current.author,
				AuthorTime: current.authorTime,
				Summary:    current.summary,
				LineNo:     lineNo,
				Content:    line[1:],
			})
			continue
		}

		key, val, ok := strings.Cut(line, " ")
		if !ok {
			continue
		}

		switch key {
		case "author":
			current.author = val
		case "author-time":
			ts, _ := strconv.ParseInt(val, 10, 64)
			current.authorTime = time.Unix(ts, 0)
		case "summary":
			current.summary = val
		}
	}

	return result
}

type commitInfo struct {
	hash       string
	author     string
	authorTime time.Time
	summary    string
}

func isHexPrefix(s string) bool {
	for _, c := range s {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
			return false
		}
	}
	return true
}
