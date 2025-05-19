package plugins

import (
	"os"

	"github.com/hashicorp/go-envparse"
)

// ParseEnvFile returns a slice of key/value pairs parsed from the file at filepath.
// As a special case, empty filepath returns nil without error.
func ParseEnvFile(filepath string) ([]string, error) {
	if filepath == "" {
		return nil, nil
	}
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()
	m, err := envparse.Parse(f)
	if err != nil {
		return nil, err
	}
	r := make([]string, 0, len(m))
	for k, v := range m {
		r = append(r, k+"="+v)
	}
	return r, nil
}
