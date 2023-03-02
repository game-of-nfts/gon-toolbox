package chain

import (
	"github.com/irisnet/core-sdk-go/types"
)

const (
	ChainIdAbbreviationIris     = "i"
	ChainIdAbbreviationStars    = "s"
	ChainIdAbbreviationJuno     = "j"
	ChainIdAbbreviationUptick   = "u"
	ChainIdAbbreviationOmniflix = "o"
)

type (
	TxResult struct {
		Events types.StringEvents
		Hash   string
		Height int64
		Time   int64
		Sender string
	}

	Class struct {
		ID      string
		Name    string
		Schema  string
		Creator string
		Uri     string
		UriHash string
		Data    string
	}

	NFT struct {
		ID      string
		Name    string
		URI     string
		Data    string
		Owner   string
		URIHash string
	}

	Chain interface {
		GetTx(txHash string) (*TxResult, error)
		GetNFT(classID, nftID string) (*NFT, error)
		HasNFT(classID, nftID string) bool
		GetClass(classID string) (*Class, error)
		HasClass(classID string) bool
	}

	Registry struct {
		chains map[string]Chain
	}
)

func NewRegistry() *Registry {
	return &Registry{
		chains: map[string]Chain{
			ChainIdAbbreviationIris:     Iris{},
			ChainIdAbbreviationStars:    Stargaze{},
			ChainIdAbbreviationJuno:     Juno{},
			ChainIdAbbreviationUptick:   Uptickd{},
			ChainIdAbbreviationOmniflix: Omniflix{},
		},
	}
}

func (cr *Registry) GetChain(chainID string) Chain {
	return cr.chains[chainID]
}
