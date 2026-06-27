package store

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func GVMDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not find home directory: %v", err)
	}
	return filepath.Join(home, ".gvm"), nil
}

func TempDir() (string, error) {
	gvmDir, err := GVMDir()
	if err != nil {
		return "", fmt.Errorf("could not find gvm directory: %v", err)
	}
	return filepath.Join(gvmDir, "tmp"), nil
}

func VersionsDir() (string, error) {
	gvmDir, err := GVMDir()
	if err != nil {
		return "", fmt.Errorf("could not find gvm directory: %v", err)
	}
	return filepath.Join(gvmDir, "versions"), nil
}

func VersionDir(version string) (string, error) {
	versionsDir, err := VersionsDir()
	if err != nil {
		return "", fmt.Errorf("could not find versions directory: %v", err)
	}
	return filepath.Join(versionsDir, version), nil
}

func Init() error {
	gvmDir, err := GVMDir()
	if err != nil {
		return fmt.Errorf("could not find gvm directory: %v", err)
	}

	dirs := []string{
		gvmDir,
		filepath.Join(gvmDir, "tmp"),
		filepath.Join(gvmDir, "versions"),
		filepath.Join(gvmDir, "alias"),
	}

	for _, dir := range dirs {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %v", dir, err)
		}
	}

	return nil
}

func IsVersionInstalled(version string) (bool, error) {
	versionDir, err := VersionDir(version)
	if err != nil {
		return false, err
	}
	_, err = os.Stat(versionDir)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check version directory: %v", err)
	}
	return true, nil
}

func ListInstalledVersions() ([]string, error) {
	versionsDir, err := VersionsDir()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(versionsDir)
	if os.IsNotExist(err) {
		return []string{}, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to read versions directory: %v", err)
	}

	var versions []string
	for _, entry := range entries {
		if entry.IsDir() {
			versions = append(versions, entry.Name())
		}
	}

	return versions, nil
}

func extractFile(f *zip.File, destDir string) error {
	destPath := filepath.Join(destDir, f.Name)
	if !strings.HasPrefix(destPath, filepath.Clean(destDir)+string(os.PathSeparator)) {
		return fmt.Errorf("invalid file path in zip: %s", f.Name)
	}

	if f.FileInfo().IsDir() {
		return os.MkdirAll(destPath, f.Mode())
	}

	err := os.MkdirAll(filepath.Dir(destPath), 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	outFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer outFile.Close()

	rc, err := f.Open()
	if err != nil {
		return fmt.Errorf("failed to open file in zip: %v", err)
	}
	defer rc.Close()

	_, err = io.Copy(outFile, rc)
	return err
}

func ExtractAndSave(zipPath string, version string) error {
	versionDir, err := VersionDir(version)
	if err != nil {
		return err
	}

	err = os.MkdirAll(versionDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create version directory: %v", err)
	}

	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %v", err)
	}
	defer r.Close()

	for _, f := range r.File {
		err := extractFile(f, versionDir)
		if err != nil {
			return fmt.Errorf("failed to extract %s: %v", f.Name, err)
		}
	}

	return nil
}

func DeleteVersion(version string) error {
	versionDir, err := VersionDir(version)
	if err != nil {
		return err
	}

	installed, err := IsVersionInstalled(version)
	if err != nil {
		return err
	}

	if !installed {
		return fmt.Errorf("version %s is not installed", version)
	}

	err = os.RemoveAll(versionDir)
	if err != nil {
		return fmt.Errorf("failed to delete version %s: %v", version, err)
	}
	return nil
}

func FindBinary(version string) (string, error) {
	versionDir, err := VersionDir(version)
	if err != nil {
		return "", err
	}

	entries, err := os.ReadDir(versionDir)
	if err != nil {
		return "", fmt.Errorf("failed to read version directory: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()

		switch runtime.GOOS {
		case "windows":
			if strings.HasSuffix(name, ".exe") && !strings.Contains(name, "console") {
				return filepath.Join(versionDir, name), nil
			}
		default:
			info, err := entry.Info()
			if err != nil {
				continue
			}
			if info.Mode()&0111 != 0 {
				return filepath.Join(versionDir, name), nil
			}
		}
	}

	return "", fmt.Errorf("could not find Godot executable in version directory for %s", version)
}

func AliasDir() (string, error) {
	gvmDir, err := GVMDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(gvmDir, "alias"), nil
}

func SetAlias(name string, version string) error {
	aliasDir, err := AliasDir()
	if err != nil {
		return err
	}

	aliasPath := filepath.Join(aliasDir, name)

	err = os.WriteFile(aliasPath, []byte(version), 0644)
	if err != nil {
		return fmt.Errorf("failed to write alias %s: %v", name, err)
	}

	return nil
}

func GetAlias(name string) (string, error) {
	aliasDir, err := AliasDir()
	if err != nil {
		return "", err
	}

	aliasPath := filepath.Join(aliasDir, name)

	content, err := os.ReadFile(aliasPath)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("alias %s does not exist", name)
	} else if err != nil {
		return "", fmt.Errorf("failed to read alias %s: %v", name, err)
	}

	return strings.TrimSpace(string(content)), nil
}

func DeleteAlias(name string) error {
	aliasDir, err := AliasDir()
	if err != nil {
		return err
	}

	aliasPath := filepath.Join(aliasDir, name)

	_, err = os.Stat(aliasPath)
	if err != nil {
		return fmt.Errorf("alias %s does not exist", name)
	}

	err = os.Remove(aliasPath)
	if err != nil {
		return fmt.Errorf("failed to delete alias %s: %v", name, err)
	}
	return nil
}

func ListAliases() (map[string]string, error) {
	aliasDir, err := AliasDir()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(aliasDir)
	if os.IsNotExist(err) {
		return map[string]string{}, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to read alias directory: %v", err)
	}

	aliases := make(map[string]string)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		version, err := GetAlias(entry.Name())
		if err != nil {
			continue
		}
		aliases[entry.Name()] = version
	}
	return aliases, nil
}
