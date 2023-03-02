package chain

type Iris struct{}

func (Iris) GetTx(txHash string) (*TxResult, error)     { return nil, nil }
func (Iris) GetNFT(classID, nftID string) (*NFT, error) { return nil, nil }
func (Iris) HasNFT(classID, nftID string) bool          { return false }
func (Iris) GetClass(classID string) (*Class, error)    { return nil, nil }
func (Iris) HasClass(classID string) bool               { return false }
