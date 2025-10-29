package pg

import "github.com/jmoiron/sqlx"

func SetConn(lock any, conn *sqlx.Conn) {
	switch v := lock.(type) {
	case *leaseLock:
		v.conn = conn
	default:
		panic("cannot set conn on unknown type")
	}
}

func GetConn(lock any) *sqlx.Conn {
	switch v := lock.(type) {
	case *leaseLock:
		return v.conn
	default:
		panic("cannot get conn on unknown type")
	}
}
