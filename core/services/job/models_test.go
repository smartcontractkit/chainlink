package job

import (
	_ "embed"
	"reflect"
	"testing"
	"time"

	"github.com/pelletier/go-toml/v2"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"
	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	evmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

func TestOCR2OracleSpec_RelayIdentifier(t *testing.T) {
	type fields struct {
		Relay       string
		ChainID     string
		RelayConfig JSONConfig
	}
	tests := []struct {
		name    string
		fields  fields
		want    types.RelayID
		wantErr bool
	}{
		{name: "err no chain id",
			fields:  fields{},
			want:    types.RelayID{},
			wantErr: true,
		},
		{
			name: "evm explicitly configured",
			fields: fields{
				Relay:   types.NetworkEVM,
				ChainID: "1",
			},
			want: types.RelayID{Network: types.NetworkEVM, ChainID: "1"},
		},
		{
			name: "evm implicitly configured",
			fields: fields{
				Relay:       types.NetworkEVM,
				RelayConfig: map[string]any{"chainID": 1},
			},
			want: types.RelayID{Network: types.NetworkEVM, ChainID: "1"},
		},
		{
			name: "evm implicitly configured with bad value",
			fields: fields{
				Relay:       types.NetworkEVM,
				RelayConfig: map[string]any{"chainID": float32(1)},
			},
			want:    types.RelayID{},
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
		Relay:                             types.NetworkEVM,
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
				Contracts: map[string]evmtypes.ChainContractReader{
					"median": {
						ContractABI: `[
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": true,
        "internalType": "address",
        "name": "requester",
        "type": "address"
      },
      {
        "indexed": false,
        "internalType": "bytes32",
        "name": "configDigest",
        "type": "bytes32"
      },
      {
        "indexed": false,
        "internalType": "uint32",
        "name": "epoch",
        "type": "uint32"
      },
      {
        "indexed": false,
        "internalType": "uint8",
        "name": "round",
        "type": "uint8"
      }
    ],
    "name": "RoundRequested",
    "type": "event"
  },
  {
    "inputs": [],
    "name": "latestTransmissionDetails",
    "outputs": [
      {
        "internalType": "bytes32",
        "name": "configDigest",
        "type": "bytes32"
      },
      {
        "internalType": "uint32",
        "name": "epoch",
        "type": "uint32"
      },
      {
        "internalType": "uint8",
        "name": "round",
        "type": "uint8"
      },
      {
        "internalType": "int192",
        "name": "latestAnswer_",
        "type": "int192"
      },
      {
        "internalType": "uint64",
        "name": "latestTimestamp_",
        "type": "uint64"
      }
    ],
    "stateMutability": "view",
    "type": "function"
  }
]
`,
						Configs: map[string]*evmtypes.ChainReaderDefinition{
							"LatestTransmissionDetails": {
								ChainSpecificName: "latestTransmissionDetails",
								OutputModifications: codec.ModifiersConfig{
									&codec.EpochToTimeModifierConfig{
										Fields: []string{"LatestTimestamp_"},
									},
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
								ReadType:          evmtypes.Event,
							},
						},
					},
				},
			},
			"codec": evmtypes.CodecConfig{
				Configs: map[string]evmtypes.ChainCodecConfig{
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
		OnchainSigningStrategy: map[string]interface{}{
			"strategyName": "single-chain",
			"config": map[string]interface{}{
				"evm":       "",
				"publicKey": "0xdeadbeef",
			},
		},
		PluginConfig: map[string]interface{}{"juelsPerFeeCoinSource": `  // data source 1
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
		t.Log("marshaled:", string(gotB))
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

func TestWorkflowSpec_Validate(t *testing.T) {
	type fields struct {
		ID            int32
		WorkflowID    string
		Workflow      string
		WorkflowOwner string
		WorkflowName  string
		CreatedAt     time.Time
		UpdatedAt     time.Time
	}
	tests := []struct {
		name        string
		fields      fields
		expectedErr error
	}{
		{
			name: "valid",
			fields: fields{
				WorkflowID:    "15c631d295ef5e32deb99a10ee6804bc4af1385568f9b3363f6552ac6dbb2cef",
				WorkflowOwner: "00000000000000000000000000000000000000aa",
				WorkflowName:  "ten bytes!",
			},
		},
		{
			name: "valid 0x prefix hex owner",
			fields: fields{
				WorkflowID:    "15c631d295ef5e32deb99a10ee6804bc4af1385568f9b3363f6552ac6dbb2cef",
				WorkflowOwner: "0x00000000000000000000000000000000000000ff",
				WorkflowName:  "ten bytes!",
			},
		},
		{
			name: "not hex owner",
			fields: fields{
				WorkflowID:    "15c631d295ef5e32deb99a10ee6804bc4af1385568f9b3363f6552ac6dbb2cef",
				WorkflowOwner: "00000000000000000000000000000000000000az",
				WorkflowName:  "ten bytes!",
			},
			expectedErr: ErrInvalidWorkflowOwner,
		},
		{
			name: "not len 40 owner",
			fields: fields{
				WorkflowID:    "15c631d295ef5e32deb99a10ee6804bc4af1385568f9b3363f6552ac6dbb2cef",
				WorkflowOwner: "0000000000",
				WorkflowName:  "ten bytes!",
			},
			expectedErr: ErrInvalidWorkflowOwner,
		},
		{
			name: "not len 10 name",
			fields: fields{
				WorkflowID:    "15c631d295ef5e32deb99a10ee6804bc4af1385568f9b3363f6552ac6dbb2cef",
				WorkflowOwner: "00000000000000000000000000000000000000aa",
				WorkflowName:  "not ten bytes!",
			},
			expectedErr: ErrInvalidWorkflowName,
		},
		{
			name: "not len 64 id",
			fields: fields{
				WorkflowID:    "15c631d295ef5e32deb99a10ee6804bc4af1385568f9b3363f",
				WorkflowOwner: "00000000000000000000000000000000000000aa",
				WorkflowName:  "ten bytes!",
			},
			expectedErr: ErrInvalidWorkflowID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WorkflowSpec{
				ID:            tt.fields.ID,
				WorkflowID:    tt.fields.WorkflowID,
				Workflow:      tt.fields.Workflow,
				WorkflowOwner: tt.fields.WorkflowOwner,
				WorkflowName:  tt.fields.WorkflowName,
				CreatedAt:     tt.fields.CreatedAt,
				UpdatedAt:     tt.fields.UpdatedAt,
			}
			err := w.Validate()
			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}
