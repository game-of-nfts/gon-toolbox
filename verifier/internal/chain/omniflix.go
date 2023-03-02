package chain

type Omniflix struct{}

func (Omniflix) GetTx(txHash string) (*TxResult, error)     { return nil, nil }
func (Omniflix) GetNFT(classID, nftID string) (*NFT, error) { return nil, nil }
func (Omniflix) HasNFT(classID, nftID string) bool          { return false }
func (Omniflix) GetClass(classID string) (*Class, error)    { return nil, nil }
func (Omniflix) HasClass(classID string) bool               { return false }
