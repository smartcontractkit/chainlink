package src

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/config"
)

// NOTE: Field names need to match what's in the JSON input file (see sample_config.json)
type TopLevelConfigSource struct {
	OracleConfig OracleConfigSource
}

type ThresholdOffchainConfig struct {
	MaxQueryLengthBytes       uint32
	MaxObservationLengthBytes uint32
	MaxReportLengthBytes      uint32
	RequestCountLimit         uint32
	RequestTotalBytesLimit    uint32
	RequireLocalRequestCheck  bool
	K                         uint32
}

type S4ReportingPluginConfig struct {
	MaxQueryLengthBytes       uint32
	MaxObservationLengthBytes uint32
	MaxReportLengthBytes      uint32
	NSnapshotShards           uint32
	MaxObservationEntries     uint32
	MaxReportEntries          uint32
	MaxDeleteExpiredEntries   uint32
}

type OracleConfigSource struct {
	MaxQueryLengthBytes       uint32
	MaxObservationLengthBytes uint32
	MaxReportLengthBytes      uint32
	MaxRequestBatchSize       uint32
	DefaultAggregationMethod  int32
	UniqueReports             bool
	ThresholdOffchainConfig   ThresholdOffchainConfig
	S4ReportingPluginConfig   S4ReportingPluginConfig
	MaxReportTotalCallbackGas uint32

	DeltaProgressMillis  uint32
	DeltaResendMillis    uint32
	DeltaRoundMillis     uint32
	DeltaGraceMillis     uint32
	DeltaStageMillis     uint32
	MaxRoundsPerEpoch    uint8
	TransmissionSchedule []int

	MaxDurationQueryMillis       uint32
	MaxDurationObservationMillis uint32
	MaxDurationReportMillis      uint32
	MaxDurationAcceptMillis      uint32
	MaxDurationTransmitMillis    uint32

	MaxFaultyOracles int
}

type NodeKeys struct {
	EthAddress            string
	P2PPeerID             string // p2p_<key>
	OCR2BundleID          string // used only in job spec
	OCR2OnchainPublicKey  string // ocr2on_evm_<key>
	OCR2OffchainPublicKey string // ocr2off_evm_<key>
	OCR2ConfigPublicKey   string // ocr2cfg_evm_<key>
	CSAPublicKey          string
}

type orc2drOracleConfig struct {
	Signers               []string `json:"signers"`
	Transmitters          []string `json:"transmitters"`
	F                     uint8    `json:"f"`
	OnchainConfig         string   `json:"onchainConfig"`
	OffchainConfigVersion uint64   `json:"offchainConfigVersion"`
	OffchainConfig        string   `json:"offchainConfig"`
}

type generateOCR2Config struct {
}

func NewGenerateOCR2ConfigCommand() *generateOCR2Config {
	return &generateOCR2Config{}
}

func (g *generateOCR2Config) Name() string {
	return "generate-ocr2config"
}

func mustParseJSONConfigFile(fileName string) (output TopLevelConfigSource) {
	return mustParseJSON[TopLevelConfigSource](fileName)
}

func mustParseKeysFile(fileName string) (output []NodeKeys) {
	return mustParseJSON[[]NodeKeys](fileName)
}

func mustParseJSON[T any](fileName string) (output T) {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	bytes, err := io.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bytes, &output)
	if err != nil {
		panic(err)
	}
	return
}

func (g *generateOCR2Config) Run(args []string) {
	fs := flag.NewFlagSet(g.Name(), flag.ContinueOnError)
	nodesFile := fs.String("nodes", "", "a file containing nodes urls, logins and passwords")
	keysFile := fs.String("keys", "", "a file containing nodes public keys")
	configFile := fs.String("config", "", "a file containing JSON config")
	chainID := fs.Int64("chainid", 80001, "chain id")
	if err := fs.Parse(args); err != nil || (*nodesFile == "" && *keysFile == "") || *configFile == "" || chainID == nil {
		fs.Usage()
		os.Exit(1)
	}

	topLevelCfg := mustParseJSONConfigFile(*configFile)
	cfg := topLevelCfg.OracleConfig
	var nca []NodeKeys
	if *keysFile != "" {
		nca = mustParseKeysFile(*keysFile)
	} else {
		nodes := mustReadNodesList(*nodesFile)
		nca = mustFetchNodesKeys(*chainID, nodes)[1:] // ignore boot node

		nodePublicKeys, err := json.MarshalIndent(nca, "", " ")
		if err != nil {
			panic(err)
		}
		filepath := filepath.Join(artefactsDir, ocr2PublicKeysJSON)
		err = os.WriteFile(filepath, nodePublicKeys, 0600)
		if err != nil {
			panic(err)
		}
		fmt.Println("Functions OCR2 public keys have been saved to:", filepath)
	}

	onchainPubKeys := []common.Address{}
	for _, n := range nca {
		onchainPubKeys = append(onchainPubKeys, common.HexToAddress(n.OCR2OnchainPublicKey))
	}

	offchainPubKeysBytes := []types.OffchainPublicKey{}
	for _, n := range nca {
		pkBytes, err := hex.DecodeString(n.OCR2OffchainPublicKey)
		if err != nil {
			panic(err)
		}

		pkBytesFixed := [ed25519.PublicKeySize]byte{}
		nCopied := copy(pkBytesFixed[:], pkBytes)
		if nCopied != ed25519.PublicKeySize {
			panic("wrong num elements copied from ocr2 offchain public key")
		}

		offchainPubKeysBytes = append(offchainPubKeysBytes, types.OffchainPublicKey(pkBytesFixed))
	}

	configPubKeysBytes := []types.ConfigEncryptionPublicKey{}
	for _, n := range nca {
		pkBytes, err := hex.DecodeString(n.OCR2ConfigPublicKey)
		helpers.PanicErr(err)

		pkBytesFixed := [ed25519.PublicKeySize]byte{}
		n := copy(pkBytesFixed[:], pkBytes)
		if n != ed25519.PublicKeySize {
			panic("wrong num elements copied")
		}

		configPubKeysBytes = append(configPubKeysBytes, types.ConfigEncryptionPublicKey(pkBytesFixed))
	}

	identities := []confighelper.OracleIdentityExtra{}
	for index := range nca {
		identities = append(identities, confighelper.OracleIdentityExtra{
			OracleIdentity: confighelper.OracleIdentity{
				OnchainPublicKey:  onchainPubKeys[index][:],
				OffchainPublicKey: offchainPubKeysBytes[index],
				PeerID:            nca[index].P2PPeerID,
				TransmitAccount:   types.Account(nca[index].EthAddress),
			},
			ConfigEncryptionPublicKey: configPubKeysBytes[index],
		})
	}

	reportingPluginConfigBytes, err := config.EncodeReportingPluginConfig(&config.ReportingPluginConfigWrapper{
		Config: &config.ReportingPluginConfig{
			MaxQueryLengthBytes:       cfg.MaxQueryLengthBytes,
			MaxObservationLengthBytes: cfg.MaxObservationLengthBytes,
			MaxReportLengthBytes:      cfg.MaxReportLengthBytes,
			MaxRequestBatchSize:       cfg.MaxRequestBatchSize,
			DefaultAggregationMethod:  config.AggregationMethod(cfg.DefaultAggregationMethod),
			UniqueReports:             cfg.UniqueReports,
			ThresholdPluginConfig: &config.ThresholdReportingPluginConfig{
				MaxQueryLengthBytes:       cfg.ThresholdOffchainConfig.MaxQueryLengthBytes,
				MaxObservationLengthBytes: cfg.ThresholdOffchainConfig.MaxObservationLengthBytes,
				MaxReportLengthBytes:      cfg.ThresholdOffchainConfig.MaxReportLengthBytes,
				RequestCountLimit:         cfg.ThresholdOffchainConfig.RequestCountLimit,
				RequestTotalBytesLimit:    cfg.ThresholdOffchainConfig.RequestTotalBytesLimit,
				RequireLocalRequestCheck:  cfg.ThresholdOffchainConfig.RequireLocalRequestCheck,
				K:                         cfg.ThresholdOffchainConfig.K,
			},
			S4PluginConfig: &config.S4ReportingPluginConfig{
				MaxQueryLengthBytes:       cfg.S4ReportingPluginConfig.MaxQueryLengthBytes,
				MaxObservationLengthBytes: cfg.S4ReportingPluginConfig.MaxObservationLengthBytes,
				MaxReportLengthBytes:      cfg.S4ReportingPluginConfig.MaxReportLengthBytes,
				NSnapshotShards:           cfg.S4ReportingPluginConfig.NSnapshotShards,
				MaxObservationEntries:     cfg.S4ReportingPluginConfig.MaxObservationEntries,
				MaxReportEntries:          cfg.S4ReportingPluginConfig.MaxReportEntries,
				MaxDeleteExpiredEntries:   cfg.S4ReportingPluginConfig.MaxDeleteExpiredEntries,
			},
			MaxReportTotalCallbackGas: cfg.MaxReportTotalCallbackGas,
		},
	})
	if err != nil {
		panic(err)
	}

	signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, err := confighelper.ContractSetConfigArgsForTests(
		time.Duration(cfg.DeltaProgressMillis)*time.Millisecond,
		time.Duration(cfg.DeltaResendMillis)*time.Millisecond,
		time.Duration(cfg.DeltaRoundMillis)*time.Millisecond,
		time.Duration(cfg.DeltaGraceMillis)*time.Millisecond,
		time.Duration(cfg.DeltaStageMillis)*time.Millisecond,
		cfg.MaxRoundsPerEpoch,
		cfg.TransmissionSchedule,
		identities,
		reportingPluginConfigBytes,
		time.Duration(cfg.MaxDurationQueryMillis)*time.Millisecond,
		time.Duration(cfg.MaxDurationObservationMillis)*time.Millisecond,
		time.Duration(cfg.MaxDurationReportMillis)*time.Millisecond,
		time.Duration(cfg.MaxDurationAcceptMillis)*time.Millisecond,
		time.Duration(cfg.MaxDurationTransmitMillis)*time.Millisecond,
		cfg.MaxFaultyOracles,
		nil, // empty onChain config
	)
	helpers.PanicErr(err)

	var signersStr []string
	var transmittersStr []string
	for i := range transmitters {
		signersStr = append(signersStr, "0x"+hex.EncodeToString(signers[i]))
		transmittersStr = append(transmittersStr, string(transmitters[i]))
	}

	config := orc2drOracleConfig{
		Signers:               signersStr,
		Transmitters:          transmittersStr,
		F:                     f,
		OnchainConfig:         "0x" + hex.EncodeToString(onchainConfig),
		OffchainConfigVersion: offchainConfigVersion,
		OffchainConfig:        "0x" + hex.EncodeToString(offchainConfig),
	}

	js, err := json.MarshalIndent(config, "", " ")
	helpers.PanicErr(err)

	filepath := filepath.Join(artefactsDir, ocr2ConfigJson)
	err = os.WriteFile(filepath, js, 0600)
	helpers.PanicErr(err)

	fmt.Println("Functions OCR2 config has been saved to:", filepath)
}
