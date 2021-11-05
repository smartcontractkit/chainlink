package configtest

const (
	// SessionSecret is the hardcoded secret solely used for test
	SessionSecret = "clsession_test_secret"
)

type MockSecretGenerator struct{}

func (m MockSecretGenerator) Generate(string) ([]byte, error) {
	return []byte(SessionSecret), nil
}
