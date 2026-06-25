//go:build !windows

package symlink

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mohamedaminearari/gvm/internal/store"
)

func Set(version string) error {
	versionDir, err := store.VersionDir(version)
	if err != nil {
		return err
	}

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

	symlinkPath := filepath.Join(gvmDir, "current")

	_, err = os.Lstat(symlinkPath)
	if err == nil {
		err = os.Remove(symlinkPath)
		if err != nil {
			return fmt.Errorf("failed to remove existing symlink: %v", err)
		}
	}

	err = os.Symlink(versionDir, symlinkPath)
	if err != nil {
		return fmt.Errorf("failed to create symlink: %v", err)
	}
	return nil
}

func Current() (string, error) {
	gvmDir, err := store.GVMDir()
	if err != nil {
		return "", err
	}

	symlinkPath := filepath.Join(gvmDir, "current")

	_, err = os.Lstat(symlinkPath)
	if os.IsNotExist(err) {
		return "", nil
	}

	target, err := os.Readlink(symlinkPath)
	if err != nil {
		return "", fmt.Errorf("failed to read symlink: %v", err)
	}

	version := filepath.Base(target)
	return version, nil
}
