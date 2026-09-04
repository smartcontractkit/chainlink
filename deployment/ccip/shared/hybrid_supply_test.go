package shared_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
)

// sixDecimalSupply is a realistic supply for a 6-decimal token, used to exercise the
// 6-to-18 boundary that a same-decimals test would miss.
var sixDecimalSupply = big.NewInt(1_206_841_810_125_844)

func mustBigInt(t *testing.T, s string) *big.Int {
	t.Helper()
	v, ok := new(big.Int).SetString(s, 10)
	require.True(t, ok, "failed to parse %s", s)
	return v
}

func TestNormalizeRemoteSupply(t *testing.T) {
	t.Parallel()

	maxUint256 := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))

	tests := []struct {
		name           string
		amount         *big.Int
		remoteDecimals uint8
		localDecimals  uint8
		want           *big.Int
		wantErrMsg     string
	}{
		{
			name:           "Success - equal decimals returns the amount unchanged",
			amount:         big.NewInt(1000),
			remoteDecimals: 18,
			localDecimals:  18,
			want:           big.NewInt(1000),
		},
		{
			name:           "Success - six decimal remote scales up to eighteen decimal local",
			amount:         sixDecimalSupply,
			remoteDecimals: 6,
			localDecimals:  18,
			want:           mustBigInt(t, "1206841810125844000000000000"),
		},
		{
			name:           "Success - eighteen decimal remote scales down to six decimal local",
			amount:         mustBigInt(t, "1206841810125844000000000000"),
			remoteDecimals: 18,
			localDecimals:  6,
			want:           sixDecimalSupply,
		},
		{
			name:           "Success - zero supply across differing decimals",
			amount:         big.NewInt(0),
			remoteDecimals: 6,
			localDecimals:  18,
			want:           big.NewInt(0),
		},
		{
			name:           "Failure - scaling down with a remainder",
			amount:         big.NewInt(1_000_001),
			remoteDecimals: 18,
			localDecimals:  12,
			wantErrMsg:     "no exact representation at 12 local decimals",
		},
		{
			name:           "Failure - nil supply",
			amount:         nil,
			remoteDecimals: 6,
			localDecimals:  18,
			wantErrMsg:     "must not be nil",
		},
		{
			name:           "Failure - negative supply",
			amount:         big.NewInt(-1),
			remoteDecimals: 6,
			localDecimals:  18,
			wantErrMsg:     "must not be negative",
		},
		{
			name:           "Failure - scaling up past uint256",
			amount:         maxUint256,
			remoteDecimals: 6,
			localDecimals:  18,
			wantErrMsg:     "overflows uint256",
		},
		{
			name:           "Failure - scale factor exceeds uint256 when scaling up",
			amount:         big.NewInt(1),
			remoteDecimals: 0,
			localDecimals:  255,
			wantErrMsg:     "scale factor exceeds uint256",
		},
		{
			name:           "Failure - scale factor exceeds uint256 when scaling down",
			amount:         big.NewInt(1),
			remoteDecimals: 255,
			localDecimals:  0,
			wantErrMsg:     "scale factor exceeds uint256",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := shared.NormalizeRemoteSupply(tt.amount, tt.remoteDecimals, tt.localDecimals)
			if tt.wantErrMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErrMsg)
				return
			}

			require.NoError(t, err)
			require.Zero(t, tt.want.Cmp(got), "want %s, got %s", tt.want, got)
		})
	}
}

// TestNormalizeRemoteSupply_DoesNotMutateInput guards against aliasing a caller's big.Int.
func TestNormalizeRemoteSupply_DoesNotMutateInput(t *testing.T) {
	t.Parallel()

	amount := new(big.Int).Set(sixDecimalSupply)

	got, err := shared.NormalizeRemoteSupply(amount, 6, 18)
	require.NoError(t, err)
	require.Zero(t, sixDecimalSupply.Cmp(amount), "input was mutated")
	require.NotSame(t, amount, got, "result aliases the input")

	got, err = shared.NormalizeRemoteSupply(amount, 18, 18)
	require.NoError(t, err)
	require.NotSame(t, amount, got, "equal-decimals result aliases the input")
}

func TestValidateRemoteChainSupply(t *testing.T) {
	t.Parallel()

	const remoteChainSelector = 124_615_329_519_749_607

	tests := []struct {
		name          string
		remoteSupply  shared.RemoteSupply
		suppliedLocal *big.Int
		localDecimals uint8
		wantErrMsg    string
	}{
		{
			name: "Success - correctly normalized six to eighteen decimals",
			remoteSupply: shared.RemoteSupply{
				Decimals: 6,
				Amount:   sixDecimalSupply,
			},
			suppliedLocal: mustBigInt(t, "1206841810125844000000000000"),
			localDecimals: 18,
		},
		{
			// The unconverted value under-backs the pool by 10^12.
			name: "Failure - remote supply passed through without conversion",
			remoteSupply: shared.RemoteSupply{
				Decimals: 6,
				Amount:   sixDecimalSupply,
			},
			suppliedLocal: sixDecimalSupply,
			localDecimals: 18,
			wantErrMsg:    "requires 1206841810125844000000000000, but the update supplies 1206841810125844",
		},
		{
			name: "Success - matching decimals accepts the supply unchanged",
			remoteSupply: shared.RemoteSupply{
				Decimals: 18,
				Amount:   big.NewInt(1000),
			},
			suppliedLocal: big.NewInt(1000),
			localDecimals: 18,
		},
		{
			name: "Success - zero supply matches a zero update",
			remoteSupply: shared.RemoteSupply{
				Decimals: 6,
				Amount:   big.NewInt(0),
			},
			suppliedLocal: big.NewInt(0),
			localDecimals: 18,
		},
		{
			name: "Failure - nil supplied amount",
			remoteSupply: shared.RemoteSupply{
				Decimals: 6,
				Amount:   sixDecimalSupply,
			},
			suppliedLocal: nil,
			localDecimals: 18,
			wantErrMsg:    "must not be nil",
		},
		{
			name: "Failure - nil remote supply amount",
			remoteSupply: shared.RemoteSupply{
				Decimals: 6,
				Amount:   nil,
			},
			suppliedLocal: big.NewInt(1000),
			localDecimals: 18,
			wantErrMsg:    "must not be nil",
		},
		{
			name: "Failure - off by one against a correct conversion",
			remoteSupply: shared.RemoteSupply{
				Decimals: 6,
				Amount:   big.NewInt(1_000_000),
			},
			suppliedLocal: mustBigInt(t, "999999999999999999"),
			localDecimals: 18,
			wantErrMsg:    "requires 1000000000000000000, but the update supplies 999999999999999999",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := shared.ValidateRemoteChainSupply(
				tt.remoteSupply, tt.suppliedLocal, tt.localDecimals, remoteChainSelector,
			)
			if tt.wantErrMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErrMsg)
				require.Contains(t, err.Error(), "124615329519749607", "error should name the remote chain")
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestVerifyRemoteDecimals(t *testing.T) {
	t.Parallel()

	const (
		ethereumSelector = uint64(5_009_297_550_715_157_269)
		solanaSelector   = uint64(124_615_329_519_749_607)
		aptosSelector    = uint64(4_741_433_654_826_277_614)
	)

	tokenAddress := make([]byte, 32)
	tokenAddress[31] = 1

	tests := []struct {
		name          string
		chainSelector uint64
		tokenAddress  []byte
		wantVerified  bool
		wantErrMsg    string
	}{
		{
			name:          "Success - empty remote token address is unverifiable",
			chainSelector: ethereumSelector,
			tokenAddress:  nil,
		},
		{
			// Each family reports unverified rather than failing when its chain is absent,
			// so a migration is never blocked by an environment that omits the remote.
			name:          "Success - EVM chain absent from environment is unverifiable",
			chainSelector: ethereumSelector,
			tokenAddress:  tokenAddress,
		},
		{
			name:          "Success - SVM chain absent from environment is unverifiable",
			chainSelector: solanaSelector,
			tokenAddress:  tokenAddress,
		},
		{
			name:          "Success - Aptos chain absent from environment is unverifiable",
			chainSelector: aptosSelector,
			tokenAddress:  tokenAddress,
		},
		{
			name:          "Failure - unknown chain selector",
			chainSelector: 1,
			tokenAddress:  tokenAddress,
			wantErrMsg:    "failed to resolve chain family",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			verified, err := shared.VerifyRemoteDecimals(
				t.Context(), cldf.Environment{}, tt.chainSelector, tt.tokenAddress, 6,
			)
			if tt.wantErrMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.wantErrMsg)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantVerified, verified)
		})
	}
}
