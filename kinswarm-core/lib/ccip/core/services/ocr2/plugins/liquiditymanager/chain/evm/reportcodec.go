package evmliquiditymanager

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"

	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/report_encoder"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

type OnchainReportCodec interface {
	Encode(models.Report) ([]byte, error)
	Decode(networkID models.NetworkSelector, rebalancerAddress models.Address, binaryReport []byte) (models.Report, report_encoder.ILiquidityManagerLiquidityInstructions, error) // todo: we should not use gethwrapper types
}

type EvmReportCodec struct {
	rebalancerABI          abi.ABI
	onchainReportArguments abi.Arguments
}

func NewEvmReportCodec() *EvmReportCodec {
	rebalancerABI := evmtypes.MustGetABI(report_encoder.ReportEncoderABI)
	exposeForEncoding, ok := rebalancerABI.Methods["exposeForEncoding"]
	if !ok {
		panic("exposeForEncoding method not found")
	}
	return &EvmReportCodec{
		rebalancerABI:          rebalancerABI,
		onchainReportArguments: exposeForEncoding.Inputs,
	}
}

func (e EvmReportCodec) Encode(r models.Report) ([]byte, error) {
	instructions, err := r.ToLiquidityInstructions()
	if err != nil {
		return nil, fmt.Errorf("converting to liquidity instructions: %w", err)
	}

	exposeMethod, ok := e.rebalancerABI.Methods["exposeForEncoding"]
	if !ok {
		return nil, fmt.Errorf("exposeForEncoding method not found")
	}

	encoded, err := exposeMethod.Inputs.Pack(instructions)
	if err != nil {
		return nil, fmt.Errorf("packing report: %w", err)
	}

	return encoded, nil
}

func (e EvmReportCodec) Decode(networkID models.NetworkSelector, rebalancerAddress models.Address, binaryReport []byte) (models.Report, report_encoder.ILiquidityManagerLiquidityInstructions, error) {
	unpacked, err := e.onchainReportArguments.Unpack(binaryReport)
	if err != nil {
		return models.Report{}, report_encoder.ILiquidityManagerLiquidityInstructions{}, fmt.Errorf("failed to unpack report: %w", err)
	}
	if len(unpacked) != 1 {
		return models.Report{}, report_encoder.ILiquidityManagerLiquidityInstructions{}, fmt.Errorf("unexpected number of arguments: %d", len(unpacked))
	}

	instructions := *abi.ConvertType(unpacked[0], new(report_encoder.ILiquidityManagerLiquidityInstructions)).(*report_encoder.ILiquidityManagerLiquidityInstructions)

	var out models.Report
	out.NetworkID = networkID
	out.LiquidityManagerAddress = rebalancerAddress
	for _, send := range instructions.SendLiquidityParams {
		out.Transfers = append(out.Transfers, models.Transfer{
			From:       networkID,
			To:         models.NetworkSelector(send.RemoteChainSelector),
			Amount:     ubig.New(send.Amount),
			BridgeData: send.BridgeData,
		})
	}

	for _, recv := range instructions.ReceiveLiquidityParams {
		out.Transfers = append(out.Transfers, models.Transfer{
			From:       models.NetworkSelector(recv.RemoteChainSelector),
			To:         networkID,
			Amount:     ubig.New(recv.Amount),
			BridgeData: recv.BridgeData,
		})
	}

	return out, instructions, nil
}
