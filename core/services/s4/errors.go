package s4

import "errors"

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrWrongSignature = errors.New("wrong signature")
	ErrSlotIdTooBig   = errors.New("slot id is too big")
	ErrPayloadTooBig  = errors.New("payload is too big")
	ErrPastExpiration = errors.New("past expiration")
	ErrVersionTooLow  = errors.New("version too low")
)
