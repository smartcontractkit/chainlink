package internal

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink-relay/pkg/types"
)

func TestNewPageToken(t *testing.T) {
	type args struct {
		t *PageToken
	}
	tests := []struct {
		name    string
		args    args
		want    *PageToken
		wantErr bool
	}{
		{
			name: "empty",
			args: args{t: &PageToken{}},
			want: &PageToken{Page: 0, Size: defaultSize},
		},
		{
			name: "page set, size unset",
			args: args{t: &PageToken{Page: 1}},
			want: &PageToken{Page: 1, Size: defaultSize},
		},
		{
			name: "page set, size set",
			args: args{t: &PageToken{Page: 3, Size: 10}},
			want: &PageToken{Page: 3, Size: 10},
		},
		{
			name: "page unset, size set",
			args: args{t: &PageToken{Size: 17}},
			want: &PageToken{Page: 0, Size: 17},
		},
	}
	for _, tt := range tests {
		enc := tt.args.t.Encode()
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPageToken(enc)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPageToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPageToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListNodeStatuses(t *testing.T) {
	testStats := []types.NodeStatus{
		types.NodeStatus{
			ChainID: "chain-1",
			Name:    "name-1",
		},
		types.NodeStatus{
			ChainID: "chain-2",
			Name:    "name-2",
		},
		types.NodeStatus{
			ChainID: "chain-3",
			Name:    "name-3",
		},
	}

	type args struct {
		pageSize  int
		pageToken string
		listFn    ListNodeStatusFn
	}
	tests := []struct {
		name               string
		args               args
		wantStats          []types.NodeStatus
		wantNext_pageToken string
		wantTotal          int
		wantErr            bool
	}{
		{
			name: "all on first page",
			args: args{
				pageSize:  10, // > length of test stats
				pageToken: "",
				listFn: func(start, end int) ([]types.NodeStatus, int, error) {
					return testStats, len(testStats), nil
				},
			},
			wantNext_pageToken: "",
			wantTotal:          len(testStats),
			wantStats:          testStats,
		},
		{
			name: "small first page",
			args: args{
				pageSize:  len(testStats) - 1,
				pageToken: "",
				listFn: func(start, end int) ([]types.NodeStatus, int, error) {
					return testStats[start:end], len(testStats), nil
				},
			},
			wantNext_pageToken: base64.RawStdEncoding.EncodeToString([]byte("page=1&size=2")), // hard coded 2 is len(testStats)-1
			wantTotal:          len(testStats),
			wantStats:          testStats[0 : len(testStats)-1],
		},
		{
			name: "second page",
			args: args{
				pageSize:  len(testStats) - 1,
				pageToken: base64.RawStdEncoding.EncodeToString([]byte("page=1&size=2")), // hard coded 2 is len(testStats)-1
				listFn: func(start, end int) ([]types.NodeStatus, int, error) {
					// note list function must do the start, end bound checking. here we are making it simple
					if end > len(testStats) {
						end = len(testStats)
					}
					return testStats[start:end], len(testStats), nil
				},
			},
			wantNext_pageToken: "",
			wantTotal:          len(testStats),
			wantStats:          testStats[len(testStats)-1:],
		},
		{
			name: "bad list fn",
			args: args{
				listFn: func(start, end int) ([]types.NodeStatus, int, error) {
					return nil, 0, fmt.Errorf("i'm a bad list fn")
				},
			},
			wantTotal: -1,
			wantErr:   true,
		},
		{
			name: "invalid token",
			args: args{
				pageToken: "invalid token",
				listFn: func(start, end int) ([]types.NodeStatus, int, error) {
					return testStats[start:end], len(testStats), nil
				},
			},
			wantTotal: -1,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStats, gotNext_pageToken, gotTotal, err := ListNodeStatuses(tt.args.pageSize, tt.args.pageToken, tt.args.listFn)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListNodeStatuses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.wantStats, gotStats)
			assert.Equal(t, tt.wantNext_pageToken, gotNext_pageToken)
			assert.Equal(t, tt.wantTotal, gotTotal)
		})
	}
}
