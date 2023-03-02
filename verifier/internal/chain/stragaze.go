package chain

type Stargaze struct{}

func (Stargaze) GetTx(txHash string) (*TxResult, error)     { return nil, nil }
func (Stargaze) GetNFT(classID, nftID string) (*NFT, error) { return nil, nil }
func (Stargaze) HasNFT(classID, nftID string) bool          { return false }
func (Stargaze) GetClass(classID string) (*Class, error)    { return nil, nil }
func (Stargaze) HasClass(classID string) bool               { return false }
