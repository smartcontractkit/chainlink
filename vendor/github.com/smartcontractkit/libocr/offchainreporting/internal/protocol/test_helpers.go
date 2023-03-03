package protocol

import "github.com/smartcontractkit/libocr/commontypes"

// Used only for testing
type XXXUnknownMessageType struct{}

// Conform to protocol.Message interface
func (XXXUnknownMessageType) process(*oracleState, commontypes.OracleID) {}
