package test

import (
	"math/big"
	"time"

	"github.com/google/uuid"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

const ConfigTOML = `[Foo]
Bar = "Baz"
`

const spec = `
answer [type=sum values=<[ $(val), 2 ]>]
answer;
`

const (
	account          = libocr.Account("testaccount")
	balanceCheck     = true
	blockHeight      = uint64(1337)
	changedInBlock   = uint64(14)
	total            = 2
	epoch            = uint32(88)
	errMsg           = "test error"
	from             = "0xabcd"
	limit            = 42
	lookbackDuration = time.Minute + 4*time.Second
	max              = 101
	n                = 12
	offset           = 11
	round            = uint8(74)
	to               = "0x1234"
	network          = "solana"
	contractID       = "contract-id"
	telemType        = "mercury"
)

var (
	amount = big.NewInt(123456789)
	chain  = types.ChainStatus{
		ID:     chainID,
		Config: ConfigTOML,
	}
	chainID            = "chain-id"
	configDigest       = libocr.ConfigDigest([32]byte{2: 10, 12: 16})
	configDigestPrefix = libocr.ConfigDigestPrefix(99)
	contractConfig     = libocr.ContractConfig{
		ConfigDigest:          configDigest,
		ConfigCount:           42,
		Signers:               []libocr.OnchainPublicKey{[]byte{15: 1}},
		Transmitters:          []libocr.Account{"foo", "bar"},
		F:                     11,
		OnchainConfig:         []byte{2: 11, 14: 22, 31: 1},
		OffchainConfigVersion: 2,
		OffchainConfig:        []byte{1: 99, 12: 55},
	}
	encoded         = []byte{5: 11}
	juelsPerFeeCoin = big.NewInt(1234)
	onchainConfig   = median.OnchainConfig{Min: big.NewInt(-12), Max: big.NewInt(1234567890987654321)}
	latestAnswer    = big.NewInt(-66)
	latestTimestamp = time.Unix(1234567890, 987654321).UTC()
	medianValue     = big.NewInt(-1042)
	nodes           = []types.NodeStatus{{
		ChainID: "foo",
		State:   "Alive",
		Config: `Name = 'bar'
URL = 'http://example.com'
`}, {
		ChainID: "foo",
		State:   "Alive",
		Config: `Name = 'baz'
URL = 'https://test.url'
`}}
	observation = libocr.Observation([]byte{21: 19})
	obs         = []libocr.AttributedObservation{{Observation: []byte{21: 19}, Observer: commontypes.OracleID(99)}}
	PluginArgs  = types.PluginArgs{
		TransmitterID: "testtransmitter",
		PluginConfig:  []byte{100: 88},
	}
	pobs      = []median.ParsedAttributedObservation{{Timestamp: 123, Value: big.NewInt(31), JuelsPerFeeCoin: big.NewInt(54), Observer: commontypes.OracleID(99)}}
	query     = []byte{42: 42}
	RelayArgs = types.RelayArgs{
		ExternalJobID: uuid.MustParse("1051429b-aa66-11ed-b0d2-5cff35dfbe67"),
		JobID:         123,
		ContractID:    "testcontract",
		New:           true,
		RelayConfig:   []byte{42: 11},
		ProviderType:  string(types.Median),
	}
	report        = libocr.Report{42: 101}
	reportContext = libocr.ReportContext{
		ReportTimestamp: libocr.ReportTimestamp{
			ConfigDigest: configDigest,
			Epoch:        epoch,
			Round:        round,
		},
		ExtraHash: [32]byte{1: 2, 3: 4, 5: 6},
	}
	reportingPluginConfig = libocr.ReportingPluginConfig{
		ConfigDigest:                            configDigest,
		OracleID:                                commontypes.OracleID(10),
		N:                                       12,
		F:                                       42,
		OnchainConfig:                           []byte{17: 11},
		OffchainConfig:                          []byte{32: 64},
		EstimatedRoundInterval:                  time.Second,
		MaxDurationQuery:                        time.Hour,
		MaxDurationObservation:                  time.Millisecond,
		MaxDurationReport:                       time.Microsecond,
		MaxDurationShouldAcceptFinalizedReport:  10 * time.Second,
		MaxDurationShouldTransmitAcceptedReport: time.Minute,
	}
	rpi = libocr.ReportingPluginInfo{
		Name:          "test",
		UniqueReports: true,
		Limits: libocr.ReportingPluginLimits{
			MaxQueryLength:       42,
			MaxObservationLength: 13,
			MaxReportLength:      17,
		},
	}
	shouldAccept   = true
	shouldReport   = true
	shouldTransmit = true
	signed         = []byte{13: 37}
	sigs           = []libocr.AttributedOnchainSignature{{Signature: []byte{9: 8, 7: 6}, Signer: commontypes.OracleID(54)}}
	value          = big.NewInt(999)
	vars           = types.Vars{
		Vars: map[string]interface{}{"foo": "baz"},
	}
	options = types.Options{
		MaxTaskDuration: 10 * time.Second,
	}
	taskResults = types.TaskResults([]types.TaskResult{
		{
			TaskValue: types.TaskValue{
				Value: "hello",
			},
			Index: 0,
		},
	})
	payload                     = []byte("oops")
	medianContractGenericMethod = "LatestTransmissionDetails"
	getLatestValueParams        = map[string]string{"param1": "value1", "param2": "value2"}
	boundContract               = types.BoundContract{
		Name:    "my median contract",
		Address: "0xBbf078A8849D74653e36E6DBBdC7e1a35E657C26",
		Pending: false,
	}
	latestValue = map[string]int{"ret1": 1, "ret2": 2}
)
