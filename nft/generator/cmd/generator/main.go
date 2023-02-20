package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/game-of-nfts/gon-toolbox/nft/generator/games/individual"
	"github.com/game-of-nfts/gon-toolbox/nft/generator/games/quiz"
	"github.com/game-of-nfts/gon-toolbox/nft/generator/games/team"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "generator",
		Short: "command",
	}

	rootCmd.AddCommand(
		individual.Commands("", ""),
		quiz.Commands("", ""),
		team.Commands("", ""),
	)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
