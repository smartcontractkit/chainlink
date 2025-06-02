package verify

import (
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/contracts/permissionless_feeds_consumer"
)

func ProofOfReserve(rpcUrl, consumerContractAddress, feedID string, untilSuccessful bool, waitTime time.Duration) error {
	sethClient, sethErr := seth.NewClientBuilder().
		WithRpcUrl(rpcUrl).
		WithReadOnlyMode().
		// do not check if there's a pending nonce nor check node's health
		WithProtections(false, false, seth.MustMakeDuration(time.Second)).
		Build()
	if sethErr != nil {
		return errors.Wrap(sethErr, "failed to connect to the Ethereum client")
	}

	address := common.HexToAddress(consumerContractAddress)

	contract, contractErr := permissionless_feeds_consumer.NewPermissionlessFeedsConsumer(address, sethClient.Client)
	if contractErr != nil {
		return errors.Wrap(contractErr, "failed to instantiate the Feeds Consumer contract")
	}

	feedID = strings.TrimPrefix(feedID, "0x")

	if len(feedID) != 64 {
		feedID = padRight(feedID, 64, '0')
	}

	fmt.Printf("Keysone Consumer contract address: %s\n", consumerContractAddress)
	fmt.Printf("Feed ID: %s\n", feedID)
	fmt.Printf("\nChecking if workflow has uplodad the value of TrueUSD asset\n")

	tickerSeconds := 10
	ticker := time.NewTicker(time.Duration(tickerSeconds) * time.Second)
	defer ticker.Stop()

	done := time.After(waitTime)

	for {
		select {
		case <-ticker.C:
			price, timestamp, priceErr := contract.GetPrice(sethClient.NewCallOpts(), common.HexToHash(feedID))
			if priceErr != nil {
				fmt.Printf("failed to read asset value: %s. Retrying in %d seconds...\n", priceErr, tickerSeconds)
			}

			if !untilSuccessful {
				return nil
			} else {
				if price.String() != "0" {
					fmt.Printf("\n✅ All good! Workflow executed successfully!\n")
					fmt.Printf("Value: %s\n", price.String())
					fmt.Printf("Timestamp: %d\n", timestamp)

					return nil
				}
				fmt.Printf("🔍 Value not updated yet, retrying in %d seconds...\n", tickerSeconds)
			}
		case <-done:
			fmt.Printf("\n❌ Workflow did not execute successfuly within %s \n", waitTime.String())
			return errors.New("workflow did not finish successfuly")
		}
	}
}

func padRight(str string, length int, padChar rune) string {
	if len(str) >= length {
		return str
	}
	return str + strings.Repeat(string(padChar), length-len(str))
}
