package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/jinzhu/gorm"
	"github.com/smartcontractkit/chainlink/core/store/assets"
	null "gopkg.in/guregu/null.v3"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/mrwonko/cron"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/ugorji/go/codec"
)

var (
	// ErrorCannotMergeNonObject is returned if a Merge was attempted on a string
	// or array JSON value
	ErrorCannotMergeNonObject = errors.New("Cannot merge, expected object '{}'")
)

// RunStatus is a string that represents the run status
type RunStatus string

const (
	// RunStatusUnstarted is the default state of any run status.
	RunStatusUnstarted = RunStatus("")
	// RunStatusInProgress is used for when a run is actively being executed.
	RunStatusInProgress = RunStatus("in_progress")
	// RunStatusPendingConfirmations is used for when a run is awaiting for block confirmations.
	RunStatusPendingConfirmations = RunStatus("pending_confirmations")
	// RunStatusPendingConnection states that the run is waiting on a connection to the block chain.
	RunStatusPendingConnection = RunStatus("pending_connection")
	// RunStatusPendingBridge is used for when a run is waiting on the completion
	// of another event.
	RunStatusPendingBridge = RunStatus("pending_bridge")
	// RunStatusPendingSleep is used for when a run is waiting on a sleep function to finish.
	RunStatusPendingSleep = RunStatus("pending_sleep")
	// RunStatusErrored is used for when a run has errored and will not complete.
	RunStatusErrored = RunStatus("errored")
	// RunStatusCompleted is used for when a run has successfully completed execution.
	RunStatusCompleted = RunStatus("completed")
)

// Unstarted returns true if the status is the initial state.
func (s RunStatus) Unstarted() bool {
	return s == RunStatusUnstarted
}

// PendingBridge returns true if the status is pending_bridge.
func (s RunStatus) PendingBridge() bool {
	return s == RunStatusPendingBridge
}

// PendingConfirmations returns true if the status is pending_confirmations.
func (s RunStatus) PendingConfirmations() bool {
	return s == RunStatusPendingConfirmations
}

// PendingConnection returns true if the status is pending_connection.
func (s RunStatus) PendingConnection() bool {
	return s == RunStatusPendingConnection
}

// PendingSleep returns true if the status is pending_sleep.
func (s RunStatus) PendingSleep() bool {
	return s == RunStatusPendingSleep
}

// Completed returns true if the status is RunStatusCompleted.
func (s RunStatus) Completed() bool {
	return s == RunStatusCompleted
}

// Errored returns true if the status is RunStatusErrored.
func (s RunStatus) Errored() bool {
	return s == RunStatusErrored
}

// Pending returns true if the status is pending external or confirmations.
func (s RunStatus) Pending() bool {
	return s.PendingBridge() || s.PendingConfirmations() || s.PendingSleep() || s.PendingConnection()
}

// Finished returns true if the status is final and can't be changed.
func (s RunStatus) Finished() bool {
	return s.Completed() || s.Errored()
}

// Runnable returns true if the status is ready to be run.
func (s RunStatus) Runnable() bool {
	return !s.Errored() && !s.Pending()
}

// CanStart returns true if the run is ready to begin processed.
func (s RunStatus) CanStart() bool {
	return s.Pending() || s.Unstarted()
}

// Value returns this instance serialized for database storage.
func (s RunStatus) Value() (driver.Value, error) {
	return string(s), nil
}

// Scan reads the database value and returns an instance.
func (s *RunStatus) Scan(value interface{}) error {
	temp, ok := value.(string)
	if !ok {
		return fmt.Errorf("Unable to convert %v of %T to RunStatus", value, value)
	}

	*s = RunStatus(temp)
	return nil
}

// JSON stores the json types string, number, bool, and null.
// Arrays and Objects are returned as their raw json types.
type JSON struct {
	gjson.Result
}

// Value returns this instance serialized for database storage.
func (j JSON) Value() (driver.Value, error) {
	s := j.String()
	if len(s) == 0 {
		return "{}", nil
	}
	return s, nil
}

// Scan reads the database value and returns an instance.
func (j *JSON) Scan(value interface{}) error {
	temp, ok := value.(string)
	if !ok {
		return fmt.Errorf("Unable to convert %v of %T to JSON", value, value)
	}

	*j = JSON{Result: gjson.Parse(temp)}
	return nil
}

// ParseJSON attempts to coerce the input byte array into valid JSON
// and parse it into a JSON object.
func ParseJSON(b []byte) (JSON, error) {
	var j JSON
	str := string(b)
	if len(str) == 0 {
		str = `{}`
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

// Merge combines the given JSON with the existing JSON.
func (j JSON) Merge(j2 JSON) (JSON, error) {
	body := j.Map()
	if body == nil || (j.Type != gjson.JSON && j.Type != gjson.Null) {
		return JSON{}, ErrorCannotMergeNonObject
	}

	for key, value := range j2.Map() {
		body[key] = value
	}

	cleaned := map[string]interface{}{}
	for k, v := range body {
		cleaned[k] = v.Value()
	}

	b, err := json.Marshal(cleaned)
	if err != nil {
		return JSON{}, err
	}

	var rval JSON
	return rval, gjson.Unmarshal(b, &rval)
}

// Empty returns true if the JSON does not exist.
func (j JSON) Empty() bool {
	return !j.Exists()
}

// Bytes returns the raw JSON.
func (j JSON) Bytes() []byte {
	return []byte(j.String())
}

// Add returns a new instance of JSON with the new value added.
func (j JSON) Add(key string, val interface{}) (JSON, error) {
	var j2 JSON
	b, err := json.Marshal(val)
	if err != nil {
		return j2, err
	}
	str := fmt.Sprintf(`{"%v":%v}`, key, string(b))
	if err = json.Unmarshal([]byte(str), &j2); err != nil {
		return j2, err
	}
	return j.Merge(j2)
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
	var b []byte
	cbor := codec.NewEncoderBytes(&b, new(codec.CborHandle))

	switch v := j.Result.Value().(type) {
	case map[string]interface{}, []interface{}, nil:
		return b, cbor.Encode(v)
	default:
		return b, fmt.Errorf("Unable to coerce JSON to CBOR for type %T", v)
	}
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
		return fmt.Errorf("Unable to convert %v of %T to WebURL", value, value)
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

// AnyTimeFromNull returns an AnyTime from a null.Time.
func AnyTimeFromNull(t null.Time) AnyTime {
	return AnyTime{Time: t.Time, Valid: t.Valid}
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
	var n json.Number
	if err := json.Unmarshal(b, &n); err != nil {
		return err
	}

	if len(n) == 0 {
		t.Valid = false
		return nil
	}

	newTime, err := dateparse.ParseAny(n.String())
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
		return fmt.Errorf("Unable to convert %v of %T to Time", value, value)
	}
}

// Cron holds the string that will represent the spec of the cron-job.
// It uses 6 fields to represent the seconds (1), minutes (2), hours (3),
// day of the month (4), month (5), and day of the week (6).
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

	_, err = cron.Parse(s)
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
	Amount             *assets.Eth    `json:"amount"`
}

// CreateKeyRequest represents a request to add an ethereum key.
type CreateKeyRequest struct {
	CurrentPassword string `json:"current_password"`
}

// Big stores large integers and can deserialize a variety of inputs.
type Big big.Int

// NewBig constructs a Big from *big.Int.
func NewBig(i *big.Int) *Big {
	if i != nil {
		b := Big(*i)
		return &b
	}
	return nil
}

// MarshalText marshals this instance to base 10 number as string.
func (b *Big) MarshalText() ([]byte, error) {
	return []byte((*big.Int)(b).Text(10)), nil
}

// MarshalJSON marshals this instance to base 10 number as string.
func (b *Big) MarshalJSON() ([]byte, error) {
	return b.MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (b *Big) UnmarshalText(input []byte) error {
	input = utils.RemoveQuotes(input)
	str := string(input)
	if utils.HasHexPrefix(str) {
		decoded, err := hexutil.DecodeBig(str)
		if err != nil {
			return err
		}
		*b = Big(*decoded)
		return nil
	}

	_, ok := b.setString(str, 10)
	if !ok {
		return fmt.Errorf("Unable to convert %s to Big", str)
	}
	return nil
}

func (b *Big) setString(s string, base int) (*Big, bool) {
	w, ok := (*big.Int)(b).SetString(s, base)
	return (*Big)(w), ok
}

// UnmarshalJSON implements encoding.JSONUnmarshaler.
func (b *Big) UnmarshalJSON(input []byte) error {
	return b.UnmarshalText(input)
}

// Value returns this instance serialized for database storage.
func (b Big) Value() (driver.Value, error) {
	return b.String(), nil
}

// Scan reads the database value and returns an instance.
func (b *Big) Scan(value interface{}) error {
	temp, ok := value.(string)
	if !ok {
		return fmt.Errorf("Unable to convert %v of %T to Big", value, value)
	}

	decoded, ok := b.setString(temp, 10)
	if !ok {
		return fmt.Errorf("Unable to set string %v of %T to base 10 big.Int for Big", value, value)
	}
	*b = *decoded
	return nil
}

// ToInt converts b to a big.Int.
func (b *Big) ToInt() *big.Int {
	return (*big.Int)(b)
}

// String returns the base 10 encoding of b.
func (b *Big) String() string {
	return b.ToInt().Text(10)
}

// Hex returns the hex encoding of b.
func (b *Big) Hex() string {
	return hexutil.EncodeBig(b.ToInt())
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
		return fmt.Errorf("Unable to convert %v of %T to AddressCollection", value, value)
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
	gorm.Model
	Name  string `gorm:"not null;unique;index"`
	Value string `gorm:"not null"`
}
