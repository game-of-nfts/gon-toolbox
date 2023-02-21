package airdrop

import (
	"errors"
	"fmt"
	sdklog "github.com/irisnet/core-sdk-go/common/log"
	sdk "github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/core-sdk-go/types/store"
	"github.com/irisnet/irismod-sdk-go/nft"
	"github.com/xuri/excelize/v2"
	"strconv"
)

const (
	SheetClass         = "class"
	SheetTokenBaseInfo = "token_base_info"
)

type AirDropper struct {
	TokenXlsx          string
	SheetClass         string
	SheetTokenBaseInfo string
	Client             *Client
	Config             *Config
}

type Config struct {
	Name          string
	Password      string
	Mnemonics     string
	Gas           uint64
	Memo          string
	GasAdjustment float64

	NodeURI  string
	GRPCAddr string
	ChainID  string
}

func NewAirDropper(tokenXlsx string, generalCfg *Config) (*AirDropper, error) {
	options := []sdk.Option{
		sdk.KeyDAOOption(store.NewMemory(nil)),
		sdk.TimeoutOption(10),
	}

	clientCfg, err := sdk.NewClientConfig(generalCfg.NodeURI, generalCfg.GRPCAddr, generalCfg.ChainID, options...)
	if err != nil {
		return nil, err
	}

	client := NewClient(clientCfg)
	client.SetLogger(sdklog.NewLogger(sdklog.Config{
		Format: sdklog.FormatJSON,
		Level:  sdklog.DebugLevel,
	}))

	_, err = client.Recover(
		generalCfg.Name,
		generalCfg.Password,
		generalCfg.Mnemonics,
		"",
	)
	if err != nil {
		return nil, err
	}

	return &AirDropper{
		TokenXlsx:          tokenXlsx,
		SheetClass:         SheetClass,
		SheetTokenBaseInfo: SheetTokenBaseInfo,
		Client:             &client,
		Config:             generalCfg,
	}, err
}

func (ad *AirDropper) ExecAirdrop() error {
	f, err := excelize.OpenFile(ad.TokenXlsx)
	if err != nil {
		return nil
	}

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	baseTx := sdk.BaseTx{
		From:          ad.Config.Name,
		Password:      ad.Config.Password,
		Gas:           ad.Config.Gas,
		Memo:          ad.Config.Memo,
		Mode:          sdk.Sync,
		GasAdjustment: ad.Config.GasAdjustment,
	}

	msgs := make([]sdk.Msg, 0)
	issueMsg, err := ad.BuildIssueMsg(f)
	if err != nil {
		return err
	}
	msgs = append(msgs, issueMsg)

	mintMsgs, err := ad.BuildMintMsgs(f)
	if err != nil {
		return err
	}
	msgs = append(msgs, mintMsgs...)

	result, err := ad.Client.BuildAndSend(msgs, baseTx)
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}

func (ad *AirDropper) BuildIssueMsg(xlsxFile *excelize.File) (sdk.Msg, error) {
	rows, err := xlsxFile.GetRows(ad.SheetClass)
	if err != nil {
		return nil, err
	}

	if len(rows) != 2 {
		return nil, errors.New("invalid class sheet, only support 2 rows")
	}

	headerRow := rows[0]
	dataRow := rows[1]
	fmt.Println("header: ", headerRow)

	mintRestricted, err := strconv.ParseBool(dataRow[5])
	if err != nil {
		return nil, err
	}
	updateRestricted, err := strconv.ParseBool(dataRow[6])
	if err != nil {
		return nil, err
	}

	return &nft.MsgIssueDenom{
		Id:               dataRow[0],
		Name:             dataRow[1],
		Schema:           dataRow[2],
		Sender:           dataRow[3],
		Symbol:           dataRow[4],
		MintRestricted:   mintRestricted,
		UpdateRestricted: updateRestricted,
		Description:      dataRow[7],
		Uri:              dataRow[8],
		UriHash:          dataRow[9],
		Data:             dataRow[10],
	}, nil
}

func (ad *AirDropper) BuildMintMsgs(xlsxFile *excelize.File) (msgs []sdk.Msg, err error) {
	rows, err := xlsxFile.GetRows(ad.SheetTokenBaseInfo)
	if err != nil {
		return nil, err
	}

	headerRow := rows[0]
	fmt.Println("header: ", headerRow)

	dataRows := rows[1:]
	for _, dataRow := range dataRows {
		msgs = append(msgs, &nft.MsgMintNFT{
			Id:        dataRow[0],
			DenomId:   dataRow[1],
			Name:      dataRow[2],
			URI:       dataRow[3],
			Sender:    dataRow[4],
			Recipient: dataRow[5],
			UriHash:   dataRow[6],
			Data:      dataRow[7],
		})
	}

	return msgs, nil
}
