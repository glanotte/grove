package gwt

import (
    "fmt"
    "github.com/spf13/cobra"
)

func newVersionCmd(version, commit, date string) *cobra.Command {
    return &cobra.Command{
        Use:   "version",
        Short: "Print version information",
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Printf("grove version %s (commit: %s, built: %s)\n", version, commit, date)
        },
    }
}