package gokzg4844

import "errors"

var (
	ErrBatchLengthCheck   = errors.New("the number of blobs, commitments, and proofs must be the same")
	ErrNonCanonicalScalar = errors.New("scalar is not canonical when interpreted as a big integer in big-endian")
)
