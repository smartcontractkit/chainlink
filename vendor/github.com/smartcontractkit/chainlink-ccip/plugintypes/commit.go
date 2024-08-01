package plugintypes

import (
	"encoding/json"
	"fmt"
	"time"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
)

// ---[ Observation ]-----------------------------------------------------------

type CommitPluginObservation struct {
	NewMsgs     []cciptypes.RampMessageHeader   `json:"newMsgs"`
	GasPrices   []cciptypes.GasPriceChain       `json:"gasPrices"`
	TokenPrices []cciptypes.TokenPrice          `json:"tokenPrices"`
	MaxSeqNums  []SeqNumChain                   `json:"maxSeqNums"`
	FChain      map[cciptypes.ChainSelector]int `json:"fChain"`
}

func NewCommitPluginObservation(
	newMsgs []cciptypes.RampMessageHeader,
	gasPrices []cciptypes.GasPriceChain,
	tokenPrices []cciptypes.TokenPrice,
	maxSeqNums []SeqNumChain,
	FChain map[cciptypes.ChainSelector]int,
) CommitPluginObservation {
	return CommitPluginObservation{
		NewMsgs:     newMsgs,
		GasPrices:   gasPrices,
		TokenPrices: tokenPrices,
		MaxSeqNums:  maxSeqNums,
		FChain:      FChain,
	}
}

func (obs CommitPluginObservation) Encode() ([]byte, error) {
	return json.Marshal(obs)
}

func DecodeCommitPluginObservation(b []byte) (CommitPluginObservation, error) {
	obs := CommitPluginObservation{}
	err := json.Unmarshal(b, &obs)
	return obs, err
}

// ---[ Outcome ]---------------------------------------------------------------

type CommitPluginOutcome struct {
	MaxSeqNums  []SeqNumChain               `json:"maxSeqNums"`
	MerkleRoots []cciptypes.MerkleRootChain `json:"merkleRoots"`
	TokenPrices []cciptypes.TokenPrice      `json:"tokenPrices"`
	GasPrices   []cciptypes.GasPriceChain   `json:"gasPrices"`
}

func NewCommitPluginOutcome(
	seqNums []SeqNumChain,
	merkleRoots []cciptypes.MerkleRootChain,
	tokenPrices []cciptypes.TokenPrice,
	gasPrices []cciptypes.GasPriceChain,
) CommitPluginOutcome {
	return CommitPluginOutcome{
		MaxSeqNums:  seqNums,
		MerkleRoots: merkleRoots,
		TokenPrices: tokenPrices,
		GasPrices:   gasPrices,
	}
}

func (o CommitPluginOutcome) Encode() ([]byte, error) {
	return json.Marshal(o)
}

// IsEmpty returns true if the CommitPluginOutcome is empty
func (o CommitPluginOutcome) IsEmpty() bool {
	return len(o.MaxSeqNums) == 0 &&
		len(o.MerkleRoots) == 0 &&
		len(o.TokenPrices) == 0 &&
		len(o.GasPrices) == 0
}

func DecodeCommitPluginOutcome(b []byte) (CommitPluginOutcome, error) {
	o := CommitPluginOutcome{}
	err := json.Unmarshal(b, &o)
	return o, err
}

func (o CommitPluginOutcome) String() string {
	return fmt.Sprintf("{MaxSeqNums: %v, MerkleRoots: %v}", o.MaxSeqNums, o.MerkleRoots)
}

// ---[ Report ]---------------------------------------------------------------

type CommitPluginReportWithMeta struct {
	Report    cciptypes.CommitPluginReport `json:"report"`
	Timestamp time.Time                    `json:"timestamp"`
	BlockNum  uint64                       `json:"blockNum"`
}

// ---[ Generic ]--------------------------------------------------------------

type SeqNumChain struct {
	ChainSel cciptypes.ChainSelector `json:"chainSel"`
	SeqNum   cciptypes.SeqNum        `json:"seqNum"`
}

func NewSeqNumChain(chainSel cciptypes.ChainSelector, seqNum cciptypes.SeqNum) SeqNumChain {
	return SeqNumChain{
		ChainSel: chainSel,
		SeqNum:   seqNum,
	}
}
