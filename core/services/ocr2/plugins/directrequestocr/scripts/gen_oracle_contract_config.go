package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/directrequestocr/config"
)

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
		"7ec48c19696fa23bc5706921e21e1a51423f24ddd0a9d0cce127d5ca0f76cb01",
		"8f4a26283996a0ff63c35e5cd34966682dbeec8817142684fe0174064fb02564",
		"401a1152c8b1f281fa3acf7a3fd2f382ae9288f51261cf80713f367bf0ed689a",
		"0f974ef13b30082565085824422f71e47aaa3a26dcd1c5fb6b1e8212cc2e91a6"}

	configPublicKeys := []string{
		"1cde63dc6b44dc44562713ce9735e68ffc47884cca966be168e318a6267b2b3a",
		"a178e98cb41454d84866a61568da410b8d996399226eb8ab2fb8ded8bad16200",
		"88947bfd1cee473c737039a6c6d0f2c5f9f33196f856e59f5f801f68f7a7db20",
		"9924feb7825163ab833ba8339ca3593cd1f63fbbd3f390a8ebbc85f8fb76191b"}

	onchainPubKeys := []common.Address{
		common.HexToAddress("659da7715b116be562ae1418ca578e1e5280245b"),
		common.HexToAddress("2723f5f184cc689a95f9457c5f500a00fe9a6304"),
		common.HexToAddress("5ca88c71e64522385b4b189bb09ef1c3db7e48fa"),
		common.HexToAddress("ebe1e7d8525911700c5543331b56118813864360")}

	peerIDs := []string{
		"12D3KooWS3gci8DQvGTXQg2376YhFQsP4qb1oUbqQFyExvSJxU4h",
		"12D3KooWDqBK99eTGik3LWXNEyWpSFytBz3V58MNymNhNKsESHU7",
		"12D3KooWRUZLZuuDvp1QYs8RnrJc1Ljn31vuerurXLDorWijF9Zm",
		"12D3KooW9rW7bpmeP88icTT2wzkoWYi9PRyCeb6nvrWPGVbHMFbA",
	}

	transmiterAddresses := []string{
		"0xDA43dABB62fAE4906CCaeA0B9c47BA4D32b887Ac",
		"0xFB1a48F0E823c20b9046ac8822c41850562424b7",
		"0x4d098a1Aae2edb8eCd5db0454296704BDec66c89",
		"0x69E7064Aa8f1BCBB7A5Cd3AD8f34d2b7dB149B58",
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

	fmt.Println("signers: ", addrArrayToPrettyString(signersStr))
	fmt.Println("transmitters: ", addrArrayToPrettyString(transmittersStr))
	fmt.Println("f:", f)
	fmt.Println("onchainConfig:", onchainConfig)
	fmt.Println("offchainConfigVersion:", offchainConfigVersion)
	fmt.Printf("offchainConfig: 0x%s\n", hex.EncodeToString(offchainConfig))
}
