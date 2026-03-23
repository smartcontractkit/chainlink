package vrfv2

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink/devenv/contracts"
	"github.com/smartcontractkit/chainlink/devenv/products"
	productvrfv2 "github.com/smartcontractkit/chainlink/devenv/products/vrfv2"
)

func deployConsumersAndFundSubs(
	ctx context.Context,
	chainClient *seth.Client,
	coord *contracts.EthereumVRFCoordinatorV2,
	link *contracts.EthereumLinkToken,
	subFundingLink float64,
	numConsumers, numSubs int,
) ([]*contracts.EthereumVRFv2LoadTestConsumer, []uint64, error) {
	consumers := make([]*contracts.EthereumVRFv2LoadTestConsumer, 0, numConsumers)
	for range numConsumers {
		c, err := contracts.DeployVRFv2LoadTestConsumer(chainClient, coord.Address())
		if err != nil {
			return nil, nil, err
		}
		consumers = append(consumers, c)
	}
	subIDs := make([]uint64, 0, numSubs)
	for range numSubs {
		receipt, err := coord.CreateSubscription()
		if err != nil {
			return nil, nil, err
		}
		subID, err := contracts.FindVRFv2SubscriptionID(receipt)
		if err != nil {
			return nil, nil, err
		}
		subIDs = append(subIDs, subID)
	}
	amount := products.EtherToWei(big.NewFloat(subFundingLink))
	for _, subID := range subIDs {
		enc, err := utils.ABIEncode(`[{"type":"uint64"}]`, subID)
		if err != nil {
			return nil, nil, err
		}
		if _, err := link.TransferAndCall(coord.Address(), amount, enc); err != nil {
			return nil, nil, err
		}
	}
	for _, subID := range subIDs {
		for _, c := range consumers {
			if err := coord.AddConsumer(subID, c.Address()); err != nil {
				return nil, nil, err
			}
		}
	}
	return consumers, subIDs, nil
}

func requestRandomnessAndWaitForFulfillment(
	ctx context.Context,
	consumer *contracts.EthereumVRFv2LoadTestConsumer,
	coord *contracts.EthereumVRFCoordinatorV2,
	keyHash [32]byte,
	subID uint64,
	minConf uint16,
	callbackGasLimit, numWords uint32,
	reqPerReq, reqDev uint16,
	fulfillTimeout time.Duration,
	keyNum int,
) (*contracts.CoordinatorRandomWordsRequested, *contracts.CoordinatorRandomWordsFulfilled, error) {
	req, err := consumer.RequestRandomnessFromKey(coord, keyHash, subID, minConf, callbackGasLimit, numWords, reqPerReq, keyNum)
	if err != nil {
		return nil, nil, err
	}
	fulfilled, err := contracts.WaitRandomWordsFulfilled(coord, req.RequestID, req.Raw.BlockNumber, fulfillTimeout)
	if err != nil {
		return req, nil, err
	}
	return req, fulfilled, nil
}

func directFundingRequestAndWait(
	ctx context.Context,
	consumer *contracts.EthereumVRFV2WrapperLoadTestConsumer,
	coord *contracts.EthereumVRFCoordinatorV2,
	subID uint64,
	minConf uint16,
	callbackGasLimit, numWords uint32,
	reqPerReq uint16,
	fulfillTimeout time.Duration,
) (*contracts.CoordinatorRandomWordsFulfilled, error) {
	req, err := consumer.RequestRandomness(coord, minConf, callbackGasLimit, numWords, reqPerReq)
	if err != nil {
		return nil, err
	}
	return contracts.WaitRandomWordsFulfilled(coord, req.RequestID, req.Raw.BlockNumber, fulfillTimeout)
}

func waitForRequestCountEqualToFulfillmentCount(
	ctx context.Context,
	consumer *contracts.EthereumVRFv2LoadTestConsumer,
	timeout time.Duration,
	wg *sync.WaitGroup,
) (reqCount *big.Int, fulCount *big.Int, err error) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	deadline := time.Now().Add(timeout)
	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return reqCount, fulCount, ctx.Err()
		case <-ticker.C:
			m, mErr := consumer.GetLoadTestMetrics(ctx)
			if mErr != nil {
				wg.Done()
				return nil, nil, mErr
			}
			reqCount, fulCount = m.RequestCount, m.FulfilmentCount
			if m.RequestCount.Cmp(m.FulfilmentCount) == 0 && m.RequestCount.Sign() > 0 {
				wg.Done()
				return m.RequestCount, m.FulfilmentCount, nil
			}
			if time.Now().After(deadline) {
				wg.Done()
				return reqCount, fulCount, fmt.Errorf("timeout waiting request==fulfillment counts (req=%s ful=%s)",
					reqCount.String(), fulCount.String())
			}
		}
	}
}

func deleteAllJobs(node *clclient.ChainlinkClient) error {
	jobs, _, err := node.ReadJobs()
	if err != nil {
		return err
	}
	for _, m := range jobs.Data {
		id, ok := m["id"].(string)
		if !ok {
			return fmt.Errorf("job missing id: %+v", m)
		}
		resp, err := node.DeleteJob(id)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
			return fmt.Errorf("delete job %s: status %d", id, resp.StatusCode)
		}
	}
	return nil
}

func getTxFromAddress(tx *types.Transaction) (string, error) {
	from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
	if err != nil {
		return "", err
	}
	return from.Hex(), nil
}

func parseFulfillTimeout(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil || d <= 0 {
		return 2 * time.Minute
	}
	return d
}

func mustKeyHash(c *productvrfv2.VRFv2) [32]byte {
	h := common.HexToHash(c.VRFKeyData.KeyHash)
	var out [32]byte
	copy(out[:], h[:])
	return out
}
