package quiz

import (
	"encoding/json"
	"fmt"

	"github.com/game-of-nfts/gon-toolbox/nft/types"
	"github.com/xuri/excelize/v2"
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
	return t.GenerateToken(t.Args.OutputPath, tokens)
}

func (t Template) ReadFromXLSX(args types.InputArgs) (types.Template, error) {
	f, err := excelize.OpenFile(args.TokenFile)
	if err != nil {
		fmt.Println(err)
		return Template{}, err
	}

	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	class, err := t.ReadClass(f)
	if err != nil {
		fmt.Println(err)
		return Template{}, err
	}

	tokenBaseInfo, err := t.ReadTokenBaseInfo(f)
	if err != nil {
		fmt.Println(err)
		return Template{}, err
	}

	tokenData, err := t.readTokenData(f)
	if err != nil {
		fmt.Println(err)
		return Template{}, err
	}

	return Template{
		BaseTemplate: types.BaseTemplate{
			SheetClass:    class,
			TokenBaseInfo: tokenBaseInfo,
		},
		TokenData: tokenData,
	}, nil
}

func (Template) readTokenData(xlsxFile *excelize.File) (infos []TokenData, err error) {
	rows, err := xlsxFile.GetRows(types.SheetTokenData)
	if err != nil {
		return nil, err
	}

	headerRow := rows[0]
	fmt.Println("header", headerRow)

	dataRows := rows[1:]

	for _, dataRow := range dataRows {
		fmt.Println("data", dataRow)
		infos = append(infos, TokenData{
			Encryption:    dataRow[0],
			EncryptedFlow: dataRow[1],
			LastRecipient: dataRow[2],
			EscapeFlow:    dataRow[3],
			Question:      dataRow[4],
		})
	}
	return
}
