package individual

import (
	"github.com/spf13/cobra"

	"github.com/game-of-nfts/gon-toolbox/nft/types"
)

func Commands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "individual",
		Short: "Generate nft of individual game type",
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
