package who

import (
	"fmt"
	"image/color"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/josh-allan/go-git/internal/ui"
	"github.com/mattn/go-runewidth"
)

var (
	hashStyle   = lipgloss.NewStyle().Foreground(ui.Purple)
	authorStyle = lipgloss.NewStyle().Foreground(ui.Green)
)

func Render(lines []BlameLine, termWidth int) string {
	if len(lines) == 0 {
		return ""
	}

	oldest := time.Now()
	newest := time.Time{}
	for _, l := range lines {
		if l.AuthorTime.Before(oldest) {
			oldest = l.AuthorTime
		}
		if l.AuthorTime.After(newest) {
			newest = l.AuthorTime
		}
	}
	timeRange := newest.Sub(oldest)
	if timeRange == 0 {
		timeRange = time.Hour
	}

	maxAuthorLen := 0
	for _, l := range lines {
		if len(l.Author) > maxAuthorLen {
			maxAuthorLen = len(l.Author)
		}
	}
	if maxAuthorLen > 20 {
		maxAuthorLen = 20
	}

	gutterWidth := 8 + 1 + maxAuthorLen + 1 + 10 + 1 + 5 + 3
	contentWidth := termWidth - gutterWidth
	if contentWidth < 20 {
		contentWidth = 20
	}

	commitColors := assignCommitColors(lines)

	var b strings.Builder
	prevHash := ""

	for _, l := range lines {
		recency := l.AuthorTime.Sub(oldest).Seconds() / timeRange.Seconds()
		dateColor := recencyColor(recency)
		dStyle := lipgloss.NewStyle().Foreground(dateColor)

		showMeta := l.Hash != prevHash
		prevHash = l.Hash

		hash := ""
		author := ""
		date := ""
		if showMeta {
			hash = hashStyle.Render(l.Hash[:8])
			authorName := runewidth.Truncate(l.Author, maxAuthorLen, "…")
			author = authorStyle.Render(fmt.Sprintf("%-*s", maxAuthorLen, authorName))
			date = dStyle.Render(l.AuthorTime.Format("2006-01-02"))
		} else {
			hash = ui.DimStyle.Render(strings.Repeat(" ", 8))
			author = strings.Repeat(" ", maxAuthorLen)
			date = strings.Repeat(" ", 10)
		}

		lineNum := ui.LineNumStyle.Render(fmt.Sprintf("%5d", l.LineNo))

		content := ui.ExpandTabs(l.Content, 4)
		content = runewidth.Truncate(content, contentWidth, "")

		contentColor := commitColors[l.Hash]
		contentStyled := lipgloss.NewStyle().Foreground(contentColor).Render(content)

		fmt.Fprintf(&b, "%s %s %s %s │ %s\n",
			hash, author, date, lineNum, contentStyled)
	}

	return b.String()
}

func assignCommitColors(lines []BlameLine) map[string]color.Color {
	seen := make(map[string]color.Color)
	idx := 0
	for _, l := range lines {
		if _, ok := seen[l.Hash]; ok {
			continue
		}
		p := ui.CommitPalette[idx%len(ui.CommitPalette)]
		seen[l.Hash] = ui.C(p.Light, p.Dark)
		idx++
	}
	return seen
}

func recencyColor(recency float64) color.Color {
	if recency > 0.8 {
		return ui.Red
	}
	if recency > 0.5 {
		return ui.Orange
	}
	if recency > 0.2 {
		return ui.DimGray
	}
	return ui.Dim
}
