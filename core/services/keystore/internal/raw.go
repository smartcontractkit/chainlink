// Package internal declares a Raw private key type,
// only available for use in the keystore sub-tree.
package internal

// Raw is a wrapper type that holds private key bytes
// and is designed to prevent accidental logging.
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

func (raw Raw) Bytes() []byte {
	return raw.bytes
}
