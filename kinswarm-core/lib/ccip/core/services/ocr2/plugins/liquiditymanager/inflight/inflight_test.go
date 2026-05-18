package inflight

import (
	"testing"

	"github.com/stretchr/testify/require"

	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

func Test_inflight_Add(t *testing.T) {
	type fields struct {
		items map[transferID]models.Transfer
	}
	type args struct {
		t models.Transfer
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		assertion func(*testing.T, *inflight)
	}{
		{
			"transfer not in map",
			fields{
				items: map[transferID]models.Transfer{},
			},
			args{
				models.Transfer{
					From:   1,
					To:     2,
					Amount: ubig.NewI(1),
				},
			},
			func(t *testing.T, i *inflight) {
				item, ok := i.transfers[transferID{From: 1, To: 2, Amount: "1"}]
				require.True(t, ok)
				require.Equal(t, models.Transfer{
					From:   1,
					To:     2,
					Amount: ubig.NewI(1),
				}, item)
			},
		},
		{
			"transfer in map",
			fields{
				items: map[transferID]models.Transfer{
					{From: 1, To: 2, Amount: "1"}: {
						From:   1,
						To:     2,
						Amount: ubig.NewI(1),
					},
				},
			},
			args{
				models.Transfer{
					From:   1,
					To:     2,
					Amount: ubig.NewI(1),
				},
			},
			func(t *testing.T, i *inflight) {
				require.Len(t, i.transfers, 1)
				item, ok := i.transfers[transferID{From: 1, To: 2, Amount: "1"}]
				require.True(t, ok)
				require.Equal(t, models.Transfer{
					From:   1,
					To:     2,
					Amount: ubig.NewI(1),
				}, item)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &inflight{
				transfers: tt.fields.items,
			}
			i.Add(tt.args.t)
			tt.assertion(t, i)
		})
	}
}

func Test_inflight_Expire(t *testing.T) {
	type fields struct {
		items map[transferID]models.Transfer
	}
	type args struct {
		pending []models.PendingTransfer
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		before    func(*testing.T, *inflight)
		assertion func(*testing.T, *inflight)
	}{
		{
			"no pending transfers",
			fields{
				items: map[transferID]models.Transfer{
					{From: 1, To: 2, Amount: "1"}: {
						From:   1,
						To:     2,
						Amount: ubig.NewI(1),
						Stage:  0,
					},
				},
			},
			args{
				[]models.PendingTransfer{},
			},
			func(t *testing.T, i *inflight) {
				require.Len(t, i.transfers, 1)
			},
			func(t *testing.T, i *inflight) {
				require.Len(t, i.transfers, 1)
			},
		},
		{
			"pending transfer with larger stage",
			fields{
				items: map[transferID]models.Transfer{
					{From: 1, To: 2, Amount: "1"}: {
						From:   1,
						To:     2,
						Amount: ubig.NewI(1),
						Stage:  0,
					},
				},
			},
			args{
				[]models.PendingTransfer{
					{
						Transfer: models.Transfer{
							From:   1,
							To:     2,
							Amount: ubig.NewI(1),
							Stage:  1,
						},
					},
				},
			},
			func(t *testing.T, i *inflight) {
				require.Len(t, i.transfers, 1)
			},
			func(t *testing.T, i *inflight) {
				require.Len(t, i.transfers, 0)
			},
		},
		{
			"pending transfer with equal stage",
			fields{
				items: map[transferID]models.Transfer{
					{From: 1, To: 2, Amount: "1"}: {
						From:   1,
						To:     2,
						Amount: ubig.NewI(1),
						Stage:  0,
					},
				},
			},
			args{
				[]models.PendingTransfer{
					{
						Transfer: models.Transfer{
							From:   1,
							To:     2,
							Amount: ubig.NewI(1),
							Stage:  0,
						},
					},
				},
			},
			func(t *testing.T, i *inflight) {
				require.Len(t, i.transfers, 1)
			},
			func(t *testing.T, i *inflight) {
				require.Len(t, i.transfers, 1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &inflight{
				transfers: tt.fields.items,
			}
			tt.before(t, i)
			i.Expire(tt.args.pending)
			tt.assertion(t, i)
		})
	}
}

func Test_inflight_GetAll(t *testing.T) {
	type fields struct {
		items map[transferID]models.Transfer
	}
	tests := []struct {
		name   string
		fields fields
		want   []models.Transfer
	}{
		{
			"empty",
			fields{
				items: map[transferID]models.Transfer{},
			},
			[]models.Transfer{},
		},
		{
			"not empty",
			fields{
				items: map[transferID]models.Transfer{
					{From: 1, To: 2, Amount: "1"}: {
						From:   1,
						To:     2,
						Amount: ubig.NewI(1),
					},
				},
			},
			[]models.Transfer{
				{
					From:   1,
					To:     2,
					Amount: ubig.NewI(1),
				},
			},
		},
		{
			"multiple",
			fields{
				items: map[transferID]models.Transfer{
					{From: 1, To: 2, Amount: "1"}: {
						From:   1,
						To:     2,
						Amount: ubig.NewI(1),
					},
					{From: 2, To: 3, Amount: "2"}: {
						From:   2,
						To:     3,
						Amount: ubig.NewI(2),
					},
				},
			},
			[]models.Transfer{
				{
					From:   1,
					To:     2,
					Amount: ubig.NewI(1),
				},
				{
					From:   2,
					To:     3,
					Amount: ubig.NewI(2),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &inflight{
				transfers: tt.fields.items,
			}
			got := i.GetAll()
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_inflight_IsInflight(t *testing.T) {
	type fields struct {
		items map[transferID]models.Transfer
	}
	type args struct {
		t models.Transfer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			"not inflight",
			fields{
				items: map[transferID]models.Transfer{},
			},
			args{
				models.Transfer{
					From:   1,
					To:     2,
					Amount: ubig.NewI(1),
				},
			},
			false,
		},
		{
			"inflight",
			fields{
				items: map[transferID]models.Transfer{
					{From: 1, To: 2, Amount: "1"}: {
						From:   1,
						To:     2,
						Amount: ubig.NewI(1),
					},
				},
			},
			args{
				models.Transfer{
					From:   1,
					To:     2,
					Amount: ubig.NewI(1),
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &inflight{
				transfers: tt.fields.items,
			}
			got := i.IsInflight(tt.args.t)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *inflight
	}{
		{
			"basic",
			&inflight{
				transfers: map[transferID]models.Transfer{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New()
			require.Equal(t, tt.want.transfers, got.transfers)
		})
	}
}
