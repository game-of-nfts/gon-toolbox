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

type TeamInfo struct {
	TeamName        string `json:"team_name"`
	IRISAddress     string `json:"iris_address"`
	StargazeAddress string `json:"stargaze_address"`
	JunoAddress     string `json:"juno_address"`
	UptickAddress   string `json:"uptick_address"`
	OmniflixAddress string `json:"omniflix_address"`
}

type TeamSelector struct {
	teams []TeamInfo
}

func NewTeamSelector(args InputArgs) (*TeamSelector, error) {
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

	as := &TeamSelector{
		teams: make([]TeamInfo, 0, 0),
	}

	rows, err := f.GetRows(ChainIDiris)
	if err != nil {
		return nil, err
	}

	for _, row := range rows[1:] {
		as.teams = append(as.teams, TeamInfo{
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

func (as *TeamSelector) PopOneAddress() string {
	if len(as.teams) == 0 {
		panic("no available address")
	}
	selectIdx := rand.Intn(len(as.teams))
	address := as.teams[selectIdx].IRISAddress

	as.teams = append(
		as.teams[0:selectIdx],
		as.teams[selectIdx:]...,
	)
	return address
}

func (as *TeamSelector) PopNTeams(n int) (teams []TeamInfo) {
	if len(as.teams)%n != 0 {
		panic("no available address")
	}

	for i := 0; i < n; i++ {
		selectIdx := rand.Intn(len(as.teams))
		teams = append(teams, as.teams[selectIdx])

		as.teams = append(
			as.teams[0:selectIdx],
			as.teams[selectIdx:]...,
		)
	}
	return
}

func (as *TeamSelector) AllTeams() []TeamInfo {
	return as.teams
}
