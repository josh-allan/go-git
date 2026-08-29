package catchup

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/josh-allan/go-git/internal/ui"
	"github.com/mattn/go-runewidth"
)

func Render(s *Summary, termWidth int, full bool) string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(ui.Purple)
	labelStyle := lipgloss.NewStyle().Foreground(ui.Green)

	b.WriteString(headerStyle.Render(fmt.Sprintf("Catchup: %s", s.Branch)))
	if s.Since != "" {
		b.WriteString(ui.DimStyle.Render(fmt.Sprintf("  (since %s)", s.Since)))
	}
	b.WriteString("\n")

	days := int(s.Span.Hours() / 24)
	spanText := "same day"
	if days == 1 {
		spanText = "1 day"
	} else if days > 1 {
		spanText = fmt.Sprintf("%d days", days)
	}
	commitWord := "commits"
	if s.Commits == 1 {
		commitWord = "commit"
	}
	authorWord := "authors"
	if s.AuthorCount == 1 {
		authorWord = "author"
	}
	fmt.Fprintf(&b, "%d %s by %d %s over %s\n\n",
		s.Commits, commitWord, s.AuthorCount, authorWord, spanText)

	if full {
		if len(s.Files) > 0 {
			b.WriteString(labelStyle.Render("Files changed"))
			b.WriteString("\n")
			renderFiles(&b, s.Files, termWidth)
			b.WriteString("\n")
		}
	} else if len(s.Dirs) > 0 {
		b.WriteString(labelStyle.Render("Areas changed"))
		b.WriteString("\n")
		renderDirs(&b, s.Dirs, termWidth)
		b.WriteString("\n")
	}

	if len(s.Authors) > 0 {
		b.WriteString(labelStyle.Render("Contributors"))
		b.WriteString("\n")
		renderAuthors(&b, s.Authors, termWidth)
		b.WriteString("\n")
	}

	if len(s.Added) > 0 {
		b.WriteString(labelStyle.Render("New files"))
		b.WriteString("\n")
		for _, f := range s.Added {
			b.WriteString("  + " + f + "\n")
		}
		b.WriteString("\n")
	}

	if len(s.Deleted) > 0 {
		b.WriteString(labelStyle.Render("Deleted files"))
		b.WriteString("\n")
		for _, f := range s.Deleted {
			b.WriteString("  - " + f + "\n")
		}
		b.WriteString("\n")
	}

	return b.String()
}

func renderDirs(b *strings.Builder, dirs []DirChurn, termWidth int) {
	maxPath := 0
	maxChurn := 0
	for _, d := range dirs {
		if len(d.Path) > maxPath {
			maxPath = len(d.Path)
		}
		if c := d.Additions + d.Deletions; c > maxChurn {
			maxChurn = c
		}
	}
	if maxPath > 30 {
		maxPath = 30
	}

	barMax := termWidth - maxPath - 25
	if barMax < 10 {
		barMax = 10
	}

	addStyle := lipgloss.NewStyle().Foreground(ui.Green)
	delStyle := lipgloss.NewStyle().Foreground(ui.Red)
	dirStyle := lipgloss.NewStyle().Bold(true)

	limit := len(dirs)
	if limit > 15 {
		limit = 15
	}

	for _, d := range dirs[:limit] {
		path := runewidth.Truncate(d.Path, maxPath, "…")
		churn := d.Additions + d.Deletions

		bar := ""
		if churn > 0 {
			barLen := ui.ScaledBarWidth(churn, maxChurn, barMax)
			addBar := d.Additions * barLen / churn
			delBar := barLen - addBar
			bar = addStyle.Render(strings.Repeat("+", addBar)) +
				delStyle.Render(strings.Repeat("-", delBar))
		}

		fileWord := "files"
		if d.FileCount == 1 {
			fileWord = "file"
		}

		fmt.Fprintf(b, "  %-*s %4d %s  %s\n",
			maxPath, dirStyle.Render(path),
			churn, bar,
			ui.DimStyle.Render(fmt.Sprintf("%d %s", d.FileCount, fileWord)),
		)

		for _, f := range d.TopFiles {
			name := f.Path
			if idx := strings.LastIndex(name, "/"); idx >= 0 {
				name = name[idx+1:]
			}
			fmt.Fprintf(b, "    %s %s\n",
				ui.DimStyle.Render(name),
				ui.DimStyle.Render(f.TopAuthor),
			)
		}
	}

	if len(dirs) > limit {
		b.WriteString(ui.DimStyle.Render(fmt.Sprintf("  ... and %d more directories\n", len(dirs)-limit)))
	}
}

func renderFiles(b *strings.Builder, files []FileChurn, termWidth int) {
	maxPath := 0
	maxChurn := 0
	for _, f := range files {
		if len(f.Path) > maxPath {
			maxPath = len(f.Path)
		}
		if c := f.Additions + f.Deletions; c > maxChurn {
			maxChurn = c
		}
	}
	if maxPath > 40 {
		maxPath = 40
	}

	barMax := termWidth - maxPath - 30
	if barMax < 10 {
		barMax = 10
	}

	addStyle := lipgloss.NewStyle().Foreground(ui.Green)
	delStyle := lipgloss.NewStyle().Foreground(ui.Red)

	for _, f := range files {
		path := runewidth.Truncate(f.Path, maxPath, "…")
		churn := f.Additions + f.Deletions

		bar := ""
		if churn > 0 {
			barLen := ui.ScaledBarWidth(churn, maxChurn, barMax)
			addBar := f.Additions * barLen / churn
			delBar := barLen - addBar
			bar = addStyle.Render(strings.Repeat("+", addBar)) +
				delStyle.Render(strings.Repeat("-", delBar))
		}

		fmt.Fprintf(b, "  %-*s %4d %s  %s\n",
			maxPath, path,
			churn, bar,
			ui.DimStyle.Render(f.TopAuthor),
		)
	}
}

func renderAuthors(b *strings.Builder, authors []AuthorSummary, termWidth int) {
	maxName := 0
	for _, a := range authors {
		if len(a.Name) > maxName {
			maxName = len(a.Name)
		}
	}
	if maxName > 20 {
		maxName = 20
	}

	topCommits := authors[0].Commits

	barMax := termWidth - maxName - 25
	if barMax < 10 {
		barMax = 10
	}

	for i, a := range authors {
		name := runewidth.Truncate(a.Name, maxName, "…")

		bar := ""
		if len(authors) > 1 {
			barLen := ui.ScaledBarWidth(a.Commits, topCommits, barMax)
			palette := ui.CommitPalette[i%len(ui.CommitPalette)]
			barColor := ui.C(palette.Light, palette.Dark)
			bar = lipgloss.NewStyle().Foreground(barColor).Render(strings.Repeat("█", barLen))
		}

		fmt.Fprintf(b, "  %-*s %4d %s  %s\n",
			maxName,
			lipgloss.NewStyle().Foreground(ui.Green).Render(name),
			a.Commits, bar,
			ui.DimStyle.Render(a.TopPath),
		)
	}
}
