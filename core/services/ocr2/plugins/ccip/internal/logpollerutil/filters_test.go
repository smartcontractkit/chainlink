package logpollerutil

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
)

func Test_FiltersDiff(t *testing.T) {
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
			gotCreated, gotDeleted := FiltersDiff(tt.args.filtersBefore, tt.args.filtersNow)
			assert.Equalf(t, tt.wantCreated, gotCreated, "filtersDiff(%v, %v)", tt.args.filtersBefore, tt.args.filtersNow)
			assert.Equalf(t, tt.wantDeleted, gotDeleted, "filtersDiff(%v, %v)", tt.args.filtersBefore, tt.args.filtersNow)
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
