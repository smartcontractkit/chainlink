package migrations

import (
	"regexp"

	"github.com/smartcontractkit/chainlink/core/store/migrations/migration0"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1559081901"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1559767166"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1560433987"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1560791143"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1560881846"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1560881855"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1560886530"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1560924400"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1564007745"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1565139192"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1565210496"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1565291711"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1565877314"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1566498796"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1566915476"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1567029116"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1568280052"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1568390387"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1568833756"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1570087128"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1570675883"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1573667511"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1573812490"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1574659987"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1575036327"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1576022702"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1579700934"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1580904019"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1581240419"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1584377646"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1585908150"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1585918589"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1586163842"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1586342453"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1586369235"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1586871710"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1586939705"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1586949323"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1586956053"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1587027516"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1587580235"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1587591248"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1587975059"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1588088353"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1588293486"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1588757164"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1588853064"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1589206996"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1589462363"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1589470036"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1590226486"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1591141873"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1591603775"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1592355365"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1594306515"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1594393769"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1594642891"
	"github.com/smartcontractkit/chainlink/core/store/migrations/migration1596021087"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	gormigrate "gopkg.in/gormigrate.v1"
)

var migrations []*gormigrate.Migration

func init() {
	migrations = []*gormigrate.Migration{
		{
			ID:      "0",
			Migrate: migration0.Migrate,
		},
		{
			ID:      "1559081901",
			Migrate: migration1559081901.Migrate,
		},
		{
			ID:      "1559767166",
			Migrate: migration1559767166.Migrate,
		},
		{
			ID:      "1560433987",
			Migrate: migration1560433987.Migrate,
		},
		{
			ID:      "1560791143",
			Migrate: migration1560791143.Migrate,
		},
		{
			ID:      "1560881846",
			Migrate: migration1560881846.Migrate,
		},
		{
			ID:      "1560886530",
			Migrate: migration1560886530.Migrate,
		},
		{
			ID:      "1560924400",
			Migrate: migration1560924400.Migrate,
		},
		{
			ID:      "1560881855",
			Migrate: migration1560881855.Migrate,
		},
		{
			ID:      "1565139192",
			Migrate: migration1565139192.Migrate,
		},
		{
			ID:      "1564007745",
			Migrate: migration1564007745.Migrate,
		},
		{
			ID:      "1565210496",
			Migrate: migration1565210496.Migrate,
		},
		{
			ID:      "1566498796",
			Migrate: migration1566498796.Migrate,
		},
		{
			ID:      "1565877314",
			Migrate: migration1565877314.Migrate,
		},
		{
			ID:      "1566915476",
			Migrate: migration1566915476.Migrate,
		},
		{
			ID:      "1567029116",
			Migrate: migration1567029116.Migrate,
		},
		{
			ID:      "1568280052",
			Migrate: migration1568280052.Migrate,
		},
		{
			ID:      "1565291711",
			Migrate: migration1565291711.Migrate,
		},
		{
			ID:      "1568390387",
			Migrate: migration1568390387.Migrate,
		},
		{
			ID:      "1568833756",
			Migrate: migration1568833756.Migrate,
		},
		{
			ID:      "1570087128",
			Migrate: migration1570087128.Migrate,
		},
		{
			ID:      "1570675883",
			Migrate: migration1570675883.Migrate,
		},
		{
			ID:      "1573667511",
			Migrate: migration1573667511.Migrate,
		},
		{
			ID:      "1573812490",
			Migrate: migration1573812490.Migrate,
		},
		{
			ID:      "1575036327",
			Migrate: migration1575036327.Migrate,
		},
		{
			ID:      "1574659987",
			Migrate: migration1574659987.Migrate,
		},
		{
			ID:      "1576022702",
			Migrate: migration1576022702.Migrate,
		},
		{
			ID:      "1579700934",
			Migrate: migration1579700934.Migrate,
		},
		{
			ID:      "1580904019",
			Migrate: migration1580904019.Migrate,
		},
		{
			ID:      "1581240419",
			Migrate: migration1581240419.Migrate,
		},
		{
			ID:      "1584377646",
			Migrate: migration1584377646.Migrate,
		},
		{
			ID:      "1585908150",
			Migrate: migration1585908150.Migrate,
		},
		{
			ID:      "1585918589",
			Migrate: migration1585918589.Migrate,
		},
		{
			ID:      "1586163842",
			Migrate: migration1586163842.Migrate,
		},
		{
			ID:      "1586342453",
			Migrate: migration1586342453.Migrate,
		}, {
			ID:      "1586369235",
			Migrate: migration1586369235.Migrate,
		},
		{
			ID:      "1586939705",
			Migrate: migration1586939705.Migrate,
		},
		{
			ID:      "1587027516",
			Migrate: migration1587027516.Migrate,
		},
		{
			ID:      "1587580235",
			Migrate: migration1587580235.Migrate,
		},
		{
			ID:      "1587591248",
			Migrate: migration1587591248.Migrate,
		},
		{
			ID:      "1587975059",
			Migrate: migration1587975059.Migrate,
		},
		{
			ID:      "1586956053",
			Migrate: migration1586956053.Migrate,
		},
		{
			ID:      "1588293486",
			Migrate: migration1588293486.Migrate,
		},
		{
			ID:      "1586949323",
			Migrate: migration1586949323.Migrate,
		},
		{
			ID:      "1588088353",
			Migrate: migration1588088353.Migrate,
		},
		{
			ID:      "1588757164",
			Migrate: migration1588757164.Migrate,
		},
		{
			ID:      "1588853064",
			Migrate: migration1588853064.Migrate,
		},
		{
			ID:      "1589470036",
			Migrate: migration1589470036.Migrate,
		},
		{
			ID:      "1586871710",
			Migrate: migration1586871710.Migrate,
		},
		{
			ID:      "1590226486",
			Migrate: migration1590226486.Migrate,
		},
		{
			ID:      "1591141873",
			Migrate: migration1591141873.Migrate,
		},
		{
			ID:      "1589206996",
			Migrate: migration1589206996.Migrate,
		},
		{
			ID:      "1589462363",
			Migrate: migration1589462363.Migrate,
		},
		{
			ID:      "1591603775",
			Migrate: migration1591603775.Migrate,
		},
		{
			ID:      "1592355365",
			Migrate: migration1592355365.Migrate,
		},
		{
			ID:      "1594393769",
			Migrate: migration1594393769.Migrate,
		},
		{
			ID:      "1594642891",
			Migrate: migration1594642891.Migrate,
		},
		{
			ID:      "1594306515",
			Migrate: migration1594306515.Migrate,
		},
		{
			ID:      "1596021087",
			Migrate: migration1596021087.Migrate,
		},
	}
}

// GORMMigrate calls through to gorm's native migrate function with minimal
// extra logic
// Useful if the migrations table doesn't exist yet but we don't care
func GORMMigrate(db *gorm.DB) error {
	options := *gormigrate.DefaultOptions
	options.UseTransaction = true

	m := gormigrate.New(db, &options, migrations)
	return m.Migrate()
}

// Migrate iterates through available migrations, running and tracking
// migrations that have not been run.
func Migrate(db *gorm.DB) error {
	return MigrateTo(db, "")
}

// MigrateTo runs all migrations up to and including the specified migration ID
func MigrateTo(db *gorm.DB, migrationID string) error {
	options := *gormigrate.DefaultOptions
	options.UseTransaction = true

	m := gormigrate.New(db, &options, migrations)

	var count int
	err := db.Table(options.TableName).Count(&count).Error
	if err != nil && !noSuchTableRegex.MatchString(err.Error()) {
		return errors.Wrap(err, "error determining migration count")
	}

	if count > len(migrations) {
		return errors.New("database is newer than current chainlink version")
	}

	if migrationID == "" {
		migrationID = migrations[len(migrations)-1].ID
	}

	err = m.MigrateTo(migrationID)
	if err != nil {
		return errors.Wrap(err, "error running migrations")
	}
	return nil
}

var (
	noSuchTableRegex = regexp.MustCompile(`^(no such table|pq: relation ".*?" does not exist)`)
)
