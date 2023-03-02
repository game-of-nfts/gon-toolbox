package chain

type Uptickd struct{}

func (Uptickd) GetTx(txHash string) (*TxResult, error)     { return nil, nil }
func (Uptickd) GetNFT(classID, nftID string) (*NFT, error) { return nil, nil }
func (Uptickd) HasNFT(classID, nftID string) bool          { return false }
func (Uptickd) GetClass(classID string) (*Class, error)    { return nil, nil }
func (Uptickd) HasClass(classID string) bool               { return false }
