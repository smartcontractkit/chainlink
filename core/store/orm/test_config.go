package orm

// TestConfig is a configuration store used for testing
type TestConfig struct {
	*Config
}

// NewTestConfig returns a new TestConfig
func NewTestConfig() *TestConfig {
	store := NewBootstrapConfigStore()
	return &TestConfig{
		Config: NewConfig(store),
	}
}

// SessionSecret returns a static session secret
func (c TestConfig) SessionSecret() ([]byte, error) {
	return []byte("clsession_test_secret"), nil
}
