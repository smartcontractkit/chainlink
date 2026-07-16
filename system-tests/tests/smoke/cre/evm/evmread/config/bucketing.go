package config

import "fmt"

// ReadBucket identifies a runtime-balanced EVM read test bucket.
type ReadBucket string

const (
	ReadBucketHeavyCalls   ReadBucket = "heavy-calls"
	ReadBucketStateQueries ReadBucket = "state-queries"
	ReadBucketTxArtifacts  ReadBucket = "tx-artifacts"
)

type readBucketDefinition struct {
	Bucket ReadBucket
	Cases  []TestCase
}

// readBucketRegistry is the single place where EVM read test cases are assigned to buckets.
// When adding a new TestCase, add it here and keep bucket runtimes balanced. Best way to do it is by
// executing the tests in CI once and asking an AI to check run details, with execution time and to
// rebalance the buckets so that they are balanced.
var readBucketRegistry = []readBucketDefinition{
	{
		Bucket: ReadBucketHeavyCalls,
		Cases: []TestCase{
			TestCaseEVMReadEvents,
			TestCaseEVMReadCallContract,
		},
	},
	{
		Bucket: ReadBucketStateQueries,
		Cases: []TestCase{
			TestCaseEVMReadBalance,
			TestCaseEVMReadHeaders,
			TestCaseEVMReadNotFoundTx,
		},
	},
	{
		Bucket: ReadBucketTxArtifacts,
		Cases: []TestCase{
			TestCaseEVMReadTransactionReceipt,
			TestCaseEVMReadBTx,
			TestCaseEVMEstimateGas,
		},
	},
}

func CasesForReadBucket(bucket ReadBucket) ([]TestCase, error) {
	for _, bucketDef := range readBucketRegistry {
		if bucketDef.Bucket != bucket {
			continue
		}

		cases := make([]TestCase, len(bucketDef.Cases))
		copy(cases, bucketDef.Cases)
		return cases, nil
	}

	return nil, fmt.Errorf("unknown EVM read bucket %q", bucket)
}

func ValidateReadBucketRegistry() error {
	assignedCases := make(map[TestCase]ReadBucket, TestCaseLen)

	for _, bucketDef := range readBucketRegistry {
		for _, testCase := range bucketDef.Cases {
			if testCase < 0 || testCase >= TestCaseLen {
				return fmt.Errorf("invalid testcase %d in bucket %q", testCase, bucketDef.Bucket)
			}

			if existingBucket, ok := assignedCases[testCase]; ok {
				return fmt.Errorf("testcase %q assigned to multiple buckets: %q and %q", testCase.String(), existingBucket, bucketDef.Bucket)
			}

			assignedCases[testCase] = bucketDef.Bucket
		}
	}

	for testCase := range TestCaseLen {
		if _, ok := assignedCases[testCase]; ok {
			continue
		}

		return fmt.Errorf("testcase %q is not assigned to any EVM read bucket", testCase.String())
	}

	return nil
}
