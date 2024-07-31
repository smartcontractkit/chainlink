package graph

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/models"
)

func TestXChainRebalancerData_Equals(t *testing.T) {
	type fields struct {
		RemoteRebalancerAddress   models.Address
		LocalBridgeAdapterAddress models.Address
		RemoteTokenAddress        models.Address
	}
	type args struct {
		other XChainLiquidityManagerData
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			"equal",
			fields{
				RemoteRebalancerAddress:   models.Address(common.HexToAddress("0x1")),
				LocalBridgeAdapterAddress: models.Address(common.HexToAddress("0x2")),
				RemoteTokenAddress:        models.Address(common.HexToAddress("0x3")),
			},
			args{
				other: XChainLiquidityManagerData{
					RemoteLiquidityManagerAddress: models.Address(common.HexToAddress("0x1")),
					LocalBridgeAdapterAddress:     models.Address(common.HexToAddress("0x2")),
					RemoteTokenAddress:            models.Address(common.HexToAddress("0x3")),
				},
			},
			true,
		},
		{
			"not equal remote rebalancer",
			fields{
				RemoteRebalancerAddress:   models.Address(common.HexToAddress("0x1")),
				LocalBridgeAdapterAddress: models.Address(common.HexToAddress("0x2")),
				RemoteTokenAddress:        models.Address(common.HexToAddress("0x3")),
			},
			args{
				other: XChainLiquidityManagerData{
					RemoteLiquidityManagerAddress: models.Address(common.HexToAddress("0x4")),
					LocalBridgeAdapterAddress:     models.Address(common.HexToAddress("0x2")),
					RemoteTokenAddress:            models.Address(common.HexToAddress("0x3")),
				},
			},
			false,
		},
		{
			"not equal local bridge",
			fields{
				RemoteRebalancerAddress:   models.Address(common.HexToAddress("0x1")),
				LocalBridgeAdapterAddress: models.Address(common.HexToAddress("0x2")),
				RemoteTokenAddress:        models.Address(common.HexToAddress("0x3")),
			},
			args{
				other: XChainLiquidityManagerData{
					RemoteLiquidityManagerAddress: models.Address(common.HexToAddress("0x1")),
					LocalBridgeAdapterAddress:     models.Address(common.HexToAddress("0x4")),
					RemoteTokenAddress:            models.Address(common.HexToAddress("0x3")),
				},
			},
			false,
		},
		{
			"not equal remote token",
			fields{
				RemoteRebalancerAddress:   models.Address(common.HexToAddress("0x1")),
				LocalBridgeAdapterAddress: models.Address(common.HexToAddress("0x2")),
				RemoteTokenAddress:        models.Address(common.HexToAddress("0x3")),
			},
			args{
				other: XChainLiquidityManagerData{
					RemoteLiquidityManagerAddress: models.Address(common.HexToAddress("0x1")),
					LocalBridgeAdapterAddress:     models.Address(common.HexToAddress("0x2")),
					RemoteTokenAddress:            models.Address(common.HexToAddress("0x4")),
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := XChainLiquidityManagerData{
				RemoteLiquidityManagerAddress: tt.fields.RemoteRebalancerAddress,
				LocalBridgeAdapterAddress:     tt.fields.LocalBridgeAdapterAddress,
				RemoteTokenAddress:            tt.fields.RemoteTokenAddress,
			}
			got := d.Equals(tt.args.other)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestData_Equals(t *testing.T) {
	type fields struct {
		Liquidity         *big.Int
		TokenAddress      models.Address
		RebalancerAddress models.Address
		XChainRebalancers map[models.NetworkSelector]XChainLiquidityManagerData
		ConfigDigest      models.ConfigDigest
		NetworkSelector   models.NetworkSelector
	}
	type args struct {
		other Data
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			"equal",
			fields{
				Liquidity:         big.NewInt(100),
				TokenAddress:      models.Address(common.HexToAddress("0x1")),
				RebalancerAddress: models.Address(common.HexToAddress("0x2")),
				XChainRebalancers: map[models.NetworkSelector]XChainLiquidityManagerData{
					models.NetworkSelector(1): {
						RemoteLiquidityManagerAddress: models.Address(common.HexToAddress("0x3")),
						LocalBridgeAdapterAddress:     models.Address(common.HexToAddress("0x4")),
						RemoteTokenAddress:            models.Address(common.HexToAddress("0x5")),
					},
					models.NetworkSelector(2): {
						RemoteLiquidityManagerAddress: models.Address(common.HexToAddress("0x6")),
						LocalBridgeAdapterAddress:     models.Address(common.HexToAddress("0x7")),
						RemoteTokenAddress:            models.Address(common.HexToAddress("0x8")),
					},
				},
				ConfigDigest:    models.ConfigDigest{ConfigDigest: types.ConfigDigest(common.HexToHash("0x9"))},
				NetworkSelector: models.NetworkSelector(3),
			},
			args{
				other: Data{
					Liquidity:               big.NewInt(100),
					TokenAddress:            models.Address(common.HexToAddress("0x1")),
					LiquidityManagerAddress: models.Address(common.HexToAddress("0x2")),
					XChainLiquidityManagers: map[models.NetworkSelector]XChainLiquidityManagerData{
						models.NetworkSelector(1): {
							RemoteLiquidityManagerAddress: models.Address(common.HexToAddress("0x3")),
							LocalBridgeAdapterAddress:     models.Address(common.HexToAddress("0x4")),
							RemoteTokenAddress:            models.Address(common.HexToAddress("0x5")),
						},
						models.NetworkSelector(2): {
							RemoteLiquidityManagerAddress: models.Address(common.HexToAddress("0x6")),
							LocalBridgeAdapterAddress:     models.Address(common.HexToAddress("0x7")),
							RemoteTokenAddress:            models.Address(common.HexToAddress("0x8")),
						},
					},
					ConfigDigest:    models.ConfigDigest{ConfigDigest: types.ConfigDigest(common.HexToHash("0x9"))},
					NetworkSelector: models.NetworkSelector(3),
				},
			},
			true,
		},
		{
			"not equal liquidity",
			fields{
				Liquidity:         big.NewInt(100),
				TokenAddress:      models.Address(common.HexToAddress("0x1")),
				RebalancerAddress: models.Address(common.HexToAddress("0x2")),
				XChainRebalancers: map[models.NetworkSelector]XChainLiquidityManagerData{
					models.NetworkSelector(1): {
						RemoteLiquidityManagerAddress: models.Address(common.HexToAddress("0x3")),
						LocalBridgeAdapterAddress:     models.Address(common.HexToAddress("0x4")),
						RemoteTokenAddress:            models.Address(common.HexToAddress("0x5")),
					},
					models.NetworkSelector(2): {
						RemoteLiquidityManagerAddress: models.Address(common.HexToAddress("0x6")),
						LocalBridgeAdapterAddress:     models.Address(common.HexToAddress("0x7")),
						RemoteTokenAddress:            models.Address(common.HexToAddress("0x8")),
					},
				},
				ConfigDigest:    models.ConfigDigest{ConfigDigest: types.ConfigDigest(common.HexToHash("0x9"))},
				NetworkSelector: models.NetworkSelector(3),
			},
			args{
				other: Data{
					Liquidity:               big.NewInt(200),
					TokenAddress:            models.Address(common.HexToAddress("0x1")),
					LiquidityManagerAddress: models.Address(common.HexToAddress("0x2")),
					XChainLiquidityManagers: map[models.NetworkSelector]XChainLiquidityManagerData{
						models.NetworkSelector(1): {
							RemoteLiquidityManagerAddress: models.Address(common.HexToAddress("0x3")),
							LocalBridgeAdapterAddress:     models.Address(common.HexToAddress("0x4")),
							RemoteTokenAddress:            models.Address(common.HexToAddress("0x5")),
						},
						models.NetworkSelector(2): {
							RemoteLiquidityManagerAddress: models.Address(common.HexToAddress("0x6")),
							LocalBridgeAdapterAddress:     models.Address(common.HexToAddress("0x7")),
							RemoteTokenAddress:            models.Address(common.HexToAddress("0x8")),
						},
					},
					ConfigDigest:    models.ConfigDigest{ConfigDigest: types.ConfigDigest(common.HexToHash("0x9"))},
					NetworkSelector: models.NetworkSelector(3),
				},
			},
			false,
		},
		{
			"not equal token address",
			fields{
				Liquidity:         big.NewInt(100),
				TokenAddress:      models.Address(common.HexToAddress("0x1")),
				RebalancerAddress: models.Address(common.HexToAddress("0x2")),
				XChainRebalancers: map[models.NetworkSelector]XChainLiquidityManagerData{
					models.NetworkSelector(1): {
						RemoteLiquidityManagerAddress: models.Address(common.HexToAddress("0x3")),
						LocalBridgeAdapterAddress:     models.Address(common.HexToAddress("0x4")),
						RemoteTokenAddress:            models.Address(common.HexToAddress("0x5")),
					},
					models.NetworkSelector(2): {
						RemoteLiquidityManagerAddress: models.Address(common.HexToAddress("0x6")),
						LocalBridgeAdapterAddress:     models.Address(common.HexToAddress("0x7")),
						RemoteTokenAddress:            models.Address(common.HexToAddress("0x8")),
					},
				},
				ConfigDigest:    models.ConfigDigest{ConfigDigest: types.ConfigDigest(common.HexToHash("0x9"))},
				NetworkSelector: models.NetworkSelector(3),
			},
			args{
				other: Data{
					Liquidity:               big.NewInt(100),
					TokenAddress:            models.Address(common.HexToAddress("0x22")),
					LiquidityManagerAddress: models.Address(common.HexToAddress("0x2")),
					XChainLiquidityManagers: map[models.NetworkSelector]XChainLiquidityManagerData{
						models.NetworkSelector(1): {
							RemoteLiquidityManagerAddress: models.Address(common.HexToAddress("0x3")),
							LocalBridgeAdapterAddress:     models.Address(common.HexToAddress("0x4")),
							RemoteTokenAddress:            models.Address(common.HexToAddress("0x5")),
						},
						models.NetworkSelector(2): {
							RemoteLiquidityManagerAddress: models.Address(common.HexToAddress("0x6")),
							LocalBridgeAdapterAddress:     models.Address(common.HexToAddress("0x7")),
							RemoteTokenAddress:            models.Address(common.HexToAddress("0x8")),
						},
					},
					ConfigDigest:    models.ConfigDigest{ConfigDigest: types.ConfigDigest(common.HexToHash("0x9"))},
					NetworkSelector: models.NetworkSelector(3),
				},
			},
			false,
		},
		{
			"not equal rebalancer address",
			fields{
				Liquidity:         big.NewInt(100),
				TokenAddress:      models.Address(common.HexToAddress("0x1")),
				RebalancerAddress: models.Address(common.HexToAddress("0x2")),
				XChainRebalancers: map[models.NetworkSelector]XChainLiquidityManagerData{
					models.NetworkSelector(1): {
						RemoteLiquidityManagerAddress: models.Address(common.HexToAddress("0x3")),
						LocalBridgeAdapterAddress:     models.Address(common.HexToAddress("0x4")),
						RemoteTokenAddress:            models.Address(common.HexToAddress("0x5")),
					},
					models.NetworkSelector(2): {
						RemoteLiquidityManagerAddress: models.Address(common.HexToAddress("0x6")),
						LocalBridgeAdapterAddress:     models.Address(common.HexToAddress("0x7")),
						RemoteTokenAddress:            models.Address(common.HexToAddress("0x8")),
					},
				},
				ConfigDigest:    models.ConfigDigest{ConfigDigest: types.ConfigDigest(common.HexToHash("0x9"))},
				NetworkSelector: models.NetworkSelector(3),
			},
			args{
				other: Data{
					Liquidity:               big.NewInt(100),
					TokenAddress:            models.Address(common.HexToAddress("0x1")),
					LiquidityManagerAddress: models.Address(common.HexToAddress("0x222")),
					XChainLiquidityManagers: map[models.NetworkSelector]XChainLiquidityManagerData{
						models.NetworkSelector(1): {
							RemoteLiquidityManagerAddress: models.Address(common.HexToAddress("0x3")),
							LocalBridgeAdapterAddress:     models.Address(common.HexToAddress("0x4")),
							RemoteTokenAddress:            models.Address(common.HexToAddress("0x5")),
						},
						models.NetworkSelector(2): {
							RemoteLiquidityManagerAddress: models.Address(common.HexToAddress("0x6")),
							LocalBridgeAdapterAddress:     models.Address(common.HexToAddress("0x7")),
							RemoteTokenAddress:            models.Address(common.HexToAddress("0x8")),
						},
					},
					ConfigDigest:    models.ConfigDigest{ConfigDigest: types.ConfigDigest(common.HexToHash("0x9"))},
					NetworkSelector: models.NetworkSelector(3),
				},
			},
			false,
		},
		{
			"not equal xchain rebalancers",
			fields{
				Liquidity:         big.NewInt(100),
				TokenAddress:      models.Address(common.HexToAddress("0x1")),
				RebalancerAddress: models.Address(common.HexToAddress("0x2")),
				XChainRebalancers: map[models.NetworkSelector]XChainLiquidityManagerData{
					models.NetworkSelector(1): {
						RemoteLiquidityManagerAddress: models.Address(common.HexToAddress("0x3")),
						LocalBridgeAdapterAddress:     models.Address(common.HexToAddress("0x4")),
						RemoteTokenAddress:            models.Address(common.HexToAddress("0x5")),
					},
					models.NetworkSelector(2): {
						RemoteLiquidityManagerAddress: models.Address(common.HexToAddress("0x6")),
						LocalBridgeAdapterAddress:     models.Address(common.HexToAddress("0x7")),
						RemoteTokenAddress:            models.Address(common.HexToAddress("0x8")),
					},
				},
				ConfigDigest:    models.ConfigDigest{ConfigDigest: types.ConfigDigest(common.HexToHash("0x9"))},
				NetworkSelector: models.NetworkSelector(3),
			},
			args{
				other: Data{
					Liquidity:               big.NewInt(100),
					TokenAddress:            models.Address(common.HexToAddress("0x1")),
					LiquidityManagerAddress: models.Address(common.HexToAddress("0x222")),
					XChainLiquidityManagers: map[models.NetworkSelector]XChainLiquidityManagerData{
						models.NetworkSelector(1): {
							RemoteLiquidityManagerAddress: models.Address(common.HexToAddress("0x33")),
							LocalBridgeAdapterAddress:     models.Address(common.HexToAddress("0x4")),
							RemoteTokenAddress:            models.Address(common.HexToAddress("0x5")),
						},
						models.NetworkSelector(2): {
							RemoteLiquidityManagerAddress: models.Address(common.HexToAddress("0x6")),
							LocalBridgeAdapterAddress:     models.Address(common.HexToAddress("0x7")),
							RemoteTokenAddress:            models.Address(common.HexToAddress("0x8")),
						},
					},
					ConfigDigest:    models.ConfigDigest{ConfigDigest: types.ConfigDigest(common.HexToHash("0x9"))},
					NetworkSelector: models.NetworkSelector(3),
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Data{
				Liquidity:               tt.fields.Liquidity,
				TokenAddress:            tt.fields.TokenAddress,
				LiquidityManagerAddress: tt.fields.RebalancerAddress,
				XChainLiquidityManagers: tt.fields.XChainRebalancers,
				ConfigDigest:            tt.fields.ConfigDigest,
				NetworkSelector:         tt.fields.NetworkSelector,
			}
			if got := d.Equals(tt.args.other); got != tt.want {
				t.Errorf("Data.Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestData_Clone(t *testing.T) {
	d := Data{
		Liquidity:               big.NewInt(100),
		TokenAddress:            models.Address(common.HexToAddress("0x1")),
		LiquidityManagerAddress: models.Address(common.HexToAddress("0x2")),
		XChainLiquidityManagers: map[models.NetworkSelector]XChainLiquidityManagerData{
			models.NetworkSelector(1): {
				RemoteLiquidityManagerAddress: models.Address(common.HexToAddress("0x3")),
				LocalBridgeAdapterAddress:     models.Address(common.HexToAddress("0x4")),
				RemoteTokenAddress:            models.Address(common.HexToAddress("0x5")),
			},
			models.NetworkSelector(2): {
				RemoteLiquidityManagerAddress: models.Address(common.HexToAddress("0x6")),
				LocalBridgeAdapterAddress:     models.Address(common.HexToAddress("0x7")),
				RemoteTokenAddress:            models.Address(common.HexToAddress("0x8")),
			},
		},
		ConfigDigest:    models.ConfigDigest{ConfigDigest: types.ConfigDigest(common.HexToHash("0x9"))},
		NetworkSelector: models.NetworkSelector(3),
	}
	clone := d.Clone()
	require.True(t, d.Equals(clone))
	clone.Liquidity.Set(big.NewInt(200))
	require.False(t, d.Equals(clone))
}
