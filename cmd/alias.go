package cmd

import (
	"fmt"
	"strings"

	"github.com/mohamedaminearari/gvm/internal/store"
	"github.com/spf13/cobra"
)

// aliasCmd represents the alias command
var aliasCmd = &cobra.Command{
	Use:   "alias <name> <version>",
	Short: "Create, list, or delete aliases for Godot versions",
	Long:  `Manage named aliases for installed versions of Godot.`,
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		delete, _ := cmd.Flags().GetBool("delete")

		if delete {
			if len(args) != 1 {
				return fmt.Errorf("--delete requires an alias name")
			}
			name := args[0]
			err := store.DeleteAlias(name)
			if err != nil {
				return err
			}
			fmt.Printf("Alias '%s' deleted.\n", name)
			return nil
		}

		if len(args) == 0 {
			aliases, err := store.ListAliases()
			if err != nil {
				return fmt.Errorf("failed to list aliases: %v", err)
			}

			if len(aliases) == 0 {
				fmt.Println("No aliases set.")
				fmt.Println("Run 'gvm alias <name> <version>' to create one.")
				return nil
			}

			fmt.Println("Aliases:")
			fmt.Println(strings.Repeat("-", 50))
			for name, version := range aliases {
				fmt.Printf("  %-20s -> %s\n", name, version)
			}
			return nil
		}
		if len(args) == 2 {
			name := args[0]
			version := args[1]

			installed, err := store.IsVersionInstalled(version)
			if err != nil {
				return err
			}
			if !installed {
				fmt.Printf("version %s is not installed, run 'gvm install %s' to install it\n", version, version)
				return nil
			}

			err = store.SetAlias(name, version)
			if err != nil {
				return err
			}

			fmt.Printf("Alias '%s' -> %s created.\n", name, version)
			return nil
		}

		if len(args) == 1 {
			name := args[0]
			version, err := store.GetAlias(name)
			if err != nil {
				return err
			}
			fmt.Printf("%s -> %s\n", name, version)
			return nil
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(aliasCmd)
	aliasCmd.Flags().Bool("delete", false, "delet an alias")
}
