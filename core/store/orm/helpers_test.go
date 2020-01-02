package orm

func (o *ORM) LockingStrategyHelperSimulateDisconnect() (error, error) {
	err1 := o.lockingStrategy.(*PostgresLockingStrategy).conn.Close()
	err2 := o.lockingStrategy.(*PostgresLockingStrategy).db.Close()
	o.lockingStrategy.(*PostgresLockingStrategy).conn = nil
	o.lockingStrategy.(*PostgresLockingStrategy).db = nil
	return err1, err2
}
