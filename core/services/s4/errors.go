package s4

import "errors"

var (
	ErrNotFound          = errors.New("not found")
	ErrWrongSignature    = errors.New("wrong signature")
	ErrSlotIdTooBig      = errors.New("slot id is too big")
	ErrPayloadTooBig     = errors.New("payload is too big")
	ErrPastExpiration    = errors.New("past expiration")
	ErrVersionTooLow     = errors.New("version too low")
	ErrExpirationTooLong = errors.New("expiration too long")
)
