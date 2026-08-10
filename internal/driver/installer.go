package driver

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed files/* files/src/*
var driverFiles embed.FS

func Install() error {
	tmpDir, err := os.MkdirTemp("", "linuwu-sense-driver-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	err = fs.WalkDir(driverFiles, "files", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relPath, _ := filepath.Rel("files", path)
		targetPath := filepath.Join(tmpDir, relPath)

		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		data, err := driverFiles.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(targetPath, data, 0644)
	})

	if err != nil {
		return fmt.Errorf("failed to extract embedded driver files: %w", err)
	}

	cmd := exec.Command("sudo", "make", "install")
	cmd.Dir = tmpDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	fmt.Println("⚡ Building and installing linuwu_sense kernel module...")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("make install failed: %w", err)
	}

	fmt.Println("✓ Driver linuwu_sense successfully installed!")
	return nil
}
