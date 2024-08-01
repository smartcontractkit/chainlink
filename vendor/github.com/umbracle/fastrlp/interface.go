package fastrlp

// Marshaler is the interface implemented by types that can marshal themselves into valid RLP messages.
type Marshaler interface {
	MarshalRLPTo(dst []byte) ([]byte, error)
	MarshalRLPWith(a *Arena) (*Value, error)
}

// Unmarshaler is the interface implemented by types that can unmarshal a RLP description of themselves
type Unmarshaler interface {
	UnmarshalRLP(buf []byte) error
	UnmarshalRLPWith(v *Value) error
}

// MarshalRLP marshals an RLP object
func MarshalRLP(m Marshaler) ([]byte, error) {
	ar := &Arena{}
	v, err := m.MarshalRLPWith(ar)
	if err != nil {
		return nil, err
	}
	return v.MarshalTo(nil), nil
}

// UnmarshalRLP unmarshals an RLP object
func UnmarshalRLP(buf []byte, m Unmarshaler) error {
	p := &Parser{}
	v, err := p.Parse(buf)
	if err != nil {
		return err
	}
	if err := m.UnmarshalRLPWith(v); err != nil {
		return err
	}
	return nil
}
