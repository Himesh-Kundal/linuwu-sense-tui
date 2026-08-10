package style

import "github.com/charmbracelet/lipgloss"

// ── Palette ─────────────────────────────────────────────────────────────────
var (
	ColorBg       = lipgloss.Color("#0D0E14")
	ColorSurface  = lipgloss.Color("#13141D")
	ColorBorder   = lipgloss.Color("#1E2030")
	ColorPrimary  = lipgloss.Color("#82AAFF")
	ColorCyan     = lipgloss.Color("#89DDFF")
	ColorGreen    = lipgloss.Color("#C3E88D")
	ColorYellow   = lipgloss.Color("#FFCB6B")
	ColorRed      = lipgloss.Color("#F07178")
	ColorPurple   = lipgloss.Color("#C792EA")
	ColorFg       = lipgloss.Color("#EEFFFF")
	ColorSubtle   = lipgloss.Color("#545C7E")
	ColorMuted    = lipgloss.Color("#3B4261")
	ColorAccent   = lipgloss.Color("#FF9E64")

	ColorSuccess = ColorGreen
	ColorWarning = ColorYellow
	ColorDanger  = ColorRed
)

// ── Typography ───────────────────────────────────────────────────────────────
var (
	StyleTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary)

	StyleSubtitle = lipgloss.NewStyle().
			Foreground(ColorSubtle)

	StyleLabel = lipgloss.NewStyle().
			Foreground(ColorFg).
			Width(26)

	StyleValue = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorCyan)

	StyleMuted = lipgloss.NewStyle().
			Foreground(ColorSubtle)

	StyleSuccess = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorGreen)

	StyleWarning = lipgloss.NewStyle().
			Foreground(ColorYellow)

	StyleDanger = lipgloss.NewStyle().
			Foreground(ColorRed)

	StyleKey = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorAccent)
)

// ── Header ───────────────────────────────────────────────────────────────────
var (
	StyleHeaderTitle = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorPrimary).
				Padding(0, 2)

	StyleHeaderModel = lipgloss.NewStyle().
				Foreground(ColorSubtle).
				Padding(0, 2)
)

// ── Tabs ─────────────────────────────────────────────────────────────────────
var (
	StyleTabActive = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorBg).
			Background(ColorPrimary).
			Padding(0, 2)

	StyleTabInactive = lipgloss.NewStyle().
				Foreground(ColorSubtle).
				Padding(0, 2)
)

// ── Panels / Boxes ───────────────────────────────────────────────────────────
var (
	StylePanel = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(0, 1)

	StylePanelFocused = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorPrimary).
				Padding(0, 1)
)

// ── List rows ────────────────────────────────────────────────────────────────
var (
	StyleRowSelected = lipgloss.NewStyle().
				Bold(true).
				Foreground(ColorBg).
				Background(ColorPrimary).
				Padding(0, 1)

	StyleRowNormal = lipgloss.NewStyle().
			Foreground(ColorFg).
			Padding(0, 1)

	StyleRowHint = lipgloss.NewStyle().
			Foreground(ColorSubtle).
			Padding(0, 1)
)

// ── Footer ───────────────────────────────────────────────────────────────────
var StyleFooter = lipgloss.NewStyle().
	Foreground(ColorSubtle).
	BorderTop(true).
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(ColorBorder).
	Padding(0, 1)

// ── Helpers ──────────────────────────────────────────────────────────────────

// Section renders a titled panel.
func Section(title, content string) string {
	hdr := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorPrimary).
		Padding(0, 1).
		Render(title)
	body := StylePanel.Render(content)
	return lipgloss.JoinVertical(lipgloss.Left, hdr, body)
}

// SectionFocused renders a titled panel with a highlighted border.
func SectionFocused(title, content string) string {
	hdr := lipgloss.NewStyle().
		Bold(true).
		Foreground(ColorBg).
		Background(ColorPrimary).
		Padding(0, 1).
		Render("  " + title + "  ")
	body := StylePanelFocused.Render(content)
	return lipgloss.JoinVertical(lipgloss.Left, hdr, body)
}

// OnOff renders a coloured ON/OFF badge.
func OnOff(active bool) string {
	if active {
		return lipgloss.NewStyle().Bold(true).Foreground(ColorGreen).Render("● ON ")
	}
	return lipgloss.NewStyle().Foreground(ColorSubtle).Render("○ OFF")
}

// KeyHint renders a key hint like "[Space] Toggle".
func KeyHint(key, desc string) string {
	k := lipgloss.NewStyle().Bold(true).Foreground(ColorAccent).Render("[" + key + "]")
	d := lipgloss.NewStyle().Foreground(ColorSubtle).Render(" " + desc)
	return k + d
}
