package config

type Password interface {
	Keystore() string
	VRF() string
}
