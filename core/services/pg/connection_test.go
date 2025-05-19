package pg

import (
	"testing"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ Getter = &mockGetter{}

type mockGetter struct {
	version int
	err     error
}

func (m *mockGetter) Get(dest interface{}, query string, args ...interface{}) error {
	if m.err != nil {
		return m.err
	}
	*(dest.(*int)) = m.version
	return nil
}

func Test_checkVersion(t *testing.T) {
	if time.Now().Year() > 2027 {
		t.Fatal("Postgres version numbers only registered until 2028, please update the postgres version check using: https://www.postgresql.org/support/versioning/ then fix this test")
	}
	t.Run("when the version is too low", func(t *testing.T) {
		m := &mockGetter{version: 100000}
		err := checkVersion(m, 110000)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "The minimum required Postgres server version is 11, you are running: 10")
	})
	t.Run("when the version is at minimum", func(t *testing.T) {
		m := &mockGetter{version: 110000}
		err := checkVersion(m, 110000)
		require.NoError(t, err)
	})
	t.Run("when the version is above minimum", func(t *testing.T) {
		m := &mockGetter{version: 110001}
		err := checkVersion(m, 110000)
		require.NoError(t, err)
		m = &mockGetter{version: 120000}
		err = checkVersion(m, 110001)
		require.NoError(t, err)
	})
	t.Run("ignores wildly small versions, 0 etc", func(t *testing.T) {
		m := &mockGetter{version: 9000}
		err := checkVersion(m, 110001)
		require.NoError(t, err)
	})
	t.Run("ignores errors", func(t *testing.T) {
		m := &mockGetter{err: errors.New("some error")}
		err := checkVersion(m, 110001)
		require.NoError(t, err)
	})
}
