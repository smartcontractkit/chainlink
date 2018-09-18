## Migrations

Every migration is encapsulated to prevent type collisions now and in the future.

Each migration contains two packages, the root package which includes the definitions
of the types being moved to, and the `old` package which contains the previous
definition of each type.

Fields that remain unchanged during a migration can be easily handled by marking
their type as `migration0.Unchanged` in both the old and new definition of each type.

When introducing a new type, be sure to include its JSON Marshaling and Unmarshaling
behavior if it is customized.

### Example

```golang
package old // previous type structure

type RunResult struct {
	JobRunID     migration0.Unchanged `json:"jobRunId" storm:"id"`
	Amount       *big.Int             `json:"amount,omitempty"` // old type
}
```

```golang
package migration1

type RunResult struct {
	JobRunID     migration0.Unchanged `json:"jobRunId" storm:"id"`
	Amount       *assets.Link         `json:"amount,omitempty"` // new type
}
```

One can now deserialize `old.RunResult`, convert to `migration1.RunResult`, and then
persist the new instance to bolt.

```golang
package migration1

type RunResult struct {
	JobRunID     migration0.Unchanged `json:"jobRunId" storm:"id"`
	Amount       *assets.Link         `json:"amount,omitempty"` // new type
}

func (m Migration) Migrate(orm *orm.ORM) error {
	var jrs []old.JobRun
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

* Be sure to add your new migration to the list of migrations in `store/migrations/migrate.go`:
