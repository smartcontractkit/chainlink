package products

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/onsi/gomega"
	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/devenv/contracts"
)

const (
	AnvilKey0                     = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	DefaultNativeTransferGasPrice = 21000
)

// WaitMinedFast is a method for Anvil's instant blocks mode to ovecrome bind.WaitMined ticker hardcode.
func WaitMinedFast(ctx context.Context, b bind.DeployBackend, txHash common.Hash) (*types.Receipt, error) {
	queryTicker := time.NewTicker(5 * time.Millisecond)
	defer queryTicker.Stop()
	for {
		receipt, err := b.TransactionReceipt(ctx, txHash)
		if err == nil {
			return receipt, nil
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-queryTicker.C:
		}
	}
}

func FundNewAddresses(ctx context.Context, keysRequired int, c *ethclient.Client, fundingAmountEth float64) ([]string, error) {
	pks := []string{}
	for range keysRequired {
		address, pk, err := seth.NewAddress()
		if err != nil {
			return nil, err
		}

		cErr := FundAddressEIP1559(ctx, c, NetworkPrivateKey(), address, fundingAmountEth)
		if cErr != nil {
			return nil, cErr
		}
		pks = append(pks, pk)
	}

	return pks, nil
}

// FundAddressEIP1559 funds an address using RPC URL, recipient address and amount of funds to send (ETH).
// Uses EIP-1559 transaction type.
func FundAddressEIP1559(ctx context.Context, c *ethclient.Client, pkey, recipientAddress string, amountOfFundsInETH float64) error {
	l := zerolog.Ctx(ctx)
	amount := new(big.Float).Mul(big.NewFloat(amountOfFundsInETH), big.NewFloat(1e18))
	amountWei, _ := amount.Int(nil)
	l.Info().Str("Addr", recipientAddress).Str("Wei", amountWei.String()).Msg("Funding Address")

	chainID, err := c.NetworkID(ctx)
	if err != nil {
		return err
	}
	privateKeyStr := strings.TrimPrefix(pkey, "0x")
	privateKey, err := crypto.HexToECDSA(privateKeyStr)
	if err != nil {
		return err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return errors.New("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := c.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return err
	}
	feeCap, err := c.SuggestGasPrice(ctx)
	if err != nil {
		return err
	}
	tipCap, err := c.SuggestGasTipCap(ctx)
	if err != nil {
		return err
	}
	recipient := common.HexToAddress(recipientAddress)
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		To:        &recipient,
		Value:     amountWei,
		Gas:       DefaultNativeTransferGasPrice,
		GasFeeCap: feeCap,
		GasTipCap: tipCap,
	})
	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainID), privateKey)
	if err != nil {
		return err
	}
	err = c.SendTransaction(ctx, signedTx)
	if err != nil {
		return err
	}
	if _, err := WaitMinedFast(ctx, c, signedTx.Hash()); err != nil {
		return err
	}
	l.Info().Str("Wei", amountWei.String()).Msg("Funded with ETH")
	return nil
}

// ETHClient creates a basic Ethereum client using PRIVATE_KEY env var and tip/cap gas settings
func ETHClient(ctx context.Context, rpcURL string, feeCapMult int64, tipCapMult int64) (*ethclient.Client, *bind.TransactOpts, string, error) {
	l := zerolog.Ctx(ctx)
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, nil, "", fmt.Errorf("could not connect to eth client: %w", err)
	}
	privateKey, err := crypto.HexToECDSA(NetworkPrivateKey())
	if err != nil {
		return nil, nil, "", fmt.Errorf("could not parse private key: %w", err)
	}
	publicKey := privateKey.PublicKey
	address := crypto.PubkeyToAddress(publicKey).String()
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, nil, "", fmt.Errorf("could not get chain ID: %w", err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, nil, "", fmt.Errorf("could not create transactor: %w", err)
	}
	fc, tc, err := multiplyEIP1559GasPrices(client, feeCapMult, tipCapMult)
	if err != nil {
		return nil, nil, "", fmt.Errorf("could not get bumped gas price: %w", err)
	}
	auth.GasFeeCap = fc
	auth.GasTipCap = tc
	l.Info().
		Str("GasFeeCap", fc.String()).
		Str("GasTipCap", tc.String()).
		Msg("Default gas prices set")
	return client, auth, address, nil
}

// multiplyEIP1559GasPrices returns bumped EIP1159 gas prices increased by multiplier
func multiplyEIP1559GasPrices(client *ethclient.Client, fcMult, tcMult int64) (*big.Int, *big.Int, error) { //nolint:revive // trivial function
	feeCap, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, nil, err
	}
	tipCap, err := client.SuggestGasTipCap(context.Background())
	if err != nil {
		return nil, nil, err
	}

	return new(big.Int).Mul(feeCap, big.NewInt(fcMult)), new(big.Int).Mul(tipCap, big.NewInt(tcMult)), nil
}

func NetworkPrivateKey() string {
	pk := os.Getenv("PRIVATE_KEY")
	if pk == "" {
		// that's the first Anvil and Geth private key, serves as a fallback for local testing if not overridden
		return AnvilKey0
	}
	return pk
}

type FundsToSendPayload struct {
	ToAddress  common.Address
	Amount     *big.Int
	PrivateKey *ecdsa.PrivateKey
	GasLimit   *int64
	GasPrice   *big.Int
	GasFeeCap  *big.Int
	GasTipCap  *big.Int
	TxTimeout  *time.Duration
}

// TODO: move to CTF?
// SendFunds sends native token amount (expressed in human-scale) from address controlled by private key
// to given address. You can override any or none of the following: gas limit, gas price, gas fee cap, gas tip cap.
// Values that are not set will be estimated or taken from config.
func SendFunds(logger zerolog.Logger, client *seth.Client, payload FundsToSendPayload) (*types.Receipt, error) {
	fromAddress, err := PrivateKeyToAddress(payload.PrivateKey)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), client.Cfg.Network.TxnTimeout.Duration())
	nonce, err := client.Client.PendingNonceAt(ctx, fromAddress)
	defer cancel()
	if err != nil {
		return nil, err
	}

	gasLimit, err := client.EstimateGasLimitForFundTransfer(fromAddress, payload.ToAddress, payload.Amount)
	if err != nil {
		transferGasFee := client.Cfg.Network.TransferGasFee
		if transferGasFee < 0 {
			return nil, fmt.Errorf("negative transfer gas fee: %d", transferGasFee)
		}
		gasLimit = uint64(transferGasFee)
	}

	gasPrice := big.NewInt(0)
	gasFeeCap := big.NewInt(0)
	gasTipCap := big.NewInt(0)

	if payload.GasLimit != nil {
		if *payload.GasLimit < 0 {
			return nil, fmt.Errorf("negative gas limit: %d", *payload.GasLimit)
		}
		gasLimit = uint64(*payload.GasLimit)
	}

	if client.Cfg.Network.EIP1559DynamicFees {
		// if any of the dynamic fees are not set, we need to either estimate them or read them from config
		if payload.GasFeeCap == nil || payload.GasTipCap == nil {
			// estimation or config reading happens here
			txOptions := client.NewTXOpts(seth.WithGasLimit(gasLimit))
			gasFeeCap = txOptions.GasFeeCap
			gasTipCap = txOptions.GasTipCap
		}

		// override with payload values if they are set
		if payload.GasFeeCap != nil {
			gasFeeCap = payload.GasFeeCap
		}

		if payload.GasTipCap != nil {
			gasTipCap = payload.GasTipCap
		}
	} else {
		if payload.GasPrice == nil {
			txOptions := client.NewTXOpts(seth.WithGasLimit(gasLimit))
			gasPrice = txOptions.GasPrice
		} else {
			gasPrice = payload.GasPrice
		}
	}

	var rawTx types.TxData

	if client.Cfg.Network.EIP1559DynamicFees {
		rawTx = &types.DynamicFeeTx{
			Nonce:     nonce,
			To:        &payload.ToAddress,
			Value:     payload.Amount,
			Gas:       gasLimit,
			GasFeeCap: gasFeeCap,
			GasTipCap: gasTipCap,
		}
	} else {
		rawTx = &types.LegacyTx{
			Nonce:    nonce,
			To:       &payload.ToAddress,
			Value:    payload.Amount,
			Gas:      gasLimit,
			GasPrice: gasPrice,
		}
	}

	signedTx, err := types.SignNewTx(payload.PrivateKey, types.LatestSignerForChainID(big.NewInt(client.ChainID)), rawTx)

	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to sign tx")
	}

	txTimeout := client.Cfg.Network.TxnTimeout.Duration()
	if payload.TxTimeout != nil {
		txTimeout = *payload.TxTimeout
	}

	logger.Debug().
		Str("From", fromAddress.Hex()).
		Str("To", payload.ToAddress.Hex()).
		Str("Amount (wei/ether)", fmt.Sprintf("%s/%s", payload.Amount, WeiToEther(payload.Amount).Text('f', -1))).
		Uint64("Nonce", nonce).
		Uint64("Gas Limit", gasLimit).
		Str("Gas Price", gasPrice.String()).
		Str("Gas Fee Cap", gasFeeCap.String()).
		Str("Gas Tip Cap", gasTipCap.String()).
		Bool("Dynamic fees", client.Cfg.Network.EIP1559DynamicFees).
		Msg("About to send funds")

	ctx, cancel = context.WithTimeout(context.Background(), txTimeout)
	defer cancel()
	err = client.Client.SendTransaction(ctx, signedTx)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to send transaction")
	}

	logger.Debug().
		Str("From", fromAddress.Hex()).
		Str("To", payload.ToAddress.Hex()).
		Str("TxHash", signedTx.Hash().String()).
		Str("Amount (wei/ether)", fmt.Sprintf("%s/%s", payload.Amount, WeiToEther(payload.Amount).Text('f', -1))).
		Uint64("Nonce", nonce).
		Uint64("Gas Limit", gasLimit).
		Str("Gas Price", gasPrice.String()).
		Str("Gas Fee Cap", gasFeeCap.String()).
		Str("Gas Tip Cap", gasTipCap.String()).
		Bool("Dynamic fees", client.Cfg.Network.EIP1559DynamicFees).
		Msg("Sent funds")

	receipt, receiptErr := client.WaitMined(ctx, logger, client.Client, signedTx)
	if receiptErr != nil {
		return nil, pkgerrors.Wrap(receiptErr, "failed to wait for transaction to be mined")
	}

	if receipt.Status == 1 {
		return receipt, nil
	}

	tx, _, err := client.Client.TransactionByHash(ctx, signedTx.Hash())
	if err != nil {
		return nil, pkgerrors.Wrap(err, "failed to get transaction by hash ")
	}

	_, err = client.Decode(tx, receiptErr)
	if err != nil {
		return nil, err
	}

	return receipt, nil
}

func PrivateKeyToAddress(privateKey *ecdsa.PrivateKey) (common.Address, error) {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}, errors.New("error casting public key to ECDSA")
	}
	return crypto.PubkeyToAddress(*publicKeyECDSA), nil
}

// EtherToWei converts an ETH float amount to wei
func EtherToWei(eth *big.Float) *big.Int {
	truncInt, _ := eth.Int(nil)
	truncInt = new(big.Int).Mul(truncInt, big.NewInt(params.Ether))
	fracStr := strings.Split(fmt.Sprintf("%.18f", eth), ".")[1]
	fracStr += strings.Repeat("0", 18-len(fracStr))
	fracInt, _ := new(big.Int).SetString(fracStr, 10)
	wei := new(big.Int).Add(truncInt, fracInt)
	return wei
}

// EtherToWei converts an ETH float amount to gwei
func EtherToGwei(eth *big.Float) *big.Int {
	truncInt, _ := eth.Int(nil)
	truncInt = new(big.Int).Mul(truncInt, big.NewInt(params.GWei))
	return truncInt
}

// WeiToEther converts a wei amount to eth float
func WeiToEther(wei *big.Int) *big.Float {
	f := new(big.Float)
	f.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	f.SetMode(big.ToNearestEven)
	fWei := new(big.Float)
	fWei.SetPrec(236) //  IEEE 754 octuple-precision binary floating-point format: binary256
	fWei.SetMode(big.ToNearestEven)
	return f.Quo(fWei.SetInt(wei), big.NewFloat(params.Ether))
}

func InitSeth(rpcURL string, privateKeys []string, chainID *uint64) (*seth.Client, error) {
	var chainClient *seth.Client
	var err error

	if os.Getenv(seth.CONFIG_FILE_ENV_VAR) != "" {
		sethCfg, sErr := seth.ReadConfig()
		if sErr != nil {
			return nil, sErr
		}

		if chainID == nil {
			return nil, errors.New("chainID of the network to use must be provided, when initialising Seth from TOML config file")
		}

		chainClient, err = seth.NewClientBuilderWithConfig(sethCfg).
			UseNetworkWithChainId(*chainID).
			WithPrivateKeys(privateKeys).
			WithRpcUrl(rpcURL).
			Build()
	} else {
		chainClient, err = seth.NewClientBuilder().
			WithPrivateKeys(privateKeys).
			WithRpcUrl(rpcURL).
			WithGasPriceEstimations(true, 0, seth.Priority_Auto, 1).
			WithProtections(true, false, seth.MustMakeDuration(1*time.Minute)).
			Build()
	}

	return chainClient, err
}

// WaitUntilChainHead blocks until the chain head is at least anchorBlock + minBlocksAfterAnchor
// On Anvil (chainID 1337) it deploys a Counter and spams Increment txs so blocks advance quickly;
// on other chains it polls block number until the target is reached.
func WaitUntilChainHead(
	ctx context.Context,
	t *testing.T,
	chainClient *seth.Client,
	anchorBlock uint64,
	minBlocksAfterAnchor int,
	chainID uint64,
	timeout time.Duration,
) {
	t.Helper()
	require.GreaterOrEqual(t, minBlocksAfterAnchor, 0, "minBlocksAfterAnchor must be non-negative")

	targetBlock := anchorBlock + uint64(minBlocksAfterAnchor) //nolint:gosec // minBlocksAfterAnchor validated non-negative above
	if chainID != 1337 {
		gomega.NewGomegaWithT(t).Eventually(func() bool {
			blk, err := chainClient.Client.BlockNumber(ctx)
			if err != nil {
				return false
			}
			return blk >= targetBlock
		}, timeout, time.Second).Should(gomega.BeTrue(),
			"timed out waiting for chain to reach block %d", targetBlock)
		return
	}

	counter, err := contracts.DeployCounterContract(chainClient)
	require.NoError(t, err, "failed to deploy counter contract for tx-spam block advancement")
	err = counter.Reset()
	require.NoError(t, err, "failed to reset counter contract for tx-spam")

	waitCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	done := make(chan struct{})
	var eg errgroup.Group
	eg.Go(func() error {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-waitCtx.Done():
				return fmt.Errorf("timeout waiting for chain to reach block %d", targetBlock)
			case <-ticker.C:
				blk, bErr := chainClient.Client.BlockNumber(waitCtx)
				if bErr != nil {
					continue
				}
				if blk >= targetBlock {
					close(done)
					return nil
				}
			}
		}
	})
	eg.Go(func() error {
		for {
			select {
			case <-done:
				return nil
			case <-waitCtx.Done():
				return fmt.Errorf("timeout while generating txs waiting for block %d", targetBlock)
			default:
				if iErr := counter.Increment(); iErr != nil {
					return iErr
				}
			}
		}
	})
	require.NoError(t, eg.Wait(), "failed while waiting for min chain head with tx-spam enabled on chainID=1337")
}
