package models

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"reflect"
	"strings"

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
func NewORM(path string) *ORM {
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
	err := orm.One("Name", name, &bt)
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

// BridgeTypeFor returns the BridgeType for a given name.
func (orm *ORM) BridgeTypeFor(name string) (BridgeType, error) {
	tt := BridgeType{}
	err := orm.One("Name", strings.ToLower(name), &tt)
	return tt, err
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

// ParseQuery parses the JSON parameters stored in a QueryObject's fields,
// creates the appropriate queries and returns them in an array
func (orm *ORM) ParseQuery(fieldValue json.RawMessage, model interface{}, lookup string) ([]q.Matcher, error) {
	var m interface{}
	var query []q.Matcher
	queryMap := map[string]func(string, interface{}) q.Matcher{
		"Eq": q.Eq, "Gt": q.Gt, "Gte": q.Gte, "In": q.In,
		"Lt": q.Lt, "Lte": q.Lte, "StrictEq": q.StrictEq,
	}

	json.Unmarshal(fieldValue, &m)
	mapKeys, ok := m.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("ParseQuery: Invalid params for %v query", lookup)
	}
	nonEmptyKeys := utils.GetStringKeys(mapKeys)

	json.Unmarshal(fieldValue, &model)

	s := reflect.Indirect(reflect.ValueOf(model))
	typeOfValue := s.Type()
	for i := 0; i < s.NumField(); i++ {
		fieldKey := typeOfValue.Field(i).Name
		fieldTag := typeOfValue.Field(i).Tag.Get("json")
		if utils.SliceIndex(nonEmptyKeys, fieldTag) == -1 {
			continue
		}
		field := s.Field(i)
		if !(field.CanInterface()) {
			return nil, fmt.Errorf("ParseQuery: Field %v can't interface", field)
		}
		curFieldValue := field.Interface()
		switch lookup {
		case "Re":
			s, ok := field.Interface().(string)
			if !ok {
				return nil, fmt.Errorf("ParseQuery: Type string required for regex query on : %v", fieldKey)
			}
			query = append(query, q.Re(fieldKey, s))

		default:
			queryFunction, ok := queryMap[lookup]
			if !ok {
				return nil, fmt.Errorf("ParseQuery: Invalid query operation %v for field %v", lookup, field)
			}
			query = append(query, queryFunction(fieldKey, curFieldValue))
		}
	}
	if len(query) != len(nonEmptyKeys) {
		return nil, fmt.Errorf("ParseQuery: Unknown field(s) for this model")
	}
	return query, nil
}

// BuildQuery builds a multi-field query based on the parameters of a QueryObject
func (orm *ORM) BuildQuery(value interface{}, model interface{}) ([]q.Matcher, error) {
	var dbSelect []q.Matcher
	s := reflect.Indirect(reflect.ValueOf(value))
	typeOfValue := s.Type()
	for i := 0; i < s.NumField(); i++ {
		fieldKey := typeOfValue.Field(i).Name
		field := s.Field(i)
		if !(field.CanInterface()) {
			return nil, fmt.Errorf("BuildQuery: Field %v can't interface", field)
		}
		fieldValue := field.Interface()
		if utils.IsZero(reflect.ValueOf(fieldValue)) {
			continue
		}

		if (fieldKey == "Not") || (fieldKey == "Or") {
			operatorQuery := QueryObject{}
			json.Unmarshal(fieldValue.(json.RawMessage), &operatorQuery)
			operatorSelect, err := orm.BuildQuery(operatorQuery, model)
			if err != nil {
				return nil, fmt.Errorf("BuildQuery: Parsing error in field %v : %v", fieldKey, err)
			}
			switch fieldKey {
			case "Not":
				dbSelect = append(dbSelect, q.Not(operatorSelect...))
			case "Or":
				dbSelect = []q.Matcher{q.Or(append([]q.Matcher{q.And(dbSelect...)}, q.And(operatorSelect...))...)}
			}
			continue
		}
		query, err := orm.ParseQuery(fieldValue.(json.RawMessage), model, fieldKey)
		if err != nil {
			return nil, fmt.Errorf("BuildQuery: Parsing error in field %v : %v", fieldKey, err)
		}
		dbSelect = append(dbSelect, query...)
	}
	return dbSelect, nil
}

// AdvancedBridgeSearch looks up Bridges according to JSON params.
func (orm *ORM) AdvancedBridgeSearch(params interface{}) ([]BridgeType, error) {
	var results []BridgeType
	var model BridgeType
	query, err := orm.BuildQuery(params, &model)
	if err != nil {
		return results, fmt.Errorf("Error building Advanced Bridge query %v", err)
	}
	err = orm.Select(query...).Find(&results)
	return results, err
}

// AdvancedJobRunSearch looks up JobRuns according to JSON params.
func (orm *ORM) AdvancedJobRunSearch(params interface{}) ([]JobRun, error) {
	var results []JobRun
	var model JobRun
	query, err := orm.BuildQuery(params, &model)
	if err != nil {
		return results, fmt.Errorf("Error building Advanced JobRun query %v", err)
	}
	err = orm.Select(query...).Find(&results)
	return results, err
}
