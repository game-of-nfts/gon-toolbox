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

	var (
		imap = make(map[string]bool)
		smap = make(map[string]bool)
		jmap = make(map[string]bool)
		umap = make(map[string]bool)
		omap = make(map[string]bool)
	)

	for _, row := range rows[1:] {
		ValidateAddress(PrefixBech32Iris, row[1])
		ValidateAddress(PrefixBech32Stars, row[2])
		ValidateAddress(PrefixBech32Juno, row[3])
		ValidateAddress(PrefixBech32Uptick, row[4])
		ValidateAddress(PrefixBech32Omniflix, row[5])

		if imap[row[1]] {
			return nil, fmt.Errorf("duplicate address: %s", row[1])
		}

		if smap[row[2]] {
			return nil, fmt.Errorf("duplicate address: %s", row[2])
		}

		if jmap[row[3]] {
			return nil, fmt.Errorf("duplicate address: %s", row[3])
		}

		if umap[row[4]] {
			return nil, fmt.Errorf("duplicate address: %s", row[4])
		}

		if omap[row[5]] {
			return nil, fmt.Errorf("duplicate address: %s", row[5])
		}

		as.users = append(as.users, UserInfo{
			TeamName:        row[0],
			IRISAddress:     row[1],
			StargazeAddress: row[2],
			JunoAddress:     row[3],
			UptickAddress:   row[4],
			OmniflixAddress: row[5],
		})
		imap[row[1]] = true
		smap[row[2]] = true
		jmap[row[3]] = true
		umap[row[4]] = true
		omap[row[5]] = true
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
		panic(fmt.Errorf("Users are divided into groups of %d, and the current number of users is not an integer multiple of %d", n, n))
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
