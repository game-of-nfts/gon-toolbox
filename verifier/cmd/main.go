package verifier

import (
	"github.com/spf13/cobra"
)

func Commands() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "verify",
		Short: "GoN evidence verify tools",
	}

	//TODO
	rootCmd.AddCommand(
	// gongen.Commands(),
	// aridrop.Commands(),
	// CryptoCommands(),
	)
	return rootCmd
}
