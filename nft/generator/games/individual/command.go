package individual

import (
	"github.com/spf13/cobra"
)

func Commands(finputFile, outputFilePath string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "individual",
		Short: "xxxx",
	}

	return cmd
}
