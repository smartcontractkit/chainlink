package types

type InfraType = string

const (
	InfraType_CRIB   InfraType = "crib"
	InfraType_Docker InfraType = "docker"
)

type InfraDetails struct {
	InfraType InfraType
	Namespace string
}
