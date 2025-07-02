package transmitter

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/por_mock_ocr3plugin/por"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

var _ ocr3types.ContractTransmitter[por.ChainSelector] = (*BasicContractTransmitter)(nil)

type BasicContractTransmitter struct {
	backend     bind.ContractTransactor
	chainID     *big.Int
	destination common.Address
	opts        *bind.TransactOpts
}

func NewBasicContractTransmitter(backend bind.ContractTransactor, chainID *big.Int, destination common.Address, privateKey ecdsa.PrivateKey) *BasicContractTransmitter {
	opts, err := bind.NewKeyedTransactorWithChainID(&privateKey, chainID)
	if err != nil {
		panic(fmt.Sprintf("bind.NewKeyedTransactorWithChainID: %v", err))
	}

	return &BasicContractTransmitter{
		backend,
		chainID,
		destination,
		opts,
	}
}

func (ct *BasicContractTransmitter) Transmit(
	ctx context.Context,
	configDigest types.ConfigDigest,
	seqNr uint64,
	reportWithInfo ocr3types.ReportWithInfo[por.ChainSelector],
	aos []types.AttributedOnchainSignature,
) error {
	// send transaction via ct.backend using ct.opts
	// ct.backend.SendTransaction(ctx, )

	// nonce, err := ct.backend.PendingNonceAt(ctx, ct.opts.From)
	// if err != nil {
	// 	return err
	// }

	// tx := ethtypes.NewTx(&ethtypes.DynamicFeeTx{
	// 	ChainID:    ct.chainID,
	// 	Nonce:      nonce,
	// 	GasTipCap:  big.NewInt(1e9),
	// 	GasFeeCap:  big.NewInt(100e9),
	// 	Gas:        100_000,
	// 	To:         &ct.destination,
	// 	Value:      big.NewInt(0),
	// 	Data:       reportWithInfo.Report,
	// 	AccessList: nil,

	// 	V: nil, R: nil, S: nil,
	// })

	// signedTx, err := ct.opts.Signer(ct.opts.From, tx)
	// if err != nil {
	// 	return err
	// }

	// err = ct.backend.SendTransaction(ct.opts.Context, signedTx)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (ct *BasicContractTransmitter) FromAccount(context.Context) (types.Account, error) {
	return types.Account(ct.opts.From.Hex()), nil
}
