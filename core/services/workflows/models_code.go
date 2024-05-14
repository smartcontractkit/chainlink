package workflows

type codeConfig struct {
	TypeMap map[string]string  `json:"type_map" jsonschema:"required"`
	Config  map[string]mapping `json:"config" jsonschema:"required"`
}
