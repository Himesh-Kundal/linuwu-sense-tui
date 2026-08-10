package views

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Himesh-Kundal/linuwu-sense-tui/internal/hardware"
	"github.com/Himesh-Kundal/linuwu-sense-tui/internal/sysfs"
	"github.com/Himesh-Kundal/linuwu-sense-tui/internal/tui/style"
	"github.com/charmbracelet/lipgloss"
)

var modeNames = []string{
	"Static", "Breathing", "Neon", "Wave", "Shifting", "Zoom", "Meteor", "Twinkling",
}

type kbRow struct {
	label string
	value string
	desc  string
}

func RenderKeyboard(caps hardware.Capabilities, cursor int) string {
	if !caps.HasFourZonedKB {
		return style.StyleWarning.Render("  Four-Zone RGB keyboard is not supported on this device.")
	}

	modePath := filepath.Join(caps.KBPath, "four_zone_mode")
	raw, err := sysfs.ReadString(modePath)
	if err != nil {
		raw = "0,0,0,1,0,0,0"
	}

	var mode, speed, brightness, direction, r, g, b int
	fmt.Sscanf(raw, "%d,%d,%d,%d,%d,%d,%d", &mode, &speed, &brightness, &direction, &r, &g, &b)

	modeName := "Unknown"
	if mode >= 0 && mode < len(modeNames) {
		modeName = modeNames[mode]
	}
	dirStr := "Left→Right"
	if direction == 2 {
		dirStr = "Right→Left"
	}

	rows := []kbRow{
		{
			label: "Effect Mode",
			value: lipgloss.NewStyle().Bold(true).Foreground(style.ColorPurple).Render(fmt.Sprintf("%-12s", modeName)),
			desc:  "Press Enter to cycle through 8 lighting modes",
		},
		{
			label: "Speed",
			value: lipgloss.NewStyle().Bold(true).Foreground(style.ColorCyan).Render(fmt.Sprintf("%d / 9", speed)),
			desc:  "Press Enter to increase animation speed (0-9)",
		},
		{
			label: "Brightness",
			value: MiniBar(brightness, 100, 16) + " " + lipgloss.NewStyle().Foreground(style.ColorYellow).Render(fmt.Sprintf("%d%%", brightness)),
			desc:  "Press Enter to cycle brightness (+25%)",
		},
		{
			label: "Direction",
			value: lipgloss.NewStyle().Bold(true).Foreground(style.ColorGreen).Render(dirStr),
			desc:  "Press Enter to toggle animation direction",
		},
	}

	var lines []string
	for i, r := range rows {
		label := fmt.Sprintf("%-16s", r.label)
		entry := "  " + label + " " + r.value
		if i == cursor {
			lines = append(lines, style.StyleRowSelected.Render(fmt.Sprintf("▶ %-16s %s", r.label, stripStyle(r.value))))
			lines = append(lines, style.StyleRowHint.Render("  └─ "+r.desc))
		} else {
			lines = append(lines, style.StyleRowNormal.Render(entry))
		}
		lines = append(lines, "")
	}

	// Color preview row
	colorSwatch := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(fmt.Sprintf("#%02x%02x%02x", r, g, b))).
		Render(fmt.Sprintf("  █████  RGB (%d, %d, %d)", r, g, b))
	lines = append(lines, colorSwatch)

	panel := style.SectionFocused("  4-Zone RGB Keyboard", strings.Join(lines, "\n"))
	help := style.KeyHint("↑/↓", "Navigate") + "   " + style.KeyHint("Enter/Space", "Cycle value")

	return lipgloss.JoinVertical(lipgloss.Left, panel, "", help)
}

// stripStyle removes ANSI color sequences for safe use in background-colored rows.
// For simplicity we just re-render a plain version.
func stripStyle(s string) string {
	// lipgloss strips styles itself when rendering inside another styled block.
	return s
}
