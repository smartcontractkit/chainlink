package utils

import (
	"strconv"
)

// TomlFloat32 represents float32 values for TOML
type TomlFloat32 float32

// UnmarshalText parses the value as a proper float32
func (t *TomlFloat32) UnmarshalText(text []byte) error {
	f32, err := strconv.ParseFloat(string(text), 32)
	if err != nil {
		return err
	}

	*t = TomlFloat32(f32)

	return nil
}

// TomlFloat64 represents float64 values for TOML
type TomlFloat64 float64

// UnmarshalText parses the value as a proper float64
func (t *TomlFloat64) UnmarshalText(text []byte) error {
	f32, err := strconv.ParseFloat(string(text), 64)
	if err != nil {
		return err
	}

	*t = TomlFloat64(f32)

	return nil
}
