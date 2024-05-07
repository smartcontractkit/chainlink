package liquiditymanager

import (
	"fmt"
	"math/big"
	"math/rand"
	"testing"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

func TestPlugin_ValidateObservation(t *testing.T) {
	testCases := []struct {
		name   string
		obs    ocrtypes.Observation
		expErr func(t *testing.T, err error)
	}{
		{
			name: "some random bytes",
			obs:  ocrtypes.Observation("abc"),
			expErr: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "empty is ok",
			obs:  ocrtypes.Observation("{}"),
		},
		{
			name: "some observation",
			obs: models.NewObservation(
				[]models.NetworkLiquidity{},
				[]models.Transfer{{}},
				[]models.PendingTransfer{},
				[]models.Transfer{},
				[]models.Edge{},
				[]models.ConfigDigestWithMeta{},
			).Encode(),
		},
		{
			name: "deduped liquidity observations",
			obs: models.NewObservation(
				[]models.NetworkLiquidity{{Network: 1, Liquidity: ubig.New(big.NewInt(1))}, {Network: 1, Liquidity: ubig.New(big.NewInt(2))}},
				[]models.Transfer{},
				[]models.PendingTransfer{},
				[]models.Transfer{},
				[]models.Edge{},
				[]models.ConfigDigestWithMeta{},
			).Encode(),
			expErr: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "deduped resolved transfers",
			obs: models.NewObservation(
				[]models.NetworkLiquidity{},
				[]models.Transfer{{From: 1}, {From: 1}},
				[]models.PendingTransfer{},
				[]models.Transfer{},
				[]models.Edge{},
				[]models.ConfigDigestWithMeta{},
			).Encode(),
			expErr: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "deduped pending transfers",
			obs: models.NewObservation(
				[]models.NetworkLiquidity{},
				[]models.Transfer{},
				[]models.PendingTransfer{{ID: "1"}, {ID: "1"}},
				[]models.Transfer{},
				[]models.Edge{},
				[]models.ConfigDigestWithMeta{},
			).Encode(),
			expErr: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "deduped inflight transfers",
			obs: models.NewObservation(
				[]models.NetworkLiquidity{},
				[]models.Transfer{},
				[]models.PendingTransfer{},
				[]models.Transfer{{From: 1}, {From: 1}},
				[]models.Edge{},
				[]models.ConfigDigestWithMeta{},
			).Encode(),
			expErr: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "deduped edges",
			obs: models.NewObservation(
				[]models.NetworkLiquidity{},
				[]models.Transfer{},
				[]models.PendingTransfer{},
				[]models.Transfer{},
				[]models.Edge{{Source: 1, Dest: 2}, {Source: 1, Dest: 2}},
				[]models.ConfigDigestWithMeta{},
			).Encode(),
			expErr: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name: "deduped config digest",
			obs: models.NewObservation(
				[]models.NetworkLiquidity{},
				[]models.Transfer{},
				[]models.PendingTransfer{},
				[]models.Transfer{},
				[]models.Edge{},
				[]models.ConfigDigestWithMeta{{NetworkSel: 1}, {NetworkSel: 1}},
			).Encode(),
			expErr: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := newPluginWithMocksAndDefaults(t)
			ao := ocrtypes.AttributedObservation{
				Observation: tc.obs,
				Observer:    commontypes.OracleID(uint8(rand.Intn(10))), // ignored by the plugin
			}
			err := p.plugin.ValidateObservation(ocr3types.OutcomeContext{}, ocrtypes.Query{}, ao)
			if tc.expErr != nil {
				tc.expErr(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func Test_validateDedupedItems(t *testing.T) {
	tests := []struct {
		name    string
		keyFn   func(models.Transfer) string
		items   []models.Transfer
		wantErr bool
	}{
		{
			name: "no duplicates",
			items: []models.Transfer{
				{From: 1},
				{From: 2},
				{From: 3},
			},
			keyFn:   dedupKeyTransfer,
			wantErr: false,
		},
		{
			name: "duplicates",
			items: []models.Transfer{
				{From: 1},
				{From: 2},
				{From: 1},
			},
			keyFn:   dedupKeyTransfer,
			wantErr: true,
		},
		{
			name:    "empty",
			items:   []models.Transfer{},
			keyFn:   dedupKeyTransfer,
			wantErr: false,
		},
		{
			name: "custom keyFn",
			keyFn: func(t models.Transfer) string {
				return fmt.Sprintf("%d", t.From)
			},
			items: []models.Transfer{
				{From: 1, To: 2},
				{From: 1, To: 3},
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validateDedupedItems(tc.keyFn, tc.items...)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
