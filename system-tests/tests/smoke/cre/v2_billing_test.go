package cre

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	crontypes "github.com/smartcontractkit/chainlink/core/scripts/cre/environment/examples/workflows/v2/cron/types"
)

func ExecuteBillingTest(t *testing.T, testEnv *TestEnvironment) {
	testLogger := framework.L
	timeout := 2 * time.Minute
	workflowFileLocation := "../../../../core/scripts/cre/environment/examples/workflows/v2/cron/main.go"
	workflowName := "cronbilling"

	priceProvider := newBillingPriceProvider(t)
	priceProvider.middleware = append(priceProvider.middleware, func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			testLogger.Info().Msg("Billing price provider received a request: " + r.Method + " " + r.URL.Path)
			h.ServeHTTP(w, r)
		})
	})
	t.Cleanup(func() {
		priceProvider.Close()
	})

	priceProvider.Start()

	require.NoError(
		t,
		startBillingStackIfIsNotRunning(testEnv.TestConfig.RelativePathToRepoRoot, testEnv.TestConfig.EnvironmentDirPath, priceProvider.URL(), testEnv),
		"failed to start Billing stack",
	)

	billingConfig, err := loadBillingStackCache(testEnv.TestConfig.RelativePathToRepoRoot)
	require.NoError(t, err, "failed to load billing config")

	dsn := billingConfig.BillingService.Output.Postgres.DSN
	db, err := sql.Open("pgx", dsn)
	require.NoError(t, err, "failed to connect to billing database")

	credits := queryCredits(t, db)

	require.Len(t, credits, 1, "expected one row in organization_credits table")
	require.Greater(t, credits[0].Credits, float32(0.0), "expected initial credits to be greater than 0")

	initialCredits := credits[0]

	testLogger.Info().Msg("Creating Cron workflow configuration file...")
	workflowConfig := crontypes.WorkflowConfig{
		Schedule: "*/30 * * * * *", // every 30 seconds
	}
	compileAndDeployWorkflow(t, testEnv, testLogger, workflowName, &workflowConfig, workflowFileLocation)

	// set up a connection to the billing database and run query until data exists
	assert.Eventually(t, func() bool {
		finalCredits := queryCredits(t, db)

		if len(finalCredits) != 1 {
			return false
		}

		credit := finalCredits[0]

		// if no credits reserved and no change in credits; nothing was billed
		if credit.Credits == initialCredits.Credits && credit.Reserved == initialCredits.Reserved {
			return false
		}

		testLogger.Info().Msg(fmt.Sprintf("Final credits: %+v", finalCredits))

		return true
	}, timeout, 10*time.Second)

	testLogger.Info().Msg("Billing test completed")
}

type billingPriceProvider struct {
	t *testing.T

	server     *httptest.Server
	handlers   map[string]http.HandlerFunc // handlers with key in form METHOD PATH: ex: "GET /prices"
	middleware []func(http.Handler) http.Handler
}

func newBillingPriceProvider(t *testing.T) *billingPriceProvider {
	t.Helper()

	provider := &billingPriceProvider{
		t: t,
	}

	provider.server = httptest.NewUnstartedServer(provider.handler(t))

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

	var last http.Handler

	if len(b.handlers) == 0 {
		last = b.handler(b.t)
	} else {
		mux := http.NewServeMux()
		for key, handler := range b.handlers {
			mux.HandleFunc(key, handler)
		}

		last = mux
	}

	for x := len(b.middleware) - 1; x >= 0; x-- {
		last = b.middleware[x](last)
	}

	return last
}

func (b *billingPriceProvider) handler(t *testing.T) http.HandlerFunc {
	t.Helper()

	return func(writer http.ResponseWriter, request *http.Request) {
		// Handle incoming requests and provide billing prices
		key := request.Method + " " + request.URL.Path
		if handler, exists := b.handlers[key]; exists {
			handler(writer, request)

			return
		}

		http.NotFound(writer, request)
	}
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
