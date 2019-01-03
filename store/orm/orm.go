package orm

import (
	"crypto/subtle"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"sort"
	"time"

	"github.com/asdine/storm"
	"github.com/asdine/storm/index"
	"github.com/asdine/storm/q"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	bolt "go.etcd.io/bbolt"
)

var (
	// ErrorInvalidCallbackSignature is returned in AllInBatches if the incorrect function signature is passed.
	ErrorInvalidCallbackSignature = errors.New("AllInBatches callback has incorrect function signature, must return bool")
	// ErrorInvalidCallbackModel is returned in AllInBatches if the model and bucket do not match types.
	ErrorInvalidCallbackModel = errors.New("AllInBatches callback has incorrect model, must match bucket")
	// ErrorNotFound is returned when finding a single value fails.
	ErrorNotFound = storm.ErrNotFound
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
	return &ORM{db}, nil
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
	if err == ErrorNotFound {
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
func (orm *ORM) FindBridge(name string) (models.BridgeType, error) {
	var bt models.BridgeType

	tt, err := models.NewTaskType(name)
	if err != nil {
		return bt, err
	}

	err = orm.One("Name", tt, &bt)
	return bt, err
}

// PendingBridgeType returns the bridge type of the current pending task,
// or error if not pending bridge.
func (orm *ORM) PendingBridgeType(jr models.JobRun) (models.BridgeType, error) {
	nextTask := jr.NextTaskRun()
	if nextTask == nil {
		return models.BridgeType{}, errors.New("Cannot find the pending bridge type of a job run with no unfinished tasks")
	}
	return orm.FindBridge(nextTask.Task.Type.String())
}

// FindJob looks up a Job by its ID.
func (orm *ORM) FindJob(id string) (models.JobSpec, error) {
	var job models.JobSpec
	err := orm.One("ID", id, &job)
	return job, err
}

// FindJobRun looks up a JobRun by its ID.
func (orm *ORM) FindJobRun(id string) (models.JobRun, error) {
	var jr models.JobRun
	err := orm.One("ID", id, &jr)
	return jr, err
}

// SaveJobRun updates UpdatedAt for a JobRun and saves it
func (orm *ORM) SaveJobRun(run *models.JobRun) error {
	run.UpdatedAt = time.Now()
	return orm.Save(run)
}

// FindServiceAgreement looks up a ServiceAgreement by its ID.
func (orm *ORM) FindServiceAgreement(id string) (models.ServiceAgreement, error) {
	var sa models.ServiceAgreement
	return sa, orm.One("ID", id, &sa)
}

// InitBucket initializes buckets and indexes before saving an object.
func (orm *ORM) InitBucket(model interface{}) error {
	return orm.Init(model)
}

// Jobs fetches all jobs.
func (orm *ORM) Jobs(cb func(models.JobSpec) bool) error {
	var bucket []models.JobSpec
	return orm.AllInBatches(&bucket, func(j models.JobSpec) bool {
		return cb(j)
	})
}

// JobRunsFor fetches all JobRuns with a given Job ID,
// sorted by their created at time.
func (orm *ORM) JobRunsFor(jobID string) ([]models.JobRun, error) {
	runs := []models.JobRun{}
	err := orm.Find("JobID", jobID, &runs) // Use Find to leverage db index
	if err == ErrorNotFound {
		return []models.JobRun{}, nil
	}
	sort.Sort(jobRunSorterAscending(runs))
	return runs, err
}

type jobRunSorterAscending []models.JobRun

func (jrs jobRunSorterAscending) Len() int      { return len(jrs) }
func (jrs jobRunSorterAscending) Swap(i, j int) { jrs[i], jrs[j] = jrs[j], jrs[i] }
func (jrs jobRunSorterAscending) Less(i, j int) bool {
	return jrs[i].CreatedAt.Sub(jrs[j].CreatedAt) > 0
}

// JobRunsCountFor returns the current number of runs for the job
func (orm *ORM) JobRunsCountFor(jobID string) (int, error) {
	query := orm.Select(q.Eq("JobID", jobID))
	return query.Count(&models.JobRun{})
}

// SaveJob saves a job to the database and adds IDs to associated tables.
func (orm *ORM) SaveJob(job *models.JobSpec) error {
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

func saveJobSpec(job *models.JobSpec, tx storm.Node) error {
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
func (orm *ORM) SaveServiceAgreement(sa *models.ServiceAgreement) error {
	tx, err := orm.Begin(true)
	if err != nil {
		return fmt.Errorf("error starting transaction: %+v", err)
	}
	defer tx.Rollback()

	if err := saveJobSpec(&sa.JobSpec, tx); err != nil {
		return fmt.Errorf("error saving service agreement: %+v", err)
	}

	sa.JobSpecID = sa.JobSpec.ID
	if err := tx.Save(sa); err != nil {
		return fmt.Errorf("error saving service agreement: %+v", err)
	}

	return tx.Commit()
}

// JobRunsWithStatus returns the JobRuns which have the passed statuses.
func (orm *ORM) JobRunsWithStatus(statuses ...models.RunStatus) ([]models.JobRun, error) {
	runs := []models.JobRun{}
	err := orm.Select(q.In("Status", statuses)).Find(&runs)
	if err == ErrorNotFound {
		return []models.JobRun{}, nil
	}

	return runs, err
}

// AnyJobWithType returns true if there is at least one job associated with
// the type name specified and false otherwise
func (orm *ORM) AnyJobWithType(taskTypeName string) (bool, error) {
	ts, err := models.NewTaskType(taskTypeName)
	if err != nil {
		return false, err
	}

	var found bool
	err = orm.Jobs(func(j models.JobSpec) bool {
		for _, t := range j.Tasks {
			if t.Type == ts {
				found = true
				return false
			}
		}
		return true
	})

	return found, err
}

// CreateTx saves the properties of an Ethereum transaction to the database.
func (orm *ORM) CreateTx(
	from common.Address,
	nonce uint64,
	to common.Address,
	data []byte,
	value *big.Int,
	gasLimit uint64,
) (*models.Tx, error) {
	tx := models.Tx{
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
func (orm *ORM) ConfirmTx(tx *models.Tx, txat *models.TxAttempt) error {
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
func (orm *ORM) AttemptsFor(id uint64) ([]models.TxAttempt, error) {
	attempts := []models.TxAttempt{}
	if err := orm.Where("TxID", id, &attempts); err != nil {
		return attempts, err
	}
	return attempts, nil
}

// AddAttempt creates a new transaction attempt and stores it
// in the database.
func (orm *ORM) AddAttempt(
	tx *models.Tx,
	etx *types.Transaction,
	blkNum uint64,
) (*models.TxAttempt, error) {
	hex, err := utils.EncodeTxToHex(etx)
	if err != nil {
		return nil, err
	}
	attempt := &models.TxAttempt{
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
	var transactions []models.Tx
	query := orm.Select(q.Eq("From", address))
	if err := query.Limit(1).OrderBy("Nonce").Reverse().Find(&transactions); err == ErrorNotFound {
		return 0, nil
	} else if err != nil {
		return 0, err
	}

	return transactions[0].Nonce, nil
}

// MarkRan will set Ran to true for a given initiator
func (orm *ORM) MarkRan(i *models.Initiator) error {
	dbtx, err := orm.Begin(true)
	if err != nil {
		return err
	}
	defer dbtx.Rollback()

	var ir models.Initiator
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
func (orm *ORM) FindUser() (models.User, error) {
	var users []models.User
	err := orm.AllByIndex("CreatedAt", &users, storm.Limit(1), storm.Reverse())
	if err != nil {
		return models.User{}, err
	}

	if len(users) == 0 {
		return models.User{}, ErrorNotFound
	}

	return users[0], nil
}

// AuthorizedUserWithSession will return the one API user if the Session ID exists
// and hasn't expired, and update session's LastUsed field.
func (orm *ORM) AuthorizedUserWithSession(sessionID string, sessionDuration time.Duration) (models.User, error) {
	if len(sessionID) == 0 {
		return models.User{}, errors.New("Session ID cannot be empty")
	}

	var session models.Session
	err := orm.One("ID", sessionID, &session)
	if err != nil {
		return models.User{}, err
	}
	now := time.Now()
	if session.LastUsed.Time.Add(sessionDuration).Before(now) {
		return models.User{}, errors.New("Session has expired")
	}
	session.LastUsed = models.Time{Time: now}
	if err := orm.Save(&session); err != nil {
		return models.User{}, err
	}
	return orm.FindUser()
}

// DeleteUser will delete the API User in the db.
func (orm *ORM) DeleteUser() (models.User, error) {
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
	err = tx.Drop(&models.Session{})
	if err != nil {
		return user, err
	}
	err = tx.Init(&models.Session{})
	if err != nil {
		return user, err
	}
	return user, tx.Commit()
}

// DeleteUserSession will erase the session ID for the sole API User.
func (orm *ORM) DeleteUserSession(sessionID string) error {
	session := models.Session{ID: sessionID}
	return orm.DeleteStruct(&session)
}

// CreateSession will check the password in the SessionRequest against
// the hashed API User password in the db.
func (orm *ORM) CreateSession(sr models.SessionRequest) (string, error) {
	user, err := orm.FindUser()
	if err != nil {
		return "", err
	}

	if !constantTimeEmailCompare(sr.Email, user.Email) {
		return "", errors.New("Invalid email")
	}

	if utils.CheckPasswordHash(sr.Password, user.HashedPassword) {
		session := models.NewSession()
		return session.ID, orm.Save(&session)
	}
	return "", errors.New("Invalid password")
}

const constantTimeEmailLength = 256

func constantTimeEmailCompare(left, right string) bool {
	length := utils.MaxInt(constantTimeEmailLength, len(left), len(right))
	leftBytes := make([]byte, length)
	rightBytes := make([]byte, length)
	copy(leftBytes, left)
	copy(rightBytes, right)
	return subtle.ConstantTimeCompare(leftBytes, rightBytes) == 1
}

// InitializeModel uses reflection on the passed klass to generate a bucket
// of the same type name.
func (orm *ORM) InitializeModel(klass interface{}) error {
	return orm.InitBucket(klass)
}

// AllInBatches iterates over every single entry in the passed bucket without holding
// the entire contents in memory, pulling down batches and streaming over each entry.
// Be sure not to use the passed bucket parameter as it is used as a buffer.
func (orm *ORM) AllInBatches(bucket interface{}, callback interface{}, optionalBatchSize ...int) error {
	skip := 0
	batchSize := 1000
	if len(optionalBatchSize) > 0 {
		batchSize = optionalBatchSize[0]
	}

	vcallback := reflect.ValueOf(callback)
	tcallback := reflect.TypeOf(callback)
	if tcallback.NumOut() != 1 || tcallback.Out(0).Kind() != reflect.Bool {
		return ErrorInvalidCallbackSignature
	}

	if tcallback.NumIn() != 1 || tcallback.In(0) != underlyingBucketType(bucket) {
		return ErrorInvalidCallbackModel
	}

	for {
		err := orm.All(bucket, storm.Limit(batchSize), storm.Skip(skip))
		if err != nil {
			return err
		}

		slice := reflect.ValueOf(bucket).Elem()
		for i := 0; i < slice.Len(); i++ {
			e := slice.Index(i)
			rval := vcallback.Call([]reflect.Value{e})[0].Bool()
			if !rval {
				return nil
			}
		}

		if slice.Len() < batchSize {
			return nil
		}

		skip += batchSize
	}
}

func underlyingBucketType(bucket interface{}) reflect.Type {
	ref := reflect.ValueOf(bucket)
	sliceType := reflect.Indirect(ref).Type()
	elemType := sliceType.Elem()
	return elemType
}

// ClearNonCurrentSessions removes all sessions but the id passed in.
func (orm *ORM) ClearNonCurrentSessions(sessionID string) error {
	var sessions []models.Session
	err := orm.Select(q.Not(q.Eq("ID", sessionID))).Find(&sessions)
	if err != nil && err != ErrorNotFound {
		return err
	}

	for _, s := range sessions {
		err := orm.DeleteStruct(&s)
		if err != nil {
			return err
		}
	}

	return nil
}

// SortType defines the different sort orders available.
type SortType int

const (
	// Ascending is the sort order going up, i.e. 1,2,3.
	Ascending SortType = iota
	// Descending is the sort order going down, i.e. 3,2,1.
	Descending
)

func stormOrder(st SortType) func(*index.Options) {
	so := func(opts *index.Options) {}
	if st == Descending {
		so = storm.Reverse()
	}
	return so
}

// SortedJobs returns many JobSpecs sorted by CreatedAt from the store adhering
// to the passed parameters.
func (orm *ORM) SortedJobs(order SortType, offset int, limit int) ([]models.JobSpec, error) {
	stormOffset := storm.Skip(offset)
	stormLimit := storm.Limit(limit)

	var jobs []models.JobSpec
	err := orm.AllByIndex("CreatedAt", &jobs, stormOrder(order), stormOffset, stormLimit)
	return jobs, err
}

// GetTxAttempts returns the last tx attempts sorted by sent at descending.
func (orm *ORM) GetTxAttempts(offset int, limit int) ([]models.TxAttempt, int, error) {
	var attempts []models.TxAttempt
	count, err := orm.Count(&models.TxAttempt{})
	if err != nil {
		return nil, 0, err
	}
	query := orm.Select().OrderBy("SentAt").Reverse().Limit(limit).Skip(offset)
	err = query.Find(&attempts)
	return attempts, count, err
}

// SortedJobRuns returns job runs ordered and filtered by the passed params.
func (orm *ORM) SortedJobRuns(order SortType, offset int, limit int) ([]models.JobRun, int, error) {
	count, err := orm.Count(&models.JobRun{})
	if err != nil {
		return nil, 0, err
	}

	query := orm.Select().OrderBy("CreatedAt").Limit(limit).Skip(offset)
	if order == Descending {
		query = query.Reverse()
	}

	var runs []models.JobRun
	err = query.Find(&runs)
	return runs, count, err
}

// SortedJobRunsFor returns job runs for a specific job spec ordered and
// filtered by the passed params.
func (orm *ORM) SortedJobRunsFor(id string, order SortType, offset int, limit int) ([]models.JobRun, int, error) {
	count, err := orm.JobRunsCountFor(id)
	if err != nil {
		return nil, 0, err
	}

	query := orm.Select(q.Eq("JobID", id)).OrderBy("CreatedAt").Limit(limit).Skip(offset)
	if order == Descending {
		query = query.Reverse()
	}

	var runs []models.JobRun
	err = query.Find(&runs)
	return runs, count, err
}

// BridgeTypes returns bridge types ordered by name filtered limited by the
// passed params.
func (orm *ORM) BridgeTypes(offset int, limit int) ([]models.BridgeType, int, error) {
	count, err := orm.Count(&models.BridgeType{})
	if err != nil {
		return nil, 0, err
	}

	var bridges []models.BridgeType
	err = orm.AllByIndex("Name", &bridges, storm.Skip(offset), storm.Limit(limit))
	return bridges, count, err
}

// SaveUser saves the user.
func (orm *ORM) SaveUser(user *models.User) error {
	return orm.Save(user)
}

// SaveSession saves the session.
func (orm *ORM) SaveSession(session *models.Session) error {
	return orm.Save(session)
}

// SaveBridgeType saves the bridge type.
func (orm *ORM) SaveBridgeType(bt *models.BridgeType) error {
	return orm.Save(bt)
}
