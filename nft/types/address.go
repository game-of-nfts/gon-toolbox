package types

import (
	"fmt"

	"github.com/irisnet/core-sdk-go/common/bech32"
)

const (
	PrefixBech32Iris     = "iaa"
	PrefixBech32Stars    = "stars"
	PrefixBech32Juno     = "juno"
	PrefixBech32Uptick   = "uptick"
	PrefixBech32Omniflix = "omniflix"
)

func ValidateAddress(hrp, address string) {
	_, err := bech32.GetFromBech32(address, hrp)
	if err != nil {
		panic(
			fmt.Errorf(
				"invalid address:%s,err:%s",
				address,
				err.Error(),
			),
		)
	}
}
