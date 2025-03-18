package integrationtest

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/sdk/freeport"
	clcommonTypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-integrations/evm/assets"
	evmtestutils "github.com/smartcontractkit/chainlink-integrations/evm/testutils"
	evmtypes "github.com/smartcontractkit/chainlink-integrations/evm/types"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/operator_factory"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/operatorforwarder/generated/authorized_forwarder"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/operatorforwarder/generated/operator"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/testhelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

// Copied from core/services/ocr2/plugins/llo/integration_test.go + svr-contracts/test/e2e/svr_test.go + integration-tests/smoke/ocr2_test.go

/*
	Steps to run:
	* `docker run --name cl-postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=dbname -p 5432:5432 -d postgres`
	* `make setup-testdb` (password is 'postgres')
	* `CL_DATABASE_URL=postgresql://chainlink_dev:insecurepassword@localhost:5432/chainlink_development_test?sslmode=disable go test -run ^TestIntegration_secondary_feed_transmission$ github.com/smartcontractkit/chainlink/v2/core/chains/evm/txm/integrationtest -v`
*/

var (
	fNodes = uint8(1)
)

func TestIntegration_secondary_feed_transmission(t *testing.T) {
	// testStartTimeStamp := time.Now()
	// multiplier := decimal.New(1, 18)
	// expirationWindow := time.Hour / time.Second

	const salt = 100

	k := big.NewInt(int64(salt))
	key := csakey.MustNewV2XXXTestingOnly(k)
	clientCSAKey := key

	contractOwner, backend := setupBlockchain(t)
	fromBlock := 1

	// Deploy link adcdress? // TODO(gg): probably not needed

	operatorFactoryAddr, _, operatorFactory, err := operator_factory.DeployOperatorFactory(contractOwner, backend.Client(), contractOwner.From) // actually: linkAddress)
	require.NoError(t, err)
	backend.Commit()
	t.Logf("Deployed OperatorFactory at %s", operatorFactoryAddr.String())

	c1 := make(chan *operator_factory.OperatorFactoryOperatorCreated)
	_, err = operatorFactory.WatchOperatorCreated(nil, c1, nil, nil, nil)
	require.NoError(t, err)

	c2 := make(chan *operator_factory.OperatorFactoryAuthorizedForwarderCreated)
	_, err = operatorFactory.WatchAuthorizedForwarderCreated(nil, c2, nil, nil, nil)
	require.NoError(t, err)

	_, err = operatorFactory.DeployNewOperatorAndForwarder(contractOwner)
	require.NoError(t, err)
	backend.Commit()
	t.Logf("Deployed Operator and Forwarder")

	var operatorInstance *operator.Operator
	select {
	case created := <-c1:
		t.Logf("Operator created at %s", created.Operator.String())
		operatorInstance, err = operator.NewOperator(created.Operator, backend.Client())
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("Timed out waiting for OperatorFactoryOperatorCreated event")
	}

	var forwarder *authorized_forwarder.AuthorizedForwarder
	select {
	case created := <-c2:
		t.Logf("Forwarder created at %s", created.Forwarder.String())
		forwarder, err = authorized_forwarder.NewAuthorizedForwarder(created.Forwarder, backend.Client()) // TODO(gg) maybe only keep address?
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("Timed out waiting for OperatorFactoryAuthorizedForwarderCreated event")
	}

	// 7. Configure the forwarder contracts

	// operator.NewOperator("address", backend.Client())
	// operatorInstance, err := contracts.LoadEthereumOperator(logger, seth, operator)
	// require.NoError(t, err, "Loading operator contract shouldn't fail")
	// forwarderInstance, err := contracts.LoadEthereumAuthorizedForwarder(seth, authorizedForwarder)
	// require.NoError(t, err, "Loading authorized forwarder contract shouldn't fail")

	// senders, err := forwarderInstance.GetAuthorizedSenders(testcontext.Get(t))
	// require.NoError(t, err, "Getting authorized senders shouldn't fail")
	// var nodesAddrs []string
	// for _, o := range nodeAddresses {
	// 	nodesAddrs = append(nodesAddrs, o.Hex())
	// }
	// require.Equal(t, nodesAddrs, senders, "Senders addresses should match node addresses")

	// owner, err := forwarderInstance.Owner(testcontext.Get(t))
	// require.NoError(t, err, "Getting authorized forwarder owner shouldn't fail")
	// require.Equal(t, operator.Hex(), owner, "Forwarder owner should match operator")

	// // actions.AcceptAuthorizedReceiversOperator(
	// // 	t, framework.L, sethClient, operators[i], forwarders[i], []common.Address{primaryAddresses[i], secondaryAddresses[i]})
	// // require.NoError(t, err, "Failed to accept authorized receivers on operator")

	// decodedTx, err := seth.Decode(tx, deployErr)
	// require.NoError(t, err, "Deploying new operator with proposed ownership with forwarder shouldn't fail")

	// for i, event := range decodedTx.Events {
	// 	require.True(t, len(event.Topics) > 0, fmt.Sprintf("Event %d should have topics", i))
	// 	switch event.Topics[0] {
	// 	case operator_factory.OperatorFactoryOperatorCreated{}.Topic().String():
	// 		if address, ok := event.EventData["operator"]; ok {
	// 			operators = append(operators, address.(common.Address))
	// 		} else {
	// 			require.Fail(t, "Operator address not found in event", event)
	// 		}
	// 	case operator_factory.OperatorFactoryAuthorizedForwarderCreated{}.Topic().String():
	// 		if address, ok := event.EventData["forwarder"]; ok {
	// 			authorizedForwarders = append(authorizedForwarders, address.(common.Address))
	// 		} else {
	// 			require.Fail(t, "Forwarder address not found in event", event)
	// 		}
	// 	}
	// }

	// donID := uint32(995544)

	// Setup the node
	port := freeport.GetOne(t)
	app, _, transmitter, kb, observedLogs := setupNode(t, port, "oracle_svr", backend, clientCSAKey, nil) // TODO(gg): fix db name?
	node := Node{app, transmitter, kb, observedLogs}
	fmt.Printf("created node with transmitter %#v\n", transmitter.String())

	// CreateTxKey creates a tx key on the Chainlink node
	// func (c *ChainlinkClient) CreateTxKey(chain string, chainId string) (*TxKey, *http.Response, error) {
	// 	txKey := &TxKey{}
	// 	framework.L.Info().Str(NodeURL, c.Config.URL).Msg("Creating Tx Key")
	// 	resp, err := c.APIClient.R().
	// 		SetPathParams(map[string]string{
	// 			"chain": chain,
	// 		}).
	// 		SetQueryParam("evmChainID", chainId).
	// 		SetResult(txKey).
	// 		Post("/v2/keys/{chain}")
	// 	if err != nil {
	// 		return nil, nil, err
	// 	}
	// 	return txKey, resp.RawResponse, err
	// }

	// set up the keys
	primaryTransmitterKey, err := node.App.GetKeyStore().Eth().Create(context.Background(), big.NewInt(int64(1337)))
	require.NoErrorf(t, err, "could not create primary transmitter key")
	secondaryTransmitterKey, err := node.App.GetKeyStore().Eth().Create(context.Background(), big.NewInt(int64(1337)))
	require.NoErrorf(t, err, "could not create secondary transmitter key")

	keys, err := node.App.GetKeyStore().Eth().GetAll(context.Background())
	require.NoError(t, err, "could not get node's eth keys")

	t.Logf("Keys are %#v", keys)

	// fund addresses

	err = fundAddressOf(primaryTransmitterKey, contractOwner, backend)
	require.NoError(t, err, "Funding primary transmitter shouldn't fail")
	backend.Commit()
	err = fundAddressOf(secondaryTransmitterKey, contractOwner, backend)
	require.NoError(t, err, "Funding secondary transmitter shouldn't fail")
	backend.Commit()
	t.Logf("Funded primary and secondary transmitter")

	_, err = operatorInstance.AcceptAuthorizedReceivers(contractOwner, []common.Address{forwarder.Address()}, []common.Address{primaryTransmitterKey.Address, secondaryTransmitterKey.Address})
	require.NoError(t, err, "Accepting authorized forwarder shouldn't fail")
	backend.Commit()
	t.Logf("Accepted authorized forwarder")

	// TODO(gg): track forwarder:
	// chainIDBigInt, ok := new(big.Int).SetString(chainID, 10)
	// require.True(t, ok, "Failed to convert chain ID to big.Int")
	// node.App.GetFeedsService().
	// _, _, err = node..TrackForwarder(chainIDBigInt, forwarders[i])

	// func (c *ChainlinkClient) TrackForwarder(chainID *big.Int, address common.Address) (*Forwarder, *http.Response, error) {
	// 	response := &Forwarder{}
	// 	request := ForwarderAttributes{
	// 		ChainID: chainID.String(),
	// 		Address: address.Hex(),
	// 	}
	// 	framework.L.Debug().Str(NodeURL, c.Config.URL).
	// 		Str("Forwarder address", (address).Hex()).
	// 		Str("Chain ID", chainID.String()).
	// 		Msg("Track forwarder")
	// 	resp, err := c.APIClient.R().
	// 		SetBody(request).
	// 		SetResult(response).
	// 		Post("/v2/nodes/evm/forwarders/track")
	// 	if err != nil {
	// 		return nil, nil, err
	// 	}
	// 	err = VerifyStatusCode(resp.StatusCode(), http.StatusCreated)
	// 	if err != nil {
	// 		return nil, nil, err
	// 	}

	// 8. Deploy dual aggregator contract
	abi, err := DualAggregatorMetaData.GetAbi()
	require.NoError(t, err, "Failed to get dual aggregator ABI")

	dualAggAddress, _, _, err := bind.DeployContract(contractOwner, *abi, common.FromHex(DualAggregatorMetaData.Bin), backend.Client(),
		common.HexToAddress(contractOwner.From.Hex()), // TODO(gg): actually linkAddress
		big.NewInt(1),                 // MinimumAnswer
		big.NewInt(50000000000000000), // MaximumAnswer
		common.Address{},              // BillingAccessController
		common.Address{},              // RequesterAccessController
		uint8(8),                      // Decimals
		"SVR test",
		common.HexToAddress("0x0000000000000000000000000000000000000000"), // secondary proxy
		uint32(30), // cutOffTime
		uint32(20), // maxSyncIterations
	)
	require.NoError(t, err, "Failed to deploy dual aggregator contract")
	backend.Commit()
	dualAggregatorInstance, err := NewDualAggregator(dualAggAddress, backend.Client())
	require.NoError(t, err, "Failed to create new dual aggregator instance")
	t.Logf("Deployed dual aggregator contract at %s", dualAggAddress.String())
	t.Logf("Primary transmitter address of node: %s", primaryTransmitterKey.EIP55Address)

	// dualAggContract, err := gethwrappers.NewDualAggregator(oevContract.Addresses[0], sethClient.Client)
	// require.NoError(t, err, "Failed to create new dual aggregator instance")
	// dualAggContracts = append(dualAggContracts, dualAggContract)
	// dualAggContractsAddresses = append(dualAggContractsAddresses, oevContract)

	// 9. Configure the dual aggregator contracts
	// S, oracleIdentities, err := getOracleIdentitiesWithKeyIndexLocal(workerNodes, 0)
	s := []int{1, 2, 2, 2}

	onchainPkBytes, err := hex.DecodeString(node.KeyBundle.OnChainPublicKey())
	require.NoError(t, err, "Failed to decode on-chain public key")

	p2pKeyId := node.App.GetConfig().P2P().PeerID()

	offchainPublicKey := node.KeyBundle.OffchainPublicKey()
	// Convert the original offchain public key (already a [32]byte) into a byte slice.
	orig := offchainPublicKey[:]

	// createVariant returns a modified version of the key by XOR-ing the last byte with tweak.
	createVariant := func(tweak byte) ocr2types.OffchainPublicKey {
		var variant [32]byte
		copy(variant[:], orig)
		variant[len(variant)-1] ^= tweak
		return variant
	}

	// createTransmitterVariant returns a modified version of the transmitter address
	// by XOR-ing its last byte with the provided tweak.
	createTransmitterVariant := func(tweak byte) string {
		addrBytes := []byte(primaryTransmitterKey.EIP55Address)
		addrBytes[len(addrBytes)-1] ^= tweak
		return string(addrBytes)
	}

	createP2PVariant := func(tweak byte) string {
		p2pBytes := []byte(p2pKeyId.Raw())
		p2pBytes[len(p2pBytes)-1] ^= tweak
		return string(p2pBytes)
	}

	oracleCurrentNode := confighelper.OracleIdentityExtra{
		OracleIdentity: confighelper.OracleIdentity{
			OnchainPublicKey:  onchainPkBytes,
			OffchainPublicKey: node.KeyBundle.OffchainPublicKey(),
			PeerID:            p2pKeyId.Raw(),
			TransmitAccount:   types.Account(primaryTransmitterKey.EIP55Address),
		},
		ConfigEncryptionPublicKey: node.KeyBundle.ConfigEncryptionPublicKey(),
	}
	t.Logf("oracleCurrentNode: %#v", oracleCurrentNode)

	mockOracle1 := confighelper.OracleIdentityExtra{
		OracleIdentity: confighelper.OracleIdentity{
			OnchainPublicKey:  onchainPkBytes,
			OffchainPublicKey: createVariant(0x02),
			PeerID:            createP2PVariant(0x02),
			TransmitAccount:   types.Account(createTransmitterVariant(0x02)),
		},
		ConfigEncryptionPublicKey: node.KeyBundle.ConfigEncryptionPublicKey(),
	}
	t.Logf("mockOracle1: %#v", mockOracle1)

	mockOracle2 := confighelper.OracleIdentityExtra{
		OracleIdentity: confighelper.OracleIdentity{
			OnchainPublicKey:  onchainPkBytes,
			OffchainPublicKey: createVariant(0x03),
			PeerID:            createP2PVariant(0x03),
			TransmitAccount:   types.Account(createTransmitterVariant(0x03)),
		},
		ConfigEncryptionPublicKey: node.KeyBundle.ConfigEncryptionPublicKey(),
	}
	t.Logf("mockOracle2: %#v", mockOracle2)

	mockOracle3 := confighelper.OracleIdentityExtra{
		OracleIdentity: confighelper.OracleIdentity{
			OnchainPublicKey:  onchainPkBytes,
			OffchainPublicKey: createVariant(0x04),
			PeerID:            createP2PVariant(0x04),
			TransmitAccount:   types.Account(createTransmitterVariant(0x04)),
		},
		ConfigEncryptionPublicKey: node.KeyBundle.ConfigEncryptionPublicKey(),
	}
	t.Logf("mockOracle3: %#v", mockOracle3)

	signerKeys, _, f, _, offchainConfigVersion, offchainConfig, err := confighelper.ContractSetConfigArgsForTests(
		30*time.Second, // deltaProgress time.Duration,
		30*time.Second, // deltaResend time.Duration,
		10*time.Second, // deltaRound time.Duration,
		20*time.Second, // deltaGrace time.Duration,
		20*time.Second, // deltaStage time.Duration,
		3,              // rMax uint8,
		s,              // s []int,
		[]confighelper.OracleIdentityExtra{oracleCurrentNode, mockOracle1, mockOracle2, mockOracle3}, // oracles []OracleIdentityExtra,
		median.OffchainConfig{
			AlphaReportInfinite: false,
			AlphaReportPPB:      1,
			AlphaAcceptInfinite: false,
			AlphaAcceptPPB:      1,
			DeltaC:              time.Minute * 30,
		}.Encode(), // reportingPluginConfig []byte,
		nil,
		5*time.Second, // maxDurationQuery time.Duration,
		5*time.Second, // maxDurationObservation time.Duration,
		5*time.Second, // maxDurationReport time.Duration,
		5*time.Second, // maxDurationShouldAcceptFinalizedReport time.Duration,
		5*time.Second, // maxDurationShouldTransmitAcceptedReport time.Duration,
		1,             // f int,
		nil,           // The median reporting plugin has an empty onchain config
	)
	require.NoError(t, err, "Failed to create contract configuration")

	// Convert signers to addresses
	signerAddresses := []common.Address{common.BytesToAddress(signerKeys[0]), common.HexToAddress("0xAD1479C185d32EB90533a08b36B3CFa5F84A0E6B"), common.HexToAddress("0xCD1479C185d32EB90533a08b36B3CFa5F84A0E6B"), common.HexToAddress("0x1D1479C185d32EB90533a08b36B3CFa5F84A0E6B")}
	require.Greater(t, len(signerAddresses), 3, "Expected more than 3 signers")

	_, err = operatorInstance.AcceptAuthorizedReceivers(contractOwner, []common.Address{forwarder.Address()}, []common.Address{primaryTransmitterKey.Address, secondaryTransmitterKey.Address})

	// Convert transmitters to addresses (needed?)
	transmitterAddresses := []common.Address{primaryTransmitterKey.EIP55Address.Address(), common.HexToAddress("0xAD1479C185d32EB90533a08b36B3CFa5F84A0E6B"), common.HexToAddress("0xCD1479C185d32EB90533a08b36B3CFa5F84A0E6B"), common.HexToAddress("0x1D1479C185d32EB90533a08b36B3CFa5F84A0E6B")}
	require.Equalf(t, len(signerAddresses), len(transmitterAddresses), "Expected the same number of signers and transmitters")

	// Replace transmitter with forwarders
	// transmitterAddresses := []common.Address{forwarder.Address()}

	onchainConfig, err := testhelpers.GenerateDefaultOCR2OnchainConfig(big.NewInt(1), big.NewInt(50000000000000000)) // MinimumAnswer MaximumAnswer
	require.NoError(t, err, "Failed to generate default ocr2 on-chain configuration")

	t.Logf("signerAddresses: %#v", signerAddresses)
	t.Logf("transmitterAddresses: %#v", transmitterAddresses)
	t.Logf("f: %d", f)
	t.Logf("ok? %t", 3*f >= uint8(len(signerAddresses)))
	_, err = dualAggregatorInstance.SetConfig(contractOwner, signerAddresses, transmitterAddresses, f, onchainConfig, offchainConfigVersion, offchainConfig)
	if err != nil {
		cerr, ok := err.(rpc.DataError)
		if !ok {
			t.Fatalf("Failed to configure dual aggregator contract: %v", err)
		}
		if cerr.ErrorData() != nil {
			t.Logf("Decoding custom ABI error from tx error")
			for k, abiError := range abi.Errors {
				data, err := hex.DecodeString(cerr.ErrorData().(string)[2:])
				if err != nil {
					t.Fatalf("Failed to decode error data: %v", err)
				}
				if len(data) < 4 {
					t.Fatalf("Error data too short: %v", data)
				}
				if bytes.Equal(data[:4], abiError.ID.Bytes()[:4]) {
					// Found a matching error
					v, err := abiError.Unpack(data)
					if err != nil {
						t.Fatalf("Failed to unpack error data: %v", err)
					}
					t.Fatalf("Failed to configure dual aggregator contract due to revert of type %s: %v", k, v)
				}
			}
		}
	}
	require.NoError(t, err, "Failed to configure dual aggregator contract")
	backend.Commit()
	backend.Commit()
	backend.Commit()
	t.Logf("Configured dual aggregator contract")

	// 3. Restart the nodes so TXMv2 can load the key for the secondary address // TODO(gg): needed?

	// 4. Fund addresses // TODO(gg): needed?

	// 5. Deploy the LINK token contract // TODO(gg): needed?

	// 6. Deploy forwarder contracts
	// var operators []common.Address
	// operators, forwarders, _ := actions.DeployForwarderContracts(
	// 	t, sethClient, common.HexToAddress(linkContract.Address()), len(workerNodes),
	// )
	// require.Equal(t, len(workerNodes), len(operators), "Number of operators does not match the number of nodes")
	// require.Equal(t, len(workerNodes), len(forwarders), "Number of authorized forwarders does not match the number of nodes")

	// offchainPublicKey, err := hex.DecodeString(strings.TrimPrefix(kb.OnChainPublicKey(), "0x"))
	// require.NoError(t, err)
	// oracles = append(oracles, confighelper.OracleIdentityExtra{
	// 	OracleIdentity: confighelper.OracleIdentity{
	// 		OnchainPublicKey:  offchainPublicKey,
	// 		TransmitAccount:   ocr2types.Account(fmt.Sprintf("%x", transmitter[:])),
	// 		OffchainPublicKey: kb.OffchainPublicKey(),
	// 		PeerID:            peerID,
	// 	},
	// 	ConfigEncryptionPublicKey: kb.ConfigEncryptionPublicKey(),
	// })

	// chainID := testutils.SimulatedChainID

	/** from Link:
		Thank you for providing the new address!
	Next, please complete the following steps to update your forwarder transmitters with the new EVM Chain Account address:
	1. Call setAuthorizedSendersOn on your operator contract:
	Using your Admin EOA, you need to submit the following values to the setAuthorizedSendersOn method on your operator contract 0xD58dd13774d6150Bc51027e3A2B46Ad5059b0EA7 on Ethereum mainnet
	https://etherscan.io/address/0xD58dd13774d6150Bc51027e3A2B46Ad5059b0EA7#writeContract#F15
	Ensure you enter all senders exactly as shown below
	targets: 0x6D53d5E35F5226a1613877e071b81217387aC6B5
	senders: 0x7663C5790E1eBf04197245d541279D13f3c2f362,0x4B2f95d9952AEd5D7Db733EF58eEdE069979f64c,0x76C07fADC35e29F0223584Fc9609Ee199b0BfC5c,0x16DBF7F4Ed84cBbA104a9305c9e614b4C20b3209
	Please reply with the tx hash. Let me know if you have questions, and thank you again.
	*/

	t.Logf("Creating feed for %s", dualAggAddress.String())
	firstKey := node.KeyBundle.Raw().Key().ID()

	// TODO(gg): put this into the actual job
	// job := fmt.Sprintf(oevJobSpec, feedNr, uuid.New().String(), contractAddress.Addresses[0].String(), "evm", firstKey, primaryAddresses[i], bootstrapPeerID.Data[0].Attributes.PeerID,
	// strings.TrimPrefix(node.DockerP2PUrl, "http://"), chainID, fromBlock, contractAddress.Addresses[0].String(), secondaryAddresses[i])

	// [relayConfig]
	// chainID = %s
	// fromBlock = %d
	// enableDualTransmission = true

	// [relayConfig.dualTransmission]
	// contractAddress = "%s"
	// transmitterAddress = "%s"

	// [relayConfig.dualTransmission.meta]
	// hint = [ "calldata" ]
	// refund = [ "0xbc1Be4cC8790b0C99cff76100E0e6d01E32C6A2C:90" ]

	// [pluginConfig]
	// juelsPerFeeCoinSource = """
	// juels_per_fee_coin [type="sum" values=<[0]>];
	// """
	// `

	pl, err := pipeline.Parse(observationSource)
	require.NoErrorf(t, err, "Failed to parse observation source")

	jb := &job.Job{
		Type:              job.OffchainReporting2,
		SchemaVersion:     1,
		Name:              null.StringFrom("SVR job 1"),
		CronSpec:          &job.CronSpec{CronSchedule: "@every 1s"},
		PipelineSpec:      &pipeline.Spec{},
		Pipeline:          *pl,
		ExternalJobID:     uuid.New(),
		ForwardingAllowed: true,
		MaxTaskDuration:   *models.NewInterval(0 * time.Second),
		OCR2OracleSpec: &job.OCR2OracleSpec{
			ContractID:           dualAggAddress.Hex(),
			Relay:                "evm",
			OCRKeyBundleID:       null.StringFrom(firstKey),
			PluginType:           clcommonTypes.Median,
			TransmitterID:        null.StringFrom(primaryTransmitterKey.Address.Hex()),
			AllowNoBootstrappers: true,
			P2PV2Bootstrappers:   []string{}, // bootstrapPeerID.Data[0].Attributes.PeerID, needed?
			RelayConfig: map[string]any{
				"chainID":                "1337",
				"fromBlock":              fromBlock,
				"enableDualTransmission": true,
				"dualTransmission": map[string]any{
					"contractAddress":    dualAggAddress.Hex(),
					"transmitterAddress": secondaryTransmitterKey.Address.Hex(),
					"meta": map[string]any{
						"hint":   []any{"calldata"},
						"refund": []any{"0xbc1Be4cC8790b0C99cff76100E0e6d01E32C6A2C:90"},
					},
				},
			},
			PluginConfig: map[string]any{
				"juelsPerFeeCoinSource": "juels_per_fee_coin [type=\"sum\" values=<[0]>]",
			},
		},
	}
	// err := helper.pipelineHelper.Jrm.CreateJob(testutils.Context(t), jb)
	err = node.App.AddJobV2(context.Background(), jb)
	require.NoError(t, err, "Failed to create feed job")

	// TODO(gg): maybe node.App.GetFeedsService().UpdateChainConfig() after setting chain config on chain
	// node.App.TxmStorageService().

	nrOfBlocks := uint64(10)
	currentBlock, err := backend.Client().BlockNumber(context.Background())
	require.NoError(t, err)

	targetBlock := big.NewInt(int64(currentBlock + nrOfBlocks))
	t.Logf("Current block is %d, waiting for %d blocks until targetBlock %d", currentBlock, nrOfBlocks, targetBlock)

	ch := make(chan *gethtypes.Header, 50)
	sub, err := backend.Client().SubscribeNewHead(context.Background(), ch)
	require.NoError(t, err)
	defer sub.Unsubscribe()

	for {
		select {
		case <-t.Context().Done():
			return
		case head := <-ch:
			t.Logf("Received block %s", head.Number.String())
			if head.Number.Cmp(targetBlock) >= 0 {
				t.Logf("Block %d has arrived, we're done", head.Number.Int64())
				return
			}
		}
	}

	/**
	NEXT STEPS

	* contract deployment + configuration:
	    logger.go:146: 2025-03-02T16:29:04.901Z	DEBUG	oracle_svr.OCR2.offchainreporting2.4afd738a-d7cd-42cd-88d9-ee960caa0e41	managed/track_config.go:46	TrackConfig: checking latestConfigDetails	{"version": "unset@unset", "jobID": 1, "jobName": "SVR job 1", "contractID": "0xbc1Be4cC8790b0C99cff76100E0e6d01E32C6A2C", "transmitterID": "0xD0203286ca243762044dc5A8636c6568b31b58A3", "evmChainID": "1337"}
	    logger.go:146: 2025-03-02T16:29:05.906Z	WARN	oracle_svr.OCR2.offchainreporting2.4afd738a-d7cd-42cd-88d9-ee960caa0e41	managed/track_config.go:110	TrackConfig: LatestConfigDetails() returned a zero configDigest. Looks like the contract has not been configured	{"version": "unset@unset", "jobID": 1, "jobName": "SVR job 1", "contractID": "0xbc1Be4cC8790b0C99cff76100E0e6d01E32C6A2C", "transmitterID": "0xD0203286ca243762044dc5A8636c6568b31b58A3", "evmChainID": "1337", "configDigest": "0000000000000000000000000000000000000000000000000000000000000000"}
	* asser on events from dual transmission (similar to svr_test
	*/

	// relayType := "evm"
	// relayConfig := fmt.Sprintf(`
	// 			chainID = "%s"
	// 			fromBlock = %d
	// 	`, chainID, fromBlock, donID)
	// addBootstrapJob(t, bootstrapNode, legacyVerifierAddr, "job-2", relayType, relayConfig)

	// 	// Channel definitions
	// 	channelDefinitions := llotypes.ChannelDefinitions{
	// 		1: {
	// 			ReportFormat: llotypes.ReportFormatEVMPremiumLegacy,
	// 			Streams: []llotypes.Stream{
	// 				{
	// 					StreamID:   ethStreamID,
	// 					Aggregator: llotypes.AggregatorMedian,
	// 				},
	// 			},
	// 			Opts: llotypes.ChannelOpts([]byte(fmt.Sprintf(`{"baseUSDFee":"0.1","expirationWindow":%d,"feedId":"0x%x","multiplier":"%s"}`, expirationWindow, quoteStreamFeedID1, multiplier.String()))),
	// 		},
	// 		2: {
	// 			ReportFormat: llotypes.ReportFormatEVMPremiumLegacy,
	// 			Streams: []llotypes.Stream{
	// 				{
	// 					StreamID:   ethStreamID,
	// 					Aggregator: llotypes.AggregatorMedian,
	// 				},
	// 			},
	// 			Opts: llotypes.ChannelOpts([]byte(fmt.Sprintf(`{"baseUSDFee":"0.1","expirationWindow":%d,"feedId":"0x%x","multiplier":"%s"}`, expirationWindow, quoteStreamFeedID2, multiplier.String()))),
	// 		},
	// 	}

	// 	pluginConfig := fmt.Sprintf(`servers = { "%s" = "%x" }
	// donID = %d
	// channelDefinitionsContractAddress = "0x%x"
	// channelDefinitionsContractFromBlock = %d`, serverURL, serverPubKey, donID, configStoreAddress, fromBlock)
	// 	addOCRJobsEVMPremiumLegacy(t, streams, serverPubKey, serverURL, legacyVerifierAddr, bootstrapPeerID, bootstrapNodePort, nodes, configStoreAddress, clientPubKeys, pluginConfig, relayType, relayConfig)

	// 	// Set config on configurator
	// 	setLegacyConfig(
	// 		t, donID, steve, backend, legacyVerifier, legacyVerifierAddr, nodes, oracles,
	// 	)

	// 	// Set config on the destination verifier
	// 	signerAddresses := make([]common.Address, len(oracles))
	// 	for i, oracle := range oracles {
	// 		signerAddresses[i] = common.BytesToAddress(oracle.OracleIdentity.OnchainPublicKey)
	// 	}
	// 	{
	// 		recipientAddressesAndWeights := []destination_verifier.CommonAddressAndWeight{}

	// 		_, err := verifier.SetConfig(steve, signerAddresses, fNodes, recipientAddressesAndWeights)
	// 		require.NoError(t, err)
	// 		backend.Commit()
	// 	}

	// 	// Expect at least one report per feed from each oracle
	// 	seen := make(map[[32]byte]map[credentials.StaticSizedPublicKey]struct{})
	// 	for _, cd := range channelDefinitions {
	// 		var opts lloevm.ReportFormatEVMPremiumLegacyOpts
	// 		err := json.Unmarshal(cd.Opts, &opts)
	// 		require.NoError(t, err)
	// 		// feedID will be deleted when all n oracles have reported
	// 		seen[opts.FeedID] = make(map[credentials.StaticSizedPublicKey]struct{}, nNodes)
	// 	}
	// 	for req := range reqs {
	// 		assert.Equal(t, uint32(llotypes.ReportFormatEVMPremiumLegacy), req.req.ReportFormat)
	// 		v := make(map[string]interface{})
	// 		err := mercury.PayloadTypes.UnpackIntoMap(v, req.req.Payload)
	// 		require.NoError(t, err)
	// 		report, exists := v["report"]
	// 		if !exists {
	// 			t.Fatalf("expected payload %#v to contain 'report'", v)
	// 		}
	// 		reportElems := make(map[string]interface{})
	// 		err = reportcodecv3.ReportTypes.UnpackIntoMap(reportElems, report.([]byte))
	// 		require.NoError(t, err)

	// 		feedID := reportElems["feedId"].([32]uint8)

	// 		if _, exists := seen[feedID]; !exists {
	// 			continue // already saw all oracles for this feed
	// 		}

	// 		var expectedBm, expectedBid, expectedAsk *big.Int
	// 		if feedID == quoteStreamFeedID1 {
	// 			expectedBm = quoteStream1.baseBenchmarkPrice.Mul(multiplier).BigInt()
	// 			expectedBid = quoteStream1.baseBid.Mul(multiplier).BigInt()
	// 			expectedAsk = quoteStream1.baseAsk.Mul(multiplier).BigInt()
	// 		} else if feedID == quoteStreamFeedID2 {
	// 			expectedBm = quoteStream2.baseBenchmarkPrice.Mul(multiplier).BigInt()
	// 			expectedBid = quoteStream2.baseBid.Mul(multiplier).BigInt()
	// 			expectedAsk = quoteStream2.baseAsk.Mul(multiplier).BigInt()
	// 		} else {
	// 			t.Fatalf("unrecognized feedID: 0x%x", feedID)
	// 		}

	// 		assert.GreaterOrEqual(t, reportElems["validFromTimestamp"].(uint32), uint32(testStartTimeStamp.Unix()))
	// 		assert.GreaterOrEqual(t, int(reportElems["observationsTimestamp"].(uint32)), int(testStartTimeStamp.Unix()))
	// 		assert.Equal(t, "33597747607000", reportElems["nativeFee"].(*big.Int).String())
	// 		assert.Equal(t, "7547169811320755", reportElems["linkFee"].(*big.Int).String())
	// 		assert.Equal(t, reportElems["observationsTimestamp"].(uint32)+uint32(expirationWindow), reportElems["expiresAt"].(uint32))
	// 		assert.Equal(t, expectedBm.String(), reportElems["benchmarkPrice"].(*big.Int).String())
	// 		assert.Equal(t, expectedBid.String(), reportElems["bid"].(*big.Int).String())
	// 		assert.Equal(t, expectedAsk.String(), reportElems["ask"].(*big.Int).String())

	// 		// emulate mercury server verifying report (local verification)
	// 		{
	// 			rv := mercuryverifier.NewVerifier()

	// 			reportSigners, err := rv.Verify(mercuryverifier.SignedReport{
	// 				RawRs:         v["rawRs"].([][32]byte),
	// 				RawSs:         v["rawSs"].([][32]byte),
	// 				RawVs:         v["rawVs"].([32]byte),
	// 				ReportContext: v["reportContext"].([3][32]byte),
	// 				Report:        v["report"].([]byte),
	// 			}, fNodes, signerAddresses)
	// 			require.NoError(t, err)
	// 			assert.GreaterOrEqual(t, len(reportSigners), int(fNodes+1))
	// 			assert.Subset(t, signerAddresses, reportSigners)
	// 		}

	// 	}

}

func setupBlockchain(t *testing.T) (*bind.TransactOpts, evmtypes.Backend) {
	// TODO(gg): maybe use seth instead?

	contractOwner := evmtestutils.MustNewSimTransactor(t) // config contract deployer and owner
	genesisData := gethtypes.GenesisAlloc{contractOwner.From: {Balance: assets.Ether(1000).ToInt()}}
	backend := cltest.NewSimulatedBackend(t, genesisData, ethconfig.Defaults.Miner.GasCeil)
	backend.Commit()
	backend.Commit() // ensure starting block number at least 1

	// // Configurator
	// configuratorAddress, _, configurator, err := configurator.DeployConfigurator(transactor, backend.Client())
	// require.NoError(t, err)
	// backend.Commit()
	// ChannelConfigStore

	return contractOwner, backend
}

func mustNewType(t string) abi.Type {
	result, err := abi.NewType(t, "", []abi.ArgumentMarshaling{})
	if err != nil {
		panic(fmt.Sprintf("Unexpected error during abi.NewType: %s", err))
	}
	return result
}

func fundAddressOf(key ethkey.KeyV2, contractOwner *bind.TransactOpts, backend evmtypes.Backend) error {

	// backend.Client().SendTransaction()
	// contractOwner.From
	// backend.Client().
	// 	// 4. Fund addresses
	// 	for i := range primaryAddresses {
	// 		require.NoError(t, ns.SendETH(sethClient.Client, pkey, primaryAddresses[i].String(), big.NewFloat(0.2)), "Failed to fund primary address")
	// 		require.NoError(t, ns.SendETH(sethClient.Client, pkey, secondaryAddresses[i].String(), big.NewFloat(0.2)), "Failed to fund secondary address")
	// 	}

	// privateKey, err := crypto.HexToECDSA(key)
	// if err != nil {
	// 	return er.Wrap(err, "failed to parse private key")
	// }
	wei := new(big.Int)
	amount := big.NewFloat(0.2)
	amountWei := new(big.Float).Mul(amount, big.NewFloat(1e18))
	amountWei.Int(wei)

	// publicKey := privateKey.Public()
	// publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	// if !ok {
	// 	return fmt.Errorf("error casting public key to ECDSA")
	// }
	// fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	backend.Client().PendingNonceAt(context.Background(), key.Address)
	nonce, err := backend.Client().PendingNonceAt(context.Background(), contractOwner.From)
	if err != nil {
		return fmt.Errorf("failed to fetch nonce: %w", err)
	}

	gasPrice, err := backend.Client().SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("failed to fetch gas price: %w", err)
	}
	gasLimit := uint64(21000) // Standard gas limit for ETH transfer

	tx := gethtypes.NewTransaction(nonce, key.Address, wei, gasLimit, gasPrice, nil)

	signedTx, err := contractOwner.Signer(contractOwner.From, tx)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %w", err)
	}

	err = backend.Client().SendTransaction(context.Background(), signedTx)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	backend.Commit()

	_, err = bind.WaitMined(context.Background(), backend.Client(), signedTx)
	return err
}

var observationSource = `
//randomness
   val1 [type="memo" value="10"]
   val2 [type="memo" value="20"]
   val3 [type="memo" value="30"]
   val4 [type="memo" value="40"]
   val5 [type="memo" value="50"]
   val6 [type="memo" value="60"]
   val7 [type="memo" value="70"]
   val8 [type="memo" value="80"]
   val9 [type="memo" value="90"]

   random1 [type="any"]
   random2 [type="any"]
   random3 [type="any"]

   val1 -> random1
   val2 -> random2
   val3 -> random3
   val4 -> random1
   val5 -> random2
   val6 -> random3
   val7 -> random1
   val8 -> random2
   val9 -> random3


   // data source 1
   ds1_multiply [type="multiply" times=100]

	// data source 2
   ds2_multiply [type="multiply" times=100]


   // data source 3
   ds3_multiply [type="multiply" times=100]


   random1 -> ds1_multiply -> answer
   random2 -> ds2_multiply -> answer
   random3 -> ds3_multiply -> answer

   answer [type=median]
`

var oevJobSpec = `
type = "offchainreporting2"
schemaVersion = 1
name = "OEV job %d"
externalJobID = "%s"
forwardingAllowed = true
maxTaskDuration = "0s"
contractID = "%s"
relay = "%s"
ocrKeyBundleID = "%s"
pluginType = "median"
transmitterID = "%s"
p2pv2Bootstrappers = ["%s@%s"]

observationSource = """
 //randomness
    val1 [type="memo" value="10"]
    val2 [type="memo" value="20"]
    val3 [type="memo" value="30"]
    val4 [type="memo" value="40"]
    val5 [type="memo" value="50"]
    val6 [type="memo" value="60"]
    val7 [type="memo" value="70"]
    val8 [type="memo" value="80"]
    val9 [type="memo" value="90"]

    random1 [type="any"]
    random2 [type="any"]
    random3 [type="any"]

    val1 -> random1
    val2 -> random2
    val3 -> random3
    val4 -> random1
    val5 -> random2
    val6 -> random3
    val7 -> random1
    val8 -> random2
    val9 -> random3


    // data source 1
    ds1_multiply [type="multiply" times=100]

     // data source 2
    ds2_multiply [type="multiply" times=100]


    // data source 3
    ds3_multiply [type="multiply" times=100]


    random1 -> ds1_multiply -> answer
    random2 -> ds2_multiply -> answer
    random3 -> ds3_multiply -> answer

    answer [type=median]
"""

[relayConfig]
chainID = %s
fromBlock = %d
enableDualTransmission = true

[relayConfig.dualTransmission]
contractAddress = "%s"
transmitterAddress = "%s"

[relayConfig.dualTransmission.meta]
hint = [ "calldata" ]
refund = [ "0xbc1Be4cC8790b0C99cff76100E0e6d01E32C6A2C:90" ]

[pluginConfig]
juelsPerFeeCoinSource = """
juels_per_fee_coin [type="sum" values=<[0]>];
"""
`
