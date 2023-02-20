package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/game-of-nfts/gon-toolbox/nft/generator/games/individual"
	"github.com/game-of-nfts/gon-toolbox/nft/generator/games/quiz"
	"github.com/game-of-nfts/gon-toolbox/nft/generator/games/team"
	"github.com/game-of-nfts/gon-toolbox/nft/types"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "generator",
		Short: "GoN testnet nft generator",
	}

	rootCmd.AddCommand(
		individual.Commands(),
		quiz.Commands(),
		team.Commands(),
	)

	pflags := rootCmd.PersistentFlags()
	pflags.AddFlagSet(types.ConfigSet)
	_ = rootCmd.MarkFlagRequired(types.FlagTokenFile)
	_ = rootCmd.MarkFlagRequired(types.FlagAddressFile)
	_ = rootCmd.MarkFlagRequired(types.FlagOutputPath)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
