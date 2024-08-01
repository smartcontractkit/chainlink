package sqlutil

import "errors"

// BatchFunc is the function to execute on each batch of records, should return the count of records affected
type BatchFunc func(offset, limit uint) (count uint, err error)

const batchSizeErr = "batch size cannot be zero"

// Batch is an iterator for batches of records
func Batch(cb BatchFunc, batchSize uint) error {
	if batchSize == 0 {
		return errors.New(batchSizeErr)
	}

	offset := uint(0)
	limit := batchSize

	for {
		count, err := cb(offset, limit)
		if err != nil {
			return err
		}

		if count < limit {
			return nil
		}

		offset += limit
	}
}
