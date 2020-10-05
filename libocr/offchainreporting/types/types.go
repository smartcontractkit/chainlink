package types

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/curve25519"
)

type OracleID int

type ConfigDigest [16]byte

func (c ConfigDigest) Hex() string {
	return fmt.Sprintf("%x", c[:])
}

func BytesToConfigDigest(b []byte) (g ConfigDigest) {
	configDigest := ConfigDigest{}
	copy(configDigest[:], b)
	return configDigest
}

type BinaryNetworkEndpoint interface {
		SendTo(payload []byte, to OracleID)
		Broadcast(payload []byte)
		Receive() <-chan BinaryMessageWithSender
		Start() error
		Close() error
}

type Bootstrapper interface {
	Start() error
	Close() error
}

type BinaryNetworkEndpointFactory interface {
	MakeEndpoint(cd ConfigDigest, peerIDs []string, bootstrappers []string, failureThreshold int) (BinaryNetworkEndpoint, error)
	PeerID() string
}

type BootstrapperFactory interface {
	MakeBootstrapper(cd ConfigDigest, peerIDs []string, bootstrappers []string, failureThreshold int) (Bootstrapper, error)
}

type BinaryMessageWithSender struct {
	Msg    []byte
	Sender OracleID
}

type Observation *big.Int

type DataSource interface {
			Observe(context.Context) (Observation, error)
}

type MonitoringEndpoint interface {
	SendLog(log []byte)
}

type ContractTransmitter interface {

		Transmit(
		ctx context.Context,
		report []byte, 		rs, ss [][32]byte, vs [32]byte, 	) error

	LatestTransmissionDetails(
		ctx context.Context,
	) (
		configDigest ConfigDigest,
		epoch uint32,
		round uint8,
		latestAnswer Observation,
		latestTimestamp time.Time,
		err error,
	)

	FromAddress() common.Address
}

type ContractConfigTracker interface {
	SubscribeToNewConfigs(ctx context.Context) (ContractConfigSubscription, error)
	LatestConfigDetails(ctx context.Context) (changedInBlock uint64, configDigest ConfigDigest, err error)
	ConfigFromLogs(ctx context.Context, changedInBlock uint64) (ContractConfig, error)

		LatestBlockHeight(ctx context.Context) (blockheight uint64, err error)
}

type ContractConfigSubscription interface {
		Configs() <-chan ContractConfig
	Close()
}

type ContractConfig struct {
	ConfigDigest         ConfigDigest
	Signers              []common.Address 	Transmitters         []common.Address
	Threshold            uint8
	EncodedConfigVersion uint64
	Encoded              []byte
}

type OffchainPublicKey ed25519.PublicKey

type OnChainSigningAddress common.Address

type SharedSecretEncryptionPublicKey [curve25519.PointSize]byte 
type PrivateKeys interface {

			SignOnChain(msg []byte) (signature []byte, err error)

			SignOffChain(msg []byte) (signature []byte, err error)

				ConfigDiffieHellman(base *[curve25519.ScalarSize]byte) (sharedPoint *[curve25519.PointSize]byte, err error)

			PublicKeyAddressOnChain() OnChainSigningAddress

		PublicKeyOffChain() OffchainPublicKey

		PublicKeyConfig() [curve25519.PointSize]byte
}
