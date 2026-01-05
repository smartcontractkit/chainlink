package devenv

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
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
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/link_token"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	"github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"golang.org/x/sync/errgroup"
	"gopkg.in/guregu/null.v4"
)

const (
	ConfigureNodesNetwork ConfigPhase = iota
	ConfigureProductContractsJobs
)

var Plog = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel).With().Fields(map[string]any{"component": "ocr2"}).Logger()

type OCR2 struct {
	OCR2                             *OCRv2OffChainOptions  `toml:"ocr2"`
	OCR2SetConfig                    *OCRv2SetConfigOptions `toml:"ocr2_set_config"`
	OCR2SetConfigOut                 *OCRv2Config           `toml:"ocr2_set_config_out"`
	OCR2MedianOffchainConfig         *MedianOffchainConfig  `toml:"ocr2_median_offchain_config"`
	EAFake                           *EAFake                `toml:"ea_fake"`
	Jobs                             *Jobs                  `toml:"jobs"`
	LinkContractAddress              string                 `toml:"link_contract_address"`
	CLNodesFundingETH                float64                `toml:"cl_nodes_funding_eth"`
	CLNodesFundingLink               float64                `toml:"cl_nodes_funding_link"`
	ChainFinalityDepth               int64                  `toml:"chain_finality_depth"`
	VerificationTimeoutSec           time.Duration          `toml:"verification_timeout_sec"`            //nolint:staticcheck
	ContractsConfigurationTimeoutSec time.Duration          `toml:"contracts_configuration_timeout_sec"` //nolint:staticcheck
	GasSettings                      *GasSettings           `toml:"gas_settings"`
	Verify                           bool                   `toml:"verify"`
	DeployedContracts                *DeployedContracts     `toml:"deployed_contracts"`
}

type DeployedContracts struct {
	OCRv2AggregatorAddr string `toml:"ocr2_aggregator_address"`
}

type GasSettings struct {
	FeeCapMultiplier int64 `toml:"fee_cap_multiplier"`
	TipCapMultiplier int64 `toml:"tip_cap_multiplier"`
}

type MedianOffchainConfig struct {
	AlphaReportPPB      uint64        `toml:"alpha_report_ppb"`
	AlphaAcceptPPB      uint64        `toml:"alpha_accept_ppb"`
	DeltaCSec           time.Duration `toml:"delta_sec"`
	AlphaReportInfinite bool          `toml:"alpha_report_infinite"`
	AlphaAcceptInfinite bool          `toml:"alpha_accept_infinite"`
}

type Jobs struct {
	ConfigPollIntervalSeconds time.Duration `toml:"config_poll_interval_sec"`
	MaxTaskDurationSec        time.Duration `toml:"max_task_duration_sec"`
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

// deployLinkAndMint is a universal action that deploys link token and mints required amount of LINK token for all the nodes.
func deployLinkAndMint(ctx context.Context, in *Cfg, c *ethclient.Client, auth *bind.TransactOpts, rootAddr string, transmitters []common.Address) (*link_token.LinkToken, error) {
	addr, tx, lt, err := link_token.DeployLinkToken(auth, c)
	if err != nil {
		return nil, fmt.Errorf("could not create link token contract: %w", err)
	}
	_, err = bind.WaitDeployed(ctx, c, tx)
	if err != nil {
		return nil, err
	}
	Plog.Info().Str("Address", addr.Hex()).Msg("Deployed link token contract")
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
		amount := new(big.Float).Mul(big.NewFloat(in.OCR2.CLNodesFundingLink), big.NewFloat(1e18))
		amountWei, _ := amount.Int(nil)
		Plog.Info().Msgf("Minting LINK for transmitter address: %s", transmitter.Hex())
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

func UpdateOCR2ConfigOffChainValues(in *Cfg, ocr2i *ocr2aggregator.OCR2Aggregator, cl []*clclient.ChainlinkClient, o2 *OCRv2SetConfigOptions) error {
	if o2 == nil {
		return nil
	}
	ctx := context.Background()
	c, auth, _, err := ETHClient(in)
	if err != nil {
		return fmt.Errorf("could not create basic eth client: %w", err)
	}
	// generating oracle identities and setting up OCRv2
	s, ids, err := getOracleIdentities(in, cl)
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
			AlphaAcceptInfinite: in.OCR2.OCR2MedianOffchainConfig.AlphaAcceptInfinite,
			AlphaReportInfinite: in.OCR2.OCR2MedianOffchainConfig.AlphaReportInfinite,
			AlphaReportPPB:      in.OCR2.OCR2MedianOffchainConfig.AlphaReportPPB,
			AlphaAcceptPPB:      in.OCR2.OCR2MedianOffchainConfig.AlphaAcceptPPB,
			DeltaC:              in.OCR2.OCR2MedianOffchainConfig.DeltaCSec * time.Second,
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
	onChainConfig, err := median.StandardOnchainConfigCodec{}.Encode(context.Background(), median.OnchainConfig{Min: in.OCR2.OCR2.MinimumAnswer, Max: in.OCR2.OCR2.MaximumAnswer})
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

func configureContracts(in *Cfg, c *ethclient.Client, auth *bind.TransactOpts, cl []*clclient.ChainlinkClient, rootAddr string, transmitters []common.Address) (*OCRv2Config, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), in.OCR2.ContractsConfigurationTimeoutSec*time.Second)
	defer cancel()
	Plog.Info().Msg("Deploying LINK token contract")
	lt, err := deployLinkAndMint(ctx, in, c, auth, rootAddr, transmitters)
	if err != nil {
		return nil, "", fmt.Errorf("could not create link token contract and mint: %w", err)
	}
	// OCRv2 Aggregator
	Plog.Info().Msg("Deploying OCRv2 aggregator contract")
	opts := in.OCR2.OCR2
	ocr2addr, tx, ocr2i, err := ocr2aggregator.DeployOCR2Aggregator(auth, c, lt.Address(), opts.MinimumAnswer, opts.MaximumAnswer, common.HexToAddress(""), common.HexToAddress(""), 18, "")
	if err != nil {
		return nil, "", fmt.Errorf("could not create ocr2 aggregator contract: %w", err)
	}
	_, err = bind.WaitDeployed(ctx, c, tx)
	if err != nil {
		return nil, "", err
	}
	Plog.Info().Str("Address", ocr2addr.String()).Msg("Deployed OCRv2 Aggregator contract")
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
	s, ids, err := getOracleIdentities(in, cl)
	if err != nil {
		return nil, "", fmt.Errorf("could not get oracle identities: %w", err)
	}
	o2 := in.OCR2.OCR2SetConfig
	signerKeys, transmitterAccounts, f, _, offchainConfigVersion, offchainConfig, err := confighelper.ContractSetConfigArgsForTests(
		o2.DeltaProgress*time.Second,
		o2.DeltaResend*time.Second,
		o2.DeltaRound*time.Second,
		o2.DeltaGrace*time.Second,
		o2.DeltaStage*time.Second,
		o2.RMax,
		s,
		ids,
		median.OffchainConfig{
			AlphaAcceptInfinite: in.OCR2.OCR2MedianOffchainConfig.AlphaAcceptInfinite,
			AlphaReportInfinite: in.OCR2.OCR2MedianOffchainConfig.AlphaReportInfinite,
			AlphaReportPPB:      in.OCR2.OCR2MedianOffchainConfig.AlphaReportPPB,
			AlphaAcceptPPB:      in.OCR2.OCR2MedianOffchainConfig.AlphaAcceptPPB,
			DeltaC:              in.OCR2.OCR2MedianOffchainConfig.DeltaCSec * time.Second,
		}.Encode(),
		nil,
		o2.MaxDurationQuery*time.Second,
		o2.MaxDurationObservation*time.Second,
		o2.MaxDurationReport*time.Second,
		o2.MaxDurationShouldAcceptFinalizedReport*time.Second,
		o2.MaxDurationShouldTransmitAcceptedReport*time.Second,
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
	onChainConfig, err := median.StandardOnchainConfigCodec{}.Encode(context.Background(), median.OnchainConfig{Min: in.OCR2.OCR2.MinimumAnswer, Max: in.OCR2.OCR2.MaximumAnswer})
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

func getOracleIdentities(_ *Cfg, clClients []*clclient.ChainlinkClient) ([]int, []confighelper.OracleIdentityExtra, error) { //nolint:gocritic
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
				return fmt.Errorf("wrong number of elements copied")
			}
			configPkBytes, err := hex.DecodeString(strings.TrimPrefix(ocr2Config.ConfigPublicKey, "ocr2cfg_evm_"))
			if err != nil {
				return err
			}
			configPkBytesFixed := [ed25519.PublicKeySize]byte{}
			n = copy(configPkBytesFixed[:], configPkBytes)
			if n != ed25519.PublicKeySize {
				return fmt.Errorf("wrong number of elements copied")
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
			Plog.Trace().
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

func configureJobs(in *Cfg, clNodes []*clclient.ChainlinkClient, ocr2Addr string) error {
	bootstrapNode := clNodes[0]
	workerNodes := clNodes[1:]
	bootstrapP2PIds, err := bootstrapNode.MustReadP2PKeys()
	if err != nil {
		return err
	}
	p2pV2Bootstrapper := fmt.Sprintf("%s@%s:%d", bootstrapP2PIds.Data[0].Attributes.PeerID, in.NodeSets[0].Out.CLNodes[0].Node.ContainerName, 6690)
	// Set the value for the jobs to report on
	bootstrapSpec := &OCR2TaskJobSpec{
		Name:    fmt.Sprintf("ocr2_bootstrap-%s", uuid.NewString()),
		JobType: "bootstrap",
		OCR2OracleSpec: OCR2OracleSpec{
			ContractID: ocr2Addr,
			Relay:      "evm",
			RelayConfig: map[string]any{
				"chainID": in.Blockchains[0].ChainID,
			},
			ContractConfigTrackerPollInterval: *NewInterval(in.OCR2.Jobs.ConfigPollIntervalSeconds * time.Second),
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

		fakeServerURL := in.FakeServer.Out.BaseURLDocker

		ea := &clclient.BridgeTypeAttributes{
			Name: fmt.Sprintf("ea-%s", uuid.NewString()),
			URL:  fmt.Sprintf("%s/%s", fakeServerURL, "ea"),
		}
		juelsBridge := &clclient.BridgeTypeAttributes{
			Name: fmt.Sprintf("juels-%s", uuid.NewString()),
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

		ocrSpec := &OCR2TaskJobSpec{
			Name:              fmt.Sprintf("ocr2-%s", uuid.NewString()),
			JobType:           "offchainreporting2",
			MaxTaskDuration:   (in.OCR2.Jobs.MaxTaskDurationSec * time.Second).String(),
			ObservationSource: clclient.ObservationSourceSpecBridge(ea),
			ForwardingAllowed: false,
			OCR2OracleSpec: OCR2OracleSpec{
				PluginType: "median",
				Relay:      "evm",
				RelayConfig: map[string]any{
					"chainID": in.Blockchains[0].ChainID,
				},
				PluginConfig: map[string]any{
					"juelsPerFeeCoinSource": fmt.Sprintf("\"\"\"%s\"\"\"", clclient.ObservationSourceSpecBridge(juelsBridge)), //nolint:gocritic
				},
				ContractConfigTrackerPollInterval: *NewInterval(in.OCR2.Jobs.ConfigPollIntervalSeconds * time.Second),
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

// OCR2ProductConfiguration is default product configuration that includes:
// - Deploying required prerequisites (LINK token, shared contracts)
// - Applying product-specific changesets
// - Creating cldf.Environment, connecting to components, see *Cfg fields
// - Generating CL nodes configs
// All the data can be added *Cfg struct like and is synced between local machine and remote environment
// so later both local and remote tests can use it.
func OCR2ProductConfiguration(in *Cfg, phase ConfigPhase) error {
	pkey := getNetworkPrivateKey()
	if pkey == "" {
		return fmt.Errorf("PRIVATE_KEY environment variable not set")
	}
	switch phase {
	case ConfigureNodesNetwork:
		Plog.Info().Msg("Applying default CL nodes configuration")
		node := in.Blockchains[0].Out.Nodes[0]
		chainID := in.Blockchains[0].ChainID
		// configure node set and generate CL nodes configs
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
`, in.OCR2.LinkContractAddress, chainID, in.OCR2.ChainFinalityDepth, node.InternalWSUrl, node.InternalHTTPUrl)
		for _, nodeSpec := range in.NodeSets[0].NodeSpecs {
			nodeSpec.Node.TestConfigOverrides = netConfig
		}
		Plog.Info().Msg("Nodes network configuration is finished")
	case ConfigureProductContractsJobs:
		Plog.Info().Msg("Connecting to CL nodes")
		nodeClients, err := clclient.New(in.NodeSets[0].Out.CLNodes)
		if err != nil {
			return err
		}
		transmitters := make([]common.Address, 0)
		ethKeyAddresses := make([]string, 0)
		for i, nc := range nodeClients {
			addr, err := nc.ReadPrimaryETHKey(in.Blockchains[0].ChainID)
			if err != nil {
				return err
			}
			ethKeyAddresses = append(ethKeyAddresses, addr.Attributes.Address)
			transmitters = append(transmitters, common.HexToAddress(addr.Attributes.Address))
			Plog.Info().
				Int("Idx", i).
				Str("ETH", addr.Attributes.Address).
				Msg("Node info")
		}
		c, auth, rootAddr, err := ETHClient(in)
		if err != nil {
			return fmt.Errorf("could not create basic eth client: %w", err)
		}
		for _, addr := range ethKeyAddresses {
			if err := FundNodeEIP1559(c, pkey, addr, in.OCR2.CLNodesFundingETH); err != nil {
				return err
			}
		}
		ocrv2Config, ocr2Addr, err := configureContracts(in, c, auth, nodeClients, rootAddr, transmitters)
		if err != nil {
			return err
		}
		in.OCR2.OCR2SetConfigOut = ocrv2Config
		if err := configureJobs(in, nodeClients, ocr2Addr); err != nil {
			return err
		}
		r := resty.New().SetBaseURL(in.FakeServer.Out.BaseURLHost)

		_, err = r.R().Post(`/trigger_deviation?result=200`)
		if err != nil {
			return fmt.Errorf("could not set ea fake values: %w", err)
		}
		Plog.Info().
			Msg("Setting fake external adapter (data feed) values")
		Plog.Info().Str("BootstrapNode", in.NodeSets[0].Out.CLNodes[0].Node.ExternalURL).Send()
		for _, n := range in.NodeSets[0].Out.CLNodes[1:] {
			Plog.Info().Str("Node", n.Node.ExternalURL).Send()
		}
		if err := verifyEnvironment(in, c, ocr2Addr); err != nil {
			return err
		}
		in.OCR2.DeployedContracts = &DeployedContracts{OCRv2AggregatorAddr: ocr2Addr}
	}
	return nil
}
