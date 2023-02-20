package team

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/game-of-nfts/gon-toolbox/nft/types"
)

const (
	splitOne = "-->"
	splitTwo = "--"
)

type TokenData struct {
	Type          string   `json:"type,omitempty"`
	Flow          string   `json:"flow,class_id"`
	Battons       []string `json:"battons,omitempty"`
	LastRecipient string   `json:"last_recipient,omitempty"`
	StartHeight   string   `json:"start_height,omitempty"`
}

type TokenDataXLSx struct {
	Type          string `json:"type,omitempty"`
	FlowID        string `json:"flow_id,class_id"`
	Flow          string `json:"flow,class_id"`
	LastRecipient string `json:"last_recipient,omitempty"`
	StartHeight   string `json:"start_height,omitempty"`
}

type Template struct {
	types.BaseTemplate
	TokenData []TokenDataXLSx
}

func NewTemplate(args types.InputArgs) (types.Template, error) {
	baseTpl, tokenDataRows, err := types.NewTemplate(args)
	if err != nil {
		return nil, err
	}

	tpl := &Template{
		BaseTemplate: baseTpl,
		TokenData:    make([]TokenDataXLSx, 0, len(baseTpl.TokenBaseInfo)),
	}

	if err = tpl.FillTokenData(tokenDataRows); err != nil {
		return nil, err
	}
	return tpl, nil
}

func (t Template) Generate() error {
	tokens := make([]types.TokenInfo, 0, len(t.TokenData))
	for i, data := range t.TokenData {
		tokenData, recipient := t.buildTokenData(data)
		tokens = append(tokens, types.TokenInfo{
			ID:        t.TokenBaseInfo[i].ID,
			ClassID:   t.TokenBaseInfo[i].ClassID,
			Name:      t.TokenBaseInfo[i].Name,
			URI:       t.TokenBaseInfo[i].URI,
			Sender:    t.Args.Sender,
			Recipient: recipient,
			UriHash:   t.TokenBaseInfo[i].UriHash,
			Data:      tokenData,
		})
	}
	return t.GenerateToken(tokens)
}

func (t Template) buildTokenData(d TokenDataXLSx) (string, string) {
	battons, teams := t.parseFlow(d.Flow)
	data := TokenData{
		Type:          d.Type,
		Flow:          d.FlowID,
		Battons:       battons,
		LastRecipient: d.LastRecipient,
		StartHeight:   d.StartHeight,
	}
	bz, err := json.Marshal(data)
	if err != nil {
		panic(errors.New("buildTokenData failed:" + err.Error()))
	}
	return string(bz), teams[0].IRISAddress
}

func (t Template) parseFlow(flow string) (battons []string, teams []types.TeamInfo) {
	paths := strings.Split(flow, splitOne)
	teams = t.PopNTeams(len(paths))

	for i, path := range paths {
		network := strings.Split(path, splitTwo)
		switch strings.TrimSpace(network[0]) {
		// case "i":
		// 	battons = append(battons, teams[i].IRISAddress)
		case "s":
			battons = append(battons, teams[i].StargazeAddress)
		case "j":
			battons = append(battons, teams[i].JunoAddress)
		case "u":
			battons = append(battons, teams[i].UptickAddress)
		case "o":
			battons = append(battons, teams[i].UptickAddress)
		}
	}
	return
}

func (t *Template) FillTokenData(dataRows [][]string) error {
	for _, dataRow := range dataRows {
		t.TokenData = append(t.TokenData, TokenDataXLSx{
			Type:          dataRow[0],
			FlowID:        dataRow[1],
			Flow:          dataRow[2],
			LastRecipient: dataRow[3],
			StartHeight:   dataRow[4],
		})
	}
	return nil
}
