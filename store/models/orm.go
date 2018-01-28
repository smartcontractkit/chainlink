package models

import (
	"log"
	"math/big"
	"path"
	"reflect"
	"strings"

	"github.com/asdine/storm"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/utils"
)

type ORM struct {
	*storm.DB
}

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

func (orm *ORM) Where(field string, value interface{}, instance interface{}) error {
	err := orm.Find(field, value, instance)
	if err == storm.ErrNotFound {
		emptySlice(instance)
		return nil
	}
	return err
}

func (orm *ORM) FindJob(id string) (*Job, error) {
	job := &Job{}
	err := orm.One("ID", id, job)
	return job, err
}

func (orm *ORM) InitBucket(model interface{}) error {
	return orm.Init(model)
}

func (orm *ORM) Jobs() ([]*Job, error) {
	var jobs []*Job
	err := orm.All(&jobs)
	return jobs, err
}

func (orm *ORM) JobRunsFor(job *Job) ([]*JobRun, error) {
	runs := []*JobRun{}
	err := orm.Where("JobID", job.ID, &runs)
	return runs, err
}

func emptySlice(to interface{}) {
	ref := reflect.ValueOf(to)
	results := reflect.MakeSlice(reflect.Indirect(ref).Type(), 0, 0)
	reflect.Indirect(ref).Set(results)
}

func (orm *ORM) SaveJob(job *Job) error {
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

func (orm *ORM) PendingJobRuns() ([]*JobRun, error) {
	runs := []*JobRun{}
	err := orm.Where("Status", StatusPending, &runs)
	return runs, err
}

func (orm *ORM) CreateTx(
	from common.Address,
	nonce uint64,
	to common.Address,
	data []byte,
	value *big.Int,
	gasLimit *big.Int,
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

func (orm *ORM) AttemptsFor(id uint64) ([]*TxAttempt, error) {
	attempts := []*TxAttempt{}
	if err := orm.Where("TxID", id, &attempts); err != nil {
		return attempts, err
	}
	return attempts, nil
}

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

func (orm *ORM) BridgeTypeFor(name string) (*BridgeType, error) {
	tt := &BridgeType{}
	err := orm.One("Name", strings.ToLower(name), tt)
	return tt, err
}
