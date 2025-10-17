package cre

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"math"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/fake"
	libcre "github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"

	ttypes "github.com/smartcontractkit/chainlink/system-tests/tests/test-helpers/configuration"
)

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
		priceURL := setupFakeBillingPriceProvider(t, testEnv.Config.Fake)

		t.Cleanup(func() {
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

		// set env vars for billing config
		cache, err := loadWorkflowRegistryCache(relativePathToRepoRoot)
		if err != nil {
			return errors.Wrap(err, "failed to load workflow registry cache")
		}

		if len(testEnv.WrappedBlockchainOutputs) == 0 {
			return errors.New("no blockchain outputs found in the test environment")
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
		os.Setenv("STREAMS_API_URL", priceURL)
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

			rpcURL := strings.Replace(selectedChain.Nodes[0].ExternalHTTPUrl, "127.0.0.1", "host.docker.internal", 1)

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
	Credits  float64
	Reserved float64
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
	require.Greater(t, credits[0].Credits, float64(0.0), "expected initial credits to be greater than 0")

	return billingAssertionState{
		Credits:  credits[0].Credits,
		Reserved: credits[0].Reserved,
		DB:       db,
	}
}

func assertBillingStateChanged(t *testing.T, initial billingAssertionState, timeout time.Duration, expectedMinChange float64) {
	t.Helper()

	// set up a connection to the billing database and run query until data exists
	const pollInterval = 2 * time.Second
	assert.Eventually(t, func() bool {
		finalCredits := queryCredits(t, initial.DB)

		if len(finalCredits) != 1 {
			return false
		}

		credit := finalCredits[0]

		framework.L.Info().
			Float64("final_credits", credit.Credits).
			Float64("initial_credits", initial.Credits).
			Float64("final_reserved", credit.Reserved).
			Float64("initial_reserved", initial.Reserved).
			Msg("checking billing credits")

		// if no credits reserved and no change in credits; nothing was billed
		if credit.Credits == initial.Credits && credit.Reserved == initial.Reserved {
			return false
		}

		if expectedMinChange > 0 {
			creditDiff := math.Floor(initial.Credits - credit.Credits)

			// credits should have decreased by at least expectedMinChange and there should be no reserved credits
			if creditDiff < expectedMinChange || credit.Reserved > 0 {
				return false
			}
		}

		return true
	}, timeout, pollInterval)
}

type billingCredit struct {
	Credits   float64
	Reserved  float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func queryCredits(t *testing.T, db *sql.DB) []billingCredit {
	t.Helper()

	query := "SELECT credits, credits_reserved, created_at, updated_at FROM billing_platform.organization_credits WHERE organization_id = 'integration-test-aggregation-org-happy-path-odd-quorum'"
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

func setupFakeBillingPriceProvider(t *testing.T, input *fake.Input) string {
	t.Helper()

	fakeProviderStarted.Do(func() {
		_, err := fake.NewFakeDataProvider(input)
		require.NoError(t, err)
	})

	host := framework.HostDockerInternal()
	url := fmt.Sprintf("%s:%d", host, input.Port)
	err := fake.Func("GET", "/api/v1/reports/bulk", func(c *gin.Context) {
		ids := strings.Split(c.Query("feedIDs"), ",")
		if len(ids) != 1 {
			c.Data(http.StatusBadRequest, "text/plain", []byte("feedIDs parameter is required"))
			return
		}

		if ids[0] == "" {
			c.Data(http.StatusBadRequest, "text/plain", []byte("mock price handler only supports one feedID"))
			return
		}

		prices := feedPriceResponses["GET /api/v1/reports/bulk"]

		priceDataStr, exists := prices[ids[0]]
		if !exists {
			c.Data(http.StatusNotFound, "text/plain", []byte("price not found"))
			return
		}

		priceData, decodeErr := hex.DecodeString(priceDataStr) // just to verify it's valid hex
		if decodeErr != nil {
			c.Data(http.StatusInternalServerError, "text/plain", []byte("failed to decode price data"))
			return
		}

		c.Data(http.StatusOK, "application/json", priceData)
	})

	require.NoError(t, err)

	return url
}
