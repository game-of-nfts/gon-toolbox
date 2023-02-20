package team

import (
	"github.com/spf13/cobra"
)

func Commands(inputFile, outputFilePath string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "team",
		Short: "xxxx",
	}

	return cmd
}
