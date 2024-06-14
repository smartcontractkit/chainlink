package pg

import (
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
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
			"SELECT foo FROM table;"},
		{"two",
			"SELECT $1 FROM table WHERE bar = $2;",
			[]interface{}{"foo", 1},
			"SELECT foo FROM table WHERE bar = 1;"},
		{"limit",
			"SELECT $1 FROM table LIMIT $2;",
			[]interface{}{"foo", Limit(10)},
			"SELECT foo FROM table LIMIT 10;"},
		{"limit-all",
			"SELECT $1 FROM table LIMIT $2;",
			[]interface{}{"foo", Limit(-1)},
			"SELECT foo FROM table LIMIT NULL;"},
		{"bytea",
			"SELECT $1 FROM table WHERE b = $2;",
			[]interface{}{"foo", []byte{0x0a}},
			"SELECT foo FROM table WHERE b = '\\x0a';"},
		{"bytea[]",
			"SELECT $1 FROM table WHERE b = $2;",
			[]interface{}{"foo", pq.ByteaArray([][]byte{{0xa}, {0xb}})},
			"SELECT foo FROM table WHERE b = ('\\x0a','\\x0b');"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := sprintQ(tt.query, tt.args)
			t.Log(tt.query, tt.args)
			t.Log(got)
			require.Equal(t, tt.exp, got)
		})
	}
}
