package por

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type BlockNumber uint64

type Blocks map[ChainSelector]BlockNumber

type BlockMintablePair struct {
	Block    BlockNumber
	Mintable *big.Int
}

type Mintables map[ChainSelector]BlockMintablePair

func (m Mintables) toString() (string, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func mintablesFromString(s string) (Mintables, error) {
	var m Mintables
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		return nil, err
	}
	return m, nil
}

type ReserveInfo struct {
	ReserveAmount *big.Int
	Timestamp     time.Time
}

type ExternalAdapterPayload struct {
	Mintables   Mintables   // The mintable amounts for each chain and its block.
	ReserveInfo ReserveInfo // The latest reserve amount and timestamp used to calculate the minting allowance above.

	LatestBlocks Blocks // The latest blocks for each chain.
}

type TransmittedReportDetails struct {
	ConfigDigest    types.ConfigDigest // The OCR3 config digest.
	SeqNr           uint64             // The OCR3 sequence number.
	LatestTimestamp time.Time          // The (on-chain specific) timestamp of the block where the latest report is included.
}

type PorOffchainConfig struct {
	MaxChains uint32 // The maximum number of chains that can be tracked by the external adapter.
}

func (p *PorOffchainConfig) Serialize() ([]byte, error) {
	return json.Marshal(p)
}

func DeserializePorOffchainConfig(data []byte) (*PorOffchainConfig, error) {
	var config PorOffchainConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

type PorReport struct {
	ConfigDigest types.ConfigDigest
	SeqNr        uint64
	Block        BlockNumber
	Mintable     *big.Int

	// The following fields might be useful in the future, but are not currently used
	// ReserveAmount *big.Int
	// ReserveTimestamp time.Time
}
