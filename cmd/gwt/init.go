package gwt

import (
    "fmt"
    "github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
    var template string

    cmd := &cobra.Command{
        Use:   "init <repo-url>",
        Short: "Initialize a bare repository for worktree management",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            repoURL := args[0]
            fmt.Printf("Initializing bare repository from %s...\n", repoURL)
            // Implementation would go here
            return nil
        },
    }

    cmd.Flags().StringVar(&template, "template", "default", "Template to use for initialization")
    return cmd
}