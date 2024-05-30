package model

import (
	"encoding/json"
	"fmt"
)

type CommitPluginObservation struct {
	NewMsgs      []CCIPMsgBaseDetails `json:"newMsgs"`
	GasPrices    []GasPriceChain      `json:"gasPrices,string"`
	TokenPrices  []TokenPrice         `json:"tokenPrices"`
	MaxSeqNums   []SeqNumChain        `json:"maxSeqNums"`
	PluginConfig CommitPluginConfig   `json:"pluginConfig"`
}

func NewCommitPluginObservation(
	newMsgs []CCIPMsgBaseDetails,
	gasPrices []GasPriceChain,
	tokenPrices []TokenPrice,
	maxSeqNums []SeqNumChain,
	pluginConfig CommitPluginConfig,
) CommitPluginObservation {
	return CommitPluginObservation{
		NewMsgs:      newMsgs,
		GasPrices:    gasPrices,
		TokenPrices:  tokenPrices,
		MaxSeqNums:   maxSeqNums,
		PluginConfig: pluginConfig,
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
	TokenPrices []TokenPrice      `json:"tokenPrices"`
	GasPrices   []GasPriceChain   `json:"gasPrices"`
}

func NewCommitPluginOutcome(
	seqNums []SeqNumChain,
	merkleRoots []MerkleRootChain,
	tokenPrices []TokenPrice,
	gasPrices []GasPriceChain,
) CommitPluginOutcome {
	return CommitPluginOutcome{
		MaxSeqNums:  seqNums,
		MerkleRoots: merkleRoots,
		TokenPrices: tokenPrices,
		GasPrices: gasPrices,
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

type PriceUpdate struct {
	TokenPriceUpdates []TokenPrice `json:"tokenPriceUpdates"`
	GasPriceUpdates  []GasPriceChain   `json:"gasPriceUpdates"`
}

type CommitPluginReport struct {
	MerkleRoots       []MerkleRootChain `json:"merkleRoots"`
	PriceUpdates	  PriceUpdate       `json:"priceUpdates"`
}

// TODO: also accept gas prices slice in the constructor
func NewCommitPluginReport(merkleRoots []MerkleRootChain, tokenPriceUpdates []TokenPrice, gasPriceUpdate []GasPriceChain) CommitPluginReport {
	return CommitPluginReport{
		MerkleRoots:       merkleRoots,
		PriceUpdates: 	  PriceUpdate{TokenPriceUpdates: tokenPriceUpdates, GasPriceUpdates: gasPriceUpdate},
	}
}

// IsEmpty returns true if the CommitPluginReport is empty
func (r CommitPluginReport) IsEmpty() bool {
	return len(r.MerkleRoots) == 0 && len(r.PriceUpdates.TokenPriceUpdates) == 0 && len(r.PriceUpdates.GasPriceUpdates) == 0
}
