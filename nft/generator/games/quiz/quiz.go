package quiz

import (
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"

	"github.com/irisnet/core-sdk-go/common/crypto"
	sdktypes "github.com/irisnet/core-sdk-go/types"

	"github.com/game-of-nfts/gon-toolbox/nft/types"
)

type TokenData struct {
	Question           string `json:"question,omitempty"`
	MnemonicsEncrypted string `json:"mnemonics_encrypted,omitempty"`
	Flow               string `json:"flow,omitempty"`
}

type Row struct {
	Question string `json:"question,omitempty"`
	Answer   string `json:"answer,omitempty"`
	Flow     string `json:"flow,omitempty"`
}

type Template struct {
	types.BaseTemplate
	Rows []Row
}

func NewTemplate(args types.InputArgs) (types.Template, error) {
	btl, err := types.NewBaseTemplate(args)
	if err != nil {
		return nil, err
	}

	tpl := &Template{
		BaseTemplate: btl,
		Rows:    make([]Row, 0),
	}

	err = tpl.PreInitialize(tpl.GetPreInitializer())
	if err != nil {
		return nil, err
	}

	err = tpl.Initialize()
	if err != nil {
		return nil, err
	}

	if err = tpl.FillRows(tpl.BaseTemplate.TokenData); err != nil {
		return nil, err
	}
	return tpl, nil
}

func (t Template) GetPreInitializer() types.PreInitializer {
	return func(tpl *types.BaseTemplate, f *excelize.File) {
		for i := 0; i < 100; i++ {
			f.SetCellValue(types.SheetTokenBaseInfo, fmt.Sprintf("A%d", i+2), fmt.Sprintf("quiz%03d", i+1))
			f.SetCellValue(types.SheetTokenBaseInfo, fmt.Sprintf("B%d", i+2), tpl.SheetClass.ID)
			f.SetCellValue(types.SheetTokenBaseInfo, fmt.Sprintf("C%d", i+2), fmt.Sprintf("quiz%03d", i+1))
		}
	}
}

func (t Template) Generate() error {
	tokens := make([]types.TokenInfo, 0, len(t.BaseTemplate.TokenData))
	for i, data := range t.Rows {
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
		fmt.Println(t.TokenBaseInfo[i].ClassID + " " + t.TokenBaseInfo[i].ID + " " + mnemonics)

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

func (t *Template) FillRows(dataRows [][]string) error {
	for _, dataRow := range dataRows {
		t.Rows = append(t.Rows, Row{
			Question: dataRow[0],
			Answer:   dataRow[1],
			Flow:     dataRow[2],
		})
	}
	return nil
}
