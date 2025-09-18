package cre

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	libcre "github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"

	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

var defaultPriceProvider *billingPriceProvider

var feedPriceResponses = map[string]map[string]string{
	"GET /api/v1/reports/bulk": {
		"0x000359843a543ee2fe414dc14c7e7920ef10f4372990b79d6361cdc0dd1ba782": "7b227265706f727473223a5b7b226665656449" +
			"44223a22307830303033353938343361353433656532666534313464633134633765373932306566313066343337323939306237" +
			"396436333631636463306464316261373832222c2276616c696446726f6d54696d657374616d70223a313735383034393034372c" +
			"226f62736572766174696f6e7354696d657374616d70223a313735383034393034372c2266756c6c5265706f7274223a22307830" +
			"30303930643965386439363736356130633439653033613661653035633832653866386465373063663137396261613633326631" +
			"38333133653534626436393030303030303030303030303030303030303030303030303030303030303030303030303030303030" +
			"30303030303030303030303030303030316530346239353030303030303030303030303030303030303030303030303030303030" +
			"30303030303030303030303030303030303030303030303030303330303030303030313030303030303030303030303030303030" +
			"30303030303030303030303030303030303030303030303030303030303030303030303030303030303030303065303030303030" +
			"30303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030" +
			"30303030323230303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030" +
			"30303030303030303030303030303030323830303030313030303030303030303030303030303030303030303030303030303030" +
			"30303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030" +
			"30303030303030303030303030303030303030303030303030303030303030303030303030303030313230303030333539383433" +
			"61353433656532666534313464633134633765373932306566313066343337323939306237396436333631636463306464316261" +
			"37383230303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030" +
			"30303030303030363863396233313730303030303030303030303030303030303030303030303030303030303030303030303030" +
			"30303030303030303030303030303030303030363863396233313730303030303030303030303030303030303030303030303030" +
			"30303030303030303030303030303030303030303030303030303034306638666365363933336530303030303030303030303030" +
			"30303030303030303030303030303030303030303030303030303030303030303030303030333035613835363631626162363730" +
			"30303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030" +
			"30303036386631343031373030303030303030303030303030303030303030303030303030303030303030303030303030303030" +
			"30303030306632643432396137316632346334303030303030303030303030303030303030303030303030303030303030303030" +
			"30303030303030303030303030303030306632643162386530346166623335643030303030303030303030303030303030303030" +
			"30303030303030303030303030303030303030303030303030303030306632643538663730323331396565396534303030303030" +
			"30303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030" +
			"30303030303032663336623231613263613833653565646336353233373734303034376363306364643363346231656362653161" +
			"38346334306237396663643732633265653563363430613738646161396264666439346462353036303964393030613437373934" +
			"39363232313932366135353231356466303230623935366237366639366235303030303030303030303030303030303030303030" +
			"30303030303030303030303030303030303030303030303030303030303030303030303030303030303032363064326264646139" +
			"64613337646131346464376634343238646563336661653964656138313931396638346436386338316132356139663435316537" +
			"31303737636262653036643937643264373835636166343232393335356631396162323137366338306533323535383930356535" +
			"336366306166646239366232366439227d5d7d",
	},
}

func loadBillingStackCache(relativePathToRepoRoot string) (*config.BillingConfig, error) {
	c := &config.BillingConfig{}
	if loadErr := c.Load(config.MustBillingStateFileAbsPath(relativePathToRepoRoot)); loadErr != nil {
		return nil, errors.Wrap(loadErr, "failed to load billing stack cache")
	}

	return c, nil
}

func startBillingStackIfIsNotRunning(t *testing.T, relativePathToRepoRoot, environmentDir string, testEnv *ttypes.TestEnvironment) error {
	if !config.BillingStateFileExists(relativePathToRepoRoot) {
		defaultPriceProvider = newBillingPriceProvider(t)

		t.Cleanup(func() {
			defaultPriceProvider.Close()

			/*
				cmd := exec.Command("go", "run", ".", "env", "billing", "stop")
				cmd.Dir = environmentDir
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmdErr := cmd.Run()
				if cmdErr != nil {
					t.Log("failed to stop Billing Platform Service:", cmdErr)
				}
			*/
		})

		defaultPriceProvider.Start()

		// set env vars for billing config
		cache, err := loadWorkflowRegistryCache(relativePathToRepoRoot)
		if err != nil {
			return errors.Wrap(err, "failed to load workflow registry cache")
		}

		if len(testEnv.WrappedBlockchainOutputs) == 0 {
			return errors.New("no blockchain outputs found in the test environment")
		}

		replaceHost := func(url string) string {
			return strings.Replace(url, "127.0.0.1", "host.docker.internal", 1)
		}

		for _, ref := range testEnv.EnvArtifact.AddressRefs {
			switch ref.Type {
			case "WorkflowRegistry":
				if cache.ChainSelector == ref.ChainSelector {
					os.Setenv("MAINNET_WORKFLOW_REGISTRY_CONTRACT_ADDRESS", ref.Address)
				}
			case "CapabilitiesRegistry":
				if cache.ChainSelector == ref.ChainSelector {
					os.Setenv("MAINNET_CAPABILITIES_REGISTRY_CONTRACT_ADDRESS", ref.Address)
				}
			default:
				continue
			}
		}

		os.Setenv("MAINNET_WORKFLOW_REGISTRY_CHAIN_SELECTOR", strconv.FormatUint(cache.ChainSelector, 10))
		os.Setenv("MAINNET_CAPABILITIES_REGISTRY_CHAIN_SELECTOR", strconv.FormatUint(cache.ChainSelector, 10))
		os.Setenv("STREAMS_API_URL", replaceHost(defaultPriceProvider.URL()))
		os.Setenv("STREAMS_API_KEY", "cannot be empty")
		os.Setenv("STREAMS_API_SECRET", "cannot be empty")
		os.Setenv("TEST_OWNERS", strings.Join(cache.WorkflowOwnersStrings(), ","))

		// Select the appropriate chain for billing service from available chains in the environment.
		// otherwise, if RPCURL is defined, billing service can be used standalone
		if len(testEnv.WrappedBlockchainOutputs) != 0 {
			var selectedChain *blockchain.Output

			for _, chain := range testEnv.WrappedBlockchainOutputs {
				if chain.ChainSelector == cache.ChainSelector {
					selectedChain = chain.BlockchainOutput
				}
			}

			if selectedChain == nil || len(selectedChain.Nodes) == 0 {
				return errors.Wrap(err, fmt.Sprintf("configured chain selector does not exist in the current topology: %d", cache.ChainSelector))
			}

			rpcURL := replaceHost(selectedChain.Nodes[0].ExternalHTTPUrl)

			os.Setenv("MAINNET_WORKFLOW_REGISTRY_RPC_URL", rpcURL)
			os.Setenv("MAINNET_CAPABILITIES_REGISTRY_RPC_URL", rpcURL)
		}

		framework.L.Info().Str("state file", config.MustBillingStateFileAbsPath(relativePathToRepoRoot)).Msg("Billing state file was not found. Starting Billing...")
		cmd := exec.Command("go", "run", ".", "env", "billing", "start")
		cmd.Dir = environmentDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmdErr := cmd.Run()
		if cmdErr != nil {
			return errors.Wrap(cmdErr, "failed to start Billing Platform Service")
		}
	}
	framework.L.Info().Msg("Billing Platform Service is running.")
	return nil
}

func loadWorkflowRegistryCache(relativePathToRepoRoot string) (*libcre.WorkflowRegistryOutput, error) {
	previousCTFconfigs := os.Getenv("CTF_CONFIGS")
	defer func() {
		setErr := os.Setenv("CTF_CONFIGS", previousCTFconfigs)
		if setErr != nil {
			framework.L.Warn().Err(setErr).Msg("failed to restore previous CTF_CONFIGS env var")
		}
	}()

	setErr := os.Setenv("CTF_CONFIGS", config.MustWorkflowRegistryStateFileAbsPath(relativePathToRepoRoot))
	if setErr != nil {
		return nil, errors.Wrap(setErr, "failed to set CTF_CONFIGS env var")
	}

	return framework.Load[libcre.WorkflowRegistryOutput](nil)
}

type billingAssertionState struct {
	Credits  float32
	Reserved float32
	DB       *sql.DB
}

func getBillingAssertionState(t *testing.T, relativePathToRepoRoot string) billingAssertionState {
	t.Helper()

	billingConfig, err := loadBillingStackCache(relativePathToRepoRoot)
	require.NoError(t, err, "failed to load billing config")

	dsn := billingConfig.BillingService.Output.Postgres.DSN
	db, err := sql.Open("pgx", dsn)
	require.NoError(t, err, "failed to connect to billing database")

	credits := queryCredits(t, db)
	require.Len(t, credits, 1, "expected one row in organization_credits table")
	require.Greater(t, credits[0].Credits, float32(0.0), "expected initial credits to be greater than 0")

	return billingAssertionState{
		Credits:  credits[0].Credits,
		Reserved: credits[0].Reserved,
		DB:       db,
	}
}

func assertBillingStateChanged(t *testing.T, initial billingAssertionState, timeout time.Duration, expectedMinChange float32) {
	t.Helper()

	// set up a connection to the billing database and run query until data exists
	assert.Eventually(t, func() bool {
		finalCredits := queryCredits(t, initial.DB)

		if len(finalCredits) != 1 {
			return false
		}

		credit := finalCredits[0]

		framework.L.Info().
			Float32("final_credits", credit.Credits).
			Float32("initial_credits", initial.Credits).
			Float32("final_reserved", credit.Reserved).
			Float32("initial_reserved", initial.Reserved).
			Msg("checking billing credits")

		// if no credits reserved and no change in credits; nothing was billed
		if credit.Credits == initial.Credits && credit.Reserved == initial.Reserved {
			return false
		}

		if expectedMinChange > 0 {
			creditDiff := float32(math.Floor(float64(initial.Credits - credit.Credits)))
			reservDiff := float32(math.Floor(float64(initial.Reserved - credit.Reserved)))

			if creditDiff < expectedMinChange && reservDiff < expectedMinChange {
				return false
			}
		}

		return true
	}, timeout, 10*time.Second)
}

type billingCredit struct {
	Credits   float32
	Reserved  float32
	CreatedAt time.Time
	UpdatedAt time.Time
}

func queryCredits(t *testing.T, db *sql.DB) []billingCredit {
	t.Helper()

	query := "SELECT credits, credits_reserved, created_at, updated_at FROM billing_platform.organization_credits WHERE organization_id = '000000000000'"
	rows, err := db.QueryContext(t.Context(), query)
	require.NoError(t, err, "failed to query billing database")

	defer func() {
		rows.Close()
		assert.NoError(t, rows.Err(), "error occurred during rows iteration")
	}()

	// query the billing database for a baseline data reference
	credits := []billingCredit{}

	for rows.Next() {
		var credit billingCredit

		scanErr := rows.Scan(&credit.Credits, &credit.Reserved, &credit.CreatedAt, &credit.UpdatedAt)
		require.NoError(t, scanErr, "failed to scan row from billing database")

		credits = append(credits, credit)
	}

	return credits
}

type billingPriceProvider struct {
	t *testing.T

	server   *httptest.Server
	handlers map[string]http.HandlerFunc // handlers with key in form METHOD PATH: ex: "GET /prices"
}

func newBillingPriceProvider(t *testing.T) *billingPriceProvider {
	t.Helper()

	provider := &billingPriceProvider{
		t:        t,
		handlers: map[string]http.HandlerFunc{},
	}

	for path, prices := range feedPriceResponses {
		provider.handlers[path] = provider.makePriceHandler(prices)
	}

	provider.server = httptest.NewUnstartedServer(nil)

	return provider
}

func (b *billingPriceProvider) URL() string {
	b.t.Helper()
	return b.server.URL
}

func (b *billingPriceProvider) Start() {
	b.t.Helper()

	b.server.Config.Handler = b.createHandler()
	b.server.Start()
}

func (b *billingPriceProvider) Close() {
	b.t.Helper()
	b.server.Close()
}

func (b *billingPriceProvider) createHandler() http.Handler {
	b.t.Helper()

	mux := http.NewServeMux()
	for key, handler := range b.handlers {
		mux.HandleFunc(key, handler)
	}

	return mux
}

func (b *billingPriceProvider) makePriceHandler(prices map[string]string) http.HandlerFunc {
	b.t.Helper()

	return func(writer http.ResponseWriter, request *http.Request) {
		if err := request.ParseForm(); err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		ids := strings.Split(request.FormValue("feedIDs"), ",")
		if len(ids) != 1 {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte("feedIDs parameter is required"))
			return
		}

		if ids[0] == "" {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte("mock price handler only supports one feedID"))
			return
		}

		priceDataStr, exists := prices[ids[0]]
		if !exists {
			writer.WriteHeader(http.StatusNotFound)
			return
		}

		priceData, err := hex.DecodeString(priceDataStr) // just to verify it's valid hex
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, _ = writer.Write(priceData)
	}
}
