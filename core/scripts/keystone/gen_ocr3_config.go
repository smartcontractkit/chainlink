package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
)

type TopLevelConfigSource struct {
	OracleConfig OracleConfigSource
}

type OracleConfigSource struct {
	MaxQueryLengthBytes       uint32
	MaxObservationLengthBytes uint32
	MaxReportLengthBytes      uint32
	MaxRequestBatchSize       uint32
	UniqueReports             bool

	DeltaProgressMillis               uint32
	DeltaResendMillis                 uint32
	DeltaInitialMillis                uint32
	DeltaRoundMillis                  uint32
	DeltaGraceMillis                  uint32
	DeltaCertifiedCommitRequestMillis uint32
	DeltaStageMillis                  uint32
	MaxRoundsPerEpoch                 uint64
	TransmissionSchedule              []int

	MaxDurationQueryMillis       uint32
	MaxDurationObservationMillis uint32
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

func main() {
	fs := flag.NewFlagSet("config_gen", flag.ExitOnError)
	configFile := fs.String("config", "config_example.json", "a file containing JSON config")
	keysFile := fs.String("keys", "public_keys_example.json", "a file containing node public keys")
	if err := fs.Parse(os.Args[1:]); err != nil || *keysFile == "" || *configFile == "" {
		fs.Usage()
		os.Exit(1)
	}

	topLevelCfg := mustParseJSONConfigFile(*configFile)
	cfg := topLevelCfg.OracleConfig
	nca := mustParseKeysFile(*keysFile)

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

	signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, err := ocr3confighelper.ContractSetConfigArgsForTests(
		time.Duration(cfg.DeltaProgressMillis)*time.Millisecond,
		time.Duration(cfg.DeltaResendMillis)*time.Millisecond,
		time.Duration(cfg.DeltaInitialMillis)*time.Millisecond,
		time.Duration(cfg.DeltaRoundMillis)*time.Millisecond,
		time.Duration(cfg.DeltaGraceMillis)*time.Millisecond,
		time.Duration(cfg.DeltaCertifiedCommitRequestMillis)*time.Millisecond,
		time.Duration(cfg.DeltaStageMillis)*time.Millisecond,
		cfg.MaxRoundsPerEpoch,
		cfg.TransmissionSchedule,
		identities,
		nil, // empty plugin config
		time.Duration(cfg.MaxDurationQueryMillis)*time.Millisecond,
		time.Duration(cfg.MaxDurationObservationMillis)*time.Millisecond,
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
	fmt.Println("Config:", string(js))
}
