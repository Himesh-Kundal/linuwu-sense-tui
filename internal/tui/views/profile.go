package views

import (
	"fmt"
	"strings"

	"github.com/Himesh-Kundal/linuwu-tui/internal/hardware"
	"github.com/Himesh-Kundal/linuwu-tui/internal/tui/style"
	"github.com/charmbracelet/lipgloss"
)

func RenderProfile(caps hardware.Capabilities) string {
	if !caps.HasPlatformProfile {
		return lipgloss.NewStyle().Foreground(style.ColorWarning).Render("ACPI Platform Profile is not supported on this kernel/system.")
	}

	current, _ := hardware.GetPlatformProfile()
	choices, _ := hardware.GetPlatformProfileChoices()

	currentBox := style.MakeBox("Active Thermal Profile",
		fmt.Sprintf("Current Profile: %s\nSupported:       %s",
			style.StyleValue.Render(strings.ToUpper(current)),
			strings.Join(choices, ", "),
		),
	)

	presets := fmt.Sprintf(
		"[1] %s   (Quiet / low noise)\n"+
			"[2] %s (Standard daily driving)\n"+
			"[3] %s (High performance)\n"+
			"[4] %s   (Maximum power & fan speeds)",
		style.StyleValue.Render("quiet"),
		style.StyleValue.Render("balanced"),
		style.StyleValue.Render("performance"),
		style.StyleValue.Render("turbo"),
	)

	presetBox := style.MakeBox("Switch Thermal Profile", presets)

	hints := lipgloss.NewStyle().Foreground(style.ColorMuted).Render("Press [1], [2], [3], or [4] to switch thermal profile.")

	return lipgloss.JoinVertical(lipgloss.Left, currentBox, presetBox, "\n", hints)
}
