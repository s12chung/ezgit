package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ezgit",
		Short: "Git Shorthands For s12chung/trailheadspodcast",
	}
}

// Execute the CLI
func Execute() {
	rootCmd := newRootCmd()
	rootCmd.AddCommand(newStartCmd())
	rootCmd.AddCommand(newPushCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
