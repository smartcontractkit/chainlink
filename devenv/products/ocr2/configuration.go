package ocr2

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	"github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"golang.org/x/sync/errgroup"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/link_token"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	"github.com/smartcontractkit/chainlink/devenv/products"

	nodeset "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
)

const (
	ConfigureNodesNetwork ConfigPhase = iota
	ConfigureProductContractsJobs
)

var L = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel).With().Fields(map[string]any{"component": "ocr2"}).Logger()

type OCR2 struct {
	OCR2                     *OCRv2OffChainOptions  `toml:"ocr2"`
	OCR2SetConfig            *OCRv2SetConfigOptions `toml:"ocr2_set_config"`
	OCR2SetConfigOut         *OCRv2Config           `toml:"ocr2_set_config_out"`
	OCR2MedianOffchainConfig *MedianOffchainConfig  `toml:"ocr2_median_offchain_config"`
	EAFake                   *EAFake                `toml:"ea_fake"`
	Jobs                     *Jobs                  `toml:"jobs"`
	LinkContractAddress      string                 `toml:"link_contract_address"`
	CLNodesFundingETH        float64                `toml:"cl_nodes_funding_eth"`
	CLNodesFundingLink       float64                `toml:"cl_nodes_funding_link"`
	ChainFinalityDepth       int64                  `toml:"chain_finality_depth"`
	VerificationTimeoutSec   int64                  `toml:"verification_timeout_sec"`
	GasSettings              *GasSettings           `toml:"gas_settings"`
	DeployedContracts        *DeployedContracts     `toml:"deployed_contracts"`
}

type DeployedContracts struct {
	OCRv2AggregatorAddr string `toml:"ocr2_aggregator_address"`
}

type GasSettings struct {
	FeeCapMultiplier int64 `toml:"fee_cap_multiplier"`
	TipCapMultiplier int64 `toml:"tip_cap_multiplier"`
}

type MedianOffchainConfig struct {
	AlphaReportPPB      uint64 `toml:"alpha_report_ppb"`
	AlphaAcceptPPB      uint64 `toml:"alpha_accept_ppb"`
	DeltaCSec           int64  `toml:"delta_sec"`
	AlphaReportInfinite bool   `toml:"alpha_report_infinite"`
	AlphaAcceptInfinite bool   `toml:"alpha_accept_infinite"`
}

type Jobs struct {
	MaxTaskDurationSec int64 `toml:"max_task_duration_sec"`
}

type EAFake struct {
	MinValue         int64 `toml:"min_value"`
	MaxValue         int64 `toml:"max_value"`
	ChangesPerMinute int64 `toml:"changes_per_minute"`
}

type ConfigPhase int

type OCRv2OffChainOptions struct {
	MinimumAnswer             *big.Int       `toml:"minimum_answer"`
	MaximumAnswer             *big.Int       `toml:"maximum_answer"`
	Description               string         `toml:"description"`
	MaximumGasPrice           uint32         `toml:"maximum_gas_price"`
	ReasonableGasPrice        uint32         `toml:"reasonable_gas_price"`
	MicroLinkPerEth           uint32         `toml:"micro_link_per_eth"`
	LinkGweiPerObservation    uint32         `toml:"link_gwei_per_observation"`
	LinkGweiPerTransmission   uint32         `toml:"link_gwei_per_transmission"`
	BillingAccessController   common.Address `toml:"billing_access_controller_addr"`
	RequesterAccessController common.Address `toml:"requester_access_controller_addr"`
	Decimals                  uint8          `toml:"decimals"`
}

type OCRv2SetConfigOptions struct {
	RMax                                    uint8         `toml:"r_max"`
	DeltaProgress                           time.Duration `toml:"delta_progress_sec"`
	DeltaResend                             time.Duration `toml:"delta_resend_sec"`
	DeltaRound                              time.Duration `toml:"delta_round_sec"`
	DeltaGrace                              time.Duration `toml:"delta_grace_sec"`
	DeltaStage                              time.Duration `toml:"delta_stage_sec"`
	MaxDurationInitialization               time.Duration `toml:"max_duration_initialization_sec"`
	MaxDurationQuery                        time.Duration `toml:"max_duration_query_sec"`
	MaxDurationObservation                  time.Duration `toml:"max_duration_observation_sec"`
	MaxDurationReport                       time.Duration `toml:"max_duration_report_sec"`
	MaxDurationShouldAcceptFinalizedReport  time.Duration `toml:"max_duration_should_accept_finalized_report_sec"`
	MaxDurationShouldTransmitAcceptedReport time.Duration `toml:"max_duration_should_transmit_accepted_report_sec"`
}

type OCRv2Config struct {
	Signers               []common.Address
	Transmitters          []common.Address
	OnchainConfig         []byte
	OffchainConfig        []byte
	OffchainConfigVersion uint64
	F                     uint8
}

type Configurator struct {
	OCR2 *OCR2 `toml:"ocr2"`
}

func NewOCR2Configurator() *Configurator {
	return &Configurator{}
}

func (m *Configurator) Load() error {
	cfg, err := products.Load[Configurator]()
	if err != nil {
		return fmt.Errorf("failed to load product config: %w", err)
	}
	m.OCR2 = cfg.OCR2
	return nil
}

func (m *Configurator) Store(path string) error {
	if err := products.Store(".", m); err != nil {
		return fmt.Errorf("failed to store product config: %w", err)
	}
	return nil
}

func (m *Configurator) GenerateCLNodesBlockchainConfig(ctx context.Context, bc *blockchain.Input) (string, error) {
	L.Info().Msg("Applying default CL nodes configuration")
	// configure node set and generate CL nodes configs
	node := bc.Out.Nodes[0]
	chainID := bc.Out.ChainID
	netConfig := fmt.Sprintf(`
       [[EVM]]
       LogPollInterval = '1s'
       BlockBackfillDepth = 100
       LinkContractAddress = '%s'
       ChainID = '%s'
       MinIncomingConfirmations = 1
       MinContractPayment = '0.0000001 link'
       FinalityDepth = %d

       [[EVM.Nodes]]
       Name = 'default'
       WsUrl = '%s'
       HttpUrl = '%s'

       [Feature]
       FeedsManager = true
       LogPoller = true
       UICSAKeys = true
       [OCR2]
       Enabled = true
       SimulateTransactions = false
       DefaultTransactionQueueDepth = 1
       [P2P.V2]
       Enabled = true
       ListenAddresses = ['0.0.0.0:6690']

   	   [Log]
   JSONConsole = true
   Level = 'debug'
   [Pyroscope]
   ServerAddress = 'http://host.docker.internal:4040'
   Environment = 'local'
   [WebServer]
          SessionTimeout = '999h0m0s'
          HTTPWriteTimeout = '3m'
   SecureCookies = false
   HTTPPort = 6688
   [WebServer.TLS]
   HTTPSPort = 0
       [WebServer.RateLimit]
       Authenticated = 5000
       Unauthenticated = 5000
   [JobPipeline]
   [JobPipeline.HTTPRequest]
   DefaultTimeout = '1m'
       [Log.File]
       MaxSize = '0b'
`, m.OCR2.LinkContractAddress,
		chainID,
		m.OCR2.ChainFinalityDepth,
		node.InternalWSUrl,
		node.InternalHTTPUrl,
	)
	L.Info().Msg("Nodes network configuration is finished")
	return netConfig, nil
}

func (m *Configurator) ConfigureJobsAndContracts(
	ctx context.Context,
	fake *fake.Input,
	bc *blockchain.Input,
	ns *nodeset.Input,
) error {
	L.Info().Msg("Connecting to CL nodes")
	cl, err := clclient.New(ns.Out.CLNodes)
	if err != nil {
		return err
	}
	pkey := getNetworkPrivateKey()
	if pkey == "" {
		return errors.New("PRIVATE_KEY environment variable not set")
	}

	transmitters := make([]common.Address, 0)
	ethKeyAddresses := make([]string, 0)
	for i, nc := range cl {
		addr, cErr := nc.ReadPrimaryETHKey(bc.Out.ChainID)
		if cErr != nil {
			return cErr
		}
		ethKeyAddresses = append(ethKeyAddresses, addr.Attributes.Address)
		transmitters = append(transmitters, common.HexToAddress(addr.Attributes.Address))
		L.Info().
			Int("Idx", i).
			Str("ETH", addr.Attributes.Address).
			Msg("Node info")
	}
	bcNode := bc.Out.Nodes[0]
	c, auth, rootAddr, err := ETHClient(
		ctx,
		bcNode.ExternalWSUrl,
		m.OCR2.GasSettings.FeeCapMultiplier,
		m.OCR2.GasSettings.TipCapMultiplier,
	)
	if err != nil {
		return fmt.Errorf("could not create basic eth client: %w", err)
	}
	for _, addr := range ethKeyAddresses {
		if cErr := FundNodeEIP1559(ctx, c, pkey, addr, m.OCR2.CLNodesFundingETH); cErr != nil {
			return cErr
		}
	}
	ocrv2Config, ocr2Addr, err := m.configureContracts(
		ctx,
		c,
		auth,
		cl,
		rootAddr,
		transmitters,
		m.OCR2.CLNodesFundingLink,
	)
	if err != nil {
		return err
	}
	m.OCR2.OCR2SetConfigOut = ocrv2Config
	if cErr := m.configureJobs(ctx, fake, bc, ns, cl, ocr2Addr); cErr != nil {
		return cErr
	}
	r := resty.New().SetBaseURL(fake.Out.BaseURLHost)

	_, err = r.R().Post(`/trigger_deviation?result=200`)
	if err != nil {
		return fmt.Errorf("could not set ea fake values: %w", err)
	}
	L.Info().
		Msg("Setting fake external adapter (data feed) values")
	m.OCR2.DeployedContracts = &DeployedContracts{OCRv2AggregatorAddr: ocr2Addr}
	return nil
}

// deployLinkAndMint is a universal action that deploys link token and mints required amount of LINK token for all the nodes.
func deployLinkAndMint(ctx context.Context, c *ethclient.Client, auth *bind.TransactOpts, rootAddr string, transmitters []common.Address, linkFunding float64) (*link_token.LinkToken, error) {
	addr, tx, lt, err := link_token.DeployLinkToken(auth, c)
	if err != nil {
		return nil, fmt.Errorf("could not create link token contract: %w", err)
	}
	_, err = bind.WaitDeployed(ctx, c, tx)
	if err != nil {
		return nil, err
	}
	L.Info().Str("Address", addr.Hex()).Msg("Deployed link token contract")
	tx, err = lt.GrantMintRole(auth, common.HexToAddress(rootAddr))
	if err != nil {
		return nil, fmt.Errorf("could not grant mint role: %w", err)
	}
	_, err = bind.WaitMined(ctx, c, tx)
	if err != nil {
		return nil, err
	}
	// mint for public keys of nodes directly instead of transferring
	for _, transmitter := range transmitters {
		amount := new(big.Float).Mul(big.NewFloat(linkFunding), big.NewFloat(1e18))
		amountWei, _ := amount.Int(nil)
		L.Info().Msgf("Minting LINK for transmitter address: %s", transmitter.Hex())
		tx, err = lt.Mint(auth, transmitter, amountWei)
		if err != nil {
			return nil, fmt.Errorf("could not transfer link token contract: %w", err)
		}
		_, err = bind.WaitMined(ctx, c, tx)
		if err != nil {
			return nil, err
		}
	}
	return lt, nil
}

func UpdateOCR2ConfigOffChainValues(ctx context.Context, bc *blockchain.Input, o *OCR2, ocr2i *ocr2aggregator.OCR2Aggregator, cl []*clclient.ChainlinkClient, o2 *OCRv2SetConfigOptions) error {
	if o2 == nil {
		return nil
	}
	c, auth, _, err := ETHClient(
		ctx,
		bc.Out.Nodes[0].ExternalHTTPUrl,
		o.GasSettings.FeeCapMultiplier,
		o.GasSettings.TipCapMultiplier,
	)
	if err != nil {
		return fmt.Errorf("could not create basic eth client: %w", err)
	}
	// generating oracle identities and setting up OCRv2
	s, ids, err := getOracleIdentities(cl)
	if err != nil {
		return fmt.Errorf("could not get oracle identities: %w", err)
	}
	signerKeys, transmitterAccounts, f, _, offchainConfigVersion, offchainConfig, err := confighelper.ContractSetConfigArgsForTests(
		o2.DeltaProgress,
		o2.DeltaResend,
		o2.DeltaRound,
		o2.DeltaGrace,
		o2.DeltaStage,
		o2.RMax,
		s,
		ids,
		median.OffchainConfig{
			AlphaAcceptInfinite: o.OCR2MedianOffchainConfig.AlphaAcceptInfinite,
			AlphaReportInfinite: o.OCR2MedianOffchainConfig.AlphaReportInfinite,
			AlphaReportPPB:      o.OCR2MedianOffchainConfig.AlphaReportPPB,
			AlphaAcceptPPB:      o.OCR2MedianOffchainConfig.AlphaAcceptPPB,
			DeltaC:              time.Duration(o.OCR2MedianOffchainConfig.DeltaCSec) * time.Second,
		}.Encode(),
		nil,
		o2.MaxDurationQuery,
		o2.MaxDurationObservation,
		o2.MaxDurationReport,
		o2.MaxDurationShouldAcceptFinalizedReport,
		o2.MaxDurationShouldTransmitAcceptedReport,
		1,
		nil, // The median reporting plugin has an empty onchain config
	)
	if err != nil {
		return fmt.Errorf("could not set config: %w", err)
	}
	signerAddresses := make([]common.Address, 0)
	for _, signer := range signerKeys {
		signerAddresses = append(signerAddresses, common.BytesToAddress(signer))
	}
	transmitterAddresses := make([]common.Address, 0)
	for _, account := range transmitterAccounts {
		transmitterAddresses = append(transmitterAddresses, common.HexToAddress(string(account)))
	}
	onChainConfig, err := median.StandardOnchainConfigCodec{}.Encode(context.Background(), median.OnchainConfig{Min: o.OCR2.MinimumAnswer, Max: o.OCR2.MaximumAnswer})
	if err != nil {
		return fmt.Errorf("could not encode onchain config: %w", err)
	}
	tx, err := ocr2i.SetConfig(auth, signerAddresses, transmitterAddresses, f, onChainConfig, offchainConfigVersion, offchainConfig)
	if err != nil {
		return fmt.Errorf("could not set OCRv2 config: %w", err)
	}
	_, err = bind.WaitMined(ctx, c, tx)
	if err != nil {
		return err
	}
	return nil
}

func (m *Configurator) configureContracts(ctx context.Context, c *ethclient.Client, auth *bind.TransactOpts, cl []*clclient.ChainlinkClient, rootAddr string, transmitters []common.Address, linkFunding float64) (*OCRv2Config, string, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()
	L.Info().Msg("Deploying LINK token contract")
	lt, err := deployLinkAndMint(ctx, c, auth, rootAddr, transmitters, linkFunding)
	if err != nil {
		return nil, "", fmt.Errorf("could not create link token contract and mint: %w", err)
	}
	// OCRv2 Aggregator
	L.Info().Msg("Deploying OCRv2 aggregator contract")
	opts := m.OCR2.OCR2
	ocr2addr, tx, ocr2i, err := ocr2aggregator.DeployOCR2Aggregator(auth, c, lt.Address(), opts.MinimumAnswer, opts.MaximumAnswer, common.HexToAddress(""), common.HexToAddress(""), 18, "")
	if err != nil {
		return nil, "", fmt.Errorf("could not create ocr2 aggregator contract: %w", err)
	}
	_, err = bind.WaitDeployed(ctx, c, tx)
	if err != nil {
		return nil, "", err
	}
	L.Info().Str("Address", ocr2addr.String()).Msg("Deployed OCRv2 Aggregator contract")
	tx, err = ocr2i.SetPayees(auth, transmitters, []common.Address{
		common.HexToAddress(rootAddr),
		common.HexToAddress(rootAddr),
		common.HexToAddress(rootAddr),
		common.HexToAddress(rootAddr),
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to set payees: %w", err)
	}
	_, err = bind.WaitMined(ctx, c, tx)
	if err != nil {
		return nil, "", err
	}
	// generating oracle identities and setting up OCRv2
	s, ids, err := getOracleIdentities(cl)
	if err != nil {
		return nil, "", fmt.Errorf("could not get oracle identities: %w", err)
	}
	ocrSetConfig := m.OCR2.OCR2SetConfig
	signerKeys, transmitterAccounts, f, _, offchainConfigVersion, offchainConfig, err := confighelper.ContractSetConfigArgsForTests(
		ocrSetConfig.DeltaProgress*time.Second,
		ocrSetConfig.DeltaResend*time.Second,
		ocrSetConfig.DeltaRound*time.Second,
		ocrSetConfig.DeltaGrace*time.Second,
		ocrSetConfig.DeltaStage*time.Second,
		ocrSetConfig.RMax,
		s,
		ids,
		median.OffchainConfig{
			AlphaAcceptInfinite: m.OCR2.OCR2MedianOffchainConfig.AlphaAcceptInfinite,
			AlphaReportInfinite: m.OCR2.OCR2MedianOffchainConfig.AlphaReportInfinite,
			AlphaReportPPB:      m.OCR2.OCR2MedianOffchainConfig.AlphaReportPPB,
			AlphaAcceptPPB:      m.OCR2.OCR2MedianOffchainConfig.AlphaAcceptPPB,
			DeltaC:              time.Duration(m.OCR2.OCR2MedianOffchainConfig.DeltaCSec) * time.Second,
		}.Encode(),
		nil,
		ocrSetConfig.MaxDurationQuery*time.Second,
		ocrSetConfig.MaxDurationObservation*time.Second,
		ocrSetConfig.MaxDurationReport*time.Second,
		ocrSetConfig.MaxDurationShouldAcceptFinalizedReport*time.Second,
		ocrSetConfig.MaxDurationShouldTransmitAcceptedReport*time.Second,
		1,
		nil, // The median reporting plugin has an empty onchain config
	)
	if err != nil {
		return nil, "", fmt.Errorf("could not set config: %w", err)
	}
	signerAddresses := make([]common.Address, 0)
	for _, signer := range signerKeys {
		signerAddresses = append(signerAddresses, common.BytesToAddress(signer))
	}
	transmitterAddresses := make([]common.Address, 0)
	for _, account := range transmitterAccounts {
		transmitterAddresses = append(transmitterAddresses, common.HexToAddress(string(account)))
	}
	onChainConfig, err := median.StandardOnchainConfigCodec{}.Encode(context.Background(), median.OnchainConfig{Min: m.OCR2.OCR2.MinimumAnswer, Max: m.OCR2.OCR2.MaximumAnswer})
	if err != nil {
		return nil, "", fmt.Errorf("could not encode onchain config: %w", err)
	}
	tx, err = ocr2i.SetConfig(auth, signerAddresses, transmitterAddresses, f, onChainConfig, offchainConfigVersion, offchainConfig)
	if err != nil {
		return nil, "", fmt.Errorf("could not set OCRv2 config: %w", err)
	}
	_, err = bind.WaitMined(ctx, c, tx)
	if err != nil {
		return nil, "", err
	}
	return &OCRv2Config{
		F:                     f,
		Signers:               signerAddresses,
		Transmitters:          transmitterAddresses,
		OnchainConfig:         onChainConfig,
		OffchainConfigVersion: offchainConfigVersion,
		OffchainConfig:        offchainConfig,
	}, ocr2addr.String(), err
}

func getOracleIdentities(clClients []*clclient.ChainlinkClient) ([]int, []confighelper.OracleIdentityExtra, error) {
	s := make([]int, len(clClients))
	oracleIdentities := make([]confighelper.OracleIdentityExtra, len(clClients))
	sharedSecretEncryptionPublicKeys := make([]types.ConfigEncryptionPublicKey, len(clClients))
	eg := &errgroup.Group{}
	for i, cl := range clClients {
		eg.Go(func() error {
			addresses, err := cl.EthAddresses()
			if err != nil {
				return err
			}
			ocr2Keys, err := cl.MustReadOCR2Keys()
			if err != nil {
				return err
			}
			var ocr2Config clclient.OCR2KeyAttributes
			for _, key := range ocr2Keys.Data {
				if key.Attributes.ChainType == "evm" {
					ocr2Config = key.Attributes
					break
				}
			}

			keys, err := cl.MustReadP2PKeys()
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
			sharedSecretEncryptionPublicKeys[i] = configPkBytesFixed
			oracleIdentities[i] = confighelper.OracleIdentityExtra{
				OracleIdentity: confighelper.OracleIdentity{
					OnchainPublicKey:  onchainPkBytes,
					OffchainPublicKey: offchainPkBytesFixed,
					PeerID:            p2pKeyID,
					TransmitAccount:   types.Account(addresses[0]),
				},
				ConfigEncryptionPublicKey: configPkBytesFixed,
			}
			s[i] = 1
			L.Trace().
				Interface("OnChainPK", onchainPkBytes).
				Interface("OffChainPK", offchainPkBytesFixed).
				Interface("ConfigPK", configPkBytesFixed).
				Str("PeerID", p2pKeyID).
				Str("Address", addresses[0]).
				Msg("Oracle identity")
			return nil
		})
	}
	return s, oracleIdentities, eg.Wait()
}

func (m *Configurator) configureJobs(ctx context.Context, fake *fake.Input, bc *blockchain.Input, ns *nodeset.Input, clNodes []*clclient.ChainlinkClient, ocr2Addr string) error {
	bootstrapNode := clNodes[0]
	workerNodes := clNodes[1:]
	bootstrapP2PIds, err := bootstrapNode.MustReadP2PKeys()
	if err != nil {
		return err
	}
	p2pV2Bootstrapper := fmt.Sprintf("%s@%s:%d", bootstrapP2PIds.Data[0].Attributes.PeerID, ns.Out.CLNodes[0].Node.ContainerName, 6690)
	// Set the value for the jobs to report on
	bootstrapSpec := &TaskJobSpec{
		Name:    "ocr2_bootstrap-" + uuid.NewString(),
		JobType: "bootstrap",
		OCR2OracleSpec: OracleSpec{
			ContractID: ocr2Addr,
			Relay:      "evm",
			RelayConfig: map[string]any{
				"chainID": bc.ChainID,
			},
			ContractConfigTrackerPollInterval: *NewInterval(5 * time.Second),
		},
	}
	_, err = bootstrapNode.MustCreateJob(bootstrapSpec)
	if err != nil {
		return fmt.Errorf("creating bootstrap job have failed: %w", err)
	}

	for _, chainlinkNode := range workerNodes {
		nodeTransmitterAddress, err := chainlinkNode.PrimaryEthAddress()
		if err != nil {
			return fmt.Errorf("getting primary ETH address from OCR node have failed: %w", err)
		}
		nodeOCRKeys, err := chainlinkNode.MustReadOCR2Keys()
		if err != nil {
			return fmt.Errorf("getting OCR keys from OCR node have failed: %w", err)
		}
		nodeOCRKeyID := nodeOCRKeys.Data[0].ID

		fakeServerURL := fake.Out.BaseURLDocker

		ea := &clclient.BridgeTypeAttributes{
			Name: "ea-" + uuid.NewString(),
			URL:  fmt.Sprintf("%s/%s", fakeServerURL, "ea"),
		}
		juelsBridge := &clclient.BridgeTypeAttributes{
			Name: "juels-" + uuid.NewString(),
			URL:  fmt.Sprintf("%s/%s", fakeServerURL, "juelsPerFeeCoinSource"),
		}
		err = chainlinkNode.MustCreateBridge(ea)
		if err != nil {
			return fmt.Errorf("creating bridge to %s on CL node failed: %w", ea.URL, err)
		}
		err = chainlinkNode.MustCreateBridge(juelsBridge)
		if err != nil {
			return fmt.Errorf("creating bridge to %s CL node failed: %w", juelsBridge.URL, err)
		}

		ocrSpec := &TaskJobSpec{
			Name:              "ocr2-" + uuid.NewString(),
			JobType:           "offchainreporting2",
			MaxTaskDuration:   (time.Duration(m.OCR2.Jobs.MaxTaskDurationSec) * time.Second).String(),
			ObservationSource: clclient.ObservationSourceSpecBridge(ea),
			ForwardingAllowed: false,
			OCR2OracleSpec: OracleSpec{
				PluginType: "median",
				Relay:      "evm",
				RelayConfig: map[string]any{
					"chainID": bc.ChainID,
				},
				PluginConfig: map[string]any{
					"juelsPerFeeCoinSource": fmt.Sprintf("\"\"\"%s\"\"\"", clclient.ObservationSourceSpecBridge(juelsBridge)),
				},
				ContractConfigTrackerPollInterval: *NewInterval(5 * time.Second),
				ContractID:                        ocr2Addr,                                // registryAddr
				OCRKeyBundleID:                    null.StringFrom(nodeOCRKeyID),           // get node ocr2config.ID
				TransmitterID:                     null.StringFrom(nodeTransmitterAddress), // node addr
				P2PV2Bootstrappers:                pq.StringArray{p2pV2Bootstrapper},       // bootstrap node key and address <p2p-key>@bootstrap:6690
			},
		}
		_, err = chainlinkNode.MustCreateJob(ocrSpec)
		if err != nil {
			return fmt.Errorf("creating OCR task job on OCR node have failed: %w", err)
		}
	}
	return nil
}
