package job

import (
	_ "embed"
	"reflect"
	"testing"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

func TestOCR2OracleSpec_RelayIdentifier(t *testing.T) {
	type fields struct {
		Relay       relay.Network
		ChainID     string
		RelayConfig JSONConfig
	}
	tests := []struct {
		name    string
		fields  fields
		want    relay.ID
		wantErr bool
	}{
		{name: "err no chain id",
			fields:  fields{},
			want:    relay.ID{},
			wantErr: true,
		},
		{
			name: "evm explicitly configured",
			fields: fields{
				Relay:   relay.EVM,
				ChainID: "1",
			},
			want: relay.ID{Network: relay.EVM, ChainID: "1"},
		},
		{
			name: "evm implicitly configured",
			fields: fields{
				Relay:       relay.EVM,
				RelayConfig: map[string]any{"chainID": 1},
			},
			want: relay.ID{Network: relay.EVM, ChainID: "1"},
		},
		{
			name: "evm implicitly configured with bad value",
			fields: fields{
				Relay:       relay.EVM,
				RelayConfig: map[string]any{"chainID": float32(1)},
			},
			want:    relay.ID{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &OCR2OracleSpec{
				Relay:       tt.fields.Relay,
				ChainID:     tt.fields.ChainID,
				RelayConfig: tt.fields.RelayConfig,
			}
			got, err := s.RelayID()
			if (err != nil) != tt.wantErr {
				t.Errorf("OCR2OracleSpec.RelayIdentifier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OCR2OracleSpec.RelayIdentifier() = %v, want %v", got, tt.want)
			}
		})
	}
}

var (
	//go:embed testdata/compact.toml
	compact string
	//go:embed testdata/pretty.toml
	pretty string
)

func TestOCR2OracleSpec(t *testing.T) {
	val := OCR2OracleSpec{
		Relay:                             relay.EVM,
		PluginType:                        types.Median,
		ContractID:                        "foo",
		OCRKeyBundleID:                    null.StringFrom("bar"),
		TransmitterID:                     null.StringFrom("baz"),
		ContractConfigConfirmations:       1,
		ContractConfigTrackerPollInterval: *models.NewInterval(time.Second),
		RelayConfig: map[string]interface{}{
			"chainID":   1337,
			"fromBlock": 42,
			"chainReader": evmtypes.ChainReaderConfig{
				ChainContractReaders: map[string]evmtypes.ChainContractReader{
					"median": {
						ContractABI: `[{"inputs":[{"internalType":"contractLinkTokenInterface","name":"link","type":"address"},{"internalType":"int192","name":"minAnswer_","type":"int192"},{"internalType":"int192","name":"maxAnswer_","type":"int192"},{"internalType":"contractAccessControllerInterface","name":"billingAccessController","type":"address"},{"internalType":"contractAccessControllerInterface","name":"requesterAccessController","type":"address"},{"internalType":"uint8","name":"decimals_","type":"uint8"},{"internalType":"string","name":"description_","type":"string"}],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"int256","name":"current","type":"int256"},{"indexed":true,"internalType":"uint256","name":"roundId","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"updatedAt","type":"uint256"}],"name":"AnswerUpdated","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"contractAccessControllerInterface","name":"old","type":"address"},{"indexed":false,"internalType":"contractAccessControllerInterface","name":"current","type":"address"}],"name":"BillingAccessControllerSet","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint32","name":"maximumGasPriceGwei","type":"uint32"},{"indexed":false,"internalType":"uint32","name":"reasonableGasPriceGwei","type":"uint32"},{"indexed":false,"internalType":"uint32","name":"observationPaymentGjuels","type":"uint32"},{"indexed":false,"internalType":"uint32","name":"transmissionPaymentGjuels","type":"uint32"},{"indexed":false,"internalType":"uint24","name":"accountingGas","type":"uint24"}],"name":"BillingSet","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint32","name":"previousConfigBlockNumber","type":"uint32"},{"indexed":false,"internalType":"bytes32","name":"configDigest","type":"bytes32"},{"indexed":false,"internalType":"uint64","name":"configCount","type":"uint64"},{"indexed":false,"internalType":"address[]","name":"signers","type":"address[]"},{"indexed":false,"internalType":"address[]","name":"transmitters","type":"address[]"},{"indexed":false,"internalType":"uint8","name":"f","type":"uint8"},{"indexed":false,"internalType":"bytes","name":"onchainConfig","type":"bytes"},{"indexed":false,"internalType":"uint64","name":"offchainConfigVersion","type":"uint64"},{"indexed":false,"internalType":"bytes","name":"offchainConfig","type":"bytes"}],"name":"ConfigSet","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"contractLinkTokenInterface","name":"oldLinkToken","type":"address"},{"indexed":true,"internalType":"contractLinkTokenInterface","name":"newLinkToken","type":"address"}],"name":"LinkTokenSet","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint256","name":"roundId","type":"uint256"},{"indexed":true,"internalType":"address","name":"startedBy","type":"address"},{"indexed":false,"internalType":"uint256","name":"startedAt","type":"uint256"}],"name":"NewRound","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint32","name":"aggregatorRoundId","type":"uint32"},{"indexed":false,"internalType":"int192","name":"answer","type":"int192"},{"indexed":false,"internalType":"address","name":"transmitter","type":"address"},{"indexed":false,"internalType":"uint32","name":"observationsTimestamp","type":"uint32"},{"indexed":false,"internalType":"int192[]","name":"observations","type":"int192[]"},{"indexed":false,"internalType":"bytes","name":"observers","type":"bytes"},{"indexed":false,"internalType":"int192","name":"juelsPerFeeCoin","type":"int192"},{"indexed":false,"internalType":"bytes32","name":"configDigest","type":"bytes32"},{"indexed":false,"internalType":"uint40","name":"epochAndRound","type":"uint40"}],"name":"NewTransmission","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"transmitter","type":"address"},{"indexed":true,"internalType":"address","name":"payee","type":"address"},{"indexed":false,"internalType":"uint256","name":"amount","type":"uint256"},{"indexed":true,"internalType":"contractLinkTokenInterface","name":"linkToken","type":"address"}],"name":"OraclePaid","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"}],"name":"OwnershipTransferRequested","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"}],"name":"OwnershipTransferred","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"transmitter","type":"address"},{"indexed":true,"internalType":"address","name":"current","type":"address"},{"indexed":true,"internalType":"address","name":"proposed","type":"address"}],"name":"PayeeshipTransferRequested","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"transmitter","type":"address"},{"indexed":true,"internalType":"address","name":"previous","type":"address"},{"indexed":true,"internalType":"address","name":"current","type":"address"}],"name":"PayeeshipTransferred","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"contractAccessControllerInterface","name":"old","type":"address"},{"indexed":false,"internalType":"contractAccessControllerInterface","name":"current","type":"address"}],"name":"RequesterAccessControllerSet","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"requester","type":"address"},{"indexed":false,"internalType":"bytes32","name":"configDigest","type":"bytes32"},{"indexed":false,"internalType":"uint32","name":"epoch","type":"uint32"},{"indexed":false,"internalType":"uint8","name":"round","type":"uint8"}],"name":"RoundRequested","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"bytes32","name":"configDigest","type":"bytes32"},{"indexed":false,"internalType":"uint32","name":"epoch","type":"uint32"}],"name":"Transmitted","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"contractAggregatorValidatorInterface","name":"previousValidator","type":"address"},{"indexed":false,"internalType":"uint32","name":"previousGasLimit","type":"uint32"},{"indexed":true,"internalType":"contractAggregatorValidatorInterface","name":"currentValidator","type":"address"},{"indexed":false,"internalType":"uint32","name":"currentGasLimit","type":"uint32"}],"name":"ValidatorConfigSet","type":"event"},{"inputs":[],"name":"acceptOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"transmitter","type":"address"}],"name":"acceptPayeeship","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"decimals","outputs":[{"internalType":"uint8","name":"","type":"uint8"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"description","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"roundId","type":"uint256"}],"name":"getAnswer","outputs":[{"internalType":"int256","name":"","type":"int256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getBilling","outputs":[{"internalType":"uint32","name":"maximumGasPriceGwei","type":"uint32"},{"internalType":"uint32","name":"reasonableGasPriceGwei","type":"uint32"},{"internalType":"uint32","name":"observationPaymentGjuels","type":"uint32"},{"internalType":"uint32","name":"transmissionPaymentGjuels","type":"uint32"},{"internalType":"uint24","name":"accountingGas","type":"uint24"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getBillingAccessController","outputs":[{"internalType":"contractAccessControllerInterface","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getLinkToken","outputs":[{"internalType":"contractLinkTokenInterface","name":"linkToken","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getRequesterAccessController","outputs":[{"internalType":"contractAccessControllerInterface","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint80","name":"roundId","type":"uint80"}],"name":"getRoundData","outputs":[{"internalType":"uint80","name":"roundId_","type":"uint80"},{"internalType":"int256","name":"answer","type":"int256"},{"internalType":"uint256","name":"startedAt","type":"uint256"},{"internalType":"uint256","name":"updatedAt","type":"uint256"},{"internalType":"uint80","name":"answeredInRound","type":"uint80"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"roundId","type":"uint256"}],"name":"getTimestamp","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getTransmitters","outputs":[{"internalType":"address[]","name":"","type":"address[]"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getValidatorConfig","outputs":[{"internalType":"contractAggregatorValidatorInterface","name":"validator","type":"address"},{"internalType":"uint32","name":"gasLimit","type":"uint32"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"latestAnswer","outputs":[{"internalType":"int256","name":"","type":"int256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"latestConfigDetails","outputs":[{"internalType":"uint32","name":"configCount","type":"uint32"},{"internalType":"uint32","name":"blockNumber","type":"uint32"},{"internalType":"bytes32","name":"configDigest","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"latestConfigDigestAndEpoch","outputs":[{"internalType":"bool","name":"scanLogs","type":"bool"},{"internalType":"bytes32","name":"configDigest","type":"bytes32"},{"internalType":"uint32","name":"epoch","type":"uint32"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"latestRound","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"latestRoundData","outputs":[{"internalType":"uint80","name":"roundId","type":"uint80"},{"internalType":"int256","name":"answer","type":"int256"},{"internalType":"uint256","name":"startedAt","type":"uint256"},{"internalType":"uint256","name":"updatedAt","type":"uint256"},{"internalType":"uint80","name":"answeredInRound","type":"uint80"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"latestTimestamp","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"latestTransmissionDetails","outputs":[{"internalType":"bytes32","name":"configDigest","type":"bytes32"},{"internalType":"uint32","name":"epoch","type":"uint32"},{"internalType":"uint8","name":"round","type":"uint8"},{"internalType":"int192","name":"latestAnswer_","type":"int192"},{"internalType":"uint64","name":"latestTimestamp_","type":"uint64"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"linkAvailableForPayment","outputs":[{"internalType":"int256","name":"availableBalance","type":"int256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"maxAnswer","outputs":[{"internalType":"int192","name":"","type":"int192"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"minAnswer","outputs":[{"internalType":"int192","name":"","type":"int192"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"transmitterAddress","type":"address"}],"name":"oracleObservationCount","outputs":[{"internalType":"uint32","name":"","type":"uint32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"transmitterAddress","type":"address"}],"name":"owedPayment","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"requestNewRound","outputs":[{"internalType":"uint80","name":"","type":"uint80"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint32","name":"maximumGasPriceGwei","type":"uint32"},{"internalType":"uint32","name":"reasonableGasPriceGwei","type":"uint32"},{"internalType":"uint32","name":"observationPaymentGjuels","type":"uint32"},{"internalType":"uint32","name":"transmissionPaymentGjuels","type":"uint32"},{"internalType":"uint24","name":"accountingGas","type":"uint24"}],"name":"setBilling","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"contractAccessControllerInterface","name":"_billingAccessController","type":"address"}],"name":"setBillingAccessController","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address[]","name":"signers","type":"address[]"},{"internalType":"address[]","name":"transmitters","type":"address[]"},{"internalType":"uint8","name":"f","type":"uint8"},{"internalType":"bytes","name":"onchainConfig","type":"bytes"},{"internalType":"uint64","name":"offchainConfigVersion","type":"uint64"},{"internalType":"bytes","name":"offchainConfig","type":"bytes"}],"name":"setConfig","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"contractLinkTokenInterface","name":"linkToken","type":"address"},{"internalType":"address","name":"recipient","type":"address"}],"name":"setLinkToken","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address[]","name":"transmitters","type":"address[]"},{"internalType":"address[]","name":"payees","type":"address[]"}],"name":"setPayees","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"contractAccessControllerInterface","name":"requesterAccessController","type":"address"}],"name":"setRequesterAccessController","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"contractAggregatorValidatorInterface","name":"newValidator","type":"address"},{"internalType":"uint32","name":"newGasLimit","type":"uint32"}],"name":"setValidatorConfig","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"to","type":"address"}],"name":"transferOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"transmitter","type":"address"},{"internalType":"address","name":"proposed","type":"address"}],"name":"transferPayeeship","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32[3]","name":"reportContext","type":"bytes32[3]"},{"internalType":"bytes","name":"report","type":"bytes"},{"internalType":"bytes32[]","name":"rs","type":"bytes32[]"},{"internalType":"bytes32[]","name":"ss","type":"bytes32[]"},{"internalType":"bytes32","name":"rawVs","type":"bytes32"}],"name":"transmit","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"typeAndVersion","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"pure","type":"function"},{"inputs":[],"name":"version","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"recipient","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"withdrawFunds","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"transmitter","type":"address"}],"name":"withdrawPayment","outputs":[],"stateMutability":"nonpayable","type":"function"}]`,
						ChainReaderDefinitions: map[string]evmtypes.ChainReaderDefinition{
							"LatestTransmissionDetails": {
								ChainSpecificName: "latestTransmissionDetails",
								OutputModifications: codec.ModifiersConfig{
									&codec.RenameModifierConfig{
										Fields: map[string]string{
											"LatestAnswer_":    "LatestAnswer",
											"LatestTimestamp_": "LatestTimestamp",
										},
									},
								},
							},
							"LatestRoundRequested": {
								ChainSpecificName: "RoundRequested",
								ReadType:          1,
							},
						},
					},
				},
			},
			"codec": evmtypes.CodecConfig{
				ChainCodecConfigs: map[string]evmtypes.ChainCodecConfig{
					"MedianReport": {
						TypeABI: `[
  {
    "Name": "Timestamp",
    "Type": "uint32"
  },
  {
    "Name": "Observers",
    "Type": "bytes32"
  },
  {
    "Name": "Observations",
    "Type": "int192[]"
  },
  {
    "Name": "JuelsPerFeeCoin",
    "Type": "int192"
  }
]
`,
					},
				},
			},
		},
		PluginConfig: map[string]interface{}{"juelsPerFeeCoinSource": `		// data source 1
		ds1          [type=bridge name="%s"];
		ds1_parse    [type=jsonparse path="data"];
		ds1_multiply [type=multiply times=2];

		// data source 2
		ds2          [type=http method=GET url="%s"];
		ds2_parse    [type=jsonparse path="data"];
		ds2_multiply [type=multiply times=2];

		ds1 -> ds1_parse -> ds1_multiply -> answer1;
		ds2 -> ds2_parse -> ds2_multiply -> answer1;

	answer1 [type=median index=0];
`,
		},
	}

	t.Run("marshal", func(t *testing.T) {
		gotB, err := toml.Marshal(val)
		require.NoError(t, err)
		require.Equal(t, compact, string(gotB))
	})

	t.Run("round-trip", func(t *testing.T) {
		var gotVal OCR2OracleSpec
		require.NoError(t, toml.Unmarshal([]byte(compact), &gotVal))
		gotB, err := toml.Marshal(gotVal)
		require.NoError(t, err)
		require.Equal(t, compact, string(gotB))
		t.Run("pretty", func(t *testing.T) {
			var gotVal OCR2OracleSpec
			require.NoError(t, toml.Unmarshal([]byte(pretty), &gotVal))
			gotB, err := toml.Marshal(gotVal)
			require.NoError(t, err)
			t.Log("marshaled compact:", string(gotB))
			require.Equal(t, compact, string(gotB))
		})
	})
}
