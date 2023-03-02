package chain

import (
	"github.com/irisnet/core-sdk-go/types"
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
		chains: make(map[string]Chain),
	}
}

func (cr *Registry) Register(chainID string, chain Chain) *Registry {
	cr.chains[chainID] = chain
	return cr
}

func (cr *Registry) GetChain(chainID string) Chain {
	return cr.chains[chainID]
}
