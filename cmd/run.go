/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/mohamedaminearari/gvm/internal/store"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run <version>",
	Short: "Run a specific Godot version without switching the active version",
	Long:  `Launches a specific version of Godot directly as a subprocess, forwarding any additional arguments to Godot.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		version := args[0]

		installed, err := store.IsVersionInstalled(version)
		if err != nil {
			return err
		}
		if !installed {
			fmt.Printf("version %s is not installed, run 'gvm install %s' to install it\n", version, version)
			return nil
		}

		binaryPath, err := store.FindBinary(version)
		if err != nil {
			return fmt.Errorf("could not find binary for version %s: %v", version, err)
		}

		godotArgs := args[1:]

		godotCmd := exec.Command(binaryPath, godotArgs...)

		godotCmd.Stdin = os.Stdin
		godotCmd.Stdout = os.Stdout
		godotCmd.Stderr = os.Stderr

		err = godotCmd.Run()
		if err != nil {
			return fmt.Errorf("godot exited with error: %v", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
