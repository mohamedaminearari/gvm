/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/mohamedaminearari/gvm/internal/store"
	"github.com/spf13/cobra"
)

// whichCmd represents the which command
var whichCmd = &cobra.Command{
	Use:   "which <version>",
	Short: "Show the path to a specific Godot Version's executable",
	Long:  `Prints the full path to the Godot executable for the given version.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		version := args[0]

		content, err := store.GetAlias(version)
		if err == nil {
			version = content
		}

		installed, err := store.IsVersionInstalled(version)
		if err != nil {
			return err
		}

		if !installed {
			return fmt.Errorf("version %s is not installed, run 'gvm install %s' to install it", version, version)
		}

		binaryPath, err := store.FindBinary(version)
		if err != nil {
			return fmt.Errorf("could not find binary for version %s: %v", version, err)
		}

		fmt.Println(binaryPath)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(whichCmd)
}
