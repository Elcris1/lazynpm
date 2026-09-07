package ui

import (
	"strings"
	"github.com/charmbracelet/lipgloss"
)

const (
	mainColor =  "#874BFD"
)

func RenderBox(title, content string, innerWidth int, isMainColor bool) string {
	

	// Styles
	borderStyle := lipgloss.NewStyle()
	titleStyle := lipgloss.NewStyle().Bold(true)

	color := lipgloss.Color(mainColor)
	if isMainColor {
		// Colors
		borderStyle = borderStyle.Foreground(color)
		titleStyle = titleStyle.Foreground(color)
	}

	// Build the top line: ╭─ TITLE ─────────╮
	styledTitle := titleStyle.Render(" " + title + " ")
	titleWidth := lipgloss.Width(styledTitle)

	// Calculate right fill width (innerWidth minus left corner, 1 dash, and title)
	rightFillWidth := innerWidth - 1 - titleWidth
	if rightFillWidth < 0 {
		rightFillWidth = 0
	}

	topLine := borderStyle.Render("╭─") + styledTitle + borderStyle.Render(strings.Repeat("─", rightFillWidth)+"╮")

	// Build middle content box (left, right, bottom borders only)
	bodyStyle := lipgloss.NewStyle().
		Border(lipgloss.Border{
			Left:        "│",
			Right:       "│",
			Bottom:      "─",
			BottomLeft:  "╰",
			BottomRight: "╯",
		}, false, true, true, true). // Top border disabled
		Width(innerWidth).
		Padding(1, 1)

	if isMainColor {
		bodyStyle = bodyStyle.BorderForeground(color)
	}

	renderedBody := bodyStyle.Render(content)

	// Stack top header line and body
	return lipgloss.JoinVertical(lipgloss.Left, topLine, renderedBody)
}

func TwoColumn(leftCol *string, rightCol *string) string {

	// Join side-by-side aligned at the top
	twoColumns := lipgloss.JoinHorizontal(lipgloss.Top, *leftCol, "  ", *rightCol)

	return twoColumns
}

func RenderText(content string, isMainOption bool) string {

	if !isMainOption {
		return content + "\n"
	}

	var titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		PaddingLeft(2).
		PaddingRight(2)

	return titleStyle.Render(content) + "\n"
}