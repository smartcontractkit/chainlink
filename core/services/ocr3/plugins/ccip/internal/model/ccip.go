package model

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type TokenPrice struct {
	TokenID types.Account `json:"tokenID"`
	Price   BigInt        `json:"price"`
}

func NewTokenPrice(tokenID types.Account, price *big.Int) TokenPrice {
	return TokenPrice{
		TokenID: tokenID,
		Price:   BigInt{price},
	}
}

type GasPriceChain struct {
	GasPrice BigInt        `json:"gasPrice"`
	ChainSel ChainSelector `json:"chainSel"`
}

func NewGasPriceChain(gasPrice *big.Int, chainSel ChainSelector) GasPriceChain {
	return GasPriceChain{
		GasPrice: BigInt{Int: gasPrice},
		ChainSel: chainSel,
	}
}

type SeqNum uint64

func (s SeqNum) String() string {
	return strconv.FormatUint(uint64(s), 10)
}

func NewSeqNumRange(start, end SeqNum) SeqNumRange {
	return SeqNumRange{start, end}
}

type SeqNumRange [2]SeqNum

func (s SeqNumRange) Start() SeqNum {
	return s[0]
}

func (s SeqNumRange) End() SeqNum {
	return s[1]
}

func (s *SeqNumRange) SetStart(v SeqNum) {
	s[0] = v
}

func (s *SeqNumRange) SetEnd(v SeqNum) {
	s[1] = v
}

// Overlaps returns true if the two ranges overlap.
func (s SeqNumRange) Overlaps(other SeqNumRange) bool {
	return s.Start() <= other.End() && other.Start() <= s.End()
}

// Contains returns true if the range contains the given sequence number.
func (s SeqNumRange) Contains(seq SeqNum) bool {
	return s.Start() <= seq && seq <= s.End()
}

func (s SeqNumRange) String() string {
	return fmt.Sprintf("[%d -> %d]", s[0], s[1])
}

type ChainSelector uint64

func (c ChainSelector) String() string {
	ch, exists := chainselectors.ChainBySelector(uint64(c))
	if !exists || ch.Name == "" {
		return fmt.Sprintf("ChainSelector(%d)", c)
	}
	return fmt.Sprintf("%d (%s)", c, ch.Name)
}

type CCIPMsg struct {
	CCIPMsgBaseDetails
}

func (c CCIPMsg) String() string {
	js, _ := json.Marshal(c)
	return string(js)
}

type CCIPMsgBaseDetails struct {
	ID          Bytes32       `json:"id"`
	SourceChain ChainSelector `json:"sourceChain,string"`
	SeqNum      SeqNum        `json:"seqNum,string"`
}

type Evm2EvmMessage struct {
}
