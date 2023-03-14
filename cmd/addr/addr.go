package main

import (
	"fmt"
	"github.com/game-of-nfts/gon-toolbox/nft/types"
	"github.com/xuri/excelize/v2"
	"os"
	"path/filepath"
	"strings"
)

const (
	GonEvidenceRootPath = "/home/yuandu/Development/GoN/gon-evidence"
	EvidenceFile        = "evidence.xlsx"
	SheetInfo           = "teams"
	OutputPath          = "/home/yuandu/Development/GoN/address.xlsx"
)

type UserInfoPlus struct {
	types.UserInfo
	Account   string
	Discord   string
	Community string
}

func main() {
	userInfos := make([]*UserInfoPlus, 0)
	accs := make([]string, 0)

	files, err := os.ReadDir(GonEvidenceRootPath)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() {
			accs = append(accs, file.Name())
		}
	}

	// read user infos
	for _, acc := range accs {
		userInfo, err := ReadUserInfo(acc)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if userInfo != nil {
			userInfos = append(userInfos, userInfo)
		}
	}

	err = WriteUserInfo(userInfos, OutputPath)
	if err != nil {
		os.Exit(1)
	}
}

func WriteUserInfo(userInfos []*UserInfoPlus, output string) error {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Print(err)
		}
	}()

	index, err := f.NewSheet(SheetInfo)
	if err != nil {
		return err
	}
	f.SetCellValue(SheetInfo, "A1", "TeamName")
	f.SetCellValue(SheetInfo, "B1", "IRISAddress")
	f.SetCellValue(SheetInfo, "C1", "StargazeAddress")
	f.SetCellValue(SheetInfo, "D1", "JunoAddress")
	f.SetCellValue(SheetInfo, "E1", "UptickAddress")
	f.SetCellValue(SheetInfo, "F1", "OmniFlixAddress")
	f.SetCellValue(SheetInfo, "G1", "Github")
	f.SetCellValue(SheetInfo, "H1", "Discord")
	f.SetCellValue(SheetInfo, "I1", "Community")

	for i, userInfo := range userInfos {
		f.SetCellValue(SheetInfo, fmt.Sprintf("A%d", i+2), strings.TrimSpace(userInfo.TeamName))
		f.SetCellValue(SheetInfo, fmt.Sprintf("B%d", i+2), strings.TrimSpace(userInfo.IRISAddress))
		f.SetCellValue(SheetInfo, fmt.Sprintf("C%d", i+2), strings.TrimSpace(userInfo.StargazeAddress))
		f.SetCellValue(SheetInfo, fmt.Sprintf("D%d", i+2), strings.TrimSpace(userInfo.JunoAddress))
		f.SetCellValue(SheetInfo, fmt.Sprintf("E%d", i+2), strings.TrimSpace(userInfo.UptickAddress))
		f.SetCellValue(SheetInfo, fmt.Sprintf("F%d", i+2), strings.TrimSpace(userInfo.OmniflixAddress))
		f.SetCellValue(SheetInfo, fmt.Sprintf("G%d", i+2), strings.TrimSpace(userInfo.Account))
		//f.SetCellValue(SheetInfo, fmt.Sprintf("H%d", i+2), userInfo.Discord)
		//f.SetCellValue(SheetInfo, fmt.Sprintf("I%d", i+2), userInfo.Community)
	}

	f.SetActiveSheet(index)
	return f.SaveAs(output)
}

func ReadUserInfo(acc string) (*UserInfoPlus, error) {
	path := filepath.Join(GonEvidenceRootPath, acc, EvidenceFile)
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	rows, err := f.GetRows("Info")
	if err != nil {
		fmt.Printf("unable to open %s\n", path)
		return nil, err
	}

	dataRows := rows[1:]
	//community := "none"
	for _, dataRow := range dataRows {
		if dataRow == nil || dataRow[0] == "team name" {
			continue
		}

		//if len(dataRow) == 8 {
		//	community = dataRow[7]
		//}

		return &UserInfoPlus{
			UserInfo: types.UserInfo{
				TeamName:        dataRow[0],
				IRISAddress:     dataRow[1],
				StargazeAddress: dataRow[2],
				JunoAddress:     dataRow[3],
				UptickAddress:   dataRow[4],
				OmniflixAddress: dataRow[5],
			},
			Account: acc,
			//Discord:   dataRow[6],
			//Community: community,
		}, nil
	}
	return nil, nil
}
