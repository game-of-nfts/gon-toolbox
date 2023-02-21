package types

import (
	"fmt"
	"math/rand"

	"github.com/xuri/excelize/v2"
)

const (
	ChainIDiris     = "gon-irishub-1"
	ChainIDstargaze = "elgafar-1"
	ChainIDjuno     = "uni-6"
	ChainIDuptick   = "uptick_7001-1"
	ChainIDflixnet  = "gon-flixnet-1"
)

type UserInfo struct {
	TeamName        string `json:"team_name"`
	IRISAddress     string `json:"iris_address"`
	StargazeAddress string `json:"stargaze_address"`
	JunoAddress     string `json:"juno_address"`
	UptickAddress   string `json:"uptick_address"`
	OmniflixAddress string `json:"omniflix_address"`
}

type UserSelector struct {
	users []UserInfo
}

func NewTeamSelector(args InputArgs) (*UserSelector, error) {
	if len(args.AddressFile) == 0 {
		return nil, nil
	}

	f, err := excelize.OpenFile(args.AddressFile)
	if err != nil {
		return nil, err
	}

	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	as := &UserSelector{
		users: make([]UserInfo, 0, 0),
	}

	rows, err := f.GetRows(SheetTeams)
	if err != nil {
		return nil, err
	}

	for _, row := range rows[1:] {
		as.users = append(as.users, UserInfo{
			TeamName:        row[0],
			IRISAddress:     row[1],
			StargazeAddress: row[2],
			JunoAddress:     row[3],
			UptickAddress:   row[4],
			OmniflixAddress: row[5],
		})
	}
	return as, nil
}

func (as *UserSelector) PopAddress() string {
	if len(as.users) == 0 {
		panic("no available address")
	}
	selectIdx := rand.Intn(len(as.users))
	address := as.users[selectIdx].IRISAddress

	as.users = append(
		as.users[0:selectIdx],
		as.users[selectIdx:]...,
	)
	return address
}

func (as *UserSelector) PopNUsers(n int) (users []UserInfo) {
	if len(as.users)%n != 0 {
		panic("no available address")
	}

	for i := 0; i < n; i++ {
		selectIdx := rand.Intn(len(as.users))
		users = append(users, as.users[selectIdx])

		as.users = append(
			as.users[0:selectIdx],
			as.users[selectIdx:]...,
		)
	}
	return
}
