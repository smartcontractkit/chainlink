package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

const (
	FIELD_GEN   int    = 3
	FIELD_PRIME string = "3618502788666131213697322783095070105623107215331596699973092056135872020481"
)

var (
	MaxFelt     = StrToFelt(FIELD_PRIME)
	asciiRegexp = regexp.MustCompile(`^([[:graph:]]|[[:space:]]){1,31}$`)
)

// Felt represents Field Element or Felt from cairo.
type Felt struct {
	*big.Int
}

// Big converts a Felt to its big.Int representation.
func (f *Felt) Big() *big.Int {
	return new(big.Int).SetBytes(f.Int.Bytes())
}

// StrToFelt converts a string containing a decimal, hexadecimal or UTF8 charset into a Felt.
func StrToFelt(str string) *Felt {
	f := new(Felt)
	if ok := f.strToFelt(str); ok {
		return f
	}
	return nil
}

func (f *Felt) strToFelt(str string) bool {
	if b, ok := new(big.Int).SetString(str, 0); ok {
		f.Int = b
		return ok
	}

	// TODO: revisit conversation on seperate 'ShortString' conversion
	if asciiRegexp.MatchString(str) {
		hexStr := hex.EncodeToString([]byte(str))
		if b, ok := new(big.Int).SetString(hexStr, 16); ok {
			f.Int = b
			return ok
		}
	}
	return false
}

// BigToFelt converts a big.Int to its Felt representation.
func BigToFelt(b *big.Int) *Felt {
	return &Felt{Int: b}
}

// BytesToFelt converts a []byte to its Felt representation.
func BytesToFelt(b []byte) *Felt {
	return &Felt{Int: new(big.Int).SetBytes(b)}
}

// String converts a Felt into its 'short string' representation.
func (f *Felt) ShortString() string {
	str := string(f.Bytes())
	if asciiRegexp.MatchString(str) {
		return str
	}
	return ""
}

// String converts a Felt into its hexadecimal string representation and implement fmt.Stringer.
func (f *Felt) String() string {
	return fmt.Sprintf("0x%x", f)
}

// MarshalJSON implements the json Marshaller interface for a Signature array to marshal types to []byte.
func (s Signature) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`["%s","%s"]`, s[0].String(), s[1].String())), nil
}

// MarshalJSON implements the json Marshaller interface for Felt to marshal types to []byte.
func (f Felt) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, f.String())), nil
}

// UnmarshalJSON implements the json Unmarshaller interface to unmarshal []byte into types.
func (f *Felt) UnmarshalJSON(p []byte) error {
	if string(p) == "null" || len(p) == 0 {
		return nil
	}

	var s string
	// parse double quotes
	if p[0] == 0x22 {
		s = string(p[1 : len(p)-1])
	} else {
		s = string(p)
	}

	if ok := f.strToFelt(s); !ok {
		return fmt.Errorf("unmarshalling big int: %s", string(p))
	}

	return nil
}

// MarshalGQL implements the gqlgen Marshaller interface to marshal Felt into an io.Writer.
func (f Felt) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(f.String()))
}

// UnmarshalGQL implements the gqlgen Unmarshaller interface to unmarshal an interface into a Felt.
func (b *Felt) UnmarshalGQL(v interface{}) error {
	switch bi := v.(type) {
	case string:
		if ok := b.strToFelt(bi); ok {
			return nil
		}
	case int:
		b.Int = big.NewInt(int64(bi))
		if b.Int != nil {
			return nil
		}
	}

	return fmt.Errorf("invalid big number")
}

// Value is used by database/sql drivers to store data in databases
func (f Felt) Value() (driver.Value, error) {
	if f.Int == nil {
		return "", nil
	}
	return f.String(), nil
}

// Scan implements the database/sql Scanner interface to read Felt from a databases.
func (f *Felt) Scan(src interface{}) error {
	var i sql.NullString
	if err := i.Scan(src); err != nil {
		return err
	}
	if !i.Valid {
		return nil
	}
	if f.Int == nil {
		f.Int = big.NewInt(0)
	}
	// Value came in a floating point format.
	if strings.ContainsAny(i.String, ".+e") {
		flt := big.NewFloat(0)
		if _, err := fmt.Sscan(i.String, f); err != nil {
			return err
		}
		f.Int, _ = flt.Int(f.Int)
	} else if _, err := fmt.Sscan(i.String, f.Int); err != nil {
		return err
	}
	return nil
}
