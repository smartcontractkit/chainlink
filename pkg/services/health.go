package services

import (
	"errors"
	"testing"
)

// CopyHealth copies health statuses from src to dest. Useful when implementing Service.HealthReport.
// If duplicate names are encountered, the errors are joined, unless testing in which case a panic is thrown.
func CopyHealth(dest, src map[string]error) {
	for name, err := range src {
		errOrig, ok := dest[name]
		if ok {
			if testing.Testing() {
				panic("service names must be unique: duplicate name: " + name)
			}
			if errOrig != nil {
				dest[name] = errors.Join(errOrig, err)
				continue
			}
		}
		dest[name] = err
	}
}
