package gongen

import (
	"github.com/spf13/cobra"

	"github.com/game-of-nfts/gon-toolbox/nft/generator/games/individual"
	"github.com/game-of-nfts/gon-toolbox/nft/generator/games/quiz"
	"github.com/game-of-nfts/gon-toolbox/nft/generator/games/team"
	"github.com/game-of-nfts/gon-toolbox/nft/types"
)

func Commands() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gen",
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
	_ = rootCmd.MarkFlagRequired(types.FlagOutputPath)
	_ = rootCmd.MarkFlagRequired(types.FlagTxSender)
	_ = rootCmd.MarkFlagRequired(types.FlagUserFile)

	return rootCmd
}
