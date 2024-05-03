package workflows

type codeConfig struct {
	TypeMap map[string]stepDefinitionType `json:"type_map" jsonschema:"required"`
	Config  map[string]mapping            `json:"config" jsonschema:"required"`
}
