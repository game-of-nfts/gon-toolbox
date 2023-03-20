package main

import (
	"os"

	"github.com/spf13/cobra"

	aridrop "github.com/game-of-nfts/gon-toolbox/nft/airdrop/cmd/airdrop"
	"github.com/game-of-nfts/gon-toolbox/nft/generator/cmd/gongen"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "gontool",
		Short: "GoN tools",
	}

	rootCmd.AddCommand(
		gongen.Commands(),
		aridrop.Commands(),
		CryptoCommands(),
	)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
