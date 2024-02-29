package pg

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/store/dialects"
)

func Test_sprintQ(t *testing.T) {
	for _, tt := range []struct {
		name  string
		query string
		args  []interface{}
		exp   string
	}{
		{"none",
			"SELECT * FROM table;",
			nil,
			"SELECT * FROM table;"},
		{"one",
			"SELECT $1 FROM table;",
			[]interface{}{"foo"},
			"SELECT 'foo' FROM table;"},
		{"two",
			"SELECT $1 FROM table WHERE bar = $2;",
			[]interface{}{"foo", 1},
			"SELECT 'foo' FROM table WHERE bar = 1;"},
		{"limit",
			"SELECT $1 FROM table LIMIT $2;",
			[]interface{}{"foo", Limit(10)},
			"SELECT 'foo' FROM table LIMIT 10;"},
		{"limit-all",
			"SELECT $1 FROM table LIMIT $2;",
			[]interface{}{"foo", Limit(-1)},
			"SELECT 'foo' FROM table LIMIT NULL;"},
		{"bytea",
			"SELECT $1 FROM table WHERE b = $2;",
			[]interface{}{"foo", []byte{0x0a}},
			"SELECT 'foo' FROM table WHERE b = '\\x0a';"},
		{"bytea[]",
			"SELECT $1 FROM table WHERE b = $2;",
			[]interface{}{"foo", pq.ByteaArray([][]byte{{0xa}, {0xb}})},
			"SELECT 'foo' FROM table WHERE b = ARRAY['\\x0a','\\x0b'];"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := sprintQ(tt.query, tt.args)
			t.Log(tt.query, tt.args)
			t.Log(got)
			require.Equal(t, tt.exp, got)
		})
	}
}

func Test_ExecQWithRowsAffected(t *testing.T) {
	db, err := sqlx.Open(string(dialects.TransactionWrappedPostgres), uuid.New().String())
	require.NoError(t, err)
	q := NewQ(db, logger.NullLogger, NewQConfig(false))

	require.NoError(t, q.ExecQ("CREATE TABLE testtable (a TEXT, b TEXT)"))

	rows, err := q.ExecQWithRowsAffected("INSERT INTO testtable (a, b) VALUES ($1, $2)", "foo", "bar")
	require.NoError(t, err)
	assert.Equal(t, int64(1), rows)

	rows, err = q.ExecQWithRowsAffected("INSERT INTO testtable (a, b) VALUES ($1, $1), ($2, $2), ($1, $2)", "foo", "bar")
	require.NoError(t, err)
	assert.Equal(t, int64(3), rows)

	rows, err = q.ExecQWithRowsAffected("delete from testtable")
	require.NoError(t, err)
	assert.Equal(t, int64(4), rows)

	rows, err = q.ExecQWithRowsAffected("delete from testtable")
	require.NoError(t, err)
	assert.Equal(t, int64(0), rows)
}
