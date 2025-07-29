package infra

type Type = string

const (
	CRIB   Type = "crib"
	Docker Type = "docker"
)

type CribProvider = string

const (
	AWS  CribProvider = "aws"
	Kind CribProvider = "kind"
)

type Input struct {
	Type string     `toml:"type" validate:"oneof=crib docker"`
	CRIB *CRIBInput `toml:"crib"`
}

type CRIBInput struct {
	Namespace string `toml:"namespace" validate:"required"`
	// absolute path to the folder with CRIB CRE
	FolderLocation string `toml:"folder_location" validate:"required"`
	Provider       string `toml:"provider" validate:"oneof=aws kind"`
	// required for cost attribution in AWS
	TeamInput *Team `toml:"team_input" validate:"required_if=Provider aws"`
}

// k8s cost attribution
type Team struct {
	Team       string `toml:"team" validate:"required"`
	Product    string `toml:"product" validate:"required"`
	CostCenter string `toml:"cost_center" validate:"required"`
	Component  string `toml:"component" validate:"required"`
}
