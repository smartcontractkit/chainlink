package pgtest

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/DATA-DOG/go-txdb"
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
	t.Helper()

	sqlDB := NewSqlDB(t)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	require.NoError(t, err)

	return gormDB
}

func NewSqlDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("txdb", t.Name())
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	return db
}
