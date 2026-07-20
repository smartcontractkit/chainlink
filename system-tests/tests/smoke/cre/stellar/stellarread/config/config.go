package config

import "fmt"

// TestCase selects which Stellar read capability action the workflow performs.
type TestCase int

const (
	TestCaseLatestLedger TestCase = iota
	TestCaseReadContract
)

func (tc TestCase) String() string {
	switch tc {
	case TestCaseLatestLedger:
		return "latest_ledger"
	case TestCaseReadContract:
		return "read_contract"
	default:
		return fmt.Sprintf("unknown TestCase: %d", int(tc))
	}
}

// Config for the Stellar read workflow, selected by TestCase.
type Config struct {
	ChainSelector     uint64 `yaml:"chainSelector"`
	WorkflowName      string `yaml:"workflowName"`
	MinLedgerSequence uint64 `yaml:"minLedgerSequence"`

	// CronSchedule overrides the workflow trigger schedule (6-field cron). Empty keeps the safe default "*/30 * * * * *".
	CronSchedule string `yaml:"cronSchedule"`

	ReadKind      TestCase `yaml:"readKind"`
	SourceAccount string   `yaml:"sourceAccount"`

	// Cases is the batch of ReadContract invocations run and asserted in a single trigger
	Cases []ReadContractAsserts `yaml:"cases"`
}

// ReadContractAsserts is one ReadContract invocation asserted within a batch run.
type ReadContractAsserts struct {
	Name                      string   `yaml:"name"`
	ContractID                string   `yaml:"contractID"`
	Function                  string   `yaml:"function"`
	ArgMode                   string   `yaml:"argMode"`
	ArgBool                   bool     `yaml:"argBool"`
	ArgU32A                   uint32   `yaml:"argU32A"`
	ArgU32B                   uint32   `yaml:"argU32B"`
	ArgString                 string   `yaml:"argString"`
	ArgBytesHex               string   `yaml:"argBytesHex"`
	ArgI64                    int64    `yaml:"argI64"`
	ArgSymbol                 string   `yaml:"argSymbol"`
	ArgVecU32                 []uint32 `yaml:"argVecU32"`
	ArgReceiverContractIDHex  string   `yaml:"argReceiverContractIDHex"`  // 32-byte hex -> ScVal Address(ContractId)
	ArgWorkflowExecutionIDHex string   `yaml:"argWorkflowExecutionIDHex"` // 32-byte hex -> ScVal Bytes / BytesN<32>
	ArgReportIDHex            string   `yaml:"argReportIDHex"`            // 2-byte hex  -> ScVal Bytes / BytesN<2>
	ExpectedResult            string   `yaml:"expectedResult"`            // host-precomputed base64 XDR ScVal; "" = shape-only assert
	ExpectError               bool     `yaml:"expectError"`
}
