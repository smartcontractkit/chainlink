package tomlutils

import (
	"strconv"
)

// Float32 represents float32 values for TOML
type Float32 float32

// UnmarshalText parses the value as a proper float32
func (t *Float32) UnmarshalText(text []byte) error {
	f32, err := strconv.ParseFloat(string(text), 32)
	if err != nil {
		return err
	}

	*t = Float32(f32)

	return nil
}

// Float64 represents float64 values for TOML
type Float64 float64

// UnmarshalText parses the value as a proper float64
func (t *Float64) UnmarshalText(text []byte) error {
	f32, err := strconv.ParseFloat(string(text), 64)
	if err != nil {
		return err
	}

	*t = Float64(f32)

	return nil
}
