## Migrations

Every migration is encapsulated in its own package to prevent type
collisions now and for the future.

It serves as a snapshot in time of the schema for bolt db.

As a result, future migrations can reference the types of previous migrations
to assist deserializing, transforming, and saving new types.

When introducing a new type, be sure to include its JSON Marshaling and Unmarshaling
behavior if it is customized.

### Example

```golang
package migration1 // previous existing migration

type RunResult struct {
	JobRunID     string      `json:"jobRunId" storm:"id"`
	Amount       *big.Int    `json:"amount,omitempty"` // old type
}
```

```golang
package migration2

type RunResult struct {
	JobRunID     string          `json:"jobRunId" storm:"id"`
	Amount       *assets.Link    `json:"amount,omitempty"` // new type
}
```

One can now deserialize `migration1.RunResult`, convert to `migration2.RunResult`, and then
persist the new instance to bolt. When updating a model, you must always refer to the most
recently used migration version of that model for correct extraction.

```golang
package migration2

...

func (m Migration) Migrate(orm *orm.ORM) error {
	var jrs []migration1.JobRun
	if err := orm.All(&jrs); err != nil {
		return fmt.Errorf("failed migration1: %v", err)
	}

	tx, err := orm.Begin(true)
	if err != nil {
		return fmt.Errorf("error starting transaction: %+v", err)
	}
	defer tx.Rollback()

	for _, jr := range jrs {
		jr2 := m.Convert(jr) // Convert to migration2.JobRun
		if err := tx.Save(&jr2); err != nil {
			return err
		}
	}

	return tx.Commit()
}
```

### Add to global registry

`store/migrations/migrate.go`

```golang
func init() {
	registerMigration(migration1.Migration{})
	registerMigration(migration2.Migration{})
	registerMigration(migration<newMigration>.Migration{})
}
```

### Helpers

Run `cldev migration` to generate a migration template, prepopulated with a current
timestamp.

* Be sure to add your new migration to the list of migrations in `migrate.go`:


```golang
package migration1536766540

import "github.com/smartcontractkit/chainlink/store/orm"

type Migration struct{}

func (m Migration) Timestamp() string {
	return "1536766540"
}

func (m Migration) Migrate(orm *orm.ORM) error {
}
```
