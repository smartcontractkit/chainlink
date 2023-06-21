package functions

const (
	LocationInline     = 0
	LocationRemote     = 1
	LanguageJavaScript = 0
)

type RequestData struct {
	Source          string   `json:"source" cbor:"source"`
	Language        int      `json:"language" cbor:"language"`
	CodeLocation    int      `json:"codeLocation" cbor:"codeLocation"`
	Secrets         []byte   `json:"secrets" cbor:"secrets"`
	SecretsLocation int      `json:"secretsLocation" cbor:"secretsLocation"`
	Args            []string `json:"args" cbor:"args"`
}
