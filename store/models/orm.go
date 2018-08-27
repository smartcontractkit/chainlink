package models

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"sort"
	"time"

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
func NewORM(path string, duration time.Duration) (*ORM, error) {
	db, err := initializeDatabase(path, duration)
	if err != nil {
		return nil, fmt.Errorf("unable to init DB: %+v", err)
	}
	orm := &ORM{db}
	orm.migrate()
	return orm, nil
}

func initializeDatabase(path string, duration time.Duration) (*storm.DB, error) {
	options := storm.BoltOptions(0600, &bolt.Options{Timeout: duration})
	db, err := storm.Open(path, options)
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

// PendingBridgeType returns the bridge type of the current pending task,
// or error if not pending bridge.
func (orm *ORM) PendingBridgeType(jr JobRun) (BridgeType, error) {
	unfinished := jr.UnfinishedTaskRuns()
	if len(unfinished) == 0 {
		return BridgeType{}, errors.New("Cannot find the pending bridge type of a job run with no unfinished tasks")
	}
	return orm.FindBridge(unfinished[0].Task.Type.String())
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

// FindServiceAgreement looks up a ServiceAgreement by its ID.
func (orm *ORM) FindServiceAgreement(id common.Hash) (ServiceAgreement, error) {
	var sa ServiceAgreement
	return sa, orm.One("ID", id.String(), &sa)
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
	err := orm.Find("JobID", jobID, &runs) // Use Find to leverage db index
	if err == storm.ErrNotFound {
		return []JobRun{}, nil
	}
	sort.Sort(jobRunSorterAscending(runs))
	return runs, err
}

type jobRunSorterAscending []JobRun

func (jrs jobRunSorterAscending) Len() int      { return len(jrs) }
func (jrs jobRunSorterAscending) Swap(i, j int) { jrs[i], jrs[j] = jrs[j], jrs[i] }
func (jrs jobRunSorterAscending) Less(i, j int) bool {
	return jrs[i].CreatedAt.Sub(jrs[j].CreatedAt) > 0
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

	if err := saveJobSpec(job, tx); err != nil {
		return err
	}
	return tx.Commit()
}

func saveJobSpec(job *JobSpec, tx storm.Node) error {
	for i := range job.Initiators {
		job.Initiators[i].JobID = job.ID
		if err := tx.Save(&job.Initiators[i]); err != nil {
			return fmt.Errorf("error saving Job Initiators: %+v", err)
		}
	}
	if err := tx.Save(job); err != nil {
		return fmt.Errorf("error saving job: %+v", err)
	}
	return nil
}

// SaveServiceAgreement saves a service agreement and it's associations to the
// database.
func (orm *ORM) SaveServiceAgreement(sa *ServiceAgreement) error {
	tx, err := orm.Begin(true)
	if err != nil {
		return fmt.Errorf("error starting transaction: %+v", err)
	}
	defer tx.Rollback()

	if err := saveJobSpec(&sa.jobSpec, tx); err != nil {
		return fmt.Errorf("error saving service agreement: %+v", err)
	}

	sa.JobSpecID = sa.jobSpec.ID
	if err := tx.Save(sa); err != nil {
		return fmt.Errorf("error saving service agreement: %+v", err)
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

// AnyJobWithType returns true if there is at least one job associated with
// the type name specified and false otherwise
func (orm *ORM) AnyJobWithType(taskTypeName string) (bool, error) {
	jobs := []JobSpec{}
	err := orm.All(&jobs)
	if err != nil {
		return false, err
	}
	ts, err := NewTaskType(taskTypeName)
	if err != nil {
		return false, err
	}
	for i := range jobs {
		for j := range jobs[i].Tasks {
			if jobs[i].Tasks[j].Type == ts {
				return true, nil
			}
		}
	}
	return false, nil
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

// FindUser will return the one API user, or an error.
func (orm *ORM) FindUser() (User, error) {
	var users []User
	err := orm.AllByIndex("CreatedAt", &users, storm.Limit(1), storm.Reverse())
	if err != nil {
		return User{}, err
	}

	if len(users) == 0 {
		return User{}, storm.ErrNotFound
	}

	return users[0], nil
}

// AuthorizedUserWithSession will return the one API user if the Session ID exists
// and hasn't expired, and update session's LastUsed field.
func (orm *ORM) AuthorizedUserWithSession(sessionID string, sessionDuration time.Duration) (User, error) {
	if len(sessionID) == 0 {
		return User{}, errors.New("Session ID cannot be empty")
	}

	var session Session
	err := orm.One("ID", sessionID, &session)
	if err != nil {
		return User{}, err
	}
	now := time.Now()
	if session.LastUsed.Time.Add(sessionDuration).Before(now) {
		return User{}, errors.New("Session has expired")
	}
	session.LastUsed = Time{Time: now}
	if err := orm.Save(&session); err != nil {
		return User{}, err
	}
	return orm.FindUser()
}

// DeleteUser will delete the API User in the db.
func (orm *ORM) DeleteUser() (User, error) {
	user, err := orm.FindUser()
	if err != nil {
		return user, err
	}

	tx, err := orm.Begin(true)
	if err != nil {
		return user, fmt.Errorf("error starting transaction: %+v", err)
	}
	defer tx.Rollback()

	err = tx.DeleteStruct(&user)
	if err != nil {
		return user, err
	}
	err = tx.Drop(&Session{})
	if err != nil {
		return user, err
	}
	err = tx.Init(&Session{})
	if err != nil {
		return user, err
	}
	return user, tx.Commit()
}

// DeleteUserSession will erase the session ID for the sole API User.
func (orm *ORM) DeleteUserSession(sessionID string) error {
	session := Session{ID: sessionID}
	return orm.DeleteStruct(&session)
}

// CreateSession will check the password in the SessionRequest against
// the hashed API User password in the db.
func (orm *ORM) CreateSession(sr SessionRequest) (string, error) {
	user, err := orm.FindUser()
	if err != nil {
		return "", err
	}

	if sr.Email != user.Email {
		return "", errors.New("Invalid email")
	}

	if utils.CheckPasswordHash(sr.Password, user.HashedPassword) {
		session := NewSession()
		return session.ID, orm.Save(&session)
	}
	return "", errors.New("Invalid password")
}
