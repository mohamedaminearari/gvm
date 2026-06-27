package cmd

import (
	"fmt"
	"strings"

	"github.com/mohamedaminearari/gvm/internal/store"
	"github.com/spf13/cobra"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all locally installed Godot versions",
	Long:  `Displays all versions of Godot that are currently installed in ~/.gvm/versions/.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		versions, err := store.ListInstalledVersions()
		if err != nil {
			return fmt.Errorf("failed to list installed versions: %v", err)
		}

		if len(versions) == 0 {
			fmt.Println("No Godot versions installed.")
			fmt.Println("Run 'gvm ls-remote' to see available versions.")
			return nil
		}

		fmt.Println("Installed Godot versions:")
		fmt.Println(strings.Repeat("-", 50))

		for _, version := range versions {
			fmt.Printf("  %s\n", version)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
