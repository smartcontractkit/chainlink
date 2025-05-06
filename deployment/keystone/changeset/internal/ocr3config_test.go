// TODO: KS-458 copied from https://github.com/smartcontractkit/chainlink/blob/65924811dc53a211613927c814d7f04fd85439a4/core/scripts/keystone/src/88_gen_ocr3_config.go#L1
// to unblock go mod issues when trying to import the scripts package
package internal

import (
	"encoding/json"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	types2 "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	types3 "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/require"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/view"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

var wantOCR3Config = `{
  "Signers": [
    "011400b35409a8d4f9a18da55c5b2bb08a3f5f68d44442052000b8834eaa062f0df4ccfe7832253920071ec14dc4f78b13ecdda10b824e2dd3b6",
    "0114008258f4c4761cc445333017608044a204fd0c006a052000247d0189f65f58be83a4e7d87ff338aaf8956e9acb9fcc783f34f9edc29d1b40",
    "011400d4dcc573e9d24a8b27a07bba670ba3a2ab36e5bb052000ba20d3da9b07663f1e8039081a514649fd61a48be2d241bc63537ee47d028fcd",
    "0114006607c140e558631407f33bafbabd103863cee876052000046faf34ebfe42510251e6098bc34fa3dd5f2de38ac07e47f2d1b34ac770639f",
    "011400a6f35436cb7bffd615cc47a0a04aa0a78696a1440520001221e131ef21014a6a99ed22376eb869746a3b5e30fd202cf79e44efaeb8c5c2",
    "011400657587eb55cecd6f90b97297b611c3024e488cc0052000425d1354a7b8180252a221040c718cac0ba0251c7efe31a2acefbba578dc2153",
    "0114004885973b2fcf061d5cdfb8f74c5139bd3056e9da0520004a94c75cb9fe8b1fba86fd4b71ad130943281fdefad10216c46eb2285d60950f",
    "011400213803bb9f9715379aaf11aadb0212369701dc0a05200096dc85670c49caa986de4ad288e680e9afb0f5491160dcbb4868ca718e194fc8",
    "0114008c2aa1e6fad88a6006dfb116eb866cbad2910314052000bddafb20cc50d89e0ae2f244908c27b1d639615d8186b28c357669de3359f208",
    "011400679296b7c1eb4948efcc87efc550940a182e610c0520004fa557850e4d5c21b3963c97414c1f37792700c4d3b8abdb904b765fd47e39bf"
  ],
  "Transmitters": [
    "0x2877F08d9c5Cc9F401F730Fa418fAE563A9a2FF3",
    "0x415aa1E9a1bcB3929ed92bFa1F9735Dc0D45AD31",
    "0xCea84bC1881F3cE14BA13Dc3a00DC1Ff3D553fF0",
    "0xA9eFB53c513E413762b2Be5299D161d8E6e7278e",
    "0x6F5cAb24Fb7412bB516b3468b9F3a9c471d25fE5",
    "0xdAd1F3F8ec690cf335D46c50EdA5547CeF875161",
    "0x19e10B063a62B1574AE19020A64fDe6419892dA6",
    "0x9ad9f3AD49e5aB0F28bD694d211a90297bD90D7f",
    "0x31B179dcF8f9036C30f04bE578793e51bF14A39E",
    "0x0b04cE574E80Da73191Ec141c0016a54A6404056"
  ],
  "F": 3,
  "OnchainConfig": "0x",
  "OffchainConfigVersion": 30,
  "OffchainConfig": "0xc80180e497d012d00180e497d012d80180a8d6b907e00180cab5ee01e80180d88ee16ff0010afa01010a82022003dacd15fc96c965c648e3623180de002b71a97cf6eeca9affb91f461dcd6ce1820220255096a3b7ade10e29c648e0b407fc486180464f713446b1da04f013df6179c8820220dba3c61e5f8bec594be481bcaf67ecea0d1c2950edb15b158ce3dbc77877def3820220b4c4993d6c15fee63800db901a8b35fa419057610962caab1c1d7bed557091278202202a4c7dec127fdd8145e48c5edb9467225098bd8c8ad1dade868325b787affbde820220283471ed66d61fbe11f64eff65d738b59a0301c9a4f846280db26c64c9fdd3f8820220aa3419628ea3536783742d17d8adf05681aa6a6bd2b206fbde78c7e5aa38586d82022001496edce35663071d74472e02119432ba059b3904d205e4358014410e4f2be3820220ad08c2a5878cada53521f4e2bb449f191ccca7899246721a0deeea19f7b83f70820220c805572b813a072067eab2087ddbee8aa719090e12890b15c01094f0d3f74a5f8a02008a02008a02008a02008a02008a02008a02008a02008a02008a020092021b08c0843d10c0843d18c0843d20c0843d2814301438901c4202081e98028094ebdc03a0028094ebdc03a8028094ebdc03b0028094ebdc03ba02f8010a20da47a8cc1c10796dd43f98ed113c648625e2e504c16ac5da9c65669e2377241b1220f5beca3bb11406079dc174183105c474c862a73c257ce8b3d9f5ca065e6264691a10805015e4203740495a23e93c1bd06ba81a10ca58ff36ffb0545dc3f800ddd6f8d0481a1076f664639ca8b5209e488895faa5460f1a104a1e89a7f2d8c89158f18856bf289c2a1a10c2f4330787831f419713ad4990e347d31a10fd403ec0797c001a2794b51d6178916d1a10e14fff88fdd3d1554ed861104ddc56a81a10b0284b9817fec2c3066c6f2651d17fc41a10b090233a67d502f78191c9e19a2a032b1a10e483414860bb612af50ee15ce8cd8ef5c00280e497d012c8028094ebdc03"
}`

var ocr3Cfg = `
{
  "MaxQueryLengthBytes": 1000000,
  "MaxObservationLengthBytes": 1000000,
  "MaxReportLengthBytes": 1000000,
  "MaxOutcomeLengthBytes": 1000000,
  "MaxReportCount": 20,
  "MaxBatchSize": 20,
  "OutcomePruningThreshold": 3600,
  "UniqueReports": true,
  "RequestTimeout": "30s",
  "DeltaProgressMillis": 5000,
  "DeltaResendMillis": 5000,
  "DeltaInitialMillis": 5000,
  "DeltaRoundMillis": 2000,
  "DeltaGraceMillis": 500,
  "DeltaCertifiedCommitRequestMillis": 1000,
  "DeltaStageMillis": 30000,
  "MaxRoundsPerEpoch": 10,
  "TransmissionSchedule": [
    10
  ],
  "MaxDurationQueryMillis": 1000,
  "MaxDurationObservationMillis": 1000,
  "MaxDurationReportMillis": 1000,
  "MaxDurationShouldAcceptMillis": 1000,
  "MaxDurationShouldTransmitMillis": 1000,
  "MaxFaultyOracles": 3
}`

func Test_configureOCR3Request_generateOCR3Config(t *testing.T) {
	nodes := loadTestData(t, "testdata/testnet_wf_view.json")

	var cfg OracleConfig
	err := json.Unmarshal([]byte(ocr3Cfg), &cfg)
	require.NoError(t, err)
	r := configureOCR3Request{
		cfg:   &cfg,
		nodes: nodes,
		chain: deployment.Chain{
			Selector: chain_selectors.ETHEREUM_TESTNET_SEPOLIA.Selector,
		},
		ocrSecrets: cldf.XXXGenerateTestOCRSecrets(),
	}
	got, err := r.generateOCR3Config()
	require.NoError(t, err)
	b, err := json.MarshalIndent(got, "", "  ")
	require.NoError(t, err)
	require.Equal(t, wantOCR3Config, string(b))

	t.Run("no multiple transmitters", func(t *testing.T) {
		cfg2 := cfg
		cfg2.TransmissionSchedule = []int{}
		for i := 1; i <= len(nodes); i++ {
			cfg2.TransmissionSchedule = append(cfg2.TransmissionSchedule, i)
		}
		r := configureOCR3Request{
			cfg:   &cfg2,
			nodes: nodes,
			chain: deployment.Chain{
				Selector: chain_selectors.ETHEREUM_TESTNET_SEPOLIA.Selector,
			},
			ocrSecrets: cldf.XXXGenerateTestOCRSecrets(),
		}
		_, err := r.generateOCR3Config()
		require.Error(t, err)
	})
	t.Run("transmitter schedule eqaul num nodes", func(t *testing.T) {
		cfg2 := cfg
		cfg2.TransmissionSchedule = []int{len(nodes) + 1}
		r := configureOCR3Request{
			cfg:   &cfg2,
			nodes: nodes,
			chain: deployment.Chain{
				Selector: chain_selectors.ETHEREUM_TESTNET_SEPOLIA.Selector,
			},
			ocrSecrets: cldf.XXXGenerateTestOCRSecrets(),
		}
		_, err := r.generateOCR3Config()
		require.Error(t, err)
	})
}

func loadTestData(t *testing.T, path string) []deployment.Node {
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	var nodeViews map[string]*view.NopView
	err = json.Unmarshal(data, &nodeViews)
	require.NoError(t, err)
	require.Len(t, nodeViews, 10)

	names := make([]string, 0)
	for k := range nodeViews {
		names = append(names, k)
	}
	sort.Strings(names)

	// in general we can map from the view to the node, but we know the test data
	var nodes []deployment.Node
	// for _, nv := range nodeViews {
	for _, name := range names {
		nv := nodeViews[name]
		node := deployment.Node{
			NodeID:         nv.NodeID,
			IsBootstrap:    nv.IsBootstrap,
			SelToOCRConfig: make(map[chain_selectors.ChainDetails]deployment.OCRConfig),
			AdminAddr:      nv.PayeeAddress,
		}
		for chain, ocrKey := range nv.OCRKeys {
			// TODO: this decoding could be shared with NodeInfo
			p, err := p2pkey.MakePeerID(ocrKey.PeerID)
			require.NoError(t, err)

			b := common.Hex2Bytes(ocrKey.OffchainPublicKey)
			var opk types2.OffchainPublicKey
			copy(opk[:], b)

			b = common.Hex2Bytes(ocrKey.ConfigEncryptionPublicKey)
			var cpk types3.ConfigEncryptionPublicKey
			copy(cpk[:], b)

			var pubkey types3.OnchainPublicKey
			if strings.HasPrefix(chain, "ethereum") {
				// convert from pubkey to address
				pubkey = common.HexToAddress(ocrKey.OnchainPublicKey).Bytes()
			} else {
				pubkey = common.Hex2Bytes(ocrKey.OnchainPublicKey)
			}

			ocrCfg := deployment.OCRConfig{
				KeyBundleID:               ocrKey.KeyBundleID,
				OffchainPublicKey:         opk,
				OnchainPublicKey:          pubkey,
				PeerID:                    p,
				TransmitAccount:           types.Account(ocrKey.TransmitAccount),
				ConfigEncryptionPublicKey: cpk,
			}
			var k chain_selectors.ChainDetails
			switch chain {
			case "aptos-testnet":
				k = chain_selectors.ChainDetails{
					ChainSelector: chain_selectors.APTOS_TESTNET.Selector,
					ChainName:     chain,
				}

			case "ethereum-testnet-sepolia":
				k = chain_selectors.ChainDetails{
					ChainSelector: chain_selectors.ETHEREUM_TESTNET_SEPOLIA.Selector,
					ChainName:     chain,
				}
			default:
				t.Fatalf("unexpected chain %s", chain)
			}
			node.SelToOCRConfig[k] = ocrCfg
		}

		nodes = append(nodes, node)
	}
	require.Len(t, nodes, 10)
	return nodes
}
