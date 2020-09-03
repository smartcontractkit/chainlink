package store

import (
	sdk "github.com/bianjieai/irita-sdk-go"
	"github.com/bianjieai/irita-sdk-go/types"
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

	IritaClient() sdk.IRITAClient
}

type IritaTxManager struct {
	Client sdk.IRITAClient
}

func NewIritaTxManager(client sdk.IRITAClient) *IritaTxManager {
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

func (itxm *IritaTxManager) IritaClient() sdk.IRITAClient {
	return itxm.Client
}
