package contracts

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/seth"

	"github.com/smartcontractkit/libocr/gethwrappers/offchainaggregator"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	ocrConfigHelper "github.com/smartcontractkit/libocr/offchainreporting/confighelper"
	ocrTypes "github.com/smartcontractkit/libocr/offchainreporting/types"

	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/authorized_forwarder"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/operator_factory"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/operator_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/link_token"
)

// EthereumOffchainAggregator represents the offchain aggregation contract
type EthereumOffchainAggregator struct {
	client  *seth.Client
	ocr     *offchainaggregator.OffchainAggregator
	address *common.Address
	l       zerolog.Logger
}

func LoadOffchainAggregator(l zerolog.Logger, seth *seth.Client, contractAddress common.Address) (EthereumOffchainAggregator, error) {
	oAbi, err := offchainaggregator.OffchainAggregatorMetaData.GetAbi()
	if err != nil {
		return EthereumOffchainAggregator{}, fmt.Errorf("failed to get OffChain Aggregator ABI: %w", err)
	}
	seth.ContractStore.AddABI("OffChainAggregator", *oAbi)
	seth.ContractStore.AddBIN("OffChainAggregator", common.FromHex(offchainaggregator.OffchainAggregatorMetaData.Bin))

	ocr, err := offchainaggregator.NewOffchainAggregator(contractAddress, seth.Client)
	if err != nil {
		return EthereumOffchainAggregator{}, fmt.Errorf("failed to instantiate OCR instance: %w", err)
	}

	return EthereumOffchainAggregator{
		client:  seth,
		ocr:     ocr,
		address: &contractAddress,
		l:       l,
	}, nil
}

func DeployOffchainAggregator(l zerolog.Logger, seth *seth.Client, linkTokenAddress common.Address, offchainOptions OffchainOptions) (EthereumOffchainAggregator, error) {
	oAbi, err := offchainaggregator.OffchainAggregatorMetaData.GetAbi()
	if err != nil {
		return EthereumOffchainAggregator{}, fmt.Errorf("failed to get OffChain Aggregator ABI: %w", err)
	}

	ocrDeploymentData, err := seth.DeployContract(
		seth.NewTXOpts(),
		"OffChainAggregator",
		*oAbi,
		common.FromHex(offchainaggregator.OffchainAggregatorMetaData.Bin),
		offchainOptions.MaximumGasPrice,
		offchainOptions.ReasonableGasPrice,
		offchainOptions.MicroLinkPerEth,
		offchainOptions.LinkGweiPerObservation,
		offchainOptions.LinkGweiPerTransmission,
		linkTokenAddress,
		offchainOptions.MinimumAnswer,
		offchainOptions.MaximumAnswer,
		offchainOptions.BillingAccessController,
		offchainOptions.RequesterAccessController,
		offchainOptions.Decimals,
		offchainOptions.Description)
	if err != nil {
		return EthereumOffchainAggregator{}, fmt.Errorf("OCR instance deployment have failed: %w", err)
	}

	ocr, err := offchainaggregator.NewOffchainAggregator(ocrDeploymentData.Address, seth.Client)
	if err != nil {
		return EthereumOffchainAggregator{}, fmt.Errorf("failed to instantiate OCR instance: %w", err)
	}

	return EthereumOffchainAggregator{
		client:  seth,
		ocr:     ocr,
		address: &ocrDeploymentData.Address,
		l:       l,
	}, nil
}

// SetPayees sets wallets for the contract to pay out to?
func (o *EthereumOffchainAggregator) SetPayees(
	transmitters, payees []string,
) error {
	var transmittersAddr, payeesAddr []common.Address
	for _, tr := range transmitters {
		transmittersAddr = append(transmittersAddr, common.HexToAddress(tr))
	}
	for _, p := range payees {
		payeesAddr = append(payeesAddr, common.HexToAddress(p))
	}

	o.l.Info().
		Str("Transmitters", fmt.Sprintf("%v", transmitters)).
		Str("Payees", fmt.Sprintf("%v", payees)).
		Str("OCR Address", o.Address()).
		Msg("Setting OCR Payees")

	_, err := o.client.Decode(o.ocr.SetPayees(o.client.NewTXOpts(), transmittersAddr, payeesAddr))
	return err
}

// SetConfig sets the payees and the offchain reporting protocol configuration
func (o *EthereumOffchainAggregator) SetConfig(
	chainlinkNodes []ChainlinkNodeWithKeysAndAddress,
	ocrConfig OffChainAggregatorConfig,
	transmitters []common.Address,
) error {
	// Gather necessary addresses and keys from our chainlink nodes to properly configure the OCR contract
	log.Info().Str("Contract Address", o.address.Hex()).Msg("Configuring OCR Contract")
	for i, node := range chainlinkNodes {
		ocrKeys, err := node.MustReadOCRKeys()
		if err != nil {
			return err
		}
		if len(ocrKeys.Data) == 0 {
			return fmt.Errorf("no OCR keys found for node %v", node)
		}
		primaryOCRKey := ocrKeys.Data[0]
		if err != nil {
			return err
		}
		p2pKeys, err := node.MustReadP2PKeys()
		if err != nil {
			return err
		}
		primaryP2PKey := p2pKeys.Data[0]

		// Need to convert the key representations
		var onChainSigningAddress [20]byte
		var configPublicKey [32]byte
		offchainSigningAddress, err := hex.DecodeString(primaryOCRKey.Attributes.OffChainPublicKey)
		if err != nil {
			return err
		}
		decodeConfigKey, err := hex.DecodeString(primaryOCRKey.Attributes.ConfigPublicKey)
		if err != nil {
			return err
		}

		// https://stackoverflow.com/questions/8032170/how-to-assign-string-to-bytes-array
		copy(onChainSigningAddress[:], common.HexToAddress(primaryOCRKey.Attributes.OnChainSigningAddress).Bytes())
		copy(configPublicKey[:], decodeConfigKey)

		oracleIdentity := ocrConfigHelper.OracleIdentity{
			TransmitAddress:       transmitters[i],
			OnChainSigningAddress: onChainSigningAddress,
			PeerID:                primaryP2PKey.Attributes.PeerID,
			OffchainPublicKey:     offchainSigningAddress,
		}
		oracleIdentityExtra := ocrConfigHelper.OracleIdentityExtra{
			OracleIdentity:                  oracleIdentity,
			SharedSecretEncryptionPublicKey: ocrTypes.SharedSecretEncryptionPublicKey(configPublicKey),
		}

		ocrConfig.OracleIdentities = append(ocrConfig.OracleIdentities, oracleIdentityExtra)
	}

	signers, transmitters, threshold, encodedConfigVersion, encodedConfig, err := ocrConfigHelper.ContractSetConfigArgs(
		ocrConfig.DeltaProgress,
		ocrConfig.DeltaResend,
		ocrConfig.DeltaRound,
		ocrConfig.DeltaGrace,
		ocrConfig.DeltaC,
		ocrConfig.AlphaPPB,
		ocrConfig.DeltaStage,
		ocrConfig.RMax,
		ocrConfig.S,
		ocrConfig.OracleIdentities,
		ocrConfig.F,
	)
	if err != nil {
		return err
	}

	// fails with error setting OCR config for contract '0x0DCd1Bf9A1b36cE34237eEaFef220932846BCD82': both gasPrice and (maxFeePerGas or maxPriorityFeePerGas) specified
	// but we only have gasPrice set... It also fails with the same error when we enable EIP-1559
	// fails when we wait for it to be minted, inside the wrapper there's no error when we call it, so it must be something inside smart contract
	// that's reverting it and maybe the error message is completely off
	_, err = o.client.Decode(o.ocr.SetConfig(o.client.NewTXOpts(), signers, transmitters, threshold, encodedConfigVersion, encodedConfig))
	return err
}

// RequestNewRound requests the OCR contract to create a new round
func (o *EthereumOffchainAggregator) RequestNewRound() error {
	o.l.Info().Str("Contract Address", o.address.Hex()).Msg("New OCR round requested")
	_, err := o.client.Decode(o.ocr.RequestNewRound(o.client.NewTXOpts()))
	return err
}

// GetLatestAnswer returns the latest answer from the OCR contract
func (o *EthereumOffchainAggregator) GetLatestAnswer(ctx context.Context) (*big.Int, error) {
	return o.ocr.LatestAnswer(&bind.CallOpts{
		From:    o.client.Addresses[0],
		Context: ctx,
	})
}

func (o *EthereumOffchainAggregator) Address() string {
	return o.address.Hex()
}

// GetLatestRound returns data from the latest round
func (o *EthereumOffchainAggregator) GetLatestRound(ctx context.Context) (*RoundData, error) {
	roundData, err := o.ocr.LatestRoundData(&bind.CallOpts{
		From:    o.client.Addresses[0],
		Context: ctx,
	})
	if err != nil {
		return nil, err
	}

	return &RoundData{
		RoundId:         roundData.RoundId,
		Answer:          roundData.Answer,
		AnsweredInRound: roundData.AnsweredInRound,
		StartedAt:       roundData.StartedAt,
		UpdatedAt:       roundData.UpdatedAt,
	}, err
}

func (o *EthereumOffchainAggregator) LatestRoundDataUpdatedAt() (*big.Int, error) {
	data, err := o.ocr.LatestRoundData(o.client.NewCallOpts())
	if err != nil {
		return nil, err
	}
	return data.UpdatedAt, nil
}

// GetRound retrieves an OCR round by the round ID
func (o *EthereumOffchainAggregator) GetRound(ctx context.Context, roundID *big.Int) (*RoundData, error) {
	roundData, err := o.ocr.GetRoundData(&bind.CallOpts{
		From:    o.client.Addresses[0],
		Context: ctx,
	}, roundID)
	if err != nil {
		return nil, err
	}

	return &RoundData{
		RoundId:         roundData.RoundId,
		Answer:          roundData.Answer,
		AnsweredInRound: roundData.AnsweredInRound,
		StartedAt:       roundData.StartedAt,
		UpdatedAt:       roundData.UpdatedAt,
	}, nil
}

// ParseEventAnswerUpdated parses the log for event AnswerUpdated
func (o *EthereumOffchainAggregator) ParseEventAnswerUpdated(eventLog types.Log) (*offchainaggregator.OffchainAggregatorAnswerUpdated, error) {
	return o.ocr.ParseAnswerUpdated(eventLog)
}

// LegacyEthereumOperatorFactory represents operator factory contract
type EthereumOperatorFactory struct {
	address         *common.Address
	client          *seth.Client
	operatorFactory *operator_factory.OperatorFactory
}

func DeployEthereumOperatorFactory(seth *seth.Client, linkTokenAddress common.Address) (EthereumOperatorFactory, error) {
	operatorAbi, err := operator_factory.OperatorFactoryMetaData.GetAbi()
	if err != nil {
		return EthereumOperatorFactory{}, fmt.Errorf("failed to get OperatorFactory ABI: %w", err)
	}
	operatorData, err := seth.DeployContract(seth.NewTXOpts(), "OperatorFactory", *operatorAbi, common.FromHex(operator_factory.OperatorFactoryMetaData.Bin), linkTokenAddress)
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
	operator *operator_wrapper.Operator
	l        zerolog.Logger
}

func LoadEthereumOperator(logger zerolog.Logger, seth *seth.Client, contractAddress common.Address) (EthereumOperator, error) {

	abi, err := operator_wrapper.OperatorMetaData.GetAbi()
	if err != nil {
		return EthereumOperator{}, err
	}
	seth.ContractStore.AddABI("EthereumOperator", *abi)
	seth.ContractStore.AddBIN("EthereumOperator", common.FromHex(operator_wrapper.OperatorMetaData.Bin))

	operator, err := operator_wrapper.NewOperator(contractAddress, seth.Client)
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

func LoadOffChainAggregatorV2(l zerolog.Logger, seth *seth.Client, contractAddress common.Address) (EthereumOffchainAggregatorV2, error) {
	oAbi, err := ocr2aggregator.OCR2AggregatorMetaData.GetAbi()
	if err != nil {
		return EthereumOffchainAggregatorV2{}, fmt.Errorf("failed to get OffChain Aggregator ABI: %w", err)
	}
	seth.ContractStore.AddABI("OffChainAggregatorV2", *oAbi)
	seth.ContractStore.AddBIN("OffChainAggregatorV2", common.FromHex(ocr2aggregator.OCR2AggregatorMetaData.Bin))

	ocr2, err := ocr2aggregator.NewOCR2Aggregator(contractAddress, seth.Client)
	if err != nil {
		return EthereumOffchainAggregatorV2{}, fmt.Errorf("failed to instantiate OCR instance: %w", err)
	}

	return EthereumOffchainAggregatorV2{
		client:   seth,
		contract: ocr2,
		address:  &contractAddress,
		l:        l,
	}, nil
}

func DeployOffchainAggregatorV2(l zerolog.Logger, seth *seth.Client, linkTokenAddress common.Address, offchainOptions OffchainOptions) (EthereumOffchainAggregatorV2, error) {
	oAbi, err := ocr2aggregator.OCR2AggregatorMetaData.GetAbi()
	if err != nil {
		return EthereumOffchainAggregatorV2{}, fmt.Errorf("failed to get OffChain Aggregator ABI: %w", err)
	}
	seth.ContractStore.AddABI("OffChainAggregatorV2", *oAbi)
	seth.ContractStore.AddBIN("OffChainAggregatorV2", common.FromHex(ocr2aggregator.OCR2AggregatorMetaData.Bin))

	ocrDeploymentData2, err := seth.DeployContract(seth.NewTXOpts(), "OffChainAggregatorV2", *oAbi, common.FromHex(ocr2aggregator.OCR2AggregatorMetaData.Bin),
		linkTokenAddress,
		offchainOptions.MinimumAnswer,
		offchainOptions.MaximumAnswer,
		offchainOptions.BillingAccessController,
		offchainOptions.RequesterAccessController,
		offchainOptions.Decimals,
		offchainOptions.Description,
	)

	if err != nil {
		return EthereumOffchainAggregatorV2{}, fmt.Errorf("OCR instance deployment have failed: %w", err)
	}

	ocr2, err := ocr2aggregator.NewOCR2Aggregator(ocrDeploymentData2.Address, seth.Client)
	if err != nil {
		return EthereumOffchainAggregatorV2{}, fmt.Errorf("failed to instantiate OCR instance: %w", err)
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
		Bytes("OnchainConfig", ocrConfig.OnchainConfig).
		Uint64("OffchainConfigVersion", ocrConfig.OffchainConfigVersion).
		Bytes("OffchainConfig", ocrConfig.OffchainConfig).
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

func DeployLinkTokenContract(client *seth.Client) (seth.DeploymentData, error) {
	linkTokenAbi, err := link_token.LinkTokenMetaData.GetAbi()
	if err != nil {
		return seth.DeploymentData{}, fmt.Errorf("failed to get LinkToken ABI: %w", err)
	}
	linkDeploymentData, err := client.DeployContract(client.NewTXOpts(), "LinkToken", *linkTokenAbi, common.FromHex(link_token.LinkTokenMetaData.Bin))
	if err != nil {
		return seth.DeploymentData{}, fmt.Errorf("LinkToken instance deployment have failed: %w", err)
	}

	return linkDeploymentData, nil
}
