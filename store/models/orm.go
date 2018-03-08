package models

import (
	"log"
	"math/big"
	"path"
	"reflect"
	"strings"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/utils"
)

// ORM contains the database object used by Chainlink.
type ORM struct {
	*storm.DB
}

// NewORM initializes a new database file at the configured path.
func NewORM(dir string) *ORM {
	path := path.Join(dir, "db.bolt")
	orm := &ORM{initializeDatabase(path)}
	orm.migrate()
	return orm
}

func initializeDatabase(path string) *storm.DB {
	db, err := storm.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	return db
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

// SaveJob saves a job to the database.
func (orm *ORM) SaveJob(job *JobSpec) error {
	tx, err := orm.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for i, initr := range job.Initiators {
		job.Initiators[i].JobID = job.ID
		initr.JobID = job.ID
		if err := tx.Save(&initr); err != nil {
			return err
		}
	}
	if err := tx.Save(job); err != nil {
		return err
	}
	return tx.Commit()
}

// PendingJobRuns returns the JobRuns which have a status of "pending".
func (orm *ORM) PendingJobRuns() ([]JobRun, error) {
	runs := []JobRun{}
	err := orm.Where("Status", StatusPending, &runs)
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

// BridgeTypeFor returns the BridgeType for a given name.
func (orm *ORM) BridgeTypeFor(name string) (BridgeType, error) {
	tt := BridgeType{}
	err := orm.One("Name", strings.ToLower(name), &tt)
	return tt, err
}
