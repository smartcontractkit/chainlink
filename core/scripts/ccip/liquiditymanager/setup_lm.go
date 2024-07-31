package main

import (
	"context"
	"flag"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/urfave/cli"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/scripts/ccip/liquiditymanager/arb"
	"github.com/smartcontractkit/chainlink/core/scripts/ccip/liquiditymanager/multienv"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/cmd"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/weth9"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/erc20"
)

var tomlConfigTemplate = `
# Arbitrum Sepolia
[[EVM]]
ChainID = "421614"
FinalityDepth = 1
LogPollInterval = "1s"
GasEstimator.LimitDefault = 3_500_000

[[EVM.Nodes]]
HTTPURL = "%s"
Name = "arbitrum_sepolia_1"
WSURL = "%s"

# Sepolia
[[EVM]]
ChainID = "11155111"
FinalityDepth = 2
LogPollInterval = "12s"
GasEstimator.LimitDefault = 3_500_000

[[EVM.Nodes]]
HTTPURL = "%s"
Name = "sepolia_1"
WSURL = "%s"

[EVM.Transactions]
ForwardersEnabled = false

[Feature]
LogPoller = true

[OCR2]
Enabled = true
ContractPollInterval = "15s"

[OCR]
Enabled = false

[P2P.V2]
ListenAddresses = ["127.0.0.1:8000"]
`

func setupLiquidityManagerNodes(e multienv.Env) {
	fs := flag.NewFlagSet("setup-liquiditymanager-nodes", flag.ExitOnError)
	l1ChainID := fs.Uint64("l1-chain-id", chainsel.ETHEREUM_TESTNET_SEPOLIA.EvmChainID, "L1 chain ID")
	l2ChainID := fs.Uint64("l2-chain-id", chainsel.ETHEREUM_TESTNET_SEPOLIA_ARBITRUM_1.EvmChainID, "L2 chain ID")
	l1TokenAddress := fs.String("l1-token-address",
		arb.ArbitrumContracts[chainsel.ETHEREUM_TESTNET_SEPOLIA.EvmChainID]["WETH"].Hex(),
		"L1 token address")
	l2TokenAddress := fs.String("l2-token-address",
		arb.ArbitrumContracts[chainsel.ETHEREUM_TESTNET_SEPOLIA_ARBITRUM_1.EvmChainID]["WETH"].Hex(),
		"L2 token address")
	apiFile := fs.String("api",
		"../../../../tools/secrets/apicredentials", "api credentials file")
	passwordFile := fs.String("password",
		"../../../../tools/secrets/password.txt", "password file")
	databasePrefix := fs.String("database-prefix",
		"postgres://postgres:postgres_password_padded_for_security@localhost:5432/liquiditymanager-test", "database prefix")
	databaseSuffixes := fs.String("database-suffixes",
		"sslmode=disable", "database parameters to be added")
	nodeCount := fs.Int("node-count", 5, "number of nodes")
	fundingAmount := fs.String("funding-amount", "100000000000000000", "amount to fund nodes") // .1 ETH
	resetDatabase := fs.Bool("reset-database", true, "boolean to reset database")

	helpers.ParseArgs(fs, os.Args[1:])

	validateEnv(e, *l1ChainID, *l2ChainID, true)

	transmitterFunding := decimal.RequireFromString(*fundingAmount).BigInt()

	uni := deployUniverse(e,
		*l1ChainID,
		*l2ChainID,
		common.HexToAddress(*l1TokenAddress),
		common.HexToAddress(*l2TokenAddress))

	fmt.Println("Configuring nodes with liquidityManager jobs...")
	var (
		onChainPublicKeys  []string
		offChainPublicKeys []string
		configPublicKeys   []string
		peerIDs            []string
		transmitters       = make(map[string][]string)
	)
	for i := 0; i < *nodeCount; i++ {
		flagSet := flag.NewFlagSet("run-liquidityManager-job-creation", flag.ExitOnError)
		flagSet.String("api", *apiFile, "api file")
		flagSet.String("password", *passwordFile, "password file")
		flagSet.String("vrfpassword", *passwordFile, "vrf password file")
		flagSet.String("bootstrapPort", fmt.Sprintf("%d", 8000), "port of bootstrap")
		flagSet.Int64("l1ChainID", int64(*l1ChainID), "the L1 chain ID")
		flagSet.Int64("l2ChainID", int64(*l2ChainID), "the L2 chain ID")
		flagSet.Bool("applyInitServerConfig", true, "override for using initServerConfig in App.Before")

		flagSet.String("job-type", "liquiditymanager", "the job type")
		flagSet.String("job-name", fmt.Sprintf("liquiditymanager-%d", i+1), "the job name")

		flagSet.String("liquidityManagerAddress", uni.L1.LiquidityManager.Hex(), "the liquidity manager address")
		flagSet.Uint64("liquidityManagerNetwork", mustGetChainByEvmID(*l1ChainID).Selector, "the liquidity manager network")

		// used by bootstrap template instantiation
		flagSet.String("contractID", uni.L1.LiquidityManager.Hex(), "the contract to get peers from")

		flagSet.Bool("dangerWillRobinson", *resetDatabase, "for resetting databases")
		flagSet.Bool("isBootstrapper", i == 0, "is first node")
		bootstrapperPeerID := ""
		if len(peerIDs) != 0 {
			bootstrapperPeerID = peerIDs[0]
		}
		flagSet.String("bootstrapperPeerID", bootstrapperPeerID, "peerID of first node")

		payload := SetupNode(e, *l1ChainID, *l2ChainID, flagSet, i, *databasePrefix, *databaseSuffixes, *resetDatabase)

		onChainPublicKeys = append(onChainPublicKeys, payload.OnChainPublicKey)
		offChainPublicKeys = append(offChainPublicKeys, payload.OffChainPublicKey)
		configPublicKeys = append(configPublicKeys, payload.ConfigPublicKey)
		peerIDs = append(peerIDs, payload.PeerID)
		for chainIDStr, transmitter := range payload.Transmitters {
			transmitters[chainIDStr] = append(transmitters[chainIDStr], transmitter)
		}
	}

	printStandardCommands(uni,
		fmt.Sprintf("%d", *l1ChainID),
		fmt.Sprintf("%d", *l2ChainID),
		*l1TokenAddress,
		*l2TokenAddress,
		onChainPublicKeys,
		offChainPublicKeys,
		configPublicKeys,
		peerIDs,
		transmitters)

	fmt.Println("Funding transmitters on L1...")
	FundNodes(e, *l1ChainID, transmitters[fmt.Sprintf("%d", *l1ChainID)], transmitterFunding)
	fmt.Println()

	fmt.Println("Funding transmitters on L2...")
	FundNodes(e, *l2ChainID, transmitters[fmt.Sprintf("%d", *l2ChainID)], transmitterFunding)
	fmt.Println()
}

func fundPoolAndLiquidityManager(
	e multienv.Env,
	chainID uint64,
	tokenAddress,
	tokenPoolAddress,
	liquidityManagerAddress common.Address,
	tokenPoolFunding *big.Int,
	liquidityManagerFunding *big.Int) {
	token, err := erc20.NewERC20(tokenAddress, e.Clients[chainID])
	helpers.PanicErr(err)

	// check if we have enough balance to transfer
	// try to deposit if token is WETH
	balance, err := token.BalanceOf(nil, e.Transactors[chainID].From)
	helpers.PanicErr(err)
	if balance.Cmp(tokenPoolFunding) < 0 {
		symbol, err2 := token.Symbol(nil)
		helpers.PanicErr(err2)
		if symbol == "WETH" {
			l1Weth, err3 := weth9.NewWETH9(tokenAddress, e.Clients[chainID])
			helpers.PanicErr(err3)

			nativeBalance, err3 := e.Clients[chainID].BalanceAt(
				context.Background(),
				e.Transactors[chainID].From,
				nil)
			helpers.PanicErr(err3)
			if nativeBalance.Cmp(tokenPoolFunding) < 0 {
				helpers.PanicErr(fmt.Errorf("not enough balance to deposit WETH"))
			}

			fmt.Println("Depositing", tokenPoolFunding.String(), "to WETH...")
			tx, err3 := l1Weth.Deposit(&bind.TransactOpts{
				From:   e.Transactors[chainID].From,
				Signer: e.Transactors[chainID].Signer,
				Value:  tokenPoolFunding,
			})
			helpers.PanicErr(err3)
			helpers.ConfirmTXMined(
				context.Background(),
				e.Clients[chainID],
				tx,
				int64(chainID),
				"Depositing", tokenPoolFunding.String(), "to WETH token at", tokenAddress.Hex())
		} else {
			helpers.PanicErr(
				fmt.Errorf("not enough balance to fund token pool, please get more tokens (address: %s)",
					tokenAddress.Hex()))
		}
	}

	fmt.Println("Funding token pool on", chainID, "with", tokenPoolFunding, "...")
	tx, err := token.Transfer(e.Transactors[chainID], tokenPoolAddress, tokenPoolFunding)
	helpers.PanicErr(err)
	helpers.ConfirmTXMined(
		context.Background(),
		e.Clients[chainID],
		tx,
		int64(chainID),
		"Transferring", tokenPoolFunding.String(), "to token pool at", tokenPoolAddress.Hex())

	fmt.Println("Funding liquidityManager on", chainID, "with", liquidityManagerFunding, "wei...")
	if err := FundNode(e, chainID, liquidityManagerAddress, liquidityManagerFunding); err != nil {
		fmt.Println("Failed to fund liquidityManager on", chainID, "with", liquidityManagerFunding, "wei:", err)
	}
}

func printStandardCommands(
	uni universe,
	l1ChainID,
	l2ChainID string,
	l1TokenAddress, l2TokenAddress string,
	onChainPublicKeys,
	offchainPublicKeys,
	configPublicKeys,
	peerIDs []string,
	transmitters map[string][]string,
) {
	fmt.Println("Contract Deployments complete\n",
		"L1 Arm:", uni.L1.Arm.Hex(), "\n",
		"L1 Arm Proxy:", uni.L1.ArmProxy.Hex(), "\n",
		"L1 Token Pool:", uni.L1.TokenPool.Hex(), "\n",
		"L1 LiquidityManager:", uni.L1.LiquidityManager.Hex(), "\n",
		"L1 Bridge Adapter:", uni.L1.BridgeAdapterAddress.Hex(), "\n",
		"L2 Arm:", uni.L2.Arm.Hex(), "\n",
		"L2 Arm Proxy:", uni.L2.ArmProxy.Hex(), "\n",
		"L2 Token Pool:", uni.L2.TokenPool.Hex(), "\n",
		"L2 LiquidityManager:", uni.L2.LiquidityManager.Hex(), "\n",
		"L2 Bridge Adapter:", uni.L2.BridgeAdapterAddress.Hex(), "\n",
		"Node launches complete\n",
		"OnChainPublicKeys:", strings.Join(onChainPublicKeys, ","), "\n",
		"OffChainPublicKeys:", strings.Join(offchainPublicKeys, ","), "\n",
		"ConfigPublicKeys:", strings.Join(configPublicKeys, ","), "\n",
		"PeerIDs:", strings.Join(peerIDs, ","), "\n",
		"Transmitters L1:", strings.Join(transmitters[l1ChainID], ","), "\n",
		"Transmitters L2:", strings.Join(transmitters[l2ChainID], ","),
	)
	fmt.Println()
	fmt.Println("Set config command:", "\n",
		"go run . set-config -l1-chain-id", l1ChainID,
		"-l2-chain-id", l2ChainID,
		"-l1-liquiditymanager-address", uni.L1.LiquidityManager.Hex(),
		"-l2-liquiditymanager-address", uni.L2.LiquidityManager.Hex(),
		"-signers", strings.Join(onChainPublicKeys, ","),
		"-offchain-pubkeys", strings.Join(offchainPublicKeys, ","),
		"-config-pubkeys", strings.Join(configPublicKeys, ","),
		"-peer-ids", strings.Join(peerIDs, ","),
		"-l1-transmitters", strings.Join(transmitters[l1ChainID], ","),
		"-l2-transmitters", strings.Join(transmitters[l2ChainID], ","),
	)
	fmt.Println()
	fmt.Println("Funding command:", "\n",
		"go run . fund-contracts -l1-chain-id", l1ChainID,
		"-l2-chain-id", l2ChainID,
		"-l1-liquiditymanager-address", uni.L1.LiquidityManager.Hex(),
		"-l2-liquiditymanager-address", uni.L2.LiquidityManager.Hex(),
		"-l1-token-address", l1TokenAddress,
		"-l2-token-address", l2TokenAddress,
		"-l1-token-pool-address", uni.L1.TokenPool.Hex(),
		"-l2-token-pool-address", uni.L2.TokenPool.Hex(),
	)
}

func SetupNode(
	e multienv.Env,
	l1ChainID, l2ChainID uint64,
	flagSet *flag.FlagSet,
	nodeIdx int,
	databasePrefix,
	databaseSuffixes string,
	resetDB bool,
) *cmd.SetupLiquidityManagerNodePayload {
	configureEnvironmentVariables(e, l1ChainID, l2ChainID, nodeIdx, databasePrefix, databaseSuffixes)

	client := newSetupClient()
	app := cmd.NewApp(client)
	ctx := cli.NewContext(app, flagSet, nil)

	defer func() {
		err := app.After(ctx)
		helpers.PanicErr(err)
	}()

	err := app.Before(ctx)
	helpers.PanicErr(err)

	if resetDB {
		resetDatabase(client, ctx)
	}

	return setupLiquidityManagerNodeFromClient(client, ctx)
}

func configureEnvironmentVariables(
	e multienv.Env,
	l1ChainID, l2ChainID uint64,
	index int,
	databasePrefix string,
	databaseSuffixes string,
) {
	// Set permitted envars for v2.
	helpers.PanicErr(os.Setenv("CL_DATABASE_URL", fmt.Sprintf("%s-%d?%s", databasePrefix, index, databaseSuffixes)))
	helpers.PanicErr(os.Setenv("CL_CONFIG", fmt.Sprintf(
		tomlConfigTemplate,
		e.HTTPURLs[l2ChainID],
		e.WSURLs[l2ChainID],
		e.HTTPURLs[l1ChainID],
		e.WSURLs[l1ChainID],
	)))

	// Unset prohibited envars for v2.
	helpers.PanicErr(os.Unsetenv("ETH_URL"))
	helpers.PanicErr(os.Unsetenv("ETH_HTTP_URL"))
	helpers.PanicErr(os.Unsetenv("ETH_CHAIN_ID"))
}

func newSetupClient() *cmd.Shell {
	prompter := cmd.NewTerminalPrompter()
	return &cmd.Shell{
		Renderer:                       cmd.RendererTable{Writer: os.Stdout},
		AppFactory:                     cmd.ChainlinkAppFactory{},
		KeyStoreAuthenticator:          cmd.TerminalKeyStoreAuthenticator{Prompter: prompter},
		FallbackAPIInitializer:         cmd.NewPromptingAPIInitializer(prompter),
		Runner:                         cmd.ChainlinkRunner{},
		PromptingSessionRequestBuilder: cmd.NewPromptingSessionRequestBuilder(prompter),
		ChangePasswordPrompter:         cmd.NewChangePasswordPrompter(),
		PasswordPrompter:               cmd.NewPasswordPrompter(),
	}
}

func resetDatabase(client *cmd.Shell, context *cli.Context) {
	helpers.PanicErr(client.ResetDatabase(context))
}

func setupLiquidityManagerNodeFromClient(
	client *cmd.Shell,
	context *cli.Context) *cmd.SetupLiquidityManagerNodePayload {
	payload, err := client.ConfigureRebalancerNode(context)
	helpers.PanicErr(err)

	return payload
}

func FundNodes(e multienv.Env, chainID uint64, transmitters []string, fundingAmount *big.Int) {
	var errs error
	for _, transmitter := range transmitters {
		errs = multierr.Append(errs, FundNode(e, chainID, common.HexToAddress(transmitter), fundingAmount))
	}
	if errs != nil {
		fmt.Println("Encountered errors funding nodes: ", errs)
	}
}

func FundNode(
	e multienv.Env,
	chainID uint64,
	toAddress common.Address,
	fundingAmount *big.Int,
) error {
	client, transactor := e.Clients[chainID], e.Transactors[chainID]

	nonce, err := client.PendingNonceAt(context.Background(), transactor.From)
	if err != nil {
		return err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}

	gasEstimate, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		From:  transactor.From,
		To:    &toAddress,
		Value: fundingAmount,
	})
	if err != nil {
		return err
	}

	tx := types.NewTx(
		&types.LegacyTx{
			Nonce:    nonce,
			GasPrice: gasPrice,
			Gas:      gasEstimate,
			To:       &toAddress,
			Value:    fundingAmount,
		},
	)
	signedTx, err := transactor.Signer(transactor.From, tx)
	if err != nil {
		return err
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return err
	}
	helpers.ConfirmTXMined(context.Background(), client, signedTx, int64(chainID), "Sending", fundingAmount.String(), "to", toAddress.Hex())
	return nil
}

func mustGetChainByEvmID(evmChainID uint64) chainsel.Chain {
	ch, exists := chainsel.ChainByEvmChainID(evmChainID)
	if !exists {
		helpers.PanicErr(fmt.Errorf("chain id %d doesn't exist in chain-selectors - forgot to add?", evmChainID))
	}
	return ch
}
