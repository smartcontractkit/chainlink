package ccip

import (
	"math"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func Test_filtersDiff(t *testing.T) {
	type args struct {
		filtersBefore []logpoller.Filter
		filtersNow    []logpoller.Filter
	}
	tests := []struct {
		name        string
		args        args
		wantCreated []logpoller.Filter
		wantDeleted []logpoller.Filter
	}{
		{
			name: "no diff, both empty",
			args: args{
				filtersBefore: []logpoller.Filter{},
				filtersNow:    []logpoller.Filter{},
			},
			wantCreated: []logpoller.Filter{},
			wantDeleted: []logpoller.Filter{},
		},
		{
			name: "no diff, both non-empty",
			args: args{
				filtersBefore: []logpoller.Filter{{Name: "a"}},
				filtersNow:    []logpoller.Filter{{Name: "a"}},
			},
			wantCreated: []logpoller.Filter{},
			wantDeleted: []logpoller.Filter{},
		},
		{
			name: "no diff, only name matters",
			args: args{
				filtersBefore: []logpoller.Filter{{Name: "a", Retention: time.Minute}},
				filtersNow:    []logpoller.Filter{{Name: "a", Retention: time.Second}},
			},
			wantCreated: []logpoller.Filter{},
			wantDeleted: []logpoller.Filter{},
		},
		{
			name: "diff for both created and deleted",
			args: args{
				filtersBefore: []logpoller.Filter{{Name: "e"}, {Name: "a"}, {Name: "b"}},
				filtersNow:    []logpoller.Filter{{Name: "a"}, {Name: "c"}, {Name: "d"}},
			},
			wantCreated: []logpoller.Filter{{Name: "c"}, {Name: "d"}},
			wantDeleted: []logpoller.Filter{{Name: "e"}, {Name: "b"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCreated, gotDeleted := filtersDiff(tt.args.filtersBefore, tt.args.filtersNow)
			assert.Equalf(t, tt.wantCreated, gotCreated, "filtersDiff(%v, %v)", tt.args.filtersBefore, tt.args.filtersNow)
			assert.Equalf(t, tt.wantDeleted, gotDeleted, "filtersDiff(%v, %v)", tt.args.filtersBefore, tt.args.filtersNow)
		})
	}
}

func Test_containsFilter(t *testing.T) {
	type args struct {
		filters []logpoller.Filter
		f       logpoller.Filter
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty",
			args: args{
				filters: []logpoller.Filter{},
				f:       logpoller.Filter{},
			},
			want: false,
		},
		{
			name: "contains",
			args: args{
				filters: []logpoller.Filter{{Name: "a"}, {Name: "b"}},
				f:       logpoller.Filter{Name: "b"},
			},
			want: true,
		},
		{
			name: "does not contain",
			args: args{
				filters: []logpoller.Filter{{Name: "a"}, {Name: "b"}},
				f:       logpoller.Filter{Name: "c"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want,
				containsFilter(tt.args.filters, tt.args.f), "containsFilter(%v, %v)", tt.args.filters, tt.args.f)
		})
	}
}

func Test_filterContainsZeroAddress(t *testing.T) {
	type args struct {
		addrs []common.Address
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "non-zero addrs",
			args: args{
				addrs: []common.Address{
					common.HexToAddress("1"),
					common.HexToAddress("2"),
					common.HexToAddress("3"),
				},
			},
			want: false,
		},
		{
			name: "empty",
			args: args{addrs: []common.Address{}},
			want: false,
		},
		{
			name: "zero addr",
			args: args{
				addrs: []common.Address{
					common.HexToAddress("1"),
					common.HexToAddress("0"),
					common.HexToAddress("2"),
					common.HexToAddress("3"),
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, filterContainsZeroAddress(tt.args.addrs), "filterContainsZeroAddress(%v)", tt.args.addrs)
		})
	}
}

func Test_mergeEpochAndRound(t *testing.T) {
	type args struct {
		epoch uint32
		round uint8
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{
			name: "zero round and epoch",
			args: args{epoch: 0, round: 0},
			want: 0,
		},
		{
			name: "avg case",
			args: args{
				epoch: 243,
				round: 15,
			},
			want: 62223,
		},
		{
			name: "largest epoch and round",
			args: args{
				epoch: math.MaxUint32,
				round: math.MaxUint8,
			},
			want: 1099511627775,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want,
				mergeEpochAndRound(tt.args.epoch, tt.args.round),
				"mergeEpochAndRound(%v, %v)", tt.args.epoch, tt.args.round)
		})
	}
}

func Test_bytesOfBytesKeccak(t *testing.T) {
	h, err := bytesOfBytesKeccak(nil)
	assert.NoError(t, err)
	assert.Equal(t, [32]byte{}, h)

	h1, err := bytesOfBytesKeccak([][]byte{{0x1}, {0x1}})
	assert.NoError(t, err)
	h2, err := bytesOfBytesKeccak([][]byte{{0x1, 0x1}})
	assert.NoError(t, err)
	assert.NotEqual(t, h1, h2)
}

func Test_contiguousReqs(t *testing.T) {
	testCases := []struct {
		min    uint64
		max    uint64
		seqNrs []uint64
		exp    bool
	}{
		{min: 5, max: 10, seqNrs: []uint64{5, 6, 7, 8, 9, 10}, exp: true},
		{min: 5, max: 10, seqNrs: []uint64{5, 7, 8, 9, 10}, exp: false},
		{min: 5, max: 10, seqNrs: []uint64{5, 6, 7, 8, 9, 10, 11}, exp: false},
		{min: 5, max: 10, seqNrs: []uint64{}, exp: false},
		{min: 1, max: 1, seqNrs: []uint64{1}, exp: true},
	}

	for _, tc := range testCases {
		res := contiguousReqs(logger.NullLogger, tc.min, tc.max, tc.seqNrs)
		assert.Equal(t, tc.exp, res)
	}
}

func Test_getMessageIDsAsHexString(t *testing.T) {
	t.Run("base", func(t *testing.T) {
		hashes := make([]common.Hash, 10)
		for i := range hashes {
			hashes[i] = common.HexToHash(strconv.Itoa(rand.Intn(100000)))
		}

		msgs := make([]evm_2_evm_offramp.InternalEVM2EVMMessage, len(hashes))
		for i := range msgs {
			msgs[i] = evm_2_evm_offramp.InternalEVM2EVMMessage{MessageId: hashes[i]}
		}

		messageIDs := getMessageIDsAsHexString(msgs)
		for i := range messageIDs {
			assert.Equal(t, hashes[i].String(), messageIDs[i])
		}
	})

	t.Run("empty", func(t *testing.T) {
		messageIDs := getMessageIDsAsHexString(nil)
		assert.Empty(t, messageIDs)
	})
}
