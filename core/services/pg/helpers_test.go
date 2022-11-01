package pg

import "github.com/smartcontractkit/sqlx"

func SetConn(lock interface{}, conn *sqlx.Conn) {
	switch v := lock.(type) {
	case *leaseLock:
		v.conn = conn
	default:
		panic("cannot set conn on unknown type")
	}
}

func GetConn(lock interface{}) *sqlx.Conn {
	switch v := lock.(type) {
	case *leaseLock:
		return v.conn
	default:
		panic("cannot get conn on unknown type")
	}
}
