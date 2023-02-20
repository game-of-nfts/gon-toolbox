package types

import (
	"errors"
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

type AddressSelector struct {
	m map[string][]string
}

func NewAddressSelector(args InputArgs) (*AddressSelector, error) {
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

	as := &AddressSelector{
		m: map[string][]string{
			ChainIDiris:     make([]string, 0, 0),
			ChainIDstargaze: make([]string, 0, 0),
			ChainIDjuno:     make([]string, 0, 0),
			ChainIDuptick:   make([]string, 0, 0),
			ChainIDflixnet:  make([]string, 0, 0),
		},
	}

	//gon-irishub-1
	as.fill(ChainIDiris, f)
	//elgafar-1
	as.fill(ChainIDstargaze, f)
	//uni-6
	as.fill(ChainIDjuno, f)
	//uptick_7001_1
	as.fill(ChainIDuptick, f)
	//gon-flixnet-1
	as.fill(ChainIDflixnet, f)
	return as, nil
}

func (as *AddressSelector) fill(chainID string, xlsxFile *excelize.File) error {
	rows, err := xlsxFile.GetRows(ChainIDiris)
	if err != nil {
		return err
	}
	for _, row := range rows[1:] {
		as.m[chainID] = append(as.m[chainID], row[0])
	}
	return nil
}

func (as *AddressSelector) Pop(chainID string) (string, error) {
	if len(as.m[chainID]) == 0 {
		return "", errors.New("no available address")
	}
	selectIdx := rand.Intn(len(as.m[chainID]))
	address := as.m[chainID][selectIdx]

	as.m[chainID] = append(
		as.m[chainID][0:selectIdx],
		as.m[chainID][selectIdx:]...,
	)
	return address, nil
}
