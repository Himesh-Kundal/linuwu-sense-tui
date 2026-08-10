package sysfs

import (
	"os"
	"strconv"
	"strings"
)

// Exists checks if a sysfs path exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// ReadString reads a string value from a sysfs path.
func ReadString(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// WriteString writes a string value to a sysfs path.
// Uses O_WRONLY so it works correctly with kernel sysfs/procfs virtual files.
func WriteString(path string, val string) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(val + "\n")
	return err
}

// ReadInt reads an integer value from a sysfs path.
func ReadInt(path string) (int, error) {
	s, err := ReadString(path)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(s)
}

// WriteInt writes an integer value to a sysfs path.
func WriteInt(path string, val int) error {
	return WriteString(path, strconv.Itoa(val))
}
