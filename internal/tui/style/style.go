package style

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	ColorPrimary   = lipgloss.Color("#00F0FF")
	ColorSecondary = lipgloss.Color("#7000FF")
	ColorAccent    = lipgloss.Color("#FF0055")
	ColorBgDark    = lipgloss.Color("#12131C")
	ColorFgLight   = lipgloss.Color("#E2E8F0")
	ColorMuted     = lipgloss.Color("#64748B")
	ColorSuccess   = lipgloss.Color("#10B981")
	ColorWarning   = lipgloss.Color("#F59E0B")
	ColorDanger    = lipgloss.Color("#EF4444")

	StyleHeaderTitle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorPrimary).
				Padding(0, 1)

	StyleHeaderModel = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorAccent).
				Align(lipgloss.Right)

	StyleTabActive = lipgloss.NewStyle().
			Bold(true).
			Background(ColorPrimary).
			Foreground(ColorBgDark).
			Padding(0, 2)

	StyleTabInactive = lipgloss.NewStyle().
				Foreground(ColorMuted).
				Padding(0, 2)

	StyleBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorPrimary).
			Padding(1, 2)

	StyleLabel = lipgloss.NewStyle().
			Foreground(ColorFgLight).
			Width(22)

	StyleValue = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary)

	StyleStatusKey = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Italic(true)
)

func MakeBox(title, content string) string {
	titledHeader := lipgloss.NewStyle().Bold(true).Foreground(ColorPrimary).Render(fmt.Sprintf(" %s ", title))
	boxed := StyleBox.Render(content)
	return titledHeader + "\n" + boxed
}
