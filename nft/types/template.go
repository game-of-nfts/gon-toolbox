package types

import (
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"strconv"
)

const (
	SheetClass         = "class"
	SheetTokenBaseInfo = "token_base_info"
	SheetTokenData     = "token_data"
	SheetToken         = "token"
	SheetTeams         = "teams"
	TemplateIndividual = "individual"
	TemplateQuiz       = "quiz"
	TemplateTeam       = "team"
)

type Class struct {
	ID               string `json:"id,omitempty"`
	Name             string `json:"name,omitempty"`
	Schema           string `json:"schema,omitempty"`
	Sender           string `json:"sender,omitempty"`
	Symbol           string `json:"symbol,omitempty"`
	MintRestricted   bool   `json:"mint_restricted,omitempty"`
	UpdateRestricted bool   `json:"update_restricted,omitempty"`
	Description      string `json:"description,omitempty"`
	Uri              string `json:"uri,omitempty"`
	UriHash          string `json:"uri_hash,omitempty"`
	Data             string `json:"data,omitempty"`
}

// TokenBaseInfo represents a row in token_base_info sheet of the input file
type TokenBaseInfo struct {
	ID      string `json:"id"`
	ClassID string `json:"class_id"`
	Name    string `json:"name,omitempty"`
	URI     string `json:"uri,omitempty"`
	UriHash string `json:"uri_hash,omitempty"`
}

// TokenInfo represents a row in token_base_info sheet of the output file
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

type Template interface {
	Generate() error
	FillRows(dataRows [][]string) error
}

type BaseTemplate struct {
	*UserSelector
	SheetClass    Class
	TokenBaseInfo []TokenBaseInfo
	TokenData     [][]string
	Args          InputArgs
}

type PreInitializer func(tpl *BaseTemplate, file *excelize.File)

func NewBaseTemplate(args InputArgs) (BaseTemplate, error) {
	tpl := BaseTemplate{Args: args}
	selector, err := NewTeamSelector(args)
	if err != nil {
		return tpl, err
	}
	tpl.UserSelector = selector
	return tpl, nil
}

// PreInitialize reads class sheet and initialize the base_token_info sheet
func (tpl *BaseTemplate) PreInitialize(piz PreInitializer) error {
	f, err := excelize.OpenFile(tpl.Args.TokenFile)
	if err != nil {
		return err
	}
	defer func() {
		err = f.Save()
		if err != nil {
			fmt.Println("Error saving file:", err)
			return
		}
	}()

	class, err := tpl.readClass(f)
	if err != nil {
		return err
	}
	tpl.SheetClass = class
	piz(tpl, f)
	return nil
}

// Initialize reads the base_token_info sheet and token_data sheet
func (tpl *BaseTemplate) Initialize() error {
	f, err := excelize.OpenFile(tpl.Args.TokenFile)
	if err != nil {
		return err
	}

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	tokenBaseInfo, err := tpl.readTokenBaseInfo(f)
	if err != nil {
		return err
	}
	tpl.TokenBaseInfo = tokenBaseInfo

	rows, err := f.GetRows(SheetTokenData)
	if err != nil {
		return err
	}

	tpl.TokenData = rows[1:]
	if len(tpl.TokenBaseInfo) != len(tpl.TokenData) {
		return errors.New("the length of token_base_info and token_data is unmatched")
	}
	return nil
}

// TODO: remove this function
func NewTemplate(args InputArgs) (BaseTemplate, [][]string, error) {
	tpl := BaseTemplate{
		Args: args,
	}

	f, err := excelize.OpenFile(args.TokenFile)
	if err != nil {
		return tpl, nil, err
	}

	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	class, err := tpl.readClass(f)
	if err != nil {
		return tpl, nil, err
	}
	tpl.SheetClass = class

	tokenBaseInfo, err := tpl.readTokenBaseInfo(f)
	if err != nil {
		return tpl, nil, err
	}
	tpl.TokenBaseInfo = tokenBaseInfo

	rows, err := f.GetRows(SheetTokenData)
	if err != nil {
		return tpl, nil, err
	}

	headerRow := rows[0]
	dataRow := rows[1:]
	PrintXLSX(SheetTokenData, headerRow, dataRow)

	if len(tokenBaseInfo) != len(dataRow) {
		return tpl, nil, errors.New("the length of token_base_info and token_data is unmatched")
	}

	selector, err := NewTeamSelector(args)
	if err != nil {
		return tpl, nil, err
	}
	tpl.UserSelector = selector

	return tpl, dataRow, nil
}

// GenerateToken generates a token file as airdrop input.
func (tpl *BaseTemplate) GenerateToken(tokens []TokenInfo) error {
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
	f.SetCellValue(SheetClass, "A2", tpl.SheetClass.ID)
	f.SetCellValue(SheetClass, "B2", tpl.SheetClass.Name)
	f.SetCellValue(SheetClass, "C2", tpl.SheetClass.Schema)
	f.SetCellValue(SheetClass, "D2", tpl.SheetClass.Sender)
	f.SetCellValue(SheetClass, "E2", tpl.SheetClass.Symbol)
	f.SetCellValue(SheetClass, "F2", tpl.SheetClass.MintRestricted)
	f.SetCellValue(SheetClass, "G2", tpl.SheetClass.UpdateRestricted)
	f.SetCellValue(SheetClass, "H2", tpl.SheetClass.Description)
	f.SetCellValue(SheetClass, "I2", tpl.SheetClass.Uri)
	f.SetCellValue(SheetClass, "J2", tpl.SheetClass.UriHash)
	f.SetCellValue(SheetClass, "K2", tpl.SheetClass.Data)

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
		f.SetCellValue(SheetToken, fmt.Sprintf("C%d", i+2), token.Name)
		f.SetCellValue(SheetToken, fmt.Sprintf("D%d", i+2), token.URI)
		f.SetCellValue(SheetToken, fmt.Sprintf("E%d", i+2), token.Sender)
		f.SetCellValue(SheetToken, fmt.Sprintf("F%d", i+2), token.Recipient)
		f.SetCellValue(SheetToken, fmt.Sprintf("G%d", i+2), token.UriHash)
		f.SetCellValue(SheetToken, fmt.Sprintf("H%d", i+2), token.Data)
	}
	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
	return f.SaveAs(tpl.Args.OutputPath + "/tokens.xlsx")
}

func (tpl *BaseTemplate) readClass(xlsxFile *excelize.File) (Class, error) {
	rows, err := xlsxFile.GetRows(SheetClass)
	if err != nil {
		return Class{}, err
	}

	if len(rows) != 2 {
		return Class{}, errors.New("invalid class sheet, only support 2 rows")
	}

	headerRow := rows[0]
	dataRow := rows[1]
	PrintXLSX(SheetClass, headerRow, dataRow)

	mintRestricted, err := strconv.ParseBool(dataRow[4])
	if err != nil {
		return Class{}, err
	}

	updateRestricted, err := strconv.ParseBool(dataRow[5])
	if err != nil {
		return Class{}, err
	}
	return Class{
		ID:               dataRow[0],
		Name:             dataRow[1],
		Schema:           dataRow[2],
		Sender:           tpl.Args.Sender,
		Symbol:           dataRow[3],
		MintRestricted:   mintRestricted,
		UpdateRestricted: updateRestricted,
		Description:      dataRow[6],
		Uri:              dataRow[7],
		UriHash:          dataRow[8],
		Data:             dataRow[9],
	}, nil
}

func (tpl *BaseTemplate) readTokenBaseInfo(xlsxFile *excelize.File) (infos []TokenBaseInfo, err error) {
	rows, err := xlsxFile.GetRows(SheetTokenBaseInfo)
	if err != nil {
		return nil, err
	}

	headerRow := rows[0]
	dataRows := rows[1:]
	PrintXLSX(SheetTokenBaseInfo, headerRow, dataRows)

	for _, dataRow := range dataRows {
		infos = append(infos, TokenBaseInfo{
			ID:      dataRow[0],
			ClassID: dataRow[1],
			Name:    dataRow[2],
			URI:     dataRow[3],
			UriHash: dataRow[4],
		})
	}
	return
}