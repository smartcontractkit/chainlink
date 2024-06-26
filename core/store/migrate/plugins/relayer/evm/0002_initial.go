package evm

import (
	"bytes"
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	"github.com/pressly/goose/v3"
)

//go:embed forwardersUp.tmpl.sql
var forwardersUpTmpl string

//go:embed  forwardersDown.tmpl.sql
var forwardersDownTmpl string

//go:embed headsUp.tmpl.sql
var headsUpTmpl string

//go:embed headsDown.tmpl.sql
var headsDownTmpl string

//go:embed key_statesUp.tmpl.sql
var keyStatesUpTmpl string

//go:embed key_statesDown.tmpl.sql
var keyStatesDownTmpl string

//go:embed log_poller_blocksUp.tmpl.sql
var logPollerBlocksUpTmpl string

//go:embed log_poller_blocksDown.tmpl.sql
var logPollerBlocksDownTmpl string

type initialMigration struct {
	upTmpl   string
	downTmpl string
	version  int64
}

var (
	forwarderMigration = initialMigration{
		upTmpl:   forwardersUpTmpl,
		downTmpl: forwardersDownTmpl,
		version:  2,
	}

	headsMigration = initialMigration{
		upTmpl:   headsUpTmpl,
		downTmpl: headsDownTmpl,
		version:  3,
	}

	keyStatesMigration = initialMigration{
		upTmpl:   keyStatesUpTmpl,
		downTmpl: keyStatesDownTmpl,
		version:  4,
	}

	logPollerBlocksMigration = initialMigration{
		upTmpl:   logPollerBlocksUpTmpl,
		downTmpl: logPollerBlocksDownTmpl,
		version:  5,
	}

	initialMigrations = []initialMigration{forwarderMigration, headsMigration, keyStatesMigration, logPollerBlocksMigration}
)

func generateGoMigration(val Cfg, m initialMigration) (*goose.Migration, error) {
	upSQL := &bytes.Buffer{}
	err := resolve(upSQL, m.upTmpl, val)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve up sql: %w", err)
	}
	upFunc := func(ctx context.Context, tx *sql.Tx) error {
		_, terr := tx.ExecContext(ctx, upSQL.String())
		return terr
	}

	downSQL := &bytes.Buffer{}
	err = resolve(downSQL, m.downTmpl, val)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve down sql: %w", err)
	}
	downFunc := func(ctx context.Context, tx *sql.Tx) error {
		_, terr := tx.ExecContext(ctx, downSQL.String())
		return terr
	}
	up := &goose.GoFunc{RunTx: upFunc}
	down := &goose.GoFunc{RunTx: downFunc}
	return goose.NewGoMigration(m.version, up, down), nil
}

func generateInitialMigrations(val Cfg) ([]*goose.Migration, error) {
	migrations := []*goose.Migration{}
	for _, m := range initialMigrations {
		mig, err := generateGoMigration(val, m)
		if err != nil {
			return nil, fmt.Errorf("failed to generate migration: %w", err)
		}
		migrations = append(migrations, mig)
	}
	return migrations, nil
}
