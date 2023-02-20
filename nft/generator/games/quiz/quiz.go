package quiz

import (
	"encoding/json"

	"github.com/game-of-nfts/gon-toolbox/nft/types"
)

type TokenData struct {
	Encryption    string `json:"encryption,omitempty"`
	EncryptedFlow string `json:"encrypted_flow,class_id"`
	LastRecipient string `json:"last_recipient,omitempty"`
	EscapeFlow    string `json:"escape_flow,omitempty"`
	Question      string `json:"question,omitempty"`
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
			Encryption:    dataRow[0],
			EncryptedFlow: dataRow[1],
			LastRecipient: dataRow[2],
			EscapeFlow:    dataRow[3],
			Question:      dataRow[4],
		})
	}
	return nil
}
