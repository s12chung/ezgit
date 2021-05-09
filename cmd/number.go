package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newNumberCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "number",
		Short: "Get ep number",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runNumber()
		},
	}
}

func runNumber() error {
	number, err := startGetLastEpisodeNumber()
	if err != nil {
		return err
	}
	fmt.Print(number)
	return nil
}
