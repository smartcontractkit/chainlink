package pgtest

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"flag"
	"fmt"
	"io"
	"net/url"
	"strings"
	"sync"
	"testing"

	"github.com/jmoiron/sqlx"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/v2/core/config/env"
	"github.com/smartcontractkit/chainlink/v2/core/store/dialects"
)

// txdb is a simplified version of https://github.com/DATA-DOG/go-txdb
//
// The original lib has various problems and is hard to understand because it
// tries to be more general. The version in this file is more tightly focused
// to our needs and should be easier to reason about and less likely to have
// subtle bugs/races.
//
// It doesn't currently support savepoints but could be made to if necessary.
//
// Transaction BEGIN/ROLLBACK effectively becomes a no-op, this should have no
// negative impact on normal test operation.
//
// If you MUST test BEGIN/ROLLBACK behaviour, you will have to configure your
// store to use the raw DialectPostgres dialect and setup a one-use database.
// See heavyweight.FullTestDB() as a convenience function to help you do this,
// but please use sparingly because as it's name implies, it is expensive.
func init() {
	testing.Init()
	if !flag.Parsed() {
		flag.Parse()
	}
	if testing.Short() {
		// -short tests don't need a DB
		return
	}
	dbURL := string(env.DatabaseURL.Get())
	if dbURL == "" {
		panic("you must provide a CL_DATABASE_URL environment variable")
	}

	parsed, err := url.Parse(dbURL)
	if err != nil {
		panic(err)
	}
	if parsed.Path == "" {
		msg := fmt.Sprintf("invalid %[1]s: `%[2]s`. You must set %[1]s env var to point to your test database. Note that the test database MUST end in `_test` to differentiate from a possible production DB. HINT: Try %[1]s=postgresql://postgres@localhost:5432/chainlink_test?sslmode=disable", env.DatabaseURL, parsed.String())
		panic(msg)
	}
	if !strings.HasSuffix(parsed.Path, "_test") {
		msg := fmt.Sprintf("cannot run tests against database named `%s`. Note that the test database MUST end in `_test` to differentiate from a possible production DB. HINT: Try %s=postgresql://postgres@localhost:5432/chainlink_test?sslmode=disable", parsed.Path[1:], env.DatabaseURL)
		panic(msg)
	}
	name := string(dialects.TransactionWrappedPostgres)
	sql.Register(name, &txDriver{
		dbURL: dbURL,
		conns: make(map[string]*conn),
	})
	sqlx.BindDriver(name, sqlx.DOLLAR)
}

var _ driver.Conn = &conn{}

var _ driver.Validator = &conn{}
var _ driver.SessionResetter = &conn{}

// txDriver is an sql driver which runs on a single transaction.
// When `Close` is called, transaction is rolled back.
type txDriver struct {
	sync.Mutex
	db    *sql.DB
	conns map[string]*conn

	dbURL string
}

func (d *txDriver) Open(dsn string) (driver.Conn, error) {
	d.Lock()
	defer d.Unlock()
	// Open real db connection if its the first call
	if d.db == nil {
		db, err := sql.Open(string(dialects.Postgres), d.dbURL)
		if err != nil {
			return nil, err
		}
		d.db = db
	}
	c, exists := d.conns[dsn]
	if !exists || !c.tryOpen() {
		tx, err := d.db.Begin()
		if err != nil {
			return nil, err
		}
		c = &conn{tx: tx, opened: 1, dsn: dsn}
		c.removeSelf = func() error {
			return d.deleteConn(c)
		}
		d.conns[dsn] = c
	}
	return c, nil
}

// deleteConn is called by a connection when it is closed via the `close` method.
// It also auto-closes the DB when the last checked out connection is closed.
func (d *txDriver) deleteConn(c *conn) error {
	// must lock here to avoid racing with Open
	d.Lock()
	defer d.Unlock()

	if d.conns[c.dsn] != c {
		return nil // already been replaced
	}
	delete(d.conns, c.dsn)
	if len(d.conns) == 0 && d.db != nil {
		if err := d.db.Close(); err != nil {
			return err
		}
		d.db = nil
	}
	return nil
}

type conn struct {
	sync.Mutex
	dsn        string
	tx         *sql.Tx // tx may be shared by many conns, definitive one lives in the map keyed by DSN on the txDriver. Do not modify from conn
	closed     bool
	opened     int
	removeSelf func() error
}

func (c *conn) Begin() (driver.Tx, error) {
	c.Lock()
	defer c.Unlock()
	if c.closed {
		panic("conn is closed")
	}
	// Begin is a noop because the transaction was already opened
	return tx{c.tx}, nil
}

// Implement the "ConnBeginTx" interface
func (c *conn) BeginTx(_ context.Context, opts driver.TxOptions) (driver.Tx, error) {
	// Context is ignored, because single transaction is shared by all callers, thus caller should not be able to
	// control it with local context
	return c.Begin()
}

// Prepare returns a prepared statement, bound to this connection.
func (c *conn) Prepare(query string) (driver.Stmt, error) {
	return c.PrepareContext(context.Background(), query)
}

// Implement the "ConnPrepareContext" interface
func (c *conn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	c.Lock()
	defer c.Unlock()
	if c.closed {
		panic("conn is closed")
	}

	// TODO: Fix context handling
	// FIXME: It is not safe to give the passed in context to the tx directly
	// because the tx is shared by many conns and cancelling the context will
	// destroy the tx which can affect other conns
	st, err := c.tx.PrepareContext(context.Background(), query)
	if err != nil {
		return nil, err
	}
	return &stmt{st, c}, nil
}

// IsValid is called prior to placing the connection into the
// connection pool by database/sql. The connection will be discarded if false is returned.
func (c *conn) IsValid() bool {
	c.Lock()
	defer c.Unlock()
	return !c.closed
}

func (c *conn) ResetSession(ctx context.Context) error {
	// Ensure bad connections are reported: From database/sql/driver:
	// If a connection is never returned to the connection pool but immediately reused, then
	// ResetSession is called prior to reuse but IsValid is not called.
	c.Lock()
	defer c.Unlock()
	if c.closed {
		return driver.ErrBadConn
	}

	return nil
}

// pgx returns nil
func (c *conn) CheckNamedValue(nv *driver.NamedValue) error {
	return nil
}

// Implement the "QueryerContext" interface
func (c *conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	c.Lock()
	defer c.Unlock()
	if c.closed {
		panic("conn is closed")
	}

	// TODO: Fix context handling
	rs, err := c.tx.QueryContext(context.Background(), query, mapNamedArgs(args)...)
	if err != nil {
		return nil, err
	}
	defer rs.Close()

	return buildRows(rs)
}

// Implement the "ExecerContext" interface
func (c *conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	c.Lock()
	defer c.Unlock()
	if c.closed {
		panic("conn is closed")
	}
	// TODO: Fix context handling
	return c.tx.ExecContext(context.Background(), query, mapNamedArgs(args)...)
}

// tryOpen attempts to increment the open count, but returns false if closed.
func (c *conn) tryOpen() bool {
	c.Lock()
	defer c.Unlock()
	if c.closed {
		return false
	}
	c.opened++
	return true
}

// Close invalidates and potentially stops any current
// prepared statements and transactions, marking this
// connection as no longer in use.
//
// Because the sql package maintains a free pool of
// connections and only calls Close when there's a surplus of
// idle connections, it shouldn't be necessary for drivers to
// do their own connection caching.
//
// Drivers must ensure all network calls made by Close
// do not block indefinitely (e.g. apply a timeout).
func (c *conn) Close() (err error) {
	if !c.close() {
		return
	}
	// Wait to remove self to avoid nesting locks.
	if err := c.removeSelf(); err != nil {
		panic(err)
	}
	return
}

func (c *conn) close() bool {
	c.Lock()
	defer c.Unlock()
	if c.closed {
		// Double close, should be a safe to make this a noop
		// PGX allows double close
		// See: https://github.com/jackc/pgx/blob/a457da8bffa4f90ad672fa093ee87f20cf06687b/conn.go#L249
		return false
	}

	c.opened--
	if c.opened > 0 {
		return false
	}
	if c.tx != nil {
		if err := c.tx.Rollback(); err != nil {
			panic(err)
		}
		c.tx = nil
	}
	c.closed = true
	return true
}

type tx struct {
	tx *sql.Tx
}

func (tx tx) Commit() error {
	// Commit is a noop because the transaction will be rolled back at the end
	return nil
}

func (tx tx) Rollback() error {
	// Rollback is a noop because the transaction will be rolled back at the end
	return nil
}

type stmt struct {
	st   *sql.Stmt
	conn *conn
}

func (s stmt) Exec(args []driver.Value) (driver.Result, error) {
	s.conn.Lock()
	defer s.conn.Unlock()
	if s.conn.closed {
		panic("conn is closed")
	}
	return s.st.Exec(mapArgs(args)...)
}

// Implement the "StmtExecContext" interface
func (s *stmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	s.conn.Lock()
	defer s.conn.Unlock()
	if s.conn.closed {
		panic("conn is closed")
	}
	// TODO: Fix context handling
	return s.st.ExecContext(context.Background(), mapNamedArgs(args)...)
}

func mapArgs(args []driver.Value) (res []interface{}) {
	res = make([]interface{}, len(args))
	for i := range args {
		res[i] = args[i]
	}
	return
}

func (s stmt) NumInput() int {
	return -1
}

func (s stmt) Query(args []driver.Value) (driver.Rows, error) {
	s.conn.Lock()
	defer s.conn.Unlock()
	if s.conn.closed {
		panic("conn is closed")
	}
	rows, err := s.st.Query(mapArgs(args)...)
	defer func() {
		err = multierr.Combine(err, rows.Close())
	}()
	if err != nil {
		return nil, err
	}
	return buildRows(rows)
}

// Implement the "StmtQueryContext" interface
func (s *stmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	s.conn.Lock()
	defer s.conn.Unlock()
	if s.conn.closed {
		panic("conn is closed")
	}
	// TODO: Fix context handling
	rows, err := s.st.QueryContext(context.Background(), mapNamedArgs(args)...)
	if err != nil {
		return nil, err
	}
	return buildRows(rows)
}

func (s stmt) Close() error {
	s.conn.Lock()
	defer s.conn.Unlock()
	return s.st.Close()
}

func buildRows(r *sql.Rows) (driver.Rows, error) {
	set := &rowSets{}
	rs := &rows{}
	if err := rs.read(r); err != nil {
		return set, err
	}
	set.sets = append(set.sets, rs)
	for r.NextResultSet() {
		rss := &rows{}
		if err := rss.read(r); err != nil {
			return set, err
		}
		set.sets = append(set.sets, rss)
	}
	return set, nil
}

// Implement the "RowsNextResultSet" interface
func (rs *rowSets) HasNextResultSet() bool {
	return rs.pos+1 < len(rs.sets)
}

// Implement the "RowsNextResultSet" interface
func (rs *rowSets) NextResultSet() error {
	if !rs.HasNextResultSet() {
		return io.EOF
	}

	rs.pos++
	return nil
}

type rows struct {
	rows     [][]driver.Value
	pos      int
	cols     []string
	colTypes []*sql.ColumnType
}

func (r *rows) Columns() []string {
	return r.cols
}

func (r *rows) ColumnTypeDatabaseTypeName(index int) string {
	return r.colTypes[index].DatabaseTypeName()
}

func (r *rows) Next(dest []driver.Value) error {
	r.pos++
	if r.pos > len(r.rows) {
		return io.EOF
	}

	for i, val := range r.rows[r.pos-1] {
		dest[i] = *(val.(*interface{}))
	}

	return nil
}

func (r *rows) Close() error {
	return nil
}

func (r *rows) read(rs *sql.Rows) error {
	var err error
	r.cols, err = rs.Columns()
	if err != nil {
		return err
	}

	r.colTypes, err = rs.ColumnTypes()
	if err != nil {
		return err
	}

	for rs.Next() {
		values := make([]interface{}, len(r.cols))
		for i := range values {
			values[i] = new(interface{})
		}
		if err := rs.Scan(values...); err != nil {
			return err
		}
		row := make([]driver.Value, len(r.cols))
		for i, v := range values {
			row[i] = driver.Value(v)
		}
		r.rows = append(r.rows, row)
	}
	return rs.Err()
}

type rowSets struct {
	sets []*rows
	pos  int
}

func (rs *rowSets) Columns() []string {
	return rs.sets[rs.pos].cols
}

func (rs *rowSets) ColumnTypeDatabaseTypeName(index int) string {
	return rs.sets[rs.pos].ColumnTypeDatabaseTypeName(index)
}

func (rs *rowSets) Close() error {
	return nil
}

// advances to next row
func (rs *rowSets) Next(dest []driver.Value) error {
	return rs.sets[rs.pos].Next(dest)
}

func mapNamedArgs(args []driver.NamedValue) (res []interface{}) {
	res = make([]interface{}, len(args))
	for i := range args {
		name := args[i].Name
		if name != "" {
			res[i] = sql.Named(name, args[i].Value)
		} else {
			res[i] = args[i].Value
		}
	}
	return
}
