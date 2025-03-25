package types

// k8s cost attribution
type TeamInput struct {
	Team       string `toml:"team" validate:"required"`
	Product    string `toml:"product" validate:"required"`
	CostCenter string `toml:"cost_center" validate:"required"`
	Component  string `toml:"component" validate:"required"`
}
