package config

// Validated configurations impose constraints that must be checked.
type Validated interface {
	// ValidateConfig returns nil if the config is valid, otherwise an error describing why it is invalid.
	//
	// For implementations:
	//  - Use package multierr to accumulate all errors, rather than returning the first encountered.
	//  - If an anonymous field also implements ValidateConfig(), it must be called explicitly!
	ValidateConfig() error
}

type SecretType int

const (
	DBSecretType SecretType = iota
)

type Secret interface {
	Validated
	Type() SecretType
	SetEnabled(on bool)
	Enabled() bool
}
