package verify

import (
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/contracts/permissionless_feeds_consumer"

	libformat "github.com/smartcontractkit/chainlink/system-tests/lib/format"
)

func ProofOfReserve(rpcURL, consumerContractAddress, feedID string, untilSuccessful bool, waitTime time.Duration) error {
	fmt.Print(libformat.DarkYellowText("🕵️‍♂️🔎 Verifying workflow execution...\n\n"))

	sethClient, sethErr := seth.NewClientBuilder().
		WithRpcUrl(rpcURL).
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

	fmt.Printf("Keystone Consumer contract address: %s\n", consumerContractAddress)
	fmt.Printf("Feed ID: %s\n", feedID)
	fmt.Printf("\nChecking if workflow has uploaded the value of TrueUSD asset\n")

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
					fmt.Printf("Value: %s\n", price.String())
					fmt.Printf("Timestamp: %d\n", timestamp)
					fmt.Print(libformat.DarkYellowText("\n✅ All good! Workflow executed successfully!\n"))

					return nil
				}
				fmt.Printf("🔍 Value not updated yet, retrying in %d seconds...\n", tickerSeconds)
			}
		case <-done:
			fmt.Print(libformat.DarkYellowText("\n❌ Workflow did not execute successfully within %s \n", waitTime.String()))
			return errors.New("workflow did not finish successfully")
		}
	}
}

func padRight(str string, length int, padChar rune) string {
	if len(str) >= length {
		return str
	}
	return str + strings.Repeat(string(padChar), length-len(str))
}
