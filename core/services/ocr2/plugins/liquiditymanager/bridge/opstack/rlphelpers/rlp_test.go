package rlphelpers_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/bridge/opstack/rlphelpers"
)

func TestDecodeRLP(t *testing.T) {
	type testCase struct {
		name     string
		inputHex string
		expected func() *rlphelpers.RawRLPOutput
	}
	tests := []testCase{
		{
			name:     "full node",
			inputHex: "0xf901f1a09082d573c1e008958410a5cdb60327ba4674a1d8022996550350b61385a2a26fa0c56ce921f714cc12af81ade7d45abcb479053e0627ee04e9cf38e1726e5581fca00fc973598f190f660fc872c69ffff607b09abebaced46893bf10d8e8fdb7a660a0ae0bd1cd7e0ec5f09aa0e82cb17b77e70a00632359d5fe9e9c7f42bbea318127a0a0b7b3d960d20416a1ced8936ab898d553f6e212b7b04801effcc7e7b19cf1a9a0b559523e82faa2e408af50372465d71ebd089b93dc16d9cbdc7a7d82f44a3c11a0febfb17ab0cff682a4ece5d91b5daca2e6dce14e3abc578915717be15fee093da0c72915d2155b058565b180294542d22c814290e73d3d7b5e084ba351b3b7e347a088cf20bf4db4aef51b599e6bf657b60fd546ae5dcf6e151d8b8bde6c60298c79a0323bc46e9470693392ba46eac91bdedbc8beed42cca7fe73dd285f6c2c133551a097cd1b854df44b3faa578a354ad3810dfba851980e51c70e29b90e321ad61071a04f680830e90cbbd7ee07c1159bbddf17912ad2490e4007f5ada73c4d16e5e8fb80a081bde169fb3bd472384ec23fdd14e2df60b982328d6b3feb67f0971d722a653ba03d66c0c21895308ca5cbcc1236ece7bf5f2fcae9ecef768a2310b98cb450b860a043d9c95db22d254c24ef42bd9c0abc826eb35ec84df97db31991da0a2839b4aa80",
			expected: func() *rlphelpers.RawRLPOutput {
				root := rlphelpers.NewRLPBuffers()
				root.Children = []*rlphelpers.RawRLPOutput{
					{Data: hexutil.MustDecode("0x9082d573c1e008958410a5cdb60327ba4674a1d8022996550350b61385a2a26f")},
					{Data: hexutil.MustDecode("0xc56ce921f714cc12af81ade7d45abcb479053e0627ee04e9cf38e1726e5581fc")},
					{Data: hexutil.MustDecode("0x0fc973598f190f660fc872c69ffff607b09abebaced46893bf10d8e8fdb7a660")},
					{Data: hexutil.MustDecode("0xae0bd1cd7e0ec5f09aa0e82cb17b77e70a00632359d5fe9e9c7f42bbea318127")},
					{Data: hexutil.MustDecode("0xa0b7b3d960d20416a1ced8936ab898d553f6e212b7b04801effcc7e7b19cf1a9")},
					{Data: hexutil.MustDecode("0xb559523e82faa2e408af50372465d71ebd089b93dc16d9cbdc7a7d82f44a3c11")},
					{Data: hexutil.MustDecode("0xfebfb17ab0cff682a4ece5d91b5daca2e6dce14e3abc578915717be15fee093d")},
					{Data: hexutil.MustDecode("0xc72915d2155b058565b180294542d22c814290e73d3d7b5e084ba351b3b7e347")},
					{Data: hexutil.MustDecode("0x88cf20bf4db4aef51b599e6bf657b60fd546ae5dcf6e151d8b8bde6c60298c79")},
					{Data: hexutil.MustDecode("0x323bc46e9470693392ba46eac91bdedbc8beed42cca7fe73dd285f6c2c133551")},
					{Data: hexutil.MustDecode("0x97cd1b854df44b3faa578a354ad3810dfba851980e51c70e29b90e321ad61071")},
					{Data: hexutil.MustDecode("0x4f680830e90cbbd7ee07c1159bbddf17912ad2490e4007f5ada73c4d16e5e8fb")},
					{},
					{Data: hexutil.MustDecode("0x81bde169fb3bd472384ec23fdd14e2df60b982328d6b3feb67f0971d722a653b")},
					{Data: hexutil.MustDecode("0x3d66c0c21895308ca5cbcc1236ece7bf5f2fcae9ecef768a2310b98cb450b860")},
					{Data: hexutil.MustDecode("0x43d9c95db22d254c24ef42bd9c0abc826eb35ec84df97db31991da0a2839b4aa")},
					{},
				}
				return root
			},
		},
		{
			name:     "short node",
			inputHex: "0xf851a05a46e7120a6b977a29dcd8b6b236ce25fac8060db19fbbe4bb32f9abfae356c2808080808080a0ee4f8880fbda3060806bd36ead0015717db1ea74f0915436a287ed9ed3091917808080808080808080",
			expected: func() *rlphelpers.RawRLPOutput {
				root := rlphelpers.NewRLPBuffers()
				root.Children = []*rlphelpers.RawRLPOutput{
					{Data: hexutil.MustDecode("0x5a46e7120a6b977a29dcd8b6b236ce25fac8060db19fbbe4bb32f9abfae356c2")},
					{},
					{},
					{},
					{},
					{},
					{},
					{Data: hexutil.MustDecode("0xee4f8880fbda3060806bd36ead0015717db1ea74f0915436a287ed9ed3091917")},
					{},
					{},
					{},
					{},
					{},
					{},
					{},
					{},
					{},
				}
				return root
			},
		},
		{
			name:     "branch node",
			inputHex: "0xf84d8080808080de9c332c35a4d03ec6ab9b3ffd06c69652ce8e02ff95537f98b7a0feb29c01808080de9c36620d27102d4799d259fe08f690c0a21b80d0a1a903db682417a23f0180808080808080",
			expected: func() *rlphelpers.RawRLPOutput {
				root := rlphelpers.NewRLPBuffers()
				root.Children = []*rlphelpers.RawRLPOutput{
					{},
					{},
					{},
					{},
					{},
					{Children: []*rlphelpers.RawRLPOutput{
						{Data: hexutil.MustDecode("0x332c35a4d03ec6ab9b3ffd06c69652ce8e02ff95537f98b7a0feb29c")},
						{Data: hexutil.MustDecode("0x01")},
					}},
					{},
					{},
					{},
					{Children: []*rlphelpers.RawRLPOutput{
						{Data: hexutil.MustDecode("0x36620d27102d4799d259fe08f690c0a21b80d0a1a903db682417a23f")},
						{Data: hexutil.MustDecode("0x01")},
					}},
					{},
					{},
					{},
					{},
					{},
					{},
					{},
				}
				return root
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			decoded := hexutil.MustDecode(test.inputHex)
			bufs := rlphelpers.NewRLPBuffers()
			err := rlp.DecodeBytes(decoded, bufs)
			require.NoError(t, err)
			require.True(t, test.expected().Equal(bufs))
		})
	}
}
