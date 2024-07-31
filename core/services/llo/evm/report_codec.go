package evm

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"

	"github.com/smartcontractkit/chainlink-data-streams/llo"
)

var (
	_      llo.ReportCodec = ReportCodec{}
	Schema                 = getSchema()
)

func getSchema() abi.Arguments {
	mustNewType := func(t string) abi.Type {
		result, err := abi.NewType(t, "", []abi.ArgumentMarshaling{})
		if err != nil {
			panic(fmt.Sprintf("Unexpected error during abi.NewType: %s", err))
		}
		return result
	}
	return abi.Arguments([]abi.Argument{
		{Name: "configDigest", Type: mustNewType("bytes32")},
		{Name: "chainId", Type: mustNewType("uint64")},
		// TODO:
		// could also include address of verifier to make things more specific.
		// downside is increased data size.
		// for now we assume that a channelId will only be registered on a single
		// verifier per chain.
		// https://smartcontract-it.atlassian.net/browse/MERC-3652
		{Name: "seqNr", Type: mustNewType("uint64")},
		{Name: "channelId", Type: mustNewType("uint32")},
		{Name: "validAfterSeconds", Type: mustNewType("uint32")},
		{Name: "validUntilSeconds", Type: mustNewType("uint32")},
		{Name: "values", Type: mustNewType("int192[]")},
		{Name: "specimen", Type: mustNewType("bool")},
	})
}

type ReportCodec struct{}

func NewReportCodec() ReportCodec {
	return ReportCodec{}
}

func (ReportCodec) Encode(report llo.Report, cd llotypes.ChannelDefinition) ([]byte, error) {
	return nil, errors.New("not implemented")
}
