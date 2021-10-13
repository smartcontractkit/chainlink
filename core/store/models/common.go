package models

import (
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"gorm.io/gorm/schema"

	"gorm.io/gorm"

	"github.com/araddon/dateparse"
	"github.com/ethereum/go-ethereum/common"
	"github.com/fxamacker/cbor/v2"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// CronParser is the global parser for crontabs.
// It accepts the standard 5 field cron syntax as well as an optional 6th field
// at the front to represent seconds.
var CronParser cron.Parser

func init() {
	cronParserSpec := cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor
	CronParser = cron.NewParser(cronParserSpec)
}

// RunStatus is a string that represents the run status
type RunStatus string

const (
	// RunStatusUnstarted is the default state of any run status.
	RunStatusUnstarted = RunStatus("unstarted")
	// RunStatusInProgress is used for when a run is actively being executed.
	RunStatusInProgress = RunStatus("in_progress")
	// RunStatusPendingIncomingConfirmations is used for when a run is awaiting for incoming block confirmations
	// e.g. waiting for the log event to be N blocks deep
	RunStatusPendingIncomingConfirmations = RunStatus("pending_incoming_confirmations")
	// RunStatusPendingConnection states that the run is waiting on a connection to the block chain.
	RunStatusPendingConnection = RunStatus("pending_connection")
	// RunStatusPendingBridge is used for when a run is waiting on the completion
	// of another event.
	RunStatusPendingBridge = RunStatus("pending_bridge")
	// RunStatusPendingSleep is used for when a run is waiting on a sleep function to finish.
	RunStatusPendingSleep = RunStatus("pending_sleep")
	// RunStatusPendingOutgoingConfirmations is used for when a run is waiting for outgoing block confirmations
	// e.g. we have sent a transaction using ethtx and are now waiting for it to be N blocks deep
	RunStatusPendingOutgoingConfirmations = RunStatus("pending_outgoing_confirmations")
	// RunStatusErrored is used for when a run has errored and will not complete.
	RunStatusErrored = RunStatus("errored")
	// RunStatusCompleted is used for when a run has successfully completed execution.
	RunStatusCompleted = RunStatus("completed")
	// RunStatusCancelled is used to indicate a run is no longer desired.
	RunStatusCancelled = RunStatus("cancelled")
)

// Unstarted returns true if the status is the initial state.
func (s RunStatus) Unstarted() bool {
	return s == RunStatusUnstarted
}

// PendingBridge returns true if the status is pending_bridge.
func (s RunStatus) PendingBridge() bool {
	return s == RunStatusPendingBridge
}

// PendingIncomingConfirmations returns true if the status is pending_incoming_confirmations.
func (s RunStatus) PendingIncomingConfirmations() bool {
	return s == RunStatusPendingIncomingConfirmations
}

// PendingConnection returns true if the status is pending_connection.
func (s RunStatus) PendingConnection() bool {
	return s == RunStatusPendingConnection
}

// PendingSleep returns true if the status is pending_sleep.
func (s RunStatus) PendingSleep() bool {
	return s == RunStatusPendingSleep
}

// PendingOutgoingConfirmations returns true if the status is pending_incoming_confirmations.
func (s RunStatus) PendingOutgoingConfirmations() bool {
	return s == RunStatusPendingOutgoingConfirmations
}

// Completed returns true if the status is RunStatusCompleted.
func (s RunStatus) Completed() bool {
	return s == RunStatusCompleted
}

// Cancelled returns true if the status is RunStatusCancelled.
func (s RunStatus) Cancelled() bool {
	return s == RunStatusCancelled
}

// Errored returns true if the status is RunStatusErrored.
func (s RunStatus) Errored() bool {
	return s == RunStatusErrored
}

// Pending returns true if the status is pending external or confirmations.
func (s RunStatus) Pending() bool {
	return s.PendingBridge() || s.PendingIncomingConfirmations() || s.PendingOutgoingConfirmations() || s.PendingSleep() || s.PendingConnection()
}

// Finished returns true if the status is final and can't be changed.
func (s RunStatus) Finished() bool {
	return s.Completed() || s.Errored() || s.Cancelled()
}

// Runnable returns true if the status is ready to be run.
func (s RunStatus) Runnable() bool {
	return !s.Errored() && !s.Pending()
}

// Value returns this instance serialized for database storage.
func (s RunStatus) Value() (driver.Value, error) {
	return string(s), nil
}

// Scan reads the database value and returns an instance.
func (s *RunStatus) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*s = RunStatus(string(v))
	case string:
		*s = RunStatus(v)
	default:
		return fmt.Errorf("unable to convert %#v of %T to RunStatus", value, value)
	}
	return nil
}

const (
	ResultKey           = "result"
	ResultCollectionKey = "__chainlink_result_collection__"
)

// JSON stores the json types string, number, bool, and null.
// Arrays and Objects are returned as their raw json types.
type JSON struct {
	gjson.Result
}

func (JSON) GormDataType() string {
	return "json"
}

// GormDBDataType gorm db data type
func (JSON) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "postgres":
		return "JSONB"
	}
	return ""
}

// Value returns this instance serialized for database storage.
func (j JSON) Value() (driver.Value, error) {
	s := j.Bytes()
	if len(s) == 0 {
		return nil, nil
	}
	return s, nil
}

// Scan reads the database value and returns an instance.
func (j *JSON) Scan(value interface{}) error {
	switch v := value.(type) {
	case string:
		*j = JSON{Result: gjson.Parse(v)}
	case []byte:
		*j = JSON{Result: gjson.ParseBytes(v)}
	default:
		return fmt.Errorf("unable to convert %v of %T to JSON", value, value)
	}
	return nil
}

func MustParseJSON(b []byte) JSON {
	var j JSON
	str := string(b)
	if len(str) == 0 {
		panic("empty byte array")
	}
	if err := json.Unmarshal([]byte(str), &j); err != nil {
		panic(err)
	}
	return j
}

// ParseJSON attempts to coerce the input byte array into valid JSON
// and parse it into a JSON object.
func ParseJSON(b []byte) (JSON, error) {
	var j JSON
	str := string(b)
	if len(str) == 0 {
		return j, nil
	}
	return j, json.Unmarshal([]byte(str), &j)
}

// UnmarshalJSON parses the JSON bytes and stores in the *JSON pointer.
func (j *JSON) UnmarshalJSON(b []byte) error {
	str := string(b)
	if !gjson.Valid(str) {
		return fmt.Errorf("invalid JSON: %v", str)
	}
	*j = JSON{gjson.Parse(str)}
	return nil
}

// MarshalJSON returns the JSON data if it already exists, returns
// an empty JSON object as bytes if not.
func (j JSON) MarshalJSON() ([]byte, error) {
	if j.Exists() {
		return j.Bytes(), nil
	}
	return []byte("{}"), nil
}

func (j *JSON) UnmarshalTOML(val interface{}) error {
	var bs []byte
	switch v := val.(type) {
	case string:
		bs = []byte(v)
	case []byte:
		bs = v
	}
	var err error
	*j, err = ParseJSON(bs)
	return err
}

// Bytes returns the raw JSON.
func (j JSON) Bytes() []byte {
	if len(j.String()) == 0 {
		return nil
	}
	return []byte(j.String())
}

// AsMap returns j as a map
func (j JSON) AsMap() (map[string]interface{}, error) {
	output := make(map[string]interface{})
	switch v := j.Result.Value().(type) {
	case map[string]interface{}:
		for key, value := range v {
			output[key] = value
		}
	case nil:
	default:
		return nil, errors.New("can only add to JSON objects or null")
	}
	return output, nil
}

// mapToJSON returns m as a JSON object, or errors
func mapToJSON(m map[string]interface{}) (JSON, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return JSON{}, err
	}
	return JSON{Result: gjson.ParseBytes(bytes)}, nil
}

// Add returns a new instance of JSON with the new value added.
func (j JSON) Add(insertKey string, insertValue interface{}) (JSON, error) {
	return j.MultiAdd(KV{insertKey: insertValue})
}

func (j JSON) PrependAtArrayKey(insertKey string, insertValue interface{}) (JSON, error) {
	curr := j.Get(insertKey).Array()
	updated := make([]interface{}, 0)
	updated = append(updated, insertValue)
	for _, c := range curr {
		updated = append(updated, c.Value())
	}
	return j.Add(insertKey, updated)
}

// KV represents a key/value pair to be added to a JSON object
type KV map[string]interface{}

// MultiAdd returns a new instance of j with the new values added.
func (j JSON) MultiAdd(keyValues KV) (JSON, error) {
	output, err := j.AsMap()
	if err != nil {
		return JSON{}, err
	}
	for key, value := range keyValues {
		output[key] = value
	}
	return mapToJSON(output)
}

// Delete returns a new instance of JSON with the specified key removed.
func (j JSON) Delete(key string) (JSON, error) {
	js, err := sjson.Delete(j.String(), key)
	if err != nil {
		return j, err
	}
	return ParseJSON([]byte(js))
}

// CBOR returns a bytes array of the JSON map or array encoded to CBOR.
func (j JSON) CBOR() ([]byte, error) {
	switch v := j.Result.Value().(type) {
	case map[string]interface{}, []interface{}, nil:
		return cbor.Marshal(v)
	default:
		var b []byte
		return b, fmt.Errorf("unable to coerce JSON to CBOR for type %T", v)
	}
}

// MarshalToMap converts a struct (typically) to a map[string] so it can be
// manipulated without repeatedly serializing/deserializing
func MarshalToMap(input interface{}) (map[string]interface{}, error) {
	bytes, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	var output map[string]interface{}
	err = json.Unmarshal(bytes, &output)
	if err != nil {
		// Technically this should be impossible
		return nil, err
	}
	return output, nil
}

// WebURL contains the URL of the endpoint.
type WebURL url.URL

// UnmarshalJSON parses the raw URL stored in JSON-encoded
// data to a URL structure and sets it to the URL field.
func (w *WebURL) UnmarshalJSON(j []byte) error {
	var v string
	err := json.Unmarshal(j, &v)
	if err != nil {
		return err
	}
	// handle no url case
	if len(v) == 0 {
		return nil
	}

	u, err := url.ParseRequestURI(v)
	if err != nil {
		return err
	}
	*w = WebURL(*u)
	return nil
}

// MarshalJSON returns the JSON-encoded string of the given data.
func (w WebURL) MarshalJSON() ([]byte, error) {
	return json.Marshal(w.String())
}

// String delegates to the wrapped URL struct or an empty string when it is nil
func (w WebURL) String() string {
	url := url.URL(w)
	return url.String()
}

// Value returns this instance serialized for database storage.
func (w WebURL) Value() (driver.Value, error) {
	return w.String(), nil
}

// Scan reads the database value and returns an instance.
func (w *WebURL) Scan(value interface{}) error {
	s, ok := value.(string)
	if !ok {
		return fmt.Errorf("unable to convert %v of %T to WebURL", value, value)
	}

	u, err := url.ParseRequestURI(s)
	if err != nil {
		return err
	}
	*w = WebURL(*u)
	return nil
}

// AnyTime holds a common field for time, and serializes it as
// a json number.
type AnyTime struct {
	time.Time
	Valid bool
}

// NewAnyTime creates a new Time.
func NewAnyTime(t time.Time) AnyTime {
	return AnyTime{Time: t, Valid: true}
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this time is null.
func (t AnyTime) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return t.Time.UTC().MarshalJSON()
}

// UnmarshalJSON parses the raw time stored in JSON-encoded
// data and stores it to the Time field.
func (t *AnyTime) UnmarshalJSON(b []byte) error {
	var str string

	var n json.Number
	if err := json.Unmarshal(b, &n); err == nil {
		str = n.String()
	} else if err := json.Unmarshal(b, &str); err != nil {
		return err
	}

	if len(str) == 0 {
		t.Valid = false
		return nil
	}

	newTime, err := dateparse.ParseAny(str)
	t.Time = newTime.UTC()
	t.Valid = true
	return err
}

// MarshalText returns null if not set, or the time.
func (t AnyTime) MarshalText() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return t.Time.MarshalText()
}

// UnmarshalText parses null or a valid time.
func (t *AnyTime) UnmarshalText(text []byte) error {
	str := string(text)
	if str == "" || str == "null" {
		t.Valid = false
		return nil
	}
	if err := t.Time.UnmarshalText(text); err != nil {
		return err
	}
	t.Valid = true
	return nil
}

// Value returns this instance serialized for database storage.
func (t AnyTime) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time, nil
}

// Scan reads the database value and returns an instance.
func (t *AnyTime) Scan(value interface{}) error {
	switch temp := value.(type) {
	case time.Time:
		t.Time = temp.UTC()
		t.Valid = true
		return nil
	case nil:
		t.Valid = false
		return nil
	default:
		return fmt.Errorf("unable to convert %v of %T to Time", value, value)
	}
}

// Cron holds the string that will represent the spec of the cron-job.
type Cron string

// UnmarshalJSON parses the raw spec stored in JSON-encoded
// data and stores it to the Cron string.
func (c *Cron) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return fmt.Errorf("Cron: %v", err)
	}
	if s == "" {
		return nil
	}

	if !strings.HasPrefix(s, "CRON_TZ=") {
		return errors.New("Cron: specs must specify a time zone using CRON_TZ, e.g. 'CRON_TZ=UTC 5 * * * *'")
	}

	_, err = CronParser.Parse(s)
	if err != nil {
		return fmt.Errorf("Cron: %v", err)
	}
	*c = Cron(s)
	return nil
}

// String returns the current Cron spec string.
func (c Cron) String() string {
	return string(c)
}

// Duration is a non-negative time duration.
type Duration struct{ d time.Duration }

func MakeDuration(d time.Duration) (Duration, error) {
	if d < time.Duration(0) {
		return Duration{}, fmt.Errorf("cannot make negative time duration: %s", d)
	}
	return Duration{d: d}, nil
}

func MustMakeDuration(d time.Duration) Duration {
	rv, err := MakeDuration(d)
	if err != nil {
		panic(err)
	}
	return rv
}

// Duration returns the value as the standard time.Duration value.
func (d Duration) Duration() time.Duration {
	return d.d
}

// Before returns the time d units before time t
func (d Duration) Before(t time.Time) time.Time {
	return t.Add(-d.Duration())
}

// Shorter returns true if and only if d is shorter than od.
func (d Duration) Shorter(od Duration) bool { return d.d < od.d }

// IsInstant is true if and only if d is of duration 0
func (d Duration) IsInstant() bool { return d.d == 0 }

// String returns a string representing the duration in the form "72h3m0.5s".
// Leading zero units are omitted. As a special case, durations less than one
// second format use a smaller unit (milli-, micro-, or nanoseconds) to ensure
// that the leading digit is non-zero. The zero duration formats as 0s.
func (d Duration) String() string {
	return d.Duration().String()
}

// MarshalJSON implements the json.Marshaler interface.
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *Duration) UnmarshalJSON(input []byte) error {
	var txt string
	err := json.Unmarshal(input, &txt)
	if err != nil {
		return err
	}
	v, err := time.ParseDuration(string(txt))
	if err != nil {
		return err
	}
	*d, err = MakeDuration(v)
	if err != nil {
		return err
	}
	return nil
}

func (d *Duration) Scan(v interface{}) (err error) {
	switch tv := v.(type) {
	case int64:
		*d, err = MakeDuration(time.Duration(tv))
		return err
	default:
		return errors.Errorf(`don't know how to parse "%s" of type %T as a `+
			`models.Duration`, tv, tv)
	}
}

func (d Duration) Value() (driver.Value, error) {
	return int64(d.d), nil
}

// Interval represents a time.Duration stored as a Postgres interval type
type Interval time.Duration

// MarshalText implements the text.Marshaler interface.
func (i Interval) MarshalText() ([]byte, error) {
	return []byte(time.Duration(i).String()), nil
}

// UnmarshalText implements the text.Unmarshaler interface.
func (i *Interval) UnmarshalText(input []byte) error {
	v, err := time.ParseDuration(string(input))
	if err != nil {
		return err
	}
	*i = Interval(v)
	return nil
}

func (i *Interval) Scan(v interface{}) error {
	if v == nil {
		*i = Interval(time.Duration(0))
		return nil
	}
	asInt64, is := v.(int64)
	if !is {
		return errors.Errorf("models.Interval#Scan() wanted int64, got %T", v)
	}
	*i = Interval(time.Duration(asInt64) * time.Nanosecond)
	return nil
}

func (i Interval) Value() (driver.Value, error) {
	return time.Duration(i).Nanoseconds(), nil
}

func (i Interval) IsZero() bool {
	return time.Duration(i) == time.Duration(0)
}

// WithdrawalRequest request to withdraw LINK.
type WithdrawalRequest struct {
	DestinationAddress common.Address `json:"address"`
	ContractAddress    common.Address `json:"contractAddress"`
	Amount             *assets.Link   `json:"amount"`
}

// SendEtherRequest represents a request to transfer ETH.
type SendEtherRequest struct {
	DestinationAddress common.Address `json:"address"`
	FromAddress        common.Address `json:"from"`
	Amount             assets.Eth     `json:"amount"`
}

// AddressCollection is an array of common.Address
// serializable to and from a database.
type AddressCollection []common.Address

// ToStrings returns this address collection as an array of strings.
func (r AddressCollection) ToStrings() []string {
	// Unable to convert copy-free without unsafe:
	// https://stackoverflow.com/a/48554123/639773
	converted := make([]string, len(r))
	for i, e := range r {
		converted[i] = e.Hex()
	}
	return converted
}

// Value returns the string value to be written to the database.
func (r AddressCollection) Value() (driver.Value, error) {
	return strings.Join(r.ToStrings(), ","), nil
}

// Scan parses the database value as a string.
func (r *AddressCollection) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("unable to convert %v of %T to AddressCollection", value, value)
	}

	if len(str) == 0 {
		return nil
	}

	arr := strings.Split(str, ",")
	collection := make(AddressCollection, len(arr))
	for i, a := range arr {
		collection[i] = common.HexToAddress(a)
	}
	*r = collection
	return nil
}

// Configuration stores key value pairs for overriding global configuration
type Configuration struct {
	ID        int64  `gorm:"primary_key"`
	Name      string `gorm:"not null;unique;index"`
	Value     string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *gorm.DeletedAt
}

// MergeExceptResult does a merge, but will never clobber the field called "result"
// On conflicting keys, rightmost inputs will clobber leftmost inputs EXCEPT if the field is named "result", in which case the leftmost result wins
// This is needed to work around idiosyncrasies in the V1 job pipeline where "result" has special meaning
func MergeExceptResult(inputs ...JSON) (JSON, error) {
	output := make(map[string]interface{})

	for _, input := range inputs {
		switch v := input.Result.Value().(type) {
		case map[string]interface{}:
			for key, value := range v {
				if key == "result" {
					if _, exists := output["result"]; exists {
						// Do not overwrite result field
						continue
					}
				}
				output[key] = value
			}
		case nil:
		default:
			return JSON{}, errors.New("can only merge JSON objects")
		}
	}

	bytes, err := json.Marshal(output)
	if err != nil {
		return JSON{}, err
	}

	return JSON{Result: gjson.ParseBytes(bytes)}, nil
}

// Merge returns a new map with all keys merged from left to right
// On conflicting keys, rightmost inputs will clobber leftmost inputs
func Merge(inputs ...JSON) (JSON, error) {
	output := make(map[string]interface{})

	for _, input := range inputs {
		switch v := input.Result.Value().(type) {
		case map[string]interface{}:
			for key, value := range v {
				output[key] = value
			}
		case nil:
		default:
			return JSON{}, errors.New("can only merge JSON objects")
		}
	}

	bytes, err := json.Marshal(output)
	if err != nil {
		return JSON{}, err
	}

	return JSON{Result: gjson.ParseBytes(bytes)}, nil
}

// Explicit type indicating a 32-byte sha256 hash
type Sha256Hash [32]byte

// MarshalJSON converts a Sha256Hash to a JSON byte slice.
func (s Sha256Hash) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// UnmarshalJSON converts a bytes slice of JSON to a TaskType.
func (s *Sha256Hash) UnmarshalJSON(input []byte) error {
	var shaHash string
	if err := json.Unmarshal(input, &shaHash); err != nil {
		return err
	}

	sha, err := Sha256HashFromHex(shaHash)
	if err != nil {
		return err
	}

	*s = sha
	return nil
}

func Sha256HashFromHex(x string) (Sha256Hash, error) {
	bs, err := hex.DecodeString(x)
	if err != nil {
		return Sha256Hash{}, err
	}
	var hash Sha256Hash
	copy(hash[:], bs)
	return hash, nil
}

func MustSha256HashFromHex(x string) Sha256Hash {
	bs, err := hex.DecodeString(x)
	if err != nil {
		panic(err)
	}
	var hash Sha256Hash
	copy(hash[:], bs)
	return hash
}

func (s Sha256Hash) String() string {
	return hex.EncodeToString(s[:])
}

func (s *Sha256Hash) UnmarshalText(bs []byte) error {
	x, err := hex.DecodeString(string(bs))
	if err != nil {
		return err
	}
	*s = Sha256Hash{}
	copy((*s)[:], x)
	return nil
}

func (s *Sha256Hash) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.Errorf("Failed to unmarshal Sha256Hash value: %v", value)
	}
	if s == nil {
		*s = Sha256Hash{}
	}
	copy((*s)[:], bytes)
	return nil
}

func (s Sha256Hash) Value() (driver.Value, error) {
	b := make([]byte, 32)
	copy(b, s[:])
	return b, nil
}
