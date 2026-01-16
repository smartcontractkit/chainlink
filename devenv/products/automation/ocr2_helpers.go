package automation

import (
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"golang.org/x/sync/errgroup"
)

func GetOracleIdentitiesWithKeyIndexLocal(
	chainlinkNodes []*clclient.ChainlinkClient,
	keyIndex int,
) ([]int, []confighelper.OracleIdentityExtra, error) {
	S := make([]int, len(chainlinkNodes))
	oracleIdentities := make([]confighelper.OracleIdentityExtra, len(chainlinkNodes))
	sharedSecretEncryptionPublicKeys := make([]types.ConfigEncryptionPublicKey, len(chainlinkNodes))
	eg := &errgroup.Group{}
	for i, cl := range chainlinkNodes {
		index, chainlinkNode := i, cl
		eg.Go(func() error {
			addresses, err := chainlinkNode.EthAddresses()
			if err != nil {
				return err
			}
			ocr2Keys, err := chainlinkNode.MustReadOCR2Keys()
			if err != nil {
				return err
			}
			var ocr2Config clclient.OCR2KeyAttributes
			for _, key := range ocr2Keys.Data {
				if strings.EqualFold(key.Attributes.ChainType, "evm") {
					ocr2Config = key.Attributes
					break
				}
			}

			keys, err := chainlinkNode.MustReadP2PKeys()
			if err != nil {
				return err
			}
			p2pKeyID := keys.Data[0].Attributes.PeerID

			offchainPkBytes, err := hex.DecodeString(strings.TrimPrefix(ocr2Config.OffChainPublicKey, "ocr2off_evm_"))
			if err != nil {
				return err
			}

			offchainPkBytesFixed := [ed25519.PublicKeySize]byte{}
			n := copy(offchainPkBytesFixed[:], offchainPkBytes)
			if n != ed25519.PublicKeySize {
				return errors.New("wrong number of elements copied")
			}

			configPkBytes, err := hex.DecodeString(strings.TrimPrefix(ocr2Config.ConfigPublicKey, "ocr2cfg_evm_"))
			if err != nil {
				return err
			}

			configPkBytesFixed := [ed25519.PublicKeySize]byte{}
			n = copy(configPkBytesFixed[:], configPkBytes)
			if n != ed25519.PublicKeySize {
				return errors.New("wrong number of elements copied")
			}

			onchainPkBytes, err := hex.DecodeString(strings.TrimPrefix(ocr2Config.OnChainPublicKey, "ocr2on_evm_"))
			if err != nil {
				return err
			}

			sharedSecretEncryptionPublicKeys[index] = configPkBytesFixed
			oracleIdentities[index] = confighelper.OracleIdentityExtra{
				OracleIdentity: confighelper.OracleIdentity{
					OnchainPublicKey:  onchainPkBytes,
					OffchainPublicKey: offchainPkBytesFixed,
					PeerID:            p2pKeyID,
					TransmitAccount:   types.Account(addresses[keyIndex]),
				},
				ConfigEncryptionPublicKey: configPkBytesFixed,
			}
			S[index] = 1
			L.Debug().
				Interface("OnChainPK", onchainPkBytes).
				Interface("OffChainPK", offchainPkBytesFixed).
				Interface("ConfigPK", configPkBytesFixed).
				Str("PeerID", p2pKeyID).
				Str("Address", addresses[keyIndex]).
				Msg("Oracle identity")
			return nil
		})
	}

	return S, oracleIdentities, eg.Wait()
}
