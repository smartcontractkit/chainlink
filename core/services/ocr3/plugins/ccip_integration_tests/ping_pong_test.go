package ccip_integration_tests

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/stretchr/testify/require"

	"golang.org/x/exp/maps"

	pp "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ping_pong_demo"
)

/*
* Test is setting up 3 chains (let's call them A, B, C), each chain deploys and starts 2 ping pong contracts for the other 2.
* A ---deploy+start---> (pingPongB, pingPongC)
* B ---deploy+start---> (pingPongA, pingPongC)
* C ---deploy+start---> (pingPongA, pingPongB)
* and then checks that each ping pong contract emitted `CCIPSendRequested` event from the expected source to destination.
* Test fails if any wiring between contracts is not correct.
 */
func TestPingPong(t *testing.T) {
	_, universes := createUniverses(t, 3)
	pingPongs := initializePingPongContracts(t, universes)
	for chainID, universe := range universes {
		for otherChain, pingPong := range pingPongs[chainID] {
			t.Log("PingPong From: ", chainID, " To: ", otherChain)
			_, err := pingPong.StartPingPong(universe.owner)
			require.NoError(t, err)
			universe.backend.Commit()

			logIter, err := universe.onramp.FilterCCIPSendRequested(&bind.FilterOpts{Start: 0}, nil)
			require.NoError(t, err)
			// Iterate until latest event
			for logIter.Next() {
			}
			log := logIter.Event
			require.Equal(t, getSelector(otherChain), log.DestChainSelector)
			require.Equal(t, pingPong.Address(), log.Message.Sender)
			chainPingPongAddr := pingPongs[otherChain][chainID].Address().Bytes()
			// With chain agnostic addresses we need to pad the address to the correct length if the receiver is zero prefixed
			paddedAddr := gethcommon.LeftPadBytes(chainPingPongAddr, len(log.Message.Receiver))
			require.Equal(t, paddedAddr, log.Message.Receiver)
		}
	}
}

// InitializeContracts initializes ping pong contracts on all chains and
// connects them all to each other.
func initializePingPongContracts(
	t *testing.T,
	chainUniverses map[uint64]onchainUniverse,
) map[uint64]map[uint64]*pp.PingPongDemo {
	pingPongs := make(map[uint64]map[uint64]*pp.PingPongDemo)
	chainIDs := maps.Keys(chainUniverses)
	// For each chain initialize N ping pong contracts, where N is the (number of chains - 1)
	for chainID, universe := range chainUniverses {
		pingPongs[chainID] = make(map[uint64]*pp.PingPongDemo)
		for _, chainToConnect := range chainIDs {
			if chainToConnect == chainID {
				continue // don't connect chain to itself
			}
			backend := universe.backend
			owner := universe.owner
			pingPongAddr, _, _, err := pp.DeployPingPongDemo(owner, backend, universe.router.Address(), universe.linkToken.Address())
			require.NoError(t, err)
			backend.Commit()
			pingPong, err := pp.NewPingPongDemo(pingPongAddr, backend)
			require.NoError(t, err)
			backend.Commit()
			// Fund the ping pong contract with LINK
			_, err = universe.linkToken.Transfer(owner, pingPong.Address(), e18Mult(10))
			backend.Commit()
			require.NoError(t, err)
			pingPongs[chainID][chainToConnect] = pingPong
		}
	}

	// Set up each ping pong contract to its counterpart on the other chain
	for chainID, universe := range chainUniverses {
		for chainToConnect, pingPong := range pingPongs[chainID] {
			_, err := pingPong.SetCounterpart(
				universe.owner,
				getSelector(chainUniverses[chainToConnect].chainID),
				// This is the address of the ping pong contract on the other chain
				pingPongs[chainToConnect][chainID].Address(),
			)
			require.NoError(t, err)
			universe.backend.Commit()
		}
	}
	return pingPongs
}
