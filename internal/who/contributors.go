package who

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/josh-allan/go-git/internal/ui"
	"github.com/josh-allan/go-git/pkg/git"
	"github.com/mattn/go-runewidth"
)

type Contributor struct {
	Author     string
	Lines      int
	Percentage float64
	LastActive time.Time
}

func RunContributors(repo *git.Repo, file string) ([]Contributor, error) {
	lines, err := Run(repo, file)
	if err != nil {
		return nil, err
	}
	if len(lines) == 0 {
		return nil, nil
	}

	type stats struct {
		name       string
		lines      int
		lastActive time.Time
	}

	byEmail := make(map[string]*stats)
	for _, l := range lines {
		email := l.AuthorEmail
		if email == "" {
			email = l.Author
		}
		s, ok := byEmail[email]
		if !ok {
			s = &stats{name: l.Author}
			byEmail[email] = s
		}
		s.lines++
		if l.AuthorTime.After(s.lastActive) {
			s.lastActive = l.AuthorTime
			s.name = l.Author
		}
	}

	total := len(lines)
	var result []Contributor
	for _, s := range byEmail {
		result = append(result, Contributor{
			Author:     s.name,
			Lines:      s.lines,
			Percentage: float64(s.lines) / float64(total) * 100,
			LastActive: s.lastActive,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Lines > result[j].Lines
	})

	return result, nil
}

func RenderContributors(contribs []Contributor, file string, termWidth int) string {
	if len(contribs) == 0 {
		return ""
	}

	maxAuthor := 0
	for _, c := range contribs {
		if len(c.Author) > maxAuthor {
			maxAuthor = len(c.Author)
		}
	}
	if maxAuthor > 30 {
		maxAuthor = 30
	}

	barMax := termWidth - maxAuthor - 28
	if barMax < 10 {
		barMax = 10
	}

	topLines := contribs[0].Lines

	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(ui.Purple)

	var b strings.Builder
	b.WriteString(headerStyle.Render(file))
	b.WriteString("\n\n")

	for i, c := range contribs {
		author := runewidth.Truncate(c.Author, maxAuthor, "…")

		barLen := ui.ScaledBarWidth(c.Lines, topLines, barMax)

		palette := ui.CommitPalette[i%len(ui.CommitPalette)]
		barColor := ui.C(palette.Light, palette.Dark)
		bar := lipgloss.NewStyle().Foreground(barColor).Render(strings.Repeat("█", barLen))

		fmt.Fprintf(&b, "  %-*s %5d %5.1f%%  %s  %s\n",
			maxAuthor,
			lipgloss.NewStyle().Foreground(ui.Green).Render(author),
			c.Lines,
			c.Percentage,
			bar,
			ui.DimStyle.Render(c.LastActive.Format("2006-01-02")),
		)
	}

	return b.String()
}
