package views

import (
	"fmt"
	"strings"

	"github.com/Himesh-Kundal/linuwu-sense-tui/internal/hardware"
	"github.com/Himesh-Kundal/linuwu-sense-tui/internal/tui/style"
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
		"[A] %s   (Quiet / low noise)\n"+
			"[B] %s (Standard daily driving)\n"+
			"[C] %s (High performance)\n"+
			"[D] %s   (Maximum power & fan speeds)",
		style.StyleValue.Render("quiet"),
		style.StyleValue.Render("balanced"),
		style.StyleValue.Render("performance"),
		style.StyleValue.Render("turbo"),
	)

	presetBox := style.MakeBox("Switch Thermal Profile", presets)

	hints := lipgloss.NewStyle().Foreground(style.ColorMuted).Render("Press [A] Quiet, [B] Balanced, [C] Performance, [D] Turbo to switch profile.")

	return lipgloss.JoinVertical(lipgloss.Left, currentBox, presetBox, "\n", hints)
}
