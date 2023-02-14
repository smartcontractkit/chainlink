package migration_test

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-env/environment"
	"github.com/smartcontractkit/chainlink-env/logging"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/testsetups"
)

type Data struct {
	ID        int       `db:"id"`
	Cfg       []byte    `db:"cfg"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Enabled   bool      `db:"enabled"`
}

func TestMain(m *testing.M) {
	logging.Init()
	os.Exit(m.Run())
}

func TestMigrationDatabase(t *testing.T) {
	testEnvironment, err := testsetups.DBMigration(&testsetups.DBMigrationSpec{
		FromSpec: testsetups.FromVersionSpec{
			Image: "public.ecr.aws/chainlink/chainlink",
			Tag:   "1.6.0-nonroot",
		},
		ToSpec: testsetups.ToVersionSpec{
			Image: "public.ecr.aws/chainlink/chainlink",
			Tag:   "1.7.1-nonroot",
		},
	})
	require.NoError(t, err, "Error setting up DBMigration test")
	// if test haven't failed after that assertion we know that migration is complete
	// check other stuff via queries if needed
	db := getDB(t, testEnvironment)
	var d []Data
	err = db.Select(&d, "select * from evm_chains;")
	require.NoError(t, err, "Error running SELECT")
	log.Info().Interface("Rows", d).Send()
}

func getDB(t *testing.T, testEnvironment *environment.Environment) *ctfClient.PostgresConnector {
	spl := strings.Split(testEnvironment.URLs["chainlink_db"][1], ":")
	port := spl[len(spl)-1]
	db, err := ctfClient.NewPostgresConnector(&ctfClient.PostgresConfig{
		Host:     "localhost",
		Port:     port,
		User:     "postgres",
		Password: "postgres",
		DBName:   "chainlink",
	})
	require.NoError(t, err, "Error connecting to postgres")
	return db
}
