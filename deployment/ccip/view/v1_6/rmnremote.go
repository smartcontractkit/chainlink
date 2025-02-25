package v1_6

import (
	"encoding/hex"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	"github.com/smartcontractkit/chainlink/deployment/common/view/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/rmn_remote"
)

type RMNRemoteCurseEntry struct {
	Subject  string `json:"subject"`
	Selector uint64 `json:"selector"`
}

type RMNRemoteView struct {
	types.ContractMetaData
	IsCursed             bool                     `json:"isCursed"`
	Config               RMNRemoteVersionedConfig `json:"config,omitempty"`
	CursedSubjectEntries []RMNRemoteCurseEntry    `json:"cursedSubjectEntries,omitempty"`
}

type RMNRemoteVersionedConfig struct {
	Version uint32            `json:"version"`
	Signers []RMNRemoteSigner `json:"signers"`
	Fsign   uint64            `json:"fSign"`
}

type RMNRemoteSigner struct {
	OnchainPublicKey string `json:"onchain_public_key"`
	NodeIndex        uint64 `json:"node_index"`
}

func mapCurseSubjects(subjects [][16]byte) []RMNRemoteCurseEntry {
	res := make([]RMNRemoteCurseEntry, 0, len(subjects))
	for _, subject := range subjects {
		res = append(res, RMNRemoteCurseEntry{
			Subject:  hex.EncodeToString(subject[:]),
			Selector: globals.SubjectToSelector(subject),
		})
	}
	return res
}

func GenerateRMNRemoteView(rmnReader *rmn_remote.RMNRemote) (RMNRemoteView, error) {
	tv, err := types.NewContractMetaData(rmnReader, rmnReader.Address())
	if err != nil {
		return RMNRemoteView{}, err
	}
	config, err := rmnReader.GetVersionedConfig(nil)
	if err != nil {
		return RMNRemoteView{}, err
	}
	rmnConfig := RMNRemoteVersionedConfig{
		Version: config.Version,
		Signers: make([]RMNRemoteSigner, 0, len(config.Config.Signers)),
		Fsign:   config.Config.FSign,
	}
	for _, signer := range config.Config.Signers {
		rmnConfig.Signers = append(rmnConfig.Signers, RMNRemoteSigner{
			OnchainPublicKey: signer.OnchainPublicKey.Hex(),
			NodeIndex:        signer.NodeIndex,
		})
	}
	isCursed, err := rmnReader.IsCursed0(nil)
	if err != nil {
		return RMNRemoteView{}, err
	}

	curseSubjects, err := rmnReader.GetCursedSubjects(nil)
	if err != nil {
		return RMNRemoteView{}, err
	}

	return RMNRemoteView{
		ContractMetaData:     tv,
		IsCursed:             isCursed,
		Config:               rmnConfig,
		CursedSubjectEntries: mapCurseSubjects(curseSubjects),
	}, nil
}
