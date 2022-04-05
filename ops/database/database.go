package database

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/lib/pq" // only place where database configs are called
	"github.com/smartcontractkit/chainlink-relay/ops/utils"
)

// Database contains the configuration for the postgres db
type Database struct {
	User, Host, Sslmode string
	Timeout, Port       int
	db                  *sql.DB
}

// Ready performs a health check on the db
func (d *Database) Ready() error {
	msg := utils.LogStatus("Waiting for health check on Postgres DB")

	// prep db connection
	conninfo := fmt.Sprintf("user=%s host=%s port=%d sslmode=%s", d.User, d.Host, d.Port, d.Sslmode)
	db, err := sql.Open("postgres", conninfo)
	if err != nil {
		return msg.Check(err)
	}
	d.db = db

	// checking if DB is ready for connection
	time.Sleep(2 * time.Second) // removing this breaks running `up` multiple times
	for i := 0; i < d.Timeout; i++ {
		err = d.db.Ping()
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second) // wait 1 second before retrying
	}
	return msg.Check(err)
}

// Create creates a database with a spcific name
func (d *Database) Create(name string) error {
	msg := utils.LogStatus(fmt.Sprintf("Creating database: %s", name))
	_, err := d.db.Exec("create database " + name)
	if err != nil && strings.Contains(err.Error(), fmt.Sprintf("pq: database \"%s\" already exists", name)) {
		msg.Exists()
		err = nil
	}
	return msg.Check(err)
}
