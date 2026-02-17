package flux

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/flux_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	nodeset "github.com/smartcontractkit/chainlink-testing-framework/framework/components/simple_node_set"
	link_token "github.com/smartcontractkit/chainlink-testing-framework/seth/contracts/bind/link"

	"github.com/smartcontractkit/chainlink/devenv/products"
)

var L = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.DebugLevel).With().Fields(map[string]any{"component": "flux"}).Logger()

type Configurator struct {
	Config []*Flux `toml:"flux"`
}

type Flux struct {
	GasSettings       *products.GasSettings `toml:"gas_settings"`
	DeployedContracts *DeployedContracts    `toml:"deployed_contracts"`
}

type AggregatorOptions struct {
	PaymentAmount *big.Int       // The amount of LINK paid to each oracle per submission, in wei (units of 10⁻¹⁸ LINK)
	Timeout       uint32         // The number of seconds after the previous round that are allowed to lapse before allowing an oracle to skip an unfinished round
	Validator     common.Address // An optional contract address for validating external validation of answers
	MinSubValue   *big.Int       // An immutable check for a lower bound of what submission values are accepted from an oracle
	MaxSubValue   *big.Int       // An immutable check for an upper bound of what submission values are accepted from an oracle
	Decimals      uint8          // The number of decimals to offset the answer by
	Description   string         // A short description of what is being reported
}

type DeployedContracts struct {
	FluxAggregator string `toml:"flux_aggregator"`
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

	tx, err = lt.Mint(auth, common.HexToAddress(rootAddr), big.NewInt(1e18))
	if err != nil {
		return fmt.Errorf("could not transfer link token contract: %w", err)
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

	fluxOptions := &AggregatorOptions{
		PaymentAmount: big.NewInt(1),
		Timeout:       180,
		Validator:     common.Address{},
		MinSubValue:   big.NewInt(0),
		MaxSubValue:   big.NewInt(1000),
		Decimals:      18,
		Description:   "flux aggregator",
	}

	fluxAggAddr, fluxAggTx, fluxAggWrapper, err := flux_aggregator_wrapper.DeployFluxAggregator(auth, c, lt.Address(),
		fluxOptions.PaymentAmount,
		fluxOptions.Timeout,
		fluxOptions.Validator,
		fluxOptions.MinSubValue,
		fluxOptions.MaxSubValue,
		fluxOptions.Decimals,
		fluxOptions.Description,
	)
	if err != nil {
		return err
	}
	_, err = bind.WaitDeployed(ctx, c, fluxAggTx)
	if err != nil {
		return err
	}

	tx, err = lt.Mint(auth, fluxAggAddr, big.NewInt(1e18))
	if err != nil {
		return fmt.Errorf("could not transfer link token contract: %w", err)
	}
	_, err = products.WaitMinedFast(ctx, c, tx.Hash())
	if err != nil {
		return err
	}

	updateTx, err := fluxAggWrapper.UpdateAvailableFunds(auth)
	if err != nil {
		return fmt.Errorf("failed to update available funds on flux aggregator: %w", err)
	}

	if _, err = products.WaitMinedFast(ctx, c, updateTx.Hash()); err != nil {
		return err
	}

	changeTx, err := fluxAggWrapper.ChangeOracles(auth, []common.Address{}, transmitters, transmitters, 3, 3, 0)
	if err != nil {
		return fmt.Errorf("failed to change oracle addresses: %w", err)
	}

	if _, err = products.WaitMinedFast(ctx, c, changeTx.Hash()); err != nil {
		return err
	}

	bta := &clclient.BridgeTypeAttributes{
		Name: "variable-" + uuid.NewString()[0:5],
		URL:  fs.Out.BaseURLDocker + "/ea",
	}

	for _, n := range cls {
		err = n.MustCreateBridge(bta)
		if err != nil {
			return fmt.Errorf("failed to create new bridge: %w", err)
		}

		fluxSpec := &clclient.FluxMonitorJobSpec{
			Name:              "flux-" + uuid.NewString()[0:5],
			ContractAddress:   fluxAggAddr.String(),
			EVMChainID:        bc[0].ChainID,
			Threshold:         0,
			AbsoluteThreshold: 0,
			PollTimerPeriod:   1 * time.Minute,
			IdleTimerDisabled: true,
			ObservationSource: clclient.ObservationSourceSpecBridge(bta),
		}
		if _, err = n.MustCreateJob(fluxSpec); err != nil {
			return fmt.Errorf("failed to create flux aggregator job: %w", err)
		}
	}

	m.Config[0].DeployedContracts = &DeployedContracts{
		FluxAggregator: fluxAggAddr.String(),
	}
	return nil
}
