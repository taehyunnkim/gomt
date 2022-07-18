package tui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	padding int = 1
)

var appStyle = lipgloss.NewStyle().
	PaddingTop(padding).
	PaddingBottom(padding).
	PaddingLeft(padding).
	PaddingRight(padding)

var contentStyle = lipgloss.NewStyle()

var boxHeaderStyle = lipgloss.NewStyle().
	Bold(true).
	BorderStyle(lipgloss.RoundedBorder()).
	BorderBottom(false).
	BorderTop(true).
	BorderLeft(true).
 	BorderRight(true).
	BorderForeground(lipgloss.Color("#FAFAFA")).
	Foreground(lipgloss.Color("#FAFAFA")).
	PaddingLeft(padding).
	PaddingRight(padding)

var borderedBoxStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#FAFAFA")).
	Foreground(lipgloss.Color("#FAFAFA")).
	PaddingTop(padding).
	PaddingBottom(padding).
	PaddingLeft(padding).
	PaddingRight(padding)
	
	
var boxStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FAFAFA"))

	
var subHeaderStyle = lipgloss.NewStyle().
	Bold(true).
	Background(lipgloss.Color("#808080")).
	Foreground(lipgloss.Color("#FFFFFF")).
	PaddingLeft(padding).
	PaddingRight(padding).
	MarginBottom(1)
