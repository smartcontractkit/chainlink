package main

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/networking"
	"github.com/smartcontractkit/libocr/offchainreporting2plus"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/por_mock_ocr3plugin/contractconfig"
	"github.com/smartcontractkit/por_mock_ocr3plugin/db"
	"github.com/smartcontractkit/por_mock_ocr3plugin/keyring"
	"github.com/smartcontractkit/por_mock_ocr3plugin/logger"
	por "github.com/smartcontractkit/por_mock_ocr3plugin/por"
	"github.com/smartcontractkit/por_mock_ocr3plugin/transmitter"
)

const basePort int = 1337

const bootstrapperIndex int = 999

const chainID = 11155111 //sepolia

// var destination common.Address = common.HexToAddress("0xc0ffeec0ffeec0ffeec0ffeec0ffeec0ffeec0ff")

var localConfig types.LocalConfig = types.LocalConfig{
	BlockchainTimeout:                  10 * time.Second,
	ContractConfigConfirmations:        1,
	SkipContractConfigConfirmations:    true,
	ContractConfigTrackerPollInterval:  10 * time.Second,
	ContractConfigLoadTimeout:          10 * time.Second,
	ContractTransmitterTransmitTimeout: 10 * time.Second,
	DatabaseTimeout:                    10 * time.Second,
	DefaultMaxDurationInitialization:   10 * time.Second,
}

const url string = "wss://sepolia.infura.io/ws/v3/de4f73b9679f41219d9a0c386367be1b"

func ethClient(rpcURL string) (client *ethclient.Client) {
	// client, err := ethclient.Dial(rpcURL)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	return nil
}

func bootstrapperConfig() networking.PeerConfig {
	return networking.PeerConfig{
		PrivKey:              contractconfig.P2pPrivateKey(bootstrapperIndex),
		Logger:               logger.NewLogger(),
		V2ListenAddresses:    []string{fmt.Sprintf("localhost:%d", basePort)},
		V2AnnounceAddresses:  []string{fmt.Sprintf("127.0.0.1:%d", basePort)},
		V2DeltaReconcile:     5 * time.Second,
		V2DeltaDial:          5 * time.Second,
		V2DiscovererDatabase: nil,
		V2EndpointConfig: networking.EndpointConfigV2{
			IncomingMessageBufferSize: 5,
			OutgoingMessageBufferSize: 5,
		},
		MetricsRegisterer:            prometheus.DefaultRegisterer,
		LatencyMetricsServiceConfigs: nil,
	}
}

func peerConfig(i int) networking.PeerConfig {
	return networking.PeerConfig{
		PrivKey:              contractconfig.P2pPrivateKey(i),
		Logger:               logger.NewLogger(),
		V2ListenAddresses:    []string{fmt.Sprintf("localhost:%d", basePort+1+i)},
		V2AnnounceAddresses:  []string{fmt.Sprintf("127.0.0.1:%d", basePort+1+i)},
		V2DeltaReconcile:     5 * time.Second,
		V2DeltaDial:          5 * time.Second,
		V2DiscovererDatabase: nil,
		V2EndpointConfig: networking.EndpointConfigV2{
			IncomingMessageBufferSize: 5,
			OutgoingMessageBufferSize: 5,
		},
		MetricsRegisterer:            prometheus.DefaultRegisterer,
		LatencyMetricsServiceConfigs: nil,
	}
}

func oracleArgs(i int, fac types.BinaryNetworkEndpointFactory) offchainreporting2plus.OCR3OracleArgs[por.ChainSelector] {
	logger := logger.NewLogger()
	return offchainreporting2plus.OCR3OracleArgs[por.ChainSelector]{
		BinaryNetworkEndpointFactory: fac,
		V2Bootstrappers: []commontypes.BootstrapperLocator{
			{
				PeerID: contractconfig.OracleIdentity(bootstrapperIndex).PeerID,
				Addrs:  []string{fmt.Sprintf("localhost:%d", basePort)},
			},
		},
		ContractConfigTracker: &contractconfig.FakeContractConfigTracker{},
		ContractTransmitter: transmitter.NewBasicContractTransmitter(
			ethClient(url),
			big.NewInt(chainID),
			contractconfig.DestinationAddress(),
			contractconfig.TransmitterPrivateKey(i),
		),
		Database:               db.NewFakeFatabase(),
		LocalConfig:            localConfig,
		Logger:                 logger,
		MonitoringEndpoint:     nil,
		MetricsRegisterer:      prometheus.DefaultRegisterer,
		OffchainConfigDigester: &contractconfig.FakeOffchainConfigDigester{},
		OffchainKeyring: &keyring.DummyOffchainKeyring{
			OffchainPrivateKey:         contractconfig.OffchainPrivateKey(i),
			ConfigEncryptionPrivateKey: contractconfig.ConfigEncryptionPrivateKey(i),
		},
		OnchainKeyring:         &keyring.DummyEVMOnchainKeyring{PrivateKey: contractconfig.OnchainPrivateKey(i)},
		ReportingPluginFactory: &por.PorReportingPluginFactory{Logger: logger},
	}
}

func main() {
	bootstrapperConfig := bootstrapperConfig()
	bootstrapperPeer, err := networking.NewPeer(bootstrapperConfig)
	if err != nil {
		panic(err)
	}
	defer bootstrapperPeer.Close()

	bootstrapper, err := offchainreporting2plus.NewBootstrapper(offchainreporting2plus.BootstrapperArgs{
		BootstrapperFactory:    bootstrapperPeer.OCR2BootstrapperFactory(),
		V2Bootstrappers:        []commontypes.BootstrapperLocator{},
		ContractConfigTracker:  &contractconfig.FakeContractConfigTracker{},
		Database:               db.NewFakeFatabase(),
		LocalConfig:            localConfig,
		Logger:                 logger.NewLogger(),
		MonitoringEndpoint:     nil,
		OffchainConfigDigester: &contractconfig.FakeOffchainConfigDigester{},
	})
	if err != nil {
		panic(err)
	}
	if err := bootstrapper.Start(); err != nil {
		panic(err)
	}
	defer bootstrapper.Close()

	oargss := []offchainreporting2plus.OCR3OracleArgs[por.ChainSelector]{}
	for i := 0; i < 4; i++ {
		peerConfig := peerConfig(i)
		peer, err := networking.NewPeer(peerConfig)
		if err != nil {
			panic(err)
		}
		defer peer.Close()

		oargss = append(oargss, oracleArgs(i, peer.OCR2BinaryNetworkEndpointFactory()))
		oracle, err := offchainreporting2plus.NewOracle(oargss[i])
		if err != nil {
			panic(err)
		}
		if err := oracle.Start(); err != nil {
			panic(err)
		}
		defer oracle.Close()
	}
	select {}
	// oargs := oracleArgs(2, peer.OCR2BinaryNetworkEndpointFactory())
	// <-time.After(10 * time.Second)
	// _ = oargs
}
