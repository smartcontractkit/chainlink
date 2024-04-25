package model

import (
	"encoding/json"
	"fmt"

	chainselectors "github.com/smartcontractkit/chain-selectors"
)

type SeqNum uint64

type SeqNumRange [2]SeqNum

func (s SeqNumRange) Start() SeqNum {
	return s[0]
}

func (s SeqNumRange) End() SeqNum {
	return s[1]
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

type NodeID string

type CCIPMsg struct {
	CCIPMsgBaseDetails
}

func (c CCIPMsg) String() string {
	js, _ := json.Marshal(c)
	return fmt.Sprintf("%s", js)
}

type CCIPMsgBaseDetails struct {
	SourceChain ChainSelector `json:"sourceChain,string"`
	SeqNum      SeqNum        `json:"seqNum,string"`
}
