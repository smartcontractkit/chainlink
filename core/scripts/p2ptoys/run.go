package main

import (
	"context"
	"crypto/ed25519"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/smartcontractkit/chainlink/core/scripts/p2ptoys/common"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/p2p"

	"github.com/smartcontractkit/libocr/ragep2p"
	ragetypes "github.com/smartcontractkit/libocr/ragep2p/types"

	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

/*
Usage:

	go run run.go --bootstrap
	go run run.go --index 0
	go run run.go --index 1

Observe nodes 0 and 1 discovering each other via the bootstrapper and exchanging messages.
*/
func main() {
	lggr, _ := logger.NewLogger()
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	var shutdownWaitGroup sync.WaitGroup

	configFile := flag.String("config", "test_keys.json", "Path to JSON config file")
	nodeIndex := flag.Int("index", 0, "Index of the key in the config file to use")
	isBootstrap := flag.Bool("bootstrap", false, "Whether to run as a bootstrapper or not")
	flag.Parse()
	config, err := common.ParseConfigFromFile(*configFile)
	if err != nil {
		lggr.Error("error parsing config ", err)
		return
	}

	// Select this node's private key and listen address from config.
	var privateKey ed25519.PrivateKey
	var peerIDs []ragetypes.PeerID
	var listenAddr string
	if *isBootstrap {
		privateKey = config.BootstrapperKeys[*nodeIndex]
		listenAddr = fmt.Sprintf("127.0.0.1:%d", common.BootstrapStartPort+*nodeIndex)
	} else {
		privateKey = config.NodeKeys[*nodeIndex]
		listenAddr = fmt.Sprintf("127.0.0.1:%d", common.NodeStartPort+*nodeIndex)
	}
	for _, key := range config.NodeKeys {
		peerID, _ := ragetypes.PeerIDFromPrivateKey(key)
		peerIDs = append(peerIDs, peerID)
	}

	reg := prometheus.NewRegistry()
	peerConfig := p2p.PeerConfig{
		PrivateKey:      privateKey,
		ListenAddresses: []string{listenAddr},
		Bootstrappers:   config.BootstrapperPeerInfos,

		DeltaReconcile:     time.Second * 5,
		DeltaDial:          time.Second * 5,
		DiscovererDatabase: p2p.NewInMemoryDiscovererDatabase(),
		MetricsRegisterer:  reg,
	}

	peer, err := p2p.NewPeer(peerConfig, lggr)
	if err != nil {
		lggr.Error("error creating peer:", err)
		return
	}
	err = peer.Start(ctx)
	if err != nil {
		lggr.Error("error starting peer:", err)
		return
	}

	peers := make(map[ragetypes.PeerID]p2ptypes.StreamConfig)
	for _, peerID := range peerIDs {
		peers[peerID] = p2ptypes.StreamConfig{
			IncomingMessageBufferSize: 1000000,
			OutgoingMessageBufferSize: 1000000,
			MaxMessageLenBytes:        100000,
			MessageRateLimiter: ragep2p.TokenBucketParams{
				Rate:     2.0,
				Capacity: 10,
			},
			BytesRateLimiter: ragep2p.TokenBucketParams{
				Rate:     100000.0,
				Capacity: 100000,
			},
		}
	}

	err = peer.UpdateConnections(peers)
	if err != nil {
		lggr.Errorw("error updating peer addresses", "error", err)
	}

	if !*isBootstrap {
		shutdownWaitGroup.Add(2)
		go sendLoop(ctx, &shutdownWaitGroup, peer, peerIDs, *nodeIndex, lggr)
		go recvLoop(ctx, &shutdownWaitGroup, peer.Receive(), lggr)
	}

	<-ctx.Done()
	err = peer.Close()
	if err != nil {
		lggr.Error("error closing peer:", err)
	}
	shutdownWaitGroup.Wait()
}

func sendLoop(ctx context.Context, shutdownWaitGroup *sync.WaitGroup, peer p2ptypes.Peer, peerIDs []ragetypes.PeerID, myId int, lggr logger.Logger) {
	defer shutdownWaitGroup.Done()
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	lastId := 0
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if lastId != myId {
				lggr.Infow("sending message", "receiver", peerIDs[lastId])
				err := peer.Send(peerIDs[lastId], []byte("hello!"))
				if err != nil {
					lggr.Errorw("error sending message", "receiver", peerIDs[lastId], "error", err)
				}
			}
			lastId++
			if lastId >= len(peerIDs) {
				lastId = 0
			}
		}
	}
}

func recvLoop(ctx context.Context, shutdownWaitGroup *sync.WaitGroup, chRecv <-chan p2ptypes.Message, lggr logger.Logger) {
	defer shutdownWaitGroup.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-chRecv:
			lggr.Info("received message from ", msg.Sender, " : ", string(msg.Payload))
		}
	}
}
