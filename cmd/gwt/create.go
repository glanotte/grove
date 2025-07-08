package gwt

import (
	"fmt"
	"github.com/spf13/cobra"
)

func newCreateCmd() *cobra.Command {
	var (
		from     string
		template string
	)

	cmd := &cobra.Command{
		Use:   "create <branch-name>",
		Short: "Create a new worktree from a branch",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			branchName := args[0]
			fmt.Printf("Creating worktree for branch %s...\n", branchName)
			// Implementation would go here
			return nil
		},
	}

	cmd.Flags().StringVar(&from, "from", "main", "Base branch to create from")
	cmd.Flags().StringVar(&template, "template", "default", "Template to use")
	return cmd
}
