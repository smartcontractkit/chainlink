package proof
import ("crypto/sha256"; "encoding/hex"; "encoding/json"; "fmt")
type ExecutionEnvelope struct {
	ExecutionID string `json:"executionId"`
	IdentityRoot string `json:"identityRoot"`
	EphemeralAddr string `json:"ephemeralAddr"`
	Intent string `json:"intent"`
	Nonce string `json:"nonce"`
}
func CreateEnvelope(root, addr, intent, nonce string, payload interface{}) (ExecutionEnvelope, error) {
	p, _ := json.Marshal(payload)
	h := sha256.Sum256(p)
	id := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%x|%s", root, addr, h, nonce)))
	return ExecutionEnvelope{hex.EncodeToString(id[:]), root, addr, intent, nonce}, nil
}
