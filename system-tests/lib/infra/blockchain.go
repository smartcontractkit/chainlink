package infra

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/rs/zerolog"
)

// TODO remove once https://smartcontract-it.atlassian.net/browse/CLO-1097 is done
func WaitForRPCEndpoint(lggr zerolog.Logger, url string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Try immediately first
	client, err := rpc.DialContext(ctx, url)
	if err == nil {
		defer client.Close()
		var blockNumber string
		if err := client.CallContext(ctx, &blockNumber, "eth_blockNumber"); err == nil {
			return nil
		}
	}

	// If immediate check fails, start periodic checks
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for RPC endpoint %s to be available", url)
		case <-ticker.C:
			lggr.Info().Msgf("waiting for %s to become available", url)
			client, err := rpc.DialContext(ctx, url)
			if err != nil {
				continue
			}

			var blockNumber string
			if err := client.CallContext(ctx, &blockNumber, "eth_blockNumber"); err != nil {
				continue
			}

			client.Close()
			// If we get here, the endpoint is responding
			return nil
		}
	}
}
