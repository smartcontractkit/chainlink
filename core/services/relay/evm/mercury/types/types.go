package types

import "fmt"

type FeedID [32]byte

func (f FeedID) String() string {
	return fmt.Sprintf("%x", f[:])
}

func (f FeedID) Hex() string {
	return f.String()
}
