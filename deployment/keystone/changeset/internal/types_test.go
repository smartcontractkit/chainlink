package internal

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

func Test_toNodeKeys(t *testing.T) {
	registryChainSel := chainsel.TEST_90000001
	registryChainID, err := chainsel.ChainIdFromSelector(registryChainSel.Selector)
	if err != nil {
		panic(err)
	}
	registryChainDetails, err := chainsel.GetChainDetailsByChainIDAndFamily(strconv.Itoa(int(registryChainID)), chainsel.FamilyEVM)
	if err != nil {
		panic(err)
	}
	if registryChainDetails.ChainName == "" {
		registryChainDetails.ChainName = strconv.FormatUint(registryChainDetails.ChainSelector, 10)
	}
	aptosChainDetails, err := chainsel.GetChainDetailsByChainIDAndFamily(strconv.Itoa(int(1)), chainsel.FamilyAptos)
	if err != nil {
		panic(err)
	}

	p2pID := p2pkey.MustNewV2XXXTestingOnly(big.NewInt(100))
	pubKey_1 := "11114981a6119ca3f932cdb8c402d71a72d672adae7849f581ecff8b8e1098e7"                    // valid csa key
	signing_1 := common.Hex2Bytes("11117293a4Cc2621b61193135a95928735e4795f")                         // valid eth address
	admin_1 := common.HexToAddress("0x1111567890123456789012345678901234567890")                      // valid eth address
	signing_2 := common.Hex2Bytes("eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee") // valid key
	var encryptionpubkey [32]byte
	if _, err := hex.Decode(encryptionpubkey[:], []byte(pubKey_1)); err != nil {
		panic(fmt.Sprintf("failed to decode pubkey %s: %v", encryptionpubkey, err))
	}
	keys := toNodeKeys(&deployment.Node{
		NodeID:    "p2p_123",
		Name:      "node 1",
		PeerID:    p2pID.PeerID(),
		CSAKey:    pubKey_1,
		AdminAddr: admin_1.String(),
		SelToOCRConfig: map[chainsel.ChainDetails]deployment.OCRConfig{
			registryChainDetails: {
				OffchainPublicKey:         types.OffchainPublicKey(common.FromHex("1111111111111111111111111111111111111111111111111111111111111111")),
				OnchainPublicKey:          signing_1,
				PeerID:                    p2pID.PeerID(),
				TransmitAccount:           types.Account(admin_1.String()),
				ConfigEncryptionPublicKey: encryptionpubkey,
				KeyBundleID:               "abcd",
			},
			aptosChainDetails: {
				TransmitAccount:  "ff",
				OnchainPublicKey: signing_2,
				KeyBundleID:      "aptos",
			},
		},
	}, registryChainSel.Selector)

	require.Equal(t, NodeKeys{
		EthAddress:            admin_1.String(),
		AptosBundleID:         "aptos",
		AptosOnchainPublicKey: hex.EncodeToString(signing_2),
		P2PPeerID:             strings.TrimPrefix(p2pID.PeerID().String(), "p2p_"),
		OCR2BundleID:          "abcd",
		OCR2OnchainPublicKey:  hex.EncodeToString(signing_1),
		OCR2OffchainPublicKey: "1111111111111111111111111111111111111111111111111111111111111111",
		OCR2ConfigPublicKey:   pubKey_1,
		CSAPublicKey:          pubKey_1,
		EncryptionPublicKey:   pubKey_1,
	}, keys)
}
