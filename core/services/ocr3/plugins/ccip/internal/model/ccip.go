package model

import (
	"encoding/json"
	"fmt"
	"math/big"

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

type GasPrice *big.Int

type GasPriceChain struct {
	GasPrice GasPrice
	ChainSel ChainSelector
}

func NewGasPriceChain(gasPrice GasPrice, chainSel ChainSelector) GasPriceChain {
	return GasPriceChain{
		GasPrice: gasPrice,
		ChainSel: chainSel,
	}
}

type SeqNum uint64

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
	return fmt.Sprintf("%s", js)
}

type CCIPMsgBaseDetails struct {
	ID          Bytes32       `json:"id"`
	SourceChain ChainSelector `json:"sourceChain,string"`
	SeqNum      SeqNum        `json:"seqNum,string"`
}
