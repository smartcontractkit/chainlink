package model

import "encoding/json"

// CommitPluginReport is placed here for reference of shared readers structure.
type CommitPluginReport struct{}

type CommitPluginObservation struct {
	NodeID  NodeID               `json:"nodeID"`
	NewMsgs []CCIPMsgBaseDetails `json:"newMsgs"`
}

func NewCommitPluginObservation(nodeID NodeID, newMsgs []CCIPMsgBaseDetails) CommitPluginObservation {
	return CommitPluginObservation{
		NodeID:  nodeID,
		NewMsgs: newMsgs,
	}
}

func (obs CommitPluginObservation) Encode() ([]byte, error) {
	return json.Marshal(obs)
}

func DecodeCommitPluginObservation(b []byte) (CommitPluginObservation, error) {
	obs := CommitPluginObservation{}
	err := json.Unmarshal(b, &obs)
	return obs, err
}
