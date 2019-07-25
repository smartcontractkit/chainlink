package config

type MockConfig struct {
	Config
}

func NewMock() *MockConfig {
	return &MockConfig{
		Config: NewConfig(),
	}
}

func (c MockConfig) SessionSecret() ([]byte, error) {
	return []byte("clsession_test_secret"), nil
}
