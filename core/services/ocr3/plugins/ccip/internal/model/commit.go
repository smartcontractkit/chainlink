package model

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type CommitPluginObservation struct {
	NewMsgs     []CCIPMsgBaseDetails `json:"newMsgs"`
	GasPrices   []GasPriceChain      `json:"gasPrices,string"`
	TokenPrices []TokenPrice         `json:"tokenPrices"`
	MaxSeqNums  []SeqNumChain        `json:"maxSeqNums"`
}

func NewCommitPluginObservation(
	newMsgs []CCIPMsgBaseDetails,
	gasPrices []GasPriceChain,
	tokenPrices []TokenPrice,
	maxSeqNums []SeqNumChain,
) CommitPluginObservation {
	return CommitPluginObservation{
		NewMsgs:     newMsgs,
		GasPrices:   gasPrices,
		TokenPrices: tokenPrices,
		MaxSeqNums:  maxSeqNums,
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

type CommitPluginOutcome struct {
	MaxSeqNums  []SeqNumChain     `json:"maxSeqNums"`
	MerkleRoots []MerkleRootChain `json:"merkleRoots"`
}

func NewCommitPluginOutcome(seqNums []SeqNumChain, merkleRoots []MerkleRootChain) CommitPluginOutcome {
	return CommitPluginOutcome{
		MaxSeqNums:  seqNums,
		MerkleRoots: merkleRoots,
	}
}

func (o CommitPluginOutcome) Encode() ([]byte, error) {
	return json.Marshal(o)
}

func DecodeCommitPluginOutcome(b []byte) (CommitPluginOutcome, error) {
	o := CommitPluginOutcome{}
	err := json.Unmarshal(b, &o)
	return o, err
}

func (o CommitPluginOutcome) String() string {
	return fmt.Sprintf("{MaxSeqNums: %v, MerkleRoots: %v}", o.MaxSeqNums, o.MerkleRoots)
}

type SeqNumChain struct {
	ChainSel ChainSelector `json:"chainSel"`
	SeqNum   SeqNum        `json:"seqNum"`
}

func NewSeqNumChain(chainSel ChainSelector, seqNum SeqNum) SeqNumChain {
	return SeqNumChain{
		ChainSel: chainSel,
		SeqNum:   seqNum,
	}
}

type MerkleRootChain struct {
	ChainSel     ChainSelector `json:"chain"`
	SeqNumsRange SeqNumRange   `json:"seqNumsRange"`
	MerkleRoot   Bytes32       `json:"merkleRoot"`
}

func NewMerkleRootChain(chainSel ChainSelector, seqNumsRange SeqNumRange, merkleRoot Bytes32) MerkleRootChain {
	return MerkleRootChain{
		ChainSel:     chainSel,
		SeqNumsRange: seqNumsRange,
		MerkleRoot:   merkleRoot,
	}
}

type CommitPluginReport struct {
	MerkleRoots  []MerkleRootChain  `json:"merkleRoots"`
	PriceUpdates []TokenPriceUpdate `json:"priceUpdates"`
}

func NewCommitPluginReport(merkleRoots []MerkleRootChain, priceUpdates []TokenPriceUpdate) CommitPluginReport {
	return CommitPluginReport{
		MerkleRoots:  merkleRoots,
		PriceUpdates: priceUpdates,
	}
}

// IsEmpty returns true if the CommitPluginReport is empty
func (r CommitPluginReport) IsEmpty() bool {
	return len(r.MerkleRoots) == 0 &&
		len(r.PriceUpdates) == 0
}

type TokenPriceUpdate struct {
	TokenID types.Account `json:"tokenID"`
	Price   *big.Int      `json:"price"`
}
