package main

import (
	"github.com/spf13/cobra"
	"os"

	"github.com/game-of-nfts/gon-toolbox/nft/airdrop/airdrop"
	"github.com/game-of-nfts/gon-toolbox/nft/types"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "airdrop",
		Short: "GoN testnet nft airdrop",
	}

	rootCmd.AddCommand(airdrop.Commands())

	pflags := rootCmd.PersistentFlags()
	pflags.AddFlagSet(types.ConfigSet)
	_ = rootCmd.MarkFlagRequired(types.FlagTokenFile)
	_ = rootCmd.MarkFlagRequired(types.FlagChainConfig)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
