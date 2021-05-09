package cmd

import (
	"github.com/spf13/cobra"
)

func newGenerateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "generate",
		Short: "Generate New Ep Files",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGenerate()
		},
	}
}

func runGenerate() error {
	return startGenerateFiles()
}
