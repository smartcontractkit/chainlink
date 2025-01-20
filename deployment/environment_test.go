package deployment

import (
	"reflect"
	"testing"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
)

func TestNode_OCRConfigForChainSelector(t *testing.T) {
	var m = map[chain_selectors.ChainDetails]OCRConfig{
		{
			ChainSelector: chain_selectors.APTOS_TESTNET.Selector,
			ChainName:     chain_selectors.APTOS_TESTNET.Name,
		}: {
			KeyBundleID: "aptos bundle 1",
		},
		{
			ChainSelector: chain_selectors.ETHEREUM_MAINNET_ARBITRUM_1.Selector,
			ChainName:     chain_selectors.ETHEREUM_MAINNET_ARBITRUM_1.Name,
		}: {
			KeyBundleID: "arb bundle 1",
		},
	}

	type fields struct {
		SelToOCRConfig map[chain_selectors.ChainDetails]OCRConfig
	}
	type args struct {
		chainSel uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   OCRConfig
		exist  bool
	}{
		{
			name: "aptos ok",
			fields: fields{
				SelToOCRConfig: m,
			},
			args: args{
				chainSel: chain_selectors.APTOS_TESTNET.Selector,
			},
			want: OCRConfig{
				KeyBundleID: "aptos bundle 1",
			},
			exist: true,
		},
		{
			name: "arb ok",
			fields: fields{
				SelToOCRConfig: m,
			},
			args: args{
				chainSel: chain_selectors.ETHEREUM_MAINNET_ARBITRUM_1.Selector,
			},
			want: OCRConfig{
				KeyBundleID: "arb bundle 1",
			},
			exist: true,
		},
		{
			name: "no exist",
			fields: fields{
				SelToOCRConfig: m,
			},
			args: args{
				chainSel: chain_selectors.WEMIX_MAINNET.Selector, // not in test data
			},
			want:  OCRConfig{},
			exist: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Node{
				SelToOCRConfig: tt.fields.SelToOCRConfig,
			}
			got, got1 := n.OCRConfigForChainSelector(tt.args.chainSel)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Node.OCRConfigForChainSelector() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.exist {
				t.Errorf("Node.OCRConfigForChainSelector() got1 = %v, want %v", got1, tt.exist)
			}
		})
	}
}
