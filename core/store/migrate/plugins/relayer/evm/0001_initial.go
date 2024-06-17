package evm

import (
	"bytes"
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"io"
	"text/template"

	"github.com/pressly/goose/v3"
)

type Cfg struct {
	Schema  string
	ChainID int
}

func RegisterSchemaMigration(val Cfg) error {
	return Register(val)
}

// go:embed initUp.tmpl.sql
var upTmpl string

func resolveUp(out io.Writer, val Cfg) error {
	if upTmpl == "" {
		return fmt.Errorf("upTmpl is empty")
	}
	return resolve(out, upTmpl, val)
}

// go:embed  initDown.tmpl.sql
var downTmpl string

func resolveDown(out io.Writer, val Cfg) error {
	return resolve(out, downTmpl, val)
}

func resolve(out io.Writer, in string, val Cfg) error {
	id := fmt.Sprintf("init_%s_%d", val.Schema, val.ChainID)
	tmpl, err := template.New(id).Parse(in)
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", in, err)
	}
	err = tmpl.Execute(out, val)
	if err != nil {
		return fmt.Errorf("failed to execute template %s: %w", in, err)
	}
	return nil
}

// Register registers the migration with goose
func Register(val Cfg) error {
	upSQL := &bytes.Buffer{}
	err := resolveUp(upSQL, val)
	if err != nil {
		return fmt.Errorf("failed to resolve up sql: %w", err)
	}
	upFunc := func(ctx context.Context, tx *sql.Tx) error {
		fmt.Printf("Executing up sql: %s\n", upSQL.String())
		panic(fmt.Sprintf("Executing up sql: %s\n", upSQL.String()))
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

func generateExec(execString string) func(ctx context.Context, tx *sql.Tx) error {
	return func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, execString)
		return err
	}
}
