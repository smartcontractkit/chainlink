package functions

type RequestData struct {
	Source          string   `json:"source" cbor:"source"`
	Language        int      `json:"language" cbor:"language"`
	CodeLocation    int      `json:"codeLocation" cbor:"codeLocation"`
	Secrets         string   `json:"secrets" cbor:"secrets"`
	SecretsLocation int      `json:"secretsLocation" cbor:"secretsLocation"`
	Args            []string `json:"args" cbor:"args"`
}
