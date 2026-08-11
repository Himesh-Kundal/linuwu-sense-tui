import sys

with open('internal/tui/app.go', 'r') as f:
    content = f.read()

target1 = """type Model struct {
	ActiveTab int
	Cursor    int // cursor within active tab
	Caps      hardware.Capabilities
	Sensors   hardware.SensorData
	Fans      views.FansModel
	Width     int
	Height    int
	StatusMsg string // transient status / error shown at the bottom
}

func NewModel() Model {
	caps := hardware.Detect()
	sensors := hardware.ReadSensors(caps)
	var fans views.FansModel
	fans.Init(caps)
	return Model{
		ActiveTab: 0,
		Cursor:    0,
		Caps:      caps,
		Sensors:   sensors,
		Fans:      fans,
	}
}"""

replacement1 = """type Model struct {
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
}"""

content = content.replace(target1, replacement1)

import re

# Match from stepKBParam down to the end of updateKBParamStep
pattern2 = re.compile(r'func \(m \*Model\) stepKBParam.*?return nil\n}', re.DOTALL)

replacement2 = """func (m *Model) stepKBParam(delta int) {
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
	val := fmt.Sprintf("%d,%d,%d,%d,%d,%d,%d\\n", m.KBState[0], m.KBState[1], m.KBState[2], m.KBState[3], m.KBState[4], m.KBState[5], m.KBState[6])
	if err := sysfs.WriteString(p, val); err != nil {
		return fmt.Errorf("write four_zone_mode %q: %w", val, err)
	}
	return nil
}"""

content = re.sub(pattern2, replacement2, content, count=1)

content = content.replace("body = views.RenderKeyboard(m.Caps, m.Cursor)", "body = views.RenderKeyboard(m.Caps, m.Cursor, m.KBState)")

with open('internal/tui/app.go', 'w') as f:
    f.write(content)
