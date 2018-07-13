package models

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	bolt "github.com/coreos/bbolt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/utils"
)

// ORM contains the database object used by Chainlink.
type ORM struct {
	*storm.DB
}

// NewORM initializes a new database file at the configured path.
func NewORM(path string) (*ORM, error) {
	db, err := initializeDatabase(path)
	if err != nil {
		return nil, fmt.Errorf("unable to init DB: %+v", err)
	}
	orm := &ORM{db}
	orm.migrate()
	return orm, nil
}

func initializeDatabase(path string) (*storm.DB, error) {
	db, err := storm.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open stormDB: %+v", err)
	}
	return db, nil
}

// GetBolt returns BoltDB from the ORM
func (orm *ORM) GetBolt() *bolt.DB {
	return orm.DB.Bolt
}

// Where fetches multiple objects with "Find" in Storm.
func (orm *ORM) Where(field string, value interface{}, instance interface{}) error {
	err := orm.Find(field, value, instance)
	if err == storm.ErrNotFound {
		emptySlice(instance)
		return nil
	}
	return err
}

func emptySlice(to interface{}) {
	ref := reflect.ValueOf(to)
	results := reflect.MakeSlice(reflect.Indirect(ref).Type(), 0, 0)
	reflect.Indirect(ref).Set(results)
}

// FindBridge looks up a Bridge by its Name.
func (orm *ORM) FindBridge(name string) (BridgeType, error) {
	var bt BridgeType

	tt, err := NewTaskType(name)
	if err != nil {
		return bt, err
	}

	err = orm.One("Name", tt, &bt)
	return bt, err
}

// FindJob looks up a Job by its ID.
func (orm *ORM) FindJob(id string) (JobSpec, error) {
	var job JobSpec
	err := orm.One("ID", id, &job)
	return job, err
}

// FindJobRun looks up a JobRun by its ID.
func (orm *ORM) FindJobRun(id string) (JobRun, error) {
	var jr JobRun
	err := orm.One("ID", id, &jr)
	return jr, err
}

// InitBucket initializes buckets and indexes before saving an object.
func (orm *ORM) InitBucket(model interface{}) error {
	return orm.Init(model)
}

// Jobs fetches all jobs.
func (orm *ORM) Jobs() ([]JobSpec, error) {
	var jobs []JobSpec
	err := orm.All(&jobs)
	return jobs, err
}

// JobRunsFor fetches all JobRuns with a given Job ID,
// sorted by their created at time.
func (orm *ORM) JobRunsFor(jobID string) ([]JobRun, error) {
	runs := []JobRun{}
	err := orm.Select(q.Eq("JobID", jobID)).OrderBy("CreatedAt").Reverse().Find(&runs)
	if err == storm.ErrNotFound {
		return []JobRun{}, nil
	}
	return runs, err
}

// JobRunsCountFor returns the current number of runs for the job
func (orm *ORM) JobRunsCountFor(jobID string) (int, error) {
	query := orm.Select(q.Eq("JobID", jobID))
	return query.Count(&JobRun{})
}

// SaveJob saves a job to the database and adds IDs to associated tables.
func (orm *ORM) SaveJob(job *JobSpec) error {
	tx, err := orm.Begin(true)
	if err != nil {
		return fmt.Errorf("error starting transaction: %+v", err)
	}
	defer tx.Rollback()

	for i := range job.Initiators {
		job.Initiators[i].JobID = job.ID
		if err := tx.Save(&job.Initiators[i]); err != nil {
			return fmt.Errorf("error saving Job Initiators: %+v", err)
		}
	}
	if err := tx.Save(job); err != nil {
		return fmt.Errorf("error saving job: %+v", err)
	}
	return tx.Commit()
}

// SaveCreationHeight stores the JobRun in the database with the given
// block number.
func (orm *ORM) SaveCreationHeight(jr JobRun, bn *IndexableBlockNumber) (JobRun, error) {
	if jr.CreationHeight != nil || bn == nil {
		return jr, nil
	}

	dup := bn.Number
	jr.CreationHeight = &dup
	return jr, orm.Save(&jr)
}

// JobRunsWithStatus returns the JobRuns which have the passed statuses.
func (orm *ORM) JobRunsWithStatus(statuses ...RunStatus) ([]JobRun, error) {
	runs := []JobRun{}
	err := orm.Select(q.In("Status", statuses)).Find(&runs)
	if err == storm.ErrNotFound {
		return []JobRun{}, nil
	}

	return runs, err
}

// CreateTx saves the properties of an Ethereum transaction to the database.
func (orm *ORM) CreateTx(
	from common.Address,
	nonce uint64,
	to common.Address,
	data []byte,
	value *big.Int,
	gasLimit uint64,
) (*Tx, error) {
	tx := Tx{
		From:     from,
		To:       to,
		Nonce:    nonce,
		Data:     data,
		Value:    value,
		GasLimit: gasLimit,
	}
	return &tx, orm.Save(&tx)
}

// ConfirmTx updates the database for the given transaction to
// show that the transaction has been confirmed on the blockchain.
func (orm *ORM) ConfirmTx(tx *Tx, txat *TxAttempt) error {
	dbtx, err := orm.Begin(true)
	if err != nil {
		return err
	}
	defer dbtx.Rollback()

	txat.Confirmed = true
	tx.TxAttempt = *txat
	if err := dbtx.Save(tx); err != nil {
		return err
	}
	if err := dbtx.Save(txat); err != nil {
		return err
	}
	return dbtx.Commit()
}

// AttemptsFor returns the Transaction Attempts (TxAttempt) for a
// given Transaction ID (TxID).
func (orm *ORM) AttemptsFor(id uint64) ([]TxAttempt, error) {
	attempts := []TxAttempt{}
	if err := orm.Where("TxID", id, &attempts); err != nil {
		return attempts, err
	}
	return attempts, nil
}

// AddAttempt creates a new transaction attempt and stores it
// in the database.
func (orm *ORM) AddAttempt(
	tx *Tx,
	etx *types.Transaction,
	blkNum uint64,
) (*TxAttempt, error) {
	hex, err := utils.EncodeTxToHex(etx)
	if err != nil {
		return nil, err
	}
	attempt := &TxAttempt{
		Hash:     etx.Hash(),
		GasPrice: etx.GasPrice(),
		Hex:      hex,
		TxID:     tx.ID,
		SentAt:   blkNum,
	}
	if !tx.Confirmed {
		tx.TxAttempt = *attempt
	}
	dbtx, err := orm.Begin(true)
	if err != nil {
		return nil, err
	}
	defer dbtx.Rollback()
	if err = dbtx.Save(tx); err != nil {
		return nil, err
	}
	if err = dbtx.Save(attempt); err != nil {
		return nil, err
	}

	return attempt, dbtx.Commit()
}

// GetLastNonce retrieves the last known nonce in the database for an account
func (orm *ORM) GetLastNonce(address common.Address) (uint64, error) {
	var transactions []Tx
	query := orm.Select(q.Eq("From", address))
	if err := query.Limit(1).OrderBy("Nonce").Reverse().Find(&transactions); err == storm.ErrNotFound {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	return transactions[0].Nonce, nil
}

// MarkRan will set Ran to true for a given initiator
func (orm *ORM) MarkRan(i *Initiator) error {
	dbtx, err := orm.Begin(true)
	if err != nil {
		return err
	}
	defer dbtx.Rollback()

	var ir Initiator
	if err := orm.One("ID", i.ID, &ir); err != nil {
		return err
	}

	if ir.Ran {
		return fmt.Errorf("Job runner: Initiator: %v cannot run more than once", ir.ID)
	}

	i.Ran = true
	if err := dbtx.Save(i); err != nil {
		return err
	}
	return dbtx.Commit()
}

// DatabaseAccessError is an error that occurs during database access.
type DatabaseAccessError struct {
	msg string
}

func (e *DatabaseAccessError) Error() string { return e.msg }

// NewDatabaseAccessError returns a database access error.
func NewDatabaseAccessError(msg string) error {
	return &DatabaseAccessError{msg}
}

func (orm *ORM) FindUser() (User, error) {
	var users []User
	err := orm.All(&users)
	if err != nil {
		return User{}, err
	}

	if len(users) == 0 {
		return User{}, storm.ErrNotFound
	}

	return users[0], nil
}

func (orm *ORM) FindUserBySession(sessionId string) (User, error) {
	user, err := orm.FindUser()
	if err != nil {
		return User{}, err
	}
	if sessionId != user.SessionID {
		return User{}, storm.ErrNotFound
	}
	return user, nil
}

func (orm *ORM) DeleteUser() (User, error) {
	user, err := orm.FindUser()
	if err != nil {
		return user, err
	}
	return user, orm.DeleteStruct(&user)
}
