/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mohamedaminearari/gvm/internal/github"
	"github.com/mohamedaminearari/gvm/internal/platform"
	"github.com/mohamedaminearari/gvm/internal/store"
	"github.com/spf13/cobra"
)

const (
	godotBuildsURL = "https://github.com/godotengine/godot-builds/releases/download"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install <version>",
	Short: "Install a specific Godot version",
	Long:  `Downloads and installs a specific version of Godot into ~/.gvm/versions/.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		version := args[0]

		info := platform.Detect()
		if !info.IsSupported() {
			return fmt.Errorf("unsupported platform: %s  %s", info.OS, info.Arch)
		}

		mono, _ := cmd.Flags().GetBool("mono")
		var suffix string
		if mono {
			suffix = info.GodotMonoAssetSuffix()
		} else {
			suffix = info.GodotAssetSuffix()
		}

		if version == "latest" {
			fmt.Println("Fetching latest stable version from Github...")
			releases, err := github.FetchReleases(false)
			if err != nil {
				return fmt.Errorf("failed to fetch feleases: %v", err)
			}
			filtered := github.FilterByAssetSuffix(releases, suffix)
			if len(filtered) == 0 {
				return fmt.Errorf("no releases found for your platform")
			}

			version = filtered[0].TagName
			fmt.Printf("Latest version is %s\n", version)
		}

		installed, err := store.IsVersionInstalled(version)
		if err != nil {
			return err
		}
		if installed {
			fmt.Printf("Godot %s is already installed. Use 'gvm use %s' to switch to it.\n", version, version)
			return nil
		}

		fileName := fmt.Sprintf("Godot_v%s_%s", version, suffix)
		downloadURL := fmt.Sprintf("%s/%s/%s", godotBuildsURL, version, fileName)

		err = store.Init()
		if err != nil {
			return fmt.Errorf("failed to initialize gvm directory: %v", err)
		}

		tempDir, err := store.TempDir()
		if err != nil {
			return err
		}

		zipPath := filepath.Join(tempDir, fileName)

		fmt.Printf("Downloading Godot %s for %s %s...\n", version, info.OS, info.Arch)
		fmt.Printf("URL: %s\n", downloadURL)

		err = downloadFile(downloadURL, zipPath)
		if err != nil {
			return fmt.Errorf("download failed: %v", err)
		}

		fmt.Println("Download complete. Extracting...")

		err = store.ExtractAndSave(zipPath, version)
		if err != nil {
			return fmt.Errorf("extraction failed: %v", err)
		}

		err = os.Remove(zipPath)
		if err != nil {
			fmt.Printf("Warning: could not remove temp file %s: %v\n", zipPath, err)
		}

		fmt.Printf("Godot %s installed successfully.\n", version)
		fmt.Printf("Run 'gvm use %s' to activate it.\n", version)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().Bool("mono", false, "Install the Mono (C#) variant of Godot")
}

func downloadFile(url string, destPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to reach download url: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return fmt.Errorf("Github API rate limit exceeded, please wait a moment and try again")
	}

	if resp.StatusCode == http.StatusForbidden {
		return fmt.Errorf("Version not found, use 'gvm ls-remote' to see available versions")
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Github API returned unexpected status: %s", resp.Status)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}
	return nil
}
