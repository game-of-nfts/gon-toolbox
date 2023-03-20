package individual

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/xuri/excelize/v2"
	"strings"
	"unicode"

	"github.com/game-of-nfts/gon-toolbox/nft/types"
)

type TokenData struct {
	Type        string `json:"type,omitempty"`
	Flow        string `json:"flow,class_id"`
	LastOwner   string `json:"last_owner,omitempty"`
	StartHeight string `json:"start_height,omitempty"`
}

type Template struct {
	types.BaseTemplate
	TokenData []TokenData
}

func NewTemplate(args types.InputArgs) (types.Template, error) {
	btl, err := types.NewBaseTemplate(args)
	if err != nil {
		return nil, err
	}

	tpl := &Template{
		BaseTemplate: btl,
		TokenData:    make([]TokenData, 0),
	}

	err = btl.PreInitialize(tpl.GetPreInitializer())
	if err != nil {
		return nil, err
	}

	err = btl.Initialize()
	if err != nil {
		return nil, err
	}

	if err = tpl.FillRows(btl.TokenData); err != nil {
		return nil, err
	}
	return tpl, nil
}

func (t Template) GetPreInitializer() types.PreInitializer {
	return func(tpl *types.BaseTemplate, f *excelize.File) {
		for i, user := range tpl.UserInfo() {
			tokenId := convertString(user.Github, tpl.SheetClass.Symbol)
			github := convertString(user.Github, "")
			tokenName := tokenId
			tokenUri := "https://github.com" + github
			tokenUriHash := sha256.Sum256([]byte(tokenUri))
			f.SetCellValue(types.SheetTokenBaseInfo, fmt.Sprintf("A%d", i+2),tokenId)
			f.SetCellValue(types.SheetTokenBaseInfo, fmt.Sprintf("B%d", i+2), tpl.SheetClass.ID)
			f.SetCellValue(types.SheetTokenBaseInfo, fmt.Sprintf("C%d", i+2), tokenName)
			f.SetCellValue(types.SheetTokenBaseInfo, fmt.Sprintf("D%d", i+2), tokenUri)
			f.SetCellValue(types.SheetTokenBaseInfo, fmt.Sprintf("E%d", i+2), hex.EncodeToString(tokenUriHash[:]))
		}
	}
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
			Sender:    t.Args.Sender,
			Recipient: t.NextAddress(),
			UriHash:   t.TokenBaseInfo[i].UriHash,
			Data:      string(bz),
		})
	}
	return t.GenerateToken(tokens)
}

func (t *Template) FillRows(dataRows [][]string) error {
	for _, dataRow := range dataRows {
		t.TokenData = append(t.TokenData, TokenData{
			Type:        dataRow[0],
			Flow:        dataRow[1],
			LastOwner:   dataRow[2],
			StartHeight: dataRow[3],
		})
	}
	return nil
}

func convertString(input, prefix string) string {
	input = strings.TrimSpace(input)

	var filtered bytes.Buffer
	for _, r := range input {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			filtered.WriteRune(r)
		}
	}

	return prefix + "/" + filtered.String()
}