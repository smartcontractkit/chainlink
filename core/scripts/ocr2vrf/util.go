package main

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/ocr2vrf/altbn_128"
	ocr2vrfdkg "github.com/smartcontractkit/ocr2vrf/dkg"
	"github.com/smartcontractkit/ocr2vrf/ocr2vrf"
	ocr2vrftypes "github.com/smartcontractkit/ocr2vrf/types"
	"github.com/urfave/cli"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/edwards25519"

	"github.com/smartcontractkit/chainlink/core/cmd"
	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/arbitrum_beacon_vrf_consumer"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/arbitrum_dkg"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/arbitrum_vrf_beacon"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/arbitrum_vrf_coordinator"
	dkgContract "github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/dkg"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/load_test_beacon_consumer"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/vrf_beacon"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/vrf_beacon_consumer"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/vrf_coordinator"
	"github.com/smartcontractkit/chainlink/core/logger"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
)

// The interfaces below define methods used by the contract wrappers.
// We use these interfaces so that we don't have to duplicate arbitrum and non-arbitrum
// calls, which are essentially identical, just slightly different go wrappers.
type beaconConsumer interface {
	TestRequestRandomness(opts *bind.TransactOpts, numWords uint16, subID uint64, confirmationDelayArg *big.Int) (*gethTypes.Transaction, error)
	TestRequestRandomnessFulfillment(opts *bind.TransactOpts, subID uint64, numWords uint16, confirmationDelayArg *big.Int, callbackGasLimit uint32, arguments []byte) (*gethTypes.Transaction, error)
	IBeaconPeriodBlocks(opts *bind.CallOpts) (*big.Int, error)
	SRequestsIDs(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error)
	TestRedeemRandomness(opts *bind.TransactOpts, requestID *big.Int) (*gethTypes.Transaction, error)
	SReceivedRandomnessByRequestID(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (*big.Int, error)
}

type beacon interface {
	SKeyID(opts *bind.CallOpts) ([32]byte, error)
	SProvingKeyHash(opts *bind.CallOpts) ([32]byte, error)
	SetConfig(opts *bind.TransactOpts, _signers []common.Address, _transmitters []common.Address, _f uint8, _onchainConfig []byte, _offchainConfigVersion uint64, _offchainConfig []byte) (*gethTypes.Transaction, error)
}

type coordinator interface {
	SetProducer(opts *bind.TransactOpts, addr common.Address) (*gethTypes.Transaction, error)
	RedeemRandomness(opts *bind.TransactOpts, requestID *big.Int) (*gethTypes.Transaction, error)
	RequestRandomness(opts *bind.TransactOpts, numWords uint16, subID uint64, confirmationDelayArg *big.Int) (*gethTypes.Transaction, error)
}

type dkg interface {
	AddClient(opts *bind.TransactOpts, keyID [32]byte, clientAddress common.Address) (*gethTypes.Transaction, error)
	RemoveClient(opts *bind.TransactOpts, keyID [32]byte, clientAddress common.Address) (*gethTypes.Transaction, error)
	SetConfig(opts *bind.TransactOpts, _signers []common.Address, _transmitters []common.Address, _f uint8, _onchainConfig []byte, _offchainConfigVersion uint64, _offchainConfig []byte) (*gethTypes.Transaction, error)
}

func deployDKG(e helpers.Environment) common.Address {
	var tx *gethTypes.Transaction
	var err error
	if helpers.IsArbitrumChainID(e.ChainID) {
		_, tx, _, err = arbitrum_dkg.DeployArbitrumDKG(e.Owner, e.Ec)
	} else {
		_, tx, _, err = dkgContract.DeployDKG(e.Owner, e.Ec)
	}

	helpers.PanicErr(err)
	return helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
}

func deployVRFCoordinator(e helpers.Environment, beaconPeriodBlocks *big.Int, linkAddress string) common.Address {
	var tx *gethTypes.Transaction
	var err error
	if helpers.IsArbitrumChainID(e.ChainID) {
		_, tx, _, err = arbitrum_vrf_coordinator.DeployArbitrumVRFCoordinator(e.Owner, e.Ec, beaconPeriodBlocks, common.HexToAddress(linkAddress))
	} else {
		_, tx, _, err = vrf_coordinator.DeployVRFCoordinator(e.Owner, e.Ec, beaconPeriodBlocks, common.HexToAddress(linkAddress))
	}

	helpers.PanicErr(err)
	return helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
}

func deployVRFBeacon(e helpers.Environment, coordinatorAddress, linkAddress, dkgAddress, keyID string) common.Address {
	keyIDBytes := decodeHexTo32ByteArray(keyID)
	var tx *gethTypes.Transaction
	var err error
	if helpers.IsArbitrumChainID(e.ChainID) {
		_, tx, _, err = arbitrum_vrf_beacon.DeployArbitrumVRFBeacon(
			e.Owner, e.Ec, common.HexToAddress(coordinatorAddress), common.HexToAddress(coordinatorAddress), common.HexToAddress(dkgAddress), keyIDBytes)
	} else {
		_, tx, _, err = vrf_beacon.DeployVRFBeacon(
			e.Owner, e.Ec, common.HexToAddress(coordinatorAddress), common.HexToAddress(coordinatorAddress), common.HexToAddress(dkgAddress), keyIDBytes)
	}

	helpers.PanicErr(err)
	return helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
}

func deployVRFBeaconCoordinatorConsumer(e helpers.Environment, coordinatorAddress string, shouldFail bool, beaconPeriodBlocks *big.Int) common.Address {
	var tx *gethTypes.Transaction
	var err error
	if helpers.IsArbitrumChainID(e.ChainID) {
		_, tx, _, err = arbitrum_beacon_vrf_consumer.DeployArbitrumBeaconVRFConsumer(
			e.Owner, e.Ec, common.HexToAddress(coordinatorAddress), shouldFail, beaconPeriodBlocks)
	} else {
		_, tx, _, err = vrf_beacon_consumer.DeployBeaconVRFConsumer(e.Owner, e.Ec, common.HexToAddress(coordinatorAddress), shouldFail, beaconPeriodBlocks)
	}

	helpers.PanicErr(err)
	return helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
}

func deployLoadTestVRFBeaconCoordinatorConsumer(e helpers.Environment, coordinatorAddress string, shouldFail bool, beaconPeriodBlocks *big.Int) common.Address {
	_, tx, _, err := load_test_beacon_consumer.DeployLoadTestBeaconVRFConsumer(e.Owner, e.Ec, common.HexToAddress(coordinatorAddress), shouldFail, beaconPeriodBlocks)
	helpers.PanicErr(err)
	return helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
}

func addClientToDKG(e helpers.Environment, dkgAddress string, keyID string, clientAddress string) {
	keyIDBytes := decodeHexTo32ByteArray(keyID)

	d := newDKG(e, common.HexToAddress(dkgAddress))

	tx, err := d.AddClient(e.Owner, keyIDBytes, common.HexToAddress(clientAddress))
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
}

func removeClientFromDKG(e helpers.Environment, dkgAddress string, keyID string, clientAddress string) {
	keyIDBytes := decodeHexTo32ByteArray(keyID)

	d := newDKG(e, common.HexToAddress(dkgAddress))

	tx, err := d.RemoveClient(e.Owner, keyIDBytes, common.HexToAddress(clientAddress))
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

	offchainConfig, err := ocr2vrfdkg.OffchainConfig(encryptionKeys, signingKeys, &altbn_128.G1{}, &ocr2vrftypes.PairingTranslation{
		Suite: &altbn_128.PairingSuite{},
	})
	helpers.PanicErr(err)
	onchainConfig, err := ocr2vrfdkg.OnchainConfig(ocr2vrfdkg.KeyID(keyIDBytes))
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

	d := newDKG(e, common.HexToAddress(dkgAddress))

	tx, err := d.SetConfig(e.Owner, helpers.ParseAddressSlice(c.onchainPubKeys), helpers.ParseAddressSlice(c.transmitters), f, onchainConfig, offchainConfigVersion, offchainConfig)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
}

func setVRFBeaconConfig(e helpers.Environment, vrfBeaconAddr string, c vrfBeaconSetConfigArgs) {
	oracleIdentities := toOraclesIdentityList(
		helpers.ParseAddressSlice(c.onchainPubKeys),
		strings.Split(c.offchainPubKeys, ","),
		strings.Split(c.configPubKeys, ","),
		strings.Split(c.peerIDs, ","),
		strings.Split(c.transmitters, ","))

	confDelays := make(map[uint32]struct{})
	for _, c := range strings.Split(c.confDelays, ",") {
		confDelay, err := strconv.ParseUint(c, 0, 32)
		helpers.PanicErr(err)
		confDelays[uint32(confDelay)] = struct{}{}
	}

	onchainConfig := ocr2vrf.OnchainConfig(confDelays)

	_, _, f, onchainConfig, offchainConfigVersion, offchainConfig, err := confighelper.ContractSetConfigArgsForTests(
		c.deltaProgress,
		c.deltaResend,
		c.deltaRound,
		c.deltaGrace,
		c.deltaStage,
		c.maxRounds,
		helpers.ParseIntSlice(c.schedule),
		oracleIdentities,
		nil, // off-chain config
		c.maxDurationQuery,
		c.maxDurationObservation,
		c.maxDurationReport,
		c.maxDurationAccept,
		c.maxDurationTransmit,
		int(c.f),
		onchainConfig)

	helpers.PanicErr(err)

	b := newVRFBeacon(e, common.HexToAddress(vrfBeaconAddr))

	tx, err := b.SetConfig(e.Owner, helpers.ParseAddressSlice(c.onchainPubKeys), helpers.ParseAddressSlice(c.transmitters), f, onchainConfig, offchainConfigVersion, offchainConfig)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
}

func setProducer(e helpers.Environment, vrfCoordinatorAddr, vrfBeaconAddr string) {
	coordinator := newVRFCoordinator(e, common.HexToAddress(vrfCoordinatorAddr))

	tx, err := coordinator.SetProducer(e.Owner, common.HexToAddress(vrfBeaconAddr))
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

func requestRandomness(e helpers.Environment, coordinatorAddress string, numWords uint16, subID uint64, confDelay *big.Int) {
	coordinator := newVRFCoordinator(e, common.HexToAddress(coordinatorAddress))

	tx, err := coordinator.RequestRandomness(e.Owner, numWords, subID, confDelay)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
}

func redeemRandomness(e helpers.Environment, coordinatorAddress string, requestID *big.Int) {
	coordinator := newVRFCoordinator(e, common.HexToAddress(coordinatorAddress))

	tx, err := coordinator.RedeemRandomness(e.Owner, requestID)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
}

func requestRandomnessFromConsumer(e helpers.Environment, consumerAddress string, numWords uint16, subID uint64, confDelay *big.Int) *big.Int {
	consumer := newBeaconConsumer(e, common.HexToAddress(consumerAddress))

	tx, err := consumer.TestRequestRandomness(e.Owner, numWords, subID, confDelay)
	helpers.PanicErr(err)
	receipt := helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)

	periodBlocks, err := consumer.IBeaconPeriodBlocks(nil)
	helpers.PanicErr(err)

	blockNumber := receipt.BlockNumber
	periodOffset := new(big.Int).Mod(blockNumber, periodBlocks)
	nextBeaconOutputHeight := new(big.Int).Sub(new(big.Int).Add(blockNumber, periodBlocks), periodOffset)

	fmt.Println("nextBeaconOutputHeight: ", nextBeaconOutputHeight)

	requestID, err := consumer.SRequestsIDs(nil, nextBeaconOutputHeight, confDelay)
	helpers.PanicErr(err)
	fmt.Println("requestID: ", requestID)

	return requestID
}

func requestRandomnessCallback(
	e helpers.Environment,
	consumerAddress string,
	numWords uint16,
	subID uint64,
	confDelay *big.Int,
	callbackGasLimit uint32,
	args []byte,
) (requestID *big.Int) {
	consumer := newBeaconConsumer(e, common.HexToAddress(consumerAddress))

	tx, err := consumer.TestRequestRandomnessFulfillment(e.Owner, subID, numWords, confDelay, callbackGasLimit, args)
	helpers.PanicErr(err)
	receipt := helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID, "TestRequestRandomnessFulfillment")

	periodBlocks, err := consumer.IBeaconPeriodBlocks(nil)
	helpers.PanicErr(err)

	blockNumber := receipt.BlockNumber
	periodOffset := new(big.Int).Mod(blockNumber, periodBlocks)
	nextBeaconOutputHeight := new(big.Int).Sub(new(big.Int).Add(blockNumber, periodBlocks), periodOffset)

	fmt.Println("nextBeaconOutputHeight: ", nextBeaconOutputHeight)

	requestID, err = consumer.SRequestsIDs(nil, nextBeaconOutputHeight, confDelay)
	helpers.PanicErr(err)
	fmt.Println("requestID: ", requestID)

	return requestID
}

func redeemRandomnessFromConsumer(e helpers.Environment, consumerAddress string, requestID *big.Int, numWords int64) {
	consumer := newBeaconConsumer(e, common.HexToAddress(consumerAddress))

	tx, err := consumer.TestRedeemRandomness(e.Owner, requestID)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)

	printRandomnessFromConsumer(consumer, requestID, numWords)
}

func printRandomnessFromConsumer(consumer beaconConsumer, requestID *big.Int, numWords int64) {
	for i := int64(0); i < numWords; i++ {
		randomness, err := consumer.SReceivedRandomnessByRequestID(nil, requestID, big.NewInt(i))
		helpers.PanicErr(err)
		fmt.Println("random words index", i, ":", randomness.String())
	}
}

func newVRFCoordinator(e helpers.Environment, addr common.Address) (c coordinator) {
	var err error
	if helpers.IsArbitrumChainID(e.ChainID) {
		c, err = arbitrum_vrf_coordinator.NewArbitrumVRFCoordinator(addr, e.Ec)
	} else {
		c, err = vrf_coordinator.NewVRFCoordinator(addr, e.Ec)
	}

	helpers.PanicErr(err)
	return
}

func newDKG(e helpers.Environment, addr common.Address) (d dkg) {
	var err error
	if helpers.IsArbitrumChainID(e.ChainID) {
		d, err = arbitrum_dkg.NewArbitrumDKG(addr, e.Ec)
	} else {
		d, err = dkgContract.NewDKG(addr, e.Ec)
	}

	helpers.PanicErr(err)
	return
}

func newBeaconConsumer(e helpers.Environment, addr common.Address) (consumer beaconConsumer) {
	var err error
	if helpers.IsArbitrumChainID(e.ChainID) {
		consumer, err = arbitrum_beacon_vrf_consumer.NewArbitrumBeaconVRFConsumer(addr, e.Ec)
	} else {
		consumer, err = vrf_beacon_consumer.NewBeaconVRFConsumer(addr, e.Ec)
	}

	helpers.PanicErr(err)
	return
}

func newLoadTestVRFBeaconCoordinatorConsumer(addr common.Address, client *ethclient.Client) *load_test_beacon_consumer.LoadTestBeaconVRFConsumer {
	consumer, err := load_test_beacon_consumer.NewLoadTestBeaconVRFConsumer(addr, client)
	helpers.PanicErr(err)
	return consumer
}

func newVRFBeacon(e helpers.Environment, addr common.Address) (b beacon) {
	var err error
	if helpers.IsArbitrumChainID(e.ChainID) {
		b, err = arbitrum_vrf_beacon.NewArbitrumVRFBeacon(addr, e.Ec)
	} else {
		b, err = vrf_beacon.NewVRFBeacon(addr, e.Ec)
	}

	helpers.PanicErr(err)
	return
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

func setupOCR2VRFNodeFromClient(client *cmd.Client, context *cli.Context) *cmd.SetupOCR2VRFNodePayload {
	payload, err := client.ConfigureOCR2VRFNode(context)
	helpers.PanicErr(err)

	return payload
}

func configureEnvironmentVariables() {
	helpers.PanicErr(os.Setenv("FEATURE_OFFCHAIN_REPORTING2", "true"))
	helpers.PanicErr(os.Setenv("SKIP_DATABASE_PASSWORD_COMPLEXITY_CHECK", "true"))
}

func resetDatabase(client *cmd.Client, context *cli.Context, index int, databasePrefix string, databaseSuffixes string) {
	helpers.PanicErr(os.Setenv("DATABASE_URL", fmt.Sprintf("%s-%d?%s", databasePrefix, index, databaseSuffixes)))
	helpers.PanicErr(client.ResetDatabase(context))
}

func newSetupClient() *cmd.Client {
	lggr, closeLggr := logger.NewLogger()
	cfg := config.NewGeneralConfig(lggr)

	prompter := cmd.NewTerminalPrompter()
	return &cmd.Client{
		Renderer:                       cmd.RendererTable{Writer: os.Stdout},
		Config:                         cfg,
		Logger:                         lggr,
		CloseLogger:                    closeLggr,
		AppFactory:                     cmd.ChainlinkAppFactory{},
		KeyStoreAuthenticator:          cmd.TerminalKeyStoreAuthenticator{Prompter: prompter},
		FallbackAPIInitializer:         cmd.NewPromptingAPIInitializer(prompter, lggr),
		Runner:                         cmd.ChainlinkRunner{},
		PromptingSessionRequestBuilder: cmd.NewPromptingSessionRequestBuilder(prompter),
		ChangePasswordPrompter:         cmd.NewChangePasswordPrompter(),
		PasswordPrompter:               cmd.NewPasswordPrompter(),
	}
}

func requestRandomnessCallbackBatch(
	e helpers.Environment,
	consumerAddress string,
	numWords uint16,
	subID uint64,
	confDelay *big.Int,
	callbackGasLimit uint32,
	args []byte,
	batchSize *big.Int,
) (requestID *big.Int) {
	consumer := newLoadTestVRFBeaconCoordinatorConsumer(common.HexToAddress(consumerAddress), e.Ec)

	tx, err := consumer.TestRequestRandomnessFulfillmentBatch(e.Owner, subID, numWords, confDelay, callbackGasLimit, args, batchSize)
	helpers.PanicErr(err)
	receipt := helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID, "TestRequestRandomnessFulfillment")

	periodBlocks, err := consumer.IBeaconPeriodBlocks(nil)
	helpers.PanicErr(err)

	blockNumber := receipt.BlockNumber
	periodOffset := new(big.Int).Mod(blockNumber, periodBlocks)
	nextBeaconOutputHeight := new(big.Int).Sub(new(big.Int).Add(blockNumber, periodBlocks), periodOffset)

	fmt.Println("nextBeaconOutputHeight: ", nextBeaconOutputHeight)

	requestID, err = consumer.SRequestsIDs(nil, nextBeaconOutputHeight, confDelay)
	helpers.PanicErr(err)
	fmt.Println("requestID: ", requestID)

	return requestID
}
