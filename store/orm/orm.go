package orm

import (
	"crypto/subtle"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // http://doc.gorm.io/database.html#connecting-to-a-database
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	_ "github.com/smartcontractkit/go-sqlite3" // http://doc.gorm.io/database.html#connecting-to-a-database
	"go.uber.org/multierr"
)

var (
	// ErrorNotFound is returned when finding a single value fails.
	ErrorNotFound = gorm.ErrRecordNotFound
)

// ORM contains the database object used by Chainlink.
type ORM struct {
	DB *gorm.DB
}

// NewORM initializes a new database file at the configured path.
func NewORM(path string) (*ORM, error) {
	db, err := initializeDatabase(path)
	if err != nil {
		return nil, fmt.Errorf("unable to init DB: %+v", err)
	}
	return &ORM{DB: db}, nil
}

func initializeDatabase(path string) (*gorm.DB, error) {
	dialect, err := DeduceDialect(path)
	if err != nil {
		return nil, err
	}
	db, err := gorm.Open(dialect, path)

	if err != nil {
		return nil, fmt.Errorf("unable to open gorm DB: %+v", err)
	}
	return db, nil
}

var validPostgresString = regexp.MustCompile("^(?:[A-Za-z]+=.* ?)+")
var validSqlitePath = regexp.MustCompile("^(.+)/([^/]+)$")
var validSqliteFilename = regexp.MustCompile(`^[\w\-. ]+$`)

// DeduceDialect returns the appropriate dialect for the passed connection string.
func DeduceDialect(path string) (string, error) {
	lower := strings.ToLower(path)
	if strings.HasPrefix(lower, "postgres://") || strings.HasPrefix(lower, "postgresql://") {
		return "postgres", nil
	} else if validPostgresString.MatchString(path) {
		return "postgres", nil
	} else if validSqlitePath.MatchString(path) || validSqliteFilename.MatchString(path) {
		return "sqlite3", nil
	}
	return "", fmt.Errorf("Unable to deduce sql dialect from path %s, please try a URI", path)
}

func ignoreRecordNotFound(db *gorm.DB) error {
	var merr error
	for _, e := range db.GetErrors() {
		if e != gorm.ErrRecordNotFound {
			merr = multierr.Append(merr, e)
		}
	}
	return merr
}

// Close closes the underlying database connection.
func (orm *ORM) Close() error {
	return orm.DB.Close()
}

// Where fetches multiple objects with "Find".
func (orm *ORM) Where(field string, value interface{}, instance interface{}) error {
	return orm.DB.Where(fmt.Sprintf("%v = ?", field), value).Find(instance).Error
}

// FindBridge looks up a Bridge by its Name.
func (orm *ORM) FindBridge(name string) (models.BridgeType, error) {
	var bt models.BridgeType
	return bt, orm.DB.First(&bt, "name = ?", name).Error
}

// PendingBridgeType returns the bridge type of the current pending task,
// or error if not pending bridge.
func (orm *ORM) PendingBridgeType(jr models.JobRun) (models.BridgeType, error) {
	nextTask := jr.NextTaskRun()
	if nextTask == nil {
		return models.BridgeType{}, errors.New("Cannot find the pending bridge type of a job run with no unfinished tasks")
	}
	return orm.FindBridge(nextTask.TaskSpec.Type.String())
}

// FindJob looks up a Job by its ID.
func (orm *ORM) FindJob(id string) (models.JobSpec, error) {
	var job models.JobSpec
	return job, orm.preloadJobs().First(&job, "id = ?", id).Error
}

// FindInitiator returns the single initiator defined by the passed ID.
func (orm *ORM) FindInitiator(ID uint) (models.Initiator, error) {
	initr := models.Initiator{}
	return initr, orm.DB.
		Set("gorm:auto_preload", true).
		First(&initr, "id = ?", ID).Error
}

func (orm *ORM) preloadJobs() *gorm.DB {
	return orm.DB.
		Preload("Initiators", func(db *gorm.DB) *gorm.DB {
			return db.Order("\"id\" asc")
		}).
		Preload("Tasks", func(db *gorm.DB) *gorm.DB {
			return db.Order("id asc")
		})
}

func (orm *ORM) preloadJobRuns() *gorm.DB {
	return orm.DB.
		Preload("Initiator").
		Preload("Overrides").
		Preload("Result").
		Preload("TaskRuns", func(db *gorm.DB) *gorm.DB {
			return db.Set("gorm:auto_preload", true).Order("task_spec_id asc")
		})
}

// FindJobRun looks up a JobRun by its ID.
func (orm *ORM) FindJobRun(id string) (models.JobRun, error) {
	var jr models.JobRun
	err := orm.preloadJobRuns().First(&jr, "id = ?", id).Error
	return jr, err
}

// SaveJobRun updates UpdatedAt for a JobRun and saves it
func (orm *ORM) SaveJobRun(run *models.JobRun) error {
	return orm.DB.Save(run).Error
}

// CreateJobRun inserts a new JobRun
func (orm *ORM) CreateJobRun(run *models.JobRun) error {
	return orm.DB.Create(run).Error
}

// FindServiceAgreement looks up a ServiceAgreement by its ID.
func (orm *ORM) FindServiceAgreement(id string) (models.ServiceAgreement, error) {
	var sa models.ServiceAgreement
	return sa, orm.DB.Set("gorm:auto_preload", true).First(&sa, "id = ?", id).Error
}

// Jobs fetches all jobs.
func (orm *ORM) Jobs(cb func(models.JobSpec) bool) error {
	offset := 0
	limit := 1000
	for {
		jobs := []models.JobSpec{}
		err := orm.preloadJobs().
			Limit(limit).
			Offset(offset).
			Find(&jobs).Error
		if err != nil {
			return err
		}
		for _, j := range jobs {
			if !cb(j) {
				return nil
			}
		}

		if len(jobs) < limit {
			return nil
		}

		offset += limit
	}
}

// JobRunsFor fetches all JobRuns with a given Job ID,
// sorted by their created at time.
func (orm *ORM) JobRunsFor(jobSpecID string) ([]models.JobRun, error) {
	runs := []models.JobRun{}
	err := orm.preloadJobRuns().
		Where("job_spec_id = ?", jobSpecID).
		Order("created_at desc").
		Find(&runs).Error
	return runs, err
}

// JobRunsCountFor returns the current number of runs for the job
func (orm *ORM) JobRunsCountFor(jobSpecID string) (int, error) {
	var count int
	err := orm.DB.
		Model(&models.JobRun{}).
		Where("job_spec_id = ?", jobSpecID).
		Count(&count).Error
	return count, err
}

// Sessions returns all sessions limited by the parameters.
func (orm *ORM) Sessions(offset, limit int) ([]models.Session, error) {
	var sessions []models.Session
	err := orm.DB.
		Set("gorm:auto_preload", true).
		Limit(limit).
		Offset(offset).
		Find(&sessions).Error
	return sessions, err
}

// CreateJob saves a job to the database and adds IDs to associated tables.
func (orm *ORM) CreateJob(job *models.JobSpec) error {
	for i := range job.Initiators {
		job.Initiators[i].JobSpecID = job.ID
	}
	return orm.DB.Create(job).Error
}

// CreateServiceAgreement saves a service agreement and it's associations to the
// database.
func (orm *ORM) CreateServiceAgreement(sa *models.ServiceAgreement) error {
	return orm.DB.Create(sa).Error
}

// JobRunsWithStatus returns the JobRuns which have the passed statuses.
func (orm *ORM) JobRunsWithStatus(statuses ...models.RunStatus) ([]models.JobRun, error) {
	runs := []models.JobRun{}
	err := orm.preloadJobRuns().Where("status IN (?)", statuses).Find(&runs).Error
	return runs, err
}

// AnyJobWithType returns true if there is at least one job associated with
// the type name specified and false otherwise
func (orm *ORM) AnyJobWithType(taskTypeName string) (bool, error) {
	db := orm.DB
	var taskSpec models.TaskSpec
	rval := db.Where("type = ?", taskTypeName).First(&taskSpec)
	found := !rval.RecordNotFound()
	return found, ignoreRecordNotFound(rval)
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
		Value:    models.NewBig(value),
		GasLimit: gasLimit,
	}
	return &tx, orm.DB.Save(&tx).Error
}

// ConfirmTx updates the database for the given transaction to
// show that the transaction has been confirmed on the blockchain.
func (orm *ORM) ConfirmTx(tx *models.Tx, txat *models.TxAttempt) error {
	txat.Confirmed = true
	tx.AssignTxAttempt(txat)
	return orm.DB.Save(tx).Save(txat).Error
}

// FindTx returns the specific transaction for the passed ID.
func (orm *ORM) FindTx(ID uint64) (*models.Tx, error) {
	tx := &models.Tx{}
	err := orm.DB.Set("gorm:auto_preload", true).First(tx, "id = ?", ID).Error
	return tx, err
}

// FindTxAttempt returns the specific transaction attempt with the hash.
func (orm *ORM) FindTxAttempt(hash common.Hash) (*models.TxAttempt, error) {
	txat := &models.TxAttempt{}
	err := orm.DB.Set("gorm:auto_preload", true).First(txat, "hash = ?", hash).Error
	return txat, err
}

// TxAttemptsFor returns the Transaction Attempts (TxAttempt) for a
// given Transaction ID (TxID).
func (orm *ORM) TxAttemptsFor(id uint64) ([]models.TxAttempt, error) {
	attempts := []models.TxAttempt{}
	err := orm.DB.
		Order("created_at asc").
		Where("tx_id = ?", id).
		Find(&attempts).Error
	return attempts, err
}

// AddTxAttempt creates a new transaction attempt and stores it
// in the database.
func (orm *ORM) AddTxAttempt(
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
		GasPrice: models.NewBig(etx.GasPrice()),
		Hex:      hex,
		TxID:     tx.ID,
		SentAt:   blkNum,
	}
	if !tx.Confirmed {
		tx.AssignTxAttempt(attempt)
	}
	err = orm.DB.Save(tx).Save(attempt).Error
	return attempt, err
}

// GetLastNonce retrieves the last known nonce in the database for an account
func (orm *ORM) GetLastNonce(address common.Address) (uint64, error) {
	var transaction models.Tx
	rval := orm.DB.Order("nonce desc").Where("\"from\" = ?", address).First(&transaction)
	return transaction.Nonce, ignoreRecordNotFound(rval)
}

// MarkRan will set Ran to true for a given initiator
func (orm *ORM) MarkRan(i *models.Initiator, ran bool) error {
	tx := orm.DB.Begin()
	var newi models.Initiator
	if err := tx.Select("ran").First(&newi, "ID = ?", i.ID).Error; err != nil {
		tx.Rollback()
		return err
	}

	if ran && newi.Ran {
		tx.Rollback()
		return fmt.Errorf("Initiator %v for job spec %s has already been run", i.ID, i.JobSpecID)
	}

	if err := tx.Model(i).UpdateColumn("ran", ran).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// FindUser will return the one API user, or an error.
func (orm *ORM) FindUser() (models.User, error) {
	user := models.User{}
	err := orm.DB.
		Set("gorm:auto_preload", true).
		Order("created_at desc").
		First(&user).Error
	return user, err
}

// AuthorizedUserWithSession will return the one API user if the Session ID exists
// and hasn't expired, and update session's LastUsed field.
func (orm *ORM) AuthorizedUserWithSession(sessionID string, sessionDuration time.Duration) (models.User, error) {
	if len(sessionID) == 0 {
		return models.User{}, errors.New("Session ID cannot be empty")
	}

	var session models.Session
	err := orm.DB.First(&session, "id = ?", sessionID).Error
	if err != nil {
		return models.User{}, err
	}
	now := time.Now()
	if session.LastUsed.Add(sessionDuration).Before(now) {
		return models.User{}, errors.New("Session has expired")
	}
	session.LastUsed = now
	if err := orm.DB.Save(&session).Error; err != nil {
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

	tx := orm.DB.Begin()
	if err := tx.Delete(&user).Error; err != nil {
		tx.Rollback()
		return user, err
	}

	if err := tx.Delete(models.Session{}).Error; err != nil {
		tx.Rollback()
		return user, err
	}

	return user, tx.Commit().Error
}

// DeleteUserSession will erase the session ID for the sole API User.
func (orm *ORM) DeleteUserSession(sessionID string) error {
	return orm.DB.Where("id = ?", sessionID).Delete(models.Session{}).Error
}

// DeleteBridgeType removes the bridge type with passed name.
func (orm *ORM) DeleteBridgeType(name models.TaskType) error {
	return orm.DB.Delete(&models.BridgeType{}, "name = ?", name).Error
}

// DeleteJobRun deletes the job run and corresponding task runs.
func (orm *ORM) DeleteJobRun(ID string) error {
	tx := orm.DB.Begin()
	if err := tx.Where("id = ?", ID).Delete(models.JobRun{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("job_run_id = ?", ID).Delete(models.TaskRun{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
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
		return session.ID, orm.DB.Save(&session).Error
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

// ClearSessions removes all sessions.
func (orm *ORM) ClearSessions() error {
	return orm.DB.Delete(models.Session{}).Error
}

// ClearNonCurrentSessions removes all sessions but the id passed in.
func (orm *ORM) ClearNonCurrentSessions(sessionID string) error {
	return orm.DB.Where("id <> ?", sessionID).Delete(models.Session{}).Error
}

// SortType defines the different sort orders available.
type SortType int

const (
	// Ascending is the sort order going up, i.e. 1,2,3.
	Ascending SortType = iota
	// Descending is the sort order going down, i.e. 3,2,1.
	Descending
)

func (s SortType) String() string {
	orderStr := "asc"
	if s == Descending {
		orderStr = "desc"
	}
	return orderStr
}

// JobsSorted returns many JobSpecs sorted by CreatedAt from the store adhering
// to the passed parameters.
func (orm *ORM) JobsSorted(order SortType, offset int, limit int) ([]models.JobSpec, int, error) {
	var count int
	err := orm.preloadJobs().Model(&models.JobSpec{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	var jobs []models.JobSpec
	rval := orm.DB.
		Set("gorm:auto_preload", true).
		Order(fmt.Sprintf("created_at %s", order.String())).
		Limit(limit).
		Offset(offset).
		Find(&jobs)
	return jobs, count, rval.Error
}

// TxFrom returns all transactions from a particular address.
func (orm *ORM) TxFrom(from common.Address) ([]models.Tx, error) {
	txs := []models.Tx{}
	return txs, orm.DB.
		Set("gorm:auto_preload", true).
		Find(&txs, "\"from\" = ?", from).Error
}

// Transactions returns all transactions limited by passed parameters.
func (orm *ORM) Transactions(offset, limit int) ([]models.Tx, error) {
	var txs []models.Tx
	err := orm.DB.
		Set("gorm:auto_preload", true).
		Order("id desc").Limit(limit).Offset(offset).
		Find(&txs).Error
	return txs, err
}

// TxAttempts returns the last tx attempts sorted by sent at descending.
func (orm *ORM) TxAttempts(offset, limit int) ([]models.TxAttempt, int, error) {
	var count int
	err := orm.DB.Model(&models.TxAttempt{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	var attempts []models.TxAttempt
	err = orm.DB.
		Set("gorm:auto_preload", true).
		Order("sent_at desc").Limit(limit).Offset(offset).
		Find(&attempts).Error
	return attempts, count, err
}

// JobRunsCount returns the total number of job runs
func (orm *ORM) JobRunsCount() (int, error) {
	var count int
	return count, orm.DB.Model(&models.JobRun{}).Count(&count).Error
}

// JobRunsSorted returns job runs ordered and filtered by the passed params.
func (orm *ORM) JobRunsSorted(order SortType, offset int, limit int) ([]models.JobRun, int, error) {
	count, err := orm.JobRunsCount()
	if err != nil {
		return nil, 0, err
	}

	var runs []models.JobRun
	err = orm.preloadJobRuns().
		Order(fmt.Sprintf("created_at %s", order.String())).
		Limit(limit).
		Offset(offset).
		Find(&runs).Error
	return runs, count, err
}

// JobRunsSortedFor returns job runs for a specific job spec ordered and
// filtered by the passed params.
func (orm *ORM) JobRunsSortedFor(id string, order SortType, offset int, limit int) ([]models.JobRun, int, error) {
	count, err := orm.JobRunsCountFor(id)
	if err != nil {
		return nil, 0, err
	}

	var runs []models.JobRun
	err = orm.preloadJobRuns().
		Order(fmt.Sprintf("created_at %s", order.String())).
		Limit(limit).
		Offset(offset).
		Find(&runs).Error
	return runs, count, err
}

// BridgeTypes returns bridge types ordered by name filtered limited by the
// passed params.
func (orm *ORM) BridgeTypes(offset int, limit int) ([]models.BridgeType, int, error) {
	db := orm.DB
	var count int
	err := db.Model(&models.BridgeType{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	var bridges []models.BridgeType
	err = db.Order("name asc").Limit(limit).Offset(offset).Find(&bridges).Error
	return bridges, count, err
}

// SaveUser saves the user.
func (orm *ORM) SaveUser(user *models.User) error {
	return orm.DB.Save(user).Error
}

// SaveSession saves the session.
func (orm *ORM) SaveSession(session *models.Session) error {
	return orm.DB.Save(session).Error
}

// CreateBridgeType saves the bridge type.
func (orm *ORM) CreateBridgeType(bt *models.BridgeType) error {
	return orm.DB.Create(bt).Error
}

// UpdateBridgeType updates the bridge type.
func (orm *ORM) UpdateBridgeType(bt *models.BridgeType) error {
	return orm.DB.Model(bt).Updates(bt).Error
}

// SaveTx saves the transaction.
func (orm *ORM) SaveTx(tx *models.Tx) error {
	return orm.DB.Save(tx).Error
}

// CreateInitiator saves the initiator.
func (orm *ORM) CreateInitiator(initr *models.Initiator) error {
	return orm.DB.Create(initr).Error
}

// SaveHead saves the indexable block number related to head tracker.
func (orm *ORM) SaveHead(n *models.IndexableBlockNumber) error {
	return orm.DB.Save(n).Error
}

// LastHead returns the last ordered IndexableBlockNumber.
func (orm *ORM) LastHead() (*models.IndexableBlockNumber, error) {
	number := &models.IndexableBlockNumber{}
	err := orm.DB.Order("digits desc, number desc").First(number).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return number, err
}

// DeleteStaleSessions deletes all sessions before the passed time.
func (orm *ORM) DeleteStaleSessions(before time.Time) error {
	return orm.DB.Where("last_used < ?", before).Delete(models.Session{}).Error
}

// SaveBulkDeleteRunTask saves the instance to the database.
func (orm *ORM) SaveBulkDeleteRunTask(task *models.BulkDeleteRunTask) error {
	return orm.DB.Save(task).Error
}

// FindBulkDeleteRunTask retrieves the instance with the id from the database.
func (orm *ORM) FindBulkDeleteRunTask(id string) (*models.BulkDeleteRunTask, error) {
	task := &models.BulkDeleteRunTask{}
	return task, orm.DB.Set("gorm:auto_preload", true).First(task, "ID = ?", id).Error
}

// BulkDeletesInProgress retrieves all bulk deletes in progress.
func (orm *ORM) BulkDeletesInProgress() ([]models.BulkDeleteRunTask, error) {
	deleteTasks := []models.BulkDeleteRunTask{}
	err := orm.DB.
		Set("gorm:auto_preload", true).
		Where("status = ?", models.BulkTaskStatusInProgress).
		Order("created_at asc").
		Find(&deleteTasks).Error
	return deleteTasks, err
}
