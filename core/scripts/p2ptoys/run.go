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

	"github.com/smartcontractkit/chainlink/core/scripts/p2ptoys/common"
	"github.com/smartcontractkit/chainlink/v2/core/logger"

	commonlogger "github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/networking"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
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
	var listenAddr string
	if *isBootstrap {
		privateKey = config.BootstrapperKeys[*nodeIndex]
		listenAddr = fmt.Sprintf("127.0.0.1:%d", common.BootstrapStartPort+*nodeIndex)
	} else {
		privateKey = config.NodeKeys[*nodeIndex]
		listenAddr = fmt.Sprintf("127.0.0.1:%d", common.NodeStartPort+*nodeIndex)
	}

	// Create a new concretePeerV2 object, which wraps RageP2P's Host and a Discoverer.
	peerConfig := networking.PeerConfig{
		PrivKey: privateKey,
		Logger:  commonlogger.NewOCRWrapper(lggr, true, func(string) {}),

		V2ListenAddresses:    []string{listenAddr},
		V2DeltaReconcile:     time.Second * 5,
		V2DeltaDial:          time.Second * 5,
		V2DiscovererDatabase: common.NewInMemoryDiscovererDatabase(),
		V2EndpointConfig: networking.EndpointConfigV2{
			IncomingMessageBufferSize: 1000000,
			OutgoingMessageBufferSize: 1000000,
		},
	}
	peer, err := networking.NewPeer(peerConfig)
	if err != nil {
		lggr.Error("error creating peer:", err)
		return
	}

	if *isBootstrap {
		// Create a Bootstrapper object.
		bootstrapper, err2 := peer.OCR2BootstrapperFactory().NewBootstrapper(
			types.ConfigDigest{},
			config.NodePeerIDsStr,
			config.BootstrapperLocators,
			1,
		)
		if err2 != nil {
			lggr.Error("error creating endpoint:", err)
			return
		}
		_ = bootstrapper.Start()
	} else {
		// Create an Endpoint object, which opens RageP2P's Streams with all other peers.
		// TODO: How can we implement out own endpoint without access to concretePeerV2 internals?
		// TODO: Can libocr expose a more generic interface for sending messages and opening streams?
		//       For example, support opening streams on demand.
		endpoint, err2 := peer.OCR2BinaryNetworkEndpointFactory().NewEndpoint(
			types.ConfigDigest{},
			config.NodePeerIDsStr,
			config.BootstrapperLocators,
			1,
			types.BinaryNetworkEndpointLimits{
				MaxMessageLength:          100000,
				MessagesRatePerOracle:     2.0,
				MessagesCapacityPerOracle: 10,
				BytesRatePerOracle:        100000.0,
				BytesCapacityPerOracle:    100000,
			},
		)
		if err2 != nil {
			lggr.Error("error creating endpoint:", err)
			return
		}
		_ = endpoint.Start()

		shutdownWaitGroup.Add(2)
		go sendLoop(ctx, &shutdownWaitGroup, endpoint, len(config.NodePeerIDs))
		go recvLoop(ctx, &shutdownWaitGroup, endpoint.Receive(), lggr)
	}

	<-ctx.Done()
	err = peer.Close()
	if err != nil {
		lggr.Error("error closing peer:", err)
	}
	shutdownWaitGroup.Wait()
}

func sendLoop(ctx context.Context, shutdownWaitGroup *sync.WaitGroup, endpoint commontypes.BinaryNetworkEndpoint, n int) {
	defer shutdownWaitGroup.Done()
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	lastId := 0
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			endpoint.SendTo([]byte("hello!"), commontypes.OracleID(lastId%n))
			lastId++
		}
	}
}

func recvLoop(ctx context.Context, shutdownWaitGroup *sync.WaitGroup, chRecv <-chan commontypes.BinaryMessageWithSender, lggr logger.Logger) {
	defer shutdownWaitGroup.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-chRecv:
			lggr.Info("received message from ", msg.Sender, " : ", string(msg.Msg))
		}
	}
}
