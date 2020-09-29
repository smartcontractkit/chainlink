package serialization_test

import (
	"math"
	"math/big"
	"reflect"
	"testing"

	"github.com/smartcontractkit/offchain-reporting/lib/networking"
	"github.com/smartcontractkit/offchain-reporting/lib/offchainreporting/internal/protocol"
	"github.com/smartcontractkit/offchain-reporting/lib/offchainreporting/internal/protocol/observation"
	"github.com/smartcontractkit/offchain-reporting/lib/offchainreporting/internal/serialization"
	"github.com/smartcontractkit/offchain-reporting/lib/offchainreporting/internal/signature"
	"github.com/smartcontractkit/offchain-reporting/lib/offchainreporting/types"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/stretchr/testify/require"
)

var obs, _ = observation.MakeObservation(big.NewInt(1))

func Test_Serialization_PropertyTests(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("MessageNewEpoch", prop.ForAll(
		func(m protocol.Message) bool {
			b, err := serialization.Serialize(m)
			require.NoError(t, err)
			m2, err := serialization.Deserialize(b)
			require.NoError(t, err)
			return m == m2
		},
		gen.Struct(reflect.TypeOf(&protocol.MessageNewEpoch{}), map[string]gopter.Gen{
			"Epoch": gen.Int(),
		}),
	))

	properties.Property("MessageObserveReq", prop.ForAll(
		func(m protocol.Message) bool {
			b, err := serialization.Serialize(m)
			require.NoError(t, err)
			m2, err := serialization.Deserialize(b)
			require.NoError(t, err)
			return m == m2
		},
		gen.Struct(reflect.TypeOf(&protocol.MessageObserveReq{}), map[string]gopter.Gen{
			"Round": gen.Int(),
			"Epoch": gen.Int(),
		}),
	))

	properties.Property("MessageObserve", prop.ForAll(
		func(m protocol.Message) bool {
			b, err := serialization.Serialize(m)
			require.NoError(t, err)
			m2, err := serialization.Deserialize(b)
			require.NoError(t, err)
			return m.(protocol.MessageObserve).Equal(m2.(protocol.MessageObserve))
		},
		gen.Struct(reflect.TypeOf(&protocol.MessageObserve{}), map[string]gopter.Gen{
			"Epoch": gen.Int(),
			"Round": gen.Int(),
			"Obs":   genObs(),
		}),
	))

	properties.Property("MessageReportReq", prop.ForAll(
		func(m protocol.MessageReportReq) bool {
			b, err := serialization.Serialize(m)
			require.NoError(t, err)
			msg2, err := serialization.Deserialize(b)
			m2 := msg2.(protocol.MessageReportReq)
			require.NoError(t, err)
			return reflect.DeepEqual(m, m2)
		},
		gen.Struct(reflect.TypeOf(&protocol.MessageReportReq{}), map[string]gopter.Gen{
			"Round":        gen.Int(),
			"Observations": gen.SliceOf(genObs()),
			"Epoch":        gen.Int(),
		}),
	))

	properties.Property("MessageReport", prop.ForAll(
		func(m protocol.MessageReport) bool {
			b, err := serialization.Serialize(m)
			require.NoError(t, err)
			msg2, err := serialization.Deserialize(b)
			m2 := msg2.(protocol.MessageReport)
			require.NoError(t, err)
			return m.Equals(m2)
		},
		gen.Struct(reflect.TypeOf(&protocol.MessageReport{}), map[string]gopter.Gen{
			"Epoch":          gen.Int(),
			"Round":          gen.Int(),
			"ContractReport": genContractReport(),
		}),
	))

	properties.Property("MessageFinal", prop.ForAll(
		func(m protocol.MessageFinal) bool {
			b, err := serialization.Serialize(m)
			require.NoError(t, err)
			msg2, err := serialization.Deserialize(b)
			m2 := msg2.(protocol.MessageFinal)
			require.NoError(t, err)
			return m.Equals(m2)
		},
		gen.Struct(reflect.TypeOf(&protocol.MessageFinal{}), map[string]gopter.Gen{
			"Epoch":  gen.Int(),
			"Leader": genOracleID(),
			"Round":  gen.Int(),
			"Report": genContractReportWithSigs(),
		}),
	))

	properties.Property("MessageFinalEcho", prop.ForAll(
		func(m protocol.MessageFinalEcho) bool {
			b, err := serialization.Serialize(m)
			require.NoError(t, err)
			msg2, err := serialization.Deserialize(b)
			m2 := msg2.(protocol.MessageFinalEcho)
			require.NoError(t, err)
			return m.Equals(m2)
		},
		gen.Struct(reflect.TypeOf(&protocol.MessageFinalEcho{}), map[string]gopter.Gen{
			"MessageFinal": gen.Struct(reflect.TypeOf(&protocol.MessageFinal{}), map[string]gopter.Gen{
				"Epoch":  gen.Int(),
				"Leader": genOracleID(),
				"Round":  gen.Int(),
				"Report": genContractReportWithSigs(),
			}),
		}),
	))
	properties.TestingRun(t)
}

func genContractReport() gopter.Gen {
	return gen.Struct(reflect.TypeOf(&protocol.ContractReport{}), map[string]gopter.Gen{
		"Ctx":    genContext(),
		"Values": gen.SliceOf(genOracleValue()),
		"Sig":    genSig(),
	})
}

func genContractReportWithSigs() gopter.Gen {
	return gen.Struct(reflect.TypeOf(&protocol.ContractReportWithSignatures{}), map[string]gopter.Gen{
		"Ctx":        genContext(),
		"Values":     gen.SliceOf(genOracleValue()),
		"Sig":        genSig(),
		"Signatures": gen.SliceOf(genSig()),
	})
}

func genContext() gopter.Gen {
	return gen.Struct(reflect.TypeOf(&signature.ReportingContext{}), map[string]gopter.Gen{
		"ConfigDigest": genConfigDigest(),
		"Epoch":        gen.Int(),
		"Round":        gen.Int(),
	})
}

func genOracleValue() gopter.Gen {
	return gen.Struct(reflect.TypeOf(&protocol.OracleValue{}), map[string]gopter.Gen{
		"ID":    genOracleID(),
		"Value": observation.GenObservationValue(),
	})
}

func genObs() gopter.Gen {
	return gen.Struct(reflect.TypeOf(&protocol.Observation{}), map[string]gopter.Gen{
		"Ctx": genCtx(),
		"OracleID": genOracleID().SuchThat(func(id types.OracleID) bool {
			return id >= 0 && id < math.MaxUint32
		}),
		"Value": observation.GenObservationValue(),
		"Sig":   genSig(),
	})
}

func genSig() gopter.Gen {
	return gen.SliceOf(genByte())
}

func genByte() gopter.Gen {
	return func(p *gopter.GenParameters) *gopter.GenResult {
		b := make([]byte, 1)
		p.Rng.Read(b)
		return gopter.NewGenResult(b[0], gopter.NoShrinker)
	}
}

func genCtx() gopter.Gen {
	return gen.Struct(reflect.TypeOf(&signature.ReportingContext{}), map[string]gopter.Gen{
		"ConfigDigest": genConfigDigest(),
		"Epoch":        gen.Int(),
		"Round":        gen.Int(),
	})
}

func genOracleID() gopter.Gen {
	return func(p *gopter.GenParameters) *gopter.GenResult {
		i := int(uint32(p.NextInt64()))
		return gopter.NewGenResult(types.OracleID(i), gopter.NoShrinker)
	}
}

func genConfigDigest() gopter.Gen {
	return func(p *gopter.GenParameters) *gopter.GenResult {
		b := make([]byte, 27)
		p.Rng.Read(b)
		return gopter.NewGenResult(types.BytesToConfigDigest(b), gopter.NoShrinker)
	}
}

func Test_Serialization_UnknownStruct(t *testing.T) {
	_, err := serialization.Serialize(protocol.XXXUnknownMessageType{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "Unable to serialize")
}

func Test_Serialization_MaxLen(t *testing.T) {
	configDigest := types.ConfigDigest{}
	for i := range configDigest {
		configDigest[i] = byte(255)
	}
	sig := make([]byte, 64)
	for i := range sig {
		sig[i] = byte(255)
	}
	observations := make([]protocol.Observation, 31)
	for i := range observations {
		v, err := observation.MakeObservation(big.NewInt(math.MaxInt64))
		require.NoError(t, err)
		observations[i] = protocol.Observation{
			Ctx: signature.ReportingContext{
				ConfigDigest: configDigest,
				Epoch:        math.MaxInt64,
				Round:        math.MaxInt64,
			},
			OracleID: math.MaxUint32,
			Value:    v,
			Sig:      sig,
		}
	}

	hugeMsg := protocol.MessageReportReq{
		Round:        math.MaxInt64,
		Observations: observations,
	}

	b, err := serialization.Serialize(hugeMsg)
	require.NoError(t, err)
	require.Less(t, len(b), 2*networking.MaxMsgLength)
}

func Test_Serialize_Deserialize_Static(t *testing.T) {
	configDigest := types.ConfigDigest{}
	for i := range configDigest {
		configDigest[i] = byte(255)
	}
	sigs := make([][]byte, 31)
	for i := range sigs {
		sig := make([]byte, 64)
		for i := range sig {
			sig[i] = byte(255)
		}
		sigs[i] = sig
	}
	values := make([]protocol.OracleValue, 31)
	for i := range values {
		v, err := observation.MakeObservation(big.NewInt(math.MaxInt64))
		require.NoError(t, err)
		values[i] = protocol.OracleValue{
			ID:    math.MaxUint32,
			Value: v,
		}
	}

	hugeMsg := protocol.MessageFinal{
		Epoch:  math.MaxInt64,
		Round:  math.MaxInt64,
		Leader: math.MaxUint32,
		Report: protocol.ContractReportWithSignatures{
			ContractReport: protocol.ContractReport{
				Ctx: signature.ReportingContext{
					ConfigDigest: configDigest,
					Epoch:        math.MaxInt64,
					Round:        math.MaxInt64,
				},
				Values: values,
				Sig:    sigs[0],
			},
			Signatures: sigs,
		},
	}

	b, err := serialization.Serialize(hugeMsg)
	require.NoError(t, err)
	require.Less(t, len(b), 2*networking.MaxMsgLength)

	deserialized, err := serialization.Deserialize(b)
	require.NoError(t, err)
	require.Equal(t, hugeMsg, deserialized)
}
