package vrfv2plus

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/google/uuid"
	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	nodeset "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"

	"github.com/smartcontractkit/chainlink/devenv/contracts"
	"github.com/smartcontractkit/chainlink/devenv/products"
)

var L = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel).With().Fields(map[string]any{"component": "vrfv2plus"}).Logger()

const (
	txKeyPassword       = "txkey-password"
	stalenessSeconds    = uint32(86400)
	gasAfterPayment     = uint32(33285)
	fallbackLinkPerUnit = int64(1e16)
)

func (m *Configurator) GenerateNodesConfig(
	ctx context.Context,
	fs *fake.Input,
	bc []*blockchain.Input,
	ns []*nodeset.Input,
) (string, error) {
	cfg := m.Config[0]

	L.Info().Msg("Pre-generating primary EVM key for VRF node")
	encJSON, addr, err := clclient.NewETHKey(txKeyPassword)
	if err != nil {
		return "", fmt.Errorf("failed to generate primary ETH key: %w", err)
	}
	m.nodeEVMKeyAddr = addr.Hex()
	m.nodeEVMKeyEncJSON = encJSON
	m.nodeEVMKeyPass = txKeyPassword

	// Generate extra TX keys (len = cfg.NumTxKeys)
	m.txKeyAddrs = make([]string, 0, cfg.NumTxKeys)
	m.txKeyEncJSONs = make([][]byte, 0, cfg.NumTxKeys)
	for i := 0; i < cfg.NumTxKeys; i++ {
		enc, a, kErr := clclient.NewETHKey(txKeyPassword)
		if kErr != nil {
			return "", fmt.Errorf("failed to generate extra TX key %d: %w", i, kErr)
		}
		m.txKeyAddrs = append(m.txKeyAddrs, a.Hex())
		m.txKeyEncJSONs = append(m.txKeyEncJSONs, enc)
		L.Info().Str("addr", a.Hex()).Int("index", i).Msg("Generated extra TX key")
	}

	// Generate BHS key if needed
	if cfg.EnableBHSJob {
		enc, a, kErr := clclient.NewETHKey(txKeyPassword)
		if kErr != nil {
			return "", fmt.Errorf("failed to generate BHS key: %w", kErr)
		}
		m.bhsKeyAddr = a.Hex()
		m.bhsKeyEncJSON = enc
		L.Info().Str("addr", a.Hex()).Msg("Generated BHS TX key")
	}

	// Generate BHF key if needed
	if cfg.EnableBHFJob {
		enc, a, kErr := clclient.NewETHKey(txKeyPassword)
		if kErr != nil {
			return "", fmt.Errorf("failed to generate BHF key: %w", kErr)
		}
		m.bhfKeyAddr = a.Hex()
		m.bhfKeyEncJSON = enc
		L.Info().Str("addr", a.Hex()).Msg("Generated BHF TX key")
	}

	baseConfig := `[Feature]
FeedsManager = true
LogPoller = true
UICSAKeys = true

[Log]
Level = 'debug'
JSONConsole = true

[Log.File]
MaxSize = '0b'

[WebServer]
AllowOrigins = '*'
HTTPPort = 6688
SecureCookies = false
HTTPWriteTimeout = '3m'
SessionTimeout = '999h0m0s'

[WebServer.RateLimit]
Authenticated = 2000
Unauthenticated = 1000

[WebServer.TLS]
HTTPSPort = 0

[OCR2]
Enabled = true

[P2P]
[P2P.V2]
Enabled = true
ListenAddresses = ['0.0.0.0:6690']
AnnounceAddresses = ['0.0.0.0:6690']
`

	// The EVM section without KeySpecific (we append those separately)
	netConfigTemplate := `
[[EVM]]
AutoCreateKey = true
MinContractPayment = 0
BlockBackfillDepth = 100
MinIncomingConfirmations = 1

ChainID = '{{.ChainID}}'

[EVM.GasEstimator]
LimitDefault = {{.TxGasLimitDefault}}
LimitMax = {{.TxGasLimitDefault}}

[[EVM.Nodes]]
Name = 'default'
WsUrl = '{{.WsURL}}'
HttpUrl = '{{.HTTPURL}}'
{{range .Keys}}
[[EVM.KeySpecific]]
Key = '{{.}}'
GasEstimator.PriceMax = '{{$.MaxGasPriceGWei}} gwei'
{{end}}`

	tmpl, err := template.New("vrfv2plus-net-config").Parse(netConfigTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse VRF net config template: %w", err)
	}

	// Collect all keys that need a KeySpecific entry
	allKeys := append([]string{m.nodeEVMKeyAddr}, m.txKeyAddrs...)
	if m.bhsKeyAddr != "" {
		allKeys = append(allKeys, m.bhsKeyAddr)
	}
	if m.bhfKeyAddr != "" {
		allKeys = append(allKeys, m.bhfKeyAddr)
	}

	type data struct {
		ChainID           string
		WsURL             string
		HTTPURL           string
		Keys              []string
		MaxGasPriceGWei   int64
		TxGasLimitDefault uint32
	}

	txGasLimitDefault := uint32(3_500_000)
	d := data{
		ChainID:           bc[0].Out.ChainID,
		WsURL:             bc[0].Out.Nodes[0].InternalWSUrl,
		HTTPURL:           bc[0].Out.Nodes[0].InternalHTTPUrl,
		Keys:              allKeys,
		MaxGasPriceGWei:   cfg.CLNodeMaxGasPriceGWei,
		TxGasLimitDefault: txGasLimitDefault,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, d); err != nil {
		return "", fmt.Errorf("failed to execute VRF net config template: %w", err)
	}

	L.Info().Msg("VRFv2Plus nodes configuration finished")

	return baseConfig + buf.String(), nil
}

func (m *Configurator) GenerateNodesSecrets(
	_ context.Context,
	_ *fake.Input,
	bc []*blockchain.Input,
	ns []*nodeset.Input,
) (string, error) {
	L.Info().Msg("Applying VRFv2Plus CL nodes secrets")

	chainID, err := strconv.ParseInt(bc[0].Out.ChainID, 10, 64)
	if err != nil {
		return "", fmt.Errorf("failed to parse chainID: %w", err)
	}

	type evmKey struct {
		JSON     string `toml:"JSON"`
		Password string `toml:"Password"`
		ID       int64  `toml:"ID"`
	}
	type evmSecrets struct {
		Keys []evmKey `toml:"Keys"`
	}
	type secretsDoc struct {
		EVM evmSecrets `toml:"EVM"`
	}

	keys := []evmKey{
		{JSON: string(m.nodeEVMKeyEncJSON), Password: m.nodeEVMKeyPass, ID: chainID},
	}
	for _, enc := range m.txKeyEncJSONs {
		keys = append(keys, evmKey{JSON: string(enc), Password: txKeyPassword, ID: chainID})
	}
	if len(m.bhsKeyEncJSON) > 0 {
		keys = append(keys, evmKey{JSON: string(m.bhsKeyEncJSON), Password: txKeyPassword, ID: chainID})
	}
	if len(m.bhfKeyEncJSON) > 0 {
		keys = append(keys, evmKey{JSON: string(m.bhfKeyEncJSON), Password: txKeyPassword, ID: chainID})
	}

	doc := secretsDoc{EVM: evmSecrets{Keys: keys}}
	out, err := toml.Marshal(doc)
	if err != nil {
		return "", fmt.Errorf("failed to marshal node secrets: %w", err)
	}
	L.Info().Int("num_keys", len(keys)).Msg("EVM keys marshalled into secrets")

	return string(out), nil
}

func (m *Configurator) ConfigureJobsAndContracts(
	ctx context.Context,
	instanceIdx int,
	_ *fake.Input,
	bc []*blockchain.Input,
	ns []*nodeset.Input,
) error {
	cfg := m.Config[instanceIdx]

	if err := validateTopology(cfg, ns[0]); err != nil {
		return err
	}

	cl, err := clclient.New(ns[0].Out.CLNodes)
	if err != nil {
		return fmt.Errorf("failed to connect to CL nodes: %w", err)
	}

	pkey := products.NetworkPrivateKey()
	if pkey == "" {
		return errors.New("PRIVATE_KEY environment variable not set")
	}

	bcNode := bc[0].Out.Nodes[0]
	c, _, _, err := products.ETHClient(
		ctx,
		bcNode.ExternalWSUrl,
		cfg.GasSettings.FeeCapMultiplier,
		cfg.GasSettings.TipCapMultiplier,
	)
	if err != nil {
		return fmt.Errorf("could not create basic eth client: %w", err)
	}

	// Fund all pre-generated addresses
	addrsToFund := append([]string{m.nodeEVMKeyAddr}, m.txKeyAddrs...)
	if m.bhsKeyAddr != "" {
		addrsToFund = append(addrsToFund, m.bhsKeyAddr)
	}
	if m.bhfKeyAddr != "" {
		addrsToFund = append(addrsToFund, m.bhfKeyAddr)
	}
	for _, addr := range addrsToFund {
		L.Info().Str("addr", addr).Float64("eth", cfg.CLNodesFundingETH).Msg("Funding EVM address")
		if fErr := products.FundAddressEIP1559(ctx, c, pkey, addr, cfg.CLNodesFundingETH); fErr != nil {
			return fmt.Errorf("failed to fund address %s: %w", addr, fErr)
		}
	}

	chainID, err := strconv.ParseUint(bc[0].Out.ChainID, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse chainID: %w", err)
	}

	chainClient, err := products.InitSeth(bcNode.ExternalWSUrl, []string{pkey}, &chainID)
	if err != nil {
		return fmt.Errorf("failed to init seth client: %w", err)
	}

	// Deploy contracts
	L.Info().Msg("Deploying BlockhashStore")
	bhs, err := contracts.DeployBlockhashStore(chainClient)
	if err != nil {
		return fmt.Errorf("failed to deploy BlockhashStore: %w", err)
	}
	cfg.DeployedContracts.BHS = bhs.Address()

	L.Info().Msg("Deploying BatchBlockhashStore")
	batchBHS, err := contracts.DeployBatchBlockhashStore(chainClient, bhs.Address())
	if err != nil {
		return fmt.Errorf("failed to deploy BatchBlockhashStore: %w", err)
	}
	cfg.DeployedContracts.BatchBHS = batchBHS.Address()

	L.Info().Msg("Deploying VRFCoordinatorV2_5")
	coord, err := contracts.DeployVRFCoordinatorV2_5(chainClient, bhs.Address())
	if err != nil {
		return fmt.Errorf("failed to deploy VRFCoordinatorV2_5: %w", err)
	}
	cfg.DeployedContracts.Coordinator = coord.Address()

	L.Info().Msg("Deploying BatchVRFCoordinatorV2Plus")
	batchCoord, err := contracts.DeployBatchVRFCoordinatorV2Plus(chainClient, coord.Address())
	if err != nil {
		return fmt.Errorf("failed to deploy BatchVRFCoordinatorV2Plus: %w", err)
	}
	cfg.DeployedContracts.BatchCoordinator = batchCoord.Address()

	L.Info().Msg("Deploying LINK token")
	linkToken, err := contracts.DeployLinkTokenContract(L, chainClient)
	if err != nil {
		return fmt.Errorf("failed to deploy LINK token: %w", err)
	}
	cfg.DeployedContracts.LinkToken = linkToken.Address()

	L.Info().Msg("Deploying Mock LINK/ETH feed")
	mockFeed, err := contracts.DeployMockLINKETHFeed(chainClient, big.NewInt(1e18))
	if err != nil {
		return fmt.Errorf("failed to deploy MockLINKETHFeed: %w", err)
	}
	cfg.DeployedContracts.MockFeed = mockFeed.Address()

	L.Info().Msg("Setting LINK and LINK/Native feed on coordinator")
	if err := coord.SetLINKAndLINKNativeFeed(linkToken.Address(), mockFeed.Address()); err != nil {
		return fmt.Errorf("SetLINKAndLINKNativeFeed failed: %w", err)
	}

	L.Info().Msg("Setting coordinator config")
	if err := coord.SetConfig(
		cfg.MinimumConfirmations,
		cfg.MaxGasLimitCoordinator,
		stalenessSeconds,
		gasAfterPayment,
		big.NewInt(fallbackLinkPerUnit),
		cfg.FlatFeeNativePPM,
		cfg.FlatFeeLinkDiscountPPM,
		cfg.NativePremiumPercentage,
		cfg.LinkPremiumPercentage,
	); err != nil {
		return fmt.Errorf("coordinator SetConfig failed: %w", err)
	}

	// Create VRF key on node
	L.Info().Msg("Creating VRF key on CL node")
	vrfKey, err := cl[0].MustCreateVRFKey()
	if err != nil {
		return fmt.Errorf("failed to create VRF key: %w", err)
	}

	provingKey, err := contracts.EncodeOnChainVRFProvingKey(vrfKey.Data.Attributes.Uncompressed)
	if err != nil {
		return fmt.Errorf("failed to encode VRF proving key: %w", err)
	}

	L.Info().Msg("Registering VRF proving key on coordinator")
	if err := coord.RegisterProvingKey(provingKey, products.EtherToGwei(big.NewFloat(float64(cfg.CLNodeMaxGasPriceGWei))).Uint64()); err != nil {
		return fmt.Errorf("failed to register proving key: %w", err)
	}

	keyHash, err := coord.HashOfKey(ctx, provingKey)
	if err != nil {
		return fmt.Errorf("failed to get key hash: %w", err)
	}
	cfg.VRFKeyData.PubKeyCompressed = vrfKey.Data.ID
	cfg.VRFKeyData.PubKeyUncompressed = vrfKey.Data.Attributes.Uncompressed
	cfg.VRFKeyData.KeyHash = fmt.Sprintf("0x%x", keyHash)

	// Build and create VRF job
	pollPeriod, err := time.ParseDuration(cfg.VRFJobPollPeriod)
	if err != nil {
		pollPeriod = 1 * time.Second
	}
	requestTimeout, err := time.ParseDuration(cfg.VRFJobRequestTimeout)
	if err != nil {
		requestTimeout = 24 * time.Hour
	}

	fromAddresses := append([]string{m.nodeEVMKeyAddr}, m.txKeyAddrs...)

	pipelineSpec := &TxPipelineSpec{
		Address:               coord.Address(),
		EstimateGasMultiplier: 1.1,
		FromAddress:           fromAddresses[0],
	}
	if cfg.VRFJobSimulationBlock != "" {
		s := cfg.VRFJobSimulationBlock
		pipelineSpec.SimulationBlock = &s
	}
	observationSource, err := pipelineSpec.String()
	if err != nil {
		return fmt.Errorf("failed to build VRF pipeline spec: %w", err)
	}

	gasMultiplier := cfg.BatchFulfillmentGasMultiplier
	if gasMultiplier == 0 {
		gasMultiplier = 1.1
	}

	jobSpec := &JobSpec{
		Name:                          "vrf-v2-plus",
		CoordinatorAddress:            coord.Address(),
		BatchCoordinatorAddress:       batchCoord.Address(),
		PublicKey:                     vrfKey.Data.ID,
		ExternalJobID:                 uuid.New().String(),
		ObservationSource:             observationSource,
		MinIncomingConfirmations:      int(cfg.MinimumConfirmations),
		FromAddresses:                 fromAddresses,
		EVMChainID:                    bc[0].Out.ChainID,
		BatchFulfillmentEnabled:       cfg.BatchFulfillmentEnabled,
		BatchFulfillmentGasMultiplier: gasMultiplier,
		BackOffInitialDelay:           15 * time.Second,
		BackOffMaxDelay:               5 * time.Minute,
		PollPeriod:                    pollPeriod,
		RequestTimeout:                requestTimeout,
	}

	L.Info().Msg("Creating VRF job on CL node")
	job, err := cl[0].MustCreateJob(jobSpec)
	if err != nil {
		return fmt.Errorf("failed to create VRF job: %w", err)
	}
	cfg.VRFKeyData.VRFJobID = job.Data.ID

	// Create BHS job if enabled (on node 1)
	if cfg.EnableBHSJob {
		bhsJob, bhsErr := cl[1].MustCreateJob(&BlockhashStoreJobSpec{
			Name:                     "bhs-vrf-v2-plus",
			ExternalJobID:            uuid.New().String(),
			CoordinatorV2Address:     coord.Address(),
			CoordinatorV2PlusAddress: coord.Address(),
			BlockhashStoreAddress:    bhs.Address(),
			FromAddresses:            []string{m.bhsKeyAddr},
			EVMChainID:               bc[0].Out.ChainID,
			WaitBlocks:               cfg.BHSJobWaitBlocks,
			LookbackBlocks:           cfg.BHSJobLookbackBlocks,
			PollPeriod:               cfg.BHSJobPollPeriod,
			RunTimeout:               cfg.BHSJobRunTimeout,
		})
		if bhsErr != nil {
			return fmt.Errorf("failed to create BHS job: %w", bhsErr)
		}
		cfg.VRFKeyData.BHSJobID = bhsJob.Data.ID
		L.Info().Str("bhs_job_id", cfg.VRFKeyData.BHSJobID).Msg("BHS job created")
	}

	// Create BHF job if enabled (on node after BHS node)
	if cfg.EnableBHFJob {
		bhfNodeIdx := 1
		if cfg.EnableBHSJob {
			bhfNodeIdx++
		}
		bhfJob, bhfErr := cl[bhfNodeIdx].MustCreateJob(&BlockhashForwarderJobSpec{
			Name:                       "bhf-vrf-v2-plus",
			ExternalJobID:              uuid.New().String(),
			ForwardingAllowed:          false,
			CoordinatorV2Address:       coord.Address(),
			CoordinatorV2PlusAddress:   coord.Address(),
			BlockhashStoreAddress:      bhs.Address(),
			BatchBlockhashStoreAddress: batchBHS.Address(),
			FromAddresses:              []string{m.bhfKeyAddr},
			EVMChainID:                 bc[0].Out.ChainID,
			WaitBlocks:                 cfg.BHFJobWaitBlocks,
			LookbackBlocks:             cfg.BHFJobLookbackBlocks,
			PollPeriod:                 cfg.BHFJobPollPeriod,
			RunTimeout:                 cfg.BHFJobRunTimeout,
		})
		if bhfErr != nil {
			return fmt.Errorf("failed to create BHF job: %w", bhfErr)
		}
		cfg.VRFKeyData.BHFJobID = bhfJob.Data.ID
		L.Info().Str("bhf_job_id", cfg.VRFKeyData.BHFJobID).Msg("BHF job created")
	}

	// Store all TX key addresses for use in tests
	cfg.VRFKeyData.TxKeyAddresses = append([]string{m.nodeEVMKeyAddr}, m.txKeyAddrs...)

	// Set up wrapper subscription
	L.Info().Msg("Creating wrapper subscription on coordinator")
	subTx, err := coord.CreateSubscription()
	if err != nil {
		return fmt.Errorf("failed to create wrapper subscription: %w", err)
	}
	receipt, err := chainClient.Client.TransactionReceipt(ctx, subTx.Hash())
	if err != nil {
		return fmt.Errorf("failed to get CreateSubscription receipt: %w", err)
	}
	wrapperSubID, err := contracts.FindSubscriptionID(receipt)
	if err != nil {
		return fmt.Errorf("failed to parse wrapper subscription ID: %w", err)
	}
	cfg.DeployedContracts.WrapperSubID = wrapperSubID.String()

	// Deploy wrapper
	L.Info().Msg("Deploying VRFV2PlusWrapper")
	wrapper, err := contracts.DeployVRFV2PlusWrapper(chainClient,
		linkToken.Address(), mockFeed.Address(), coord.Address(), wrapperSubID)
	if err != nil {
		return fmt.Errorf("failed to deploy VRFV2PlusWrapper: %w", err)
	}
	cfg.DeployedContracts.Wrapper = wrapper.Address()

	L.Info().Msg("Adding wrapper as consumer on wrapper subscription")
	if err := coord.AddConsumer(wrapperSubID, wrapper.Address()); err != nil {
		return fmt.Errorf("failed to add wrapper as consumer: %w", err)
	}

	L.Info().Msg("Configuring wrapper")
	if err := wrapper.SetConfig(
		cfg.WrapperGasOverhead,
		cfg.CoordinatorGasOverheadNative,
		cfg.CoordinatorGasOverheadLink,
		cfg.CoordinatorGasOverheadPerWord,
		cfg.CoordinatorNativePremiumPct,
		cfg.CoordinatorLinkPremiumPct,
		keyHash,
		10, // maxNumWords
		stalenessSeconds,
		big.NewInt(fallbackLinkPerUnit),
		cfg.FlatFeeNativePPM,
		cfg.FlatFeeLinkDiscountPPM,
	); err != nil {
		return fmt.Errorf("wrapper SetConfig failed: %w", err)
	}

	wrapperNativeFund := products.EtherToWei(big.NewFloat(cfg.SubFundingAmountNative))
	L.Info().Str("amount", wrapperNativeFund.String()).Msg("Funding wrapper sub with native")
	if err := coord.FundSubscriptionWithNative(wrapperSubID, wrapperNativeFund); err != nil {
		return fmt.Errorf("failed to fund wrapper sub with native: %w", err)
	}

	wrapperLinkFund := products.EtherToWei(big.NewFloat(cfg.SubFundingAmountLink))
	encodedSubID, err := encodeSubID(wrapperSubID)
	if err != nil {
		return fmt.Errorf("failed to encode wrapper subID: %w", err)
	}
	L.Info().Str("amount", wrapperLinkFund.String()).Msg("Funding wrapper sub with LINK")
	if _, err := linkToken.TransferAndCall(coord.Address(), wrapperLinkFund, encodedSubID); err != nil {
		return fmt.Errorf("failed to fund wrapper sub with LINK: %w", err)
	}

	L.Info().Msg("Deploying VRFV2PlusWrapperLoadTestConsumer")
	wrapperConsumer, err := contracts.DeployVRFV2PlusWrapperLoadTestConsumer(chainClient, wrapper.Address())
	if err != nil {
		return fmt.Errorf("failed to deploy wrapper consumer: %w", err)
	}
	cfg.DeployedContracts.WrapperConsumer = wrapperConsumer.Address()

	consumerLinkFund := wrapperConsumerLinkFundJuels()
	L.Info().Str("amount", consumerLinkFund.String()).Msg("Funding wrapper consumer with LINK")
	if err := linkToken.Transfer(wrapperConsumer.Address(), consumerLinkFund); err != nil {
		return fmt.Errorf("failed to fund wrapper consumer with LINK: %w", err)
	}

	L.Info().Msg("Funding wrapper consumer with native ETH")
	if err := products.FundAddressEIP1559(ctx, c, pkey, wrapperConsumer.Address(), 1.0); err != nil {
		return fmt.Errorf("failed to fund wrapper consumer with native: %w", err)
	}

	L.Info().
		Str("Coordinator", cfg.DeployedContracts.Coordinator).
		Str("Wrapper", cfg.DeployedContracts.Wrapper).
		Str("WrapperConsumer", cfg.DeployedContracts.WrapperConsumer).
		Str("KeyHash", cfg.VRFKeyData.KeyHash).
		Str("VRFJobID", cfg.VRFKeyData.VRFJobID).
		Strs("TxKeyAddresses", cfg.VRFKeyData.TxKeyAddresses).
		Msg("VRFv2Plus setup complete")

	return nil
}

// validateTopology ensures the nodeset matches the product's job topology.
//
// VRFv2Plus always has exactly one VRF node. Auxiliary nodes are only valid
// when the corresponding job flag is set:
//   - node 1 (BHS node): requires EnableBHSJob = true
//
// Allowing arbitrary node counts would silently misconfigure the environment:
// a second VRF node is meaningless (only one key/job is deployed), and a BHS
// node without EnableBHSJob would sit idle and confuse debugging.
func validateTopology(cfg *VRFv2Plus, ns *nodeset.Input) error {
	got := len(ns.NodeSpecs)
	want := 1
	if cfg.EnableBHSJob {
		want++
	}
	if cfg.EnableBHFJob {
		want++
	}
	if got != want {
		return fmt.Errorf("topology mismatch: nodeset has %d node(s), want %d (enable_bhs_job=%v, enable_bhf_job=%v)",
			got, want, cfg.EnableBHSJob, cfg.EnableBHFJob)
	}
	return nil
}

// encodeSubID ABI-encodes a uint256 subscription ID for use in TransferAndCall.
func encodeSubID(subID *big.Int) ([]byte, error) {
	b := make([]byte, 32)
	subIDBytes := subID.Bytes()
	if len(subIDBytes) > 32 {
		return nil, errors.New("subID too large for uint256")
	}
	copy(b[32-len(subIDBytes):], subIDBytes)
	return b, nil
}

// wrapperConsumerLinkFundJuels returns 5 LINK in juels.
func wrapperConsumerLinkFundJuels() *big.Int {
	return new(big.Int).Mul(big.NewInt(5), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
}
