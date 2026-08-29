package differ

import (
	"strconv"
	"strings"
)

type DiffFile struct {
	OldName string
	NewName string
	Hunks   []Hunk
}

type Hunk struct {
	OldStart int
	OldCount int
	NewStart int
	NewCount int
	Lines    []Line
}

type LineType int

const (
	LineContext LineType = iota
	LineAdded
	LineRemoved
)

type Line struct {
	Type    LineType
	Content string
}

func Parse(raw string) []DiffFile {
	var files []DiffFile
	var current *DiffFile
	var currentHunk *Hunk

	for _, line := range strings.Split(raw, "\n") {
		switch {
		case strings.HasPrefix(line, "diff --git"):
			files = append(files, DiffFile{})
			current = &files[len(files)-1]
			currentHunk = nil

		case strings.HasPrefix(line, "--- "):
			if current != nil {
				current.OldName = strings.TrimPrefix(line, "--- ")
			}

		case strings.HasPrefix(line, "+++ "):
			if current != nil {
				current.NewName = strings.TrimPrefix(line, "+++ ")
			}

		case strings.HasPrefix(line, "@@"):
			if current == nil {
				continue
			}
			h := parseHunkHeader(line)
			current.Hunks = append(current.Hunks, h)
			currentHunk = &current.Hunks[len(current.Hunks)-1]

		default:
			if currentHunk == nil {
				continue
			}
			switch {
			case strings.HasPrefix(line, "+"):
				currentHunk.Lines = append(currentHunk.Lines, Line{LineAdded, line[1:]})
			case strings.HasPrefix(line, "-"):
				currentHunk.Lines = append(currentHunk.Lines, Line{LineRemoved, line[1:]})
			default:
				content := line
				if len(content) > 0 && content[0] == ' ' {
					content = content[1:]
				}
				currentHunk.Lines = append(currentHunk.Lines, Line{LineContext, content})
			}
		}
	}

	return files
}

func parseHunkHeader(line string) Hunk {
	// @@ -oldStart,oldCount +newStart,newCount @@
	line = strings.TrimPrefix(line, "@@ ")
	parts := strings.SplitN(line, " @@", 2)
	if len(parts) == 0 {
		return Hunk{}
	}

	ranges := strings.Fields(parts[0])
	h := Hunk{}

	if len(ranges) >= 1 {
		h.OldStart, h.OldCount = parseRange(strings.TrimPrefix(ranges[0], "-"))
	}
	if len(ranges) >= 2 {
		h.NewStart, h.NewCount = parseRange(strings.TrimPrefix(ranges[1], "+"))
	}

	return h
}

func parseRange(s string) (int, int) {
	parts := strings.SplitN(s, ",", 2)
	start, _ := strconv.Atoi(parts[0])
	count := 1
	if len(parts) == 2 {
		count, _ = strconv.Atoi(parts[1])
	}
	return start, count
}
