package models

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	chainsel "github.com/smartcontractkit/chain-selectors"

	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
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
	isEvm, err := chainsel.IsEvm(uint64(n))
	if err == nil && isEvm {
		return NetworkTypeEvm
	}

	return NetworkTypeUnknown
}

type NetworkType string

// ProposedTransfer is a transfer that is proposed by the rebalancing algorithm.
type ProposedTransfer struct {
	From   NetworkSelector
	To     NetworkSelector
	Amount *ubig.Big
}

// Transfer is a ProposedTransfer that has had a lot of its information resolved.
type Transfer struct {
	// From identifies the network where the tokens are originating from.
	From NetworkSelector
	// To identifies the network where the tokens are headed to.
	To NetworkSelector
	// Sender is an address on the From network that is sending the tokens.
	// Typically this will be the rebalancer contract on the From network.
	Sender Address
	// Receiver is an address on the To network that will be receiving the tokens.
	// Typically this will be the rebalancer contract on the To network.
	Receiver Address
	// LocalTokenAddress is the address of the token on the From network.
	LocalTokenAddress Address
	// RemoteTokenAddress is the address of the token on the To network.
	RemoteTokenAddress Address
	// Amount is the amount of tokens being transferred.
	Amount *ubig.Big
	// Date is the date when the transfer was initiated.
	Date time.Time
	// BridgeData is any additional data that needs to be sent to the bridge
	// in order to make the transfer succeed.
	BridgeData hexutil.Bytes
	// NativeBridgeFee is the fee that the bridge charges for the transfer.
	NativeBridgeFee *ubig.Big
	// todo: consider adding some unique id field
}

func NewTransfer(from, to NetworkSelector, amount *big.Int, date time.Time, bridgeData []byte) Transfer {
	return Transfer{
		From:       from,
		To:         to,
		Amount:     ubig.New(amount),
		Date:       date,
		BridgeData: bridgeData,
	}
}

func (t Transfer) Equals(other Transfer) bool {
	return t.From == other.From &&
		t.To == other.To &&
		t.Sender == other.Sender &&
		t.Receiver == other.Receiver &&
		t.LocalTokenAddress == other.LocalTokenAddress &&
		t.RemoteTokenAddress == other.RemoteTokenAddress &&
		t.Amount.Cmp(other.Amount) == 0 &&
		t.Date.Equal(other.Date) &&
		bytes.Equal(t.BridgeData, other.BridgeData)
}

func (t Transfer) String() string {
	return fmt.Sprintf("{From: %d, To: %d, Amount: %s, Sender: %s, Receiver: %s, LocalTokenAddress: %s, RemoteTokenAddress: %s, BridgeData: %s, NativeBridgeFee: %s}",
		t.From,
		t.To,
		t.Amount.String(),
		t.Sender.String(),
		t.Receiver.String(),
		t.LocalTokenAddress.String(),
		t.RemoteTokenAddress.String(),
		hexutil.Encode(t.BridgeData),
		t.NativeBridgeFee.String(),
	)
}

// PendingTransfer is a Transfer whose status has been resolved.
type PendingTransfer struct {
	Transfer
	Status TransferStatus
	ID     string
}

func (p PendingTransfer) Hash() ([32]byte, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return [32]byte{}, fmt.Errorf("marshal: %w", err)
	}
	return sha256.Sum256(b), nil
}

func (p PendingTransfer) String() string {
	return fmt.Sprintf("PendingTransfer{ID: %s, Transfer: %s, Status: %s}", p.ID, p.Transfer.String(), p.Status)
}

func NewPendingTransfer(tr Transfer) PendingTransfer {
	return PendingTransfer{
		Transfer: tr,
		Status:   TransferStatusNotReady,
	}
}

type TransferStatus string

const (
	// TransferStatusNotReady indicates that the transfer is in-flight, but has either not been auto-finalized (e.g L1 -> L2 transfers)
	// or is not ready to finalize on-chain (e.g L2 -> L1 transfers).
	TransferStatusNotReady = "not-ready"
	// TransferStatusReady indicates that the transfer is in-flight but ready to be finalized on-chain.
	TransferStatusReady = "ready"
	// TransferStatusFinalized indicates that the transfer has been finalized on-chain, but not yet executed.
	TransferStatusFinalized = "finalized"
	// TransferStatusExecuted indicates that the transfer has been finalized and executed. This is a terminal state.
	TransferStatusExecuted = "executed"
)

type Edge struct {
	Source NetworkSelector
	Dest   NetworkSelector
}

func NewEdge(source, dest NetworkSelector) Edge {
	return Edge{
		Source: source,
		Dest:   dest,
	}
}
