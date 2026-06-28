package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "gvm",
	Version: "v1.0.0",
	Short:   "Godot Version Manager - A CLI tool to manage multiple Godot engine versions",
	Long: `GVM is a cross-platform CLI tool for managing multiple versions of the Godot game engine.

It allows you to install, switch, and manage Godot versions on Windows, Linux, and macOS, modeled after nvm (Node Version Manager)

Examples:
	gvm ls-remote                   List all available Godot versions for your OS
	gvm install 4.3-stable          Install a specific Godot version
	gvm use 4.3-stable              Switch to an installed version
	gvm current                     Show the currently active version
	gvm ls                          List all locally installed versions
	gvm uninstall 4.3-stable        Remove an installed version
	gvm alias myproject 4.3-stable  Create a named alias for a version
	gvm which 4.3-stable            Show the path to a version's executable
	gvm run 4.3-stable              Launch a specific version directly`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
