package aridrop

import (
	"github.com/spf13/cobra"

	"github.com/game-of-nfts/gon-toolbox/nft/airdrop/airdrop"
	"github.com/game-of-nfts/gon-toolbox/nft/types"
)

func Commands() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "airdrop",
		Short: "GoN testnet nft airdrop",
	}

	rootCmd.AddCommand(airdrop.Commands())

	pflags := rootCmd.PersistentFlags()
	pflags.AddFlagSet(types.ConfigSet)
	_ = rootCmd.MarkFlagRequired(types.FlagTokenFile)
	_ = rootCmd.MarkFlagRequired(types.FlagChainConfig)

	return rootCmd
}
