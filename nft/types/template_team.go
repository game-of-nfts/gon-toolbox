package nft

import (
	"encoding/json"
	"fmt"

	"github.com/xuri/excelize/v2"
)

type TokenDataTeam struct {
	Type          string `json:"type,omitempty"`
	Flow          string `json:"flow,class_id"`
	Battons       string `json:"battons,omitempty"`
	LastRecipient string `json:"last_recipient,omitempty"`
	StartHeight   string `json:"start_height,omitempty"`
}

type InputTemplateTeam struct {
	InputTemplate
	TokenData []TokenDataTeam
}

func (t InputTemplateTeam) Generate(outputFile string) error {
	tokens := make([]TokenInfo, 0, len(t.TokenData))
	for i, data := range t.TokenData {
		bz, err := json.Marshal(data)
		if err != nil {
			return err
		}
		tokens = append(tokens, TokenInfo{
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
	return GenerateToken(outputFile, t.SheetClass, tokens)
}

func (t InputTemplateTeam) FromXLSX(file string) (Template, error) {
	f, err := excelize.OpenFile(file)
	if err != nil {
		fmt.Println(err)
		return InputTemplateTeam{}, err
	}

	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	class, err := t.readClass(f)
	if err != nil {
		fmt.Println(err)
		return InputTemplateTeam{}, err
	}

	tokenBaseInfo, err := t.readTokenBaseInfo(f)
	if err != nil {
		fmt.Println(err)
		return InputTemplateTeam{}, err
	}

	tokenData, err := t.readTokenData(f)
	if err != nil {
		fmt.Println(err)
		return InputTemplateTeam{}, err
	}

	return InputTemplateTeam{
		InputTemplate: InputTemplate{
			SheetClass:    class,
			TokenBaseInfo: tokenBaseInfo,
		},
		TokenData: tokenData,
	}, nil
}

func (InputTemplateTeam) readTokenData(xlsxFile *excelize.File) (infos []TokenDataTeam, err error) {
	rows, err := xlsxFile.GetRows(SheetTokenData)
	if err != nil {
		return nil, err
	}

	headerRow := rows[0]
	fmt.Println("header", headerRow)

	dataRows := rows[1:]

	for _, dataRow := range dataRows {
		fmt.Println("data", dataRow)
		infos = append(infos, TokenDataTeam{
			Type:          dataRow[0],
			Flow:          dataRow[1],
			Battons:       dataRow[2],
			LastRecipient: dataRow[3],
			StartHeight:   dataRow[4],
		})
	}
	return
}
