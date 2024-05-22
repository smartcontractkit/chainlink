package median_test

import (
	"math/big"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	libocr "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	errorlogtest "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/core/services/errorlog/test"
	chainreadertest "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/chainreader/test"
	ocr2test "github.com/smartcontractkit/chainlink-common/pkg/loop/internal/relayer/pluginprovider/ocr2/test"
)

const ConfigTOML = `[Foo]
Bar = "Baz"
`

const (
	lookbackDuration = time.Minute + 4*time.Second
	max              = 101
	n                = 12
)

var (
	MedianFactoryServer = staticMedianFactoryServer{
		staticPluginMedianConfig: staticPluginMedianConfig{
			provider:                   MedianProvider,
			dataSource:                 DataSource,
			juelsPerFeeCoinDataSource:  JuelsPerFeeCoinDataSource,
			gasPriceSubunitsDataSource: GasPriceSubunitsDataSource,
			errorLog:                   errorlogtest.ErrorLog,
		},
	}

	MedianProvider = staticMedianProvider{
		staticMedianProviderConfig: staticMedianProviderConfig{
			offchainDigester:    ocr2test.OffchainConfigDigester,
			contractTracker:     ocr2test.ContractConfigTracker,
			contractTransmitter: ocr2test.ContractTransmitter,
			reportCodec:         staticReportCodec{},
			medianContract: staticMedianContract{
				staticMedianContractConfig: staticMedianContractConfig{
					configDigest:     libocr.ConfigDigest([32]byte{1: 1, 11: 8}),
					epoch:            7,
					round:            11,
					latestAnswer:     big.NewInt(123),
					latestTimestamp:  time.Unix(1234567890, 987654321).UTC(),
					lookbackDuration: lookbackDuration,
				},
			},
			onchainConfigCodec: staticOnchainConfigCodec{},
			chainReader:        chainreadertest.ChainReader,
		},
	}
)

var (
	encodedOnchainConfig = []byte{5: 11}
	juelsPerFeeCoin      = big.NewInt(1234)
	gasPriceSubunits     = big.NewInt(777)
	onchainConfig        = median.OnchainConfig{Min: big.NewInt(-12), Max: big.NewInt(1234567890987654321)}
	medianValue          = big.NewInt(-1042)

	pobs = []median.ParsedAttributedObservation{{Timestamp: 123, Value: big.NewInt(31), JuelsPerFeeCoin: big.NewInt(54), GasPriceSubunits: big.NewInt(77), Observer: commontypes.OracleID(99)}}

	report        = libocr.Report{42: 101}
	reportContext = libocr.ReportContext{
		ReportTimestamp: libocr.ReportTimestamp{
			ConfigDigest: libocr.ConfigDigest([32]byte{1: 7, 31: 3}),
			Epoch:        79,
			Round:        17,
		},
		ExtraHash: [32]byte{1: 2, 3: 4, 5: 6},
	}

	reportingPluginConfig = libocr.ReportingPluginConfig{
		ConfigDigest:                            libocr.ConfigDigest{}, //testpluginprovider.ConfigDigest,
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
	value = big.NewInt(999)
)
