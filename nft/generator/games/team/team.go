package team

import (
	"encoding/json"

	"github.com/game-of-nfts/gon-toolbox/nft/types"
)

type TokenData struct {
	Type          string `json:"type,omitempty"`
	Flow          string `json:"flow,class_id"`
	Battons       string `json:"battons,omitempty"`
	LastRecipient string `json:"last_recipient,omitempty"`
	StartHeight   string `json:"start_height,omitempty"`
}

type Template struct {
	types.BaseTemplate
	TokenData []TokenData
}

func NewTemplate(args types.InputArgs) (types.Template, error) {
	baseTpl, tokenDataRows, err := types.NewTemplate(args)
	if err != nil {
		return nil, err
	}

	tpl := &Template{
		BaseTemplate: baseTpl,
		TokenData:    make([]TokenData, 0, len(baseTpl.TokenBaseInfo)),
	}

	if err = tpl.FillTokenData(tokenDataRows); err != nil {
		return nil, err
	}
	return tpl, nil
}

func (t Template) Generate() error {
	tokens := make([]types.TokenInfo, 0, len(t.TokenData))
	for i, data := range t.TokenData {
		bz, err := json.Marshal(data)
		if err != nil {
			return err
		}
		tokens = append(tokens, types.TokenInfo{
			ID:        t.TokenBaseInfo[i].ID,
			ClassID:   t.TokenBaseInfo[i].ClassID,
			Name:      t.TokenBaseInfo[i].Name,
			URI:       t.TokenBaseInfo[i].URI,
			Sender:    t.TokenBaseInfo[i].Sender,
			Recipient: t.TokenBaseInfo[i].Recipient,
			UriHash:   t.TokenBaseInfo[i].UriHash,
			Data:      string(bz),
		})
	}
	return t.GenerateToken(tokens)
}

func (t *Template) FillTokenData(dataRows [][]string) error {
	for _, dataRow := range dataRows {
		t.TokenData = append(t.TokenData, TokenData{
			Type:          dataRow[0],
			Flow:          dataRow[1],
			Battons:       dataRow[2],
			LastRecipient: dataRow[3],
			StartHeight:   dataRow[4],
		})
	}
	return nil
}
