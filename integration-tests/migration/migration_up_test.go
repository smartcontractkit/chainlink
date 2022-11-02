package migration_test

import (
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-env/environment"
	ctfClient "github.com/smartcontractkit/chainlink-testing-framework/client"
	"github.com/smartcontractkit/chainlink-testing-framework/testsetups"
	"github.com/smartcontractkit/chainlink-testing-framework/utils"
	"github.com/smartcontractkit/chainlink/integration-tests/actions"
)

type Data struct {
	ID        int       `db:"id"`
	Cfg       []byte    `db:"cfg"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Enabled   bool      `db:"enabled"`
}

func getDB(e *environment.Environment) *ctfClient.PostgresConnector {
	spl := strings.Split(e.URLs["chainlink_db"][1], ":")
	port := spl[len(spl)-1]
	db, err := ctfClient.NewPostgresConnector(&ctfClient.PostgresConfig{
		Host:     "localhost",
		Port:     port,
		User:     "postgres",
		Password: "postgres",
		DBName:   "chainlink",
	})
	Expect(err).ShouldNot(HaveOccurred())
	return db
}

// Migration template boiler, for now it's semi-automatic, integrating into CI to make it fully automatic
var _ = Describe("Migration up test suite @db-migration", func() {
	var (
		err error
		e   *environment.Environment
	)

	AfterEach(func() {
		By("Tearing down the environment")
		err = actions.TeardownSuite(e, utils.ProjectRoot, nil, nil, nil)
		Expect(err).ShouldNot(HaveOccurred(), "Environment teardown shouldn't fail")
	})

	Describe("Migration up succeeds @db-migration-up", func() {
		It("Migrated successfully", func() {
			e, err = testsetups.DBMigration(&testsetups.DBMigrationSpec{
				FromSpec: testsetups.FromVersionSpec{
					Image: "public.ecr.aws/chainlink/chainlink",
					Tag:   "1.6.0-nonroot",
				},
				ToSpec: testsetups.ToVersionSpec{
					Image: "public.ecr.aws/chainlink/chainlink",
					Tag:   "1.7.1-nonroot",
				},
			})
			Expect(err).ShouldNot(HaveOccurred())
			// if test haven't failed after that assertion we know that migration is complete
			// check other stuff via queries if needed
			db := getDB(e)
			var d []Data
			err = db.Select(&d, "select * from evm_chains;")
			Expect(err).ShouldNot(HaveOccurred())
			log.Info().Interface("Rows", d).Send()
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
