package irita

import (
	sdk "github.com/bianjieai/irita-sdk-go"
	"github.com/bianjieai/irita-sdk-go/types"
	"github.com/bianjieai/irita-sdk-go/types/store"
	"github.com/smartcontractkit/chainlink/core/store/orm"
)

var (
	Client    sdk.IRITAClient
	initiated bool
	password  string
)

func GetClient(config *orm.Config) sdk.IRITAClient {
	if initiated {
		return Client
	} else {
		return newClient(config)
	}
}

func newClient(config *orm.Config) sdk.IRITAClient {
	options := []types.Option{
		types.KeyDAOOption(store.NewFileDAO(config.IritaKeyDao())),
		types.TimeoutOption(10),
	}

	cfg, err := types.NewClientConfig(
		config.IritaURL(),
		config.IritaChainID(),
		options...,
	)
	if err != nil {
		panic(err)
	}

	return sdk.NewIRITAClient(cfg)
}

func SetPassword(config *orm.Config, pwd string) error {
	if _, err := GetClient(config).Key.Show(config.IritaKeyName(), pwd); err != nil {
		return err
	}
	password = pwd
	return nil
}

func GetPassword() string {
	return password
}
