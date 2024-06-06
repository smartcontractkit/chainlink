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
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/edwards25519"
	"go.dedis.ch/kyber/v3/pairing"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-vrf/altbn_128"
	"github.com/smartcontractkit/chainlink-vrf/dkg"
	"github.com/smartcontractkit/chainlink-vrf/ocr2vrf"
	ocr2vrftypes "github.com/smartcontractkit/chainlink-vrf/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/authorized_forwarder"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	dkgContract "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ocr2vrf/generated/dkg"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ocr2vrf/generated/load_test_beacon_consumer"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ocr2vrf/generated/vrf_beacon"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ocr2vrf/generated/vrf_beacon_consumer"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ocr2vrf/generated/vrf_coordinator"

	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
)

var (
	suite              pairing.Suite = &altbn_128.PairingSuite{}
	g1                               = suite.G1()
	g2                               = suite.G2()
	tomlConfigTemplate               = `
	[P2P.V2]
	ListenAddresses = ["127.0.0.1:8000"]

	[Feature]
	LogPoller = true

	[OCR2]
	Enabled = true

	[[EVM]]
	FinalityDepth = 10
	ChainID = '%d'

	[EVM.Transactions]
	ForwardersEnabled = %t

	[EVM.HeadTracker]
	HistoryDepth = 1

	[[EVM.Nodes]]
	Name = "chain1"
	HTTPURL = "%s"
	WSURL = "%s"
	`
)

func deployDKG(e helpers.Environment) common.Address {
	_, tx, _, err := dkgContract.DeployDKG(e.Owner, e.Ec)
	helpers.PanicErr(err)
	return helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
}

func deployVRFCoordinator(e helpers.Environment, beaconPeriodBlocks *big.Int, linkAddress, linkEthFeed string) (common.Address, *vrf_coordinator.VRFCoordinator) {
	_, tx, coordinator, err := vrf_coordinator.DeployVRFCoordinator(
		e.Owner,
		e.Ec,
		beaconPeriodBlocks,
		common.HexToAddress(linkAddress),
	)
	helpers.PanicErr(err)
	return helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID), coordinator
}

func configureVRFCoordinator(e helpers.Environment, coordinator *vrf_coordinator.VRFCoordinator, maxCbGasLimit, maxCbArgsLen uint32) *gethtypes.Receipt {
	tx, err := coordinator.SetCallbackConfig(e.Owner, vrf_coordinator.VRFCoordinatorCallbackConfig{
		MaxCallbackGasLimit:        maxCbGasLimit,
		MaxCallbackArgumentsLength: maxCbArgsLen,
	})
	helpers.PanicErr(err)
	return helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID, "vrf coordinator setConfig")
}

func deployAuthorizedForwarder(e helpers.Environment, link common.Address, owner common.Address) common.Address {
	_, tx, _, err := authorized_forwarder.DeployAuthorizedForwarder(e.Owner, e.Ec, link, owner, common.Address{}, []byte{})
	helpers.PanicErr(err)
	return helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
}

func setAuthorizedSenders(e helpers.Environment, forwarder common.Address, senders []common.Address) {
	f, err := authorized_forwarder.NewAuthorizedForwarder(forwarder, e.Ec)
	helpers.PanicErr(err)
	tx, err := f.SetAuthorizedSenders(e.Owner, senders)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
}

func deployVRFBeacon(e helpers.Environment, coordinatorAddress, linkAddress, dkgAddress, keyID string) common.Address {
	keyIDBytes := decodeHexTo32ByteArray(keyID)
	_, tx, _, err := vrf_beacon.DeployVRFBeacon(e.Owner, e.Ec, common.HexToAddress(linkAddress), common.HexToAddress(coordinatorAddress), common.HexToAddress(dkgAddress), keyIDBytes)
	helpers.PanicErr(err)
	return helpers.ConfirmContractDeployed(context.Background(), e.Ec, tx, e.ChainID)
}

func deployVRFBeaconCoordinatorConsumer(e helpers.Environment, coordinatorAddress string, shouldFail bool, beaconPeriodBlocks *big.Int) common.Address {
	_, tx, _, err := vrf_beacon_consumer.DeployBeaconVRFConsumer(e.Owner, e.Ec, common.HexToAddress(coordinatorAddress), shouldFail, beaconPeriodBlocks)
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

	offchainConfig, err := dkg.OffchainConfig(encryptionKeys, signingKeys, &altbn_128.G1{}, &ocr2vrftypes.PairingTranslation{
		Suite: &altbn_128.PairingSuite{},
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

func (c *vrfBeaconSetConfigArgs) setVRFBeaconConfig(e helpers.Environment, vrfBeaconAddr string) {
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
		ocr2vrf.OffchainConfig(&c.coordinatorConfig), // off-chain config
		c.maxDurationQuery,
		c.maxDurationObservation,
		c.maxDurationReport,
		c.maxDurationAccept,
		c.maxDurationTransmit,
		int(c.f),
		onchainConfig)

	helpers.PanicErr(err)

	beacon := newVRFBeacon(common.HexToAddress(vrfBeaconAddr), e.Ec)

	tx, err := beacon.SetConfig(e.Owner, helpers.ParseAddressSlice(c.onchainPubKeys), helpers.ParseAddressSlice(c.transmitters), f, onchainConfig, offchainConfigVersion, offchainConfig)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
}

func setProducer(e helpers.Environment, vrfCoordinatorAddr, vrfBeaconAddr string) {
	coordinator := newVRFCoordinator(common.HexToAddress(vrfCoordinatorAddr), e.Ec)

	tx, err := coordinator.SetProducer(e.Owner, common.HexToAddress(vrfBeaconAddr))
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
}

func createSubscription(e helpers.Environment, vrfCoordinatorAddr string) {
	coordinator := newVRFCoordinator(common.HexToAddress(vrfCoordinatorAddr), e.Ec)

	tx, err := coordinator.CreateSubscription(e.Owner)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
}

func getSubscription(e helpers.Environment, vrfCoordinatorAddr string, subId *big.Int) vrf_coordinator.GetSubscription {
	coordinator := newVRFCoordinator(common.HexToAddress(vrfCoordinatorAddr), e.Ec)

	sub, err := coordinator.GetSubscription(nil, subId)
	helpers.PanicErr(err)
	return sub
}

// returns subscription ID that belongs to the given owner. Returns result found first
func findSubscriptionID(e helpers.Environment, vrfCoordinatorAddr string) *big.Int {
	// Use most recent 500 blocks as search window.
	head, err := e.Ec.BlockNumber(context.Background())
	helpers.PanicErr(err)
	fopts := &bind.FilterOpts{
		Start: head - 500,
	}

	coordinator := newVRFCoordinator(common.HexToAddress(vrfCoordinatorAddr), e.Ec)
	subscriptionIterator, err := coordinator.FilterSubscriptionCreated(fopts, nil, []common.Address{e.Owner.From})
	helpers.PanicErr(err)

	if !subscriptionIterator.Next() {
		helpers.PanicErr(fmt.Errorf("expected at leats 1 subID for the given owner %s", e.Owner.From.Hex()))
	}
	return subscriptionIterator.Event.SubId
}

func addConsumer(e helpers.Environment, vrfCoordinatorAddr, consumerAddr string, subId *big.Int) {
	coordinator := newVRFCoordinator(common.HexToAddress(vrfCoordinatorAddr), e.Ec)

	tx, err := coordinator.AddConsumer(e.Owner, subId, common.HexToAddress(consumerAddr))
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
}

func setPayees(e helpers.Environment, vrfBeaconAddr string, transmitters, payees []common.Address) {
	beacon := newVRFBeacon(common.HexToAddress(vrfBeaconAddr), e.Ec)

	tx, err := beacon.SetPayees(e.Owner, transmitters, payees)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
}

func eoaFundSubscription(e helpers.Environment, coordinatorAddress, linkAddress string, amount, subID *big.Int) {
	linkToken, err := link_token_interface.NewLinkToken(common.HexToAddress(linkAddress), e.Ec)
	helpers.PanicErr(err)
	bal, err := linkToken.BalanceOf(nil, e.Owner.From)
	helpers.PanicErr(err)
	fmt.Println("Initial account balance:", bal, e.Owner.From.String(), "Funding amount:", amount.String())
	b, err := utils.ABIEncode(`[{"type":"uint256"}]`, subID)
	helpers.PanicErr(err)
	tx, err := linkToken.TransferAndCall(e.Owner, common.HexToAddress(coordinatorAddress), amount, b)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID, fmt.Sprintf("sub ID: %d", subID))
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

func requestRandomness(e helpers.Environment, coordinatorAddress string, numWords uint16, subID, confDelay *big.Int) {
	coordinator := newVRFCoordinator(common.HexToAddress(coordinatorAddress), e.Ec)

	tx, err := coordinator.RequestRandomness(e.Owner, confDelay, numWords, confDelay, nil)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
}

func redeemRandomness(e helpers.Environment, coordinatorAddress string, requestID, subID *big.Int) {
	coordinator := newVRFCoordinator(common.HexToAddress(coordinatorAddress), e.Ec)

	tx, err := coordinator.RedeemRandomness(e.Owner, subID, requestID, nil)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)
}

func requestRandomnessFromConsumer(e helpers.Environment, consumerAddress string, numWords uint16, subID, confDelay *big.Int) *big.Int {
	consumer := newVRFBeaconCoordinatorConsumer(common.HexToAddress(consumerAddress), e.Ec)

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

func readRandomness(
	e helpers.Environment,
	consumerAddress string,
	requestID *big.Int,
	numWords int) {
	consumer := newVRFBeaconCoordinatorConsumer(common.HexToAddress(consumerAddress), e.Ec)
	for i := 0; i < numWords; i++ {
		r, err := consumer.SReceivedRandomnessByRequestID(nil, requestID, big.NewInt(int64(i)))
		helpers.PanicErr(err)
		fmt.Println("random word", i, ":", r.String())
	}
}

func requestRandomnessCallback(
	e helpers.Environment,
	consumerAddress string,
	numWords uint16,
	subID, confDelay *big.Int,
	callbackGasLimit uint32,
	args []byte,
) (requestID *big.Int) {
	consumer := newVRFBeaconCoordinatorConsumer(common.HexToAddress(consumerAddress), e.Ec)

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

func redeemRandomnessFromConsumer(e helpers.Environment, consumerAddress string, subID, requestID *big.Int, numWords int64) {
	consumer := newVRFBeaconCoordinatorConsumer(common.HexToAddress(consumerAddress), e.Ec)

	tx, err := consumer.TestRedeemRandomness(e.Owner, subID, requestID)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(context.Background(), e.Ec, tx, e.ChainID)

	printRandomnessFromConsumer(consumer, requestID, numWords)
}

func printRandomnessFromConsumer(consumer *vrf_beacon_consumer.BeaconVRFConsumer, requestID *big.Int, numWords int64) {
	for i := int64(0); i < numWords; i++ {
		randomness, err := consumer.SReceivedRandomnessByRequestID(nil, requestID, big.NewInt(0))
		helpers.PanicErr(err)
		fmt.Println("random words index", i, ":", randomness.String())
	}
}

func newVRFCoordinator(addr common.Address, client *ethclient.Client) *vrf_coordinator.VRFCoordinator {
	coordinator, err := vrf_coordinator.NewVRFCoordinator(addr, client)
	helpers.PanicErr(err)
	return coordinator
}

func newDKG(addr common.Address, client *ethclient.Client) *dkgContract.DKG {
	dkg, err := dkgContract.NewDKG(addr, client)
	helpers.PanicErr(err)
	return dkg
}

func newVRFBeaconCoordinatorConsumer(addr common.Address, client *ethclient.Client) *vrf_beacon_consumer.BeaconVRFConsumer {
	consumer, err := vrf_beacon_consumer.NewBeaconVRFConsumer(addr, client)
	helpers.PanicErr(err)
	return consumer
}

func newLoadTestVRFBeaconCoordinatorConsumer(addr common.Address, client *ethclient.Client) *load_test_beacon_consumer.LoadTestBeaconVRFConsumer {
	consumer, err := load_test_beacon_consumer.NewLoadTestBeaconVRFConsumer(addr, client)
	helpers.PanicErr(err)
	return consumer
}

func newVRFBeacon(addr common.Address, client *ethclient.Client) *vrf_beacon.VRFBeacon {
	beacon, err := vrf_beacon.NewVRFBeacon(addr, client)
	helpers.PanicErr(err)
	return beacon
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

func setupOCR2VRFNodeFromClient(client *cmd.Shell, context *cli.Context, e helpers.Environment) *cmd.SetupOCR2VRFNodePayload {
	payload, err := client.ConfigureOCR2VRFNode(context, e.Owner, e.Ec)
	helpers.PanicErr(err)

	return payload
}

func configureEnvironmentVariables(useForwarder bool, chainID int64, wsUrl string, ethURL string, index int, databasePrefix string, databaseSuffixes string) {
	// Set permitted envars for v2.
	helpers.PanicErr(os.Setenv("CL_DATABASE_URL", fmt.Sprintf("%s-%d?%s", databasePrefix, index, databaseSuffixes)))
	helpers.PanicErr(os.Setenv("CL_CONFIG", fmt.Sprintf(tomlConfigTemplate, chainID, useForwarder, ethURL, wsUrl)))

	// Unset prohibited envars for v2.
	helpers.PanicErr(os.Unsetenv("ETH_URL"))
	helpers.PanicErr(os.Unsetenv("ETH_HTTP_URL"))
	helpers.PanicErr(os.Unsetenv("ETH_CHAIN_ID"))
}

func resetDatabase(client *cmd.Shell, context *cli.Context) {
	helpers.PanicErr(client.ResetDatabase(context))
}

func newSetupClient() *cmd.Shell {
	prompter := cmd.NewTerminalPrompter()
	return &cmd.Shell{
		Renderer:                       cmd.RendererTable{Writer: os.Stdout},
		AppFactory:                     cmd.ChainlinkAppFactory{},
		KeyStoreAuthenticator:          cmd.TerminalKeyStoreAuthenticator{Prompter: prompter},
		FallbackAPIInitializer:         cmd.NewPromptingAPIInitializer(prompter),
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
	subID, confDelay *big.Int,
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

func printLoadtestResults(e helpers.Environment, consumerAddress string) {
	consumer := newLoadTestVRFBeaconCoordinatorConsumer(common.HexToAddress(consumerAddress), e.Ec)

	totalRequests, err := consumer.STotalRequests(nil)
	helpers.PanicErr(err)

	totalFulfilled, err := consumer.STotalFulfilled(nil)
	helpers.PanicErr(err)

	avgBlocksInMil, err := consumer.SAverageFulfillmentInMillions(nil)
	helpers.PanicErr(err)

	slowestBlocks, err := consumer.SSlowestFulfillment(nil)
	helpers.PanicErr(err)

	fastestBlock, err := consumer.SFastestFulfillment(nil)
	helpers.PanicErr(err)

	slowestRequest, err := consumer.SSlowestRequestID(nil)
	helpers.PanicErr(err)

	pendingRequests, err := consumer.PendingRequests(nil)
	helpers.PanicErr(err)

	fmt.Println("Total Requests: ", totalRequests.Uint64())
	fmt.Println("Total Fulfilled: ", totalFulfilled.Uint64())
	fmt.Println("Average Fulfillment Delay in Blocks: ", float64(avgBlocksInMil.Uint64())/1000000)
	fmt.Println("Slowest Fulfillment Delay in Blocks: ", slowestBlocks.Uint64())
	fmt.Println("Slowest Request ID: ", slowestRequest.Uint64())
	fmt.Println("Fastest Fulfillment Delay in Blocks: ", fastestBlock.Uint64())
	fmt.Println("Pending Requests: ", pendingRequests)
}
