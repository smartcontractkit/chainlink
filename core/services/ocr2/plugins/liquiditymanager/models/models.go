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

func (n NetworkSelector) ChainID() uint64 {
	chainID, _ := chainsel.ChainIdFromSelector(uint64(n))
	return chainID
}

type NetworkType string

func (n NetworkSelector) Chain() (chainsel.Chain, bool) {
	return chainsel.ChainBySelector(uint64(n))
}

func (n NetworkSelector) String() string {
	chain, b := chainsel.ChainBySelector(uint64(n))
	if !b {
		return fmt.Sprintf("Unknown(%d)", n)
	}
	return chain.Name
}

// ProposedTransfer is a transfer that is proposed by the rebalancing algorithm.
type ProposedTransfer struct {
	From   NetworkSelector
	To     NetworkSelector
	Amount *ubig.Big
	Status TransferStatus
}

func (p ProposedTransfer) FromNetwork() NetworkSelector {
	return p.From
}

func (p ProposedTransfer) ToNetwork() NetworkSelector {
	return p.To
}

func (p ProposedTransfer) TransferAmount() *big.Int {
	return p.Amount.ToInt()
}

func (p ProposedTransfer) TransferStatus() TransferStatus {
	if p.Status != "" {
		return p.Status
	}
	return TransferStatusProposed
}

func (p ProposedTransfer) String() string {
	return fmt.Sprintf("from:%v to:%v amount:%s", p.From, p.To, p.Amount.String())
}

type ProposedTransfers []ProposedTransfer

func (p ProposedTransfers) Len() int      { return len(p) }
func (p ProposedTransfers) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p ProposedTransfers) Less(i, j int) bool {
	if p[i].From == p[j].From {
		return p[i].To < p[j].To
	}
	return p[i].From < p[j].From
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
	// Stage is the stage of the transfer.
	// This is primarily used for correct inflight tracking and expiry.
	// In particular, only transfers with a strictly higher stage can expire a transfer with a lower stage.
	Stage int
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

func (t Transfer) FromNetwork() NetworkSelector {
	return t.From
}

func (t Transfer) ToNetwork() NetworkSelector {
	return t.To
}

func (t Transfer) TransferAmount() *big.Int {
	return t.Amount.ToInt()
}

func (t Transfer) TransferStatus() TransferStatus {
	return TransferStatusInflight
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
		bytes.Equal(t.BridgeData, other.BridgeData) &&
		t.NativeBridgeFee.Cmp(other.NativeBridgeFee) == 0 &&
		t.Stage == other.Stage
}

func (t Transfer) String() string {
	return fmt.Sprintf("{From: %s, To: %s, Amount: %s, Sender: %s, Receiver: %s, LocalTokenAddress: %s, RemoteTokenAddress: %s, BridgeData: %s, NativeBridgeFee: %s, Stage: %d}",
		t.From.String(),
		t.To.String(),
		t.Amount.String(),
		t.Sender.String(),
		t.Receiver.String(),
		t.LocalTokenAddress.String(),
		t.RemoteTokenAddress.String(),
		hexutil.Encode(t.BridgeData),
		t.NativeBridgeFee.String(),
		t.Stage,
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

func (p PendingTransfer) FromNetwork() NetworkSelector {
	return p.Transfer.From
}

func (p PendingTransfer) ToNetwork() NetworkSelector {
	return p.Transfer.To
}

func (p PendingTransfer) TransferAmount() *big.Int {
	return p.Transfer.Amount.ToInt()
}

func (p PendingTransfer) TransferStatus() TransferStatus {
	return p.Status
}

func NewPendingTransfer(tr Transfer) PendingTransfer {
	return PendingTransfer{
		Transfer: tr,
		Status:   TransferStatusNotReady,
	}
}

type TransferStatus string

// Proposed and Inflight are used for transfers that are not yet on-chain. (not deducted from the sender on chain)
const (
	// TransferStatusProposed indicates that the transfer has been proposed by the rebalancing algorithm.
	TransferStatusProposed = "proposed"
	// TransferStatusInflight indicates that the transfer is in-flight, but has not yet been included on-chain.
	TransferStatusInflight = "inflight"
)

// the below statuses represent transfers that would have already been started on-chain. (already deducted from the sender on chain)
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
