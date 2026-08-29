package ui

import (
	"strings"

	"github.com/mattn/go-runewidth"
)

func ExpandTabs(s string, tabWidth int) string {
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

// ScaledBarWidth returns a bar length for value given maxValue and available
// width. Small values produce small bars; bars only compress when maxValue
// exceeds the available width.
func ScaledBarWidth(value, maxValue, availableWidth int) int {
	if availableWidth <= 0 {
		return 0
	}
	scale := maxValue
	if scale < availableWidth {
		scale = availableWidth
	}
	bar := value * availableWidth / scale
	if bar > availableWidth {
		bar = availableWidth
	}
	if bar < 1 && value > 0 {
		bar = 1
	}
	return bar
}

func WrapText(s string, width int) []string {
	if runewidth.StringWidth(s) <= width {
		return []string{s}
	}

	var lines []string
	runes := []rune(s)
	for len(runes) > 0 {
		w := 0
		cut := 0
		for i, r := range runes {
			rw := runewidth.RuneWidth(r)
			if w+rw > width {
				break
			}
			w += rw
			cut = i + 1
		}
		if cut == 0 {
			cut = 1
		}
		lines = append(lines, string(runes[:cut]))
		runes = runes[cut:]
	}
	return lines
}
