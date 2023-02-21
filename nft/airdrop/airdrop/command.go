package airdrop

import (
	"github.com/BurntSushi/toml"
	"github.com/game-of-nfts/gon-toolbox/nft/types"
	"github.com/spf13/cobra"
)

func Commands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exec",
		Short: "execute NFT airdrop",
		RunE: func(cmd *cobra.Command, args []string) error {
			inputFile, err := cmd.Flags().GetString(types.FlagTokenFile)
			if err != nil {
				return err
			}
			configFile, err := cmd.Flags().GetString(types.FlagChainConfig)
			if err != nil {
				return err
			}

			cfg := &Config{}
			_, err = toml.DecodeFile(configFile, cfg)
			if err != nil {
				return err
			}
			ad, err := NewAirDropper(inputFile, cfg)
			if err != nil {
				return err
			}

			return ad.ExecAirdrop()
		},
	}
	return cmd
}
