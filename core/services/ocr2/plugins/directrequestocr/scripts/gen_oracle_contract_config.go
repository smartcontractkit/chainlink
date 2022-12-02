package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

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

func addrArrayToPrettyString(addrs []string) string {
	var quoted []string
	for _, addr := range addrs {
		quoted = append(quoted, fmt.Sprintf("\"%s\"", addr))
	}
	return "[" + strings.Join(quoted, ",") + "]"
}

func main() {
	// NOTE: replace values below with actual keys fetched from DON members
	offchainPubKeysHex := []string{
		"7ec48c19696fa23bc5706921e21e1a51423f24ddd0a9d0cce127d5ca0f76cb01", // replace this
		"7ec48c19696fa23bc5706921e21e1a51423f24ddd0a9d0cce127d5ca0f76cb01",
		"7ec48c19696fa23bc5706921e21e1a51423f24ddd0a9d0cce127d5ca0f76cb01",
		"7ec48c19696fa23bc5706921e21e1a51423f24ddd0a9d0cce127d5ca0f76cb01"}

	configPublicKeys := []string{
		"1cde63dc6b44dc44562713ce9735e68ffc47884cca966be168e318a6267b2b3a", // replace this
		"1cde63dc6b44dc44562713ce9735e68ffc47884cca966be168e318a6267b2b3a",
		"1cde63dc6b44dc44562713ce9735e68ffc47884cca966be168e318a6267b2b3a",
		"1cde63dc6b44dc44562713ce9735e68ffc47884cca966be168e318a6267b2b3a"}

	onchainPubKeys := []common.Address{
		common.HexToAddress("659da7715b116be562ae1418ca578e1e5280245b"), // replace this (on-chain pub key)
		common.HexToAddress("659da7715b116be562ae1418ca578e1e5280245b"),
		common.HexToAddress("659da7715b116be562ae1418ca578e1e5280245b"),
		common.HexToAddress("659da7715b116be562ae1418ca578e1e5280245b")}

	peerIDs := []string{
		"12D3KooWS3gci8DQvGTXQg2376YhFQsP4qb1oUbqQFyExvSJxU4h", // replace this
		"12D3KooWS3gci8DQvGTXQg2376YhFQsP4qb1oUbqQFyExvSJxU4h",
		"12D3KooWS3gci8DQvGTXQg2376YhFQsP4qb1oUbqQFyExvSJxU4h",
		"12D3KooWS3gci8DQvGTXQg2376YhFQsP4qb1oUbqQFyExvSJxU4h",
	}

	transmiterAddresses := []string{
		"0xDA43dABB62fAE4906CCaeA0B9c47BA4D32b887Ac", // replace this
		"0xDA43dABB62fAE4906CCaeA0B9c47BA4D32b887Ac",
		"0xDA43dABB62fAE4906CCaeA0B9c47BA4D32b887Ac",
		"0xDA43dABB62fAE4906CCaeA0B9c47BA4D32b887Ac",
	}

	offchainPubKeysBytes := []types.OffchainPublicKey{}
	for _, pkHex := range offchainPubKeysHex {
		pkBytes, err := hex.DecodeString(pkHex)
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
	for _, pkHex := range configPublicKeys {
		pkBytes, err := hex.DecodeString(pkHex)
		if err != nil {
			panic(err)
		}

		pkBytesFixed := [ed25519.PublicKeySize]byte{}
		n := copy(pkBytesFixed[:], pkBytes)
		if n != ed25519.PublicKeySize {
			panic("wrong num elements copied")
		}

		configPubKeysBytes = append(configPubKeysBytes, types.ConfigEncryptionPublicKey(pkBytesFixed))
	}

	S := []int{}
	identities := []confighelper.OracleIdentityExtra{}

	for index := range configPublicKeys {
		S = append(S, 1)
		identities = append(identities, confighelper.OracleIdentityExtra{
			OracleIdentity: confighelper.OracleIdentity{
				OnchainPublicKey:  onchainPubKeys[index][:],
				OffchainPublicKey: offchainPubKeysBytes[index],
				PeerID:            peerIDs[index],
				TransmitAccount:   types.Account(transmiterAddresses[index]),
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

	signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, err := confighelper.ContractSetConfigArgsForTests(
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
	if err != nil {
		panic(err)
	}

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

	file, err := json.MarshalIndent(config, "", " ")

	if err != nil {
		panic(err)
	}

	filename := flag.String("o", "OCR2DROracleConfig.json", "output file")
	flag.Parse()

	err = os.WriteFile(*filename, file, 0600)

	if err != nil {
		panic(err)
	}

	fmt.Println("signers: ", addrArrayToPrettyString(signersStr))
	fmt.Println("transmitters: ", addrArrayToPrettyString(transmittersStr))
	fmt.Println("f:", f)
	fmt.Println("onchainConfig:", onchainConfig)
	fmt.Println("offchainConfigVersion:", offchainConfigVersion)
	fmt.Printf("offchainConfig: 0x%s\n", hex.EncodeToString(offchainConfig))
}
