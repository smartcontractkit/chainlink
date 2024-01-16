package models

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

type Address common.Address

func (a *Address) String() string {
	return common.Address(*a).Hex()
}

func (a *Address) UnmarshalJSON(input []byte) error {
	ta := common.Address(*a)
	err := ta.UnmarshalJSON(input)
	if err != nil {
		return err
	}
	*a = Address(ta)
	return nil
}

func (a Address) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, common.Address(a).Hex())), nil
}

type NetworkSelector uint64

const (
	NetworkTypeUnknown = "unknown"
	NetworkTypeEvm     = "evm"
	NetworkTypeSolana  = "sol"
)

func (n NetworkSelector) Type() NetworkType {
	switch n {
	case 1, 2, 3, 1337, 1338, 1339, 1340: // todo: use some lib
		return NetworkTypeEvm
	case 4:
		return NetworkTypeSolana
	default:
		return NetworkTypeUnknown
	}
}

func (n *NetworkSelector) UnmarshalJSON(input []byte) error {
	var i uint64
	err := json.Unmarshal(input, &i)
	if err != nil {
		return err
	}
	*n = NetworkSelector(i)
	return nil
}

func (n NetworkSelector) MarshalJSON() ([]byte, error) {
	return json.Marshal(uint64(n))
}

type NetworkType string

type Transfer struct {
	From   NetworkSelector
	To     NetworkSelector
	Amount *big.Int
	Date   time.Time
	// todo: consider adding some unique id field
}

func NewTransfer(from, to NetworkSelector, amount *big.Int, date time.Time) Transfer {
	return Transfer{
		From:   from,
		To:     to,
		Amount: amount,
		Date:   date,
	}
}

func (t Transfer) Equals(other Transfer) bool {
	return t.From == other.From &&
		t.To == other.To &&
		t.Amount.Cmp(other.Amount) == 0 &&
		t.Date.Equal(other.Date)
}

type PendingTransfer struct {
	Transfer
	Status TransferStatus
}

func (p PendingTransfer) Hash() ([32]byte, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return [32]byte{}, fmt.Errorf("marshal: %w", err)
	}
	return sha256.Sum256(b), nil
}

func (p PendingTransfer) String() string {
	return fmt.Sprintf("PendingTransfer{Transfer: %s, Status: %s}", p.Transfer.String(), p.Status)
}

func NewPendingTransfer(tr Transfer) PendingTransfer {
	return PendingTransfer{
		Transfer: tr,
		Status:   TransferStatusNotReady,
	}
}

type TransferStatus string

const (
	TransferStatusNotReady  = "not-ready"
	TransferStatusReady     = "ready"
	TransferStatusFinalized = "finalized"
	TransferStatusExecuted  = "executed"
)

func (t Transfer) String() string {
	return fmt.Sprintf("%v->%v %s", t.From, t.To, t.Amount.String())
}

type ReportMetadata struct {
	Transfers               []Transfer
	LiquidityManagerAddress Address
	NetworkID               NetworkSelector
}

func NewReportMetadata(transfers []Transfer, lmAddr Address, networkID NetworkSelector) ReportMetadata {
	return ReportMetadata{
		Transfers:               transfers,
		LiquidityManagerAddress: lmAddr,
		NetworkID:               networkID,
	}
}

func (r ReportMetadata) Encode() []byte {
	b, err := json.Marshal(r)
	if err != nil {
		panic(fmt.Errorf("report meta %#v encoding unexpected internal error: %w", r, err))
	}
	return b
}

func (r ReportMetadata) OnchainEncode() ([]byte, error) {
	return onchainReportArguments.Pack(big.NewInt(int64(r.NetworkID)), common.Address(r.LiquidityManagerAddress))
}

func (r ReportMetadata) GetDestinationChain() relay.ID {
	return relay.NewID(relay.EVM, fmt.Sprintf("%d", r.NetworkID))
}

func (r ReportMetadata) String() string {
	return fmt.Sprintf("ReportMetadata{Transfers: %v, LiquidityManagerAddress: %s, NetworkID: %d}", r.Transfers, r.LiquidityManagerAddress, r.NetworkID)
}

func DecodeReportMetadata(b []byte) (ReportMetadata, error) {
	var meta ReportMetadata
	err := json.Unmarshal(b, &meta)
	return meta, err
}

func DecodeReport(b []byte) (ReportMetadata, error) {
	unpacked, err := onchainReportArguments.Unpack(b)
	if err != nil {
		return ReportMetadata{}, fmt.Errorf("failed to unpack report: %w", err)
	}
	if len(unpacked) != 2 {
		return ReportMetadata{}, fmt.Errorf("unexpected number of arguments: %d", len(unpacked))
	}
	var out ReportMetadata
	chainID := *abi.ConvertType(unpacked[0], new(*big.Int)).(**big.Int)
	lqmgrAddr := *abi.ConvertType(unpacked[1], new(common.Address)).(*common.Address)
	out.NetworkID = NetworkSelector(chainID.Int64())
	out.LiquidityManagerAddress = Address(lqmgrAddr)
	return out, nil
}

var onchainReportArguments = ReportSchema()

func ReportSchema() abi.Arguments {
	mustNewType := func(t string) abi.Type {
		result, err := abi.NewType(t, "", []abi.ArgumentMarshaling{})
		if err != nil {
			panic(fmt.Sprintf("Unexpected error during abi.NewType: %s", err))
		}
		return result
	}
	return abi.Arguments([]abi.Argument{
		{
			Name: "chainId",
			Type: mustNewType("uint256"),
		},
		{
			Name: "liquidityManagerAddress",
			Type: mustNewType("address"),
		},
		// TODO: add transfer operations
		// Once the onchain code is merged we can remove this function
		// and just use the gethwrapper
	})
}
