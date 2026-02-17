package directrequest

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/oracle_wrapper"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/test_api_consumer_wrapper"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/link_token"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	nodeset "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"

	"github.com/smartcontractkit/chainlink/devenv/products"
)

var L = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel).With().Fields(map[string]any{"component": "direct_request"}).Logger()

type Configurator struct {
	Config []*DirectRequest `toml:"direct_request"`
}

type DirectRequest struct {
	GasSettings *products.GasSettings `toml:"gas_settings"`
	Out         *Out                  `toml:"out"`
}

type Out struct {
	JobID    string `toml:"job_id"`
	Oracle   string `toml:"oracle"`
	Consumer string `toml:"consumer"`
}

func NewConfigurator() *Configurator {
	return &Configurator{}
}

func (m *Configurator) Load() error {
	cfg, err := products.Load[Configurator]()
	if err != nil {
		return fmt.Errorf("failed to load product config: %w", err)
	}
	m.Config = cfg.Config
	return nil
}

func (m *Configurator) Store(path string, instanceIdx int) error {
	if err := products.Store(".", m); err != nil {
		return fmt.Errorf("failed to store product config: %w", err)
	}
	return nil
}

func (m *Configurator) GenerateNodesConfig(
	ctx context.Context,
	fs *fake.Input,
	bc []*blockchain.Input,
	ns []*nodeset.Input,
) (string, error) {
	return products.DefaultLegacyCLNodeConfig(bc)
}

func (m *Configurator) GenerateNodesSecrets(
	_ context.Context,
	_ *fake.Input,
	_ []*blockchain.Input,
	_ []*nodeset.Input,
) (string, error) {
	return "", nil
}

func (m *Configurator) ConfigureJobsAndContracts(
	ctx context.Context,
	instanceIdx int,
	fs *fake.Input,
	bc []*blockchain.Input,
	ns []*nodeset.Input,
) error {
	L.Info().Msg("Connecting to CL nodes")
	cls, err := clclient.New(ns[0].Out.CLNodes)
	if err != nil {
		return fmt.Errorf("failed to connect to CL nodes: %w", err)
	}
	c, auth, rootAddr, err := products.ETHClient(ctx, bc[0].Out.Nodes[0].ExternalWSUrl, m.Config[0].GasSettings.FeeCapMultiplier, m.Config[0].GasSettings.TipCapMultiplier)
	if err != nil {
		return fmt.Errorf("failed to connect to blockchain :%w", err)
	}
	addr, tx, lt, err := link_token.DeployLinkToken(auth, c)
	if err != nil {
		return fmt.Errorf("could not create link token contract: %w", err)
	}
	_, err = bind.WaitDeployed(ctx, c, tx)
	if err != nil {
		return err
	}
	L.Info().Str("Address", addr.Hex()).Msg("Deployed link token contract")
	tx, err = lt.GrantMintRole(auth, common.HexToAddress(rootAddr))
	if err != nil {
		return fmt.Errorf("could not grant mint role: %w", err)
	}
	_, err = products.WaitMinedFast(ctx, c, tx.Hash())
	if err != nil {
		return err
	}

	transmitters := make([]common.Address, 0)
	for _, nc := range cls {
		addr, cErr := nc.ReadPrimaryETHKey(bc[0].Out.ChainID)
		if cErr != nil {
			return cErr
		}
		transmitters = append(transmitters, common.HexToAddress(addr.Attributes.Address))
	}
	pkey := products.NetworkPrivateKey()
	if pkey == "" {
		return errors.New("PRIVATE_KEY environment variable not set")
	}
	for _, addr := range transmitters {
		if cErr := products.FundAddressEIP1559(ctx, c, pkey, addr.String(), 10); cErr != nil {
			return cErr
		}
	}

	oracleAddr, oracleTx, _, err := oracle_wrapper.DeployOracle(auth, c, addr)
	if err != nil {
		return err
	}
	_, err = bind.WaitDeployed(ctx, c, oracleTx)
	if err != nil {
		return err
	}
	consumerAddr, consumerTx, _, err := test_api_consumer_wrapper.DeployTestAPIConsumer(auth, c, addr)
	if err != nil {
		return err
	}
	_, err = bind.WaitDeployed(ctx, c, consumerTx)
	if err != nil {
		return err
	}

	L.Info().Msgf("Minting LINK for consumer address: %s", consumerAddr)
	amount := new(big.Float).Mul(big.NewFloat(100), big.NewFloat(1e18))
	amountWei, _ := amount.Int(nil)
	tx, err = lt.Mint(auth, consumerAddr, amountWei)
	if err != nil {
		return fmt.Errorf("could not mint link token: %w", err)
	}
	_, err = products.WaitMinedFast(ctx, c, tx.Hash())
	if err != nil {
		return err
	}

	jobUUID := uuid.NewString()

	bta := clclient.BridgeTypeAttributes{
		Name: "direct_request_bridge-" + jobUUID,
		URL:  fs.Out.BaseURLDocker + "/direct_request_response",
	}
	if err = cls[0].MustCreateBridge(&bta); err != nil {
		return errors.New("failed to create bridge job: %w")
	}

	os := &clclient.DirectRequestTxPipelineSpec{
		BridgeTypeAttributes: bta,
		DataPath:             "data,result",
	}
	observationSource, err := os.String()
	if err != nil {
		return fmt.Errorf("failed to render Runlog job spec: %w", err)
	}

	_, err = cls[0].MustCreateJob(&clclient.DirectRequestJobSpec{
		Name:                     "direct_request-" + uuid.NewString(),
		MinIncomingConfirmations: "0",
		ContractAddress:          oracleAddr.String(),
		EVMChainID:               bc[0].ChainID,
		ExternalJobID:            jobUUID,
		ObservationSource:        observationSource,
	})
	if err != nil {
		return fmt.Errorf("failed to create runlog job: %w", err)
	}
	m.Config[0].Out = &Out{
		JobID:    strings.Replace(jobUUID, "-", "", 4),
		Oracle:   oracleAddr.String(),
		Consumer: consumerAddr.String(),
	}
	return nil
}
