package interfacetests

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/smartcontractkit/libocr/commontypes"
)

type BasicTester interface {
	Setup(t *testing.T)
	Name() string
	GetAccountBytes(i int) []byte
}

type testcase struct {
	name string
	test func(t *testing.T)
}

func runTests(t *testing.T, tester BasicTester, tests []testcase) {
	for _, test := range tests {
		t.Run(test.name+" for "+tester.Name(), func(t *testing.T) {
			tester.Setup(t)
			test.test(t)
		})
	}
}

type InnerTestStruct struct {
	I int
	S string
}

type MidLevelTestStruct struct {
	FixedBytes [2]byte
	Inner      InnerTestStruct
}

type TestStruct struct {
	Field          *int32
	DifferentField string
	OracleID       commontypes.OracleID
	OracleIDs      [32]commontypes.OracleID
	Account        []byte
	Accounts       [][]byte
	BigField       *big.Int
	NestedStruct   MidLevelTestStruct
}

type TestStructWithExtraField struct {
	TestStruct
	ExtraField int
}

type TestStructMissingField struct {
	DifferentField string
	OracleID       commontypes.OracleID
	OracleIDs      [32]commontypes.OracleID
	Account        []byte
	Accounts       [][]byte
	BigField       *big.Int
	NestedStruct   MidLevelTestStruct
}

// compatibleTestStruct has fields in a different order
type compatibleTestStruct struct {
	Account        []byte
	Accounts       [][]byte
	BigField       *big.Int
	DifferentField string
	Field          int32
	NestedStruct   MidLevelTestStruct
	OracleID       commontypes.OracleID
	OracleIDs      [32]commontypes.OracleID
}

type LatestParams struct {
	I int
}

type FilterEventParams struct {
	Field int32
}

func CreateTestStruct(i int, tester BasicTester) TestStruct {
	s := fmt.Sprintf("field%v", i)
	fv := int32(i)
	return TestStruct{
		Field:          &fv,
		DifferentField: s,
		OracleID:       commontypes.OracleID(i + 1),
		OracleIDs:      [32]commontypes.OracleID{commontypes.OracleID(i + 2), commontypes.OracleID(i + 3)},
		Account:        tester.GetAccountBytes(i + 3),
		Accounts:       [][]byte{tester.GetAccountBytes(i + 4), tester.GetAccountBytes(i + 5)},
		BigField:       big.NewInt(int64((i + 1) * (i + 2))),
		NestedStruct: MidLevelTestStruct{
			FixedBytes: [2]byte{uint8(i), uint8(i + 1)},
			Inner: InnerTestStruct{
				I: i,
				S: s,
			},
		},
	}
}
