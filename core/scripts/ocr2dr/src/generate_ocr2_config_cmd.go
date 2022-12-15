package src

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/directrequestocr/config"
)

type orc2drOracleConfig struct {
	Signers               []string `json:"signers"`
	Transmitters          []string `json:"transmitters"`
	F                     uint8    `json:"f"`
	OnchainConfig         []string `json:"onchainConfig"`
	OffchainConfigVersion uint64   `json:"offchainConfigVersion"`
	OffchainConfig        string   `json:"offchainConfig"`
}

type generateOCR2Config struct {
}

func NewGenerateOCR2ConfigCommand() *generateOCR2Config {
	return &generateOCR2Config{}
}

func (g *generateOCR2Config) Name() string {
	return "generate-ocr2-config"
}

func (g *generateOCR2Config) Run(args []string) {
	fs := flag.NewFlagSet(g.Name(), flag.ContinueOnError)
	nodesFile := fs.String("nodes", "", "a file containing nodes urls, logins and passwords")
	chainID := fs.Int64("chainid", 80001, "chain id")
	err := fs.Parse(args)
	if err != nil || nodesFile == nil || *nodesFile == "" || chainID == nil {
		fs.Usage()
		os.Exit(1)
	}

	nodes := mustReadNodesList(*nodesFile)
	nca := mustFetchNodesKeys(*chainID, nodes)[1:] // ignore boot node

	onchainPubKeys := []common.Address{}
	for _, n := range nca {
		onchainPubKeys = append(onchainPubKeys, common.HexToAddress(n.ocr2OnchainPublicKey))
	}

	offchainPubKeysBytes := []types.OffchainPublicKey{}
	for _, n := range nca {
		pkBytes, err := hex.DecodeString(n.ocr2OffchainPublicKey)
		if err != nil {
			panic(err)
		}

		pkBytesFixed := [ed25519.PublicKeySize]byte{}
		n := copy(pkBytesFixed[:], pkBytes)
		if n != ed25519.PublicKeySize {
			panic("wrong num elements copied")
		}

		offchainPubKeysBytes = append(offchainPubKeysBytes, types.OffchainPublicKey(pkBytesFixed))
	}

	configPubKeysBytes := []types.ConfigEncryptionPublicKey{}
	for _, n := range nca {
		pkBytes, err := hex.DecodeString(n.ocr2ConfigPublicKey)
		helpers.PanicErr(err)

		pkBytesFixed := [ed25519.PublicKeySize]byte{}
		n := copy(pkBytesFixed[:], pkBytes)
		if n != ed25519.PublicKeySize {
			panic("wrong num elements copied")
		}

		configPubKeysBytes = append(configPubKeysBytes, types.ConfigEncryptionPublicKey(pkBytesFixed))
	}

	S := []int{}
	identities := []confighelper.OracleIdentityExtra{}

	for index := range nca {
		S = append(S, 1)
		identities = append(identities, confighelper.OracleIdentityExtra{
			OracleIdentity: confighelper.OracleIdentity{
				OnchainPublicKey:  onchainPubKeys[index][:],
				OffchainPublicKey: offchainPubKeysBytes[index],
				PeerID:            nca[index].p2pPeerID,
				TransmitAccount:   types.Account(nca[index].ethAddress),
			},
			ConfigEncryptionPublicKey: configPubKeysBytes[index],
		})
	}

	reportingPluginConfigBytes, err := config.EncodeReportingPluginConfig(&config.ReportingPluginConfigWrapper{
		Config: &config.ReportingPluginConfig{
			MaxQueryLengthBytes:       10_000,
			MaxObservationLengthBytes: 10_000,
			MaxReportLengthBytes:      10_000,
			MaxRequestBatchSize:       10,
			DefaultAggregationMethod:  config.AggregationMethod_AGGREGATION_MODE,
			UniqueReports:             true,
		},
	})
	if err != nil {
		panic(err)
	}

	signers, transmitters, f, _, offchainConfigVersion, offchainConfig, err := confighelper.ContractSetConfigArgsForTests(
		30*time.Second, // deltaProgress
		10*time.Second, // deltaResend
		10*time.Second, // deltaRound
		2*time.Second,  // deltaGrace
		30*time.Second, // deltaStage
		5,              // RMax (maxRounds)
		S,              // S (schedule of randomized transmission order)
		identities,
		reportingPluginConfigBytes,
		5*time.Second, // maxDurationQuery
		5*time.Second, // maxDurationObservation
		5*time.Second, // maxDurationReport
		5*time.Second, // maxDurationAccept
		5*time.Second, // maxDurationTransmit
		1,             // f (max faulty oracles)
		nil,           // empty onChain config
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
		OnchainConfig:         make([]string, 0),
		OffchainConfigVersion: offchainConfigVersion,
		OffchainConfig:        "0x" + hex.EncodeToString(offchainConfig),
	}

	js, err := json.MarshalIndent(config, "", " ")
	helpers.PanicErr(err)

	filepath := filepath.Join(artefactsDir, ocr2ConfigJson)
	err = os.WriteFile(filepath, js, 0600)
	helpers.PanicErr(err)

	fmt.Println("OCR2 config has been saved to:", filepath)
}
