package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/assets"
	clnull "github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/utils"
	null "gopkg.in/guregu/null.v4"
)

// JobSpecRequest represents a schema for the incoming job spec request as used by the API.
type JobSpecRequest struct {
	Name       string             `json:"name"`
	Initiators []InitiatorRequest `json:"initiators"`
	Tasks      []TaskSpecRequest  `json:"tasks"`
	StartAt    null.Time          `json:"startAt"`
	EndAt      null.Time          `json:"endAt"`
	MinPayment *assets.Link       `json:"minPayment,omitempty"`
}

// InitiatorRequest represents a schema for incoming initiator requests as used by the API.
type InitiatorRequest struct {
	JobSpecId       JobID  `json:"jobSpecId"`
	Type            string `json:"type"`
	InitiatorParams `json:"params,omitempty"`
}

// TaskSpecRequest represents a schema for incoming TaskSpec requests as used by the API.
type TaskSpecRequest struct {
	Type                             TaskType      `json:"type"`
	MinRequiredIncomingConfirmations clnull.Uint32 `json:"confirmations"`
	Params                           JSON          `json:"params"`
}

// JobSpec is the definition for all the work to be carried out by the node
// for a given contract. It contains the Initiators, Tasks (which are the
// individual steps to be carried out), StartAt, EndAt, and CreatedAt fields.
type JobSpec struct {
	ID         JobID          `json:"id,omitempty" gorm:"primary_key;not null"`
	Name       string         `json:"name"`
	CreatedAt  time.Time      `json:"createdAt" gorm:"index"`
	Initiators []Initiator    `json:"initiators"`
	MinPayment *assets.Link   `json:"minPayment,omitempty"`
	Tasks      []TaskSpec     `json:"tasks"`
	StartAt    null.Time      `json:"startAt" gorm:"index"`
	EndAt      null.Time      `json:"endAt" gorm:"index"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	UpdatedAt  time.Time      `json:"-"`
	Errors     []JobSpecError `json:"-" gorm:"foreignkey:JobSpecID;->"`
}

// GetID returns the ID of this structure for jsonapi serialization.
func (j JobSpec) GetID() string {
	return j.ID.String()
}

// GetName returns the pluralized "type" of this structure for jsonapi serialization.
func (j JobSpec) GetName() string {
	return "specs"
}

// SetID is used to set the ID of this structure when deserializing from jsonapi documents.
func (j *JobSpec) SetID(value string) error {
	return j.ID.UnmarshalText([]byte(value))
}

// NewJob initializes a new job by generating a unique ID and setting
// the CreatedAt field to the time of invokation.
func NewJob() JobSpec {
	return JobSpec{
		ID:        NewJobID(),
		CreatedAt: time.Now(),
	}
}

// NewJobFromRequest creates a JobSpec from the corresponding parameters in a
// JobSpecRequest
func NewJobFromRequest(jsr JobSpecRequest) JobSpec {
	jobSpec := NewJob()
	if jsr.Name != "" {
		jobSpec.Name = jsr.Name
	}
	for _, initr := range jsr.Initiators {
		init := NewInitiatorFromRequest(initr, jobSpec)
		jobSpec.Initiators = append(jobSpec.Initiators, init)
	}
	for _, task := range jsr.Tasks {
		jobSpec.Tasks = append(jobSpec.Tasks, TaskSpec{
			JobSpecID:                        jobSpec.ID,
			Type:                             task.Type,
			MinRequiredIncomingConfirmations: task.MinRequiredIncomingConfirmations,
			Params:                           task.Params,
		})
	}

	jobSpec.EndAt = jsr.EndAt
	jobSpec.StartAt = jsr.StartAt
	jobSpec.MinPayment = jsr.MinPayment
	return jobSpec
}

// Archived returns true if the job spec has been soft deleted
func (j JobSpec) Archived() bool {
	return j.DeletedAt.Valid
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

// InitiatorExternal finds the Job Spec's Initiator field associated with the
// External Initiator's name using a case insensitive search.
//
// Returns nil if not found.
func (j JobSpec) InitiatorExternal(name string) *Initiator {
	for _, i := range j.InitiatorsFor(InitiatorExternal) {
		if strings.EqualFold(i.Name, name) {
			found := i
			return &found
		}
	}
	return nil
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
	// InitiatorExternal for tasks in a job to be trigger by an external party.
	InitiatorExternal = "external"
	// InitiatorFluxMonitor for tasks in a job to be run on price deviation
	// or request for a new round of prices.
	InitiatorFluxMonitor = "fluxmonitor"
	// InitiatorRandomnessLog for tasks from a VRF specific contract
	InitiatorRandomnessLog = "randomnesslog"
)

// Initiator could be thought of as a trigger, defines how a Job can be
// started, or rather, how a JobRun can be created from a Job.
// Initiators will have their own unique ID, but will be associated
// to a parent JobID.
type Initiator struct {
	ID        int64 `json:"id" gorm:"primary_key;auto_increment"`
	JobSpecID JobID `json:"jobSpecId"`

	// Type is one of the Initiator* string constants defined just above.
	Type            string    `json:"type" gorm:"index;not null"`
	CreatedAt       time.Time `json:"createdAt" gorm:"index"`
	InitiatorParams `json:"params,omitempty"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
	UpdatedAt       time.Time      `json:"-"`
}

// InitiatorParams is a collection of the possible parameters that different
// Initiators may require.
type InitiatorParams struct {
	// Common parameters
	Address common.Address `json:"address,omitempty" gorm:"index"`
	Name    string         `json:"name,omitempty"`

	// Cron parameters
	Schedule Cron `json:"schedule,omitempty"`

	// RunAt parameters.
	Time AnyTime `json:"time,omitempty"`
	Ran  bool    `json:"ran,omitempty"`

	// External initiator job parameters.
	Body *JSON `json:"body,omitempty" gorm:"column:params"`

	// Log specific job parameters.
	Requesters AddressCollection `json:"requesters,omitempty" gorm:"type:text"`
	FromBlock  *utils.Big        `json:"fromBlock,omitempty" gorm:"type:varchar(255)"`
	ToBlock    *utils.Big        `json:"toBlock,omitempty" gorm:"type:varchar(255)"`
	Topics     Topics            `json:"topics,omitempty"`
	// JobIDTopicFilter, if present, is used in addition to the job's actual ID when filtering
	// initiator logs
	JobIDTopicFilter JobID `json:"jobIDTopicFilter,omitempty"`

	// Flux monitior specific parameters.
	RequestData JSON    `json:"requestData,omitempty" gorm:"type:text"`
	Feeds       Feeds   `json:"feeds,omitempty" gorm:"type:text"`
	Precision   int32   `json:"precision,omitempty" gorm:"type:smallint"`
	Threshold   float32 `json:"threshold,omitempty"`
	// AbsoluteThreshold is the maximum absolute change allowed in a fluxmonitored
	// value before a new round should be kicked off, so that the current value
	// can be reported on-chain.
	AbsoluteThreshold float32         `json:"absoluteThreshold" gorm:"type:float;not null"`
	PollTimer         PollTimerConfig `json:"pollTimer,omitempty" gorm:"type:jsonb"`
	IdleTimer         IdleTimerConfig `json:"idleTimer,omitempty" gorm:"type:jsonb"`
}

type PollTimerConfig struct {
	Disabled bool     `json:"disabled,omitempty"`
	Period   Duration `json:"period,omitempty"`
}

// Value is defined so that we can store PollTimerConfig as JSONB, because
// of an error with GORM where it has trouble with nested structs as JSONB.
// See https://github.com/jinzhu/gorm/issues/2704
func (ptc PollTimerConfig) Value() (driver.Value, error) {
	b, err := json.Marshal(ptc)
	if err != nil {
		return nil, err
	}
	return b, err
}

// Scan is defined so that we can read PollTimerConfig as JSONB, because
// of an error with GORM where it has trouble with nested structs as JSONB.
// See https://github.com/jinzhu/gorm/issues/2704
func (ptc *PollTimerConfig) Scan(value interface{}) error {
	if value == nil {
		*ptc = PollTimerConfig{}
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("invalid Scan Source")
	}
	return json.Unmarshal(b, ptc)
}

type IdleTimerConfig struct {
	Disabled bool     `json:"disabled,omitempty"`
	Duration Duration `json:"duration,omitempty"`
}

// Value is defined so that we can store IdleTimerConfig as JSONB, because
// of an error with GORM where it has trouble with nested structs as JSONB.
// See https://github.com/jinzhu/gorm/issues/2704
func (itc IdleTimerConfig) Value() (driver.Value, error) {
	b, err := json.Marshal(itc)
	if err != nil {
		return nil, err
	}
	return b, err
}

// Scan is defined so that we can read IdleTimerConfig as JSONB, because
// of an error with GORM where it has trouble with nested structs as JSONB.
// See https://github.com/jinzhu/gorm/issues/2704
func (itc *IdleTimerConfig) Scan(value interface{}) error {
	if value == nil {
		*itc = IdleTimerConfig{}
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("invalid Scan Source")
	}
	return json.Unmarshal(b, itc)
}

// Topics handle the serialization of ethereum log topics to and from the data store.
type Topics [][]common.Hash

// Scan coerces the value returned from the data store to the proper data
// in this instance.
func (t *Topics) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		err := json.Unmarshal(v, &t)
		if err != nil {
			return errors.Wrapf(err, "Unable to convert %v of %T to Topics", value, value)
		}
		return nil
	default:
		return fmt.Errorf("unable to convert %v of %T to Topics", value, value)
	}
}

// Value returns this instance serialized for database storage.
func (t Topics) Value() (driver.Value, error) {
	j, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	return string(j), nil
}

// NewInitiatorFromRequest creates an Initiator from the corresponding
// parameters in a InitiatorRequest
func NewInitiatorFromRequest(
	initr InitiatorRequest,
	_ JobSpec,
) Initiator {
	ret := Initiator{
		JobSpecID: initr.JobSpecId,
		// Type must be downcast to comply with Initiator
		// deserialization logic. Ideally, Initiator.Type should be its
		// own type (InitiatorType) that handles deserialization
		// validation.
		Type:            strings.ToLower(initr.Type),
		InitiatorParams: initr.InitiatorParams,
	}
	return ret
}

// IsLogInitiated Returns true if triggered by event logs.
func (i Initiator) IsLogInitiated() bool {
	for _, logType := range LogBasedChainlinkJobInitiators {
		if i.Type == logType {
			return true
		}
	}
	return false
}

// Feeds holds the json of the feeds parameter in the job spec. It is an array of
// URL strings and/or objects containing the names of bridges
type Feeds = JSON

// TaskSpec is the definition of work to be carried out. The
// Type will be an adapter, and the Params will contain any
// additional information that adapter would need to operate.
type TaskSpec struct {
	ID                               int64         `gorm:"primary_key"`
	JobSpecID                        JobID         `json:"jobSpecId"`
	Type                             TaskType      `json:"type" gorm:"index;not null"`
	MinRequiredIncomingConfirmations clnull.Uint32 `json:"confirmations" gorm:"column:confirmations"`
	Params                           JSON          `json:"params" gorm:"type:text"`
	CreatedAt                        time.Time
	UpdatedAt                        time.Time
	DeletedAt                        gorm.DeletedAt
}

// TaskType defines what Adapter a TaskSpec will use.
type TaskType string

// NewTaskType returns a formatted Task type.
func NewTaskType(val string) (TaskType, error) {
	re := regexp.MustCompile("^[a-zA-Z0-9-_]*$")
	if !re.MatchString(val) {
		return TaskType(""), fmt.Errorf("task type validation: name %v contains invalid characters", val)
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
	temp, ok := value.(string)
	if !ok {
		return fmt.Errorf("unable to convert %v of %T to TaskType", value, value)
	}

	*t = TaskType(temp)
	return nil
}
