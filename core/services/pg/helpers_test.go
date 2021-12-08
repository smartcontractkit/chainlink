package pg

import "github.com/smartcontractkit/sqlx"

func GetConn(ll LeaseLock) *sqlx.Conn {
	return ll.(*leaseLock).conn
}
