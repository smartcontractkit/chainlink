package dhtrouter

import (
	"context"
	"fmt"
	"sync"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting/types"

	"github.com/pkg/errors"
)

// This code is borrowed from the go-ipfs bootstrap process

// Returns error only if all connections failed
func tryConnectToBootstrappers(ctx context.Context, ph host.Host, peers []peer.AddrInfo) error {
	if len(peers) == 0 {
		return nil
	}

	errs := make(chan error, len(peers))
	var wg sync.WaitGroup
	for _, p := range peers {
		if ph.ID() == p.ID {
			errs <- errors.New("will not connect to self")
			continue
		}

		// performed asynchronously because when performed synchronously, if
		// one `Connect` call hangs, subsequent calls are more likely to
		// fail/abort due to an expiring context.
		// Also, performed asynchronously for dial speed.

		wg.Add(1)
		go func(p peer.AddrInfo) {
			defer wg.Done()
			err := ph.Connect(ctx, p)
			errs <- err
		}(p)
	}
	wg.Wait()

	// our failure condition is when no connection attempt succeeded.
	// So drain the errs channel, counting the results.
	close(errs)
	count := 0
	var err error
	for err = range errs {
		if err != nil {
			count++
		}
	}
	if count == len(peers) {
		return fmt.Errorf("failed to bootstrap. Last error: %w", err)
	}
	return nil
}

func newDHT(ctx context.Context, config DHTNodeConfig, aclHost ACLHost) (*dht.IpfsDHT, error) {
	// create a kadDHT
	protocolID := config.ProtocolID()

	// we set the bucket size large enough so that all peers are in the same bucket.
	const BucketSize = 2 * types.MaxOracles

	kadDHT, err := dht.New(ctx, aclHost,
		dht.BucketSize(BucketSize), // K in the Kademlia paper
		dht.NamespacedValidator(ValidatorNamespace, AnnouncementValidator{}), // THIS IS CRITICAL. WE MUST USE AnnouncementValidator.
		dht.ProtocolPrefix(config.prefix),                                    // stands for off-chain reporting
		dht.ProtocolExtension(config.extension),
		dht.BootstrapPeers(config.bootstrapNodes...),
		dht.Mode(dht.ModeServer), // this must be set in order for DHT to work in internal networks
		dht.DisableProviders(),
		dht.QueryFilter(ACLQueryFilter(aclHost.GetACL(), protocolID, config.logger)),
		dht.RoutingTableFilter(ACLRoutingTableFilter(aclHost.GetACL(), protocolID, config.logger)),
		dht.Concurrency(2*config.failureThreshold+1), // query 2f+1, so that at least f+1 will eventually respond
		dht.Resiliency(config.failureThreshold+1),    // wait for f+1 response so that at least one of them is from an honest player
	)
	if err != nil {
		return nil, errors.Wrap(err, "could not create new DHT")
	}

	// this simply starts the routing table manager thread
	err = kadDHT.Bootstrap(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error bootstrapping dht")
	}

	config.logger.Info("DHT initialized", commontypes.LogFields{
		"id":             "DHT",
		"protocolID":     protocolID,
		"bootstrapNodes": config.bootstrapNodes,
	})
	return kadDHT, nil
}
