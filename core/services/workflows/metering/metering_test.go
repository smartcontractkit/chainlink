package metering

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder/beholdertest"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/metrics"
	billing "github.com/smartcontractkit/chainlink-protos/billing/go"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"
	eventspb "github.com/smartcontractkit/chainlink-protos/workflows/go/events"

	"github.com/smartcontractkit/chainlink/v2/core/platform"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/events"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
)

const (
	testAccountID           = "accountId"
	testWorkflowID          = "workflowId"
	testWorkflowExecutionID = "workflowExecutionId"
	dummyRegistryAddress    = "0x123"
	dummyChainSelector      = "16015286601757825753" // Ethereum Sepolia chain selector
	workflowV1              = "v1"
	workflowV2              = "v2"
)

var (
	defaultLabels = map[string]string{
		platform.KeyWorkflowOwner:       "accountId",
		platform.KeyWorkflowID:          "workflowId",
		platform.KeyWorkflowName:        "workflowName",
		platform.KeyWorkflowVersion:     workflowV2,
		platform.KeyWorkflowExecutionID: "workflowExecutionId",
		platform.KeyDonID:               "42",
		platform.KeyDonF:                "1",
		platform.KeyDonN:                "3",
		platform.KeyP2PID:               "peerId",
		platform.KeyTriggerID:           "triggerId",
	}
	testUnitA      = billing.ResourceType_name[int32(billing.ResourceType_RESOURCE_TYPE_COMPUTE)]
	testUnitB      = billing.ResourceType_name[int32(billing.ResourceType_RESOURCE_TYPE_UNSPECIFIED)]
	testUnitC      = billing.ResourceType_name[int32(billing.ResourceType_RESOURCE_TYPE_NETWORK)]
	testUnitGas    = "GAS.5009297550715157269" // ETH mainnet
	validConfig, _ = values.NewMap(map[string]any{
		RatiosKey: map[string]string{
			testUnitA: "0.4",
			testUnitB: "0.6",
		},
	})
	two = decimal.NewFromInt(2)
)

func Test_Report(t *testing.T) {
	t.Parallel()

	t.Run("error if incorrect labels", func(t *testing.T) {
		t.Parallel()

		_, err := NewReport(t.Context(), map[string]string{}, logger.Nop(), defaultMetrics(t), dummyRegistryAddress, dummyChainSelector, workflowV1)
		require.ErrorIs(t, err, ErrMissingLabels)
	})
}

func Test_Report_MeteringMode(t *testing.T) {
	t.Parallel()

	t.Run("Deduct", func(t *testing.T) {
		t.Parallel()

		t.Run("returns empty limits and no error in metering mode", func(t *testing.T) {
			t.Parallel()

			emptyUserSpendLimit := decimal.NewNullDecimal(decimal.Zero)
			report := newTestReport(t, logger.Nop())

			require.NoError(t, report.Reserve())

			limits, err := report.Deduct("ref1", ByDerivedAvailability(emptyUserSpendLimit, 1, capabilities.CapabilityInfo{}, nil))

			require.NoError(t, err)
			assert.Empty(t, limits)
		})

		t.Run("does not modify local balance in metering mode", func(t *testing.T) {
			t.Parallel()

			report := newTestReport(t, logger.Nop())

			require.NoError(t, report.Reserve())

			balanceBefore := report.balance.balance
			_, err := report.Deduct("ref1", ByResource(testUnitA, "", two))

			require.NoError(t, err)
			_, err = report.Deduct("ref2", ByDerivedAvailability(decimal.NewNullDecimal(decimal.Zero), 1, capabilities.CapabilityInfo{}, nil))
			require.NoError(t, err)

			balanceAfter := report.balance.balance
			assert.Equal(t, balanceBefore, balanceAfter)
		})

		t.Run("switches to metering mode", func(t *testing.T) {
			t.Parallel()

			t.Run("if only one spend type and rate does not exist", func(t *testing.T) {
				t.Parallel()

				lggr, logs := logger.TestObserved(t, zapcore.DebugLevel)
				report := newTestReport(t, lggr)

				require.NoError(t, report.Reserve())

				info := capabilities.CapabilityInfo{
					SpendTypes: []capabilities.CapabilitySpendType{capabilities.CapabilitySpendType(testUnitB)},
				}

				// ratios and spend types should match
				config, _ := values.NewMap(map[string]any{
					RatiosKey: map[string]string{
						testUnitB: "1",
					},
				})

				limits, err := report.Deduct("ref1", ByDerivedAvailability(decimal.NewNullDecimal(decimal.Zero), 1, info, config))

				require.NoError(t, err)
				assert.Empty(t, limits)
				assert.True(t, report.meteringMode)
				assert.Len(t, logs.FilterMessage("switching to metering mode").All(), 1)
			})

			t.Run("if ratio and spend type lengths do not match", func(t *testing.T) {
				t.Parallel()

				lggr, logs := logger.TestObserved(t, zapcore.DebugLevel)
				report := newTestReport(t, lggr)

				require.NoError(t, report.Reserve())

				info := capabilities.CapabilityInfo{
					SpendTypes: []capabilities.CapabilitySpendType{capabilities.CapabilitySpendType(testUnitA), capabilities.CapabilitySpendType(testUnitB), capabilities.CapabilitySpendType(testUnitC)},
				}

				// 3 spend types and 2 ratios creates the mismatch
				limits, err := report.Deduct("ref1", ByDerivedAvailability(decimal.NewNullDecimal(decimal.Zero), 1, info, validConfig))

				require.NoError(t, err)
				assert.Empty(t, limits)
				assert.True(t, report.meteringMode)
				assert.Len(t, logs.FilterMessage("switching to metering mode").All(), 1)
			})

			t.Run("if multiple spend types and ratio does not exist", func(t *testing.T) {
				t.Parallel()

				lggr, logs := logger.TestObserved(t, zapcore.DebugLevel)
				report := newTestReport(t, lggr)

				require.NoError(t, report.Reserve())

				info := capabilities.CapabilityInfo{
					SpendTypes: []capabilities.CapabilitySpendType{capabilities.CapabilitySpendType(testUnitA), capabilities.CapabilitySpendType(testUnitC)},
				}

				// spend types and ratios should not match and return an error
				limits, err := report.Deduct("ref1", ByDerivedAvailability(decimal.NewNullDecimal(decimal.Zero), 1, info, validConfig))

				require.NoError(t, err)
				assert.Empty(t, limits)
				assert.True(t, report.meteringMode)
				assert.Len(t, logs.FilterMessage("switching to metering mode").All(), 1)
			})

			t.Run("if multiple spend types and rate does not exist", func(t *testing.T) {
				t.Parallel()

				lggr, logs := logger.TestObserved(t, zapcore.DebugLevel)
				report := newTestReport(t, lggr)

				require.NoError(t, report.Reserve())

				info := capabilities.CapabilityInfo{
					SpendTypes: []capabilities.CapabilitySpendType{capabilities.CapabilitySpendType(testUnitA), capabilities.CapabilitySpendType(testUnitB)},
				}

				limits, err := report.Deduct("ref1", ByDerivedAvailability(decimal.NewNullDecimal(decimal.Zero), 1, info, validConfig))

				require.NoError(t, err)
				assert.Empty(t, limits)
				assert.True(t, report.meteringMode)
				assert.Len(t, logs.FilterMessage("switching to metering mode").All(), 1)
			})
		})
	})

	t.Run("Settle does not modify local balance in metering mode", func(t *testing.T) {
		t.Parallel()

		report := newTestReport(t, logger.Nop())

		require.NoError(t, report.Reserve())

		balanceBefore := report.balance.balance

		_, err := report.Deduct("ref1", ByResource(testUnitA, "", two))
		require.NoError(t, err)

		steps := capabilities.ResponseMetadata{Metering: []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "xyz", SpendUnit: testUnitA, SpendValue: "2"},
		}}
		require.NoError(t, report.Settle("ref1", steps))

		balanceAfter := report.balance.balance
		require.Equal(t, balanceBefore, balanceAfter)
	})

	t.Run("single execution switches to metering mode for multiple executions", func(t *testing.T) {
		t.Parallel()

		report := newTestReport(t, logger.Nop())

		require.NoError(t, report.Reserve())

		stepRef := "ref1"

		_, err := report.Deduct(stepRef, ByResource(testUnitA, "", decimal.NewFromInt(1)))
		require.NoError(t, err)

		balanceBefore := report.balance.balance
		require.NoError(t, report.Settle(stepRef, capabilities.ResponseMetadata{Metering: []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "xyz", SpendUnit: billing.ResourceType_RESOURCE_TYPE_COMPUTE.String(), SpendValue: "42"},
			{Peer2PeerID: "abc", SpendUnit: billing.ResourceType_RESOURCE_TYPE_COMPUTE.String(), SpendValue: "44"},
			{Peer2PeerID: "abc", SpendUnit: testUnitGas, SpendValue: "0.000001"}, // 1000 gwei as a decimal
			{Peer2PeerID: "lmno", SpendUnit: billing.ResourceType_RESOURCE_TYPE_COMPUTE.String(), SpendValue: "12"},
			{Peer2PeerID: "lmno", SpendUnit: testUnitGas, SpendValue: "0.000001"}, // 1000 gwei as a decimal
		}}))

		balanceAfter := report.balance.balance

		assert.Equal(t, balanceBefore.String(), balanceAfter.String())
		assert.True(t, report.meteringMode)
	})
}

func Test_medianSpend(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		input    []decimal.Decimal
		expected string
	}{
		{
			name: "MedianSpend returns median for a list of int spend values",
			input: []decimal.Decimal{
				decimal.NewFromInt(1),
				decimal.NewFromInt(2),
				decimal.NewFromInt(3),
			},
			expected: "2",
		},
		{
			name: "MedianSpend returns median for a list of float spend values",
			input: []decimal.Decimal{
				decimal.NewFromFloat(0.1),
				decimal.NewFromFloat(0.2),
				decimal.NewFromFloat(0.3),
			},
			expected: "0.2",
		},
		{
			name: "MedianSpend returns median single spend value",
			input: []decimal.Decimal{
				decimal.NewFromInt(1),
			},
			expected: "1",
		},
		{
			name: "MedianSpend returns median even number of spend values",
			input: []decimal.Decimal{
				decimal.NewFromInt(2),
				decimal.NewFromInt(2),
				decimal.NewFromInt(4),
				decimal.NewFromInt(4),
			},
			expected: "3",
		},
		{
			name: "MedianSpend returns median odd number of spend values",
			input: []decimal.Decimal{
				decimal.NewFromInt(1),
				decimal.NewFromInt(13),
				decimal.NewFromInt(50),
				decimal.NewFromInt(51),
				decimal.NewFromInt(100),
			},
			expected: "50",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, medianSpend(tc.input).String())
		})
	}
}

func Test_Report_Reserve(t *testing.T) {
	t.Parallel()

	t.Run("marks the execution as reserved", func(t *testing.T) {
		t.Parallel()

		report := newTestReport(t, logger.Nop())

		require.NoError(t, report.Reserve())

		_, err := report.Deduct("ref1", ByResource(testUnitA, "", two))
		require.NoError(t, err)
	})
}

func Test_Report_Deduct(t *testing.T) {
	t.Parallel()

	one := decimal.NewFromInt(1)
	two := decimal.NewFromInt(2)

	t.Run("returns an error if not initialized", func(t *testing.T) {
		t.Parallel()

		report := newTestReport(t, logger.Nop())
		_, err := report.Deduct("ref1", ByResource(testUnitA, "", one))

		require.ErrorIs(t, err, ErrNoReserve)
	})

	t.Run("returns an error if step already exists", func(t *testing.T) {
		t.Parallel()

		report := newTestReport(t, logger.Nop())

		require.NoError(t, report.Reserve())

		_, err := report.Deduct("ref1", ByResource(testUnitA, "", two))
		require.NoError(t, err)

		_, err = report.Deduct("ref1", ByResource(testUnitA, "", one))
		require.ErrorIs(t, err, ErrStepDeductExists)
	})

	t.Run("deduct runs in metering mode without a rate card", func(t *testing.T) {
		t.Parallel()

		report := newTestReport(t, logger.Nop())

		require.NoError(t, report.Reserve())

		config, _ := values.NewMap(map[string]any{
			RatiosKey: map[string]string{
				testUnitA:   "0.4",
				testUnitGas: "0.6",
			},
		})

		info := capabilities.CapabilityInfo{
			SpendTypes: []capabilities.CapabilitySpendType{
				capabilities.CapabilitySpendType(testUnitA),
				capabilities.CapabilitySpendType(testUnitGas),
			},
		}

		emptyLimit := decimal.NewNullDecimal(decimal.Zero)
		emptyLimit.Valid = false
		limits, err := report.Deduct("ref1", ByDerivedAvailability(emptyLimit, 1, info, config))

		require.NoError(t, err)
		require.NotNil(t, limits)
		require.Empty(t, limits, "no rate card: no spend limits are derived")
		assert.True(t, report.meteringMode)

	})

	t.Run("empty limits for empty spend types", func(t *testing.T) {
		t.Parallel()

		report := newTestReport(t, logger.Nop())

		require.NoError(t, report.Reserve())

		config, _ := values.NewMap(map[string]any{
			RatiosKey: map[string]string{
				testUnitA: "0.4",
				testUnitC: "0.6",
			},
		})

		limits, err := report.Deduct("ref1", ByDerivedAvailability(decimal.NewNullDecimal(decimal.Zero), 1, capabilities.CapabilityInfo{}, config))

		require.NoError(t, err)
		assert.Empty(t, limits)

	})

	emptyUserSpendLimit := decimal.NewNullDecimal(decimal.Zero)
	emptyUserSpendLimit.Valid = false

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		t.Run("if open slots is 0", func(t *testing.T) {
			t.Parallel()

			report := newTestReport(t, logger.Nop())

			require.NoError(t, report.Reserve())

			_, err := report.Deduct("ref1", ByDerivedAvailability(emptyUserSpendLimit, 0, capabilities.CapabilityInfo{}, nil))
			require.ErrorIs(t, ErrNoOpenCalls, err)
		})

		t.Run("if reserve is not called first", func(t *testing.T) {
			t.Parallel()

			report := newTestReport(t, logger.Nop())
			_, err := report.Deduct("ref1", ByDerivedAvailability(emptyUserSpendLimit, 1, capabilities.CapabilityInfo{}, nil))

			require.ErrorIs(t, ErrNoReserve, err)
		})
	})
}

func Test_Report_Settle(t *testing.T) {
	t.Parallel()

	t.Run("returns an error if not initialized", func(t *testing.T) {
		t.Parallel()

		report := newTestReport(t, logger.Nop())

		require.ErrorIs(t, report.Settle("ref1", capabilities.ResponseMetadata{}), ErrNoReserve)
	})

	t.Run("returns an error if Deduct is not called first", func(t *testing.T) {
		t.Parallel()

		report := newTestReport(t, logger.Nop())

		require.NoError(t, report.Reserve())
		require.ErrorIs(t, report.Settle("ref1", capabilities.ResponseMetadata{}), ErrNoDeduct)
	})

	t.Run("returns an error if step already exists", func(t *testing.T) {
		t.Parallel()

		report := newTestReport(t, logger.Nop())

		require.NoError(t, report.Reserve())

		steps := capabilities.ResponseMetadata{Metering: []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "abc", SpendUnit: testUnitA, SpendValue: "1"},
		}}

		_, err := report.Deduct("ref1", ByResource(testUnitA, "", decimal.NewFromInt(2)))
		require.NoError(t, err)
		require.NoError(t, report.Settle("ref1", steps))
		require.ErrorIs(t, report.Settle("ref1", steps), ErrStepSpendExists)
	})

	t.Run("ignores invalid spend values", func(t *testing.T) {
		t.Parallel()

		lggr, logs := logger.TestObserved(t, zapcore.InfoLevel)
		report := newTestReport(t, lggr)

		require.NoError(t, report.Reserve())

		steps := capabilities.ResponseMetadata{Metering: []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "xyz", SpendUnit: testUnitA, SpendValue: "????"},
			{Peer2PeerID: "abc", SpendUnit: testUnitA, SpendValue: "1"},
		}}

		_, err := report.Deduct("ref1", ByResource(testUnitA, "", decimal.NewFromInt(2)))
		require.NoError(t, err)

		require.NoError(t, report.Settle("ref1", steps))
		assert.Len(t, logs.All(), 1)
	})

	t.Run("does not error when spend exceeds reservation", func(t *testing.T) {
		t.Parallel()

		lggr, logs := logger.TestObserved(t, zapcore.InfoLevel)
		report := newTestReport(t, lggr)

		require.NoError(t, report.Reserve())

		steps := capabilities.ResponseMetadata{Metering: []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "xyz", SpendUnit: testUnitA, SpendValue: "2"},
		}}

		_, err := report.Deduct("ref1", ByResource(testUnitA, "", decimal.NewFromInt(1)))
		require.NoError(t, err)

		require.NoError(t, report.Settle("ref1", steps))
		assert.True(t, report.isMeteringMode(), "no rate card: settles in metering mode")
		assert.Empty(t, logs.All())
	})

	t.Run("successfully settles gas token usage as single execution", func(t *testing.T) {
		t.Parallel()

		lggr, logs := logger.TestObserved(t, zapcore.InfoLevel)
		report := newTestReport(t, lggr)

		require.NoError(t, report.Reserve())

		config, _ := values.NewMap(map[string]any{
			RatiosKey: map[string]string{
				testUnitGas: "1.0",
			},
		})

		info := capabilities.CapabilityInfo{
			SpendTypes: []capabilities.CapabilitySpendType{
				capabilities.CapabilitySpendType(testUnitGas),
			},
		}

		value := decimal.NewNullDecimal(decimal.Zero)
		value.Valid = false

		_, err := report.Deduct("ref1", ByDerivedAvailability(value, 1, info, config))
		require.NoError(t, err)

		steps := capabilities.ResponseMetadata{Metering: []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "xyz", SpendUnit: testUnitGas, SpendValue: "0.000700000000000000"}, // should convert to 3.5 credits
		}, CapDON_N: 42}

		require.NoError(t, report.Settle("ref1", steps))

		assert.True(t, report.isMeteringMode(), "no rate card: settles in metering mode")

		assert.Empty(t, logs.All())
	})

	t.Run("successfully settles zero value rate", func(t *testing.T) {
		t.Parallel()

		lggr, logs := logger.TestObserved(t, zapcore.InfoLevel)
		report := newTestReport(t, lggr)

		require.NoError(t, report.Reserve())

		config, _ := values.NewMap(map[string]any{
			RatiosKey: map[string]string{
				testUnitA:   "0.5",
				testUnitGas: "0.5",
			},
		})

		info := capabilities.CapabilityInfo{
			SpendTypes: []capabilities.CapabilitySpendType{
				capabilities.CapabilitySpendType(testUnitA),
				capabilities.CapabilitySpendType(testUnitGas),
			},
		}

		value := decimal.NewNullDecimal(decimal.Zero)
		value.Valid = false

		_, err := report.Deduct("ref1", ByDerivedAvailability(value, 1, info, config))
		require.NoError(t, err)

		steps := capabilities.ResponseMetadata{Metering: []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "xyz", SpendUnit: testUnitA, SpendValue: "0.000007"},
			{Peer2PeerID: "xyz", SpendUnit: testUnitGas, SpendValue: "0.000700000000000000"}, // should convert to 0 credits
		}, CapDON_N: 42}

		require.NoError(t, report.Settle("ref1", steps))

		assert.Equal(t, "0", report.balance.GetSpent().String())

		assert.Empty(t, logs.All())
	})
}

func Test_Report_FormatReport(t *testing.T) {
	t.Parallel()

	t.Run("does not contain metadata", func(t *testing.T) {
		t.Parallel()

		report := newTestReport(t, logger.Nop())

		require.NoError(t, report.Reserve())

		meteringReport := report.FormatReport()
		require.Equal(t, &eventspb.WorkflowMetadata{
			WorkflowOwner:           "accountId",
			WorkflowName:            "workflowName",
			Version:                 workflowV2,
			WorkflowID:              "workflowId",
			WorkflowExecutionID:     "workflowExecutionId",
			DonID:                   42,
			DonF:                    1,
			DonN:                    3,
			P2PID:                   "peerId",
			WorkflowRegistryAddress: "0x123",
			WorkflowRegistryChain:   "16015286601757825753",
			EngineVersion:           workflowV2,
			Trigger: &eventspb.TriggerDetail{
				TriggerID: "triggerId",
			},
		}, meteringReport.Metadata)
	})

	t.Run("includes orgID in metadata", func(t *testing.T) {
		t.Parallel()

		testOrgID := "org-123"
		labels := map[string]string{
			platform.KeyWorkflowOwner:       "accountId",
			platform.KeyWorkflowID:          "workflowId",
			platform.KeyWorkflowVersion:     workflowV2,
			platform.KeyWorkflowExecutionID: "workflowExecutionId",
			platform.KeyDonID:               "42",
			platform.KeyDonF:                "1",
			platform.KeyDonN:                "3",
			platform.KeyP2PID:               "peerId",
			platform.KeyTriggerID:           "triggerId",
			platform.KeyOrganizationID:      testOrgID,
		}

		report, err := NewReport(t.Context(), labels, logger.Nop(), defaultMetrics(t), dummyRegistryAddress, dummyChainSelector, workflowV2)
		require.NoError(t, err)

		require.NoError(t, report.Reserve())

		meteringReport := report.FormatReport()
		require.Equal(t, testOrgID, meteringReport.Metadata.OrgID)
	})

	t.Run("contains all step data", func(t *testing.T) {
		t.Parallel()

		numSteps := 100
		report := newTestReport(t, logger.Nop())

		require.NoError(t, report.Reserve())

		expected := map[string]*eventspb.MeteringReportStep{}

		for i := range numSteps {
			stepRef := strconv.Itoa(i)

			_, err := report.Deduct(stepRef, ByResource(testUnitA, "", decimal.NewFromInt(1)))
			require.NoError(t, err)

			require.NoError(t, report.Settle(stepRef, capabilities.ResponseMetadata{Metering: []capabilities.MeteringNodeDetail{
				{Peer2PeerID: "xyz", SpendUnit: billing.ResourceType_RESOURCE_TYPE_COMPUTE.String(), SpendValue: "42"},
			}, CapDON_N: 10}))

			expected[stepRef] = &eventspb.MeteringReportStep{
				Nodes: []*eventspb.MeteringReportNodeDetail{
					{
						Peer_2PeerId:  "xyz",
						SpendUnit:     billing.ResourceType_RESOURCE_TYPE_COMPUTE.String(),
						SpendValue:    "42",
						SpendValueCre: "0.0000000000",
					},
				},
				AggSpend: []*eventspb.AggregatedSpendDetail{
					{
						SpendValue:    "42.0000000000",
						SpendUnit:     billing.ResourceType_RESOURCE_TYPE_COMPUTE.String(),
						SpendValueCre: "0.0000000000",
					},
				},
				AggSpendValue:    "42.0000000000",
				AggSpendUnit:     "RESOURCE_TYPE_COMPUTE",
				AggSpendValueCre: "0.0000000000",
				CapdonN:          10,
			}
		}

		assert.Equal(t, expected, report.FormatReport().Steps)
	})

	t.Run("aggregated value shows in the report", func(t *testing.T) {
		t.Parallel()

		numSteps := 100
		report := newTestReport(t, logger.Nop())

		require.NoError(t, report.Reserve())

		expected := map[string]*eventspb.MeteringReportStep{}

		for i := range numSteps {
			stepRef := strconv.Itoa(i)

			_, err := report.Deduct(stepRef, ByResource(testUnitA, "", decimal.NewFromInt(1)))
			require.NoError(t, err)

			require.NoError(t, report.Settle(stepRef, capabilities.ResponseMetadata{Metering: []capabilities.MeteringNodeDetail{
				{Peer2PeerID: "xyz", SpendUnit: billing.ResourceType_RESOURCE_TYPE_COMPUTE.String(), SpendValue: "42"},
				{Peer2PeerID: "abc", SpendUnit: billing.ResourceType_RESOURCE_TYPE_COMPUTE.String(), SpendValue: "44"},
				{Peer2PeerID: "lmno", SpendUnit: billing.ResourceType_RESOURCE_TYPE_COMPUTE.String(), SpendValue: "12"},
			}}))

			expected[stepRef] = &eventspb.MeteringReportStep{
				Nodes: []*eventspb.MeteringReportNodeDetail{
					{
						Peer_2PeerId:  "xyz",
						SpendUnit:     billing.ResourceType_RESOURCE_TYPE_COMPUTE.String(),
						SpendValue:    "42",
						SpendValueCre: "0.0000000000",
					},
					{
						Peer_2PeerId:  "abc",
						SpendUnit:     billing.ResourceType_RESOURCE_TYPE_COMPUTE.String(),
						SpendValue:    "44",
						SpendValueCre: "0.0000000000",
					},
					{
						Peer_2PeerId:  "lmno",
						SpendUnit:     billing.ResourceType_RESOURCE_TYPE_COMPUTE.String(),
						SpendValue:    "12",
						SpendValueCre: "0.0000000000",
					},
				},
				AggSpendValue:    "42.0000000000",
				AggSpendUnit:     "RESOURCE_TYPE_COMPUTE",
				AggSpendValueCre: "0.0000000000",
				CapdonN:          1,
				AggSpend: []*eventspb.AggregatedSpendDetail{
					{
						SpendValue:    "42.0000000000", // median of 42, 44, 12
						SpendUnit:     billing.ResourceType_RESOURCE_TYPE_COMPUTE.String(),
						SpendValueCre: "0.0000000000",
					},
				},
			}
		}

		assert.Equal(t, expected, report.FormatReport().Steps)
	})

	t.Run("single execution does not aggregate", func(t *testing.T) {
		t.Parallel()

		numSteps := 1
		report := newTestReport(t, logger.Nop())

		require.NoError(t, report.Reserve())

		expected := map[string]*eventspb.MeteringReportStep{}

		for i := range numSteps {
			stepRef := strconv.Itoa(i)

			_, err := report.Deduct(stepRef, ByResource(testUnitA, "", decimal.NewFromInt(1)))
			require.NoError(t, err)

			require.NoError(t, report.Settle(stepRef, capabilities.ResponseMetadata{Metering: []capabilities.MeteringNodeDetail{
				{Peer2PeerID: "xyz", SpendUnit: billing.ResourceType_RESOURCE_TYPE_COMPUTE.String(), SpendValue: "42"},
				{Peer2PeerID: "abc", SpendUnit: billing.ResourceType_RESOURCE_TYPE_COMPUTE.String(), SpendValue: "44"},
				{Peer2PeerID: "abc", SpendUnit: testUnitGas, SpendValue: "0.000001"}, // 1000 gwei as a decimal
				{Peer2PeerID: "lmno", SpendUnit: billing.ResourceType_RESOURCE_TYPE_COMPUTE.String(), SpendValue: "12"},
				{Peer2PeerID: "lmno", SpendUnit: testUnitGas, SpendValue: "0.000001"}, // 1000 gwei as a decimal
			}}))

			expected[stepRef] = &eventspb.MeteringReportStep{
				Nodes: []*eventspb.MeteringReportNodeDetail{
					{
						Peer_2PeerId:  "xyz",
						SpendUnit:     billing.ResourceType_RESOURCE_TYPE_COMPUTE.String(),
						SpendValue:    "42",
						SpendValueCre: "0.0000000000",
					},
					{
						Peer_2PeerId:  "abc",
						SpendUnit:     billing.ResourceType_RESOURCE_TYPE_COMPUTE.String(),
						SpendValue:    "44",
						SpendValueCre: "0.0000000000",
					},
					{
						Peer_2PeerId:  "lmno",
						SpendUnit:     billing.ResourceType_RESOURCE_TYPE_COMPUTE.String(),
						SpendValue:    "12",
						SpendValueCre: "0.0000000000",
					},
					{
						Peer_2PeerId:  "abc",
						SpendUnit:     testUnitGas,
						SpendValue:    "0.000001",
						SpendValueCre: "0.0000000000",
					},
					{
						Peer_2PeerId:  "lmno",
						SpendUnit:     testUnitGas,
						SpendValue:    "0.000001",
						SpendValueCre: "0.0000000000",
					},
				},
				AggSpendValue:    "1000000000000.0000000000",
				AggSpendUnit:     "GAS.5009297550715157269",
				AggSpendValueCre: "0.0000000000",
				CapdonN:          1,
				AggSpend: []*eventspb.AggregatedSpendDetail{
					{
						SpendValue:    "42.0000000000", // median of 42, 44, 12
						SpendUnit:     billing.ResourceType_RESOURCE_TYPE_COMPUTE.String(),
						SpendValueCre: "0.0000000000",
					},
					{
						SpendValue:    "1000000000000.0000000000", // converted to wei before median is taken
						SpendUnit:     testUnitGas,
						SpendValueCre: "0.0000000000",
					},
				},
			}
		}

		assert.Equal(t, expected, report.FormatReport().Steps)
	})

	t.Run("gas spend with spend_value_in_gas_units skips shift", func(t *testing.T) {
		t.Parallel()

		numSteps := 1
		report := newTestReport(t, logger.Nop())

		require.NoError(t, report.Reserve())

		expected := map[string]*eventspb.MeteringReportStep{}

		for i := range numSteps {
			stepRef := strconv.Itoa(i)

			_, err := report.Deduct(stepRef, ByResource(testUnitA, "", decimal.NewFromInt(1)))
			require.NoError(t, err)

			require.NoError(t, report.Settle(stepRef, capabilities.ResponseMetadata{Metering: []capabilities.MeteringNodeDetail{
				{Peer2PeerID: "xyz", SpendUnit: billing.ResourceType_RESOURCE_TYPE_COMPUTE.String(), SpendValue: "42"},
				{Peer2PeerID: "abc", SpendUnit: testUnitGas, SpendValue: "0.000001", SpendValueInGasUnits: "1000000000000"}, // 1000 gwei in wei, no shift needed
			}}))

			expected[stepRef] = &eventspb.MeteringReportStep{
				Nodes: []*eventspb.MeteringReportNodeDetail{
					{
						Peer_2PeerId:  "xyz",
						SpendUnit:     billing.ResourceType_RESOURCE_TYPE_COMPUTE.String(),
						SpendValue:    "42",
						SpendValueCre: "0.0000000000",
					},
					{
						Peer_2PeerId: "abc",
						SpendUnit:    testUnitGas,
						// native fixed-point value is surfaced in the legacy field for
						// consumers that only read SpendValue (e.g. billing)
						SpendValue:           "1000000000000",
						SpendValueCre:        "0.0000000000",
						SpendValueInGasUnits: "1000000000000",
					},
				},
				AggSpendValue:    "1000000000000.0000000000",
				AggSpendUnit:     "GAS.5009297550715157269",
				AggSpendValueCre: "0.0000000000",
				CapdonN:          1,
				AggSpend: []*eventspb.AggregatedSpendDetail{
					{
						SpendValue:    "42.0000000000",
						SpendUnit:     billing.ResourceType_RESOURCE_TYPE_COMPUTE.String(),
						SpendValueCre: "0.0000000000",
					},
					{
						SpendValue:    "1000000000000.0000000000",
						SpendUnit:     testUnitGas,
						SpendValueCre: "0.0000000000",
					},
				},
			}
		}

		assert.Equal(t, expected, report.FormatReport().Steps)
	})

	t.Run("mixed legacy and new gas nodes switches to metering mode", func(t *testing.T) {
		t.Parallel()

		report := newTestReport(t, logger.Nop())

		require.NoError(t, report.Reserve())

		_, err := report.Deduct("ref1", ByResource(testUnitA, "", decimal.NewFromInt(1)))
		require.NoError(t, err)

		// One node sends the new SpendValueInGasUnits, another sends only legacy SpendValue
		require.NoError(t, report.Settle("ref1", capabilities.ResponseMetadata{Metering: []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "abc", SpendUnit: testUnitGas, SpendValue: "0.000001", SpendValueInGasUnits: "1000000000000"},
			{Peer2PeerID: "lmno", SpendUnit: testUnitGas, SpendValue: "0.000001"},
		}}))

		assert.True(t, report.meteringMode)
	})

	t.Run("ignores invalid spend_value_in_gas_units", func(t *testing.T) {
		t.Parallel()

		lggr, logs := logger.TestObserved(t, zapcore.InfoLevel)
		report := newTestReport(t, lggr)

		require.NoError(t, report.Reserve())

		_, err := report.Deduct("ref1", ByResource(testUnitA, "", decimal.NewFromInt(1)))
		require.NoError(t, err)

		// One node sends invalid SpendValueInGasUnits (skipped), another sends valid value
		require.NoError(t, report.Settle("ref1", capabilities.ResponseMetadata{Metering: []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "abc", SpendUnit: testUnitGas, SpendValue: "0.000000000000000001", SpendValueInGasUnits: "????"},
			{Peer2PeerID: "xyz", SpendUnit: testUnitGas, SpendValue: "0.000000000000000001", SpendValueInGasUnits: "1"},
		}}))

		assert.Len(t, logs.All(), 1)
	})

	// Aptos specific regression test
	t.Run("non-18-decimal gas token uses gas units, not legacy shift", func(t *testing.T) {
		t.Parallel()

		const aptosUnit = "GAS.4741433654826277614"

		report := newTestReport(t, logger.Nop())

		require.NoError(t, report.Reserve())

		_, err := report.Deduct("ref1", ByResource(testUnitA, "", decimal.NewFromInt(1)))
		require.NoError(t, err)

		// 0.001158 APT = 115800 octas
		require.NoError(t, report.Settle("ref1", capabilities.ResponseMetadata{Metering: []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "abc", SpendUnit: aptosUnit, SpendValue: "0.001158", SpendValueInGasUnits: "115800"},
		}}))

		expected := map[string]*eventspb.MeteringReportStep{
			"ref1": {
				Nodes: []*eventspb.MeteringReportNodeDetail{
					{
						Peer_2PeerId: "abc",
						SpendUnit:    aptosUnit,
						// native fixed-point value is surfaced in the legacy field
						SpendValue:           "115800",
						SpendValueCre:        "0.0000000000",
						SpendValueInGasUnits: "115800",
					},
				},
				AggSpendValue:    "115800.0000000000",
				AggSpendUnit:     aptosUnit,
				AggSpendValueCre: "0.0000000000",
				CapdonN:          1,
				AggSpend: []*eventspb.AggregatedSpendDetail{
					{
						SpendValue:    "115800.0000000000",
						SpendUnit:     aptosUnit,
						SpendValueCre: "0.0000000000",
					},
				},
			},
		}

		assert.Equal(t, expected, report.FormatReport().Steps)
		assert.True(t, report.meteringMode)
	})

	t.Run("invalid native gas units are dropped and not promoted to spend value", func(t *testing.T) {
		t.Parallel()

		lggr, logs := logger.TestObserved(t, zapcore.WarnLevel)
		report := newTestReport(t, lggr)

		require.NoError(t, report.Reserve())

		_, err := report.Deduct("ref1", ByResource(testUnitA, "", decimal.NewFromInt(1)))
		require.NoError(t, err)

		require.NoError(t, report.Settle("ref1", capabilities.ResponseMetadata{Metering: []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "frac", SpendUnit: testUnitGas, SpendValue: "0.000005", SpendValueInGasUnits: "0.5"},      // fractional: dropped
			{Peer2PeerID: "neg", SpendUnit: testUnitGas, SpendValue: "0.000005", SpendValueInGasUnits: "-100"},      // negative: dropped
			{Peer2PeerID: "ok", SpendUnit: testUnitGas, SpendValue: "0.000005", SpendValueInGasUnits: "5000000000"}, // valid
		}}))

		steps := report.FormatReport().Steps
		require.Len(t, steps, 1)
		require.Len(t, steps["ref1"].Nodes, 3)

		// only the valid detail feeds the aggregate
		assert.Equal(t, "5000000000.0000000000", steps["ref1"].AggSpendValue)
		assert.Equal(t, "0.0000000000", steps["ref1"].AggSpendValueCre)

		for _, node := range steps["ref1"].Nodes {
			if node.Peer_2PeerId == "ok" {
				assert.Equal(t, "5000000000", node.SpendValue)
			} else {
				assert.Equal(t, "0.000005", node.SpendValue, "peer %s", node.Peer_2PeerId)
				assert.Equal(t, "0.0000000000", node.SpendValueCre, "peer %s", node.Peer_2PeerId)
			}
		}

		// one warn per invalid detail, none for the valid one
		assert.Len(t, logs.All(), 2)
		assert.True(t, report.meteringMode)
	})
	t.Run("legacy gas spend without gas units falls back to shift 18", func(t *testing.T) {
		t.Parallel()

		const aptosUnit = "GAS.4741433654826277614"

		report := newTestReport(t, logger.Nop())

		require.NoError(t, report.Reserve())

		_, err := report.Deduct("ref1", ByResource(testUnitA, "", decimal.NewFromInt(1)))
		require.NoError(t, err)

		require.NoError(t, report.Settle("ref1", capabilities.ResponseMetadata{Metering: []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "abc", SpendUnit: aptosUnit, SpendValue: "0.001158"},
		}}))

		expected := map[string]*eventspb.MeteringReportStep{
			"ref1": {
				Nodes: []*eventspb.MeteringReportNodeDetail{
					{
						Peer_2PeerId:  "abc",
						SpendUnit:     aptosUnit,
						SpendValue:    "0.001158",     // no native value available: legacy field is preserved
						SpendValueCre: "0.0000000000", // Shift(18) on an 8-decimal token inflates by 10^10
					},
				},
				AggSpendValue:    "1158000000000000.0000000000",
				AggSpendUnit:     aptosUnit,
				AggSpendValueCre: "0.0000000000",
				CapdonN:          1,
				AggSpend: []*eventspb.AggregatedSpendDetail{
					{
						SpendValue:    "1158000000000000.0000000000",
						SpendUnit:     aptosUnit,
						SpendValueCre: "0.0000000000",
					},
				},
			},
		}

		assert.Equal(t, expected, report.FormatReport().Steps)
		assert.True(t, report.meteringMode)
	})
}

func Test_parseNativeGasUnits(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "integer", input: "115800", want: "115800"},
		{name: "large integer", input: "34440000000000", want: "34440000000000"},
		{name: "zero", input: "0", want: "0"},
		{name: "exponent with integer value rejected", input: "1e18", wantErr: true},
		{name: "max exponent rejected", input: "1e2147483647", wantErr: true},
		{name: "min exponent rejected", input: "1e-2147483648", wantErr: true},
		{name: "fractional rejected", input: "0.001158", wantErr: true},
		{name: "fractional wei rejected", input: "0.5", wantErr: true},
		{name: "negative rejected", input: "-5", wantErr: true},
		{name: "unparseable rejected", input: "????", wantErr: true},
		{name: "empty rejected", input: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := parseNativeGasUnits(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got.String())
		})
	}
}

func Test_Report_StepsWithoutRateCard(t *testing.T) {
	t.Parallel()

	// No rate card: executions run in metering mode — no credit accounting,
	// no spend limits — while steps and reports are still recorded.
	t.Run("deduct and settle record steps without touching the balance", func(t *testing.T) {
		t.Parallel()

		report := newTestReport(t, logger.Nop())

		require.NoError(t, report.Reserve())
		assert.True(t, report.isMeteringMode())

		limits, err := report.Deduct("step1", ByResource(testUnitA, "", decimal.NewFromInt(2)))
		require.NoError(t, err)
		assert.Empty(t, limits)

		require.NoError(t, report.Settle("step1", capabilities.ResponseMetadata{Metering: []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "node1", SpendUnit: testUnitA, SpendValue: "2"},
		}}))

		assert.Equal(t, "0", report.balance.GetSpent().String())
		assert.Contains(t, report.FormatReport().Steps, "step1")
	})
}

func Test_Report_EmitReceipt(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		rec := beholdertest.NewRecorder()
		var durableAttrs []any
		emitter := events.NewEmitter(rec, events.WithDurableEmit(func(_ context.Context, _ []byte, attrKVs ...any) error {
			durableAttrs = attrKVs
			return nil
		}))

		report := newTestReport(t, logger.Nop(), WithEventEmitter(emitter))

		require.NoError(t, report.Reserve())

		require.NoError(t, report.EmitReceipt(t.Context()))

		meteringEntity := fmt.Sprintf("%s.%s", events.ProtoPkg, events.MeteringReportEntity)
		assert.Equal(t, 1, rec.Len(t, "beholder_entity", meteringEntity))

		// Receipts must also reach the durable emitter, not just beholder.
		assert.Equal(t, []any{"source", "platform", "type", meteringEntity}, durableAttrs)

		messages := rec.Messages(t, "beholder_entity", meteringEntity)

		for _, msg := range messages {
			entity := msg.Attrs["beholder_entity"]
			if entity == fmt.Sprintf("%s.%s", events.ProtoPkg, events.MeteringReportEntity) {
				var report eventspb.MeteringReport
				require.NoError(t, proto.Unmarshal(msg.Body, &report))
				assert.Equal(t, testWorkflowID, report.Metadata.WorkflowID)
				assert.NotEmpty(t, report.Metadata.WorkflowExecutionID)
				assert.Equal(t, testAccountID, report.Metadata.WorkflowOwner)
			}
		}
	})

	t.Run("returns an error if not initialized", func(t *testing.T) {
		t.Parallel()

		report := newTestReport(t, logger.Nop())

		require.ErrorIs(t, report.EmitReceipt(t.Context()), ErrNoReserve)
	})
}

func Test_MeterReports(t *testing.T) {
	t.Parallel()

	workflowExecutionID1 := "exec1"
	capabilityCall1 := "ref1"

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		metrics := defaultMetrics(t)
		mrs := NewReports(testAccountID, testWorkflowID, logger.Nop(), defaultLabels, metrics, dummyRegistryAddress, dummyChainSelector, workflowV2)

		r, err := mrs.Start(t.Context(), workflowExecutionID1)
		require.NoError(t, err)

		require.NoError(t, r.Reserve())

		_, err = r.Deduct(capabilityCall1, ByResource(testUnitA, "", decimal.NewFromInt(1)))
		require.NoError(t, err)

		require.NoError(t, r.Settle(capabilityCall1, capabilities.ResponseMetadata{Metering: []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "1", SpendUnit: testUnitA, SpendValue: "0.8"},
			{Peer2PeerID: "2", SpendUnit: testUnitA, SpendValue: "0.9"},
			{Peer2PeerID: "3", SpendUnit: testUnitA, SpendValue: "1"},
			{Peer2PeerID: "4", SpendUnit: testUnitA, SpendValue: "1"},
		}}))
		require.NoError(t, mrs.End(t.Context(), workflowExecutionID1))
	})

	t.Run("happy path in metering mode", func(t *testing.T) {
		t.Parallel()

		metrics := defaultMetrics(t)
		mrs := NewReports(testAccountID, testWorkflowID, logger.Nop(), defaultLabels, metrics, dummyRegistryAddress, dummyChainSelector, workflowV2)

		r, err := mrs.Start(t.Context(), workflowExecutionID1)
		require.NoError(t, err)

		require.NoError(t, r.Reserve())

		_, err = r.Deduct(capabilityCall1, ByResource(testUnitA, "", decimal.NewFromInt(1)))
		require.NoError(t, err)

		require.NoError(t, r.Settle(capabilityCall1, capabilities.ResponseMetadata{Metering: []capabilities.MeteringNodeDetail{
			{Peer2PeerID: "1", SpendUnit: testUnitA, SpendValue: "1"},
			{Peer2PeerID: "2", SpendUnit: testUnitA, SpendValue: "1"},
			{Peer2PeerID: "3", SpendUnit: testUnitA, SpendValue: "1"},
			{Peer2PeerID: "4", SpendUnit: testUnitA, SpendValue: "1"},
		}}))
		require.NoError(t, mrs.End(t.Context(), workflowExecutionID1))
	})
}

func Test_MeterReports_Length(t *testing.T) {
	t.Parallel()

	em, err := monitoring.InitMonitoringResources()
	require.NoError(t, err)
	metrics := monitoring.NewWorkflowsMetricLabeler(metrics.NewLabeler(), em)
	mrs := NewReports("", "", logger.Nop(), defaultLabels, metrics, dummyRegistryAddress, dummyChainSelector, workflowV2)

	_, err = mrs.Start(t.Context(), "exec1")
	require.NoError(t, err)

	mr, err := mrs.Start(t.Context(), "exec2")
	require.NoError(t, err)

	_, err = mrs.Start(t.Context(), "exec3")
	require.NoError(t, err)
	assert.Equal(t, 3, mrs.Len())

	require.NoError(t, mr.Reserve())
	require.NoError(t, mrs.End(t.Context(), "exec2"))
	assert.Equal(t, 2, mrs.Len())
}

func Test_MeterReports_Start(t *testing.T) {
	t.Parallel()

	t.Run("can only start report once", func(t *testing.T) {
		t.Parallel()

		metrics := defaultMetrics(t)
		mrs := NewReports("", "", logger.Nop(), defaultLabels, metrics, dummyRegistryAddress, dummyChainSelector, workflowV2)

		_, err := mrs.Start(t.Context(), "exec1")
		require.NoError(t, err)

		_, err = mrs.Start(t.Context(), "exec1")
		require.ErrorIs(t, err, ErrReportExists)
	})
}

func Test_MeterReports_Get(t *testing.T) {
	t.Parallel()

	t.Run("returns when report exists", func(t *testing.T) {
		t.Parallel()

		metrics := defaultMetrics(t)
		mrs := NewReports("", "", logger.Nop(), defaultLabels, metrics, dummyRegistryAddress, dummyChainSelector, workflowV2)

		_, err := mrs.Start(t.Context(), "exec1")
		require.NoError(t, err)

		report, exists := mrs.Get("exec1")
		require.True(t, exists)
		require.NotEmpty(t, report)
	})

	t.Run("returns when no report exists", func(t *testing.T) {
		t.Parallel()

		metrics := defaultMetrics(t)
		mrs := NewReports("", "", logger.Nop(), defaultLabels, metrics, dummyRegistryAddress, dummyChainSelector, workflowV2)

		report, exists := mrs.Get("exec1")
		require.False(t, exists)
		require.Nil(t, report)
	})
}

func Test_MeterReports_End(t *testing.T) {
	t.Parallel()

	t.Run("can only end existing report", func(t *testing.T) {
		t.Parallel()

		metrics := defaultMetrics(t)
		mrs := NewReports("", "", logger.Nop(), defaultLabels, metrics, dummyRegistryAddress, dummyChainSelector, workflowV2)

		require.ErrorIs(t, mrs.End(t.Context(), "exec1"), ErrReportNotFound)
	})

	t.Run("cleans up report after emit", func(t *testing.T) {
		t.Parallel()

		metrics := defaultMetrics(t)
		mrs := NewReports("", "", logger.Nop(), defaultLabels, metrics, dummyRegistryAddress, dummyChainSelector, workflowV2)

		mr, err := mrs.Start(t.Context(), "exec1")
		require.NoError(t, err)
		assert.Len(t, mrs.reports, 1)

		require.NoError(t, mr.Reserve())
		require.NoError(t, mrs.End(t.Context(), "exec1"))
		assert.Empty(t, mrs.reports)
	})

	t.Run("cleans up report on failed emission", func(t *testing.T) {
		t.Parallel()

		emitErr := errors.New("emit failed")
		emitter := events.NewEmitter(beholdertest.NewRecorder(), events.WithDurableEmit(func(_ context.Context, _ []byte, _ ...any) error {
			return emitErr
		}))

		metrics := defaultMetrics(t)
		mrs := NewReports("", "", logger.Nop(), defaultLabels, metrics, dummyRegistryAddress, dummyChainSelector, workflowV2, WithEventEmitter(emitter))

		mr, err := mrs.Start(t.Context(), "exec1")
		require.NoError(t, err)
		assert.Len(t, mrs.reports, 1)

		require.NoError(t, mr.Reserve())
		require.ErrorIs(t, mrs.End(t.Context(), "exec1"), emitErr)
		assert.Empty(t, mrs.reports)
	})
}

func TestRatiosFromConfig(t *testing.T) {
	t.Parallel()

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		ratios, err := ratiosFromConfig(capabilities.CapabilityInfo{
			SpendTypes: []capabilities.CapabilitySpendType{
				capabilities.CapabilitySpendType(testUnitA),
				capabilities.CapabilitySpendType(testUnitB),
			},
		}, validConfig)

		require.NoError(t, err)
		require.Len(t, ratios, 2)

		assert.Contains(t, ratios, capabilities.CapabilitySpendType(testUnitA))
		assert.Contains(t, ratios, capabilities.CapabilitySpendType(testUnitB))
	})

	t.Run("automatic ratio of 1 for single spend type", func(t *testing.T) {
		t.Parallel()

		ratios, err := ratiosFromConfig(capabilities.CapabilityInfo{
			SpendTypes: []capabilities.CapabilitySpendType{capabilities.CapabilitySpendType(testUnitA)},
		}, nil)

		require.NoError(t, err)
		require.Len(t, ratios, 1)

		assert.Contains(t, ratios, capabilities.CapabilitySpendType(testUnitA))
	})

	t.Run("error when spend ratios key does not exist", func(t *testing.T) {
		t.Parallel()

		ratios, err := ratiosFromConfig(capabilities.CapabilityInfo{}, new(values.Map))
		require.ErrorIs(t, err, ErrInvalidRatios)
		assert.Empty(t, ratios)
	})

	t.Run("error when spend ratios fails to unwrap to map", func(t *testing.T) {
		t.Parallel()

		config := &values.Map{
			Underlying: map[string]values.Value{
				"spendRatios": &values.String{Underlying: "invalid"},
			},
		}

		ratios, err := ratiosFromConfig(capabilities.CapabilityInfo{}, config)
		require.ErrorIs(t, err, ErrInvalidRatios)
		assert.Empty(t, ratios)
	})

	t.Run("error when spend type is not in ratios map", func(t *testing.T) {
		t.Parallel()

		ratios, err := ratiosFromConfig(capabilities.CapabilityInfo{
			SpendTypes: []capabilities.CapabilitySpendType{
				capabilities.CapabilitySpendType(testUnitA),
				capabilities.CapabilitySpendType(testUnitC),
			},
		}, validConfig)

		require.ErrorIs(t, err, ErrInvalidRatios)
		assert.Empty(t, ratios)
	})

	t.Run("error when ratio contains invalid data type", func(t *testing.T) {
		t.Parallel()

		config, _ := values.NewMap(map[string]any{
			RatiosKey: map[string]any{
				testUnitA: "0.2",
				testUnitB: 5,
			},
		})

		ratios, err := ratiosFromConfig(capabilities.CapabilityInfo{
			SpendTypes: []capabilities.CapabilitySpendType{
				capabilities.CapabilitySpendType(testUnitA),
				capabilities.CapabilitySpendType(testUnitB),
			},
		}, config)

		require.ErrorIs(t, err, ErrInvalidRatios)
		assert.Empty(t, ratios)
	})

	t.Run("error when ratio contains invalid decimal", func(t *testing.T) {
		t.Parallel()

		config, _ := values.NewMap(map[string]any{
			RatiosKey: map[string]any{
				testUnitA: "invalid",
				testUnitB: "0.2",
			},
		})

		ratios, err := ratiosFromConfig(capabilities.CapabilityInfo{
			SpendTypes: []capabilities.CapabilitySpendType{
				capabilities.CapabilitySpendType(testUnitA),
				capabilities.CapabilitySpendType(testUnitB),
			},
		}, config)

		require.ErrorIs(t, err, ErrInvalidRatios)
		assert.Empty(t, ratios)
	})
}

func newTestReport(t *testing.T, lggr logger.Logger, opts ...Option) *Report {
	t.Helper()

	meteringReport, err := NewReport(t.Context(), defaultLabels, lggr, defaultMetrics(t), dummyRegistryAddress, dummyChainSelector, workflowV2, opts...)
	require.NoError(t, err)

	return meteringReport
}

func defaultMetrics(t *testing.T) *monitoring.WorkflowsMetricLabeler {
	em, err := monitoring.InitMonitoringResources()
	require.NoError(t, err)

	return monitoring.NewWorkflowsMetricLabeler(metrics.NewLabeler(), em)
}

func Test_Report_ConcurrentSettleAndFormatReport(t *testing.T) {
	t.Parallel()

	report := newTestReport(t, logger.Nop())
	require.NoError(t, report.Reserve())

	numRefs := 10
	for i := range numRefs {
		ref := fmt.Sprintf("ref%d", i)
		_, err := report.Deduct(ref, ByResource(testUnitA, "capID", decimal.NewFromInt(1)))
		require.NoError(t, err)
	}

	done := make(chan struct{})
	go func() {
		for range 100 {
			_ = report.FormatReport()
			_ = report.EmitReceipt(t.Context())
		}
		close(done)
	}()

	for i := range numRefs {
		go func(idx int) {
			ref := fmt.Sprintf("ref%d", idx)
			metadata := capabilities.ResponseMetadata{
				Metering: []capabilities.MeteringNodeDetail{
					{Peer2PeerID: "peer1", SpendUnit: testUnitA, SpendValue: "1"},
				},
			}
			_ = report.Settle(ref, metadata)
		}(i)
	}

	<-done
}
