package feeds

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func SimTransactOpts() *bind.TransactOpts {
	return &bind.TransactOpts{Signer: func(address common.Address, transaction *types.Transaction) (*types.Transaction, error) {
		return transaction, nil
	}, From: common.HexToAddress("0x0"), NoSend: true, GasLimit: 1_000_000}
}

type GnosisPayload struct {
	To   common.Address
	Data []byte
}

// TODO: probably want a batch size?
func TransferOwnership(chainState FeedsOnChainState, timelock common.Address) ([]GnosisPayload, error) {
	var payloads []GnosisPayload
	for addr, feed := range chainState.Feeds {
		payload, err := feed.TransferOwnership(SimTransactOpts(), timelock)
		if err != nil {
			return nil, err
		}
		payloads = append(payloads, GnosisPayload{
			To:   addr,
			Data: payload.Data(),
		})
	}
	return payloads, nil
}
