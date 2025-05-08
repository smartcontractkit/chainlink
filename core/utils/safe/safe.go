package safe

import (
	"fmt"
)

func IntToUint64(n int) (uint64, error) {
	if n < 0 {
		return 0, fmt.Errorf("cannot convert negative int to uint64")
	}

	return uint64(n), nil
}
