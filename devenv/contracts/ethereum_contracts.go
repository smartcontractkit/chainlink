package contracts

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/counter"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/i_automation_registry_master_wrapper_2_3"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/mock_ethusd_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/shared/generated/initial/weth9"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/functions/generated/functions_coordinator"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/functions/generated/functions_load_test_client"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/functions/generated/functions_router"
	iregistry22 "github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/i_automation_registry_master_wrapper_2_2"
	iregistry21 "github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/mock_ethlink_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/mock_gas_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/oracle_wrapper"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/test_api_consumer_wrapper"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/operatorforwarder/generated/operator"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/i_automation_registry_master_wrapper_2_2"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/keeper_registry_wrapper1_1"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/keeper_registry_wrapper1_2"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/keeper_registry_wrapper1_3"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/keeper_registry_wrapper2_0"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/operatorforwarder/generated/authorized_forwarder"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/operatorforwarder/generated/operator_factory"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	eth_contracts "github.com/smartcontractkit/chainlink/devenv/contracts/ethereum"
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

type EthereumFunctionsLoadStats struct {
	LastRequestID string
	LastResponse  string
	LastError     string
	Total         uint32
	Succeeded     uint32
	Errored       uint32
	Empty         uint32
}

func Bytes32ToSlice(a [32]byte) (r []byte) {
	r = append(r, a[:]...)
	return
}

// func ChainlinkK8sClientToChainlinkNodeWithKeysAndAddress(k8sNodes []*nodeclient.ChainlinkK8sClient) []ChainlinkNodeWithKeysAndAddress {
// 	var nodesAsInterface = make([]ChainlinkNodeWithKeysAndAddress, len(k8sNodes))
// 	for i, node := range k8sNodes {
// 		nodesAsInterface[i] = node
// 	}

// 	return nodesAsInterface
// }

// func ChainlinkClientToChainlinkNodeWithKeysAndAddress(k8sNodes []*nodeclient.ChainlinkClient) []ChainlinkNodeWithKeysAndAddress {
// 	var nodesAsInterface = make([]ChainlinkNodeWithKeysAndAddress, len(k8sNodes))
// 	for i, node := range k8sNodes {
// 		nodesAsInterface[i] = node
// 	}

// 	return nodesAsInterface
// }

// func V2OffChainAgrregatorToOffChainAggregatorWithRounds(contracts []OffchainAggregatorV2) []OffChainAggregatorWithRounds {
// 	var contractsAsInterface = make([]OffChainAggregatorWithRounds, len(contracts))
// 	for i, contract := range contracts {
// 		contractsAsInterface[i] = contract
// 	}

// 	return contractsAsInterface
// }

// func V1OffChainAgrregatorToOffChainAggregatorWithRounds(contracts []OffchainAggregator) []OffChainAggregatorWithRounds {
// 	var contractsAsInterface = make([]OffChainAggregatorWithRounds, len(contracts))
// 	for i, contract := range contracts {
// 		contractsAsInterface[i] = contract
// 	}

// 	return contractsAsInterface
// }

func GetRegistryContractABI(version eth_contracts.KeeperRegistryVersion) (*abi.ABI, error) {
	var (
		contractABI *abi.ABI
		err         error
	)
	switch version {
	case eth_contracts.RegistryVersion_1_0, eth_contracts.RegistryVersion_1_1:
		contractABI, err = keeper_registry_wrapper1_1.KeeperRegistryMetaData.GetAbi()
	case eth_contracts.RegistryVersion_1_2:
		contractABI, err = keeper_registry_wrapper1_2.KeeperRegistryMetaData.GetAbi()
	case eth_contracts.RegistryVersion_1_3:
		contractABI, err = keeper_registry_wrapper1_3.KeeperRegistryMetaData.GetAbi()
	case eth_contracts.RegistryVersion_2_0:
		contractABI, err = keeper_registry_wrapper2_0.KeeperRegistryMetaData.GetAbi()
	case eth_contracts.RegistryVersion_2_1:
		contractABI, err = iregistry21.IKeeperRegistryMasterMetaData.GetAbi()
	case eth_contracts.RegistryVersion_2_2:
		contractABI, err = iregistry22.IAutomationRegistryMasterMetaData.GetAbi()
	default:
		contractABI, err = keeper_registry_wrapper2_0.KeeperRegistryMetaData.GetAbi()
	}

	return contractABI, err
}

// LegacyEthereumOperatorFactory represents operator factory contract
type EthereumOperatorFactory struct {
	address         *common.Address
	client          *seth.Client
	operatorFactory *operator_factory.OperatorFactory
}

func DeployEthereumOperatorFactory(seth *seth.Client, linkTokenAddress common.Address) (EthereumOperatorFactory, error) {
	abi, err := operator_factory.OperatorFactoryMetaData.GetAbi()
	if err != nil {
		return EthereumOperatorFactory{}, fmt.Errorf("failed to get OperatorFactory ABI: %w", err)
	}
	operatorData, err := seth.DeployContract(seth.NewTXOpts(), "OperatorFactory", *abi, common.FromHex(operator_factory.OperatorFactoryMetaData.Bin), linkTokenAddress)
	if err != nil {
		return EthereumOperatorFactory{}, fmt.Errorf("OperatorFactory instance deployment have failed: %w", err)
	}

	operatorFactory, err := operator_factory.NewOperatorFactory(operatorData.Address, seth.Client)
	if err != nil {
		return EthereumOperatorFactory{}, fmt.Errorf("failed to instantiate OperatorFactory instance: %w", err)
	}

	return EthereumOperatorFactory{
		address:         &operatorData.Address,
		client:          seth,
		operatorFactory: operatorFactory,
	}, nil
}

func (e *EthereumOperatorFactory) ParseAuthorizedForwarderCreated(eventLog types.Log) (*operator_factory.OperatorFactoryAuthorizedForwarderCreated, error) {
	return e.operatorFactory.ParseAuthorizedForwarderCreated(eventLog)
}

func (e *EthereumOperatorFactory) ParseOperatorCreated(eventLog types.Log) (*operator_factory.OperatorFactoryOperatorCreated, error) {
	return e.operatorFactory.ParseOperatorCreated(eventLog)
}

func (e *EthereumOperatorFactory) Address() string {
	return e.address.Hex()
}

func (e *EthereumOperatorFactory) DeployNewOperatorAndForwarder() (*types.Transaction, error) {
	return e.operatorFactory.DeployNewOperatorAndForwarder(e.client.NewTXOpts())
}

// EthereumOperator represents operator contract
type EthereumOperator struct {
	address  *common.Address
	client   *seth.Client
	operator *operator.Operator
	l        zerolog.Logger
}

func LoadEthereumOperator(logger zerolog.Logger, seth *seth.Client, contractAddress common.Address) (EthereumOperator, error) {
	abi, err := operator.OperatorMetaData.GetAbi()
	if err != nil {
		return EthereumOperator{}, err
	}
	seth.ContractStore.AddABI("EthereumOperator", *abi)
	seth.ContractStore.AddBIN("EthereumOperator", common.FromHex(operator.OperatorMetaData.Bin))

	operator, err := operator.NewOperator(contractAddress, seth.Client)
	if err != nil {
		return EthereumOperator{}, err
	}

	return EthereumOperator{
		address:  &contractAddress,
		client:   seth,
		operator: operator,
		l:        logger,
	}, nil
}

func (e *EthereumOperator) Address() string {
	return e.address.Hex()
}

func (e *EthereumOperator) AcceptAuthorizedReceivers(forwarders []common.Address, eoa []common.Address) error {
	e.l.Info().
		Str("ForwardersAddresses", fmt.Sprint(forwarders)).
		Str("EoaAddresses", fmt.Sprint(eoa)).
		Msg("Accepting Authorized Receivers")
	_, err := e.client.Decode(e.operator.AcceptAuthorizedReceivers(e.client.NewTXOpts(), forwarders, eoa))
	return err
}

// EthereumAuthorizedForwarder represents authorized forwarder contract
type EthereumAuthorizedForwarder struct {
	address             *common.Address
	client              *seth.Client
	authorizedForwarder *authorized_forwarder.AuthorizedForwarder
}

func LoadEthereumAuthorizedForwarder(seth *seth.Client, contractAddress common.Address) (EthereumAuthorizedForwarder, error) {
	abi, err := authorized_forwarder.AuthorizedForwarderMetaData.GetAbi()
	if err != nil {
		return EthereumAuthorizedForwarder{}, err
	}
	seth.ContractStore.AddABI("AuthorizedForwarder", *abi)
	seth.ContractStore.AddBIN("AuthorizedForwarder", common.FromHex(authorized_forwarder.AuthorizedForwarderMetaData.Bin))

	authorizedForwarder, err := authorized_forwarder.NewAuthorizedForwarder(contractAddress, seth.Client)
	if err != nil {
		return EthereumAuthorizedForwarder{}, fmt.Errorf("failed to instantiate AuthorizedForwarder instance: %w", err)
	}

	return EthereumAuthorizedForwarder{
		address:             &contractAddress,
		client:              seth,
		authorizedForwarder: authorizedForwarder,
	}, nil
}

// Owner return authorized forwarder owner address
func (e *EthereumAuthorizedForwarder) Owner(_ context.Context) (string, error) {
	owner, err := e.authorizedForwarder.Owner(e.client.NewCallOpts())

	return owner.Hex(), err
}

func (e *EthereumAuthorizedForwarder) GetAuthorizedSenders(ctx context.Context) ([]string, error) {
	opts := &bind.CallOpts{
		From:    e.client.Addresses[0],
		Context: ctx,
	}
	authorizedSenders, err := e.authorizedForwarder.GetAuthorizedSenders(opts)
	if err != nil {
		return nil, err
	}
	var sendersAddrs []string
	for _, o := range authorizedSenders {
		sendersAddrs = append(sendersAddrs, o.Hex())
	}
	return sendersAddrs, nil
}

func (e *EthereumAuthorizedForwarder) Address() string {
	return e.address.Hex()
}

type EthereumOffchainAggregatorV2 struct {
	address  *common.Address
	client   *seth.Client
	contract *ocr2aggregator.OCR2Aggregator
	l        zerolog.Logger
}

func LoadOffchainAggregatorV2(l zerolog.Logger, seth *seth.Client, address common.Address) (EthereumOffchainAggregatorV2, error) {
	contractAbi, err := ocr2aggregator.OCR2AggregatorMetaData.GetAbi()
	if err != nil {
		return EthereumOffchainAggregatorV2{}, fmt.Errorf("failed to get OffChain Aggregator v2 ABI: %w", err)
	}
	seth.ContractStore.AddABI("OffChainAggregatorV2", *contractAbi)
	seth.ContractStore.AddBIN("OffChainAggregatorV2", common.FromHex(ocr2aggregator.OCR2AggregatorMetaData.Bin))

	ocr2, err := ocr2aggregator.NewOCR2Aggregator(address, seth.Client)
	if err != nil {
		return EthereumOffchainAggregatorV2{}, fmt.Errorf("failed to instantiate OCRv2 instance: %w", err)
	}

	return EthereumOffchainAggregatorV2{
		client:   seth,
		contract: ocr2,
		address:  &address,
		l:        l,
	}, nil
}

func DeployOffchainAggregatorV2(l zerolog.Logger, seth *seth.Client, linkTokenAddress common.Address, offchainOptions OffchainOptions) (EthereumOffchainAggregatorV2, error) {
	contractAbi, err := ocr2aggregator.OCR2AggregatorMetaData.GetAbi()
	if err != nil {
		return EthereumOffchainAggregatorV2{}, fmt.Errorf("failed to get OffChain Aggregator v2 ABI: %w", err)
	}
	seth.ContractStore.AddABI("OffChainAggregatorV2", *contractAbi)
	seth.ContractStore.AddBIN("OffChainAggregatorV2", common.FromHex(ocr2aggregator.OCR2AggregatorMetaData.Bin))

	ocrDeploymentData2, err := seth.DeployContract(seth.NewTXOpts(), "OffChainAggregatorV2", *contractAbi, common.FromHex(ocr2aggregator.OCR2AggregatorMetaData.Bin),
		linkTokenAddress,
		offchainOptions.MinimumAnswer,
		offchainOptions.MaximumAnswer,
		offchainOptions.BillingAccessController,
		offchainOptions.RequesterAccessController,
		offchainOptions.Decimals,
		offchainOptions.Description,
	)

	if err != nil {
		return EthereumOffchainAggregatorV2{}, fmt.Errorf("OCRv2 instance deployment have failed: %w", err)
	}

	ocr2, err := ocr2aggregator.NewOCR2Aggregator(ocrDeploymentData2.Address, MustNewWrappedContractBackend(nil, seth))
	if err != nil {
		return EthereumOffchainAggregatorV2{}, fmt.Errorf("failed to instantiate OCRv2 instance: %w", err)
	}

	return EthereumOffchainAggregatorV2{
		client:   seth,
		contract: ocr2,
		address:  &ocrDeploymentData2.Address,
		l:        l,
	}, nil
}

func (e *EthereumOffchainAggregatorV2) Address() string {
	return e.address.Hex()
}

func (e *EthereumOffchainAggregatorV2) RequestNewRound() error {
	_, err := e.client.Decode(e.contract.RequestNewRound(e.client.NewTXOpts()))
	return err
}

func (e *EthereumOffchainAggregatorV2) GetLatestAnswer(ctx context.Context) (*big.Int, error) {
	return e.contract.LatestAnswer(&bind.CallOpts{
		From:    e.client.Addresses[0],
		Context: ctx,
	})
}

func (e *EthereumOffchainAggregatorV2) GetLatestRound(ctx context.Context) (*RoundData, error) {
	data, err := e.contract.LatestRoundData(&bind.CallOpts{
		From:    e.client.Addresses[0],
		Context: ctx,
	})
	if err != nil {
		return nil, err
	}
	return &RoundData{
		RoundId:         data.RoundId,
		StartedAt:       data.StartedAt,
		UpdatedAt:       data.UpdatedAt,
		AnsweredInRound: data.AnsweredInRound,
		Answer:          data.Answer,
	}, nil
}

func (e *EthereumOffchainAggregatorV2) GetRound(ctx context.Context, roundID *big.Int) (*RoundData, error) {
	data, err := e.contract.GetRoundData(&bind.CallOpts{
		From:    e.client.Addresses[0],
		Context: ctx,
	}, roundID)
	if err != nil {
		return nil, err
	}
	return &RoundData{
		RoundId:         data.RoundId,
		StartedAt:       data.StartedAt,
		UpdatedAt:       data.UpdatedAt,
		AnsweredInRound: data.AnsweredInRound,
		Answer:          data.Answer,
	}, nil
}

func (e *EthereumOffchainAggregatorV2) SetPayees(transmitters, payees []string) error {
	e.l.Info().
		Str("Transmitters", fmt.Sprintf("%v", transmitters)).
		Str("Payees", fmt.Sprintf("%v", payees)).
		Str("OCRv2 Address", e.Address()).
		Msg("Setting OCRv2 Payees")

	var addTransmitters, addrPayees []common.Address
	for _, t := range transmitters {
		addTransmitters = append(addTransmitters, common.HexToAddress(t))
	}
	for _, p := range payees {
		addrPayees = append(addrPayees, common.HexToAddress(p))
	}

	_, err := e.client.Decode(e.contract.SetPayees(e.client.NewTXOpts(), addTransmitters, addrPayees))
	return err
}

func (e *EthereumOffchainAggregatorV2) SetConfig(ocrConfig *OCRv2Config) error {
	e.l.Info().
		Str("Address", e.Address()).
		Interface("Signers", ocrConfig.Signers).
		Interface("Transmitters", ocrConfig.Transmitters).
		Uint8("F", ocrConfig.F).
		Str("OnchainConfig", string(ocrConfig.OnchainConfig)).
		Uint64("OffchainConfigVersion", ocrConfig.OffchainConfigVersion).
		Str("OffchainConfig", string(ocrConfig.OffchainConfig)).
		Msg("Setting OCRv2 Config")

	_, err := e.client.Decode(e.contract.SetConfig(
		e.client.NewTXOpts(),
		ocrConfig.Signers,
		ocrConfig.Transmitters,
		ocrConfig.F,
		ocrConfig.OnchainConfig,
		ocrConfig.OffchainConfigVersion,
		ocrConfig.OffchainConfig,
	))
	return err
}

func (e *EthereumOffchainAggregatorV2) ParseEventAnswerUpdated(log types.Log) (*ocr2aggregator.OCR2AggregatorAnswerUpdated, error) {
	return e.contract.ParseAnswerUpdated(log)
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

func DeployOracle(seth *seth.Client, linkAddr string) (Oracle, error) {
	abi, err := oracle_wrapper.OracleMetaData.GetAbi()
	if err != nil {
		return &EthereumOracle{}, fmt.Errorf("failed to get Oracle ABI: %w", err)
	}
	seth.ContractStore.AddABI("Oracle", *abi)
	seth.ContractStore.AddBIN("Oracle", common.FromHex(oracle_wrapper.OracleMetaData.Bin))

	oracleDeploymentData, err := seth.DeployContract(seth.NewTXOpts(), "Oracle", *abi, common.FromHex(oracle_wrapper.OracleMetaData.Bin),
		common.HexToAddress(linkAddr),
	)

	if err != nil {
		return &EthereumOracle{}, fmt.Errorf("Oracle instance deployment have failed: %w", err)
	}

	oracle, err := oracle_wrapper.NewOracle(oracleDeploymentData.Address, MustNewWrappedContractBackend(nil, seth))
	if err != nil {
		return &EthereumOracle{}, fmt.Errorf("Oracle to instantiate FluxAggregator instance: %w", err)
	}

	return &EthereumOracle{
		client:  seth,
		address: &oracleDeploymentData.Address,
		oracle:  oracle,
	}, nil
}

// EthereumOracle oracle for "directrequest" job tests
type EthereumOracle struct {
	address *common.Address
	client  *seth.Client
	oracle  *oracle_wrapper.Oracle
}

func (e *EthereumOracle) Address() string {
	return e.address.Hex()
}

func (e *EthereumOracle) Fund(_ *big.Float) error {
	panic("do not use this function, use actions.SendFunds() instead, otherwise we will have to deal with circular dependencies")
}

// SetFulfillmentPermission sets fulfillment permission for particular address
func (e *EthereumOracle) SetFulfillmentPermission(address string, allowed bool) error {
	_, err := e.client.Decode(e.oracle.SetFulfillmentPermission(e.client.NewTXOpts(), common.HexToAddress(address), allowed))
	return err
}

func DeployAPIConsumer(seth *seth.Client, linkAddr string) (APIConsumer, error) {
	abi, err := test_api_consumer_wrapper.TestAPIConsumerMetaData.GetAbi()
	if err != nil {
		return &EthereumAPIConsumer{}, fmt.Errorf("failed to get TestAPIConsumer ABI: %w", err)
	}
	seth.ContractStore.AddABI("TestAPIConsumer", *abi)
	seth.ContractStore.AddBIN("TestAPIConsumer", common.FromHex(test_api_consumer_wrapper.TestAPIConsumerMetaData.Bin))

	consumerDeploymentData, err := seth.DeployContract(seth.NewTXOpts(), "TestAPIConsumer", *abi, common.FromHex(test_api_consumer_wrapper.TestAPIConsumerMetaData.Bin),
		common.HexToAddress(linkAddr),
	)

	if err != nil {
		return &EthereumAPIConsumer{}, fmt.Errorf("TestAPIConsumer instance deployment have failed: %w", err)
	}

	consumer, err := test_api_consumer_wrapper.NewTestAPIConsumer(consumerDeploymentData.Address, MustNewWrappedContractBackend(nil, seth))
	if err != nil {
		return &EthereumAPIConsumer{}, fmt.Errorf("failed to instantiate TestAPIConsumer instance: %w", err)
	}

	return &EthereumAPIConsumer{
		client:   seth,
		address:  &consumerDeploymentData.Address,
		consumer: consumer,
	}, nil
}

// EthereumAPIConsumer API consumer for job type "directrequest" tests
type EthereumAPIConsumer struct {
	address  *common.Address
	client   *seth.Client
	consumer *test_api_consumer_wrapper.TestAPIConsumer
}

func (e *EthereumAPIConsumer) Address() string {
	return e.address.Hex()
}

func (e *EthereumAPIConsumer) RoundID(ctx context.Context) (*big.Int, error) {
	return e.consumer.CurrentRoundID(&bind.CallOpts{
		From:    e.client.Addresses[0],
		Context: ctx,
	})
}

func (e *EthereumAPIConsumer) Fund(_ *big.Float) error {
	panic("do not use this function, use actions.SendFunds() instead, otherwise we will have to deal with circular dependencies")
}

func (e *EthereumAPIConsumer) Data(ctx context.Context) (*big.Int, error) {
	return e.consumer.Data(&bind.CallOpts{
		From:    e.client.Addresses[0],
		Context: ctx,
	})
}

// CreateRequestTo creates request to an oracle for particular jobID with params
func (e *EthereumAPIConsumer) CreateRequestTo(
	oracleAddr string,
	jobID [32]byte,
	payment *big.Int,
	url string,
	path string,
	times *big.Int,
) error {
	_, err := e.client.Decode(e.consumer.CreateRequestTo(e.client.NewTXOpts(), common.HexToAddress(oracleAddr), jobID, payment, url, path, times))
	return err
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

func LoadFunctionsCoordinator(seth *seth.Client, addr string) (FunctionsCoordinator, error) {
	abi, err := functions_coordinator.FunctionsCoordinatorMetaData.GetAbi()
	if err != nil {
		return &EthereumFunctionsCoordinator{}, fmt.Errorf("failed to get FunctionsCoordinator ABI: %w", err)
	}
	seth.ContractStore.AddABI("FunctionsCoordinator", *abi)
	seth.ContractStore.AddBIN("FunctionsCoordinator", common.FromHex(functions_coordinator.FunctionsCoordinatorMetaData.Bin))

	instance, err := functions_coordinator.NewFunctionsCoordinator(common.HexToAddress(addr), seth.Client)
	if err != nil {
		return &EthereumFunctionsCoordinator{}, fmt.Errorf("failed to instantiate FunctionsCoordinator instance: %w", err)
	}

	return &EthereumFunctionsCoordinator{
		client:   seth,
		instance: instance,
		address:  common.HexToAddress(addr),
	}, err
}

type EthereumFunctionsCoordinator struct {
	address  common.Address
	client   *seth.Client
	instance *functions_coordinator.FunctionsCoordinator
}

func (e *EthereumFunctionsCoordinator) GetThresholdPublicKey() ([]byte, error) {
	return e.instance.GetThresholdPublicKey(e.client.NewCallOpts())
}

func (e *EthereumFunctionsCoordinator) GetDONPublicKey() ([]byte, error) {
	return e.instance.GetDONPublicKey(e.client.NewCallOpts())
}

func (e *EthereumFunctionsCoordinator) Address() string {
	return e.address.Hex()
}

func LoadFunctionsRouter(l zerolog.Logger, seth *seth.Client, addr string) (FunctionsRouter, error) {
	abi, err := functions_router.FunctionsRouterMetaData.GetAbi()
	if err != nil {
		return &EthereumFunctionsRouter{}, fmt.Errorf("failed to get FunctionsRouter ABI: %w", err)
	}
	seth.ContractStore.AddABI("FunctionsRouter", *abi)
	seth.ContractStore.AddBIN("FunctionsRouter", common.FromHex(functions_router.FunctionsRouterMetaData.Bin))

	instance, err := functions_router.NewFunctionsRouter(common.HexToAddress(addr), seth.Client)
	if err != nil {
		return &EthereumFunctionsRouter{}, fmt.Errorf("failed to instantiate FunctionsRouter instance: %w", err)
	}

	return &EthereumFunctionsRouter{
		client:   seth,
		instance: instance,
		address:  common.HexToAddress(addr),
		l:        l,
	}, err
}

type EthereumFunctionsRouter struct {
	address  common.Address
	client   *seth.Client
	instance *functions_router.FunctionsRouter
	l        zerolog.Logger
}

func (e *EthereumFunctionsRouter) Address() string {
	return e.address.Hex()
}

func (e *EthereumFunctionsRouter) CreateSubscriptionWithConsumer(consumer string) (uint64, error) {
	tx, err := e.client.Decode(e.instance.CreateSubscriptionWithConsumer(e.client.NewTXOpts(), common.HexToAddress(consumer)))
	if err != nil {
		return 0, err
	}

	if tx.Receipt == nil {
		return 0, errors.New("transaction did not err, but the receipt is nil")
	}
	for _, l := range tx.Receipt.Logs {
		e.l.Info().Interface("Log", common.Bytes2Hex(l.Data)).Send()
	}
	topicsMap := map[string]any{}

	fabi, err := abi.JSON(strings.NewReader(functions_router.FunctionsRouterABI))
	if err != nil {
		return 0, err
	}
	for _, ev := range fabi.Events {
		e.l.Info().Str("EventName", ev.Name).Send()
	}
	topicOneInputs := abi.Arguments{fabi.Events["SubscriptionCreated"].Inputs[0]}
	topicOneHash := []common.Hash{tx.Receipt.Logs[0].Topics[1:][0]}
	if err := abi.ParseTopicsIntoMap(topicsMap, topicOneInputs, topicOneHash); err != nil {
		return 0, fmt.Errorf("failed to decode topic value, err: %w", err)
	}
	e.l.Info().Interface("NewTopicsDecoded", topicsMap).Send()
	if topicsMap["subscriptionId"] == 0 {
		return 0, errors.New("failed to decode subscription ID after creation")
	}
	return topicsMap["subscriptionId"].(uint64), nil
}

func DeployFunctionsLoadTestClient(seth *seth.Client, router string) (FunctionsLoadTestClient, error) {
	operatorAbi, err := functions_load_test_client.FunctionsLoadTestClientMetaData.GetAbi()
	if err != nil {
		return &EthereumFunctionsLoadTestClient{}, fmt.Errorf("failed to get FunctionsLoadTestClient ABI: %w", err)
	}
	data, err := seth.DeployContract(seth.NewTXOpts(), "FunctionsLoadTestClient", *operatorAbi, common.FromHex(functions_load_test_client.FunctionsLoadTestClientMetaData.Bin), common.HexToAddress(router))
	if err != nil {
		return &EthereumFunctionsLoadTestClient{}, fmt.Errorf("FunctionsLoadTestClient instance deployment have failed: %w", err)
	}

	instance, err := functions_load_test_client.NewFunctionsLoadTestClient(data.Address, seth.Client)
	if err != nil {
		return &EthereumFunctionsLoadTestClient{}, fmt.Errorf("failed to instantiate FunctionsLoadTestClient instance: %w", err)
	}

	return &EthereumFunctionsLoadTestClient{
		client:   seth,
		instance: instance,
		address:  data.Address,
	}, nil
}

// LoadFunctionsLoadTestClient returns deployed on given address FunctionsLoadTestClient contract instance
func LoadFunctionsLoadTestClient(seth *seth.Client, addr string) (FunctionsLoadTestClient, error) {
	abi, err := functions_load_test_client.FunctionsLoadTestClientMetaData.GetAbi()
	if err != nil {
		return &EthereumFunctionsLoadTestClient{}, fmt.Errorf("failed to get FunctionsLoadTestClient ABI: %w", err)
	}
	seth.ContractStore.AddABI("FunctionsLoadTestClient", *abi)
	seth.ContractStore.AddBIN("FunctionsLoadTestClient", common.FromHex(functions_load_test_client.FunctionsLoadTestClientMetaData.Bin))

	instance, err := functions_load_test_client.NewFunctionsLoadTestClient(common.HexToAddress(addr), seth.Client)
	if err != nil {
		return &EthereumFunctionsLoadTestClient{}, fmt.Errorf("failed to instantiate FunctionsLoadTestClient instance: %w", err)
	}

	return &EthereumFunctionsLoadTestClient{
		client:   seth,
		instance: instance,
		address:  common.HexToAddress(addr),
	}, err
}

type EthereumFunctionsLoadTestClient struct {
	address  common.Address
	client   *seth.Client
	instance *functions_load_test_client.FunctionsLoadTestClient
}

func (e *EthereumFunctionsLoadTestClient) Address() string {
	return e.address.Hex()
}

func (e *EthereumFunctionsLoadTestClient) GetStats() (*EthereumFunctionsLoadStats, error) {
	lr, lbody, lerr, total, succeeded, errored, empty, err := e.instance.GetStats(e.client.NewCallOpts())
	if err != nil {
		return nil, err
	}
	return &EthereumFunctionsLoadStats{
		LastRequestID: string(Bytes32ToSlice(lr)),
		LastResponse:  string(lbody),
		LastError:     string(lerr),
		Total:         total,
		Succeeded:     succeeded,
		Errored:       errored,
		Empty:         empty,
	}, nil
}

func (e *EthereumFunctionsLoadTestClient) ResetStats() error {
	_, err := e.client.Decode(e.instance.ResetStats(e.client.NewTXOpts()))
	return err
}

func (e *EthereumFunctionsLoadTestClient) SendRequest(times uint32, source string, encryptedSecretsReferences []byte, args []string, subscriptionId uint64, jobId [32]byte) error {
	_, err := e.client.Decode(e.instance.SendRequest(e.client.NewTXOpts(), times, source, encryptedSecretsReferences, args, subscriptionId, jobId))
	return err
}

func (e *EthereumFunctionsLoadTestClient) SendRequestWithDONHostedSecrets(times uint32, source string, slotID uint8, slotVersion uint64, args []string, subscriptionId uint64, donID [32]byte) error {
	_, err := e.client.Decode(e.instance.SendRequestWithDONHostedSecrets(e.client.NewTXOpts(), times, source, slotID, slotVersion, args, subscriptionId, donID))
	return err
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
