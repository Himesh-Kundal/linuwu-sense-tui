package views

import (
	"fmt"
	"strings"

	"github.com/Himesh-Kundal/linuwu-sense-tui/internal/hardware"
	"github.com/Himesh-Kundal/linuwu-sense-tui/internal/tui/style"
	"github.com/charmbracelet/lipgloss"
)

type profileEntry struct {
	name  string
	label string
	desc  string
}

var profiles = []profileEntry{
	{"quiet", "Quiet", "Minimum fan noise. Throttles CPU/GPU for silent operation."},
	{"balanced", "Balanced", "Standard daily-use profile. Balanced power and noise."},
	{"performance", "Performance", "High performance. Fans ramp up as needed."},
	{"turbo", "Turbo", "Maximum CPU/GPU boost. Fans at high speed."},
}

func RenderProfile(caps hardware.Capabilities, cursor int) string {
	if !caps.HasPlatformProfile {
		return style.StyleDanger.Render("  ⚠  ACPI Platform Profile is not supported on this system.")
	}

	current, _ := hardware.GetPlatformProfile()
	choices, _ := hardware.GetPlatformProfileChoices()

	// ── Current indicator ──
	currentLine := style.StyleLabel.Render("Active Profile") + "  " +
		lipgloss.NewStyle().Bold(true).Foreground(style.ColorPrimary).Render(strings.ToUpper(current))
	supportedLine := style.StyleMuted.Render("Supported:  " + strings.Join(choices, "  ·  "))

	infoPanel := style.Section("  Thermal Profile", currentLine+"\n"+supportedLine)

	// ── Profile list ──
	var lines []string
	for i, p := range profiles {
		isActive := p.name == current
		isCursor := i == cursor

		nameStr := "  " + p.label
		if isActive {
			nameStr = "✓ " + p.label
		}

		if isCursor {
			row := style.StyleRowSelected.Render(fmt.Sprintf("▶ %-14s", p.label))
			hint := style.StyleRowHint.Render("  └─ " + p.desc)
			lines = append(lines, row, hint, "")
		} else {
			var labelStyle lipgloss.Style
			if isActive {
				labelStyle = lipgloss.NewStyle().Bold(true).Foreground(style.ColorGreen)
			} else {
				labelStyle = lipgloss.NewStyle().Foreground(style.ColorFg)
			}
			_ = nameStr
			lines = append(lines, style.StyleRowNormal.Render(
				"  "+labelStyle.Render(fmt.Sprintf("%-14s", p.label))+style.StyleMuted.Render("  "+p.desc),
			), "")
		}
	}

	listPanel := style.SectionFocused("  Switch Profile", strings.Join(lines, "\n"))
	help := style.KeyHint("↑/↓", "Navigate") + "   " + style.KeyHint("Enter/Space", "Apply profile")

	return lipgloss.JoinVertical(lipgloss.Left, infoPanel, "", listPanel, "", help)
}
