package quiz

import (
	"github.com/spf13/cobra"
)

func Commands(finputFile, outputFilePath string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "quiz",
		Short: "xxxx",
	}

	return cmd
}
