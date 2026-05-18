package intent
import ("encoding/hex"; "encoding/json"; "github.com/ethereum/go-ethereum/crypto"; "kinproof/internal/proof")
func Sign(pkHex string, e proof.ExecutionEnvelope) (string, error) {
	pk, _ := hex.DecodeString(pkHex)
	priv, _ := crypto.ToECDSA(pk)
	msg, _ := json.Marshal(e)
	hash := crypto.Keccak256Hash(msg)
	sig, _ := crypto.Sign(hash.Bytes(), priv)
	return hex.EncodeToString(sig), nil
}
