package eth

import (
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/orm"
)

var ExposedAppendLogChannel = appendLogChannel

func MakeLogBroadcast(log Log, _orm *orm.ORM) LogBroadcast {
	return &logBroadcast{log: log, orm: _orm, consumerID: models.NewID()}
}
