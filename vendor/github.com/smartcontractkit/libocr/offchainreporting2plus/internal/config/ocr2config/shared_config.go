package ocr2config

import (
	"bytes"
	cryptorand "crypto/rand"
	"fmt"
	"math"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ethcontractconfig"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"golang.org/x/crypto/sha3"
)

// SharedConfig is the configuration shared by all oracles running an instance
// of the protocol. It's disseminated through the smart contract,
// but parts of it are encrypted so that only oracles can access them.
type SharedConfig struct {
	PublicConfig
	SharedSecret *[config.SharedSecretSize]byte
}

func (c *SharedConfig) LeaderSelectionKey() [16]byte {
	var result [16]byte
	h := sha3.NewLegacyKeccak256()
	h.Write(c.SharedSecret[:])
	h.Write([]byte("chainlink offchain reporting v1 leader selection key"))

	copy(result[:], h.Sum(nil))
	return result
}

func (c *SharedConfig) TransmissionOrderKey() [16]byte {
	var result [16]byte
	h := sha3.NewLegacyKeccak256()
	h.Write(c.SharedSecret[:])
	h.Write([]byte("chainlink offchain reporting v1 transmission order key"))

	copy(result[:], h.Sum(nil))
	return result
}

func SharedConfigFromContractConfig(
	skipResourceExhaustionChecks bool,
	change types.ContractConfig,
	offchainKeyring types.OffchainKeyring,
	onchainKeyring types.OnchainKeyring,
	peerID string,
	transmitAccount types.Account,
) (SharedConfig, commontypes.OracleID, error) {
	publicConfig, encSharedSecret, err := publicConfigFromContractConfig(skipResourceExhaustionChecks, change)
	if err != nil {
		return SharedConfig{}, 0, err
	}

	oracleID := commontypes.OracleID(math.MaxUint8)
	{
		onchainPublicKey := onchainKeyring.PublicKey()
		offchainPublicKey := offchainKeyring.OffchainPublicKey()
		var found bool
		for i, identity := range publicConfig.OracleIdentities {
			if bytes.Equal(identity.OnchainPublicKey, onchainPublicKey) {
				if identity.OffchainPublicKey != offchainPublicKey {
					return SharedConfig{}, 0, errors.Errorf(
						"OnchainPublicKey %x in publicConfig matches "+
							"mine, but OffchainPublicKey does not: %v (config) vs %v (mine)",
						onchainPublicKey, identity.OffchainPublicKey, offchainPublicKey)
				}
				if identity.PeerID != peerID {
					return SharedConfig{}, 0, errors.Errorf(
						"OnchainPublicKey %x in publicConfig matches "+
							"mine, but PeerID does not: %v (config) vs %v (mine)",
						onchainPublicKey, identity.PeerID, peerID)
				}
				if identity.TransmitAccount != transmitAccount {
					return SharedConfig{}, 0, errors.Errorf(
						"OnchainPublicKey %x in publicConfig matches "+
							"mine, but TransmitAccount does not: %v (config) vs %v (mine)",
						onchainPublicKey, identity.TransmitAccount, transmitAccount)
				}
				oracleID = commontypes.OracleID(i)
				found = true
			}
		}

		if !found {
			return SharedConfig{},
				0,
				fmt.Errorf("could not find my OnchainPublicKey %x in publicConfig", onchainPublicKey)
		}
	}

	x, err := encSharedSecret.Decrypt(oracleID, offchainKeyring)
	if err != nil {
		return SharedConfig{}, 0, fmt.Errorf("could not decrypt shared secret: %w", err)
	}

	return SharedConfig{
		publicConfig,
		x,
	}, oracleID, nil

}

func XXXContractSetConfigArgsFromSharedConfig(
	c SharedConfig,
	sharedSecretEncryptionPublicKeys []types.ConfigEncryptionPublicKey,
) (
	signers []types.OnchainPublicKey,
	transmitters []types.Account,
	f uint8,
	onchainConfig []byte,
	offchainConfigVersion uint64,
	offchainConfig_ []byte,
	err error,
) {
	offChainPublicKeys := []types.OffchainPublicKey{}
	peerIDs := []string{}
	for _, identity := range c.OracleIdentities {
		signers = append(signers, identity.OnchainPublicKey)
		transmitters = append(transmitters, identity.TransmitAccount)
		offChainPublicKeys = append(offChainPublicKeys, identity.OffchainPublicKey)
		peerIDs = append(peerIDs, identity.PeerID)
	}
	f = uint8(c.F)
	onchainConfig = c.OnchainConfig
	offchainConfigVersion = config.OCR2OffchainConfigVersion
	offchainConfig_ = (offchainConfig{
		c.DeltaProgress,
		c.DeltaResend,
		c.DeltaRound,
		c.DeltaGrace,
		c.DeltaStage,
		c.RMax,
		c.S,
		offChainPublicKeys,
		peerIDs,
		c.ReportingPluginConfig,
		c.MaxDurationQuery,
		c.MaxDurationObservation,
		c.MaxDurationReport,
		c.MaxDurationShouldAcceptFinalizedReport,
		c.MaxDurationShouldTransmitAcceptedReport,
		config.XXXEncryptSharedSecret(
			sharedSecretEncryptionPublicKeys,
			c.SharedSecret,
			cryptorand.Reader,
		),
	}).serialize()
	err = nil
	return
}

func XXXContractSetConfigArgsFromSharedConfigEthereum(
	c SharedConfig,
	sharedSecretEncryptionPublicKeys []types.ConfigEncryptionPublicKey,
) (
	setConfigArgs ethcontractconfig.SetConfigArgs,
	err error,
) {
	signerOnchainPublicKeys, transmitterAccounts, f, onchainConfig, offchainConfigVersion, offchainConfig, err :=
		XXXContractSetConfigArgsFromSharedConfig(c, sharedSecretEncryptionPublicKeys)
	if err != nil {
		return ethcontractconfig.SetConfigArgs{}, err
	}

	var signers []common.Address
	for _, signer := range signerOnchainPublicKeys {
		if len(signer) != 20 {
			return ethcontractconfig.SetConfigArgs{}, fmt.Errorf("OnChainPublicKey has wrong length for address")
		}
		signers = append(signers, common.BytesToAddress(signer))
	}

	var transmitters []common.Address
	for _, transmitter := range transmitterAccounts {
		if !common.IsHexAddress(string(transmitter)) {
			return ethcontractconfig.SetConfigArgs{}, fmt.Errorf("TransmitAccount is not a valid Ethereum address")
		}
		transmitters = append(transmitters, common.HexToAddress(string(transmitter)))
	}

	return ethcontractconfig.SetConfigArgs{
		signers,
		transmitters,
		f,
		onchainConfig,
		offchainConfigVersion,
		offchainConfig,
	}, nil
}
