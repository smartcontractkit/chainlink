package protocol

import "github.com/smartcontractkit/chainlink/libocr/offchainreporting/types"

type XXXUnknownMessageType struct{}

func (XXXUnknownMessageType) process(*oracleState, types.OracleID) {}
