package types

import (
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

const (
	FlagTokenFile   = "token-file"
	FlagUserFile    = "user-file"
	FlagOutputPath  = "output-path"
	FlagTxSender    = "sender"
	FlagChainConfig = "config-file"
)

// common flagsets to add to various functions
var (
	ConfigSet = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	ConfigSet.String(FlagTokenFile, "", "Generate the source files needed for nft")
	ConfigSet.String(FlagUserFile, "", "User information file")
	ConfigSet.String(FlagOutputPath, "", "NFT output path")
	ConfigSet.String(FlagTxSender, "", "Issuer address(bech32) of nft")
	ConfigSet.String(FlagChainConfig, "", "chain config file")
}

type InputArgs struct {
	TokenFile   string
	AddressFile string
	OutputPath  string
	Sender      string
}

func ReadArgs(cmd *cobra.Command) (InputArgs, error) {
	tokenFile, err := cmd.Flags().GetString(FlagTokenFile)
	if err != nil {
		return InputArgs{}, err
	}

	addressFile, err := cmd.Flags().GetString(FlagUserFile)
	if err != nil {
		return InputArgs{}, err
	}

	outputPath, err := cmd.Flags().GetString(FlagOutputPath)
	if err != nil {
		return InputArgs{}, err
	}

	sender, err := cmd.Flags().GetString(FlagTxSender)
	if err != nil {
		return InputArgs{}, err
	}
	ValidateAddress(PrefixBech32Iris, sender)

	return InputArgs{
		TokenFile:   tokenFile,
		AddressFile: addressFile,
		OutputPath:  outputPath,
		Sender:      sender,
	}, nil

}
