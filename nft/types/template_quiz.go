package nft

import (
	"encoding/json"
	"fmt"

	"github.com/xuri/excelize/v2"
)

type TokenDataQuiz struct {
	Encryption    string `json:"encryption,omitempty"`
	EncryptedFlow string `json:"encrypted_flow,class_id"`
	LastRecipient string `json:"last_recipient,omitempty"`
	EscapeFlow    string `json:"escape_flow,omitempty"`
	Question      string `json:"question,omitempty"`
}

type InputTemplateQuiz struct {
	InputTemplate
	TokenData []TokenDataQuiz
}

func (t InputTemplateQuiz) Generate(outputFile string) error {
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

func (t InputTemplateQuiz) FromXLSX(file string) (Template, error) {
	f, err := excelize.OpenFile(file)
	if err != nil {
		fmt.Println(err)
		return InputTemplateQuiz{}, err
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
		return InputTemplateQuiz{}, err
	}

	tokenBaseInfo, err := t.readTokenBaseInfo(f)
	if err != nil {
		fmt.Println(err)
		return InputTemplateQuiz{}, err
	}

	tokenData, err := t.readTokenData(f)
	if err != nil {
		fmt.Println(err)
		return InputTemplateQuiz{}, err
	}

	return InputTemplateQuiz{
		InputTemplate: InputTemplate{
			SheetClass:    class,
			TokenBaseInfo: tokenBaseInfo,
		},
		TokenData: tokenData,
	}, nil
}

func (InputTemplateQuiz) readTokenData(xlsxFile *excelize.File) (infos []TokenDataQuiz, err error) {
	rows, err := xlsxFile.GetRows(SheetTokenData)
	if err != nil {
		return nil, err
	}

	headerRow := rows[0]
	fmt.Println("header", headerRow)

	dataRows := rows[1:]

	for _, dataRow := range dataRows {
		fmt.Println("data", dataRow)
		infos = append(infos, TokenDataQuiz{
			Encryption:    dataRow[0],
			EncryptedFlow: dataRow[1],
			LastRecipient: dataRow[2],
			EscapeFlow:    dataRow[3],
			Question:      dataRow[4],
		})
	}
	return
}
