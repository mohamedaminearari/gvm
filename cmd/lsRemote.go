package cmd

import (
	"fmt"
	"strings"

	"github.com/mohamedaminearari/gvm/internal/github"
	"github.com/mohamedaminearari/gvm/internal/platform"
	"github.com/spf13/cobra"
)

var includePrerelease bool

// lsRemoteCmd represents the lsRemote command
var lsRemoteCmd = &cobra.Command{
	Use:   "ls-remote",
	Short: "List available Godot versions from Github",
	Long:  `Fetches and displayes all available Godot versions for your current OS and architecture.`,
	RunE: func(cmd *cobra.Command, args []string) error {
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

		fmt.Println("Fetching available Godot versions from Github...")
		releases, err := github.FetchReleases(includePrerelease)
		if err != nil {
			return fmt.Errorf("failed to fetch releases: %v", err)
		}

		filtered := github.FilterByAssetSuffix(releases, suffix)
		if len(filtered) == 0 {
			fmt.Printf("No releases found for your platform: %s %s\n", info.OS, info.Arch)
			return nil
		}

		variant := "standard"
		if mono {
			variant = "mono"
		}

		fmt.Printf("Available Godot versions for %s %s (%s):\n", info.OS, info.Arch, variant)
		fmt.Println(strings.Repeat("-", 50))

		for _, release := range filtered {
			version := strings.TrimPrefix(release.TagName, "v")
			if release.Prerelease {
				fmt.Printf("  %s (pre-release)\n", version)
			} else {
				fmt.Printf("  %s\n", version)
			}
		}

		fmt.Printf("Total %d versions found\n", len(filtered))
		fmt.Println("Run 'gvm install <version>' to install a version")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(lsRemoteCmd)
	lsRemoteCmd.Flags().BoolVar(&includePrerelease, "pre", false, "Include pre-release versions")
	lsRemoteCmd.Flags().Bool("mono", false, "Show versions available for the Mono (C#) variant")
}
