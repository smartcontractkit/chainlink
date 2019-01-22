package models

import (
	"crypto/subtle"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jinzhu/gorm"
	"github.com/smartcontractkit/chainlink/store/assets"
	"github.com/smartcontractkit/chainlink/utils"
	null "gopkg.in/guregu/null.v3"
)

// JobSpec is the definition for all the work to be carried out by the node
// for a given contract. It contains the Initiators, Tasks (which are the
// individual steps to be carried out), StartAt, EndAt, and CreatedAt fields.
type JobSpec struct {
	ID         string      `json:"id" gorm:"primary_key;not null"`
	CreatedAt  Time        `json:"createdAt" gorm:"index"`
	Initiators []Initiator `json:"initiators"`
	Tasks      []TaskSpec  `json:"tasks"`
	StartAt    null.Time   `json:"startAt" gorm:"index"`
	EndAt      null.Time   `json:"endAt" gorm:"index"`
}

// JobSpecRequest represents a schema for the incoming job spec request as used by the API.
type JobSpecRequest struct {
	Initiators []Initiator `json:"initiators"`
	Tasks      []TaskSpec  `json:"tasks"`
	StartAt    null.Time   `json:"startAt" gorm:"index"`
	EndAt      null.Time   `json:"endAt" gorm:"index"`
}

// GetID returns the ID of this structure for jsonapi serialization.
func (j JobSpec) GetID() string {
	return j.ID
}

// GetName returns the pluralized "type" of this structure for jsonapi serialization.
func (j JobSpec) GetName() string {
	return "specs"
}

// SetID is used to set the ID of this structure when deserializing from jsonapi documents.
func (j *JobSpec) SetID(value string) error {
	j.ID = value
	return nil
}

// NewJob initializes a new job by generating a unique ID and setting
// the CreatedAt field to the time of invokation.
func NewJob() JobSpec {
	return JobSpec{
		ID:        utils.NewBytes32ID(),
		CreatedAt: Time{Time: time.Now()},
	}
}

// NewJobFromRequest creates a JobSpec from the corresponding parameters in a
// JobSpecRequest
func NewJobFromRequest(jsr JobSpecRequest) JobSpec {
	jobSpec := NewJob()
	jobSpec.Initiators = jsr.Initiators
	jobSpec.Tasks = jsr.Tasks
	jobSpec.EndAt = jsr.EndAt
	jobSpec.StartAt = jsr.StartAt
	return jobSpec
}

// NewRun initializes the job by creating the IDs for the job
// and all associated tasks, and setting the CreatedAt field.
func (j JobSpec) NewRun(i Initiator) JobRun {
	jrid := utils.NewBytes32ID()
	taskRuns := make([]TaskRun, len(j.Tasks))
	for i, task := range j.Tasks {
		trid := utils.NewBytes32ID()
		taskRuns[i] = TaskRun{
			ID:       trid,
			JobRunID: jrid,
			TaskSpec: task,
			Result:   RunResult{CachedTaskRunID: trid, CachedJobRunID: jrid},
		}
	}

	now := time.Now()
	return JobRun{
		ID:          jrid,
		JobSpecID:   j.ID,
		CreatedAt:   now,
		UpdatedAt:   now,
		TaskRuns:    taskRuns,
		Initiator:   i,
		InitiatorID: i.ID,
		Status:      RunStatusUnstarted,
		Result:      RunResult{CachedJobRunID: jrid},
	}
}

// InitiatorsFor returns an array of Initiators for the given list of
// Initiator types.
func (j JobSpec) InitiatorsFor(types ...string) []Initiator {
	list := []Initiator{}
	for _, initr := range j.Initiators {
		for _, t := range types {
			if initr.Type == t {
				list = append(list, initr)
			}
		}
	}
	return list
}

// WebAuthorized returns true if the "web" initiator is present.
func (j JobSpec) WebAuthorized() bool {
	for _, initr := range j.Initiators {
		if initr.Type == InitiatorWeb {
			return true
		}
	}
	return false
}

// IsLogInitiated Returns true if any of the job's initiators are triggered by event logs.
func (j JobSpec) IsLogInitiated() bool {
	for _, initr := range j.Initiators {
		if initr.IsLogInitiated() {
			return true
		}
	}
	return false
}

// Ended returns true if the job has ended.
func (j JobSpec) Ended(t time.Time) bool {
	if !j.EndAt.Valid {
		return false
	}
	return t.After(j.EndAt.Time)
}

// Started returns true if the job has started.
func (j JobSpec) Started(t time.Time) bool {
	if !j.StartAt.Valid {
		return true
	}
	return t.After(j.StartAt.Time) || t.Equal(j.StartAt.Time)
}

// Types of Initiators (see Initiator struct just below.)
const (
	// InitiatorRunLog for tasks in a job to watch an ethereum address
	// and expect a JSON payload from a log event.
	InitiatorRunLog = "runlog"
	// InitiatorCron for tasks in a job to be ran on a schedule.
	InitiatorCron = "cron"
	// InitiatorEthLog for tasks in a job to use the Ethereum blockchain.
	InitiatorEthLog = "ethlog"
	// InitiatorRunAt for tasks in a job to be ran once.
	InitiatorRunAt = "runat"
	// InitiatorWeb for tasks in a job making a web request.
	InitiatorWeb = "web"
	// InitiatorServiceAgreementExecutionLog for tasks in a job to watch a
	// Solidity Coordinator contract and expect a payload from a log event.
	InitiatorServiceAgreementExecutionLog = "execagreement"
)

// Initiator could be thought of as a trigger, defines how a Job can be
// started, or rather, how a JobRun can be created from a Job.
// Initiators will have their own unique ID, but will be associated
// to a parent JobID.
type Initiator struct {
	ID        uint   `json:"id" gorm:"primary_key;auto_increment"`
	JobSpecID string `json:"jobSpecId" gorm:"index"`
	// Type is one of the Initiator* string constants defined just above.
	Type            string `json:"type" gorm:"index;not null"`
	InitiatorParams `json:"params,omitempty"`
	CreatedAt       time.Time `gorm:"index"`
}

// InitiatorParams is a collection of the possible parameters that different
// Initiators may require.
type InitiatorParams struct {
	Schedule   Cron              `json:"schedule,omitempty"`
	Time       Time              `json:"time,omitempty"`
	Ran        bool              `json:"ran,omitempty"`
	Address    common.Address    `json:"address,omitempty" gorm:"index"`
	Requesters AddressCollection `json:"requesters,omitempty" gorm:"type:text"`
}

// UnmarshalJSON parses the raw initiator data and updates the
// initiator as long as the type is valid.
func (i *Initiator) UnmarshalJSON(input []byte) error {
	type Alias Initiator
	var aux Alias
	if err := json.Unmarshal(input, &aux); err != nil {
		return err
	}

	*i = Initiator(aux)
	i.Type = strings.ToLower(aux.Type)
	return nil
}

// IsLogInitiated Returns true if triggered by event logs.
func (i Initiator) IsLogInitiated() bool {
	return i.Type == InitiatorEthLog || i.Type == InitiatorRunLog ||
		i.Type == InitiatorServiceAgreementExecutionLog
}

// TaskSpec is the definition of work to be carried out. The
// Type will be an adapter, and the Params will contain any
// additional information that adapter would need to operate.
type TaskSpec struct {
	gorm.Model
	JobSpecID     string   `json:"-" gorm:"index"`
	Type          TaskType `json:"type" gorm:"index;not null"`
	Confirmations uint64   `json:"confirmations"`
	Params        JSON     `json:"params" gorm:"type:text"`
}

// TaskType defines what Adapter a TaskSpec will use.
type TaskType string

// NewTaskType returns a formatted Task type.
func NewTaskType(val string) (TaskType, error) {
	re := regexp.MustCompile("^[a-zA-Z0-9-_]*$")
	if !re.MatchString(val) {
		return TaskType(""), fmt.Errorf("Task Type validation: name %v contains invalid characters", val)
	}

	return TaskType(strings.ToLower(val)), nil
}

// MustNewTaskType instantiates a new TaskType, and panics if a bad input is provided.
func MustNewTaskType(val string) TaskType {
	tt, err := NewTaskType(val)
	if err != nil {
		panic(fmt.Sprintf("%v is not a valid TaskType", val))
	}
	return tt
}

// UnmarshalJSON converts a bytes slice of JSON to a TaskType.
func (t *TaskType) UnmarshalJSON(input []byte) error {
	var aux string
	if err := json.Unmarshal(input, &aux); err != nil {
		return err
	}
	tt, err := NewTaskType(aux)
	*t = tt
	return err
}

// MarshalJSON converts a TaskType to a JSON byte slice.
func (t TaskType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// String returns this TaskType as a string.
func (t TaskType) String() string {
	return string(t)
}

// Value returns this instance serialized for database storage.
func (t TaskType) Value() (driver.Value, error) {
	return string(t), nil
}

// Scan reads the database value and returns an instance.
func (t *TaskType) Scan(value interface{}) error {
	temp, ok := value.([]uint8)
	if !ok {
		return fmt.Errorf("Unable to convert %v of %T to TaskType", value, value)
	}

	*t = TaskType(temp)
	return nil
}

// BridgeType is used for external adapters and has fields for
// the name of the adapter and its URL.
type BridgeType struct {
	Name                   TaskType    `json:"name" gorm:"unique_index"`
	URL                    WebURL      `json:"url"`
	Confirmations          uint64      `json:"confirmations"`
	IncomingToken          string      `json:"incomingToken"`
	OutgoingToken          string      `json:"outgoingToken"`
	MinimumContractPayment assets.Link `json:"minimumContractPayment" gorm:"type:varchar(255)"`
}

// GetID returns the ID of this structure for jsonapi serialization.
func (bt BridgeType) GetID() string {
	return bt.Name.String()
}

// GetName returns the pluralized "type" of this structure for jsonapi serialization.
func (bt BridgeType) GetName() string {
	return "bridges"
}

// SetID is used to set the ID of this structure when deserializing from jsonapi documents.
func (bt *BridgeType) SetID(value string) error {
	name, err := NewTaskType(value)
	bt.Name = name
	return err
}

// Authenticate returns true if the passed token matches its IncomingToken, or
// returns false with an error.
func (bt BridgeType) Authenticate(token string) (bool, error) {
	var asbytes [utils.NewBytes32Length]byte
	copy(asbytes[:], token)
	if subtle.ConstantTimeCompare(asbytes[:], []byte(bt.IncomingToken)) == 1 {
		return true, nil
	}
	return false, fmt.Errorf("Incorrect access token for %s", bt.Name)
}
