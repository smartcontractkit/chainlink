package deployment

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/require"
)

func TestLabeledAddresses_And(t *testing.T) {
	tests := []struct {
		name   string
		give   LabeledAddresses
		labels []string
		want   LabeledAddresses
	}{
		{
			name: "No labels returns unlabeled",
			give: LabeledAddresses{
				"addr1": {Type: "type1", Version: *semver.MustParse("1.0.0"), Labels: nil},
				"addr2": {Type: "type2", Version: *semver.MustParse("2.0.0"), Labels: NewLabelSet("label1")},
			},
			labels: nil,
			want: LabeledAddresses{
				"addr1": {Type: "type1", Version: *semver.MustParse("1.0.0"), Labels: nil},
			},
		},
		{
			name: "One label",
			give: LabeledAddresses{
				"addr1": {Type: "type1", Version: *semver.MustParse("1.0.0"), Labels: NewLabelSet("label1")},
				"addr2": {Type: "type2", Version: *semver.MustParse("2.0.0"), Labels: NewLabelSet("label2")},
				"addr3": {Type: "type3", Version: *semver.MustParse("3.0.0"), Labels: NewLabelSet("label1", "label2")},
			},
			labels: []string{"label1"},
			want: LabeledAddresses{
				"addr1": {Type: "type1", Version: *semver.MustParse("1.0.0"), Labels: NewLabelSet("label1")},
				"addr3": {Type: "type3", Version: *semver.MustParse("3.0.0"), Labels: NewLabelSet("label1", "label2")},
			},
		},
		{
			name: "Multiple labels",
			give: LabeledAddresses{
				"addr1": {Type: "type1", Version: *semver.MustParse("1.0.0"), Labels: NewLabelSet("label1", "label2")},
				"addr2": {Type: "type2", Version: *semver.MustParse("2.0.0"), Labels: NewLabelSet("label1")},
				"addr3": {Type: "type3", Version: *semver.MustParse("3.0.0"), Labels: NewLabelSet("label2")},
			},
			labels: []string{"label1", "label2"},
			want: LabeledAddresses{
				"addr1": {Type: "type1", Version: *semver.MustParse("1.0.0"), Labels: NewLabelSet("label1", "label2")},
			},
		},
		{
			name:   "Empty LabeledAddresses",
			give:   LabeledAddresses{},
			labels: []string{"label1"},
			want:   LabeledAddresses{},
		},
		{
			name: "No matching labels",
			give: LabeledAddresses{
				"addr1": {Type: "type1", Version: *semver.MustParse("1.0.0"), Labels: NewLabelSet("label1")},
			},
			labels: []string{"label2"},
			want:   LabeledAddresses{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.give.And(tt.labels...)
			require.Equal(t, len(tt.want), len(got), "len of got != len of want")

			for addr, lsa := range tt.want {
				lsb, exists := got[addr]
				require.Truef(t, exists, "key %s is missing from got", addr)
				require.True(t, lsa.Equal(lsb), "label sets differ from got and want")
			}
		})
	}
}
