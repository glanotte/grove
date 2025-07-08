package gwt

import (
	"fmt"
	"github.com/spf13/cobra"
)

func newSwitchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "switch <worktree-name>",
		Short: "Switch to a different worktree (outputs cd command)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			worktreeName := args[0]
			// This would output the path for shell integration
			fmt.Printf("/path/to/worktree/%s\n", worktreeName)
			return nil
		},
	}
}
