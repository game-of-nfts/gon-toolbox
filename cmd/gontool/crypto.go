package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/game-of-nfts/gon-toolbox/nft/types"
)

func CryptoCommands() *cobra.Command {
	cryptoCmd := &cobra.Command{
		Use:   "aes",
		Short: "GoN crypto tools",
	}

	cryptoCmd.AddCommand(
		EncryptCommand(),
		DecryptCommand(),
	)
	return cryptoCmd
}

func EncryptCommand() *cobra.Command {
	enCmd := &cobra.Command{
		Use:     "encrypt",
		Short:   "GoN encrypt command",
		Example: "encrypt [key] [text]",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			encryptedText, err := types.Encrypt(args[0], args[1])
			if err != nil {
				return err
			}
			fmt.Println(encryptedText)
			return nil
		},
	}
	return enCmd
}

func DecryptCommand() *cobra.Command {
	decrypt := &cobra.Command{
		Use:     "decrypt",
		Short:   "GoN decrypt command",
		Example: "decrypt [key] [text]",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			dncryptedText, err := types.Decrypt(args[0], args[1])
			if err != nil {
				return err
			}
			fmt.Println(dncryptedText)
			return nil
		},
	}
	return decrypt
}
