package types

import (
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

const (
	FlagTokenFile   = "token-file"
	FlagAddressFile = "address-file"
	FlagOutputPath  = "output-path"
	FlagTxSender    = "sender"
)

// common flagsets to add to various functions
var (
	ConfigSet = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	ConfigSet.String(FlagTokenFile, "", "bech32 encoded account address")
	ConfigSet.String(FlagAddressFile, "", "description of account")
	ConfigSet.String(FlagOutputPath, "", "bech32 encoded account address")
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

	addressFile, err := cmd.Flags().GetString(FlagAddressFile)
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

	return InputArgs{
		TokenFile:   tokenFile,
		AddressFile: addressFile,
		OutputPath:  outputPath,
		Sender:      sender,
	}, nil

}
