package cmd

import (
	"testing"
)

func TestMatchTestNameInCmd(t *testing.T) {
	tests := []struct {
		cmd      string
		testName string
		expected bool
	}{
		{"go test -test.run ^TestExample$", "TestExample", true},
		{"go test -test.run ^TestExample$", "TestAnother", false},
		{"go test -test.run ^TestExample$ -v", "TestExample", true},
		{"go test -test.run ^TestExamplePart$", "TestExample", false},
		{"go test -test.run ^TestWithNumbers123$", "TestWithNumbers123", true},
		{"go test -test.run ^Test_With_Underscores$", "Test_With_Underscores", true},
		{"go test -test.run ^Test-With-Dash$", "Test-With-Dash", true},
		{"go test -test.run ^TestWithSpace Space$", "TestWithSpace Space", true},
		{"go test -test.run ^TestWithNewline\nNewline$", "TestWithNewline\nNewline", true},
		{"go test -test.run ^TestOne$|^TestTwo$", "TestOne", true},
		{"go test -test.run ^TestOne$|^TestTwo$", "TestTwo", true},
		{"go test -test.run ^TestOne$|^TestTwo$", "TestThree", false},
		{"go test -test.run TestOne|TestTwo", "TestTwo", true},
		{"go test -test.run TestOne|TestTwo", "TestOne", true},
		{"go test -test.run TestOne|TestTwo|TestThree", "TestFour", false},
		{"go test -test.run ^TestOne$|TestTwo$", "TestTwo", true},
		{"go test -test.run ^TestOne$|TestTwo|TestThree$", "TestThree", true},
		{"go test -test.run TestOne|TestTwo|TestThree", "TestOne", true},
		{"go test -test.run TestOne|TestTwo|TestThree", "TestThree", true},
		{"go test -test.run ^TestA$|^TestB$|^TestC$", "TestA", true},
		{"go test -test.run ^TestA$|^TestB$|^TestC$", "TestB", true},
		{"go test -test.run ^TestA$|^TestB$|^TestC$", "TestD", false},
		{"go test -test.run TestA|^TestB$|TestC", "TestB", true},
		{"go test -test.run ^TestA|^TestB|TestC$", "TestA", true},
		{"go test -test.run ^TestA|^TestB|TestC$", "TestC", true},
		{"go test -test.run ^TestA|^TestB|TestC$", "TestD", false},
	}

	for _, tt := range tests {
		result := matchTestNameInCmd(tt.cmd, tt.testName)
		if result != tt.expected {
			t.Errorf("matchTestNameInCmd(%s, %s) = %t; expected %t", tt.cmd, tt.testName, result, tt.expected)
		}
	}
}
