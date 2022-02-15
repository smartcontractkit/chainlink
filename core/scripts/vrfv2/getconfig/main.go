package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kelseyhightower/envconfig"
	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_coordinator_v2"
)

type config struct {
	EthURLs           URLMap                    `envconfig:"ETH_URLS"`
	ChainIDs          map[string]int64          `envconfig:"CHAIN_IDS"`
	VRFV2Coordinators map[string]common.Address `envconfig:"VRFV2_COORDINATORS"`
}

type URLMap map[string]string

func (u *URLMap) Decode(input string) error {
	out := make(URLMap)
	// split on commas first to extract name/url pairs
	items := strings.Split(input, ",")
	for _, item := range items {
		// split on colon
		split := strings.SplitN(item, ":", 2)
		network := split[0]
		url := split[1]
		out[network] = url
	}
	*u = out
	return nil
}

type coordinatorConfig struct {
	CoordinatorAddress string                          `toml:"coordinatorAddress"`
	LinkAddress        string                          `toml:"linkAddress"`
	BHSAddress         string                          `toml:"bhsAddress"`
	LinkEthAddress     string                          `toml:"linkEthFeedAddress"`
	BaseConfig         vrf_coordinator_v2.GetConfig    `toml:"baseConfig"`
	FeeConfig          vrf_coordinator_v2.GetFeeConfig `toml:"feeConfig"`
}

func getEthClients(urls map[string]string) (clients map[string]*ethclient.Client, err error) {
	clients = make(map[string]*ethclient.Client)
	for network, url := range urls {
		client, err := ethclient.Dial(url)
		if err != nil {
			return nil, errors.Wrapf(err, "ethclient dial %s (network: %s)", url, network)
		}
		clients[network] = client
	}
	return
}

func getVRFV2Coordinators(
	addresses map[string]common.Address,
	clients map[string]*ethclient.Client,
) (coordinators map[string]*vrf_coordinator_v2.VRFCoordinatorV2, err error) {
	coordinators = make(map[string]*vrf_coordinator_v2.VRFCoordinatorV2)
	for network, address := range addresses {
		client, ok := clients[network]
		if !ok {
			return nil, fmt.Errorf("no eth client available for network '%s', did you forget to provide one?", network)
		}
		coordinator, err := vrf_coordinator_v2.NewVRFCoordinatorV2(address, client)
		if err != nil {
			return nil, errors.Wrap(err, "create vrf coordinator v2")
		}
		coordinators[network] = coordinator
	}
	return
}

func getCoordinatorConfigs(coordinators map[string]*vrf_coordinator_v2.VRFCoordinatorV2) (configs map[string]coordinatorConfig, err error) {
	configs = make(map[string]coordinatorConfig)
	for network, coordinator := range coordinators {
		baseConfig, err := coordinator.GetConfig(&bind.CallOpts{Context: context.Background()})
		if err != nil {
			return nil, errors.Wrap(err, "get base config")
		}

		feeConfig, err := coordinator.GetFeeConfig(&bind.CallOpts{Context: context.Background()})
		if err != nil {
			return nil, errors.Wrap(err, "get fee config")
		}

		bhs, err := coordinator.BLOCKHASHSTORE(&bind.CallOpts{Context: context.Background()})
		if err != nil {
			return nil, errors.Wrap(err, "get bhs address")
		}

		link, err := coordinator.LINK(&bind.CallOpts{Context: context.Background()})
		if err != nil {
			return nil, errors.Wrap(err, "get link address")
		}

		feed, err := coordinator.LINKETHFEED(&bind.CallOpts{Context: context.Background()})
		if err != nil {
			return nil, errors.Wrap(err, "get linketh address")
		}

		configs[network] = coordinatorConfig{
			CoordinatorAddress: coordinator.Address().Hex(),
			BHSAddress:         bhs.Hex(),
			LinkAddress:        link.Hex(),
			LinkEthAddress:     feed.Hex(),
			BaseConfig:         baseConfig,
			FeeConfig:          feeConfig,
		}
	}
	return
}

func main() {
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err)
	}

	clients, err := getEthClients(cfg.EthURLs)
	if err != nil {
		log.Fatal(err)
	}

	coordinators, err := getVRFV2Coordinators(cfg.VRFV2Coordinators, clients)
	if err != nil {
		log.Fatal(err)
	}

	configs, err := getCoordinatorConfigs(coordinators)
	if err != nil {
		log.Fatal(err)
	}

	t, err := toml.Marshal(configs)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(t))
}
