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

	fmt.Println("🛠️  Building linuwu_sense kernel module...")
	buildCmd := exec.Command("make", "all")
	buildCmd.Dir = tmpDir
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("module build failed (make all): %w", err)
	}

	fmt.Println("⚡ Installing linuwu_sense kernel module (requires sudo)...")
	installCmd := exec.Command("sudo", "-E", "make", "install")
	installCmd.Dir = tmpDir
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	installCmd.Stdin = os.Stdin

	if err := installCmd.Run(); err != nil {
		return fmt.Errorf("make install failed: %w", err)
	}

	fmt.Println("\n✓ Driver linuwu_sense successfully installed and loaded!")
	return nil
}
