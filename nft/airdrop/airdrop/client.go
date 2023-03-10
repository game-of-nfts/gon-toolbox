package airdrop

import (
	"github.com/irisnet/core-sdk-go/bank"
	"github.com/irisnet/core-sdk-go/client"
	"github.com/irisnet/core-sdk-go/common/codec"
	cdctypes "github.com/irisnet/core-sdk-go/common/codec/types"
	cryptocodec "github.com/irisnet/core-sdk-go/common/crypto/codec"
	"github.com/irisnet/core-sdk-go/gov"
	"github.com/irisnet/core-sdk-go/staking"
	"github.com/irisnet/core-sdk-go/types"
	txtypes "github.com/irisnet/core-sdk-go/types/tx"
	"github.com/irisnet/irismod-sdk-go/nft"
	"github.com/tendermint/tendermint/libs/log"
)

type Client struct {
	logger         log.Logger
	moduleManager  map[string]types.Module
	encodingConfig types.EncodingConfig

	types.BaseClient
	Bank bank.Client
	NFT  nft.Client
}

func NewClient(cfg types.ClientConfig) Client {
	encodingConfig := makeEncodingConfig()

	// create a instance of baseClient
	baseClient := client.NewBaseClient(cfg, encodingConfig, nil)
	bankClient := bank.NewClient(baseClient, encodingConfig.Marshaler)
	stakingClient := staking.NewClient(baseClient, encodingConfig.Marshaler)
	govClient := gov.NewClient(baseClient, encodingConfig.Marshaler)
	nftClient := nft.NewClient(baseClient, encodingConfig.Marshaler)

	client := &Client{
		logger:         baseClient.Logger(),
		BaseClient:     baseClient,
		moduleManager:  make(map[string]types.Module),
		encodingConfig: encodingConfig,
		Bank:           bankClient,
		NFT:            nftClient,
	}

	client.RegisterModule(
		bankClient,
		stakingClient,
		govClient,
		nftClient,
	)
	return *client
}

func (client Client) SetLogger(logger log.Logger) {
	client.BaseClient.SetLogger(logger)
}

func (client Client) Codec() *codec.LegacyAmino {
	return client.encodingConfig.Amino
}

func (client Client) AppCodec() codec.Marshaler {
	return client.encodingConfig.Marshaler
}

func (client Client) EncodingConfig() types.EncodingConfig {
	return client.encodingConfig
}

func (client Client) Manager() types.BaseClient {
	return client.BaseClient
}

func (client Client) RegisterModule(ms ...types.Module) {
	for _, m := range ms {
		m.RegisterInterfaceTypes(client.encodingConfig.InterfaceRegistry)
	}
}

func (client Client) Module(name string) types.Module {
	return client.moduleManager[name]
}

func makeEncodingConfig() types.EncodingConfig {
	amino := codec.NewLegacyAmino()
	interfaceRegistry := cdctypes.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	txCfg := txtypes.NewTxConfig(marshaler, txtypes.DefaultSignModes)

	encodingConfig := types.EncodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Marshaler:         marshaler,
		TxConfig:          txCfg,
		Amino:             amino,
	}
	RegisterLegacyAminoCodec(encodingConfig.Amino)
	RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}

// RegisterLegacyAminoCodec registers the sdk message type.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*types.Msg)(nil), nil)
	cdc.RegisterInterface((*types.Tx)(nil), nil)
	cryptocodec.RegisterCrypto(cdc)
}

// RegisterInterfaces registers the sdk message type.
func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterInterface("cosmos.v1beta1.Msg", (*types.Msg)(nil))
	txtypes.RegisterInterfaces(registry)
	cryptocodec.RegisterInterfaces(registry)
}
