package orm

import "github.com/smartcontractkit/chainlink/core/gracefulpanic"

func (o *ORM) LockingStrategyHelperSimulateDisconnect() (error, error) {
	err1 := o.lockingStrategy.(*PostgresLockingStrategy).conn.Close()
	err2 := o.lockingStrategy.(*PostgresLockingStrategy).db.Close()
	o.lockingStrategy.(*PostgresLockingStrategy).conn = nil
	o.lockingStrategy.(*PostgresLockingStrategy).db = nil
	return err1, err2
}

func (o *ORM) ShutdownSignal() gracefulpanic.Signal {
	return o.shutdownSignal
}
