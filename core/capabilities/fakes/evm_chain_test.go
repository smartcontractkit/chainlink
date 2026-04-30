package fakes

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	evmcappb "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/chain-capabilities/evm"
)

// transferEventSig is keccak256("Transfer(address,address,uint256)"), used as a
// stand-in event signature (topic0) throughout the filter tests.
var transferEventSig = mustDecodeHex("ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")

func mustDecodeHex(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}

var (
	testAddr  = []byte{0x11, 0x22, 0x33}
	otherAddr = []byte{0xAA, 0xBB, 0xCC}
	topic1    = []byte{0x01, 0x02, 0x03}
	topic2    = []byte{0x04, 0x05, 0x06}
	topic3    = []byte{0x07, 0x08, 0x09}
)

func makeLog(addr []byte, topics ...[]byte) *evmcappb.Log {
	return &evmcappb.Log{
		Address: addr,
		Topics:  topics,
	}
}

func makeFilter(addr []byte, topics ...*evmcappb.TopicValues) *evmcappb.FilterLogTriggerRequest {
	return &evmcappb.FilterLogTriggerRequest{
		Addresses: func() [][]byte {
			if addr == nil {
				return nil
			}
			return [][]byte{addr}
		}(),
		Topics: topics,
	}
}

func tv(values ...[]byte) *evmcappb.TopicValues {
	return &evmcappb.TopicValues{Values: values}
}

func wildcardTV() *evmcappb.TopicValues {
	return &evmcappb.TopicValues{} // empty Values = wildcard
}

func TestFakeEVMLogMatchesFilter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		log         *evmcappb.Log
		filter      *evmcappb.FilterLogTriggerRequest
		wantErr     bool
		errContains string
	}{
		{
			name:   "exact match on address and topic0",
			log:    makeLog(testAddr, transferEventSig),
			filter: makeFilter(testAddr, tv(transferEventSig)),
		},
		{
			// This is the most common real-world user error: a developer (especially
			// using the TypeScript SDK or a raw call) registers a log trigger without
			// specifying topic0.  Without topic0 the trigger fires for every event
			// emitted by the contract, not just the intended one.  The fake surfaces
			// this mistake immediately instead of letting it slip through to staging.
			name:        "missing topic0 — empty Topics slice",
			log:         makeLog(testAddr, transferEventSig),
			filter:      makeFilter(testAddr /* no Topics */),
			wantErr:     true,
			errContains: "missing topic0",
		},
		{
			// Same mistake, slightly different form: Topics slice present but slot 0
			// has no Values, i.e. it was left as a wildcard.
			name:        "missing topic0 — slot 0 has empty Values",
			log:         makeLog(testAddr, transferEventSig),
			filter:      makeFilter(testAddr, wildcardTV()),
			wantErr:     true,
			errContains: "missing topic0",
		},
		{
			name:        "address mismatch",
			log:         makeLog(otherAddr, transferEventSig),
			filter:      makeFilter(testAddr, tv(transferEventSig)),
			wantErr:     true,
			errContains: "does not match any of the addresses",
		},
		{
			name:        "topic0 mismatch",
			log:         makeLog(testAddr, topic1),
			filter:      makeFilter(testAddr, tv(transferEventSig)),
			wantErr:     true,
			errContains: "log topic 0 does not match",
		},
		{
			// No address constraint in the filter — any address should match.
			name:   "wildcard address (no Addresses in filter)",
			log:    makeLog(otherAddr, transferEventSig),
			filter: makeFilter(nil, tv(transferEventSig)),
		},
		{
			// Slots 1-3 are allowed to be wildcards even when topic0 is set.
			name:   "wildcard slot 1 matches any indexed arg",
			log:    makeLog(testAddr, transferEventSig, topic1),
			filter: makeFilter(testAddr, tv(transferEventSig), wildcardTV()),
		},
		{
			name:   "OR semantics within a topic slot",
			log:    makeLog(testAddr, transferEventSig, topic2),
			filter: makeFilter(testAddr, tv(transferEventSig), tv(topic1, topic2)),
		},
		{
			name:        "AND semantics across topic slots — second slot mismatch",
			log:         makeLog(testAddr, transferEventSig, topic1),
			filter:      makeFilter(testAddr, tv(transferEventSig), tv(topic2)),
			wantErr:     true,
			errContains: "log topic 1 does not match",
		},
		{
			name:        "log has fewer topics than filter requires",
			log:         makeLog(testAddr, transferEventSig), // no topic1
			filter:      makeFilter(testAddr, tv(transferEventSig), tv(topic1)),
			wantErr:     true,
			errContains: "log topics length",
		},
		{
			// A real customer scenario: the filter constrains topic0, topic1, and
			// topic3, but intentionally leaves topic2 as a wildcard (null) so that
			// any value in that indexed arg position is accepted.  The filter still
			// must match a log that carries all four topics.
			name: "wildcard topic2 with constrained topic0, topic1, and topic3",
			log:  makeLog(testAddr, transferEventSig, topic1, topic2, topic3),
			filter: makeFilter(testAddr,
				tv(transferEventSig), // slot 0: event signature
				tv(topic1),           // slot 1: constrained
				wildcardTV(),         // slot 2: wildcard — any value accepted
				tv(topic3),           // slot 3: constrained
			),
		},
		{
			// Same layout but topic3 on the log does not match the filter value.
			name: "wildcard topic2 — topic3 mismatch",
			log:  makeLog(testAddr, transferEventSig, topic1, topic2, topic2 /* wrong topic3 */),
			filter: makeFilter(testAddr,
				tv(transferEventSig),
				tv(topic1),
				wildcardTV(),
				tv(topic3),
			),
			wantErr:     true,
			errContains: "log topic 3 does not match",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := fakeEVMLogMatchesFilter(tc.log, tc.filter)
			if tc.wantErr {
				require.Error(t, err)
				if tc.errContains != "" {
					assert.Contains(t, err.Error(), tc.errContains)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}
