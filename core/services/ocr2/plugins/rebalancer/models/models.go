package models

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
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
	From       NetworkSelector
	To         NetworkSelector
	Amount     *big.Int
	Date       time.Time
	BridgeData []byte
	// todo: consider adding some unique id field
}

func NewTransfer(from, to NetworkSelector, amount *big.Int, date time.Time, bridgeData []byte) Transfer {
	return Transfer{
		From:       from,
		To:         to,
		Amount:     amount,
		Date:       date,
		BridgeData: bridgeData,
	}
}

func (t Transfer) Equals(other Transfer) bool {
	return t.From == other.From &&
		t.To == other.To &&
		t.Amount.Cmp(other.Amount) == 0 &&
		t.Date.Equal(other.Date) &&
		bytes.Equal(t.BridgeData, other.BridgeData)
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
