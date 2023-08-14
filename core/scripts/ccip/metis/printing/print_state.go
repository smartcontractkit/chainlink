package printing

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/exp/slices"

	"github.com/smartcontractkit/chainlink/core/scripts/ccip/dione"
	"github.com/smartcontractkit/chainlink/core/scripts/ccip/rhea"
	"github.com/smartcontractkit/chainlink/core/scripts/ccip/rhea/deployments"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/arm_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/burn_mint_erc677"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/ping_pong_demo"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/price_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip"
)

func PrintCCIPState(source *rhea.EvmDeploymentConfig, destination *rhea.EvmDeploymentConfig) {
	printPoolBalances(source)
	printPoolBalances(destination)

	printSupportedTokensCheck(source, destination)

	printDappSanityCheck(source)
	printDappSanityCheck(destination)

	printRampSanityCheck(source, destination.LaneConfig.OnRamp, rhea.GetCCIPChainSelector(destination.ChainConfig.EvmChainId))
	printRampSanityCheck(destination, source.LaneConfig.OnRamp, rhea.GetCCIPChainSelector(source.ChainConfig.EvmChainId))

	checkPriceRegistrySet(source, destination)

	printPaused(source)
	printPaused(destination)

	printRateLimitingStatus(source)
	printRateLimitingStatus(destination)
}

func SetupAllLanesReadOnly(logger logger.Logger) {
	err := deployments.Prod_AvaxFujiToOptimismGoerli.SetupReadOnlyChain(logger.Named(ccip.ChainName(int64(deployments.Prod_AvaxFujiToOptimismGoerli.ChainConfig.EvmChainId))))
	if err != nil {
		panic(err)
	}
	err = deployments.Prod_OptimismGoerliToAvaxFuji.SetupReadOnlyChain(logger.Named(ccip.ChainName(int64(deployments.Prod_OptimismGoerliToAvaxFuji.ChainConfig.EvmChainId))))
	if err != nil {
		panic(err)
	}
}

func PrintTokenSupportAllChains(logger logger.Logger) {
	SetupAllLanesReadOnly(logger)
	err := deployments.Prod_SepoliaToOptimismGoerli.SetupReadOnlyChain(logger.Named(ccip.ChainName(int64(deployments.Prod_SepoliaToOptimismGoerli.ChainConfig.EvmChainId))))
	if err != nil {
		log.Fatal(err)
	}
	err = deployments.Prod_OptimismGoerliToSepolia.SetupReadOnlyChain(logger.Named(ccip.ChainName(int64(deployments.Prod_OptimismGoerliToSepolia.ChainConfig.EvmChainId))))
	if err != nil {
		log.Fatal(err)
	}
	err = deployments.Prod_SepoliaToAvaxFuji.SetupReadOnlyChain(logger.Named(ccip.ChainName(int64(deployments.Prod_SepoliaToAvaxFuji.ChainConfig.EvmChainId))))
	if err != nil {
		log.Fatal(err)
	}
	err = deployments.Prod_AvaxFujiToSepolia.SetupReadOnlyChain(logger.Named(ccip.ChainName(int64(deployments.Prod_AvaxFujiToSepolia.ChainConfig.EvmChainId))))
	if err != nil {
		log.Fatal(err)
	}
	err = deployments.Prod_AvaxFujiToOptimismGoerli.SetupReadOnlyChain(logger.Named(ccip.ChainName(int64(deployments.Prod_AvaxFujiToOptimismGoerli.ChainConfig.EvmChainId))))
	if err != nil {
		log.Fatal(err)
	}
	err = deployments.Prod_OptimismGoerliToAvaxFuji.SetupReadOnlyChain(logger.Named(ccip.ChainName(int64(deployments.Prod_OptimismGoerliToAvaxFuji.ChainConfig.EvmChainId))))
	if err != nil {
		log.Fatal(err)
	}

	printSupportedTokensCheck(&deployments.Prod_SepoliaToOptimismGoerli, &deployments.Prod_OptimismGoerliToSepolia)
	printSupportedTokensCheck(&deployments.Prod_SepoliaToAvaxFuji, &deployments.Prod_AvaxFujiToSepolia)
	printSupportedTokensCheck(&deployments.Prod_AvaxFujiToOptimismGoerli, &deployments.Prod_OptimismGoerliToAvaxFuji)
}

type CCIPTXStatus struct {
	message      *evm_2_evm_onramp.EVM2EVMOnRampCCIPSendRequested
	commitReport *commit_store.CommitStoreReportAccepted
	execStatus   *evm_2_evm_offramp.EVM2EVMOffRampExecutionStateChanged
}

type ExecutionStatus uint8

const (
	Untouched  ExecutionStatus = 0
	InProgress ExecutionStatus = 1
	Success    ExecutionStatus = 2
	Failed     ExecutionStatus = 3
)

func (e ExecutionStatus) String() string {
	switch e {
	case Untouched:
		return "Untouched"
	case InProgress:
		return "InProgress"
	case Success:
		return "Success"
	case Failed:
		return "Failed"
	default:
		return fmt.Sprintf("%d", int(e))
	}
}

func printBool(b bool) string {
	if b {
		return "✅"
	}
	return "❌"
}

func printBoolNeutral(b bool) string {
	if b {
		return "✅"
	}
	return "➖"
}

func PrintTxStatuses(source *rhea.EvmDeploymentConfig, destination *rhea.EvmDeploymentConfig) {
	onRamp, err := evm_2_evm_onramp.NewEVM2EVMOnRamp(source.LaneConfig.OnRamp, source.Client)
	helpers.PanicErr(err)

	block, err := source.Client.BlockNumber(context.Background())
	helpers.PanicErr(err)

	sendRequested, err := onRamp.FilterCCIPSendRequested(&bind.FilterOpts{
		Start: block - 9990,
	})
	helpers.PanicErr(err)

	txs := make(map[uint64]*CCIPTXStatus)
	maxSeqNum := uint64(0)
	minSeqNum := uint64(1)
	var seqNums []uint64

	for sendRequested.Next() {
		txs[sendRequested.Event.Message.SequenceNumber] = &CCIPTXStatus{
			message: sendRequested.Event,
		}
		if sendRequested.Event.Message.SequenceNumber > maxSeqNum {
			maxSeqNum = sendRequested.Event.Message.SequenceNumber
		}
		if minSeqNum == 1 {
			minSeqNum = sendRequested.Event.Message.SequenceNumber
		}
		seqNums = append(seqNums, sendRequested.Event.Message.SequenceNumber)
	}

	commitStore, err := commit_store.NewCommitStore(destination.LaneConfig.CommitStore, destination.Client)
	helpers.PanicErr(err)

	block, err = destination.Client.BlockNumber(context.Background())
	helpers.PanicErr(err)

	reports, err := commitStore.FilterReportAccepted(&bind.FilterOpts{
		Start: block - 9990,
	})
	helpers.PanicErr(err)

	for reports.Next() {
		for j := reports.Event.Report.Interval.Min; j <= reports.Event.Report.Interval.Max; j++ {
			if _, ok := txs[j]; ok {
				txs[j].commitReport = reports.Event
			}
		}
	}

	offRamp, err := evm_2_evm_offramp.NewEVM2EVMOffRamp(destination.LaneConfig.OffRamp, destination.Client)
	helpers.PanicErr(err)

	stateChanges, err := offRamp.FilterExecutionStateChanged(
		&bind.FilterOpts{
			Start: block - 9990,
		},
		seqNums,
		[][32]byte{})
	helpers.PanicErr(err)

	for stateChanges.Next() {
		if _, ok := txs[stateChanges.Event.SequenceNumber]; !ok {
			txs[stateChanges.Event.SequenceNumber] = &CCIPTXStatus{}
			if stateChanges.Event.SequenceNumber > maxSeqNum {
				maxSeqNum = stateChanges.Event.SequenceNumber
			}
		}
		txs[stateChanges.Event.SequenceNumber].execStatus = stateChanges.Event
	}

	var sb strings.Builder
	sb.WriteString("\n")
	tableHeaders := []string{"SequenceNumber", "Committed in block", "Execution status", "Executed in block", "Nonce"}
	headerLengths := []int{18, 18, 20, 18, 18}

	sb.WriteString(generateHeader(tableHeaders, headerLengths))

	if minSeqNum > 1 {
		sb.WriteString(fmt.Sprintf("| %18d | %18d | %41s | %18s | \n", 1, minSeqNum-1, "Probably > 10k blocks in the past", ""))
	}

	for i := minSeqNum; i <= maxSeqNum; i++ {
		tx := txs[i]
		committedAt := "-"
		if tx == nil {
			sb.WriteString(fmt.Sprintf("| %18d | %18s | %41s | %18s | \n", i, "TX MISSING", "", ""))
			continue
		}
		if tx.commitReport != nil {
			committedAt = strconv.Itoa(int(tx.commitReport.Raw.BlockNumber))
		}

		if tx.message == nil {
			sb.WriteString(fmt.Sprintf("| %18s | %18s | %20v | %18d | %18s | \n", "MISSING", committedAt, ExecutionStatus(tx.execStatus.State), tx.execStatus.Raw.BlockNumber, "-"))
		} else if tx.execStatus != nil {
			sb.WriteString(fmt.Sprintf("| %18d | %18s | %20v | %18d | %18d | %s \n",
				tx.message.Message.SequenceNumber,
				committedAt,
				ExecutionStatus(tx.execStatus.State),
				tx.execStatus.Raw.BlockNumber,
				tx.message.Message.Nonce,
				helpers.ExplorerLink(int64(destination.ChainConfig.EvmChainId), tx.execStatus.Raw.TxHash)))
		} else {
			sb.WriteString(fmt.Sprintf("| %18d | %18s | %20v | %18s | %18d | %s \n",
				tx.message.Message.SequenceNumber,
				committedAt,
				"-",
				"-",
				tx.message.Message.Nonce,
				""))
		}
	}
	sb.WriteString(generateSeparator(headerLengths))

	destination.Logger.Info(sb.String())
}

func printDappSanityCheck(source *rhea.EvmDeploymentConfig) {
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("Dapp sanity checks for %s\n", ccip.ChainName(int64(source.ChainConfig.EvmChainId))))

	tableHeaders := []string{"Dapp", "Router Set"}
	headerLengths := []int{30, 14}

	sb.WriteString(generateHeader(tableHeaders, headerLengths))

	if source.LaneConfig.PingPongDapp != common.HexToAddress("") {
		pingDapp, err := ping_pong_demo.NewPingPongDemo(source.LaneConfig.PingPongDapp, source.Client)
		helpers.PanicErr(err)
		router, err := pingDapp.GetRouter(&bind.CallOpts{})
		helpers.PanicErr(err)
		sb.WriteString(fmt.Sprintf("| %-30s | %14s |\n", "Ping dapp sender", printBool(router == source.ChainConfig.Router)))
	}

	sb.WriteString(generateSeparator(headerLengths))

	source.Logger.Info(sb.String())
}

func printRampSanityCheck(chain *rhea.EvmDeploymentConfig, sourceOnRamp common.Address, remoteChainSelector uint64) {
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("Ramp checks for %s\n", ccip.ChainName(int64(chain.ChainConfig.EvmChainId))))

	tableHeaders := []string{"Contract", "Config correct"}
	headerLengths := []int{30, 14}

	sb.WriteString(generateHeader(tableHeaders, headerLengths))

	arm, err := arm_contract.NewARMContract(chain.ChainConfig.ARM, chain.Client)
	helpers.PanicErr(err)
	badSignal, err := arm.IsCursed(&bind.CallOpts{})
	helpers.PanicErr(err)

	sb.WriteString(fmt.Sprintf("| %-30s | %14s |\n", "ARM healthy", printBool(!badSignal)))

	onRamp, err := evm_2_evm_onramp.NewEVM2EVMOnRamp(chain.LaneConfig.OnRamp, chain.Client)
	helpers.PanicErr(err)
	dynamicOnRampConfig, err := onRamp.GetDynamicConfig(&bind.CallOpts{})
	helpers.PanicErr(err)
	staticOnRampConfig, err := onRamp.GetStaticConfig(&bind.CallOpts{})
	helpers.PanicErr(err)
	sb.WriteString(fmt.Sprintf("| %-30s | %14s |\n", "OnRamp Router set", printBool(dynamicOnRampConfig.Router == chain.ChainConfig.Router)))
	sb.WriteString(fmt.Sprintf("| %-30s | %14s |\n", "OnRamp chainSelector valid",
		printBool(staticOnRampConfig.ChainSelector == rhea.GetCCIPChainSelector(chain.ChainConfig.EvmChainId))))
	sb.WriteString(fmt.Sprintf("| %-30s | %14s |\n", "OnRamp destChainSelector valid",
		printBool(staticOnRampConfig.DestChainSelector == remoteChainSelector)))

	offRamp, err := evm_2_evm_offramp.NewEVM2EVMOffRamp(chain.LaneConfig.OffRamp, chain.Client)
	helpers.PanicErr(err)
	dynamicOffRampConfig, err := offRamp.GetDynamicConfig(&bind.CallOpts{})
	helpers.PanicErr(err)
	staticOffRampConfig, err := offRamp.GetStaticConfig(&bind.CallOpts{})
	helpers.PanicErr(err)

	sb.WriteString(fmt.Sprintf("| %-30s | %14s |\n", "OffRamp Router set", printBool(dynamicOffRampConfig.Router == chain.ChainConfig.Router)))
	sb.WriteString(fmt.Sprintf("| %-30s | %14s |\n", "OffRamp chainSelector valid",
		printBool(staticOffRampConfig.ChainSelector == rhea.GetCCIPChainSelector(chain.ChainConfig.EvmChainId))))
	sb.WriteString(fmt.Sprintf("| %-30s | %14s |\n", "OffRamp sourceChainSelector valid",
		printBool(staticOffRampConfig.SourceChainSelector == remoteChainSelector)))

	configDetails, err := offRamp.LatestConfigDetails(&bind.CallOpts{})
	helpers.PanicErr(err)
	sb.WriteString(fmt.Sprintf("| %-30s | %14s |\n", "OffRamp OCR2 configured", printBool(configDetails.ConfigCount != 0)))

	commitStore, err := commit_store.NewCommitStore(chain.LaneConfig.CommitStore, chain.Client)
	helpers.PanicErr(err)

	blobConfigDetails, err := commitStore.LatestConfigDetails(&bind.CallOpts{})
	helpers.PanicErr(err)
	sb.WriteString(fmt.Sprintf("| %-30s | %14s |\n", "CommitStore OCR2 configured", printBool(blobConfigDetails.ConfigCount != 0)))

	router, err := router.NewRouter(chain.ChainConfig.Router, chain.Client)
	helpers.PanicErr(err)

	offRamps, err := router.GetOffRamps(&bind.CallOpts{})
	helpers.PanicErr(err)

	isRamp := false
	for _, ramp := range offRamps {
		if ramp.OffRamp == chain.LaneConfig.OffRamp {
			isRamp = true
			break
		}
	}

	sb.WriteString(fmt.Sprintf("| %-30s | %14s |\n", "Router has offRamp Set", printBool(isRamp)))

	sb.WriteString(generateSeparator(headerLengths))

	chain.Logger.Info(sb.String())
}

func printRateLimitingStatus(chain *rhea.EvmDeploymentConfig) {
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("Rate limits for %s\n", ccip.ChainName(int64(chain.ChainConfig.EvmChainId))))

	tableHeaders := []string{"Contract", "Tokens"}
	headerLengths := []int{25, 42}

	sb.WriteString(generateHeader(tableHeaders, headerLengths))

	onRamp, err := evm_2_evm_onramp.NewEVM2EVMOnRamp(chain.LaneConfig.OnRamp, chain.Client)
	helpers.PanicErr(err)
	onRampRateLimiterState, err := onRamp.CurrentRateLimiterState(&bind.CallOpts{})
	helpers.PanicErr(err)

	sb.WriteString(fmt.Sprintf("| %-25s | %42d |\n", "onramp", onRampRateLimiterState.Tokens))

	offRamp, err := evm_2_evm_offramp.NewEVM2EVMOffRamp(chain.LaneConfig.OffRamp, chain.Client)
	helpers.PanicErr(err)
	offRampRateLimiterState, err := offRamp.CurrentRateLimiterState(&bind.CallOpts{})
	helpers.PanicErr(err)

	sb.WriteString(fmt.Sprintf("| %-25s | %42d |\n", "offramp", offRampRateLimiterState.Tokens))

	sb.WriteString(generateSeparator(headerLengths))
	chain.Logger.Info(sb.String())
}

func printPaused(chain *rhea.EvmDeploymentConfig) {
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("Paused addresses for %s\n", ccip.ChainName(int64(chain.ChainConfig.EvmChainId))))

	tableHeaders := []string{"Contract", "Address", "Running"}
	headerLengths := []int{25, 42, 14}

	sb.WriteString(generateHeader(tableHeaders, headerLengths))

	for _, tokenConfig := range chain.ChainConfig.SupportedTokens {
		if tokenConfig.Pool == common.HexToAddress("") {
			continue
		}
		sb.WriteString(fmt.Sprintf("| %-25s | %42s |\n", "token pool", tokenConfig.Pool.Hex()))
	}

	commitStore, err := commit_store.NewCommitStore(chain.LaneConfig.CommitStore, chain.Client)
	helpers.PanicErr(err)
	paused, err := commitStore.Paused(&bind.CallOpts{})
	helpers.PanicErr(err)

	sb.WriteString(fmt.Sprintf("| %-25s | %42s | %14s |\n", "commitStore", commitStore.Address(), printBool(!paused)))

	sb.WriteString(generateSeparator(headerLengths))
	chain.Logger.Info(sb.String())
}

func PrintNodeBalances(chain *rhea.EvmDeploymentConfig, addresses []common.Address) {
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("Node balances for %s\n", ccip.ChainName(int64(chain.ChainConfig.EvmChainId))))

	tableHeaders := []string{"Sender", "Balance"}
	headerLengths := []int{42, 18}

	sb.WriteString(generateHeader(tableHeaders, headerLengths))

	for _, sender := range addresses {
		balanceAt, err := chain.Client.BalanceAt(context.Background(), sender, nil)
		helpers.PanicErr(err)

		sb.WriteString(fmt.Sprintf("| %42s |   %-16s |\n", sender.Hex(), new(big.Float).Quo(new(big.Float).SetInt(balanceAt), big.NewFloat(1e18)).String()))
	}

	sb.WriteString(generateSeparator(headerLengths))
	chain.Logger.Info(sb.String())
}

func printPoolBalances(chain *rhea.EvmDeploymentConfig) {
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("Pool balances for %s\n", ccip.ChainName(int64(chain.ChainConfig.EvmChainId))))

	tableHeaders := []string{"Token", "Pool", "Balance", "Onramp", "OffRamp", "Price"}
	headerLengths := []int{32, 42, 20, 9, 9, 10}

	sb.WriteString(generateHeader(tableHeaders, headerLengths))

	priceRegistry, err := price_registry.NewPriceRegistry(chain.ChainConfig.PriceRegistry, chain.Client)
	helpers.PanicErr(err)

	for tokenName, tokenConfig := range chain.ChainConfig.SupportedTokens {
		if tokenConfig.Pool == common.HexToAddress("") {
			sb.WriteString(fmt.Sprintf("| %-32s | No pool found\n", tokenName))
			continue
		}

		tokenPool, err := lock_release_token_pool.NewLockReleaseTokenPool(tokenConfig.Pool, chain.Client)
		helpers.PanicErr(err)

		tokenAddress, err := tokenPool.GetToken(&bind.CallOpts{})
		helpers.PanicErr(err)

		tokenInstance, err := burn_mint_erc677.NewBurnMintERC677(tokenAddress, chain.Client)
		helpers.PanicErr(err)

		tokenPrice, err := priceRegistry.GetTokenPrice(&bind.CallOpts{}, tokenAddress)
		helpers.PanicErr(err)

		balance, err := tokenInstance.BalanceOf(&bind.CallOpts{}, tokenConfig.Pool)
		helpers.PanicErr(err)

		isAllowedOnRamp, err := tokenPool.IsOnRamp(&bind.CallOpts{}, chain.LaneConfig.OnRamp)
		helpers.PanicErr(err)

		isAllowedOffRamp, err := tokenPool.IsOffRamp(&bind.CallOpts{}, chain.LaneConfig.OffRamp)
		helpers.PanicErr(err)

		if tokenAddress != tokenConfig.Token {
			sb.WriteString(fmt.Sprintf("| %-32s | TOKEN CONFIG MISMATCH ❌ | expected %s | pool token %s |\n", tokenName, tokenConfig.Token.Hex(), tokenAddress.Hex()))
		} else {
			sb.WriteString(fmt.Sprintf("| %-32s | %s | %20d | %9s | %9s | %10s |\n", tokenName, tokenConfig.Pool, balance, printBool(isAllowedOnRamp), printBool(isAllowedOffRamp), tokenPrice.Value.String()))
		}
	}

	sb.WriteString(generateSeparator(headerLengths))

	chain.Logger.Info(sb.String())
}

func printSupportedTokensCheck(source *rhea.EvmDeploymentConfig, destination *rhea.EvmDeploymentConfig) {
	sourceRouter, err := router.NewRouter(source.ChainConfig.Router, source.Client)
	helpers.PanicErr(err)

	sourceTokens, err := sourceRouter.GetSupportedTokens(&bind.CallOpts{}, rhea.GetCCIPChainSelector(destination.ChainConfig.EvmChainId))
	helpers.PanicErr(err)

	destRouter, err := router.NewRouter(destination.ChainConfig.Router, destination.Client)
	helpers.PanicErr(err)

	destTokens, err := destRouter.GetSupportedTokens(&bind.CallOpts{}, rhea.GetCCIPChainSelector(source.ChainConfig.EvmChainId))
	helpers.PanicErr(err)

	var sb strings.Builder
	sb.WriteString("\nToken matching\n")

	tableHeaders := []string{"", "Source " + ccip.ChainName(int64(source.ChainConfig.EvmChainId)), "Destination " + ccip.ChainName(int64(destination.ChainConfig.EvmChainId))}
	headerLengths := []int{20, 49, 49}

	sb.WriteString(generateSeparator(headerLengths))
	sb.WriteString("|")
	for i, header := range tableHeaders {
		sb.WriteString(fmt.Sprintf(" %-"+strconv.Itoa(headerLengths[i])+"s |", header))
	}
	sb.WriteString("\n")

	tableHeaders = []string{"Token", "FeeToken", "FeeAmount", "Transfer", "Pool", "FeeToken", "FeeAmount", "Transfer", "Pool"}
	headerLengths = []int{20, 9, 13, 9, 9, 9, 13, 9, 9}

	sb.WriteString(generateHeader(tableHeaders, headerLengths))

	for _, token := range rhea.GetAllTokens() {
		var sourceEnabled, isSourcePool, isSourceFeeToken bool
		sourceFeeAmount := "➖"

		if _, ok := source.ChainConfig.SupportedTokens[token]; ok {
			sourceEnabled = slices.Contains(sourceTokens, source.ChainConfig.SupportedTokens[token].Token)
			isSourcePool = source.ChainConfig.SupportedTokens[token].Pool != common.HexToAddress("")

			sourceFee, err := sourceRouter.GetFee(&bind.CallOpts{}, rhea.GetCCIPChainSelector(destination.ChainConfig.EvmChainId), router.ClientEVM2AnyMessage{
				Receiver:     common.HexToAddress("").Bytes(),
				Data:         []byte{},
				TokenAmounts: []router.ClientEVMTokenAmount{},
				FeeToken:     source.ChainConfig.SupportedTokens[token].Token,
				ExtraArgs:    []byte{},
			})

			if isSourceFeeToken = err == nil; isSourceFeeToken {
				sourceFeeAmount = dione.EthBalanceToString(sourceFee)
			}
		}

		var destEnabled, isDestPool, isDestFeeToken bool
		destFeeAmount := "➖"

		if _, ok := destination.ChainConfig.SupportedTokens[token]; ok {
			destEnabled = slices.Contains(destTokens, destination.ChainConfig.SupportedTokens[token].Token)
			isDestPool = destination.ChainConfig.SupportedTokens[token].Pool != common.HexToAddress("")

			destFee, err := destRouter.GetFee(&bind.CallOpts{}, rhea.GetCCIPChainSelector(source.ChainConfig.EvmChainId), router.ClientEVM2AnyMessage{
				Receiver:     common.HexToAddress("").Bytes(),
				Data:         []byte{},
				TokenAmounts: []router.ClientEVMTokenAmount{},
				FeeToken:     destination.ChainConfig.SupportedTokens[token].Token,
				ExtraArgs:    []byte{},
			})

			if isDestFeeToken = err == nil; isDestFeeToken {
				destFeeAmount = dione.EthBalanceToString(destFee)
			}
		}

		boolParser := printBoolNeutral
		if sourceEnabled || destEnabled {
			boolParser = printBool
		}

		sb.WriteString(fmt.Sprintf("| %-20s | %9s | %13s | %9s | %9s | %9s | %13s | %9s | %9s |\n",
			token,
			printBoolNeutral(isSourceFeeToken),
			sourceFeeAmount,
			boolParser(sourceEnabled),
			boolParser(isSourcePool),
			printBoolNeutral(isDestFeeToken),
			destFeeAmount,
			boolParser(destEnabled),
			boolParser(isDestPool),
		))
	}

	sb.WriteString(generateSeparator(headerLengths))

	source.Logger.Info(sb.String())
}

func checkPriceRegistrySet(source *rhea.EvmDeploymentConfig, destination *rhea.EvmDeploymentConfig) {
	var sb strings.Builder

	tableHeaders := []string{"Token", "Remote ChainID", "ConfigSet"}
	headerLengths := []int{20, 14, 9}

	sb.WriteString(fmt.Sprintf("PriceRegistry token config for %s\n", ccip.ChainName(int64(source.ChainConfig.EvmChainId))))

	sb.WriteString(generateHeader(tableHeaders, headerLengths))

	feeManager, err := price_registry.NewPriceRegistry(source.ChainConfig.PriceRegistry, source.Client)
	helpers.PanicErr(err)

	for _, tokenName := range source.ChainConfig.FeeTokens {
		token := source.ChainConfig.SupportedTokens[tokenName].Token
		_, err = feeManager.GetTokenAndGasPrices(&bind.CallOpts{}, token, rhea.GetCCIPChainSelector(destination.ChainConfig.EvmChainId))
		if err != nil {
			sb.WriteString(fmt.Sprintf("| %-20s | %14d | %9s |\n", tokenName, destination.ChainConfig.EvmChainId, printBool(false)))
		}
		sb.WriteString(fmt.Sprintf("| %-20s | %14d | %9s |\n", tokenName, destination.ChainConfig.EvmChainId, printBool(true)))
	}
	sb.WriteString(generateSeparator(headerLengths))

	sb.WriteString(fmt.Sprintf("PriceRegistry token config for %s\n", ccip.ChainName(int64(destination.ChainConfig.EvmChainId))))
	sb.WriteString(generateHeader(tableHeaders, headerLengths))
	feeManager, err = price_registry.NewPriceRegistry(destination.ChainConfig.PriceRegistry, destination.Client)
	helpers.PanicErr(err)

	for _, tokenName := range destination.ChainConfig.FeeTokens {
		token := destination.ChainConfig.SupportedTokens[tokenName].Token
		_, err = feeManager.GetTokenAndGasPrices(&bind.CallOpts{}, token, rhea.GetCCIPChainSelector(source.ChainConfig.EvmChainId))
		if err != nil {
			sb.WriteString(fmt.Sprintf("| %-20s | %14d | %9s |\n", tokenName, source.ChainConfig.EvmChainId, printBool(false)))
		}
		sb.WriteString(fmt.Sprintf("| %-20s | %14d | %9s |\n", tokenName, source.ChainConfig.EvmChainId, printBool(true)))
	}

	sb.WriteString(generateSeparator(headerLengths))

	source.Logger.Info(sb.String())
}

func generateHeader(headers []string, headerLengths []int) string {
	var sb strings.Builder

	sb.WriteString(generateSeparator(headerLengths))
	sb.WriteString("|")
	for i, header := range headers {
		sb.WriteString(fmt.Sprintf(" %-"+strconv.Itoa(headerLengths[i])+"s |", header))
	}
	sb.WriteString("\n")
	sb.WriteString(generateSeparator(headerLengths))

	return sb.String()
}

func generateSeparator(headerLengths []int) string {
	length := 1

	for _, headerLength := range headerLengths {
		length += headerLength + 3
	}
	return strings.Repeat("─", length) + "\n"
}

// PrintJobSpecs prints the job spec for each node and CCIP spec type, as well as a bootstrap spec.
func PrintJobSpecs(env dione.Environment, sourceClient rhea.EvmDeploymentConfig, destClient rhea.EvmDeploymentConfig, version string) {
	don := dione.NewOfflineDON(env, nil)
	// jobparams for the lane
	jobParams := dione.NewCCIPJobSpecParams(&sourceClient.ChainConfig, sourceClient.LaneConfig, &destClient.ChainConfig, destClient.LaneConfig, version)

	bootstrapSpec := jobParams.BootstrapJob(destClient.LaneConfig.CommitStore.Hex())
	specString, err := bootstrapSpec.String()
	helpers.PanicErr(err)
	jobs := fmt.Sprintf("# BootstrapSpec%s", specString)

	commitJobSpec, err := jobParams.CommitJobSpec()
	helpers.PanicErr(err)
	committingChainID := commitJobSpec.OCR2OracleSpec.RelayConfig["chainID"].(uint64)
	executionSpec, err := jobParams.ExecutionJobSpec()
	helpers.PanicErr(err)
	execChainID := executionSpec.OCR2OracleSpec.RelayConfig["chainID"].(uint64)
	for i, oracle := range don.Config.Nodes {
		jobs += fmt.Sprintf("\n// [Node %d]\n", i)
		evmKeyBundle := dione.GetOCRkeysForChainType(oracle.OCRKeys, "evm")
		transmitterIDs := oracle.EthKeys

		// set node specific values
		commitJobSpec.OCR2OracleSpec.OCRKeyBundleID.SetValid(evmKeyBundle.ID)
		commitJobSpec.OCR2OracleSpec.TransmitterID.SetValid(transmitterIDs[fmt.Sprintf("%v", committingChainID)])
		specString, err := commitJobSpec.String()
		helpers.PanicErr(err)
		jobs += fmt.Sprintf("\n# CCIP commit spec%s", specString)

		// set node specific values
		executionSpec.OCR2OracleSpec.OCRKeyBundleID.SetValid(evmKeyBundle.ID)
		executionSpec.OCR2OracleSpec.TransmitterID.SetValid(transmitterIDs[fmt.Sprintf("%v", execChainID)])
		specString, err = executionSpec.String()
		helpers.PanicErr(err)
		jobs += fmt.Sprintf("\n# CCIP execution spec%s", specString)
	}
	fmt.Println(jobs)
}

func PrintBidirectionalTokenSupportState(source *rhea.EvmDeploymentConfig, destination *rhea.EvmDeploymentConfig) {
	printTokenSupportState(source, destination)
	printTokenSupportState(destination, source)
}

type TokenPool interface {
	IsOnRamp(opts *bind.CallOpts, onRamp common.Address) (bool, error)
	IsOffRamp(opts *bind.CallOpts, offRamp common.Address) (bool, error)
}

func printTokenSupportState(source *rhea.EvmDeploymentConfig, destination *rhea.EvmDeploymentConfig) {
	TOKEN := rhea.LINK
	sourceTokenConfig := source.ChainConfig.SupportedTokens[TOKEN]
	destTokenConfig := destination.ChainConfig.SupportedTokens[TOKEN]

	onRamp, err := evm_2_evm_onramp.NewEVM2EVMOnRamp(source.LaneConfig.OnRamp, source.Client)
	helpers.PanicErr(err)

	offRamp, err := evm_2_evm_offramp.NewEVM2EVMOffRamp(destination.LaneConfig.OffRamp, destination.Client)
	helpers.PanicErr(err)

	var sourcePool, destPool TokenPool

	if sourceTokenConfig.TokenPoolType == rhea.BurnMint {
		sourcePool, err = burn_mint_token_pool.NewBurnMintTokenPool(sourceTokenConfig.Pool, source.Client)
		helpers.PanicErr(err)
	} else {
		sourcePool, err = lock_release_token_pool.NewLockReleaseTokenPool(sourceTokenConfig.Pool, source.Client)
		helpers.PanicErr(err)
	}

	if destTokenConfig.TokenPoolType == rhea.BurnMint {
		destPool, err = burn_mint_token_pool.NewBurnMintTokenPool(destTokenConfig.Pool, destination.Client)
		helpers.PanicErr(err)
	} else {
		destPool, err = lock_release_token_pool.NewLockReleaseTokenPool(destTokenConfig.Pool, destination.Client)
		helpers.PanicErr(err)
	}

	tableHeaders := []string{"onRampAllowsToken", "onRampPoolCorrect", "offRampAllowsToken", "offRampPoolCorrect", "sourcePoolAllowsOnRamp", "destPoolAllowsOffRamp"}
	headerLengths := []int{20, 20, 20, 20, 20, 20}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("\nToken support for %s -> %s \n", ccip.ChainName(int64(source.ChainConfig.EvmChainId)), ccip.ChainName(int64(destination.ChainConfig.EvmChainId))))
	sb.WriteString(generateHeader(tableHeaders, headerLengths))

	// Check
	// onRamp token pool allowed
	// onRamp token pool is correct
	// offRamp token pool allowed
	// offRamp token pool is correct
	// source pool onRamp allowed
	// dest pool offRamp allowed
	var onRampAllowsPool, onRampPoolCorrect, offRampAllowsPool, offRampPoolCorrect, sourcePoolAllowsOnRamp, destPoolAllowsOffRamp bool

	poolBySourceToken, err := onRamp.GetPoolBySourceToken(&bind.CallOpts{}, sourceTokenConfig.Token)
	if err == nil && poolBySourceToken != common.HexToAddress("") {
		onRampAllowsPool = true
	}
	if poolBySourceToken == sourceTokenConfig.Pool {
		onRampPoolCorrect = true
	}

	destPoolBySourceToken, err := offRamp.GetPoolBySourceToken(&bind.CallOpts{}, sourceTokenConfig.Token)
	if err == nil && destPoolBySourceToken != common.HexToAddress("") {
		offRampAllowsPool = true
	}
	if destPoolBySourceToken == destTokenConfig.Pool {
		offRampPoolCorrect = true
	}

	sourcePoolAllowsOnRamp, err = sourcePool.IsOnRamp(&bind.CallOpts{}, onRamp.Address())
	helpers.PanicErr(err)
	destPoolAllowsOffRamp, err = destPool.IsOffRamp(&bind.CallOpts{}, offRamp.Address())
	helpers.PanicErr(err)

	sb.WriteString(fmt.Sprintf("| %20s | %20s | %20s | %20s | %20s | %20s |\n",
		printBool(onRampAllowsPool),
		printBool(onRampPoolCorrect),
		printBool(offRampAllowsPool),
		printBool(offRampPoolCorrect),
		printBool(sourcePoolAllowsOnRamp),
		printBool(destPoolAllowsOffRamp),
	))

	sb.WriteString(generateSeparator(headerLengths))

	source.Logger.Info(sb.String())
}
