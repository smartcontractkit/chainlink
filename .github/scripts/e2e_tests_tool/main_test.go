package main

import (
	"testing"
)

func TestFilterTestsByID(t *testing.T) {
	tests := []Test{
		{ID: "run_all_in_ocr_tests_go", Name: "Run all ocr_tests.go", TestType: "docker"},
		{ID: "run_all_in_ocr2_tests_go", Name: "Run TestOCRv2Request in ocr2_test.go", TestType: "docker"},
		{ID: "run_all_in_ocr3_tests_go", Name: "Run TestOCRv2Basic in ocr2_test.go", TestType: "k8s_remote_runner"},
	}

	cases := []struct {
		description string
		inputIDs    string
		expectedLen int
	}{
		{"Filter by single ID", "run_all_in_ocr_tests_go", 1},
		{"Filter by multiple IDs", "run_all_in_ocr_tests_go,run_all_in_ocr2_tests_go", 2},
		{"Wildcard to include all", "*", 3},
		{"Empty ID string to include all", "", 3},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			filtered := FilterTests(tests, "", "", "", c.inputIDs)
			if len(filtered) != c.expectedLen {
				t.Errorf("FilterTests(%s) returned %d tests, expected %d", c.description, len(filtered), c.expectedLen)
			}
		})
	}
}

func TestFilterTestsIntegration(t *testing.T) {
	tests := []Test{
		{ID: "run_all_in_ocr_tests_go", Name: "Run all ocr_tests.go", TestType: "docker", Trigger: []string{"nightly"}},
		{ID: "run_all_in_ocr2_tests_go", Name: "Run TestOCRv2Request in ocr2_test.go", TestType: "docker", Trigger: []string{"push"}},
		{ID: "run_all_in_ocr3_tests_go", Name: "Run TestOCRv2Basic in ocr2_test.go", TestType: "k8s_remote_runner", Trigger: []string{"push"}},
	}

	cases := []struct {
		description   string
		inputNames    string
		inputTrigger  string
		inputTestType string
		inputIDs      string
		expectedLen   int
	}{
		{"Filter by test type and ID", "", "", "docker", "run_all_in_ocr2_tests_go", 1},
		{"Filter by trigger and test type", "", "push", "docker", "*", 1},
		{"No filters applied", "", "", "", "*", 3},
		{"Filter mismatching all criteria", "", "nightly", "k8s_remote_runner", "run_all_in_ocr_tests_go", 0},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			filtered := FilterTests(tests, c.inputNames, c.inputTrigger, c.inputTestType, c.inputIDs)
			if len(filtered) != c.expectedLen {
				t.Errorf("FilterTests(%s) returned %d tests, expected %d", c.description, len(filtered), c.expectedLen)
			}
		})
	}
}
