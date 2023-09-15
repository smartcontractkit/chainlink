package offchain

import (
	"encoding/hex"
	"encoding/json"
)

type SignedReport struct {
	Report    []byte
	Rs        [][32]byte
	Ss        [][32]byte
	Vs        [32]byte
	RequestID [32]byte
	Result    []byte
	Error     []byte
}

type TransmitterHook interface {
	Transmit(signedReport *SignedReport)
	ReportChannel() chan *SignedReport
}

type transmitterHook struct {
	chReports chan *SignedReport
}

func NewTransmitterHook() TransmitterHook {
	return &transmitterHook{chReports: make(chan *SignedReport)}
}

func (t *transmitterHook) Transmit(signedReport *SignedReport) {
	t.chReports <- signedReport
}

func (t *transmitterHook) ReportChannel() chan *SignedReport {
	return t.chReports
}

func (f *SignedReport) MarshalJSON() ([]byte, error) {
	rsStr := []string{}
	ssStr := []string{}
	for id, r := range f.Rs {
		rsStr = append(rsStr, "0x"+hex.EncodeToString(r[:]))
		ssStr = append(ssStr, "0x"+hex.EncodeToString(f.Ss[id][:]))
	}
	return json.Marshal(&struct {
		Report    string   `json:"report"`
		Rs        []string `json:"rs"`
		Ss        []string `json:"ss"`
		Vs        string   `json:"vs"`
		RequestID string   `json:"requestId"`
		Result    string   `json:"result"`
		Error     string   `json:"error"`
	}{
		Report:    "0x" + hex.EncodeToString(f.Report),
		Rs:        rsStr,
		Ss:        ssStr,
		Vs:        "0x" + hex.EncodeToString(f.Vs[:]),
		RequestID: "0x" + hex.EncodeToString(f.RequestID[:]),
		Result:    "0x" + hex.EncodeToString(f.Result),
		Error:     "0x" + hex.EncodeToString(f.Error),
	})
}
