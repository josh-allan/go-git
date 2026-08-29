package who

import (
	"fmt"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/josh-allan/go-git/internal/ui"
	"github.com/mattn/go-runewidth"
)

func RenderSummary(files []FileSummary, termWidth int) string {
	if len(files) == 0 {
		return ""
	}

	oldest := time.Now()
	newest := time.Time{}
	for _, f := range files {
		if f.AuthorTime.Before(oldest) {
			oldest = f.AuthorTime
		}
		if f.AuthorTime.After(newest) {
			newest = f.AuthorTime
		}
	}
	timeRange := newest.Sub(oldest)
	if timeRange == 0 {
		timeRange = time.Hour
	}

	maxAuthor := 0
	maxPath := 0
	for _, f := range files {
		if len(f.Author) > maxAuthor {
			maxAuthor = len(f.Author)
		}
		if len(f.Path) > maxPath {
			maxPath = len(f.Path)
		}
	}
	if maxAuthor > 20 {
		maxAuthor = 20
	}
	if maxPath > 40 {
		maxPath = 40
	}

	gutterWidth := maxPath + 1 + maxAuthor + 1 + 10 + 3
	summaryWidth := termWidth - gutterWidth
	if summaryWidth < 10 {
		summaryWidth = 10
	}

	var b strings.Builder
	for _, f := range files {
		recency := f.AuthorTime.Sub(oldest).Seconds() / timeRange.Seconds()
		dateColor := recencyColor(recency)
		dStyle := lipgloss.NewStyle().Foreground(dateColor)

		path := runewidth.Truncate(f.Path, maxPath, "…")
		author := runewidth.Truncate(f.Author, maxAuthor, "…")
		summary := runewidth.Truncate(f.Summary, summaryWidth, "…")

		fmt.Fprintf(&b, "%s %s %s │ %s\n",
			lipgloss.NewStyle().Foreground(ui.Purple).Render(fmt.Sprintf("%-*s", maxPath, path)),
			lipgloss.NewStyle().Foreground(ui.Green).Render(fmt.Sprintf("%-*s", maxAuthor, author)),
			dStyle.Render(f.AuthorTime.Format("2006-01-02")),
			ui.DimStyle.Render(summary),
		)
	}

	return b.String()
}
