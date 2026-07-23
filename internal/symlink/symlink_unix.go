//go:build !windows

package symlink

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mohamedaminearari/gvm/internal/store"
)

const (
	symlinkFileName = "godot"
)

func Set(version string) error {
	installed, err := store.IsVersionInstalled(version)
	if err != nil {
		return err
	}
	if !installed {
		return fmt.Errorf("version %s is not installed, run 'gvm install %s' to install it", version, version)
	}

	binDir, err := store.BinDir()
	if err != nil {
		return err
	}

	binaryPath, err := store.FindBinary(version)
	if err != nil {
		return err
	}

	symlinkPath := filepath.Join(binDir, symlinkFileName)

	_, err = os.Lstat(symlinkPath)
	if err == nil {
		err = os.Remove(symlinkPath)
		if err != nil {
			return fmt.Errorf("failed to remove existing symlink: %v", err)
		}
	}

	err = os.Symlink(binaryPath, symlinkPath)
	if err != nil {
		return fmt.Errorf("failed to create symlink: %v", err)
	}
	return nil
}

func Current() (string, error) {
	binDir, err := store.BinDir()
	if err != nil {
		return "", err
	}

	symlinkPath := filepath.Join(binDir, symlinkFileName)

	_, err = os.Lstat(symlinkPath)
	if os.IsNotExist(err) {
		return "", nil
	}

	target, err := os.Readlink(symlinkPath)
	if err != nil {
		return "", fmt.Errorf("failed to read symlink: %v", err)
	}

	version := filepath.Base(filepath.Dir(target))
	return version, nil
}
