package actions

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink/core/services/keystore/chaintype"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
)

func DeployOCRv2Contracts(
	numberOfContracts int,
	linkTokenContract contracts.LinkToken,
	contractDeployer contracts.ContractDeployer,
	chainlinkNodes []*client.Chainlink,
	client blockchain.EVMClient,
) ([]contracts.OffchainAggregatorV2, error) {
	var ocrInstances []contracts.OffchainAggregatorV2
	for contractCount := 0; contractCount < numberOfContracts; contractCount++ {
		ocrInstance, err := contractDeployer.DeployOffchainAggregatorV2(
			linkTokenContract.Address(),
			contracts.DefaultOffChainAggregatorOptions(),
		)
		if err != nil {
			return nil, fmt.Errorf("OCRv2 instance deployment have failed: %w", err)
		}
		ocrInstances = append(ocrInstances, ocrInstance)
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			err = client.WaitForEvents()
			if err != nil {
				return nil, fmt.Errorf("failed to wait for OCRv2 contract deployments: %w", err)
			}
		}
	}
	err := client.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("error waiting for OCRv2 contract deployments: %w", err)
	}

	// Gather transmitter and address payees
	var transmitters, payees []string
	for _, node := range chainlinkNodes[1:] {
		addr, err := node.PrimaryEthAddress()
		if err != nil {
			return nil, fmt.Errorf("error getting node's primary ETH address: %w", err)
		}
		transmitters = append(transmitters, addr)
		payees = append(payees, client.GetDefaultWallet().Address())
	}

	// Set Payees
	for contractCount, ocrInstance := range ocrInstances {
		err = ocrInstance.SetPayees(transmitters, payees)
		if err != nil {
			return nil, fmt.Errorf("error settings OCR payees: %w", err)
		}
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			err = client.WaitForEvents()
			if err != nil {
				return nil, fmt.Errorf("failed to wait for setting OCR payees: %w", err)
			}
		}
	}
	err = client.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("error waiting for OCR contracts to set payees and transmitters: %w", err)
	}

	// Set Config
	transmitterAddresses, err := ChainlinkNodeAddresses(chainlinkNodes[1:])
	if err != nil {
		return nil, fmt.Errorf("getting node common addresses should not fail: %w", err)
	}
	for contractCount, ocrInstance := range ocrInstances {
		// Exclude the first node, which will be used as a bootstrapper
		// err = ocrInstance.SetConfig(
		// 	chainlinkNodes[1:],
		// 	contracts.DefaultOffChainAggregatorConfig(len(chainlinkNodes[1:])),
		// 	transmitterAddresses,
		// )
		// if err != nil {
		// 	return nil, fmt.Errorf("error setting OCR config for contract '%s': %w", ocrInstance.Address(), err)
		// }
		if (contractCount+1)%ContractDeploymentInterval == 0 { // For large amounts of contract deployments, space things out some
			err = client.WaitForEvents()
			if err != nil {
				return nil, fmt.Errorf("failed to wait for setting OCR config: %w", err)
			}
		}
	}
	err = client.WaitForEvents()
	if err != nil {
		return nil, fmt.Errorf("error waiting for OCR contracts to set config: %w", err)
	}
	return ocrInstances, nil
}

func BuildGeneralOCR2Config(
	t *testing.T,
	chainlinkNodes []*client.Chainlink,
	deltaProgress time.Duration,
	deltaResend time.Duration,
	deltaRound time.Duration,
	deltaGrace time.Duration,
	deltaStage time.Duration,
	rMax uint8,
	s []int,
	reportingPluginConfig []byte,
	maxDurationQuery time.Duration,
	maxDurationObservation time.Duration,
	maxDurationReport time.Duration,
	maxDurationShouldAcceptFinalizedReport time.Duration,
	maxDurationShouldTransmitAcceptedReport time.Duration,
	f int,
	onchainConfig []byte,
) contracts.OCRConfig {
	l := utils.GetTestLogger(t)
	_, oracleIdentities := GetOracleIdentities(t, chainlinkNodes)

	signerOnchainPublicKeys, transmitterAccounts, f_, onchainConfig_, offchainConfigVersion, offchainConfig, err := confighelper.ContractSetConfigArgsForTests(
		deltaProgress,
		deltaResend,
		deltaRound,
		deltaGrace,
		deltaStage,
		rMax,
		s,
		oracleIdentities,
		reportingPluginConfig,
		maxDurationQuery,
		maxDurationObservation,
		maxDurationReport,
		maxDurationShouldAcceptFinalizedReport,
		maxDurationShouldTransmitAcceptedReport,
		f,
		onchainConfig,
	)
	require.NoError(t, err, "Shouldn't fail ContractSetConfigArgsForTests")

	var signers []common.Address
	for _, signer := range signerOnchainPublicKeys {
		require.Equal(t, 20, len(signer), "OnChainPublicKey has wrong length for address")
		signers = append(signers, common.BytesToAddress(signer))
	}

	var transmitters []common.Address
	for _, transmitter := range transmitterAccounts {
		require.True(t, common.IsHexAddress(string(transmitter)), "TransmitAccount is not a valid Ethereum address")
		transmitters = append(transmitters, common.HexToAddress(string(transmitter)))
	}

	l.Info().Msg("Done building OCR2 config")
	return contracts.OCRConfig{
		Signers:               signers,
		Transmitters:          transmitters,
		F:                     f_,
		OnchainConfig:         onchainConfig_,
		OffchainConfigVersion: offchainConfigVersion,
		OffchainConfig:        offchainConfig,
	}
}

func getOracleIdentities(t *testing.T, chainlinkNodes []*client.Chainlink) ([]int, []confighelper.OracleIdentityExtra, error) {
	l := utils.GetTestLogger(t)
	S := make([]int, len(chainlinkNodes))
	oracleIdentities := make([]confighelper.OracleIdentityExtra, len(chainlinkNodes))
	sharedSecretEncryptionPublicKeys := make([]types.ConfigEncryptionPublicKey, len(chainlinkNodes))
	eg := &errgroup.Group{}
	for i, cl := range chainlinkNodes {
		index, chainlinkNode := i, cl
		eg.Go(func() error {
			address, err := chainlinkNode.PrimaryEthAddress()
			require.NoError(t, err, "Shouldn't fail getting primary ETH address from OCR node: index %d", index)
			ocr2Keys, err := chainlinkNode.MustReadOCR2Keys()
			require.NoError(t, err, "Shouldn't fail reading OCR2 keys from node")
			var ocr2Config client.OCR2KeyAttributes
			for _, key := range ocr2Keys.Data {
				if key.Attributes.ChainType == string(chaintype.EVM) {
					ocr2Config = key.Attributes
					break
				}
			}

			keys, err := chainlinkNode.MustReadP2PKeys()
			require.NoError(t, err, "Shouldn't fail reading P2P keys from node")
			p2pKeyID := keys.Data[0].Attributes.PeerID

			offchainPkBytes, err := hex.DecodeString(strings.TrimPrefix(ocr2Config.OffChainPublicKey, "ocr2off_evm_"))
			require.NoError(t, err, "failed to decode %s: %v", ocr2Config.OffChainPublicKey, err)

			offchainPkBytesFixed := [ed25519.PublicKeySize]byte{}
			n := copy(offchainPkBytesFixed[:], offchainPkBytes)
			require.Equal(t, ed25519.PublicKeySize, n, "Wrong number of elements copied")

			configPkBytes, err := hex.DecodeString(strings.TrimPrefix(ocr2Config.ConfigPublicKey, "ocr2cfg_evm_"))
			require.NoError(t, err, "failed to decode %s: %v", ocr2Config.ConfigPublicKey, err)

			configPkBytesFixed := [ed25519.PublicKeySize]byte{}
			n = copy(configPkBytesFixed[:], configPkBytes)
			require.Equal(t, ed25519.PublicKeySize, n, "Wrong number of elements copied")

			onchainPkBytes, err := hex.DecodeString(strings.TrimPrefix(ocr2Config.OnChainPublicKey, "ocr2on_evm_"))
			require.NoError(t, err, "failed to decode %s: %v", ocr2Config.OnChainPublicKey, err)

			sharedSecretEncryptionPublicKeys[index] = configPkBytesFixed
			oracleIdentities[index] = confighelper.OracleIdentityExtra{
				OracleIdentity: confighelper.OracleIdentity{
					OnchainPublicKey:  onchainPkBytes,
					OffchainPublicKey: offchainPkBytesFixed,
					PeerID:            p2pKeyID,
					TransmitAccount:   types.Account(address),
				},
				ConfigEncryptionPublicKey: configPkBytesFixed,
			}
			S[index] = 1
			return nil
		})
	}

	return S, oracleIdentities, eg.Wait()
}
