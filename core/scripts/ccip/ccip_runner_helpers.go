package main

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/smartcontractkit/chainlink/core/scripts/ccip/dione"
	"github.com/smartcontractkit/chainlink/core/scripts/ccip/rhea"
	"github.com/smartcontractkit/chainlink/core/scripts/ccip/rhea/deployments"
	"github.com/smartcontractkit/chainlink/core/scripts/ccip/shared"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/evm_2_evm_onramp"
)

var laneMapping = map[dione.Environment]map[rhea.Chain]map[rhea.Chain]rhea.EvmDeploymentConfig{
	dione.StagingAlpha: deployments.AlphaChainMapping,
	dione.StagingBeta:  deployments.BetaChainMapping,
	dione.Production:   deployments.ProdChainMapping,
}

var chainMapping = map[dione.Environment]map[rhea.Chain]rhea.EvmDeploymentConfig{
	dione.StagingAlpha: deployments.AlphaChains,
	dione.StagingBeta:  deployments.BetaChains,
	dione.Production:   deployments.ProdChains,
}

func checkOwnerKey(t *testing.T) string {
	ownerKey := os.Getenv("OWNER_KEY")
	if ownerKey == "" {
		t.Log("No key given, this test will be skipped. This is intended behaviour for automated testing.")
		t.SkipNow()
	}

	return ownerKey
}

// DoForEachChain can be run concurrently as all calls target a single chain
func DoForEachChain(t *testing.T, env dione.Environment, f func(chain rhea.EvmDeploymentConfig)) {
	ownerKey := checkOwnerKey(t)
	var wg sync.WaitGroup
	for chnName, chn := range chainMapping[env] {
		wg.Add(1)
		go func(chainName rhea.Chain, chain rhea.EvmDeploymentConfig) {
			defer wg.Done()
			t.Logf("Running function for chain %s", chainName)
			chain.SetupChain(t, ownerKey)
			f(chain)
		}(chnName, chn)
	}
	wg.Wait()
}

func DoForEachLane(t *testing.T, env dione.Environment, f func(source rhea.EvmDeploymentConfig, destination rhea.EvmDeploymentConfig)) {
	ownerKey := checkOwnerKey(t)
	for sourceChain, sourceMap := range laneMapping[env] {
		for destChain := range sourceMap {
			t.Logf("Running function for lane %s -> %s", sourceChain, destChain)

			source := laneMapping[env][sourceChain][destChain]
			dest := laneMapping[env][destChain][sourceChain]

			source.SetupChain(t, ownerKey)
			dest.SetupChain(t, ownerKey)

			f(source, dest)
		}
	}
}

func DoForEachBidirectionalLane(t *testing.T, env dione.Environment, f func(source rhea.EvmDeploymentConfig, destination rhea.EvmDeploymentConfig)) {
	ownerKey := checkOwnerKey(t)
	completed := make(map[rhea.Chain]map[rhea.Chain]interface{})

	for sourceChain, sourceMap := range laneMapping[env] {
		for destChain := range sourceMap {
			// Skip if we already processed the lane from the other side
			if destMap, ok := completed[destChain]; ok {
				if _, ok := destMap[sourceChain]; ok {
					continue
				}
			}

			t.Logf("Running function for lane %s <-> %s", sourceChain, destChain)

			source := laneMapping[env][sourceChain][destChain]
			dest := laneMapping[env][destChain][sourceChain]

			source.SetupChain(t, ownerKey)
			dest.SetupChain(t, ownerKey)

			f(source, dest)

			if _, ok := completed[sourceChain]; !ok {
				completed[sourceChain] = make(map[rhea.Chain]interface{})
			}
			if _, ok := completed[destChain]; !ok {
				completed[destChain] = make(map[rhea.Chain]interface{})
			}
			completed[sourceChain][destChain] = true
			completed[destChain][sourceChain] = true
		}
	}
}

// WaitForCrossChainSendRequest checks on chain for a successful onramp send event with the given tx hash.
// If not immediately found it will keep retrying in intervals of the globally specified RetryTiming.
func WaitForCrossChainSendRequest(source SourceClient, fromBlockNum uint64, txhash common.Hash) *evm_2_evm_onramp.EVM2EVMOnRampCCIPSendRequested {
	filter := bind.FilterOpts{Start: fromBlockNum}
	source.logger.Infof("Waiting for cross chain send... ")

	for {
		iterator, err := source.OnRamp.FilterCCIPSendRequested(&filter)
		helpers.PanicErr(err)
		for iterator.Next() {
			if iterator.Event.Raw.TxHash.Hex() == txhash.Hex() {
				source.logger.Infof("Cross chain send event found in tx: %s ", helpers.ExplorerLink(int64(source.ChainId), txhash))
				return iterator.Event
			}
		}
		time.Sleep(shared.RetryTiming)
	}
}

func GetCurrentBlockNumber(chain *ethclient.Client) uint64 {
	blockNumber, err := chain.BlockNumber(context.Background())
	helpers.PanicErr(err)
	return blockNumber
}
