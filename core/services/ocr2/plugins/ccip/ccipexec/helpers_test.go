package ccipexec

import (
	"testing"

	"github.com/stretchr/testify/assert"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
)

func Test_CommitReportWithSendRequests_uniqueSenders(t *testing.T) {
	messageFn := func(address cciptypes.Address) cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta {
		return cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{EVM2EVMMessage: cciptypes.EVM2EVMMessage{Sender: address}}
	}

	tests := []struct {
		name             string
		sendRequests     []cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta
		expUniqueSenders int
		expSendersOrder  []cciptypes.Address
	}{
		{
			name: "all unique senders",
			sendRequests: []cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{
				messageFn(cciptypes.Address(utils.RandomAddress().String())),
				messageFn(cciptypes.Address(utils.RandomAddress().String())),
				messageFn(cciptypes.Address(utils.RandomAddress().String())),
			},
			expUniqueSenders: 3,
		},
		{
			name: "some senders are the same",
			sendRequests: []cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{
				messageFn("0x1"),
				messageFn("0x2"),
				messageFn("0x1"),
				messageFn("0x2"),
				messageFn("0x3"),
			},
			expUniqueSenders: 3,
			expSendersOrder: []cciptypes.Address{
				cciptypes.Address("0x1"),
				cciptypes.Address("0x2"),
				cciptypes.Address("0x3"),
			},
		},
		{
			name: "all senders are the same",
			sendRequests: []cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{
				messageFn("0x1"),
				messageFn("0x1"),
				messageFn("0x1"),
			},
			expUniqueSenders: 1,
			expSendersOrder: []cciptypes.Address{
				cciptypes.Address("0x1"),
			},
		},
		{
			name: "order is preserved",
			sendRequests: []cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{
				messageFn("0x3"),
				messageFn("0x1"),
				messageFn("0x3"),
				messageFn("0x2"),
				messageFn("0x2"),
				messageFn("0x1"),
			},
			expUniqueSenders: 3,
			expSendersOrder: []cciptypes.Address{
				cciptypes.Address("0x3"),
				cciptypes.Address("0x1"),
				cciptypes.Address("0x2"),
			},
		},
		{
			name:             "no senders",
			sendRequests:     []cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{},
			expUniqueSenders: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rep := commitReportWithSendRequests{sendRequestsWithMeta: tt.sendRequests}
			uniqueSenders := rep.uniqueSenders()

			assert.Len(t, uniqueSenders, tt.expUniqueSenders)
			if tt.expSendersOrder != nil {
				assert.Equal(t, tt.expSendersOrder, uniqueSenders)
			}
		})
	}
}
