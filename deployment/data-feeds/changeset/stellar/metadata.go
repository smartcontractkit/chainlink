package stellar

import (
	"encoding/hex"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

// ContractMetadata mirrors a cache contract's feed configuration so it is
// queryable off-chain. On-chain state stays authoritative; only the feed
// config changesets mutate this.
type ContractMetadata struct {
	Cache string                  `json:"cache,omitempty"` // proxy only: current cache target
	Feeds map[string]FeedMetadata `json:"feeds,omitempty"` // cache only, keyed by 0x data id
}

// FeedMetadata mirrors one configured feed.
type FeedMetadata struct {
	Description string           `json:"description"`
	Decimals    uint32           `json:"decimals"`
	Permissions []FeedPermission `json:"permissions"`
}

// dataIDHex canonicalizes a feed id to full-width 0x hex.
func dataIDHex(id [16]byte) string {
	return "0x" + hex.EncodeToString(id[:])
}

// decimalsFromID mirrors the contract: byte 7 in [0x20,0x60] encodes decimals.
func decimalsFromID(id [16]byte) uint32 {
	if b := id[7]; b >= 0x20 && b <= 0x60 {
		return uint32(b - 0x20)
	}
	return 0
}

// metadataOutput applies mutate to the contract's latest stored metadata and
// returns a ChangesetOutput carrying the updated record.
func metadataOutput(env cldf.Environment, chainSel uint64, address string, mutate func(*ContractMetadata)) (cldf.ChangesetOutput, error) {
	meta := loadMetadata(env, chainSel, address)
	mutate(&meta)
	var out cldf.ChangesetOutput
	out.DataStore = datastore.NewMemoryDataStore()
	return out, out.DataStore.ContractMetadata().Upsert(datastore.ContractMetadata{
		ChainSelector: chainSel,
		Address:       address,
		Metadata:      meta,
	})
}

// loadMetadata returns the contract's stored metadata, or a zero value.
func loadMetadata(env cldf.Environment, chainSel uint64, address string) ContractMetadata {
	rec, err := env.DataStore.ContractMetadata().Get(datastore.NewContractMetadataKey(chainSel, address))
	if err != nil {
		return ContractMetadata{}
	}
	meta, err := datastore.As[ContractMetadata](rec.Metadata)
	if err != nil {
		return ContractMetadata{}
	}
	return meta
}
