package views

import (
	"fmt"
	"path/filepath"

	"github.com/Himesh-Kundal/linuwu-sense-tui/internal/hardware"
	"github.com/Himesh-Kundal/linuwu-sense-tui/internal/sysfs"
	"github.com/Himesh-Kundal/linuwu-sense-tui/internal/tui/style"
	"github.com/charmbracelet/lipgloss"
)

var ModeNames = []string{
	"Static", "Breathing", "Neon", "Wave", "Shifting", "Zoom", "Meteor", "Twinkling",
}

func RenderKeyboard(caps hardware.Capabilities) string {
	if !caps.HasFourZonedKB {
		return lipgloss.NewStyle().Foreground(style.ColorWarning).Render("Four-Zone RGB keyboard is not supported on this device.")
	}

	modePath := filepath.Join(caps.KBPath, "four_zone_mode")
	val, err := sysfs.ReadString(modePath)
	if err != nil {
		val = "Error reading state"
	}

	var mode, speed, brightness, direction, r, g, b int
	fmt.Sscanf(val, "%d,%d,%d,%d,%d,%d,%d", &mode, &speed, &brightness, &direction, &r, &g, &b)

	modeName := "Unknown"
	if mode >= 0 && mode < len(ModeNames) {
		modeName = ModeNames[mode]
	}

	dirStr := "Left-to-Right"
	if direction == 1 {
		dirStr = "Right-to-Left"
	}

	content := fmt.Sprintf(
		"Raw Sysfs: [%s]\n\n"+
			"%s %s (%d)\n"+
			"%s %d\n"+
			"%s %d%%\n"+
			"%s %s\n"+
			"%s (%d, %d, %d)",
		val,
		style.StyleLabel.Render("[M] Effect Mode:"), style.StyleValue.Render(modeName), mode,
		style.StyleLabel.Render("[S] Speed (0-9):"), speed,
		style.StyleLabel.Render("[B] Brightness:"), brightness,
		style.StyleLabel.Render("[D] Direction:"), dirStr,
		style.StyleLabel.Render("Color (R,G,B):"), r, g, b,
	)

	box := style.MakeBox("4-Zone RGB Keyboard", content)

	hints := lipgloss.NewStyle().Foreground(style.ColorMuted).Render(
		"Controls:\n" +
			"  [M] Cycle Mode  |  [S] Cycle Speed  |  [B] Adjust Brightness  |  [D] Toggle Direction",
	)

	return lipgloss.JoinVertical(lipgloss.Left, box, "\n", hints)
}
