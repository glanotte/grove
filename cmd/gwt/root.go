package gwt

import (
	"github.com/spf13/cobra"
)

func NewRootCmd(version, commit, date string) *cobra.Command {
	var configFile string

	rootCmd := &cobra.Command{
		Use:   "grove",
		Short: "Git worktree manager with Docker and template support",
		Long: `grove is a CLI tool for managing git worktrees with template support,
Docker integration, and automatic web serving configuration.`,
	}

	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is .grove/config.yaml)")

	rootCmd.AddCommand(
		newInitCmd(),
		newCreateCmd(),
		newListCmd(),
		newRemoveCmd(),
		newSwitchCmd(),
		newVersionCmd(version, commit, date),
	)

	return rootCmd
}
