package cmd

import (
	"fmt"

	"github.com/mohamedaminearari/gvm/internal/store"
	"github.com/mohamedaminearari/gvm/internal/symlink"
	"github.com/spf13/cobra"
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use <version>",
	Short: "Switch to a specific Godot version",
	Long:  `Updates the active version of Godot by updating the symlink (or wrapper for Windows) in ~/.gvm/bin.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		version := args[0]

		content, err := store.GetAlias(version)
		if err == nil {
			version = content
		}

		err = symlink.Set(version)
		if err != nil {
			return fmt.Errorf("failed to switch version: %v", err)
		}

		fmt.Printf("Now using Godot %s\n", version)
		fmt.Println("Make sure ~/.gvm/bin is in your PATH to use the 'godot' command.")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(useCmd)
}
