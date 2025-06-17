package safe

import (
	"math"
	"math/big"
	"testing"
)

func TestIntToUint64(t *testing.T) {
	// Test Case 1: Positive integer within the range of uint64
	t.Run("PositiveInteger", func(t *testing.T) {
		input := 12345
		expectedOutput := uint64(12345)
		output, err := IntToUint64(input)
		if err != nil {
			t.Errorf("Test Case PositiveInteger failed: expected no error, got %v", err)
		}
		if output != expectedOutput {
			t.Errorf("Test Case PositiveInteger failed: expected %d, got %d", expectedOutput, output)
		}
	})

	// Test Case 2: Negative integer
	t.Run("NegativeInteger", func(t *testing.T) {
		input := -10
		_, err := IntToUint64(input)
		if err == nil {
			t.Errorf("Test Case NegativeInteger failed: expected an error, got nil")
		}
		expectedErrorMessage := "cannot convert negative int to uint64"
		if err.Error() != expectedErrorMessage {
			t.Errorf("Test Case NegativeInteger failed: expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
		}
	})

	// Test Case 3: math.MaxInt64 + n will be interpreted as negative
	// because the first byte is non-zero
	t.Run("IntegerOverflow", func(t *testing.T) {
		maxInt := big.NewInt(math.MaxInt64)
		overflowInt2 := new(big.Int).Add(maxInt, big.NewInt(1))
		_, err := IntToUint64(int(overflowInt2.Int64()))
		if err == nil {
			t.Errorf("Test Case IntegerOverflow failed: expected an error, got nil")
		}
		expectedErrorMessage := "cannot convert negative int to uint64"
		if err.Error() != expectedErrorMessage {
			t.Errorf("Test Case IntegerOverflow failed: expected error message '%s', got '%s'", expectedErrorMessage, err.Error())
		}
	})

	// Test Case 4: Zero
	t.Run("Zero", func(t *testing.T) {
		input := 0
		expectedOutput := uint64(0)
		output, err := IntToUint64(input)
		if err != nil {
			t.Errorf("Test Case Zero failed: expected no error, got %v", err)
		}
		if output != expectedOutput {
			t.Errorf("Test Case Zero failed: expected %d, got %d", expectedOutput, output)
		}
	})
}
