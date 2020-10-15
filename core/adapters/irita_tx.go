package adapters

import (
	"github.com/irisnet/service-sdk-go/service"
	"github.com/irisnet/service-sdk-go/types"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/irita"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

type IritaTx struct {
	BaseTx types.BaseTx
	ToAddr string
	Amount types.DecCoins
}

func (itx *IritaTx) TaskType() models.TaskType {
	return TaskTypeIritaTx
}

func (itx *IritaTx) Perform(input models.RunInput, store *strpkg.Store) models.RunOutput {
	baseTx := types.BaseTx{
		From:     store.Config.IritaKeyName(),
		Gas:      200000,
		Memo:     "respond by chainlink provider",
		Mode:     types.Sync,
		Password: irita.GetPassword(),
	}

	serviceRequset := strpkg.GetServiceMemory()[input.JobRunID().String()]
	strpkg.DeleteFromMemory(input.JobRunID().String())

	provider, _ := types.AccAddressFromBech32(serviceRequset.Provider)

	msgs := []types.Msg{&service.MsgRespondService{
		RequestId: types.MustHexBytesFrom(serviceRequset.RequestResponse.ID),
		Provider:  provider,
		Output:    input.Result().String(),
		Result:    `{"code":200,"message":""}`,
	}}

	result, err := store.TxManagerIrita.IritaClient().SendBatch(msgs, baseTx)
	if err != nil {
		println(err.Error())
		return models.NewRunOutputCompleteWithResult(err)
	}

	logger.Info(input.Result().String())
	return models.NewRunOutputCompleteWithResult(result)
}
