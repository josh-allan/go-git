package differ

import (
	"fmt"
	"strings"
)

func RenderStat(files []DiffFile, termWidth int) string {
	if len(files) == 0 {
		return ""
	}

	maxNameLen := 0
	type fileStat struct {
		name    string
		added   int
		removed int
	}

	stats := make([]fileStat, len(files))
	maxChanges := 0

	for i, f := range files {
		name := f.NewName
		if strings.HasPrefix(name, "b/") {
			name = name[2:]
		}
		if len(name) > maxNameLen {
			maxNameLen = len(name)
		}

		var added, removed int
		for _, h := range f.Hunks {
			for _, l := range h.Lines {
				switch l.Type {
				case LineAdded:
					added++
				case LineRemoved:
					removed++
				}
			}
		}
		stats[i] = fileStat{name, added, removed}
		if added+removed > maxChanges {
			maxChanges = added + removed
		}
	}

	barWidth := termWidth - maxNameLen - 15
	if barWidth < 10 {
		barWidth = 10
	}

	var b strings.Builder
	totalAdded, totalRemoved := 0, 0

	for _, s := range stats {
		totalAdded += s.added
		totalRemoved += s.removed

		total := s.added + s.removed
		bar := ""
		if maxChanges > 0 && total > 0 {
			width := total * barWidth / maxChanges
			if width < 1 {
				width = 1
			}
			addBar := s.added * width / total
			remBar := width - addBar

			bar = addedStyle.Render(strings.Repeat("+", addBar)) +
				removedStyle.Render(strings.Repeat("-", remBar))
		}

		b.WriteString(fmt.Sprintf(" %-*s | %4d %s\n", maxNameLen, s.name, total, bar))
	}

	b.WriteString(fmt.Sprintf(" %d files changed, %d insertions(+), %d deletions(-)\n",
		len(files), totalAdded, totalRemoved))

	return b.String()
}
