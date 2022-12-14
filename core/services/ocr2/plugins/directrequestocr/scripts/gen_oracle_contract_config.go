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
		"f19f36c376839c74c637efb04666db5735d2f83e5c78ecac5870f4518cfe1a1c", // replace this
		"73e72ebb7a5fe651bdacb8b94100d426ca77ffb7be7b1b8680dcdfbf78fada56",
		"ec3d11a5497cd58de6d8833721b6c2dbf571e488cac6b8707efbbd652b716431",
		"7985152a2865e0b0f777dbe871c4409c014d3717c24c0c208cf2805cf7dfa870"}

	configPublicKeys := []string{
		"9e19519f1372f99118cd794823e4bdbe108571eb9aaca4f2911a6cdf96759d0c", // replace this
		"4abb228fc3a537c5f5130283f58e2f7aadc4b767baa261d7126abd53299c5c2e",
		"263624fee591e6ca5075982e2bfcc6782d94331f4452fae937f73a8ae5125a57",
		"1760d2d1f1fded9440f38bdd416d2eef9e791f675ac004c4d459f7fd327ea244"}

	onchainPubKeys := []common.Address{
		common.HexToAddress("12024999d75b69e7edc04383889ab2ffbdc409e9"), // replace this (on-chain pub key)
		common.HexToAddress("a4ad9d5a6c6171ebbceb4eaa149d6ac0cd85b3ad"),
		common.HexToAddress("219d033aeb68213fc98ff28d613a48d10c337b9b"),
		common.HexToAddress("ad095d87f24a69d01f45681c1c6868c80fec5703")}

	peerIDs := []string{
		"12D3KooWM8MSbL9cD9JiUT5HfiBionGEuEnTBBL6MyT8BhR5XmQz", // replace this
		"12D3KooW9rXtuvYx2RgXeBkAYMVAAcQ8LG5eFCi6msZCCuYt1FQa",
		"12D3KooWKUDHb4C7iGFbwzuRtTqZnTpwoRwTM7hKz68gS28qvAV2",
		"12D3KooWQymo2wcBMiWVVNaTMLNhozACrxt3wbMKDg6vjKzVz6Kp",
	}

	transmiterAddresses := []string{
		"0xD3c1b2710332FDA16001433c6198AF7e55EeD627", // replace this
		"0xf318D60885C95B3262fa70fbc6e912282Eb2925b",
		"0x5d7589b2D20c4532F94eb30a713b571F0bc82477",
		"0x12BeF77C49792935d46A42F6eA4C2a582Da627C2",
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
