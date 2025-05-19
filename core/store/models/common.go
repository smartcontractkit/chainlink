package models

import (
	"bytes"
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"github.com/tidwall/gjson"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
)

// CronParser is the global parser for crontabs.
// It accepts the standard 5 field cron syntax as well as an optional 6th field
// at the front to represent seconds.
var CronParser cron.Parser

func init() {
	cronParserSpec := cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor
	CronParser = cron.NewParser(cronParserSpec)
}

// JSON stores the json types string, number, bool, and null.
// Arrays and Objects are returned as their raw json types.
type JSON struct {
	gjson.Result
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

// Cron holds the string that will represent the spec of the cron-job.
type Cron string

// UnmarshalJSON parses the raw spec stored in JSON-encoded
// data and stores it to the Cron string.
func (c *Cron) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return fmt.Errorf("Cron: %w", err)
	}
	if s == "" {
		return nil
	}

	if !strings.HasPrefix(s, "CRON_TZ=") {
		return errors.New("Cron: specs must specify a time zone using CRON_TZ, e.g. 'CRON_TZ=UTC 5 * * * *'")
	}

	_, err = CronParser.Parse(s)
	if err != nil {
		return fmt.Errorf("Cron: %w", err)
	}
	*c = Cron(s)
	return nil
}

// String returns the current Cron spec string.
func (c Cron) String() string {
	return string(c)
}

// Interval represents a time.Duration stored as a Postgres interval type
type Interval time.Duration

// NewInterval creates Interval for specified duration
func NewInterval(d time.Duration) *Interval {
	i := new(Interval)
	*i = Interval(d)
	return i
}

func (i Interval) Duration() time.Duration {
	return time.Duration(i)
}

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

// SendEtherRequest represents a request to transfer ETH.
type SendEtherRequest struct {
	DestinationAddress common.Address `json:"address"`
	FromAddress        common.Address `json:"from"`
	Amount             assets.Eth     `json:"amount"`
	EVMChainID         *big.Big       `json:"evmChainID"`
	AllowHigherAmounts bool           `json:"allowHigherAmounts"`
	SkipWaitTxAttempt  bool           `json:"skipWaitTxAttempt"`
	WaitAttemptTimeout *time.Duration `json:"waitAttemptTimeout"`
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

var EmptySha256Hash = new(Sha256Hash)

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

func (s *Sha256Hash) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

func (s *Sha256Hash) UnmarshalText(bs []byte) (err error) {
	*s, err = Sha256HashFromHex(string(bs))
	return
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

// ServiceHeader is an HTTP header to include in POST to log service.
type ServiceHeader struct {
	Header string
	Value  string
}

func (h *ServiceHeader) UnmarshalText(input []byte) error {
	parts := strings.SplitN(string(input), ":", 2)
	h.Header = parts[0]
	if len(parts) > 1 {
		h.Value = strings.TrimSpace(parts[1])
	}
	return h.Validate()
}

func (h *ServiceHeader) MarshalText() ([]byte, error) {
	var b bytes.Buffer
	fmt.Fprintf(&b, "%s: %s", h.Header, h.Value)
	return b.Bytes(), nil
}

type ServiceHeaders []ServiceHeader

func (sh *ServiceHeaders) UnmarshalText(input []byte) error {
	if sh == nil {
		return errors.New("Cannot unmarshal to a nil receiver")
	}

	headers := string(input)

	var parsedHeaders []ServiceHeader
	if headers != "" {
		headerLines := strings.Split(headers, "\\")
		for _, header := range headerLines {
			keyValue := strings.Split(header, "||")
			if len(keyValue) != 2 {
				return errors.Errorf("invalid headers provided for the audit logger. Value, single pair split on || required, got: %s", keyValue)
			}
			h := ServiceHeader{
				Header: keyValue[0],
				Value:  keyValue[1],
			}

			if err := h.Validate(); err != nil {
				return err
			}
			parsedHeaders = append(parsedHeaders, h)
		}
	}

	*sh = parsedHeaders
	return nil
}

func (sh *ServiceHeaders) MarshalText() ([]byte, error) {
	if sh == nil {
		return nil, errors.New("Cannot marshal to a nil receiver")
	}

	sb := strings.Builder{}
	for _, header := range *sh {
		sb.WriteString(header.Header)
		sb.WriteString("||")
		sb.WriteString(header.Value)
		sb.WriteString("\\")
	}

	serialized := sb.String()

	if len(serialized) > 0 {
		serialized = serialized[:len(serialized)-1]
	}

	return []byte(serialized), nil
}

// We act slightly more strictly than the HTTP specifications
// technically allow instead following the guidelines of
// cloudflare transforms.
// https://developers.cloudflare.com/rules/transform/request-header-modification/reference/header-format
var (
	headerNameRegex  = regexp.MustCompile(`^[A-Za-z\-]+$`)
	headerValueRegex = regexp.MustCompile("^[A-Za-z_ :;.,\\/\"'?!(){}[\\]@<>=\\-+*#$&`|~^%]+$")
)

func (h ServiceHeader) Validate() (err error) {
	if !headerNameRegex.MatchString(h.Header) {
		err = multierr.Append(err, errors.Errorf("invalid header name: %s", h.Header))
	}

	if !headerValueRegex.MatchString(h.Value) {
		err = multierr.Append(err, errors.Errorf("invalid header value: %s", h.Value))
	}
	return
}
