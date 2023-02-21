package team

import (
	"github.com/game-of-nfts/gon-toolbox/nft/types"
	"github.com/spf13/cobra"
)

func Commands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "team",
		Short: "Generate nft of team game type",
		RunE: func(cmd *cobra.Command, args []string) error {
			myArgs, err := types.ReadArgs(cmd)
			if err != nil {
				return err
			}

			tpl, err := NewTemplate(myArgs)
			if err != nil {
				return err
			}
			return tpl.Generate()
		},
	}
	_ = cmd.MarkFlagRequired(types.FlagUserFile)
	return cmd
}
