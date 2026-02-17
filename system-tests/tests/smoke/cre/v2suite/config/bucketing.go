package config

import "fmt"

type SuiteScenario int

const (
	SuiteScenarioProofOfReserve SuiteScenario = iota
	SuiteScenarioVaultDON
	SuiteScenarioCronBeholder
	SuiteScenarioHTTPTriggerAction
	SuiteScenarioHTTPActionCRUD
	SuiteScenarioDONTime
	SuiteScenarioConsensus
	SuiteScenarioLen
)

func (s SuiteScenario) String() string {
	switch s {
	case SuiteScenarioProofOfReserve:
		return "ProofOfReserve"
	case SuiteScenarioVaultDON:
		return "VaultDON"
	case SuiteScenarioCronBeholder:
		return "CronBeholder"
	case SuiteScenarioHTTPTriggerAction:
		return "HTTPTriggerAction"
	case SuiteScenarioHTTPActionCRUD:
		return "HTTPActionCRUD"
	case SuiteScenarioDONTime:
		return "DONTime"
	case SuiteScenarioConsensus:
		return "Consensus"
	default:
		return fmt.Sprintf("unknown SuiteScenario: %d", s)
	}
}

// SuiteBucket identifies a runtime-balanced bucket for the old V2 suite scenarios.
type SuiteBucket string

const (
	SuiteBucketA SuiteBucket = "suite-bucket-a"
	SuiteBucketB SuiteBucket = "suite-bucket-b"
	SuiteBucketC SuiteBucket = "suite-bucket-c"
)

type suiteBucketDefinition struct {
	Bucket    SuiteBucket
	Scenarios []SuiteScenario
}

// suiteBucketRegistry is the single place where old V2 suite scenarios are assigned to buckets.
// When adding a new scenario, add it here and keep bucket runtimes balanced. Best way to do it is by
// executing the tests in CI once and asking an AI to check run details, with execution time and to
// rebalance the buckets so that they are balanced.
var suiteBucketRegistry = []suiteBucketDefinition{
	{
		Bucket: SuiteBucketA,
		Scenarios: []SuiteScenario{
			SuiteScenarioProofOfReserve,
			SuiteScenarioHTTPTriggerAction,
			SuiteScenarioDONTime,
			SuiteScenarioConsensus,
		},
	},
	{
		Bucket: SuiteBucketB,
		Scenarios: []SuiteScenario{
			SuiteScenarioVaultDON,
		},
	},
	{
		Bucket: SuiteBucketC,
		Scenarios: []SuiteScenario{
			SuiteScenarioCronBeholder,
			SuiteScenarioHTTPActionCRUD,
		},
	},
}

func ScenariosForSuiteBucket(bucket SuiteBucket) ([]SuiteScenario, error) {
	for _, bucketDef := range suiteBucketRegistry {
		if bucketDef.Bucket != bucket {
			continue
		}

		scenarios := make([]SuiteScenario, len(bucketDef.Scenarios))
		copy(scenarios, bucketDef.Scenarios)
		return scenarios, nil
	}

	return nil, fmt.Errorf("unknown V2 suite bucket %q", bucket)
}

func ValidateSuiteBucketRegistry() error {
	assignedScenarios := make(map[SuiteScenario]SuiteBucket, SuiteScenarioLen)

	for _, bucketDef := range suiteBucketRegistry {
		for _, scenario := range bucketDef.Scenarios {
			if scenario < 0 || scenario >= SuiteScenarioLen {
				return fmt.Errorf("invalid scenario %d in bucket %q", scenario, bucketDef.Bucket)
			}

			if existingBucket, ok := assignedScenarios[scenario]; ok {
				return fmt.Errorf("scenario %q assigned to multiple buckets: %q and %q", scenario.String(), existingBucket, bucketDef.Bucket)
			}

			assignedScenarios[scenario] = bucketDef.Bucket
		}
	}

	for scenario := SuiteScenario(0); scenario < SuiteScenarioLen; scenario++ {
		if _, ok := assignedScenarios[scenario]; ok {
			continue
		}

		return fmt.Errorf("scenario %q is not assigned to any V2 suite bucket", scenario.String())
	}

	return nil
}
