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

	colorHex := fmt.Sprintf("#%02x%02x%02x", r, g, b)
	colorSwatch := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(colorHex)).
		Render(fmt.Sprintf("████████  RGB(%d, %d, %d)", r, g, b))

	redBar := colorBar(r, 255, 14, style.ColorRed)
	greenBar := colorBar(g, 255, 14, style.ColorGreen)
	blueBar := colorBar(b, 255, 14, style.ColorPrimary)

	rows := []kbRow{
		{
			label: "Effect Mode",
			value: lipgloss.NewStyle().Bold(true).Foreground(style.ColorPurple).Render(fmt.Sprintf("%-12s", modeName)),
			desc:  "Press Enter or ←/→ to cycle lighting modes (Static, Breathing, Neon, Wave, etc.)",
		},
		{
			label: "Speed",
			value: lipgloss.NewStyle().Bold(true).Foreground(style.ColorCyan).Render(fmt.Sprintf("%d / 9", speed)),
			desc:  "Press Enter or ←/→ to adjust animation speed (0-9)",
		},
		{
			label: "Brightness",
			value: MiniBar(brightness, 100, 14) + " " + lipgloss.NewStyle().Foreground(style.ColorYellow).Render(fmt.Sprintf("%d%%", brightness)),
			desc:  "Press Enter or ←/→ to adjust brightness (0-100%)",
		},
		{
			label: "Direction",
			value: lipgloss.NewStyle().Bold(true).Foreground(style.ColorGreen).Render(dirStr),
			desc:  "Press Enter or ←/→ to toggle animation direction",
		},
		{
			label: "Red Channel",
			value: redBar + " " + lipgloss.NewStyle().Foreground(style.ColorRed).Render(fmt.Sprintf("%d", r)),
			desc:  "Press Enter or ←/→ to adjust Red intensity (0-255)",
		},
		{
			label: "Green Channel",
			value: greenBar + " " + lipgloss.NewStyle().Foreground(style.ColorGreen).Render(fmt.Sprintf("%d", g)),
			desc:  "Press Enter or ←/→ to adjust Green intensity (0-255)",
		},
		{
			label: "Blue Channel",
			value: blueBar + " " + lipgloss.NewStyle().Foreground(style.ColorPrimary).Render(fmt.Sprintf("%d", b)),
			desc:  "Press Enter or ←/→ to adjust Blue intensity (0-255)",
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

	// Live combined preview
	lines = append(lines, "  Active Color: "+colorSwatch)

	panel := style.SectionFocused("  4-Zone RGB Keyboard", strings.Join(lines, "\n"))
	help := style.KeyHint("↑/↓", "Navigate") + "   " + style.KeyHint("←/→ or Enter", "Adjust value")

	return lipgloss.JoinVertical(lipgloss.Left, panel, "", help)
}

func colorBar(val, max, width int, color lipgloss.Color) string {
	if val < 0 {
		val = 0
	}
	if val > max {
		val = max
	}
	fill := (val * width) / max
	empty := width - fill
	filled := lipgloss.NewStyle().Foreground(color).Render(strings.Repeat("█", fill))
	rest := lipgloss.NewStyle().Foreground(style.ColorMuted).Render(strings.Repeat("░", empty))
	return filled + rest
}

func stripStyle(s string) string {
	return s
}
