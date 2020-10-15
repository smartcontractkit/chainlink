package irita

import (
	iservicesdk "github.com/irisnet/service-sdk-go"
	"github.com/irisnet/service-sdk-go/types"
	"github.com/irisnet/service-sdk-go/types/store"
	"github.com/smartcontractkit/chainlink/core/store/orm"
)

var (
	Client    iservicesdk.ServiceClient
	initiated bool
	password  string
)

func GetClient(config *orm.Config) iservicesdk.ServiceClient {
	if initiated {
		return Client
	} else {
		return newClient(config)
	}
}

func newClient(config *orm.Config) iservicesdk.ServiceClient {
	options := []types.Option{
		types.KeyDAOOption(store.NewFileDAO(config.IritaKeyDao())),
		types.TimeoutOption(10),
	}

	cfg, err := types.NewClientConfig(
		config.IritaURL(),
		config.IritaGRPCAddr(),
		config.IritaChainID(),
		options...,
	)
	if err != nil {
		panic(err)
	}

	return iservicesdk.NewServiceClient(cfg)
}

func SetPassword(config *orm.Config, pwd string) error {
	if _, _, err := GetClient(config).Find(config.IritaKeyName(), pwd); err != nil {
		return err
	}
	password = pwd
	return nil
}

func GetPassword() string {
	return password
}
