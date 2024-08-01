package config

import (
	"bytes"
	cryptorand "crypto/rand"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting/types"
	"golang.org/x/crypto/sha3"
)

// SharedConfig is the configuration shared by all oracles running an instance
// of the protocol. It's disseminated through the smart contract,
// but parts of it are encrypted so that only oracles can access them.
type SharedConfig struct {
	PublicConfig
	SharedSecret *[SharedSecretSize]byte
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
	chainID *big.Int,
	skipChainSpecificChecks bool,
	change types.ContractConfig,
	privateKeys types.PrivateKeys,
	peerID string,
	transmitAddress common.Address,
) (SharedConfig, commontypes.OracleID, error) {
	publicConfig, encSharedSecret, err := publicConfigFromContractConfig(chainID, skipChainSpecificChecks, change)
	if err != nil {
		return SharedConfig{}, 0, err
	}

	oracleID := commontypes.OracleID(math.MaxUint8)
	{
		var found bool
		for i, identity := range publicConfig.OracleIdentities {
			address := privateKeys.PublicKeyAddressOnChain()
			offchainPublicKey := privateKeys.PublicKeyOffChain()
			if identity.OnChainSigningAddress == address {
				if !bytes.Equal(identity.OffchainPublicKey, offchainPublicKey) {
					return SharedConfig{}, 0, errors.Errorf(
						"OnChainSigningAddress (0x%x) in publicConfig matches "+
							"mine, but OffchainPublicKey does not: %v (config) vs %v (mine)",
						address, identity.OffchainPublicKey, offchainPublicKey)
				}
				if identity.PeerID != peerID {
					return SharedConfig{}, 0, errors.Errorf(
						"OnChainSigningAddress (0x%x) in publicConfig matches "+
							"mine, but PeerID does not: %v (config) vs %v (mine)",
						address, identity.PeerID, peerID)
				}
				if identity.TransmitAddress != transmitAddress {
					return SharedConfig{}, 0, errors.Errorf(
						"OnChainSigningAddress (0x%x) in publicConfig matches "+
							"mine, but TransmitAddress does not: 0x%x (config) vs 0x%x (mine)",
						address, identity.TransmitAddress, transmitAddress)
				}
				oracleID = commontypes.OracleID(i)
				found = true
			}
		}

		if !found {
			return SharedConfig{},
				0,
				errors.Errorf("Could not find my OnChainSigningAddress 0x%x in publicConfig", privateKeys.PublicKeyAddressOnChain())
		}
	}

	x, err := encSharedSecret.Decrypt(oracleID, privateKeys)
	if err != nil {
		return SharedConfig{}, 0, errors.Wrapf(err, "could not decrypt shared secret")
	}

	return SharedConfig{
		publicConfig,
		x,
	}, oracleID, nil

}

func XXXContractSetConfigArgsFromSharedConfig(
	c SharedConfig,
	sharedSecretEncryptionPublicKeys []types.SharedSecretEncryptionPublicKey,
) (
	signers []common.Address,
	transmitters []common.Address,
	threshold uint8,
	encodedConfigVersion uint64,
	encodedConfig []byte,
	err error,
) {
	offChainPublicKeys := []types.OffchainPublicKey{}
	peerIDs := []string{}
	for _, identity := range c.OracleIdentities {
		signers = append(signers, common.Address(identity.OnChainSigningAddress))
		transmitters = append(transmitters, identity.TransmitAddress)
		offChainPublicKeys = append(offChainPublicKeys, identity.OffchainPublicKey)
		peerIDs = append(peerIDs, identity.PeerID)
	}
	threshold = uint8(c.F)
	encodedConfigVersion = 1
	encodedConfig = (setConfigEncodedComponents{
		c.DeltaProgress,
		c.DeltaResend,
		c.DeltaRound,
		c.DeltaGrace,
		c.DeltaC,
		c.AlphaPPB,
		c.DeltaStage,
		c.RMax,
		c.S,
		offChainPublicKeys,
		peerIDs,
		XXXEncryptSharedSecret(
			sharedSecretEncryptionPublicKeys,
			c.SharedSecret,
			cryptorand.Reader,
		),
	}).encode()
	err = nil
	return
}
