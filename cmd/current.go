/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/mohamedaminearari/gvm/internal/symlink"
	"github.com/spf13/cobra"
)

// currentCmd represents the current command
var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show the currently active Godot version",
	Long:  `Displays the version of Godot that is currently active via the ~/.gvm symlink.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		version, err := symlink.Current()
		if err != nil {
			return fmt.Errorf("failed to read current version: %v", err)
		}

		if version == "" {
			fmt.Println("No active Godot version set.")
			fmt.Println("Run 'gvm use <version>' to activate one.")
			return nil
		}

		fmt.Printf("Current Godot version: %s\n", version)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(currentCmd)
}
