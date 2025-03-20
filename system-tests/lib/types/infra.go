package types

type InfraType = string

const (
	CRIB   InfraType = "crib"
	Docker InfraType = "docker"
)

type CribProvider = string

const (
	AWS  CribProvider = "aws"
	Kind CribProvider = "kind"
)

type InfraInput struct {
	InfraType string     `toml:"type" validate:"oneof=crib docker"`
	CRIB      *CRIBInput `toml:"crib"`
}

type CRIBInput struct {
	Namespace string `toml:"namespace" validate:"required"`
	// absolute path to the folder with CRIB CRE
	FolderLocation string `toml:"folder_location" validate:"required"`
	Provider       string `toml:"provider" validate:"oneof=aws kind"`
	// required for cost attribution in AWS
	TeamInput *TeamInput `toml:"team_input" validate:"required_if=Provider aws"`
}
