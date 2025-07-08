package gwt

import (
	"fmt"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	var format string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all worktrees with their status",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Listing worktrees...")
			// Implementation would go here
			return nil
		},
	}

	cmd.Flags().StringVar(&format, "format", "table", "Output format (table|json)")
	return cmd
}
