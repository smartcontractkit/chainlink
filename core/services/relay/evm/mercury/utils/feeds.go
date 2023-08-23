package utils

import (
	"encoding/binary"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type FeedVersion uint16

const (
	_ FeedVersion = iota
	REPORT_V1
	REPORT_V2
	REPORT_V3
	_
)

type FeedID [32]byte

func (f FeedID) Hex() string { return (utils.Hash)(f).Hex() }

func (f FeedID) String() string { return (utils.Hash)(f).String() }

func (f *FeedID) UnmarshalText(input []byte) error {
	return (*utils.Hash)(f).UnmarshalText(input)
}

func (f FeedID) Version() FeedVersion {
	return FeedVersion(binary.BigEndian.Uint16(f[:2]))
}

func (f FeedID) IsV1() bool { return f.Version() == REPORT_V1 }
func (f FeedID) IsV2() bool { return f.Version() == REPORT_V2 }
func (f FeedID) IsV3() bool { return f.Version() == REPORT_V3 }
