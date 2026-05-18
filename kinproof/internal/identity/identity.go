package identity
import ("encoding/json"; "os"; "path/filepath"; "time"; "kinproof/internal/hd")
type ID struct { Address string `json:"address"`; PrivateKey string `json:"privateKey"`; Index uint32 `json:"index"`; CreatedAt int64 `json:"createdAt"` }
func Rotate(h *hd.SovereignHD, dir string, i uint32) (*ID, error) {
	addr, priv, _ := h.DeriveEphemeral(i)
	id := &ID{addr, priv, i, time.Now().UnixMilli()}
	os.MkdirAll(dir, 0700)
	data, _ := json.MarshalIndent(id, "", "  ")
	os.WriteFile(filepath.Join(dir, "latest_identity.json"), data, 0600)
	return id, nil
}
