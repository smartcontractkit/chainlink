package internal_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"testing"
	"time"

	"chainlink/core/internal/cltest"
	"chainlink/core/internal/gethwrappers/generated/flux_aggregator_wrapper"
	"chainlink/core/internal/gethwrappers/generated/link_token_interface"
	"chainlink/core/store/models"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	goEthereumEth "github.com/ethereum/go-ethereum/eth"
	"github.com/stretchr/testify/require"
)

// fluxAggregator represents the universe with which the aggregator contract
// interacts
type fluxAggregator struct {
	aggregatorContract        *flux_aggregator_wrapper.FluxAggregator
	aggregatorContractAddress common.Address
	linkContract              *link_token_interface.LinkToken
	// Abstraction representation of the ethereum blockchain
	backend       *backends.SimulatedBackend
	aggregatorABI abi.ABI
	// Cast of participants
	sergey  *bind.TransactOpts // Owns all the LINK initially
	neil    *bind.TransactOpts // Node operator Flux Monitor Oracle
	ned     *bind.TransactOpts // Node operator Flux Monitor Oracle
	nallory *bind.TransactOpts // Node operator Flux Monitor Oracle (Baddie.)
}

// newIdentity returns a go-ethereum abstraction of an ethereum account for
// interacting with contract golang wrappers
func newIdentity(t *testing.T) *bind.TransactOpts {
	key, err := crypto.GenerateKey()
	require.NoError(t, err, "failed to generate ethereum identity")
	return bind.NewKeyedTransactor(key)
}

// deployFluxAggregator returns a fully initialized fluxAggregator universe. The
// arguments match the arguments of the same name in the FluxAggregator
// constructor.
func deployFluxAggregator(t *testing.T, paymentAmount *big.Int, timeout uint32,
	decimals uint8, description [32]byte) fluxAggregator {
	var f fluxAggregator
	f.sergey = newIdentity(t)
	f.neil = newIdentity(t)
	f.ned = newIdentity(t)
	f.nallory = cltest.OracleTransactor
	oneEth := big.NewInt(1000000000000000000)
	genesisData := core.GenesisAlloc{
		f.sergey.From:  {Balance: oneEth},
		f.neil.From:    {Balance: oneEth},
		f.ned.From:     {Balance: oneEth},
		f.nallory.From: {Balance: oneEth},
	}
	gasLimit := goEthereumEth.DefaultConfig.Miner.GasCeil
	f.backend = backends.NewSimulatedBackend(genesisData, gasLimit)
	var err error
	f.aggregatorABI, err = abi.JSON(strings.NewReader(
		flux_aggregator_wrapper.FluxAggregatorABI))
	require.NoError(t, err, "could not parse FluxAggregator ABI")
	var linkAddress common.Address
	linkAddress, _, f.linkContract, err = link_token_interface.DeployLinkToken(
		f.sergey, f.backend)
	require.NoError(t, err,
		"failed to deploy link contract to simulated ethereum blockchain")
	f.backend.Commit()
	// FluxAggregator contract subtracts timeout from block timestamp, which will
	// be less than the timeout, leading to a SafeMath error. Wait for longer than
	// the timeout... Golang is unpleasant about mixing int64 and time.Duration in
	// arithmetic operations, so do everything as int64 and then convert.
	waitTimeMs := int64(timeout * 1000)
	time.Sleep(time.Duration((waitTimeMs + waitTimeMs/20) * int64(time.Millisecond)))
	f.aggregatorContractAddress, _, f.aggregatorContract, err =
		flux_aggregator_wrapper.DeployFluxAggregator(f.sergey, f.backend,
			linkAddress, paymentAmount, timeout, decimals, description)
	f.backend.Commit() // Must commit contract to chain before we can fund with LINK
	require.NoError(t, err,
		"failed to deploy FluxAggregator contract to simulated ethereum blockchain")
	_, err = f.linkContract.Transfer(f.sergey, f.aggregatorContractAddress,
		oneEth) // Actually, LINK
	require.NoError(t, err, "failed to fund FluxAggregator contract with LINK")
	_, err = f.aggregatorContract.UpdateAvailableFunds(f.sergey)
	require.NoError(t, err, "failed to update aggregator's vailableFunds field")
	oracleList := []common.Address{f.neil.From, f.ned.From, f.nallory.From}
	_, err = f.aggregatorContract.AddOracles(
		f.sergey, oracleList, oracleList, 2, 3, 2)
	f.backend.Commit()
	return f
}

//- have the malicious node start the next round. UpdateAnswer with the next round ID.
//- successfully close the round through the submissions of the other nodes
//- have the malicious node try to start another round repeatedly until the roundDelay is reached, making sure that it isn't successful
//- finally, ensure it can start a legitimate round after roundDelay is reached
func TestFluxMonitorAntiSpamLogic(t *testing.T) {
	// Comments starting with "-" describe the steps this test executes.

	//- deploy a brand new FM contract
	var description [32]byte
	copy(description[:], "exactly thirty-two characters!!!")
	fee := big.NewInt(100) // Amount paid by FA contract, in LINK-wei
	fa := deployFluxAggregator(t, fee, 1, 8, description)

	//- have one of the fake nodes start a round. UpdateAnswer
	_, err := fa.aggregatorContract.UpdateAnswer(fa.neil, big.NewInt(1),
		big.NewInt(1))
	require.NoError(t, err, "failed to initialize first flux aggregation round")
	config, cfgCleanup := cltest.NewConfig(t)
	defer cfgCleanup()
	app, cleanup := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t,
		config, fa.backend)
	defer cleanup()
	require.NoError(t, app.StartAndConnect(),
		"failed to start chainlink application")
	minFee := app.Store.Config.MinimumContractPayment().ToInt()
	require.True(t, fee.Cmp(minFee) >= 0, "fee paid by FluxAggregator (%d) must "+
		"at least match MinimumContractPayment (%s). (Which is currently set in "+
		"cltest.go.)", fee, minFee)
	// Have server respond with 102 for price when FM checks external price
	// adapter for deviation. 102 is enough deviation to trigger a job run.
	priceResponse := `{"data":{"result": 102}}`
	mockServer, assertCalled := cltest.NewHTTPMockServer(t, http.StatusOK, "POST",
		priceResponse)
	defer assertCalled()
	// Create FM Job, and wait for job run to start because the above criteria initiates a run.
	buffer := cltest.MustReadFile(t, "testdata/flux_monitor_job.json")
	var job models.JobSpec
	err = json.Unmarshal(buffer, &job)
	require.NoError(t, err)
	job.Initiators[0].InitiatorParams.Feeds = cltest.JSONFromString(t, fmt.Sprintf(`["%s"]`, mockServer.URL))
	job.Initiators[0].InitiatorParams.PollingInterval = models.Duration(15 * time.Second)
	job.Initiators[0].InitiatorParams.Address = fa.aggregatorContractAddress
	j := cltest.CreateJobSpecViaWeb(t, app, job)

	// - successfully close the round through the submissions of the other nodes.
	// UpdateAnswer with the same round ID.
	jrs := cltest.WaitForRuns(t, j, app.Store, 1) // Submit answer from
	fmt.Println("jrs", jrs)
}

// XAU/XAG happened partly because you can update the entire state all at once.
// Having to add oracles one-by-one slows you down, so you can avoid some
// mistakes.
