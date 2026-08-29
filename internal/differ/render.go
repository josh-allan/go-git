package differ

import (
	"fmt"
	"os"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/mattn/go-runewidth"
)

var (
	hasDarkBG = lipgloss.HasDarkBackground(os.Stdin, os.Stderr)
	lightDark = lipgloss.LightDark(hasDarkBG)

	fileStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lightDark(lipgloss.Color("#6C5CE7"), lipgloss.Color("#A29BFE")))

	hunkStyle = lipgloss.NewStyle().
			Foreground(lightDark(lipgloss.Color("#636E72"), lipgloss.Color("#B2BEC3")))

	addedStyle = lipgloss.NewStyle().
			Foreground(lightDark(lipgloss.Color("#00B894"), lipgloss.Color("#55EFC4")))

	removedStyle = lipgloss.NewStyle().
			Foreground(lightDark(lipgloss.Color("#D63031"), lipgloss.Color("#FF7675")))

	lineNumStyle = lipgloss.NewStyle().
			Foreground(lightDark(lipgloss.Color("#636E72"), lipgloss.Color("#636E72")))

	separatorStyle = lipgloss.NewStyle().
			Foreground(lightDark(lipgloss.Color("#636E72"), lipgloss.Color("#636E72")))
)

type sidePair struct {
	leftNum   int
	leftLine  *Line
	rightNum  int
	rightLine *Line
}

func Render(files []DiffFile, termWidth int) string {
	if termWidth < 40 {
		termWidth = 80
	}

	var b strings.Builder

	for i, file := range files {
		if i > 0 {
			b.WriteString("\n")
		}
		renderFile(&b, file, termWidth)
	}

	return b.String()
}

func renderFile(b *strings.Builder, file DiffFile, termWidth int) {
	name := strings.TrimPrefix(file.NewName, "b/")

	b.WriteString(fileStyle.Render("── " + name + " "))
	remaining := termWidth - len(name) - 4
	if remaining > 0 {
		b.WriteString(fileStyle.Render(strings.Repeat("─", remaining)))
	}
	b.WriteString("\n")

	for _, hunk := range file.Hunks {
		renderHunk(b, hunk, termWidth)
	}
}

func renderHunk(b *strings.Builder, hunk Hunk, termWidth int) {
	header := fmt.Sprintf("@@ -%d,%d +%d,%d @@", hunk.OldStart, hunk.OldCount, hunk.NewStart, hunk.NewCount)
	b.WriteString(hunkStyle.Render(header))
	b.WriteString("\n")

	pairs := alignLines(hunk)
	colWidth := (termWidth - 3) / 2 // 3 for " │ " separator
	numWidth := 4

	for _, p := range pairs {
		left := formatSide(p.leftNum, p.leftLine, numWidth, colWidth)
		right := formatSide(p.rightNum, p.rightLine, numWidth, colWidth)
		b.WriteString(left)
		b.WriteString(separatorStyle.Render(" │ "))
		b.WriteString(right)
		b.WriteString("\n")
	}
}

func alignLines(hunk Hunk) []sidePair {
	var pairs []sidePair
	oldNum := hunk.OldStart
	newNum := hunk.NewStart

	i := 0
	for i < len(hunk.Lines) {
		line := hunk.Lines[i]

		switch line.Type {
		case LineContext:
			pairs = append(pairs, sidePair{
				leftNum: oldNum, leftLine: &hunk.Lines[i],
				rightNum: newNum, rightLine: &hunk.Lines[i],
			})
			oldNum++
			newNum++
			i++

		case LineRemoved:
			// collect consecutive removed lines, then pair with following added lines
			var removed []int
			for i < len(hunk.Lines) && hunk.Lines[i].Type == LineRemoved {
				removed = append(removed, i)
				i++
			}
			var added []int
			for i < len(hunk.Lines) && hunk.Lines[i].Type == LineAdded {
				added = append(added, i)
				i++
			}

			maxLen := len(removed)
			if len(added) > maxLen {
				maxLen = len(added)
			}

			for j := range maxLen {
				p := sidePair{}
				if j < len(removed) {
					p.leftNum = oldNum
					p.leftLine = &hunk.Lines[removed[j]]
					oldNum++
				}
				if j < len(added) {
					p.rightNum = newNum
					p.rightLine = &hunk.Lines[added[j]]
					newNum++
				}
				pairs = append(pairs, p)
			}

		case LineAdded:
			pairs = append(pairs, sidePair{
				rightNum: newNum, rightLine: &hunk.Lines[i],
			})
			newNum++
			i++
		}
	}

	return pairs
}

func formatSide(num int, line *Line, numWidth, colWidth int) string {
	if line == nil {
		return strings.Repeat(" ", colWidth)
	}

	numStr := fmt.Sprintf("%*d", numWidth, num)
	content := expandTabs(line.Content, 4)

	contentWidth := colWidth - numWidth - 1
	content = runewidth.Truncate(content, contentWidth, "")

	padding := contentWidth - runewidth.StringWidth(content)
	if padding < 0 {
		padding = 0
	}

	var style lipgloss.Style
	switch line.Type {
	case LineAdded:
		style = addedStyle
	case LineRemoved:
		style = removedStyle
	default:
		style = lipgloss.NewStyle()
	}

	return lineNumStyle.Render(numStr) + " " + style.Render(content) + strings.Repeat(" ", padding)
}

func expandTabs(s string, tabWidth int) string {
	var b strings.Builder
	col := 0
	for _, r := range s {
		if r == '\t' {
			spaces := tabWidth - (col % tabWidth)
			b.WriteString(strings.Repeat(" ", spaces))
			col += spaces
		} else {
			b.WriteRune(r)
			col++
		}
	}
	return b.String()
}
