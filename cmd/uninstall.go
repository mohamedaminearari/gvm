/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/mohamedaminearari/gvm/internal/store"
	"github.com/spf13/cobra"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall <version>",
	Short: "Uninstall a specific Godot version ",
	Long:  `Removes an installed version of Godot from ~/.gvm/versions/.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		version := args[0]

		installed, err := store.IsVersionInstalled(version)
		if err != nil {
			return err
		}
		if !installed {
			return fmt.Errorf("version %s is not installed, use 'gvm ls' to see installed versions", version)
		}

		err = store.DeleteVersion(version)
		if err != nil {
			return fmt.Errorf("failed to uninstall version %s: %v", version, err)
		}

		fmt.Printf("Godot %s uninstalled successfully.\n", version)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
