package contracts

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/counter"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/i_automation_registry_master_wrapper_2_2"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/i_automation_registry_master_wrapper_2_3"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/keeper_registry_wrapper1_1"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/keeper_registry_wrapper1_2"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/keeper_registry_wrapper1_3"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/keeper_registry_wrapper2_0"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/mock_ethlink_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/mock_ethusd_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/mock_gas_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/simple_log_upkeep_counter_wrapper"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/weth9"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"
)

// OCRv2Config represents the config for the OCRv2 contract
type OCRv2Config struct {
	Signers               []common.Address
	Transmitters          []common.Address
	F                     uint8
	OnchainConfig         []byte
	TypedOnchainConfig21  i_keeper_registry_master_wrapper_2_1.IAutomationV21PlusCommonOnchainConfigLegacy
	TypedOnchainConfig22  i_automation_registry_master_wrapper_2_2.AutomationRegistryBase22OnchainConfig
	TypedOnchainConfig23  i_automation_registry_master_wrapper_2_3.AutomationRegistryBase23OnchainConfig
	OffchainConfigVersion uint64
	OffchainConfig        []byte
	BillingTokens         []common.Address
	BillingConfigs        []i_automation_registry_master_wrapper_2_3.AutomationRegistryBase23BillingConfig
}

func Bytes32ToSlice(a [32]byte) (r []byte) {
	r = append(r, a[:]...)
	return
}

func GetRegistryContractABI(version KeeperRegistryVersion) (*abi.ABI, error) {
	var (
		contractABI *abi.ABI
		err         error
	)
	switch version {
	case RegistryVersion_1_0, RegistryVersion_1_1:
		contractABI, err = keeper_registry_wrapper1_1.KeeperRegistryMetaData.GetAbi()
	case RegistryVersion_1_2:
		contractABI, err = keeper_registry_wrapper1_2.KeeperRegistryMetaData.GetAbi()
	case RegistryVersion_1_3:
		contractABI, err = keeper_registry_wrapper1_3.KeeperRegistryMetaData.GetAbi()
	case RegistryVersion_2_0:
		contractABI, err = keeper_registry_wrapper2_0.KeeperRegistryMetaData.GetAbi()
	case RegistryVersion_2_1:
		contractABI, err = i_keeper_registry_master_wrapper_2_1.IKeeperRegistryMasterMetaData.GetAbi()
	case RegistryVersion_2_2:
		contractABI, err = i_automation_registry_master_wrapper_2_2.IAutomationRegistryMasterMetaData.GetAbi()
	default:
		contractABI, err = keeper_registry_wrapper2_0.KeeperRegistryMetaData.GetAbi()
	}

	return contractABI, err
}

// EthereumLinkToken represents a LinkToken address
type EthereumLinkToken struct {
	client   *seth.Client
	instance *link_token_interface.LinkToken
	address  common.Address
	l        zerolog.Logger
}

func (l *EthereumLinkToken) Decimals() uint {
	return 18
}

func DeployLinkTokenContract(l zerolog.Logger, client *seth.Client) (*EthereumLinkToken, error) {
	linkTokenAbi, err := link_token_interface.LinkTokenMetaData.GetAbi()
	if err != nil {
		return &EthereumLinkToken{}, fmt.Errorf("failed to get LinkToken ABI: %w", err)
	}
	linkDeploymentData, err := client.DeployContract(client.NewTXOpts(), "LinkToken", *linkTokenAbi, common.FromHex(link_token_interface.LinkTokenMetaData.Bin))
	if err != nil {
		return &EthereumLinkToken{}, fmt.Errorf("LinkToken instance deployment have failed: %w", err)
	}

	linkToken, err := link_token_interface.NewLinkToken(linkDeploymentData.Address, MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumLinkToken{}, fmt.Errorf("failed to instantiate LinkToken instance: %w", err)
	}

	return &EthereumLinkToken{
		client:   client,
		instance: linkToken,
		address:  linkDeploymentData.Address,
		l:        l,
	}, nil
}

func LoadLinkTokenContract(l zerolog.Logger, client *seth.Client, address common.Address) (*EthereumLinkToken, error) {
	loader := seth.NewContractLoader[link_token_interface.LinkToken](client)
	instance, err := loader.LoadContract("LinkToken", address, link_token_interface.LinkTokenMetaData.GetAbi, link_token_interface.NewLinkToken)

	if err != nil {
		return &EthereumLinkToken{}, fmt.Errorf("failed to instantiate LinkToken instance: %w", err)
	}

	return &EthereumLinkToken{
		client:   client,
		instance: instance,
		address:  address,
		l:        l,
	}, nil
}

// Fund the LINK Token contract with ETH to distribute the token
func (l *EthereumLinkToken) Fund(_ *big.Float) error {
	panic("do not use this function, use actions.SendFunds instead")
}

func (l *EthereumLinkToken) BalanceOf(ctx context.Context, addr string) (*big.Int, error) {
	return l.instance.BalanceOf(&bind.CallOpts{
		From:    l.client.Addresses[0],
		Context: ctx,
	}, common.HexToAddress(addr))
}

// Name returns the name of the link token
func (l *EthereumLinkToken) Name(ctx context.Context) (string, error) {
	return l.instance.Name(&bind.CallOpts{
		From:    l.client.Addresses[0],
		Context: ctx,
	})
}

func (l *EthereumLinkToken) Address() string {
	return l.address.Hex()
}

func (l *EthereumLinkToken) Approve(to string, amount *big.Int) error {
	l.l.Info().
		Str("From", l.client.Addresses[0].Hex()).
		Str("To", to).
		Str("Amount", amount.String()).
		Msg("Approving LINK Transfer")
	_, err := l.client.Decode(l.instance.Approve(l.client.NewTXOpts(), common.HexToAddress(to), amount))
	return err
}

func (l *EthereumLinkToken) Transfer(to string, amount *big.Int) error {
	l.l.Info().
		Str("From", l.client.Addresses[0].Hex()).
		Str("To", to).
		Str("Amount", amount.String()).
		Msg("Transferring LINK")
	_, err := l.client.Decode(l.instance.Transfer(l.client.NewTXOpts(), common.HexToAddress(to), amount))
	return err
}

func (l *EthereumLinkToken) TransferAndCall(to string, amount *big.Int, data []byte) (*types.Transaction, error) {
	l.l.Info().
		Str("From", l.client.Addresses[0].Hex()).
		Str("To", to).
		Str("Amount", amount.String()).
		Msg("Transferring and Calling LINK")
	decodedTx, err := l.client.Decode(l.instance.TransferAndCall(l.client.NewTXOpts(), common.HexToAddress(to), amount, data))
	if err != nil {
		return nil, err
	}
	return decodedTx.Transaction, nil
}

func (l *EthereumLinkToken) TransferAndCallFromKey(to string, amount *big.Int, data []byte, keyNum int) (*types.Transaction, error) {
	l.l.Info().
		Str("From", l.client.Addresses[keyNum].Hex()).
		Str("To", to).
		Str("Amount", amount.String()).
		Msg("Transferring and Calling LINK")
	decodedTx, err := l.client.Decode(l.instance.TransferAndCall(l.client.NewTXKeyOpts(keyNum), common.HexToAddress(to), amount, data))
	if err != nil {
		return nil, err
	}
	return decodedTx.Transaction, nil
}

// EthereumMockETHLINKFeed represents mocked ETH/LINK feed contract
type EthereumMockETHLINKFeed struct {
	client  *seth.Client
	feed    *mock_ethlink_aggregator_wrapper.MockETHLINKAggregator
	address *common.Address
}

func (v *EthereumMockETHLINKFeed) Address() string {
	return v.address.Hex()
}

func (v *EthereumMockETHLINKFeed) LatestRoundData() (*big.Int, error) {
	data, err := v.feed.LatestRoundData(&bind.CallOpts{
		From:    v.client.Addresses[0],
		Context: context.Background(),
	})
	if err != nil {
		return nil, err
	}
	return data.Ans, nil
}

func (v *EthereumMockETHLINKFeed) LatestRoundDataUpdatedAt() (*big.Int, error) {
	data, err := v.feed.LatestRoundData(&bind.CallOpts{
		From:    v.client.Addresses[0],
		Context: context.Background(),
	})
	if err != nil {
		return nil, err
	}
	return data.UpdatedAt, nil
}

func DeployMockLINKETHFeed(client *seth.Client, answer *big.Int) (MockLINKETHFeed, error) {
	abi, err := mock_ethlink_aggregator_wrapper.MockETHLINKAggregatorMetaData.GetAbi()
	if err != nil {
		return &EthereumMockETHLINKFeed{}, fmt.Errorf("failed to get MockLINKETHFeed ABI: %w", err)
	}
	data, err := client.DeployContract(client.NewTXOpts(), "MockLINKETHFeed", *abi, common.FromHex(mock_ethlink_aggregator_wrapper.MockETHLINKAggregatorMetaData.Bin), answer)
	if err != nil {
		return &EthereumMockETHLINKFeed{}, fmt.Errorf("MockLINKETHFeed instance deployment have failed: %w", err)
	}

	instance, err := mock_ethlink_aggregator_wrapper.NewMockETHLINKAggregator(data.Address, MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumMockETHLINKFeed{}, fmt.Errorf("failed to instantiate MockLINKETHFeed instance: %w", err)
	}

	return &EthereumMockETHLINKFeed{
		address: &data.Address,
		client:  client,
		feed:    instance,
	}, nil
}

func LoadMockLINKETHFeed(client *seth.Client, address common.Address) (MockLINKETHFeed, error) {
	abi, err := mock_ethlink_aggregator_wrapper.MockETHLINKAggregatorMetaData.GetAbi()
	if err != nil {
		return &EthereumMockETHLINKFeed{}, fmt.Errorf("failed to get MockLINKETHFeed ABI: %w", err)
	}
	client.ContractStore.AddABI("MockLINKETHFeed", *abi)
	client.ContractStore.AddBIN("MockLINKETHFeed", common.FromHex(mock_ethlink_aggregator_wrapper.MockETHLINKAggregatorMetaData.Bin))

	instance, err := mock_ethlink_aggregator_wrapper.NewMockETHLINKAggregator(address, MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumMockETHLINKFeed{}, fmt.Errorf("failed to instantiate MockLINKETHFeed instance: %w", err)
	}

	return &EthereumMockETHLINKFeed{
		address: &address,
		client:  client,
		feed:    instance,
	}, nil
}

// EthereumMockGASFeed represents mocked Gas feed contract
type EthereumMockGASFeed struct {
	client  *seth.Client
	feed    *mock_gas_aggregator_wrapper.MockGASAggregator
	address *common.Address
}

func (v *EthereumMockGASFeed) Address() string {
	return v.address.Hex()
}

func DeployMockGASFeed(client *seth.Client, answer *big.Int) (MockGasFeed, error) {
	abi, err := mock_gas_aggregator_wrapper.MockGASAggregatorMetaData.GetAbi()
	if err != nil {
		return &EthereumMockGASFeed{}, fmt.Errorf("failed to get MockGasFeed ABI: %w", err)
	}
	data, err := client.DeployContract(client.NewTXOpts(), "MockGasFeed", *abi, common.FromHex(mock_gas_aggregator_wrapper.MockGASAggregatorMetaData.Bin), answer)
	if err != nil {
		return &EthereumMockGASFeed{}, fmt.Errorf("MockGasFeed instance deployment have failed: %w", err)
	}

	instance, err := mock_gas_aggregator_wrapper.NewMockGASAggregator(data.Address, MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumMockGASFeed{}, fmt.Errorf("failed to instantiate MockGasFeed instance: %w", err)
	}

	return &EthereumMockGASFeed{
		address: &data.Address,
		client:  client,
		feed:    instance,
	}, nil
}

func LoadMockGASFeed(client *seth.Client, address common.Address) (MockGasFeed, error) {
	abi, err := mock_gas_aggregator_wrapper.MockGASAggregatorMetaData.GetAbi()
	if err != nil {
		return &EthereumMockGASFeed{}, fmt.Errorf("failed to get MockGasFeed ABI: %w", err)
	}
	client.ContractStore.AddABI("MockGasFeed", *abi)
	client.ContractStore.AddBIN("MockGasFeed", common.FromHex(mock_gas_aggregator_wrapper.MockGASAggregatorMetaData.Bin))

	instance, err := mock_gas_aggregator_wrapper.NewMockGASAggregator(address, MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumMockGASFeed{}, fmt.Errorf("failed to instantiate MockGasFeed instance: %w", err)
	}

	return &EthereumMockGASFeed{
		address: &address,
		client:  client,
		feed:    instance,
	}, nil
}

func DeployMultiCallContract(client *seth.Client) (common.Address, error) {
	abi, err := abi.JSON(strings.NewReader(MultiCallABI))
	if err != nil {
		return common.Address{}, err
	}

	data, err := client.DeployContract(client.NewTXOpts(), "MultiCall", abi, common.FromHex(MultiCallBIN))
	if err != nil {
		return common.Address{}, fmt.Errorf("MultiCall instance deployment have failed: %w", err)
	}

	return data.Address, nil
}

// EthereumWETHToken represents a WETH address
type EthereumWETHToken struct {
	client   *seth.Client
	instance *weth9.WETH9
	address  common.Address
	l        zerolog.Logger
}

func DeployWETHTokenContract(l zerolog.Logger, client *seth.Client) (*EthereumWETHToken, error) {
	wethTokenAbi, err := weth9.WETH9MetaData.GetAbi()
	if err != nil {
		return &EthereumWETHToken{}, fmt.Errorf("failed to get WETH token ABI: %w", err)
	}
	wethDeploymentData, err := client.DeployContract(client.NewTXOpts(), "WETHToken", *wethTokenAbi, common.FromHex(weth9.WETH9MetaData.Bin))
	if err != nil {
		return &EthereumWETHToken{}, fmt.Errorf("WETH token instance deployment failed: %w", err)
	}

	wethToken, err := weth9.NewWETH9(wethDeploymentData.Address, MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumWETHToken{}, fmt.Errorf("failed to instantiate WETHToken instance: %w", err)
	}

	return &EthereumWETHToken{
		client:   client,
		instance: wethToken,
		address:  wethDeploymentData.Address,
		l:        l,
	}, nil
}

func LoadWETHTokenContract(l zerolog.Logger, client *seth.Client, address common.Address) (*EthereumWETHToken, error) {
	abi, err := weth9.WETH9MetaData.GetAbi()
	if err != nil {
		return &EthereumWETHToken{}, fmt.Errorf("failed to get WETH token ABI: %w", err)
	}

	client.ContractStore.AddABI("WETHToken", *abi)
	client.ContractStore.AddBIN("WETHToken", common.FromHex(weth9.WETH9MetaData.Bin))

	wethToken, err := weth9.NewWETH9(address, MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumWETHToken{}, fmt.Errorf("failed to instantiate WETHToken instance: %w", err)
	}

	return &EthereumWETHToken{
		client:   client,
		instance: wethToken,
		address:  address,
		l:        l,
	}, nil
}

// Fund the WETH Token contract with ETH to distribute the token
func (l *EthereumWETHToken) Fund(_ *big.Float) error {
	panic("do not use this function, use actions_seth.SendFunds instead")
}

func (l *EthereumWETHToken) Decimals() uint {
	return 18
}

func (l *EthereumWETHToken) BalanceOf(ctx context.Context, addr string) (*big.Int, error) {
	return l.instance.BalanceOf(&bind.CallOpts{
		From:    l.client.Addresses[0],
		Context: ctx,
	}, common.HexToAddress(addr))
}

// Name returns the name of the weth token
func (l *EthereumWETHToken) Name(ctx context.Context) (string, error) {
	return l.instance.Name(&bind.CallOpts{
		From:    l.client.Addresses[0],
		Context: ctx,
	})
}

func (l *EthereumWETHToken) Address() string {
	return l.address.Hex()
}

func (l *EthereumWETHToken) Approve(to string, amount *big.Int) error {
	l.l.Info().
		Str("From", l.client.Addresses[0].Hex()).
		Str("To", to).
		Str("Amount", amount.String()).
		Msg("Approving WETH Transfer")
	_, err := l.client.Decode(l.instance.Approve(l.client.NewTXOpts(), common.HexToAddress(to), amount))
	return err
}

func (l *EthereumWETHToken) Transfer(to string, amount *big.Int) error {
	l.l.Info().
		Str("From", l.client.Addresses[0].Hex()).
		Str("To", to).
		Str("Amount", amount.String()).
		Msg("Transferring WETH")
	_, err := l.client.Decode(l.instance.Transfer(l.client.NewTXOpts(), common.HexToAddress(to), amount))
	return err
}

// EthereumMockETHUSDFeed represents mocked ETH/USD feed contract
// For the integration tests, we also use this ETH/USD feed for LINK/USD feed since they have the same structure
type EthereumMockETHUSDFeed struct {
	client  *seth.Client
	feed    *mock_ethusd_aggregator_wrapper.MockETHUSDAggregator
	address *common.Address
}

func (l *EthereumMockETHUSDFeed) Decimals() uint {
	return 8
}

func (l *EthereumMockETHUSDFeed) Address() string {
	return l.address.Hex()
}

func (l *EthereumMockETHUSDFeed) LatestRoundData() (*big.Int, error) {
	data, err := l.feed.LatestRoundData(&bind.CallOpts{
		From:    l.client.Addresses[0],
		Context: context.Background(),
	})
	if err != nil {
		return nil, err
	}
	return data.Ans, nil
}

func (l *EthereumMockETHUSDFeed) LatestRoundDataUpdatedAt() (*big.Int, error) {
	data, err := l.feed.LatestRoundData(&bind.CallOpts{
		From:    l.client.Addresses[0],
		Context: context.Background(),
	})
	if err != nil {
		return nil, err
	}
	return data.UpdatedAt, nil
}

func DeployMockETHUSDFeed(client *seth.Client, answer *big.Int) (MockETHUSDFeed, error) {
	abi, err := mock_ethusd_aggregator_wrapper.MockETHUSDAggregatorMetaData.GetAbi()
	if err != nil {
		return &EthereumMockETHUSDFeed{}, fmt.Errorf("failed to get MockETHUSDFeed ABI: %w", err)
	}
	data, err := client.DeployContract(client.NewTXOpts(), "MockETHUSDFeed", *abi, common.FromHex(mock_ethusd_aggregator_wrapper.MockETHUSDAggregatorMetaData.Bin), answer)
	if err != nil {
		return &EthereumMockETHUSDFeed{}, fmt.Errorf("MockETHUSDFeed instance deployment have failed: %w", err)
	}

	instance, err := mock_ethusd_aggregator_wrapper.NewMockETHUSDAggregator(data.Address, MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumMockETHUSDFeed{}, fmt.Errorf("failed to instantiate MockETHUSDFeed instance: %w", err)
	}

	return &EthereumMockETHUSDFeed{
		address: &data.Address,
		client:  client,
		feed:    instance,
	}, nil
}

func LoadMockETHUSDFeed(client *seth.Client, address common.Address) (MockETHUSDFeed, error) {
	abi, err := mock_ethusd_aggregator_wrapper.MockETHUSDAggregatorMetaData.GetAbi()
	if err != nil {
		return &EthereumMockETHUSDFeed{}, fmt.Errorf("failed to get MockETHUSDFeed ABI: %w", err)
	}
	client.ContractStore.AddABI("MockETHUSDFeed", *abi)
	client.ContractStore.AddBIN("MockETHUSDFeed", common.FromHex(mock_ethusd_aggregator_wrapper.MockETHUSDAggregatorMetaData.Bin))

	instance, err := mock_ethusd_aggregator_wrapper.NewMockETHUSDAggregator(address, MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumMockETHUSDFeed{}, fmt.Errorf("failed to instantiate MockETHUSDFeed instance: %w", err)
	}

	return &EthereumMockETHUSDFeed{
		address: &address,
		client:  client,
		feed:    instance,
	}, nil
}

type Counter struct {
	client   *seth.Client
	instance *counter.Counter
	address  common.Address
}

func DeployCounterContract(client *seth.Client) (*Counter, error) {
	abi, err := counter.CounterMetaData.GetAbi()
	if err != nil {
		return &Counter{}, fmt.Errorf("failed to get Counter ABI: %w", err)
	}
	linkDeploymentData, err := client.DeployContract(client.NewTXOpts(), "Counter", *abi, common.FromHex(counter.CounterMetaData.Bin))
	if err != nil {
		return &Counter{}, fmt.Errorf("Counter instance deployment have failed: %w", err)
	}

	instance, err := counter.NewCounter(linkDeploymentData.Address, MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &Counter{}, fmt.Errorf("failed to instantiate Counter instance: %w", err)
	}

	return &Counter{
		client:   client,
		instance: instance,
		address:  linkDeploymentData.Address,
	}, nil
}

func (c *Counter) Address() string {
	return c.address.Hex()
}

func (c *Counter) Increment() error {
	_, err := c.client.Decode(c.instance.Increment(
		c.client.NewTXOpts(),
	))
	return err
}

func (c *Counter) Reset() error {
	_, err := c.client.Decode(c.instance.Reset(
		c.client.NewTXOpts(),
	))
	return err
}

func (c *Counter) Count() (*big.Int, error) {
	data, err := c.instance.Count(&bind.CallOpts{
		From:    c.client.Addresses[0],
		Context: context.Background(),
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}

type EthereumAutomationSimpleLogCounterConsumer struct {
	client   *seth.Client
	consumer *simple_log_upkeep_counter_wrapper.SimpleLogUpkeepCounter
	address  *common.Address
}

func (v *EthereumAutomationSimpleLogCounterConsumer) Address() string {
	return v.address.Hex()
}

func (v *EthereumAutomationSimpleLogCounterConsumer) Start() error {
	return nil
}

func (v *EthereumAutomationSimpleLogCounterConsumer) Counter(ctx context.Context) (*big.Int, error) {
	return v.consumer.Counter(&bind.CallOpts{
		From:    v.client.MustGetRootKeyAddress(),
		Context: ctx,
	})
}

func DeployAutomationSimpleLogTriggerConsumerFromKey(client *seth.Client, isStreamsLookup bool, keyNum int) (KeeperConsumer, error) {
	abi, err := simple_log_upkeep_counter_wrapper.SimpleLogUpkeepCounterMetaData.GetAbi()
	if err != nil {
		return &EthereumAutomationSimpleLogCounterConsumer{}, fmt.Errorf("failed to get SimpleLogUpkeepCounter ABI: %w", err)
	}
	data, err := client.DeployContract(client.NewTXKeyOpts(keyNum), "SimpleLogUpkeepCounter", *abi, common.FromHex(simple_log_upkeep_counter_wrapper.SimpleLogUpkeepCounterMetaData.Bin), isStreamsLookup)
	if err != nil {
		return &EthereumAutomationSimpleLogCounterConsumer{}, fmt.Errorf("SimpleLogUpkeepCounter instance deployment have failed: %w", err)
	}

	instance, err := simple_log_upkeep_counter_wrapper.NewSimpleLogUpkeepCounter(data.Address, MustNewWrappedContractBackend(nil, client))
	if err != nil {
		return &EthereumAutomationSimpleLogCounterConsumer{}, fmt.Errorf("failed to instantiate SimpleLogUpkeepCounter instance: %w", err)
	}

	return &EthereumAutomationSimpleLogCounterConsumer{
		client:   client,
		consumer: instance,
		address:  &data.Address,
	}, nil
}
