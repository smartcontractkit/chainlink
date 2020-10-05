package dhtrouter

import (
	"context"
	"fmt"
	"sync"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/types"

	"github.com/pkg/errors"
)

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

		wg.Add(1)
		go func(p peer.AddrInfo) {
			defer wg.Done()
			var err error
			err = ph.Connect(ctx, p)
			errs <- err
		}(p)
	}
	wg.Wait()

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
	protocolID := config.ProtocolID()

	const BucketSize = 64

	kadDHT, err := dht.New(ctx, aclHost,
		dht.BucketSize(BucketSize), dht.NamespacedValidator(ValidatorNamespace, AnnouncementValidator{}),
		dht.ProtocolPrefix(config.prefix), dht.ProtocolExtension(config.extension),
		dht.BootstrapPeers(config.bootstrapNodes...),
		dht.Mode(dht.ModeServer), dht.DisableProviders(),
		dht.QueryFilter(ACLQueryFilter(aclHost.GetACL(), protocolID, config.logger)),
		dht.RoutingTableFilter(ACLRoutingTableFilter(aclHost.GetACL(), protocolID, config.logger)),
		dht.Concurrency(2*config.failureThreshold+1), dht.Resiliency(config.failureThreshold+1))
	if err != nil {
		return nil, errors.Wrap(err, "could not create new DHT")
	}

	err = kadDHT.Bootstrap(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "error bootstrapping dht")
	}

	config.logger.Info("DHT initialized", types.LogFields{
		"id":             "DHT",
		"protocolID":     protocolID,
		"bootstrapNodes": config.bootstrapNodes,
	})
	return kadDHT, nil
}
