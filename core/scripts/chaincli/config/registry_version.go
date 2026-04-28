package config

// RegistryVersion identifies a keeper registry contract version for chaincli tooling.
type RegistryVersion int32

const (
	RegistryVersion_2_0 RegistryVersion = 4
	RegistryVersion_2_1                 = 5
)

func (rv RegistryVersion) String() string {
	switch rv {
	case RegistryVersion_2_0:
		return "v2.0"
	case RegistryVersion_2_1:
		return "v2.1"
	default:
		return "unknown registry version"
	}
}
