package pgtest

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("You must provide a DATABASE_URL environment variable")
	}

	txdb.Register("txdb", "postgres", dbURL)
}

func NewGormDB(t *testing.T) *gorm.DB {
	sqlDB := NewSqlDB(t)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
		DSN:  uuid.NewV4().String(),
	}), &gorm.Config{})
	require.NoError(t, err)

	return gormDB
}

func NewSqlDB(t *testing.T) *sql.DB {
	db, err := sql.Open("txdb", uuid.NewV4().String())
	require.NoError(t, err)
	t.Cleanup(func() { assert.NoError(t, db.Close()) })

	// There is a bug to do with context cancellation somewhere in txdb, sql or
	// gorm. If you try to use the DB "too quickly" using a .WithContext then
	// cancel the context, the transaction state gets poisoned or lost somehow
	// and subsequent queries fail with "sql: transaction has already been
	// committed or rolled back" (although postgres does not log any errors).
	// Calling runtime.Gosched() sometimes helps. Calling SELECT 1 here seems
	// to reliably fix it. Created an issue to track here:
	// https://github.com/DATA-DOG/go-txdb/issues/43
	// runtime.Gosched()
	_, err = db.Exec(`SELECT 1`)
	require.NoError(t, err)

	return db
}
