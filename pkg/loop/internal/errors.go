package internal

import (
	"fmt"
	"math"
)

type ErrConnAccept struct {
	ID   uint32
	Name string
	Err  error
}

func (e ErrConnAccept) Error() string {
	return fmt.Sprintf("failed to accept %s server connection %d: %s", e.Name, e.ID, e.Err)
}

func (e ErrConnAccept) Unwrap() error {
	return e.Err
}

type ErrConnDial struct {
	ID   uint32
	Name string
	Err  error
}

func (e ErrConnDial) Error() string {
	return fmt.Sprintf("failed to dial %s client connection %d: %s", e.Name, e.ID, e.Err)
}

func (e ErrConnDial) Unwrap() error {
	return e.Err
}

type ErrConfigDigestLen int

func (e ErrConfigDigestLen) Error() string {
	return fmt.Sprintf("invalid ConfigDigest len %d: must be 32", e)
}

type ErrUint8Bounds struct {
	U    uint32
	Name string
}

func (e ErrUint8Bounds) Error() string {
	return fmt.Sprintf("expected uint8 %s (max %d) but got %d", e.Name, math.MaxUint8, e.U)
}
