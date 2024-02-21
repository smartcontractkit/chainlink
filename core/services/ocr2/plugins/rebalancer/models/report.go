package models

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	chainsel "github.com/smartcontractkit/chain-selectors"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/rebalancer/generated/rebalancer_report_encoder"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

var (
	rebalancerABI          = evmtypes.MustGetABI(rebalancer_report_encoder.RebalancerReportEncoderABI)
	onchainReportArguments abi.Arguments
)

func init() {
	exposeForEncoding, ok := rebalancerABI.Methods["exposeForEncoding"]
	if !ok {
		panic("exposeForEncoding method not found")
	}
	onchainReportArguments = exposeForEncoding.Inputs
}

// ConfigDigest wraps ocrtypes.ConfigDigest and adds json encoding support
type ConfigDigest struct {
	ocrtypes.ConfigDigest
}

func (c *ConfigDigest) UnmarshalJSON(b []byte) error {
	var hexStr string
	if err := json.Unmarshal(b, &hexStr); err != nil {
		return err
	}
	b, err := hex.DecodeString(hexStr)
	if err != nil {
		return fmt.Errorf("decoding config digest hex string: %w", err)
	}
	if len(b) != 32 {
		return fmt.Errorf("config digest must be 32 bytes, got %d bytes", len(b))
	}
	var digestBytes [32]byte
	copy(digestBytes[:], b)
	*c = ConfigDigest{
		ConfigDigest: ocrtypes.ConfigDigest(digestBytes),
	}
	return nil
}

func (c ConfigDigest) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(c.ConfigDigest[:]))
}

func (c ConfigDigest) ToOCRConfigDigest() ocrtypes.ConfigDigest {
	return c.ConfigDigest
}

type Report struct {
	Transfers               []Transfer
	LiquidityManagerAddress Address

	// NetworkID is the network ID that this report is going to be posted to
	NetworkID NetworkSelector
	// ConfigDigest is the latest config digest of the contract that the report
	// is going to be posted to.
	ConfigDigest ConfigDigest
}

func NewReport(transfers []Transfer, lmAddr Address, networkID NetworkSelector, configDigest ocrtypes.ConfigDigest) Report {
	return Report{
		Transfers:               transfers,
		LiquidityManagerAddress: lmAddr,
		NetworkID:               networkID,
		ConfigDigest: ConfigDigest{
			ConfigDigest: configDigest,
		},
	}
}

func (r Report) Encode() []byte {
	b, err := json.Marshal(r)
	if err != nil {
		panic(fmt.Errorf("report meta %#v encoding unexpected internal error: %w", r, err))
	}
	return b
}

func (r Report) OnchainEncode() ([]byte, error) {
	instructions, err := r.ToLiquidityInstructions()
	if err != nil {
		return nil, fmt.Errorf("converting to liquidity instructions: %w", err)
	}
	exposeMethod, ok := rebalancerABI.Methods["exposeForEncoding"]
	if !ok {
		return nil, fmt.Errorf("exposeForEncoding method not found")
	}
	encoded, err := exposeMethod.Inputs.Pack(instructions)
	if err != nil {
		return nil, fmt.Errorf("packing report: %w", err)
	}
	return encoded, nil
}

func (r Report) ToLiquidityInstructions() (rebalancer_report_encoder.IRebalancerLiquidityInstructions, error) {
	var sendInstructions []rebalancer_report_encoder.IRebalancerSendLiquidityParams
	var receiveInstructions []rebalancer_report_encoder.IRebalancerReceiveLiquidityParams
	for _, tr := range r.Transfers {
		if r.NetworkID == tr.From {
			sendInstructions = append(sendInstructions, rebalancer_report_encoder.IRebalancerSendLiquidityParams{
				Amount:              tr.Amount.ToInt(),
				RemoteChainSelector: uint64(tr.To),
				BridgeData:          tr.BridgeData,
				NativeBridgeFee:     tr.NativeBridgeFee.ToInt(),
			})
		} else if r.NetworkID == tr.To {
			receiveInstructions = append(receiveInstructions, rebalancer_report_encoder.IRebalancerReceiveLiquidityParams{
				Amount:              tr.Amount.ToInt(),
				RemoteChainSelector: uint64(tr.From),
				BridgeData:          tr.BridgeData,
			})
		} else {
			return rebalancer_report_encoder.IRebalancerLiquidityInstructions{},
				fmt.Errorf("transfer %+v is not related to network %d", tr, r.NetworkID)
		}
	}
	return rebalancer_report_encoder.IRebalancerLiquidityInstructions{
		SendLiquidityParams:    sendInstructions,
		ReceiveLiquidityParams: receiveInstructions,
	}, nil
}

func (r Report) GetDestinationChain() relay.ID {
	networkID := r.NetworkID

	ch, exists := chainsel.ChainBySelector(uint64(r.NetworkID))
	if exists {
		networkID = NetworkSelector(ch.EvmChainID)
	}

	return relay.NewID(relay.EVM, fmt.Sprintf("%d", networkID))
}

func (r Report) GetDestinationConfigDigest() ocrtypes.ConfigDigest {
	return r.ConfigDigest.ToOCRConfigDigest()
}

func (r Report) String() string {
	return fmt.Sprintf("Report{Transfers: %v, RebalancerAddress: %s, NetworkID: %d, ConfigDigest: %s}",
		r.Transfers, r.LiquidityManagerAddress.String(), r.NetworkID, r.ConfigDigest.Hex())
}

func DecodeReportMetadata(b []byte) (Report, error) {
	var meta Report
	err := json.Unmarshal(b, &meta)
	return meta, err
}

func DecodeReport(networkID NetworkSelector, rebalancerAddress Address, binaryReport []byte) (Report, rebalancer_report_encoder.IRebalancerLiquidityInstructions, error) {
	unpacked, err := onchainReportArguments.Unpack(binaryReport)
	if err != nil {
		return Report{}, rebalancer_report_encoder.IRebalancerLiquidityInstructions{}, fmt.Errorf("failed to unpack report: %w", err)
	}
	if len(unpacked) != 1 {
		return Report{}, rebalancer_report_encoder.IRebalancerLiquidityInstructions{}, fmt.Errorf("unexpected number of arguments: %d", len(unpacked))
	}

	instructions := *abi.ConvertType(unpacked[0], new(rebalancer_report_encoder.IRebalancerLiquidityInstructions)).(*rebalancer_report_encoder.IRebalancerLiquidityInstructions)

	var out Report
	out.NetworkID = networkID
	out.LiquidityManagerAddress = rebalancerAddress
	for _, send := range instructions.SendLiquidityParams {
		out.Transfers = append(out.Transfers, Transfer{
			From:       networkID,
			To:         NetworkSelector(send.RemoteChainSelector),
			Amount:     ubig.New(send.Amount),
			BridgeData: send.BridgeData,
		})
	}

	for _, recv := range instructions.ReceiveLiquidityParams {
		out.Transfers = append(out.Transfers, Transfer{
			From:       NetworkSelector(recv.RemoteChainSelector),
			To:         networkID,
			Amount:     ubig.New(recv.Amount),
			BridgeData: recv.BridgeData,
		})
	}

	return out, instructions, nil
}
