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
	Use:   "run <version> [args...]",
	Short: "Run a specific Godot version without switching the active version",
	Long:  `Launches a specific version of Godot directly as a subprocess, forwarding any additional arguments to Godot.`,
	Args:  cobra.MinimumNArgs(1),
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
}
