package vrfv2

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

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	nodeset "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"

	"github.com/smartcontractkit/chainlink/devenv/contracts"
	"github.com/smartcontractkit/chainlink/devenv/products"
)

var L = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel).With().Fields(map[string]any{"component": "vrfv2"}).Logger()

const txKeyPassword = "txkey-password"

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

	m.txKeyAddrs = make([]string, 0, cfg.NumTxKeys)
	m.txKeyEncJSONs = make([][]byte, 0, cfg.NumTxKeys)
	for i := 0; i < cfg.NumTxKeys; i++ {
		enc, a, kErr := clclient.NewETHKey(txKeyPassword)
		if kErr != nil {
			return "", fmt.Errorf("failed to generate extra TX key %d: %w", i, kErr)
		}
		m.txKeyAddrs = append(m.txKeyAddrs, a.Hex())
		m.txKeyEncJSONs = append(m.txKeyEncJSONs, enc)
	}

	if cfg.EnableBHSJob {
		enc, a, kErr := clclient.NewETHKey(txKeyPassword)
		if kErr != nil {
			return "", fmt.Errorf("failed to generate BHS TX key: %w", kErr)
		}
		m.bhsKeyAddr = a.Hex()
		m.bhsKeyEncJSON = enc
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

	tmpl, err := template.New("vrfv2-net-config").Parse(netConfigTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse VRF net config template: %w", err)
	}

	allKeys := append([]string{m.nodeEVMKeyAddr}, m.txKeyAddrs...)
	if m.bhsKeyAddr != "" {
		allKeys = append(allKeys, m.bhsKeyAddr)
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

	return baseConfig + buf.String(), nil
}

func (m *Configurator) GenerateNodesSecrets(
	_ context.Context,
	_ *fake.Input,
	bc []*blockchain.Input,
	_ []*nodeset.Input,
) (string, error) {
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

	addrsToFund := append([]string{m.nodeEVMKeyAddr}, m.txKeyAddrs...)
	if m.bhsKeyAddr != "" {
		addrsToFund = append(addrsToFund, m.bhsKeyAddr)
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

	L.Info().Msg("Deploying BlockhashStore")
	bhs, err := contracts.DeployBlockhashStore(chainClient)
	if err != nil {
		return fmt.Errorf("failed to deploy BlockhashStore: %w", err)
	}
	cfg.DeployedContracts.BHS = bhs.Address()

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

	L.Info().Msg("Deploying VRFCoordinatorV2")
	coord, err := contracts.DeployVRFCoordinatorV2(chainClient, linkToken.Address(), bhs.Address(), mockFeed.Address())
	if err != nil {
		return fmt.Errorf("failed to deploy VRFCoordinatorV2: %w", err)
	}
	cfg.DeployedContracts.Coordinator = coord.Address()

	L.Info().Msg("Deploying BatchVRFCoordinatorV2")
	batchCoord, err := contracts.DeployBatchVRFCoordinatorV2(chainClient, coord.Address())
	if err != nil {
		return fmt.Errorf("failed to deploy BatchVRFCoordinatorV2: %w", err)
	}
	cfg.DeployedContracts.BatchCoordinator = batchCoord.Address()

	fallbackWei, err := contracts.FallbackWeiBigInt(cfg.FallbackWeiPerUnitLink)
	if err != nil {
		return fmt.Errorf("invalid fallback_wei_per_unit_link: %w", err)
	}

	feeConfig := vrf_coordinator_v2.VRFCoordinatorV2FeeConfig{
		FulfillmentFlatFeeLinkPPMTier1: cfg.FulfillmentFlatFeeLinkPPMTier1,
		FulfillmentFlatFeeLinkPPMTier2: cfg.FulfillmentFlatFeeLinkPPMTier2,
		FulfillmentFlatFeeLinkPPMTier3: cfg.FulfillmentFlatFeeLinkPPMTier3,
		FulfillmentFlatFeeLinkPPMTier4: cfg.FulfillmentFlatFeeLinkPPMTier4,
		FulfillmentFlatFeeLinkPPMTier5: cfg.FulfillmentFlatFeeLinkPPMTier5,
		ReqsForTier2:                   big.NewInt(cfg.ReqsForTier2),
		ReqsForTier3:                   big.NewInt(cfg.ReqsForTier3),
		ReqsForTier4:                   big.NewInt(cfg.ReqsForTier4),
		ReqsForTier5:                   big.NewInt(cfg.ReqsForTier5),
	}

	L.Info().Msg("Setting VRFCoordinatorV2 config")
	if err := coord.SetConfig(
		cfg.MinimumConfirmations,
		cfg.MaxGasLimitCoordinator,
		cfg.StalenessSeconds,
		cfg.GasAfterPaymentCalculation,
		fallbackWei,
		feeConfig,
	); err != nil {
		return fmt.Errorf("coordinator SetConfig failed: %w", err)
	}

	L.Info().Msg("Creating VRF key on CL node")
	vrfKey, err := cl[0].MustCreateVRFKey()
	if err != nil {
		return fmt.Errorf("failed to create VRF key: %w", err)
	}

	provingKey, err := contracts.EncodeOnChainVRFProvingKey(vrfKey.Data.Attributes.Uncompressed)
	if err != nil {
		return fmt.Errorf("failed to encode VRF proving key: %w", err)
	}

	rootAddr := chainClient.MustGetRootKeyAddress()
	L.Info().Str("oracle", rootAddr.Hex()).Msg("Registering VRF proving key on coordinator")
	if err := coord.RegisterProvingKey(rootAddr.Hex(), provingKey); err != nil {
		return fmt.Errorf("failed to register proving key: %w", err)
	}

	keyHash, err := coord.HashOfKey(ctx, provingKey)
	if err != nil {
		return fmt.Errorf("failed to get key hash: %w", err)
	}
	cfg.VRFKeyData.PubKeyCompressed = vrfKey.Data.ID
	cfg.VRFKeyData.PubKeyUncompressed = vrfKey.Data.Attributes.Uncompressed
	cfg.VRFKeyData.KeyHash = fmt.Sprintf("0x%x", keyHash)

	pollPeriod, err := time.ParseDuration(cfg.VRFJobPollPeriod)
	if err != nil {
		pollPeriod = time.Second
	}
	requestTimeout, err := time.ParseDuration(cfg.VRFJobRequestTimeout)
	if err != nil {
		requestTimeout = 24 * time.Hour
	}

	fromAddresses := append([]string{m.nodeEVMKeyAddr}, m.txKeyAddrs...)

	pipelineSpec := &TxPipelineSpec{
		Address:               coord.Address(),
		EstimateGasMultiplier: cfg.VRFJobEstimateGasMultiplier,
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

	batchGasMult := cfg.VRFJobBatchFulfillmentGasMultiplier
	if batchGasMult == 0 {
		batchGasMult = 1.1
	}

	jobSpec := &JobSpec{
		Name:                          "vrf-v2",
		CoordinatorAddress:            coord.Address(),
		BatchCoordinatorAddress:       batchCoord.Address(),
		PublicKey:                     vrfKey.Data.ID,
		ExternalJobID:                 uuid.New().String(),
		ObservationSource:             observationSource,
		MinIncomingConfirmations:      int(cfg.MinimumConfirmations),
		FromAddresses:                 fromAddresses,
		EVMChainID:                    bc[0].Out.ChainID,
		ForwardingAllowed:             cfg.VRFJobForwardingAllowed,
		BatchFulfillmentEnabled:       cfg.VRFJobBatchFulfillmentEnabled,
		BatchFulfillmentGasMultiplier: batchGasMult,
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

	if cfg.EnableBHSJob {
		coordinatorAddr := coord.Address()
		bhsJob, bhsErr := cl[1].MustCreateJob(&BlockhashStoreJobSpec{
			Name:                     "bhs-vrf-v2",
			ExternalJobID:            uuid.New().String(),
			CoordinatorV2Address:     coordinatorAddr,
			CoordinatorV2PlusAddress: coordinatorAddr,
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

	cfg.VRFKeyData.TxKeyAddresses = append([]string{m.nodeEVMKeyAddr}, m.txKeyAddrs...)

	L.Info().
		Str("Coordinator", cfg.DeployedContracts.Coordinator).
		Str("BatchCoordinator", cfg.DeployedContracts.BatchCoordinator).
		Str("KeyHash", cfg.VRFKeyData.KeyHash).
		Str("VRFJobID", cfg.VRFKeyData.VRFJobID).
		Strs("TxKeyAddresses", cfg.VRFKeyData.TxKeyAddresses).
		Msg("VRFv2 setup complete")

	return nil
}

func validateTopology(cfg *VRFv2, ns *nodeset.Input) error {
	got := len(ns.NodeSpecs)
	want := 1
	if cfg.EnableBHSJob {
		want++
	}
	if got != want {
		return fmt.Errorf("topology mismatch: nodeset has %d node(s), want %d (enable_bhs_job=%v)",
			got, want, cfg.EnableBHSJob)
	}
	return nil
}
