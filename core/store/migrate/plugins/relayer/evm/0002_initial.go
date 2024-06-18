package evm

import (
	"bytes"
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"io"

	"github.com/pressly/goose/v3"
)

//go:embed initUp.tmpl.sql
var upTmpl string

func resolveUp(out io.Writer, val Cfg) error {
	if upTmpl == "" {
		return fmt.Errorf("upTmpl is empty")
	}
	return resolve(out, upTmpl, val)
}

//go:embed  initDown.tmpl.sql
var downTmpl string

func resolveDown(out io.Writer, val Cfg) error {
	return resolve(out, downTmpl, val)
}

// Register0002 registers the migration with goose
func Register0002(val Cfg) error {
	upSQL := &bytes.Buffer{}
	err := resolveUp(upSQL, val)
	if err != nil {
		return fmt.Errorf("failed to resolve up sql: %w", err)
	}
	upFunc := func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, upSQL.String())
		return err
	}

	downSQL := &bytes.Buffer{}
	err = resolveDown(downSQL, val)
	if err != nil {
		return fmt.Errorf("failed to resolve down sql: %w", err)
	}
	downFunc := func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, downSQL.String())
		return err
	}
	goose.AddMigrationContext(upFunc, downFunc)
	return nil
}
