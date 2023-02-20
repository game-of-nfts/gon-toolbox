package nft

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/xuri/excelize/v2"
)

const (
	SheetClass         = "class"
	SheetTokenBaseInfo = "token_base_info"
	SheetTokenData     = "token_data"
	SheetToken         = "token"
	TemplateIndividual = "individual"
	TemplateQuiz       = "quiz"
	TemplateTeam       = "team"
)

var TemplateReader = map[string]Template{
	TemplateIndividual: InputTemplateIndividual{},
	TemplateQuiz:       InputTemplateQuiz{},
	TemplateTeam:       InputTemplateTeam{},
}

type Class struct {
	ID               string `json:"id,omitempty"`
	Name             string `json:"name,omitempty"`
	Schema           string `json:"schema,omitempty"`
	Creator          string `json:"creator,omitempty"`
	Symbol           string `json:"symbol,omitempty"`
	MintRestricted   bool   `json:"mint_restricted,omitempty"`
	UpdateRestricted bool   `json:"update_restricted,omitempty"`
	Description      string `json:"description,omitempty"`
	Uri              string `json:"uri,omitempty"`
	UriHash          string `json:"uri_hash,omitempty"`
	Data             string `json:"data,omitempty"`
}

type TokenBaseInfo struct {
	ID        string `json:"id"`
	ClassID   string `json:"class_id"`
	Name      string `json:"name,omitempty"`
	URI       string `json:"uri,omitempty"`
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	UriHash   string `json:"uri_hash,omitempty"`
}

type TokenInfo struct {
	ID        string `json:"id"`
	ClassID   string `json:"class_id"`
	Name      string `json:"name,omitempty"`
	URI       string `json:"uri,omitempty"`
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	UriHash   string `json:"uri_hash,omitempty"`
	Data      string `json:"data,omitempty"`
}

type InputTemplate struct {
	SheetClass    Class
	TokenBaseInfo []TokenBaseInfo
}

type Template interface {
	FromXLSX(file string) (Template, error)
}

func GenerateToken(outputFile string, class Class, tokens []TokenInfo) error {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Create a class sheet.
	index, err := f.NewSheet(SheetClass)
	if err != nil {
		return err
	}

	// Set class header
	f.SetCellValue(SheetClass, "A1", "ID")
	f.SetCellValue(SheetClass, "B1", "Name")
	f.SetCellValue(SheetClass, "C1", "Schema")
	f.SetCellValue(SheetClass, "D1", "Sender")
	f.SetCellValue(SheetClass, "E1", "Symbol")
	f.SetCellValue(SheetClass, "F1", "MintRestricted")
	f.SetCellValue(SheetClass, "G1", "UpdateRestricted")
	f.SetCellValue(SheetClass, "H1", "Description")
	f.SetCellValue(SheetClass, "I1", "Uri")
	f.SetCellValue(SheetClass, "J1", "UriHash")
	f.SetCellValue(SheetClass, "K1", "Data")

	// Set class data
	f.SetCellValue(SheetClass, "A2", class.ID)
	f.SetCellValue(SheetClass, "B2", class.Name)
	f.SetCellValue(SheetClass, "C2", class.Schema)
	f.SetCellValue(SheetClass, "D2", class.Creator)
	f.SetCellValue(SheetClass, "E2", class.Symbol)
	f.SetCellValue(SheetClass, "F2", class.MintRestricted)
	f.SetCellValue(SheetClass, "G2", class.UpdateRestricted)
	f.SetCellValue(SheetClass, "H2", class.Description)
	f.SetCellValue(SheetClass, "I2", class.Uri)
	f.SetCellValue(SheetClass, "J2", class.UriHash)
	f.SetCellValue(SheetClass, "K2", class.Data)

	// Create a token sheet.
	_, err = f.NewSheet(SheetToken)
	if err != nil {
		return err
	}

	// Set token header
	f.SetCellValue(SheetToken, "A1", "ID")
	f.SetCellValue(SheetToken, "B1", "ClassID")
	f.SetCellValue(SheetToken, "C1", "Name")
	f.SetCellValue(SheetToken, "D1", "URI")
	f.SetCellValue(SheetToken, "E1", "Sender")
	f.SetCellValue(SheetToken, "F1", "Recipient")
	f.SetCellValue(SheetToken, "G1", "UriHash")
	f.SetCellValue(SheetToken, "H1", "Data")

	for i, token := range tokens {
		f.SetCellValue(SheetToken, fmt.Sprintf("A%d", i+2), token.ID)
		f.SetCellValue(SheetToken, fmt.Sprintf("B%d", i+2), token.ClassID)
		f.SetCellValue(SheetToken, fmt.Sprintf("C%d", i+2), token.URI)
		f.SetCellValue(SheetToken, fmt.Sprintf("D%d", i+2), token.Sender)
		f.SetCellValue(SheetToken, fmt.Sprintf("E%d", i+2), token.Recipient)
		f.SetCellValue(SheetToken, fmt.Sprintf("F%d", i+2), token.Recipient)
		f.SetCellValue(SheetToken, fmt.Sprintf("G%d", i+2), token.UriHash)
		f.SetCellValue(SheetToken, fmt.Sprintf("H%d", i+2), token.Data)
	}
	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
	return f.SaveAs(outputFile)
}

func (InputTemplate) readClass(xlsxFile *excelize.File) (Class, error) {
	rows, err := xlsxFile.GetRows(SheetClass)
	if err != nil {
		return Class{}, err
	}

	if len(rows) != 2 {
		return Class{}, errors.New("invalid class sheet, only support 2 rows")
	}

	headerRow := rows[0]
	dataRow := rows[1]

	fmt.Println("header", headerRow)

	mintRestricted, err := strconv.ParseBool(dataRow[5])
	if err != nil {
		return Class{}, err
	}

	updateRestricted, err := strconv.ParseBool(dataRow[6])
	if err != nil {
		return Class{}, err
	}
	return Class{
		ID:               dataRow[0],
		Name:             dataRow[1],
		Schema:           dataRow[2],
		Creator:          dataRow[3],
		Symbol:           dataRow[4],
		MintRestricted:   mintRestricted,
		UpdateRestricted: updateRestricted,
		Description:      dataRow[7],
		Uri:              dataRow[8],
		UriHash:          dataRow[9],
		Data:             dataRow[10],
	}, nil
}

func (InputTemplate) readTokenBaseInfo(xlsxFile *excelize.File) (infos []TokenBaseInfo, err error) {
	rows, err := xlsxFile.GetRows(SheetTokenBaseInfo)
	if err != nil {
		return nil, err
	}

	headerRow := rows[0]
	fmt.Println("header", headerRow)

	dataRows := rows[1:]

	for _, dataRow := range dataRows {
		infos = append(infos, TokenBaseInfo{
			ID:        dataRow[0],
			ClassID:   dataRow[1],
			Name:      dataRow[2],
			URI:       dataRow[3],
			Sender:    dataRow[4],
			Recipient: dataRow[5],
			UriHash:   dataRow[6],
		})
	}
	return
}
