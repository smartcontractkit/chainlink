// Package internal declares a Raw private key type,
// only available for use in the keystore sub-tree.
package internal

// Raw is a wrapper type that holds private key bytes
// and is designed to prevent accidental logging.
// The only way to access the internal bytes (without reflection) is to use Bytes,
// which is only available to ancestor packages of the parent keystore/ directory.
type Raw struct {
	bytes []byte
}

func NewRaw(b []byte) Raw {
	return Raw{bytes: b}
}

func (raw Raw) String() string {
	return "<Raw Private Key>"
}

func (raw Raw) GoString() string {
	return raw.String()
}

// Bytes is a func for accessing the internal bytes field of Raw.
// It is not declared as a method, because that would allow access from callers which cannot otherwise access this internal package.
func Bytes(raw Raw) []byte { return raw.bytes }

// RawBytes is a helper to use Bytes with keys that have a Raw() method.
func RawBytes(key interface{ Raw() Raw }) []byte { return Bytes(key.Raw()) }
