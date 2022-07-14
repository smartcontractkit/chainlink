package main

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/ocr2vrf/altbn_128"
	"github.com/smartcontractkit/ocr2vrf/dkg"
	dkgContract "github.com/smartcontractkit/ocr2vrf/gethwrappers/dkg"
	"github.com/smartcontractkit/ocr2vrf/gethwrappers/vrf"
	ocr2vrftypes "github.com/smartcontractkit/ocr2vrf/types"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/edwards25519"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
)

func deployDKG(e helpers.Environment) common.Address {
	_, tx, _, err := dkgContract.DeployDKG(e.Owner, e.Ec)
	helpers.PanicErr(err)
	return helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
}

func addClientToDKG(e helpers.Environment, dkgAddress string, keyID string, clientAddress string) {
	keyIDBytes := decodeHexTo32ByteArray(keyID)

	dkg, err := dkgContract.NewDKG(common.HexToAddress(dkgAddress), e.Ec)
	helpers.PanicErr(err)

	tx, err := dkg.AddClient(e.Owner, keyIDBytes, common.HexToAddress(clientAddress))
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
}

func removeClientFromDKG(e helpers.Environment, dkgAddress string, keyID string, clientAddress string) {
	keyIDBytes := decodeHexTo32ByteArray(keyID)

	dkg, err := dkgContract.NewDKG(common.HexToAddress(dkgAddress), e.Ec)
	helpers.PanicErr(err)

	tx, err := dkg.RemoveClient(e.Owner, keyIDBytes, common.HexToAddress(clientAddress))
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
}

func setDKGConfig(e helpers.Environment, dkgAddress string, c dkgSetConfigArgs) {
	oracleIdentities := toOraclesIdentityList(
		helpers.ParseAddressSlice(c.onchainPubKeys),
		strings.Split(c.offchainPubKeys, ","),
		strings.Split(c.configPubKeys, ","),
		strings.Split(c.peerIDs, ","),
		strings.Split(c.transmitters, ","))

	ed25519Suite := edwards25519.NewBlakeSHA256Ed25519()
	var signingKeys []kyber.Point
	for _, signingKey := range strings.Split(c.dkgSigningPubKeys, ",") {
		signingKeyBytes, err := hex.DecodeString(signingKey)
		helpers.PanicErr(err)
		signingKeyPoint := ed25519Suite.Point()
		helpers.PanicErr(signingKeyPoint.UnmarshalBinary(signingKeyBytes))
		signingKeys = append(signingKeys, signingKeyPoint)
	}

	altbn128Suite := &altbn_128.PairingSuite{}
	var encryptionKeys []kyber.Point
	for _, encryptionKey := range strings.Split(c.dkgEncryptionPubKeys, ",") {
		encryptionKeyBytes, err := hex.DecodeString(encryptionKey)
		helpers.PanicErr(err)
		encryptionKeyPoint := altbn128Suite.G1().Point()
		helpers.PanicErr(encryptionKeyPoint.UnmarshalBinary(encryptionKeyBytes))
		encryptionKeys = append(encryptionKeys, encryptionKeyPoint)
	}

	keyIDBytes := decodeHexTo32ByteArray(c.keyID)

	offchainConfig, err := dkg.OffchainConfig(dkg.EncryptionPublicKeys(encryptionKeys), dkg.SigningPublicKeys(signingKeys), &altbn_128.G1{}, &ocr2vrftypes.PairingTranslation{
		&altbn_128.PairingSuite{},
	})
	helpers.PanicErr(err)
	onchainConfig, err := dkg.OnchainConfig(dkg.KeyID(keyIDBytes))
	helpers.PanicErr(err)

	fmt.Println("dkg offchain config:", hex.EncodeToString(offchainConfig))
	fmt.Println("dkg onchain config:", hex.EncodeToString(onchainConfig))

	_, _, f, onchainConfig, offchainConfigVersion, offchainConfig, err := confighelper.ContractSetConfigArgsForTests(
		c.deltaProgress,
		c.deltaResend,
		c.deltaRound,
		c.deltaGrace,
		c.deltaStage,
		c.maxRounds,
		helpers.ParseIntSlice(c.schedule),
		oracleIdentities,
		offchainConfig,
		c.maxDurationQuery,
		c.maxDurationObservation,
		c.maxDurationReport,
		c.maxDurationAccept,
		c.maxDurationTransmit,
		int(c.f),
		onchainConfig)

	helpers.PanicErr(err)

	dkg := newDKG(common.HexToAddress(dkgAddress), e.Ec)

	tx, err := dkg.SetConfig(e.Owner, helpers.ParseAddressSlice(c.onchainPubKeys), helpers.ParseAddressSlice(c.transmitters), f, onchainConfig, offchainConfigVersion, offchainConfig)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
}

func toOraclesIdentityList(onchainPubKeys []common.Address, offchainPubKeys, configPubKeys, peerIDs, transmitters []string) []confighelper.OracleIdentityExtra {
	offchainPubKeysBytes := []types.OffchainPublicKey{}
	for _, pkHex := range offchainPubKeys {
		pkBytes, err := hex.DecodeString(pkHex)
		helpers.PanicErr(err)
		pkBytesFixed := [ed25519.PublicKeySize]byte{}
		n := copy(pkBytesFixed[:], pkBytes)
		if n != ed25519.PublicKeySize {
			panic("wrong num elements copied")
		}

		offchainPubKeysBytes = append(offchainPubKeysBytes, types.OffchainPublicKey(pkBytesFixed))
	}

	configPubKeysBytes := []types.ConfigEncryptionPublicKey{}
	for _, pkHex := range configPubKeys {
		pkBytes, err := hex.DecodeString(pkHex)
		helpers.PanicErr(err)

		pkBytesFixed := [ed25519.PublicKeySize]byte{}
		n := copy(pkBytesFixed[:], pkBytes)
		if n != ed25519.PublicKeySize {
			panic("wrong num elements copied")
		}

		configPubKeysBytes = append(configPubKeysBytes, types.ConfigEncryptionPublicKey(pkBytesFixed))
	}

	o := []confighelper.OracleIdentityExtra{}
	for index := range configPubKeys {
		o = append(o, confighelper.OracleIdentityExtra{
			OracleIdentity: confighelper.OracleIdentity{
				OnchainPublicKey:  onchainPubKeys[index][:],
				OffchainPublicKey: offchainPubKeysBytes[index],
				PeerID:            peerIDs[index],
				TransmitAccount:   types.Account(transmitters[index]),
			},
			ConfigEncryptionPublicKey: configPubKeysBytes[index],
		})
	}
	return o
}

func deployVRF(e helpers.Environment, dkgAddress string, keyID string) common.Address {
	var keyIDBytes [32]byte
	copy(keyIDBytes[:], keyID)
	fmt.Printf("Key ID in bytes: 0x%x \n", keyIDBytes)

	_, tx, _, err := vrf.DeployVRF(e.Owner, e.Ec, common.HexToAddress(dkgAddress), keyIDBytes)
	helpers.PanicErr(err)
	return helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
}

func newVRF(addr common.Address, client *ethclient.Client) *vrf.VRF {
	vrf, err := vrf.NewVRF(addr, client)
	helpers.PanicErr(err)
	return vrf
}

func newDKG(addr common.Address, client *ethclient.Client) *dkgContract.DKG {
	dkg, err := dkgContract.NewDKG(addr, client)
	helpers.PanicErr(err)
	return dkg
}

func decodeHexTo32ByteArray(val string) (byteArray [32]byte) {
	decoded, err := hex.DecodeString(val)
	helpers.PanicErr(err)
	if len(decoded) != 32 {
		panic(fmt.Sprintf("expected value to be 32 bytes but received %d bytes", len(decoded)))
	}
	copy(byteArray[:], decoded)
	return
}
