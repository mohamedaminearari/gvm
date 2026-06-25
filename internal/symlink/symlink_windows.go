//go:build windows

package symlink

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mohamedaminearari/gvm/internal/store"
)

const (
	batFileName = "current.bat"
)

func Set(version string) error {
	installed, err := store.IsVersionInstalled(version)
	if err != nil {
		return err
	}
	if !installed {
		return fmt.Errorf("version %s is not installed, run 'gvm install %s' to install it", version, version)
	}

	gvmDir, err := store.GVMDir()
	if err != nil {
		return err
	}

	binaryPath, err := findBinary(version)
	if err != nil {
		return err
	}

	batContent := fmt.Sprintf("@echo off\n\"%s\" %%*\n", binaryPath)
	batPath := filepath.Join(gvmDir, batFileName)

	err = os.WriteFile(batPath, []byte(batContent), 0755)
	if err != nil {
		return fmt.Errorf("failed to write current.bat: %v", err)
	}

	return nil
}

func Current() (string, error) {
	gvmDir, err := store.GVMDir()
	if err != nil {
		return "", err
	}

	batPath := filepath.Join(gvmDir, batFileName)

	_, err = os.Stat(batPath)
	if os.IsNotExist(err) {
		return "", nil
	}

	content, err := os.ReadFile(batPath)
	if err != nil {
		return "", fmt.Errorf("failed to read current.bat: %v", err)
	}

	version, err := parseVersionFormat(string(content))
	if err != nil {
		return "", err
	}

	return version, nil
}

func findBinary(version string) (string, error) {
	versionDir, err := store.VersionDir(version)
	if err != nil {
		return "", err
	}

	entries, err := os.ReadDir(versionDir)
	if err != nil {
		return "", fmt.Errorf("failed to read version directory: %v", err)
	}

	for _, entry := range entries {
		name := entry.Name()
		if strings.HasSuffix(name, ".exe") && !strings.Contains(name, "console") {
			return filepath.Join(versionDir, name), nil
		}
	}

	return "", fmt.Errorf("could not find Godot executable in version directory for %s", version)
}

func parseVersionFormat(content string) (string, error) {
	marker := `versions\`
	idx := strings.Index(content, marker)
	if idx == -1 {
		return "", fmt.Errorf("could not parse version from current.bat")
	}

	rest := content[idx+len(marker):]
	end := strings.Index(rest, `\`)
	if end == -1 {
		return "", fmt.Errorf("could not parse version from current.bat")
	}

	return rest[:end], nil
}
