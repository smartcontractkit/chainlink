package llo

import (
	"fmt"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"google.golang.org/protobuf/proto"
)

type OffchainConfig struct {
	// We use the offchainconfig of the plugin to tell the plugin the
	// configdigest of its predecessor protocol instance.
	//
	// NOTE: Set here:
	// https://github.com/smartcontractkit/mercury-v1-sketch/blob/f52c0f823788f86c1aeaa9ba1eee32a85b981535/onchain/src/ConfigurationStore.sol#L13
	// TODO: This needs to be implemented alongside staging/production
	// switchover support: https://smartcontract-it.atlassian.net/browse/MERC-3386
	PredecessorConfigDigest *types.ConfigDigest
	// TODO: Billing
	// https://smartcontract-it.atlassian.net/browse/MERC-1189
	// QUESTION: Previously we stored ExpiryWindow and BaseUSDFeeCents in offchain
	// config, but those might be channel specific so need to move to
	// channel definition
	// ExpirationWindow uint32          `json:"expirationWindow"` // Integer number of seconds
	// BaseUSDFee       decimal.Decimal `json:"baseUSDFee"`       // Base USD fee
}

func DecodeOffchainConfig(b []byte) (o OffchainConfig, err error) {
	pbuf := &LLOOffchainConfigProto{}
	err = proto.Unmarshal(b, pbuf)
	if err != nil {
		return o, fmt.Errorf("failed to decode offchain config: expected protobuf (got: 0x%x); %w", b, err)
	}
	if len(pbuf.PredecessorConfigDigest) > 0 {
		var predecessorConfigDigest types.ConfigDigest
		predecessorConfigDigest, err = types.BytesToConfigDigest(pbuf.PredecessorConfigDigest)
		if err != nil {
			return o, err
		}
		o.PredecessorConfigDigest = &predecessorConfigDigest
	}
	return
}

func (c OffchainConfig) Encode() ([]byte, error) {
	pbuf := LLOOffchainConfigProto{}
	if c.PredecessorConfigDigest != nil {
		pbuf.PredecessorConfigDigest = c.PredecessorConfigDigest[:]
	}
	return proto.Marshal(&pbuf)
}
