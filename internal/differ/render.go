package differ

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/josh-allan/go-git/internal/ui"
	"github.com/mattn/go-runewidth"
)

var (
	fileStyle    = lipgloss.NewStyle().Bold(true).Foreground(ui.Purple)
	hunkStyle    = lipgloss.NewStyle().Foreground(ui.DimGray)
	addedStyle   = lipgloss.NewStyle().Foreground(ui.Green)
	removedStyle = lipgloss.NewStyle().Foreground(ui.Red)
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
	colWidth := (termWidth - 3) / 2
	numWidth := 4

	for _, p := range pairs {
		leftLines := wrapSide(p.leftNum, p.leftLine, numWidth, colWidth)
		rightLines := wrapSide(p.rightNum, p.rightLine, numWidth, colWidth)

		rows := len(leftLines)
		if len(rightLines) > rows {
			rows = len(rightLines)
		}

		blank := strings.Repeat(" ", colWidth)
		for row := range rows {
			left := blank
			if row < len(leftLines) {
				left = leftLines[row]
			}
			right := blank
			if row < len(rightLines) {
				right = rightLines[row]
			}
			b.WriteString(left)
			b.WriteString(ui.SeparatorStyle.Render(" │ "))
			b.WriteString(right)
			b.WriteString("\n")
		}
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

func wrapSide(num int, line *Line, numWidth, colWidth int) []string {
	if line == nil {
		return []string{strings.Repeat(" ", colWidth)}
	}

	contentWidth := colWidth - numWidth - 1
	if contentWidth < 1 {
		contentWidth = 1
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

	content := ui.ExpandTabs(line.Content, 4)
	chunks := ui.WrapText(content, contentWidth)
	if len(chunks) == 0 {
		chunks = []string{""}
	}

	result := make([]string, len(chunks))
	for i, chunk := range chunks {
		prefix := strings.Repeat(" ", numWidth)
		if i == 0 {
			prefix = fmt.Sprintf("%*d", numWidth, num)
		}

		padding := contentWidth - runewidth.StringWidth(chunk)
		if padding < 0 {
			padding = 0
		}

		result[i] = ui.LineNumStyle.Render(prefix) + " " + style.Render(chunk) + strings.Repeat(" ", padding)
	}
	return result
}
