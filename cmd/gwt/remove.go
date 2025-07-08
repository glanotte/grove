package gwt

import (
	"fmt"
	"github.com/spf13/cobra"
)

func newRemoveCmd() *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "remove <worktree-name>",
		Short: "Remove a worktree and its associated resources",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			worktreeName := args[0]
			fmt.Printf("Removing worktree %s...\n", worktreeName)
			// Implementation would go here
			return nil
		},
	}

	cmd.Flags().BoolVar(&force, "force", false, "Force removal even if there are uncommitted changes")
	return cmd
}
