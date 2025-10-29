package config

type Config struct {
	URL      string `json:"url"`
	TestCase string `json:"testCase"`         // Identifies which test case to run
	Method   string `json:"method,omitempty"` // HTTP method for the test
	Body     string `json:"body,omitempty"`
}
