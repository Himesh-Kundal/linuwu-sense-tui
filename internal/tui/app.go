package tui

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/Himesh-Kundal/linuwu-sense-tui/internal/hardware"
	"github.com/Himesh-Kundal/linuwu-sense-tui/internal/sysfs"
	"github.com/Himesh-Kundal/linuwu-sense-tui/internal/tui/style"
	"github.com/Himesh-Kundal/linuwu-sense-tui/internal/tui/views"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ── Constants ────────────────────────────────────────────────────────────────

const numTabs = 5

var tabNames = [numTabs]string{"Dashboard", "Fans", "Power", "Keyboard", "Profiles"}
var tabIcons = [numTabs]string{"◈", "⌬", "⚡", "⌨", "◎"}

// ── Model ────────────────────────────────────────────────────────────────────

type TickMsg time.Time

type Model struct {
	ActiveTab int
	Cursor    int // cursor within active tab
	Caps      hardware.Capabilities
	Sensors   hardware.SensorData
	Fans      views.FansModel
	Width     int
	Height    int
	StatusMsg string // transient status / error shown at the bottom
	KBState   [7]int // [mode, speed, brightness, direction, red, green, blue] — tracked in-memory
}

func NewModel() Model {
	caps := hardware.Detect()
	sensors := hardware.ReadSensors(caps)
	var fans views.FansModel
	fans.Init(caps)
	m := Model{
		ActiveTab: 0,
		Cursor:    0,
		Caps:      caps,
		Sensors:   sensors,
		Fans:      fans,
	}
	// Seed keyboard state from sysfs (one-time read; readback is unreliable)
	if caps.HasFourZonedKB {
		p := filepath.Join(caps.KBPath, "four_zone_mode")
		if raw, err := sysfs.ReadString(p); err == nil {
			fmt.Sscanf(raw, "%d,%d,%d,%d,%d,%d,%d", &m.KBState[0], &m.KBState[1], &m.KBState[2], &m.KBState[3], &m.KBState[4], &m.KBState[5], &m.KBState[6])
		}
		// Ensure direction is valid
		if m.KBState[3] <= 0 {
			m.KBState[3] = 1
		}
	}
	return m
}

func (m Model) Init() tea.Cmd {
	return tickCmd()
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// ── Update ────────────────────────────────────────────────────────────────────

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case TickMsg:
		m.Sensors = hardware.ReadSensors(m.Caps)
		return m, tickCmd()

	case tea.MouseMsg:
		if msg.Type == tea.MouseLeft {
			m.handleMouse(msg.X, msg.Y)
		}

	case tea.KeyMsg:
		switch msg.String() {
		// ── Global: quit
		case "q", "ctrl+c":
			return m, tea.Quit

		// ── Global: tab navigation  (Tab / Shift+Tab / F1-F5)
		case "tab":
			m.ActiveTab = (m.ActiveTab + 1) % numTabs
			m.Cursor = 0
		case "shift+tab":
			m.ActiveTab = (m.ActiveTab + numTabs - 1) % numTabs
			m.Cursor = 0
		case "f1":
			m.setTab(0)
		case "f2":
			m.setTab(1)
		case "f3":
			m.setTab(2)
		case "f4":
			m.setTab(3)
		case "f5":
			m.setTab(4)

		// ── Within tab: cursor movement
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			m.Cursor = clamp(m.Cursor+1, 0, m.maxCursor()-1)

		// ── Within tab: activate
		case " ", "enter":
			m.activate()

		// ── Tab-specific left/right & number shortcuts
		case "left", "h":
			if m.ActiveTab == 1 && m.Caps.HasFanSpeed {
				m.Fans.CPUSpeed = clamp(m.Fans.CPUSpeed-5, 0, 100)
				m.Fans.GPUSpeed = clamp(m.Fans.GPUSpeed-5, 0, 100)
				hardware.SetFanSpeed(m.Caps, m.Fans.CPUSpeed, m.Fans.GPUSpeed)
			} else if m.ActiveTab == 3 {
				m.activateKeyboardPrev()
			}
		case "right", "l":
			if m.ActiveTab == 1 && m.Caps.HasFanSpeed {
				m.Fans.CPUSpeed = clamp(m.Fans.CPUSpeed+5, 0, 100)
				m.Fans.GPUSpeed = clamp(m.Fans.GPUSpeed+5, 0, 100)
				hardware.SetFanSpeed(m.Caps, m.Fans.CPUSpeed, m.Fans.GPUSpeed)
			} else if m.ActiveTab == 3 {
				m.activateKeyboardNext()
			}
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			if m.ActiveTab == 1 && m.Caps.HasFanSpeed {
				val := (int(msg.String()[0] - '0')) * 10
				if msg.String() == "9" {
					val = 100
				}
				m.Fans.CPUSpeed, m.Fans.GPUSpeed = val, val
				hardware.SetFanSpeed(m.Caps, val, val)
			}
		}
	}
	return m, nil
}

func (m *Model) setTab(t int) {
	m.ActiveTab = t
	m.Cursor = 0
}

// maxCursor returns how many selectable rows the current tab has.
func (m *Model) maxCursor() int {
	switch m.ActiveTab {
	case 1: // Fans: nothing to navigate (use ←/→)
		return 1
	case 2: // Power
		return m.countPowerItems()
	case 3: // Keyboard
		return 7 // mode, speed, brightness, direction, Red, Green, Blue
	case 4: // Profiles
		choices, err := hardware.GetPlatformProfileChoices()
		if err != nil || len(choices) == 0 {
			return 1
		}
		return len(choices)
	default:
		return 1
	}
}

func (m *Model) countPowerItems() int {
	n := 0
	if m.Caps.HasBatteryLimiter {
		n++
	}
	if m.Caps.HasBatteryCalibration {
		n++
	}
	if m.Caps.HasUSBCharging {
		n++
	}
	if m.Caps.HasBacklightTimeout {
		n++
	}
	if m.Caps.HasLCDOverride {
		n++
	}
	if m.Caps.HasBootAnimationSound {
		n++
	}
	if n == 0 {
		return 1
	}
	return n
}

// activate fires the action for the cursor'd item in the active tab.
func (m *Model) activate() {
	switch m.ActiveTab {
	case 2: // Power
		m.activatePower()
	case 3: // Keyboard
		m.activateKeyboard()
	case 4: // Profiles
		if m.Caps.HasPlatformProfile {
			choices, err := hardware.GetPlatformProfileChoices()
			if err == nil && m.Cursor < len(choices) {
				hardware.SetPlatformProfile(choices[m.Cursor])
			}
		}
	}
}

func (m *Model) activatePower() {
	// Build ordered slice of enabled features to map cursor → attr.
	type item struct{ attr string }
	var items []item
	if m.Caps.HasBatteryLimiter {
		items = append(items, item{"battery_limiter"})
	}
	if m.Caps.HasBatteryCalibration {
		items = append(items, item{"battery_calibration"})
	}
	if m.Caps.HasUSBCharging {
		items = append(items, item{"_usb"}) // special
	}
	if m.Caps.HasBacklightTimeout {
		items = append(items, item{"backlight_timeout"})
	}
	if m.Caps.HasLCDOverride {
		items = append(items, item{"lcd_override"})
	}
	if m.Caps.HasBootAnimationSound {
		items = append(items, item{"boot_animation_sound"})
	}
	if m.Cursor >= len(items) {
		return
	}
	attr := items[m.Cursor].attr
	if attr == "_usb" {
		m.cycleUSBCharging()
	} else {
		m.toggleBattery(attr)
	}
}

func (m *Model) activateKeyboard() {
	m.stepKBParam(1)
}

func (m *Model) activateKeyboardNext() {
	m.stepKBParam(1)
}

func (m *Model) activateKeyboardPrev() {
	m.stepKBParam(-1)
}

func (m *Model) stepKBParam(delta int) {
	switch m.Cursor {
	case 0: // mode (0-7)
		m.KBState[0] = (m.KBState[0] + delta + 8) % 8
	case 1: // speed (0-9)
		m.KBState[1] = (m.KBState[1] + delta + 10) % 10
	case 2: // brightness (0-100 in steps of 25)
		next := m.KBState[2] + delta*25
		if next > 100 {
			next = 0
		} else if next < 0 {
			next = 100
		}
		m.KBState[2] = next
	case 3: // direction (toggle 1 <-> 2)
		if m.KBState[3] == 1 {
			m.KBState[3] = 2
		} else {
			m.KBState[3] = 1
		}
	case 4: // Red (0-255)
		next := m.KBState[4] + delta*15
		if next > 255 {
			next = 0
		} else if next < 0 {
			next = 255
		}
		m.KBState[4] = next
	case 5: // Green (0-255)
		next := m.KBState[5] + delta*15
		if next > 255 {
			next = 0
		} else if next < 0 {
			next = 255
		}
		m.KBState[5] = next
	case 6: // Blue (0-255)
		next := m.KBState[6] + delta*15
		if next > 255 {
			next = 0
		} else if next < 0 {
			next = 255
		}
		m.KBState[6] = next
	}

	// Ensure direction is always valid
	if m.KBState[3] <= 0 || m.KBState[3] > 2 {
		m.KBState[3] = 1
	}

	if err := m.writeKBState(); err != nil {
		m.StatusMsg = "⚠ KB write error: " + err.Error()
	} else {
		m.StatusMsg = ""
	}
}

// writeKBState writes the in-memory KBState to the sysfs four_zone_mode file.
func (m *Model) writeKBState() error {
	if !m.Caps.HasFourZonedKB {
		return fmt.Errorf("four-zone KB not supported")
	}
	p := filepath.Join(m.Caps.KBPath, "four_zone_mode")
	val := fmt.Sprintf("%d,%d,%d,%d,%d,%d,%d\n", m.KBState[0], m.KBState[1], m.KBState[2], m.KBState[3], m.KBState[4], m.KBState[5], m.KBState[6])
	if err := sysfs.WriteString(p, val); err != nil {
		return fmt.Errorf("write four_zone_mode %q: %w", val, err)
	}
	return nil
}

// handleMouse maps a click to an action.
func (m *Model) handleMouse(x, y int) {
	// Tab bar sits on row 1.
	if y == 1 {
		// Each tab is 14 chars wide approximately (icon + name + padding).
		tab := x / 14
		if tab < numTabs {
			m.setTab(tab)
		}
		return
	}
	// Content area begins at row 4; treat click as cursor selection + activate.
	if y >= 4 {
		row := y - 4
		if row >= 0 && row < m.maxCursor() {
			m.Cursor = row
			m.activate()
		}
	}
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func (m *Model) toggleBattery(attr string) {
	if !m.Caps.ModuleLoaded {
		return
	}
	p := filepath.Join(m.Caps.SensePath, attr)
	val, _ := sysfs.ReadInt(p)
	next := 0
	if val == 0 {
		next = 1
	}
	sysfs.WriteInt(p, next)
}

func (m *Model) cycleUSBCharging() {
	if !m.Caps.HasUSBCharging {
		return
	}
	p := filepath.Join(m.Caps.SensePath, "usb_charging")
	val, _ := sysfs.ReadInt(p)
	steps := []int{0, 10, 20, 30}
	next := steps[0]
	for i, s := range steps {
		if val == s {
			next = steps[(i+1)%len(steps)]
			break
		}
	}
	sysfs.WriteInt(p, next)
}

func clamp(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

// ── View ──────────────────────────────────────────────────────────────────────

func (m Model) View() string {
	// ── Header ──
	title := style.StyleHeaderTitle.Render("  LINUWU-SENSE TUI")
	model := style.StyleHeaderModel.Render(string(m.Caps.Model) + "  ")

	headerWidth := m.Width
	if headerWidth == 0 {
		headerWidth = 80
	}
	gapW := headerWidth - lipgloss.Width(title) - lipgloss.Width(model)
	if gapW < 0 {
		gapW = 0
	}
	gap := lipgloss.NewStyle().Width(gapW).Render("")
	header := lipgloss.JoinHorizontal(lipgloss.Top, title, gap, model)
	divider := lipgloss.NewStyle().Foreground(style.ColorBorder).Render(
		"─────────────────────────────────────────────────────────────────────────────────",
	)

	// ── Tab bar ──
	var tabs []string
	for i := 0; i < numTabs; i++ {
		label := fmt.Sprintf(" %s %s ", tabIcons[i], tabNames[i])
		if i == m.ActiveTab {
			tabs = append(tabs, style.StyleTabActive.Render(label))
		} else {
			tabs = append(tabs, style.StyleTabInactive.Render(label))
		}
	}
	tabBar := lipgloss.JoinHorizontal(lipgloss.Left, tabs...)

	// ── Body ──
	var body string
	switch m.ActiveTab {
	case 0:
		body = views.RenderDashboard(m.Caps, m.Sensors)
	case 1:
		body = views.RenderFans(&m.Fans, m.Caps, m.Sensors)
	case 2:
		body = views.RenderPower(m.Caps, m.Cursor)
	case 3:
		body = views.RenderKeyboard(m.Caps, m.Cursor, m.KBState)
	case 4:
		body = views.RenderProfile(m.Caps, m.Cursor)
	}

	// ── Footer ──
	hints := footerHints(m.ActiveTab)
	footer := style.StyleFooter.Render(hints)

	rows := []string{header, divider, tabBar, divider, "", body, ""}
	if m.StatusMsg != "" {
		rows = append(rows, lipgloss.NewStyle().Foreground(style.ColorRed).Bold(true).Render("  "+m.StatusMsg))
	}
	rows = append(rows, footer)

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func footerHints(tab int) string {
	kh := style.KeyHint
	base := kh("Tab", "Next Tab") + "  " + kh("Shift+Tab", "Prev Tab") + "  " + kh("F1-F5", "Jump Tab") + "  " + kh("q", "Quit")
	switch tab {
	case 1:
		return kh("←/→", "Fan speed ±5%") + "  " + kh("0", "Auto") + "  " + kh("9", "Max") + "  " + base
	case 2, 3, 4:
		return kh("↑/↓", "Navigate") + "  " + kh("Space/Enter", "Toggle/Select") + "  " + base
	}
	return base
}
