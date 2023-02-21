package txmgr

// A way to export and test private methods and variables in the txmgr orm
type Eorm struct {
	*orm
}

func NewEorm(o ORM) *Eorm {
	return &Eorm{o.(*orm)}
}

var UpdateEthTxAttemptUnbroadcast = updateEthTxAttemptUnbroadcast

var UpdateEthTxUnconfirm = updateEthTxUnconfirm

var DeleteEthReceipts = deleteEthReceipts
