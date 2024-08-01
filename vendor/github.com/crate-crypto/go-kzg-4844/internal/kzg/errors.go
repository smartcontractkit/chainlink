package kzg

import "errors"

var (
	ErrInvalidNumDigests              = errors.New("number of digests is not the same as the number of polynomials")
	ErrInvalidPolynomialSize          = errors.New("invalid polynomial size (larger than SRS or == 0)")
	ErrVerifyOpeningProof             = errors.New("can't verify opening proof")
	ErrPolynomialMismatchedSizeDomain = errors.New("domain size does not equal the number of evaluations in the polynomial")
	ErrMinSRSSize                     = errors.New("minimum srs size is 2")
)
