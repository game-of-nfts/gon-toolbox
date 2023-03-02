package chain

type Juno struct{}

func (Juno) GetTx(txHash string) (*TxResult, error)     { return nil, nil }
func (Juno) GetNFT(classID, nftID string) (*NFT, error) { return nil, nil }
func (Juno) HasNFT(classID, nftID string) bool          { return false }
func (Juno) GetClass(classID string) (*Class, error)    { return nil, nil }
func (Juno) HasClass(classID string) bool               { return false }
