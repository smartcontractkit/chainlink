package cre

import (
	"database/sql"
	"math"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
)

func loadBillingStackCache(relativePathToRepoRoot string) (*config.BillingConfig, error) {
	c := &config.BillingConfig{}
	if loadErr := c.Load(config.MustBillingStateFileAbsPath(relativePathToRepoRoot)); loadErr != nil {
		return nil, errors.Wrap(loadErr, "failed to load billing stack cache")
	}

	return c, nil
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
