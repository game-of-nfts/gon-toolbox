package quiz

import (
	"encoding/json"

	"github.com/irisnet/core-sdk-go/common/crypto"
	sdktypes "github.com/irisnet/core-sdk-go/types"

	"github.com/game-of-nfts/gon-toolbox/nft/types"
)

type TokenData struct {
	Question           string `json:"question,omitempty"`
	MnemonicsEncrypted string `json:"mnemonics_encrypted,omitempty"`
	Flow               string `json:"flow,omitempty"`
}

type TokenDataXlsx struct {
	Question string `json:"question,omitempty"`
	Answer   string `json:"answer,omitempty"`
	Flow     string `json:"flow,omitempty"`
}

type Template struct {
	types.BaseTemplate
	TokenData []TokenDataXlsx
}

func NewTemplate(args types.InputArgs) (types.Template, error) {
	baseTpl, tokenDataRows, err := types.NewTemplate(args)
	if err != nil {
		return nil, err
	}

	tpl := &Template{
		BaseTemplate: baseTpl,
		TokenData:    make([]TokenDataXlsx, 0, len(baseTpl.TokenBaseInfo)),
	}
	if err = tpl.FillTokenData(tokenDataRows); err != nil {
		return nil, err
	}
	return tpl, nil
}

func (t Template) Generate() error {
	tokens := make([]types.TokenInfo, 0, len(t.TokenData))
	for i, data := range t.TokenData {
		keyManager, err := crypto.NewAlgoKeyManager("secp256k1")
		if err != nil {
			return err
		}
		mnemonics, _ := keyManager.Generate()
		pubKey := keyManager.ExportPubKey()
		address := sdktypes.AccAddress(pubKey.Address().Bytes()).String()

		mnemonicsEncrypted, err := types.Encrypt(data.Answer, mnemonics)
		if err != nil {
			return err
		}

		metadata := TokenData{
			Question:           data.Question,
			MnemonicsEncrypted: mnemonicsEncrypted,
			Flow:               data.Flow,
		}
		bz, err := json.Marshal(metadata)
		if err != nil {
			return err
		}

		tokens = append(tokens, types.TokenInfo{
			ID:        t.TokenBaseInfo[i].ID,
			ClassID:   t.TokenBaseInfo[i].ClassID,
			Name:      t.TokenBaseInfo[i].Name,
			URI:       t.TokenBaseInfo[i].URI,
			Sender:    t.Args.Sender,
			Recipient: address,
			UriHash:   t.TokenBaseInfo[i].UriHash,
			Data:      string(bz),
		})
	}
	return t.GenerateToken(tokens)
}

func (t *Template) FillTokenData(dataRows [][]string) error {
	for _, dataRow := range dataRows {
		t.TokenData = append(t.TokenData, TokenDataXlsx{
			Question: dataRow[0],
			Answer:   dataRow[1],
			Flow:     dataRow[2],
		})
	}
	return nil
}
