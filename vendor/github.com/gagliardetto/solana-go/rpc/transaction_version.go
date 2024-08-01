package rpc

import "strconv"

type TransactionVersion int

const (
	LegacyTransactionVersion TransactionVersion = -1
	legacyVersion                               = `"legacy"`
)

func (a *TransactionVersion) UnmarshalJSON(b []byte) error {
	// Ignore null, like in the main JSON package.
	s := string(b)
	if s == "null" || s == `""` || s == legacyVersion {
		*a = LegacyTransactionVersion
		return nil
	}

	v, err := strconv.Atoi(s)
	if err != nil {
		return err
	}
	*a = TransactionVersion(v)
	return nil
}

func (a TransactionVersion) MarshalJSON() ([]byte, error) {
	if a == LegacyTransactionVersion {
		return []byte(legacyVersion), nil
	} else {
		return []byte(strconv.Itoa(int(a))), nil
	}
}
