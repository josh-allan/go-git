package differ

import (
	"fmt"
	"strings"

	"github.com/josh-allan/go-git/internal/ui"
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
		name := strings.TrimPrefix(f.NewName, "b/")
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

	maxBarWidth := termWidth - maxNameLen - 15
	if maxBarWidth < 10 {
		maxBarWidth = 10
	}
	var b strings.Builder
	totalAdded, totalRemoved := 0, 0

	for _, s := range stats {
		totalAdded += s.added
		totalRemoved += s.removed

		total := s.added + s.removed
		bar := ""
		if total > 0 {
			width := ui.ScaledBarWidth(total, maxChanges, maxBarWidth)
			addBar := s.added * width / total
			remBar := width - addBar

			bar = addedStyle.Render(strings.Repeat("+", addBar)) +
				removedStyle.Render(strings.Repeat("-", remBar))
		}

		fmt.Fprintf(&b, " %-*s | %4d %s\n", maxNameLen, s.name, total, bar)
	}

	fmt.Fprintf(&b, " %d files changed, %d insertions(+), %d deletions(-)\n",
		len(files), totalAdded, totalRemoved)

	return b.String()
}
