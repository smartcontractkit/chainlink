package migration1536521223

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/chainlink/store/assets"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/orm"
	null "gopkg.in/guregu/null.v3"
)

type Migration struct{}

func (m Migration) Timestamp() string {
	return "1536521223"
}

func (m Migration) Migrate(orm *orm.ORM) error {
	orm.InitializeModel(&JobSpec{})
	orm.InitializeModel(&JobRun{})
	orm.InitializeModel(&Initiator{})
	orm.InitializeModel(&Tx{})
	orm.InitializeModel(&TxAttempt{})
	orm.InitializeModel(&BridgeType{})
	orm.InitializeModel(&IndexableBlockNumber{})
	orm.InitializeModel(&User{})
	orm.InitializeModel(&Session{})
	orm.InitializeModel(&ServiceAgreement{})
	return nil
}

type JobSpec struct {
	ID        string      `json:"id" storm:"id,unique"`
	CreatedAt models.Time `json:"createdAt" storm:"index"`
}

type RunStatus string

type RunResult struct {
	JobRunID     string      `json:"jobRunId"`
	Data         models.JSON `json:"data"`
	Status       RunStatus   `json:"status"`
	ErrorMessage null.String `json:"error"`
	Amount       *big.Int    `json:"amount,omitempty"`
}

type TaskType string

type TaskSpec struct {
	Type          TaskType    `json:"type" storm:"index"`
	Confirmations uint64      `json:"confirmations"`
	Params        models.JSON `json:"-"`
}

type TaskRun struct {
	ID     string    `json:"id" storm:"id,unique"`
	Result RunResult `json:"result"`
	Status RunStatus `json:"status"`
	Task   TaskSpec  `json:"task"`
}

type JobRun struct {
	ID             string       `json:"id" storm:"id,unique"`
	JobID          string       `json:"jobId" storm:"index"`
	Result         RunResult    `json:"result" storm:"inline"`
	Status         RunStatus    `json:"status" storm:"index"`
	TaskRuns       []TaskRun    `json:"taskRuns" storm:"inline"`
	CreatedAt      time.Time    `json:"createdAt" storm:"index"`
	CompletedAt    null.Time    `json:"completedAt"`
	Initiator      Initiator    `json:"initiator"`
	CreationHeight *hexutil.Big `json:"creationHeight"`
	Overrides      RunResult    `json:"overrides"`
}

type Cron string

type Initiator struct {
	ID       int            `json:"id" storm:"id,increment"`
	JobID    string         `json:"jobId" storm:"index"`
	Type     string         `json:"type" storm:"index"`
	Schedule Cron           `json:"schedule,omitempty"`
	Time     models.Time    `json:"time,omitempty"`
	Ran      bool           `json:"ran,omitempty"`
	Address  common.Address `json:"address,omitempty" storm:"index"`
}

type Tx struct {
	ID       uint64         `storm:"id,increment,index"`
	From     common.Address `storm:"index"`
	To       common.Address
	Data     []byte
	Nonce    uint64 `storm:"index"`
	Value    *big.Int
	GasLimit uint64
	TxAttempt
}

type TxAttempt struct {
	Hash      common.Hash `storm:"id,unique"`
	TxID      uint64      `storm:"index"`
	GasPrice  *big.Int
	Confirmed bool
	Hex       string
	SentAt    uint64
}

type BridgeType struct {
	Name                 TaskType      `json:"name" storm:"id,unique"`
	URL                  models.WebURL `json:"url"`
	DefaultConfirmations uint64        `json:"defaultConfirmations"`
	IncomingToken        string        `json:"incomingToken"`
	OutgoingToken        string        `json:"outgoingToken"`
}

type IndexableBlockNumber struct {
	Number hexutil.Big `json:"number" storm:"id,unique"`
	Digits int         `json:"digits" storm:"index"`
	Hash   common.Hash `json:"hash"`
}

type User struct {
	Email          string      `json:"email" storm:"id,unique"`
	HashedPassword string      `json:"hashedPassword"`
	CreatedAt      models.Time `json:"createdAt" storm:"index"`
}

type Session struct {
	ID       string      `json:"id" storm:"id,unique"`
	LastUsed models.Time `json:"lastUsed" storm:"index"`
}

type Encumbrance struct {
	Payment    *assets.Link          `json:"payment"`
	Expiration uint64                `json:"expiration"`
	Oracles    []models.EIP55Address `json:"oracles"`
}

type ServiceAgreement struct {
	CreatedAt   models.Time      `json:"createdAt" storm:"index"`
	Encumbrance Encumbrance      `json:"encumbrance" storm:"inline"`
	ID          string           `json:"id" storm:"id,unique"`
	JobSpecID   string           `json:"jobSpecID"`
	RequestBody string           `json:"requestBody"`
	Signature   models.Signature `json:"signature"`
	JobSpec     JobSpec
}
