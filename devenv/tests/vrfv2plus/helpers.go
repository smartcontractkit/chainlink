package vrfv2plus

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_coordinator_v2_5"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink/devenv/contracts"
	"github.com/smartcontractkit/chainlink/devenv/products"
	productvrfv2plus "github.com/smartcontractkit/chainlink/devenv/products/vrfv2plus"
)

const (
	defaultCallbackGasLimit = uint32(300000)
	defaultNumWords         = uint32(1)
	defaultRequestCount     = uint16(1)
	defaultFulfillTimeout   = 3 * time.Minute
)

// createAndFundSub creates a new subscription on the coordinator, funds it with LINK and native,
// and returns the subscription ID.
func createAndFundSub(
	ctx context.Context,
	chainClient *seth.Client,
	coord *contracts.EthereumVRFCoordinatorV2_5,
	linkToken contracts.LinkToken,
	fundLinkEth float64,
	fundNativeEth float64,
) (*big.Int, error) {
	subTx, err := coord.CreateSubscription()
	if err != nil {
		return nil, fmt.Errorf("CreateSubscription failed: %w", err)
	}
	receipt, err := chainClient.Client.TransactionReceipt(ctx, subTx.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get CreateSubscription receipt: %w", err)
	}
	subID, err := contracts.FindSubscriptionID(receipt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse subscription ID: %w", err)
	}

	// Fund with native
	nativeWei := products.EtherToWei(big.NewFloat(fundNativeEth))
	if err := coord.FundSubscriptionWithNative(subID, nativeWei); err != nil {
		return nil, fmt.Errorf("FundSubscriptionWithNative failed: %w", err)
	}

	// Fund with LINK via TransferAndCall
	linkJuels := products.EtherToWei(big.NewFloat(fundLinkEth))
	encodedSubID, err := encodeSubID(subID)
	if err != nil {
		return nil, fmt.Errorf("failed to encode subID: %w", err)
	}
	if _, err := linkToken.TransferAndCall(coord.Address(), linkJuels, encodedSubID); err != nil {
		return nil, fmt.Errorf("TransferAndCall (LINK funding) failed: %w", err)
	}

	return subID, nil
}

// requestAndWait sends a VRF randomness request from a consumer and waits for it to be fulfilled.
func requestAndWait(
	ctx context.Context,
	t *testing.T,
	consumer *contracts.EthereumVRFv2PlusLoadTestConsumer,
	coord *contracts.EthereumVRFCoordinatorV2_5,
	keyHash [32]byte,
	subID *big.Int,
	isNative bool,
	minConfs uint16,
	timeout time.Duration,
) *vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsFulfilled {
	t.Helper()

	requestID, err := consumer.RequestRandomness(
		keyHash,
		subID,
		minConfs,
		defaultCallbackGasLimit,
		isNative,
		defaultNumWords,
		defaultRequestCount,
	)
	require.NoError(t, err, "RequestRandomness failed")
	require.NotNil(t, requestID, "requestID should not be nil")

	// TODO try the original wait first? and only fallback to event later?
	// randomWordsFulfilledEvent, err := coordinator.WaitForRandomWordsFulfilledEvent(
	// 	contracts.RandomWordsFulfilledEventFilter{
	// 		SubIDs:     []*big.Int{subID},
	// 		RequestIds: []*big.Int{requestId},
	// 		Timeout:    randomWordsFulfilledEventTimeout,
	// 	},
	// )

	var fulfilled *vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsFulfilled
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		event, fErr := coord.FilterRandomWordsFulfilled(
			&bind.FilterOpts{Context: ctx},
			requestID,
		)
		if fErr != nil {
			return false
		}
		fulfilled = event
		return true
	}, timeout, 5*time.Second).Should(gomega.BeTrue(),
		"timed out waiting for RandomWordsFulfilled event for requestID %s", requestID)

	return fulfilled
}

// requestAndWaitWrapper sends a VRF request via the wrapper consumer and waits for fulfillment.
func requestAndWaitWrapper(
	ctx context.Context,
	t *testing.T,
	consumer *contracts.EthereumVRFV2PlusWrapperLoadTestConsumer,
	coord *contracts.EthereumVRFCoordinatorV2_5,
	isNative bool,
	minConfs uint16,
	timeout time.Duration,
) *vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsFulfilled {
	t.Helper()

	var requestID *big.Int
	var err error
	if isNative {
		requestID, err = consumer.RequestRandomWordsNative(minConfs, defaultCallbackGasLimit, defaultNumWords, defaultRequestCount)
	} else {
		requestID, err = consumer.RequestRandomWords(minConfs, defaultCallbackGasLimit, defaultNumWords, defaultRequestCount)
	}
	require.NoError(t, err, "RequestRandomWords failed")
	require.NotNil(t, requestID, "requestID should not be nil")

	var fulfilled *vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsFulfilled
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		event, fErr := coord.FilterRandomWordsFulfilled(
			&bind.FilterOpts{Context: ctx},
			requestID,
		)
		if fErr != nil {
			return false
		}
		fulfilled = event
		return true
	}, timeout, 5*time.Second).Should(gomega.BeTrue(),
		"timed out waiting for RandomWordsFulfilled event for requestID %s", requestID)

	return fulfilled
}

// encodeSubID ABI-encodes a uint256 subscription ID for use in TransferAndCall.
func encodeSubID(subID *big.Int) ([]byte, error) {
	b := make([]byte, 32)
	subIDBytes := subID.Bytes()
	if len(subIDBytes) > 32 {
		return nil, errors.New("subID too large for uint256")
	}
	copy(b[32-len(subIDBytes):], subIDBytes)
	return b, nil
}

// reconcileConfiguredFunding ensures wrapper subscription and wrapper consumer balances
// are at least the configured funding levels before each subtest.
func reconcileConfiguredFunding(
	ctx context.Context,
	t *testing.T,
	chainClient *seth.Client,
	coord *contracts.EthereumVRFCoordinatorV2_5,
	linkToken contracts.LinkToken,
	c *productvrfv2plus.VRFv2Plus,
) {
	t.Helper()

	targetLink := products.EtherToWei(big.NewFloat(c.SubFundingAmountLink))
	targetNative := products.EtherToWei(big.NewFloat(c.SubFundingAmountNative))

	// Ensure wrapper subscription balances.
	wrapperSubID, ok := new(big.Int).SetString(c.DeployedContracts.WrapperSubID, 10)
	require.True(t, ok, "failed to parse wrapper sub id: %s", c.DeployedContracts.WrapperSubID)

	sub, err := coord.GetSubscription(ctx, wrapperSubID)
	require.NoError(t, err, "failed to get wrapper subscription")

	if sub.NativeBalance.Cmp(targetNative) < 0 {
		deltaNative := new(big.Int).Sub(targetNative, sub.NativeBalance)
		err = coord.FundSubscriptionWithNative(wrapperSubID, deltaNative)
		require.NoError(t, err, "failed to top up wrapper sub native balance")
	}

	if sub.Balance.Cmp(targetLink) < 0 {
		deltaLink := new(big.Int).Sub(targetLink, sub.Balance)
		encodedSubID, eErr := encodeSubID(wrapperSubID)
		require.NoError(t, eErr, "failed to encode wrapper sub id")
		_, err = linkToken.TransferAndCall(coord.Address(), deltaLink, encodedSubID)
		require.NoError(t, err, "failed to top up wrapper sub LINK balance")
	}

	// Ensure wrapper consumer balances for direct-funding requests.
	wrapperConsumerAddr := c.DeployedContracts.WrapperConsumer
	currentConsumerNative, err := chainClient.Client.BalanceAt(ctx, common.HexToAddress(wrapperConsumerAddr), nil)
	require.NoError(t, err, "failed to get wrapper consumer native balance")

	if currentConsumerNative.Cmp(targetNative) < 0 {
		deltaNative := new(big.Int).Sub(targetNative, currentConsumerNative)
		privateKey, pErr := crypto.HexToECDSA(strings.TrimPrefix(products.NetworkPrivateKey(), "0x"))
		require.NoError(t, pErr, "failed to parse funding private key")
		feeCap, fErr := chainClient.Client.SuggestGasPrice(ctx)
		require.NoError(t, fErr, "failed to suggest gas fee cap")
		tipCap, tErr := chainClient.Client.SuggestGasTipCap(ctx)
		require.NoError(t, tErr, "failed to suggest gas tip cap")

		_, err = products.SendFunds(zerolog.Nop(), chainClient, products.FundsToSendPayload{
			ToAddress:  common.HexToAddress(wrapperConsumerAddr),
			Amount:     deltaNative,
			PrivateKey: privateKey,
			GasFeeCap:  feeCap,
			GasTipCap:  tipCap,
		})
		require.NoError(t, err, "failed to top up wrapper consumer native balance")
	}

	currentConsumerLink, err := linkToken.BalanceOf(ctx, wrapperConsumerAddr)
	require.NoError(t, err, "failed to get wrapper consumer LINK balance")
	if currentConsumerLink.Cmp(targetLink) < 0 {
		deltaLink := new(big.Int).Sub(targetLink, currentConsumerLink)
		err = linkToken.Transfer(wrapperConsumerAddr, deltaLink)
		require.NoError(t, err, "failed to top up wrapper consumer LINK balance")
	}

	// Final sanity checks so tests fail early with a clear reason.
	finalSub, err := coord.GetSubscription(ctx, wrapperSubID)
	require.NoError(t, err, "failed to re-read wrapper subscription")
	require.GreaterOrEqual(t, finalSub.Balance.Cmp(targetLink), 0, "wrapper sub LINK below configured funding target")
	require.GreaterOrEqual(t, finalSub.NativeBalance.Cmp(targetNative), 0, "wrapper sub native below configured funding target")

	finalConsumerNative, err := chainClient.Client.BalanceAt(ctx, common.HexToAddress(wrapperConsumerAddr), nil)
	require.NoError(t, err, "failed to re-read wrapper consumer native balance")
	require.GreaterOrEqual(t, finalConsumerNative.Cmp(targetNative), 0, "wrapper consumer native below configured funding target")

	finalConsumerLink, err := linkToken.BalanceOf(ctx, wrapperConsumerAddr)
	require.NoError(t, err, "failed to re-read wrapper consumer LINK balance")
	require.GreaterOrEqual(t, finalConsumerLink.Cmp(targetLink), 0, "wrapper consumer LINK below configured funding target")
}
