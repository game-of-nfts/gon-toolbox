package nft

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/xuri/excelize/v2"
)

type TokenDataIndividual struct {
	Type        string `json:"type,omitempty"`
	Flow        string `json:"flow,class_id"`
	LastBatton  string `json:"last_batton,omitempty"`
	StartHeight string `json:"start_height,omitempty"`
}

type InputTemplateIndividual struct {
	InputTemplate
	TokenData []TokenDataIndividual
}

func (t InputTemplateIndividual) Generate(outputFile string) error {
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

func (t InputTemplateIndividual) FromXLSX(file string) (Template, error) {
	f, err := excelize.OpenFile(file)
	if err != nil {
		return InputTemplateIndividual{}, err
	}

	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	class, err := t.readClass(f)
	if err != nil {
		return InputTemplateIndividual{}, err
	}

	tokenBaseInfo, err := t.readTokenBaseInfo(f)
	if err != nil {
		return InputTemplateIndividual{}, err
	}

	tokenData, err := t.readTokenData(f)
	if err != nil {
		return InputTemplateIndividual{}, err
	}

	if len(tokenData) != len(tokenBaseInfo) {
		return nil, errors.New("the lenght of tokenData and tokenBaseInfo is unmatch")
	}

	return InputTemplateIndividual{
		InputTemplate: InputTemplate{
			SheetClass:    class,
			TokenBaseInfo: tokenBaseInfo,
		},
		TokenData: tokenData,
	}, nil
}

func (InputTemplateIndividual) readTokenData(xlsxFile *excelize.File) (infos []TokenDataIndividual, err error) {
	rows, err := xlsxFile.GetRows(SheetTokenData)
	if err != nil {
		return nil, err
	}

	headerRow := rows[0]
	fmt.Println("header", headerRow)

	dataRows := rows[1:]

	for _, dataRow := range dataRows {
		fmt.Println("data", dataRow)
		infos = append(infos, TokenDataIndividual{
			Type:        dataRow[0],
			Flow:        dataRow[1],
			LastBatton:  dataRow[2],
			StartHeight: dataRow[3],
		})
	}
	return
}
