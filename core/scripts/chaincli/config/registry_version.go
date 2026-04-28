package config

import "fmt"

// RegistryVersion identifies a keeper registry contract version for chaincli tooling.
type RegistryVersion int32

const (
	RegistryVersion_1_0 RegistryVersion = iota
	RegistryVersion_1_1
	RegistryVersion_1_2
	RegistryVersion_1_3
	RegistryVersion_2_0
	RegistryVersion_2_1
)

func (rv RegistryVersion) String() string {
	switch rv {
	case RegistryVersion_1_0, RegistryVersion_1_1, RegistryVersion_1_2, RegistryVersion_1_3:
		return fmt.Sprintf("v1.%d", rv)
	case RegistryVersion_2_0:
		return "v2.0"
	case RegistryVersion_2_1:
		return "v2.1"
	default:
		return "unknown registry version"
	}
}
