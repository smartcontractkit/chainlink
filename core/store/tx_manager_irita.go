package store

import (
	iservicesdk "github.com/irisnet/service-sdk-go"
	"github.com/irisnet/service-sdk-go/types"
)

type TxManagerIrita interface {
	SendTx(
		toAddr string,
		amount types.DecCoins,
		baseTx types.BaseTx,
	) (
		types.ResultTx,
		types.Error,
	)

	IritaClient() iservicesdk.ServiceClient
}

type IritaTxManager struct {
	Client iservicesdk.ServiceClient
}

func NewIritaTxManager(client iservicesdk.ServiceClient) *IritaTxManager {
	return &IritaTxManager{
		Client: client,
	}
}

func (itxm *IritaTxManager) SendTx(
	toAddr string,
	amount types.DecCoins,
	baseTx types.BaseTx,
) (
	types.ResultTx,
	types.Error,
) {
	return itxm.Client.Bank.Send(toAddr, amount, baseTx)
}

func (itxm *IritaTxManager) IritaClient() iservicesdk.ServiceClient {
	return itxm.Client
}
