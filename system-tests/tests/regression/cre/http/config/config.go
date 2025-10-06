package config

type Config struct {
	AuthorizedKey string `json:"authorizedKey"`
	URL           string `json:"url"`
	TestCase      string `json:"testCase"` // Identifies which test case to run
}
