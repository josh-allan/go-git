package ui

import (
	"image/color"
	"os"

	"charm.land/lipgloss/v2"
)

var (
	HasDarkBG = lipgloss.HasDarkBackground(os.Stdin, os.Stderr)
	LightDark = lipgloss.LightDark(HasDarkBG)
)

func C(light, dark string) color.Color {
	return LightDark(lipgloss.Color(light), lipgloss.Color(dark))
}

var (
	Purple   = C("#6C5CE7", "#A29BFE")
	Green    = C("#00B894", "#55EFC4")
	Red      = C("#D63031", "#FF7675")
	Orange   = C("#E17055", "#FAB1A0")
	Gray     = C("#636E72", "#636E72")
	DimGray  = C("#636E72", "#B2BEC3")
	Dim      = C("#B2BEC3", "#636E72")
	StatusFg = C("#FAFAFA", "#FAFAFA")
	StatusBg = C("#636E72", "#2D3436")
)

var (
	LineNumStyle   = lipgloss.NewStyle().Foreground(Gray)
	DimStyle       = lipgloss.NewStyle().Foreground(Dim)
	SeparatorStyle = lipgloss.NewStyle().Foreground(Gray)
)

var CommitPalette = []struct{ Light, Dark string }{
	{"#FAFAFA", "#DFE6E9"},
	{"#74B9FF", "#A3D8F4"},
	{"#A29BFE", "#C3B1E1"},
	{"#FDCB6E", "#F6E58D"},
	{"#55EFC4", "#BADFDB"},
	{"#FAB1A0", "#E8A598"},
	{"#81ECEC", "#C4FAF8"},
	{"#FF7675", "#FFB8B8"},
}
