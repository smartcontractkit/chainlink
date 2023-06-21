package functions

const (
	LocationInline     = 0
	LocationRemote     = 1
	LanguageJavaScript = 0
)

type RequestData struct {
	Source          string   `cbor:"source"`
	Language        int      `cbor:"language"`
	CodeLocation    int      `cbor:"codeLocation"`
	Secrets         []byte   `cbor:"secrets"`
	SecretsLocation int      `cbor:"secretsLocation"`
	Args            []string `cbor:"args"`
}

type AdapterRequestData struct {
	Source          string   `json:"source"`
	Language        int      `json:"language"`
	CodeLocation    int      `json:"codeLocation"`
	Secrets         string   `json:"secrets"`
	SecretsLocation int      `json:"secretsLocation"`
	Args            []string `json:"args"`
}
