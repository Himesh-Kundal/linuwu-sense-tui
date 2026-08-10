package views

import (
	"fmt"
	"strings"

	"github.com/Himesh-Kundal/linuwu-sense-tui/internal/hardware"
	"github.com/Himesh-Kundal/linuwu-sense-tui/internal/tui/style"
	"github.com/charmbracelet/lipgloss"
)

// profileDesc gives a human description for known ACPI profile names.
var profileDesc = map[string]string{
	"quiet":           "Minimum fan noise, throttles CPU/GPU for silent operation.",
	"low-power":       "Low power mode, throttles CPU/GPU for silent operation.",
	"balanced":        "Standard daily-use. Balanced power and noise.",
	"performance":     "High performance. Fans ramp up as needed.",
	"balanced-performance": "High performance with balanced fan curve.",
	"turbo":           "Maximum CPU/GPU boost. Fans at high speed.",
}

func RenderProfile(caps hardware.Capabilities, cursor int) string {
	if !caps.HasPlatformProfile {
		return style.StyleDanger.Render("  ⚠  ACPI Platform Profile is not supported on this system.")
	}

	current, _ := hardware.GetPlatformProfile()
	choices, _ := hardware.GetPlatformProfileChoices()

	if len(choices) == 0 {
		return style.StyleWarning.Render("  ⚠  No platform profiles found on this system.")
	}

	// ── Current indicator ──
	currentLine := style.StyleLabel.Render("Active Profile") + "  " +
		lipgloss.NewStyle().Bold(true).Foreground(style.ColorPrimary).Render(strings.ToUpper(current))
	supportedLine := style.StyleMuted.Render("Available:  " + strings.Join(choices, "  ·  "))

	infoPanel := style.Section("  Thermal Profile", currentLine+"\n"+supportedLine)

	// Clamp cursor to actual choices length
	if cursor >= len(choices) {
		cursor = len(choices) - 1
	}

	// ── Profile list from actual system choices ──
	var lines []string
	for i, name := range choices {
		isActive := name == current
		isCursor := i == cursor

		desc, ok := profileDesc[name]
		if !ok {
			desc = fmt.Sprintf("Apply '%s' thermal profile.", name)
		}

		label := name
		if isActive {
			label = "✓ " + name
		} else {
			label = "  " + name
		}

		if isCursor {
			row := style.StyleRowSelected.Render(fmt.Sprintf("▶ %-22s", name))
			hint := style.StyleRowHint.Render("  └─ " + desc)
			lines = append(lines, row, hint, "")
		} else {
			var nameStyle lipgloss.Style
			if isActive {
				nameStyle = lipgloss.NewStyle().Bold(true).Foreground(style.ColorGreen)
			} else {
				nameStyle = lipgloss.NewStyle().Foreground(style.ColorFg)
			}
			_ = label
			lines = append(lines, style.StyleRowNormal.Render(
				"  "+nameStyle.Render(fmt.Sprintf("%-22s", name))+style.StyleMuted.Render("  "+desc),
			), "")
		}
	}

	listPanel := style.SectionFocused("  Switch Profile", strings.Join(lines, "\n"))
	help := style.KeyHint("↑/↓", "Navigate") + "   " + style.KeyHint("Enter/Space", "Apply profile")

	return lipgloss.JoinVertical(lipgloss.Left, infoPanel, "", listPanel, "", help)
}
