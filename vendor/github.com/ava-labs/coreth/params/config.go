// (c) 2019-2020, Ava Labs, Inc.
//
// This file is a derived work, based on the go-ethereum library whose original
// notices appear below.
//
// It is distributed under a license compatible with the licensing terms of the
// original code from which it is derived.
//
// Much love to the original authors for their work.
// **********
// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package params

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ava-labs/coreth/precompile"
	"github.com/ava-labs/coreth/utils"
	"github.com/ethereum/go-ethereum/common"
)

// Avalanche ChainIDs
var (
	// AvalancheMainnetChainID ...
	AvalancheMainnetChainID = big.NewInt(43114)
	// AvalancheFujiChainID ...
	AvalancheFujiChainID = big.NewInt(43113)
	// AvalancheLocalChainID ...
	AvalancheLocalChainID = big.NewInt(43112)

	errNonGenesisForkByHeight = errors.New("coreth only supports forking by height at the genesis block")
)

var (
	// AvalancheMainnetChainConfig is the configuration for Avalanche Main Network
	AvalancheMainnetChainConfig = &ChainConfig{
		ChainID:                         AvalancheMainnetChainID,
		HomesteadBlock:                  big.NewInt(0),
		DAOForkBlock:                    big.NewInt(0),
		DAOForkSupport:                  true,
		EIP150Block:                     big.NewInt(0),
		EIP150Hash:                      common.HexToHash("0x2086799aeebeae135c246c65021c82b4e15a2c451340993aacfd2751886514f0"),
		EIP155Block:                     big.NewInt(0),
		EIP158Block:                     big.NewInt(0),
		ByzantiumBlock:                  big.NewInt(0),
		ConstantinopleBlock:             big.NewInt(0),
		PetersburgBlock:                 big.NewInt(0),
		IstanbulBlock:                   big.NewInt(0),
		MuirGlacierBlock:                big.NewInt(0),
		ApricotPhase1BlockTimestamp:     big.NewInt(time.Date(2021, time.March, 31, 14, 0, 0, 0, time.UTC).Unix()),
		ApricotPhase2BlockTimestamp:     big.NewInt(time.Date(2021, time.May, 10, 11, 0, 0, 0, time.UTC).Unix()),
		ApricotPhase3BlockTimestamp:     big.NewInt(time.Date(2021, time.August, 24, 14, 0, 0, 0, time.UTC).Unix()),
		ApricotPhase4BlockTimestamp:     big.NewInt(time.Date(2021, time.September, 22, 21, 0, 0, 0, time.UTC).Unix()),
		ApricotPhase5BlockTimestamp:     big.NewInt(time.Date(2021, time.December, 2, 18, 0, 0, 0, time.UTC).Unix()),
		ApricotPhasePre6BlockTimestamp:  big.NewInt(time.Date(2022, time.September, 5, 1, 30, 0, 0, time.UTC).Unix()),
		ApricotPhase6BlockTimestamp:     big.NewInt(time.Date(2022, time.September, 6, 20, 0, 0, 0, time.UTC).Unix()),
		ApricotPhasePost6BlockTimestamp: big.NewInt(time.Date(2022, time.September, 7, 3, 0, 0, 0, time.UTC).Unix()),
		BanffBlockTimestamp:             big.NewInt(time.Date(2022, time.October, 18, 16, 0, 0, 0, time.UTC).Unix()),
		// TODO Add Banff and Cortina timestamps
	}

	// AvalancheFujiChainConfig is the configuration for the Fuji Test Network
	AvalancheFujiChainConfig = &ChainConfig{
		ChainID:                         AvalancheFujiChainID,
		HomesteadBlock:                  big.NewInt(0),
		DAOForkBlock:                    big.NewInt(0),
		DAOForkSupport:                  true,
		EIP150Block:                     big.NewInt(0),
		EIP150Hash:                      common.HexToHash("0x2086799aeebeae135c246c65021c82b4e15a2c451340993aacfd2751886514f0"),
		EIP155Block:                     big.NewInt(0),
		EIP158Block:                     big.NewInt(0),
		ByzantiumBlock:                  big.NewInt(0),
		ConstantinopleBlock:             big.NewInt(0),
		PetersburgBlock:                 big.NewInt(0),
		IstanbulBlock:                   big.NewInt(0),
		MuirGlacierBlock:                big.NewInt(0),
		ApricotPhase1BlockTimestamp:     big.NewInt(time.Date(2021, time.March, 26, 14, 0, 0, 0, time.UTC).Unix()),
		ApricotPhase2BlockTimestamp:     big.NewInt(time.Date(2021, time.May, 5, 14, 0, 0, 0, time.UTC).Unix()),
		ApricotPhase3BlockTimestamp:     big.NewInt(time.Date(2021, time.August, 16, 19, 0, 0, 0, time.UTC).Unix()),
		ApricotPhase4BlockTimestamp:     big.NewInt(time.Date(2021, time.September, 16, 21, 0, 0, 0, time.UTC).Unix()),
		ApricotPhase5BlockTimestamp:     big.NewInt(time.Date(2021, time.November, 24, 15, 0, 0, 0, time.UTC).Unix()),
		ApricotPhasePre6BlockTimestamp:  big.NewInt(time.Date(2022, time.September, 6, 20, 0, 0, 0, time.UTC).Unix()),
		ApricotPhase6BlockTimestamp:     big.NewInt(time.Date(2022, time.September, 6, 20, 0, 0, 0, time.UTC).Unix()),
		ApricotPhasePost6BlockTimestamp: big.NewInt(time.Date(2022, time.September, 7, 6, 0, 0, 0, time.UTC).Unix()),
		BanffBlockTimestamp:             big.NewInt(time.Date(2022, time.October, 3, 14, 0, 0, 0, time.UTC).Unix()),
		// TODO add Cortina timestamp
	}

	// AvalancheLocalChainConfig is the configuration for the Avalanche Local Network
	AvalancheLocalChainConfig = &ChainConfig{
		ChainID:                         AvalancheLocalChainID,
		HomesteadBlock:                  big.NewInt(0),
		DAOForkBlock:                    big.NewInt(0),
		DAOForkSupport:                  true,
		EIP150Block:                     big.NewInt(0),
		EIP150Hash:                      common.HexToHash("0x2086799aeebeae135c246c65021c82b4e15a2c451340993aacfd2751886514f0"),
		EIP155Block:                     big.NewInt(0),
		EIP158Block:                     big.NewInt(0),
		ByzantiumBlock:                  big.NewInt(0),
		ConstantinopleBlock:             big.NewInt(0),
		PetersburgBlock:                 big.NewInt(0),
		IstanbulBlock:                   big.NewInt(0),
		MuirGlacierBlock:                big.NewInt(0),
		ApricotPhase1BlockTimestamp:     big.NewInt(0),
		ApricotPhase2BlockTimestamp:     big.NewInt(0),
		ApricotPhase3BlockTimestamp:     big.NewInt(0),
		ApricotPhase4BlockTimestamp:     big.NewInt(0),
		ApricotPhase5BlockTimestamp:     big.NewInt(0),
		ApricotPhasePre6BlockTimestamp:  big.NewInt(0),
		ApricotPhase6BlockTimestamp:     big.NewInt(0),
		ApricotPhasePost6BlockTimestamp: big.NewInt(0),
		BanffBlockTimestamp:             big.NewInt(0),
	}

	TestChainConfig             = &ChainConfig{big.NewInt(1), big.NewInt(0), nil, false, big.NewInt(0), common.Hash{}, big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)}
	TestLaunchConfig            = &ChainConfig{big.NewInt(1), big.NewInt(0), nil, false, big.NewInt(0), common.Hash{}, big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), nil, nil, nil, nil, nil, nil, nil, nil, nil, nil}
	TestApricotPhase1Config     = &ChainConfig{big.NewInt(1), big.NewInt(0), nil, false, big.NewInt(0), common.Hash{}, big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), nil, nil, nil, nil, nil, nil, nil, nil, nil}
	TestApricotPhase2Config     = &ChainConfig{big.NewInt(1), big.NewInt(0), nil, false, big.NewInt(0), common.Hash{}, big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), nil, nil, nil, nil, nil, nil, nil, nil}
	TestApricotPhase3Config     = &ChainConfig{big.NewInt(1), big.NewInt(0), nil, false, big.NewInt(0), common.Hash{}, big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), nil, nil, nil, nil, nil, nil, nil}
	TestApricotPhase4Config     = &ChainConfig{big.NewInt(1), big.NewInt(0), nil, false, big.NewInt(0), common.Hash{}, big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), nil, nil, nil, nil, nil, nil}
	TestApricotPhase5Config     = &ChainConfig{big.NewInt(1), big.NewInt(0), nil, false, big.NewInt(0), common.Hash{}, big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), nil, nil, nil, nil, nil}
	TestApricotPhasePre6Config  = &ChainConfig{big.NewInt(1), big.NewInt(0), nil, false, big.NewInt(0), common.Hash{}, big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), nil, nil, nil, nil}
	TestApricotPhase6Config     = &ChainConfig{big.NewInt(1), big.NewInt(0), nil, false, big.NewInt(0), common.Hash{}, big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), nil, nil, nil}
	TestApricotPhasePost6Config = &ChainConfig{big.NewInt(1), big.NewInt(0), nil, false, big.NewInt(0), common.Hash{}, big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), nil, nil}
	TestBanffChainConfig        = &ChainConfig{big.NewInt(1), big.NewInt(0), nil, false, big.NewInt(0), common.Hash{}, big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), nil}
	TestCortinaChainConfig      = &ChainConfig{big.NewInt(1), big.NewInt(0), nil, false, big.NewInt(0), common.Hash{}, big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0)}
	TestRules                   = TestChainConfig.AvalancheRules(new(big.Int), new(big.Int))
)

// ChainConfig is the core config which determines the blockchain settings.
//
// ChainConfig is stored in the database on a per block basis. This means
// that any network, identified by its genesis block, can have its own
// set of configuration options.
type ChainConfig struct {
	ChainID *big.Int `json:"chainId"` // chainId identifies the current chain and is used for replay protection

	HomesteadBlock *big.Int `json:"homesteadBlock,omitempty"` // Homestead switch block (nil = no fork, 0 = already homestead)

	DAOForkBlock   *big.Int `json:"daoForkBlock,omitempty"`   // TheDAO hard-fork switch block (nil = no fork)
	DAOForkSupport bool     `json:"daoForkSupport,omitempty"` // Whether the nodes supports or opposes the DAO hard-fork

	// EIP150 implements the Gas price changes (https://github.com/ethereum/EIPs/issues/150)
	EIP150Block *big.Int    `json:"eip150Block,omitempty"` // EIP150 HF block (nil = no fork)
	EIP150Hash  common.Hash `json:"eip150Hash,omitempty"`  // EIP150 HF hash (needed for header only clients as only gas pricing changed)

	EIP155Block *big.Int `json:"eip155Block,omitempty"` // EIP155 HF block
	EIP158Block *big.Int `json:"eip158Block,omitempty"` // EIP158 HF block

	ByzantiumBlock      *big.Int `json:"byzantiumBlock,omitempty"`      // Byzantium switch block (nil = no fork, 0 = already on byzantium)
	ConstantinopleBlock *big.Int `json:"constantinopleBlock,omitempty"` // Constantinople switch block (nil = no fork, 0 = already activated)
	PetersburgBlock     *big.Int `json:"petersburgBlock,omitempty"`     // Petersburg switch block (nil = same as Constantinople)
	IstanbulBlock       *big.Int `json:"istanbulBlock,omitempty"`       // Istanbul switch block (nil = no fork, 0 = already on istanbul)
	MuirGlacierBlock    *big.Int `json:"muirGlacierBlock,omitempty"`    // Eip-2384 (bomb delay) switch block (nil = no fork, 0 = already activated)

	// Avalanche Network Upgrades
	ApricotPhase1BlockTimestamp *big.Int `json:"apricotPhase1BlockTimestamp,omitempty"` // Apricot Phase 1 Block Timestamp (nil = no fork, 0 = already activated)
	// Apricot Phase 2 Block Timestamp (nil = no fork, 0 = already activated)
	// Apricot Phase 2 includes a modified version of the Berlin Hard Fork from Ethereum
	ApricotPhase2BlockTimestamp *big.Int `json:"apricotPhase2BlockTimestamp,omitempty"`
	// Apricot Phase 3 introduces dynamic fees and a modified version of the London Hard Fork from Ethereum (nil = no fork, 0 = already activated)
	ApricotPhase3BlockTimestamp *big.Int `json:"apricotPhase3BlockTimestamp,omitempty"`
	// Apricot Phase 4 introduces the notion of a block fee to the dynamic fee algorithm (nil = no fork, 0 = already activated)
	ApricotPhase4BlockTimestamp *big.Int `json:"apricotPhase4BlockTimestamp,omitempty"`
	// Apricot Phase 5 introduces a batch of atomic transactions with a maximum atomic gas limit per block. (nil = no fork, 0 = already activated)
	ApricotPhase5BlockTimestamp *big.Int `json:"apricotPhase5BlockTimestamp,omitempty"`
	// Apricot Phase Pre-6 deprecates the NativeAssetCall precompile (soft). (nil = no fork, 0 = already activated)
	ApricotPhasePre6BlockTimestamp *big.Int `json:"apricotPhasePre6BlockTimestamp,omitempty"`
	// Apricot Phase 6 deprecates the NativeAssetBalance and NativeAssetCall precompiles. (nil = no fork, 0 = already activated)
	ApricotPhase6BlockTimestamp *big.Int `json:"apricotPhase6BlockTimestamp,omitempty"`
	// Apricot Phase Post-6 deprecates the NativeAssetCall precompile (soft). (nil = no fork, 0 = already activated)
	ApricotPhasePost6BlockTimestamp *big.Int `json:"apricotPhasePost6BlockTimestamp,omitempty"`
	// Banff TODO comment. (nil = no fork, 0 = already activated)
	BanffBlockTimestamp *big.Int `json:"banffBlockTimestamp,omitempty"`
	// Cortina TODO comment. (nil = no fork, 0 = already activated)
	CortinaBlockTimestamp *big.Int `json:"cortinaBlockTimestamp,omitempty"`
}

// String implements the fmt.Stringer interface.
func (c *ChainConfig) String() string {
	var banner string

	banner += fmt.Sprintf("Chain ID:  %v\n", c.ChainID)
	banner += "Consensus: Dummy Consensus Engine\n\n"

	// Create a list of forks with a short description of them. Forks that only
	// makes sense for mainnet should be optional at printing to avoid bloating
	// the output for testnets and private networks.
	banner += "Hard Forks:\n"
	banner += fmt.Sprintf(" - Homestead:                   %-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/homestead.md)\n", c.HomesteadBlock)
	if c.DAOForkBlock != nil {
		banner += fmt.Sprintf(" - DAO Fork:                    %-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/dao-fork.md)\n", c.DAOForkBlock)
	}
	banner += fmt.Sprintf(" - Tangerine Whistle (EIP 150): %-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/tangerine-whistle.md)\n", c.EIP150Block)
	banner += fmt.Sprintf(" - Spurious Dragon/1 (EIP 155): %-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/spurious-dragon.md)\n", c.EIP155Block)
	banner += fmt.Sprintf(" - Spurious Dragon/2 (EIP 158): %-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/spurious-dragon.md)\n", c.EIP155Block)
	banner += fmt.Sprintf(" - Byzantium:                   %-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/byzantium.md)\n", c.ByzantiumBlock)
	banner += fmt.Sprintf(" - Constantinople:              %-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/constantinople.md)\n", c.ConstantinopleBlock)
	banner += fmt.Sprintf(" - Petersburg:                  %-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/petersburg.md)\n", c.PetersburgBlock)
	banner += fmt.Sprintf(" - Istanbul:                    %-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/istanbul.md)\n", c.IstanbulBlock)
	if c.MuirGlacierBlock != nil {
		banner += fmt.Sprintf(" - Muir Glacier:                %-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/muir-glacier.md)\n", c.MuirGlacierBlock)
	}
	banner += fmt.Sprintf(" - Apricot Phase 1 Timestamp:        %-8v (https://github.com/ava-labs/avalanchego/releases/tag/v1.3.0)\n", c.ApricotPhase1BlockTimestamp)
	banner += fmt.Sprintf(" - Apricot Phase 2 Timestamp:        %-8v (https://github.com/ava-labs/avalanchego/releases/tag/v1.4.0)\n", c.ApricotPhase2BlockTimestamp)
	banner += fmt.Sprintf(" - Apricot Phase 3 Timestamp:        %-8v (https://github.com/ava-labs/avalanchego/releases/tag/v1.5.0)\n", c.ApricotPhase3BlockTimestamp)
	banner += fmt.Sprintf(" - Apricot Phase 4 Timestamp:        %-8v (https://github.com/ava-labs/avalanchego/releases/tag/v1.6.0)\n", c.ApricotPhase4BlockTimestamp)
	banner += fmt.Sprintf(" - Apricot Phase 5 Timestamp:        %-8v (https://github.com/ava-labs/avalanchego/releases/tag/v1.7.0)\n", c.ApricotPhase5BlockTimestamp)
	banner += fmt.Sprintf(" - Apricot Phase P6 Timestamp        %-8v (https://github.com/ava-labs/avalanchego/releases/tag/v1.8.0)\n", c.ApricotPhasePre6BlockTimestamp)
	banner += fmt.Sprintf(" - Apricot Phase 6 Timestamp:        %-8v (https://github.com/ava-labs/avalanchego/releases/tag/v1.8.0)\n", c.ApricotPhase6BlockTimestamp)
	banner += fmt.Sprintf(" - Apricot Phase Post-6 Timestamp:   %-8v (https://github.com/ava-labs/avalanchego/releases/tag/v1.8.0\n", c.ApricotPhasePost6BlockTimestamp)
	banner += fmt.Sprintf(" - Banff Timestamp:                  %-8v (https://github.com/ava-labs/avalanchego/releases/tag/v1.9.0)\n", c.BanffBlockTimestamp)
	banner += fmt.Sprintf(" - Cortina Timestamp:                %-8v (https://github.com/ava-labs/avalanchego/releases/tag/v1.10.0)\n", c.CortinaBlockTimestamp)
	banner += "\n"
	return banner
}

// IsHomestead returns whether num is either equal to the homestead block or greater.
func (c *ChainConfig) IsHomestead(num *big.Int) bool {
	return utils.IsForked(c.HomesteadBlock, num)
}

// IsDAOFork returns whether num is either equal to the DAO fork block or greater.
func (c *ChainConfig) IsDAOFork(num *big.Int) bool {
	return utils.IsForked(c.DAOForkBlock, num)
}

// IsEIP150 returns whether num is either equal to the EIP150 fork block or greater.
func (c *ChainConfig) IsEIP150(num *big.Int) bool {
	return utils.IsForked(c.EIP150Block, num)
}

// IsEIP155 returns whether num is either equal to the EIP155 fork block or greater.
func (c *ChainConfig) IsEIP155(num *big.Int) bool {
	return utils.IsForked(c.EIP155Block, num)
}

// IsEIP158 returns whether num is either equal to the EIP158 fork block or greater.
func (c *ChainConfig) IsEIP158(num *big.Int) bool {
	return utils.IsForked(c.EIP158Block, num)
}

// IsByzantium returns whether num is either equal to the Byzantium fork block or greater.
func (c *ChainConfig) IsByzantium(num *big.Int) bool {
	return utils.IsForked(c.ByzantiumBlock, num)
}

// IsConstantinople returns whether num is either equal to the Constantinople fork block or greater.
func (c *ChainConfig) IsConstantinople(num *big.Int) bool {
	return utils.IsForked(c.ConstantinopleBlock, num)
}

// IsMuirGlacier returns whether num is either equal to the Muir Glacier (EIP-2384) fork block or greater.
func (c *ChainConfig) IsMuirGlacier(num *big.Int) bool {
	return utils.IsForked(c.MuirGlacierBlock, num)
}

// IsPetersburg returns whether num is either
// - equal to or greater than the PetersburgBlock fork block,
// - OR is nil, and Constantinople is active
func (c *ChainConfig) IsPetersburg(num *big.Int) bool {
	return utils.IsForked(c.PetersburgBlock, num) || c.PetersburgBlock == nil && utils.IsForked(c.ConstantinopleBlock, num)
}

// IsIstanbul returns whether num is either equal to the Istanbul fork block or greater.
func (c *ChainConfig) IsIstanbul(num *big.Int) bool {
	return utils.IsForked(c.IstanbulBlock, num)
}

// Avalanche Upgrades:

// IsApricotPhase1 returns whether [blockTimestamp] represents a block
// with a timestamp after the Apricot Phase 1 upgrade time.
func (c *ChainConfig) IsApricotPhase1(blockTimestamp *big.Int) bool {
	return utils.IsForked(c.ApricotPhase1BlockTimestamp, blockTimestamp)
}

// IsApricotPhase2 returns whether [blockTimestamp] represents a block
// with a timestamp after the Apricot Phase 2 upgrade time.
func (c *ChainConfig) IsApricotPhase2(blockTimestamp *big.Int) bool {
	return utils.IsForked(c.ApricotPhase2BlockTimestamp, blockTimestamp)
}

// IsApricotPhase3 returns whether [blockTimestamp] represents a block
// with a timestamp after the Apricot Phase 3 upgrade time.
func (c *ChainConfig) IsApricotPhase3(blockTimestamp *big.Int) bool {
	return utils.IsForked(c.ApricotPhase3BlockTimestamp, blockTimestamp)
}

// IsApricotPhase4 returns whether [blockTimestamp] represents a block
// with a timestamp after the Apricot Phase 4 upgrade time.
func (c *ChainConfig) IsApricotPhase4(blockTimestamp *big.Int) bool {
	return utils.IsForked(c.ApricotPhase4BlockTimestamp, blockTimestamp)
}

// IsApricotPhase5 returns whether [blockTimestamp] represents a block
// with a timestamp after the Apricot Phase 5 upgrade time.
func (c *ChainConfig) IsApricotPhase5(blockTimestamp *big.Int) bool {
	return utils.IsForked(c.ApricotPhase5BlockTimestamp, blockTimestamp)
}

// IsApricotPhasePre6 returns whether [blockTimestamp] represents a block
// with a timestamp after the Apricot Phase Pre 6 upgrade time.
func (c *ChainConfig) IsApricotPhasePre6(blockTimestamp *big.Int) bool {
	return utils.IsForked(c.ApricotPhasePre6BlockTimestamp, blockTimestamp)
}

// IsApricotPhase6 returns whether [blockTimestamp] represents a block
// with a timestamp after the Apricot Phase 6 upgrade time.
func (c *ChainConfig) IsApricotPhase6(blockTimestamp *big.Int) bool {
	return utils.IsForked(c.ApricotPhase6BlockTimestamp, blockTimestamp)
}

// IsApricotPhasePost6 returns whether [blockTimestamp] represents a block
// with a timestamp after the Apricot Phase 6 Post upgrade time.
func (c *ChainConfig) IsApricotPhasePost6(blockTimestamp *big.Int) bool {
	return utils.IsForked(c.ApricotPhasePost6BlockTimestamp, blockTimestamp)
}

// IsBanff returns whether [blockTimestamp] represents a block
// with a timestamp after the Banff upgrade time.
func (c *ChainConfig) IsBanff(blockTimestamp *big.Int) bool {
	return utils.IsForked(c.BanffBlockTimestamp, blockTimestamp)
}

// IsCortina returns whether [blockTimestamp] represents a block
// with a timestamp after the Cortina upgrade time.
func (c *ChainConfig) IsCortina(blockTimestamp *big.Int) bool {
	return utils.IsForked(c.CortinaBlockTimestamp, blockTimestamp)
}

// CheckCompatible checks whether scheduled fork transitions have been imported
// with a mismatching chain configuration.
func (c *ChainConfig) CheckCompatible(newcfg *ChainConfig, height uint64, timestamp uint64) *ConfigCompatError {
	bNumber := new(big.Int).SetUint64(height)
	bTimestamp := new(big.Int).SetUint64(timestamp)

	// Iterate checkCompatible to find the lowest conflict.
	var lasterr *ConfigCompatError
	for {
		err := c.checkCompatible(newcfg, bNumber, bTimestamp)
		if err == nil || (lasterr != nil && err.RewindTo == lasterr.RewindTo) {
			break
		}
		lasterr = err
		bNumber.SetUint64(err.RewindTo)
	}
	return lasterr
}

// CheckConfigForkOrder checks that we don't "skip" any forks, geth isn't pluggable enough
// to guarantee that forks can be implemented in a different order than on official networks
func (c *ChainConfig) CheckConfigForkOrder() error {
	type fork struct {
		name     string
		block    *big.Int
		optional bool // if true, the fork may be nil and next fork is still allowed
	}
	var lastFork fork
	for _, cur := range []fork{
		{name: "homesteadBlock", block: c.HomesteadBlock},
		{name: "daoForkBlock", block: c.DAOForkBlock, optional: true},
		{name: "eip150Block", block: c.EIP150Block},
		{name: "eip155Block", block: c.EIP155Block},
		{name: "eip158Block", block: c.EIP158Block},
		{name: "byzantiumBlock", block: c.ByzantiumBlock},
		{name: "constantinopleBlock", block: c.ConstantinopleBlock},
		{name: "petersburgBlock", block: c.PetersburgBlock},
		{name: "istanbulBlock", block: c.IstanbulBlock},
		{name: "muirGlacierBlock", block: c.MuirGlacierBlock, optional: true},
	} {
		if cur.block != nil && common.Big0.Cmp(cur.block) != 0 {
			return errNonGenesisForkByHeight
		}
		if lastFork.name != "" {
			// Next one must be higher number
			if lastFork.block == nil && cur.block != nil {
				return fmt.Errorf("unsupported fork ordering: %v not enabled, but %v enabled at %v",
					lastFork.name, cur.name, cur.block)
			}
			if lastFork.block != nil && cur.block != nil {
				if lastFork.block.Cmp(cur.block) > 0 {
					return fmt.Errorf("unsupported fork ordering: %v enabled at %v, but %v enabled at %v",
						lastFork.name, lastFork.block, cur.name, cur.block)
				}
			}
		}
		// If it was optional and not set, then ignore it
		if !cur.optional || cur.block != nil {
			lastFork = cur
		}
	}

	// Note: ApricotPhase1 and ApricotPhase2 override the rules set by block number
	// hard forks. In Avalanche, hard forks must take place via block timestamps instead
	// of block numbers since blocks are produced asynchronously. Therefore, we do not
	// check that the block timestamps for Apricot Phase1 and Phase2 in the same way as for
	// the block number forks since it would not be a meaningful comparison.
	// Instead, we check only that Apricot Phases are enabled in order.
	lastFork = fork{}
	for _, cur := range []fork{
		{name: "apricotPhase1BlockTimestamp", block: c.ApricotPhase1BlockTimestamp},
		{name: "apricotPhase2BlockTimestamp", block: c.ApricotPhase2BlockTimestamp},
		{name: "apricotPhase3BlockTimestamp", block: c.ApricotPhase3BlockTimestamp},
		{name: "apricotPhase4BlockTimestamp", block: c.ApricotPhase4BlockTimestamp},
		{name: "apricotPhase5BlockTimestamp", block: c.ApricotPhase5BlockTimestamp},
		{name: "apricotPhasePre6BlockTimestamp", block: c.ApricotPhasePre6BlockTimestamp},
		{name: "apricotPhase6BlockTimestamp", block: c.ApricotPhase6BlockTimestamp},
		{name: "apricotPhasePost6BlockTimestamp", block: c.ApricotPhasePost6BlockTimestamp},
		{name: "banffBlockTimestamp", block: c.BanffBlockTimestamp},
		{name: "cortinaBlockTimestamp", block: c.CortinaBlockTimestamp},
	} {
		if lastFork.name != "" {
			// Next one must be higher number
			if lastFork.block == nil && cur.block != nil {
				return fmt.Errorf("unsupported fork ordering: %v not enabled, but %v enabled at %v",
					lastFork.name, cur.name, cur.block)
			}
			if lastFork.block != nil && cur.block != nil {
				if lastFork.block.Cmp(cur.block) > 0 {
					return fmt.Errorf("unsupported fork ordering: %v enabled at %v, but %v enabled at %v",
						lastFork.name, lastFork.block, cur.name, cur.block)
				}
			}
		}
		// If it was optional and not set, then ignore it
		if !cur.optional || cur.block != nil {
			lastFork = cur
		}
	}
	// TODO(aaronbuchwald) check that avalanche block timestamps are at least possible with the other rule set changes
	// additional change: require that block number hard forks are either 0 or nil since they should not
	// be enabled at a specific block number.

	return nil
}

func (c *ChainConfig) checkCompatible(newcfg *ChainConfig, lastHeight *big.Int, lastTimestamp *big.Int) *ConfigCompatError {
	if isForkIncompatible(c.HomesteadBlock, newcfg.HomesteadBlock, lastHeight) {
		return newCompatError("Homestead fork block", c.HomesteadBlock, newcfg.HomesteadBlock)
	}
	if isForkIncompatible(c.DAOForkBlock, newcfg.DAOForkBlock, lastHeight) {
		return newCompatError("DAO fork block", c.DAOForkBlock, newcfg.DAOForkBlock)
	}
	if c.IsDAOFork(lastHeight) && c.DAOForkSupport != newcfg.DAOForkSupport {
		return newCompatError("DAO fork support flag", c.DAOForkBlock, newcfg.DAOForkBlock)
	}
	if isForkIncompatible(c.EIP150Block, newcfg.EIP150Block, lastHeight) {
		return newCompatError("EIP150 fork block", c.EIP150Block, newcfg.EIP150Block)
	}
	if isForkIncompatible(c.EIP155Block, newcfg.EIP155Block, lastHeight) {
		return newCompatError("EIP155 fork block", c.EIP155Block, newcfg.EIP155Block)
	}
	if isForkIncompatible(c.EIP158Block, newcfg.EIP158Block, lastHeight) {
		return newCompatError("EIP158 fork block", c.EIP158Block, newcfg.EIP158Block)
	}
	if c.IsEIP158(lastHeight) && !configNumEqual(c.ChainID, newcfg.ChainID) {
		return newCompatError("EIP158 chain ID", c.EIP158Block, newcfg.EIP158Block)
	}
	if isForkIncompatible(c.ByzantiumBlock, newcfg.ByzantiumBlock, lastHeight) {
		return newCompatError("Byzantium fork block", c.ByzantiumBlock, newcfg.ByzantiumBlock)
	}
	if isForkIncompatible(c.ConstantinopleBlock, newcfg.ConstantinopleBlock, lastHeight) {
		return newCompatError("Constantinople fork block", c.ConstantinopleBlock, newcfg.ConstantinopleBlock)
	}
	if isForkIncompatible(c.PetersburgBlock, newcfg.PetersburgBlock, lastHeight) {
		// the only case where we allow Petersburg to be set in the past is if it is equal to Constantinople
		// mainly to satisfy fork ordering requirements which state that Petersburg fork be set if Constantinople fork is set
		if isForkIncompatible(c.ConstantinopleBlock, newcfg.PetersburgBlock, lastHeight) {
			return newCompatError("Petersburg fork block", c.PetersburgBlock, newcfg.PetersburgBlock)
		}
	}
	if isForkIncompatible(c.IstanbulBlock, newcfg.IstanbulBlock, lastHeight) {
		return newCompatError("Istanbul fork block", c.IstanbulBlock, newcfg.IstanbulBlock)
	}
	if isForkIncompatible(c.MuirGlacierBlock, newcfg.MuirGlacierBlock, lastHeight) {
		return newCompatError("Muir Glacier fork block", c.MuirGlacierBlock, newcfg.MuirGlacierBlock)
	}
	if isForkIncompatible(c.ApricotPhase1BlockTimestamp, newcfg.ApricotPhase1BlockTimestamp, lastTimestamp) {
		return newCompatError("ApricotPhase1 fork block timestamp", c.ApricotPhase1BlockTimestamp, newcfg.ApricotPhase1BlockTimestamp)
	}
	if isForkIncompatible(c.ApricotPhase2BlockTimestamp, newcfg.ApricotPhase2BlockTimestamp, lastTimestamp) {
		return newCompatError("ApricotPhase2 fork block timestamp", c.ApricotPhase2BlockTimestamp, newcfg.ApricotPhase2BlockTimestamp)
	}
	if isForkIncompatible(c.ApricotPhase3BlockTimestamp, newcfg.ApricotPhase3BlockTimestamp, lastTimestamp) {
		return newCompatError("ApricotPhase3 fork block timestamp", c.ApricotPhase3BlockTimestamp, newcfg.ApricotPhase3BlockTimestamp)
	}
	if isForkIncompatible(c.ApricotPhase4BlockTimestamp, newcfg.ApricotPhase4BlockTimestamp, lastTimestamp) {
		return newCompatError("ApricotPhase4 fork block timestamp", c.ApricotPhase4BlockTimestamp, newcfg.ApricotPhase4BlockTimestamp)
	}
	if isForkIncompatible(c.ApricotPhase5BlockTimestamp, newcfg.ApricotPhase5BlockTimestamp, lastTimestamp) {
		return newCompatError("ApricotPhase5 fork block timestamp", c.ApricotPhase5BlockTimestamp, newcfg.ApricotPhase5BlockTimestamp)
	}
	// TODO: add Phase 6 checks

	// TODO activate isForkIncompatible checks
	// if isForkIncompatible(c.BanffBlockTimestamp, newcfg.BanffBlockTimestamp, lastTimestamp) {
	// 	return newCompatError("Banff fork block timestamp", c.BanffBlockTimestamp, newcfg.BanffBlockTimestamp)
	// }
	// if isForkIncompatible(c.CortinaBlockTimestamp, newcfg.CortinaBlockTimestamp, lastTimestamp) {
	// 	return newCompatError("Cortina fork block timestamp", c.CortinaBlockTimestamp, newcfg.CortinaBlockTimestamp)
	// }
	return nil
}

// isForkIncompatible returns true if a fork scheduled at s1 cannot be rescheduled to
// block s2 because head is already past the fork.
func isForkIncompatible(s1, s2, head *big.Int) bool {
	return (utils.IsForked(s1, head) || utils.IsForked(s2, head)) && !configNumEqual(s1, s2)
}

func configNumEqual(x, y *big.Int) bool {
	if x == nil {
		return y == nil
	}
	if y == nil {
		return x == nil
	}
	return x.Cmp(y) == 0
}

// ConfigCompatError is raised if the locally-stored blockchain is initialised with a
// ChainConfig that would alter the past.
type ConfigCompatError struct {
	What string
	// block numbers of the stored and new configurations
	StoredConfig, NewConfig *big.Int
	// the block number to which the local chain must be rewound to correct the error
	RewindTo uint64
}

func newCompatError(what string, storedblock, newblock *big.Int) *ConfigCompatError {
	var rew *big.Int
	switch {
	case storedblock == nil:
		rew = newblock
	case newblock == nil || storedblock.Cmp(newblock) < 0:
		rew = storedblock
	default:
		rew = newblock
	}
	err := &ConfigCompatError{what, storedblock, newblock, 0}
	if rew != nil && rew.Sign() > 0 {
		err.RewindTo = rew.Uint64() - 1
	}
	return err
}

func (err *ConfigCompatError) Error() string {
	return fmt.Sprintf("mismatching %s in database (have %d, want %d, rewindto %d)", err.What, err.StoredConfig, err.NewConfig, err.RewindTo)
}

// Rules wraps ChainConfig and is merely syntactic sugar or can be used for functions
// that do not have or require information about the block.
//
// Rules is a one time interface meaning that it shouldn't be used in between transition
// phases.
type Rules struct {
	ChainID                                                 *big.Int
	IsHomestead, IsEIP150, IsEIP155, IsEIP158               bool
	IsByzantium, IsConstantinople, IsPetersburg, IsIstanbul bool

	// Rules for Avalanche releases
	IsApricotPhase1, IsApricotPhase2, IsApricotPhase3, IsApricotPhase4, IsApricotPhase5 bool
	IsApricotPhasePre6, IsApricotPhase6, IsApricotPhasePost6                            bool
	IsBanff                                                                             bool
	IsCortina                                                                           bool

	// Precompiles maps addresses to stateful precompiled contracts that are enabled
	// for this rule set.
	// Note: none of these addresses should conflict with the address space used by
	// any existing precompiles.
	Precompiles map[common.Address]precompile.StatefulPrecompiledContract
}

// Rules ensures c's ChainID is not nil.
func (c *ChainConfig) rules(num *big.Int) Rules {
	chainID := c.ChainID
	if chainID == nil {
		chainID = new(big.Int)
	}
	return Rules{
		ChainID:          new(big.Int).Set(chainID),
		IsHomestead:      c.IsHomestead(num),
		IsEIP150:         c.IsEIP150(num),
		IsEIP155:         c.IsEIP155(num),
		IsEIP158:         c.IsEIP158(num),
		IsByzantium:      c.IsByzantium(num),
		IsConstantinople: c.IsConstantinople(num),
		IsPetersburg:     c.IsPetersburg(num),
		IsIstanbul:       c.IsIstanbul(num),
	}
}

// AvalancheRules returns the Avalanche modified rules to support Avalanche
// network upgrades
func (c *ChainConfig) AvalancheRules(blockNum, blockTimestamp *big.Int) Rules {
	rules := c.rules(blockNum)

	rules.IsApricotPhase1 = c.IsApricotPhase1(blockTimestamp)
	rules.IsApricotPhase2 = c.IsApricotPhase2(blockTimestamp)
	rules.IsApricotPhase3 = c.IsApricotPhase3(blockTimestamp)
	rules.IsApricotPhase4 = c.IsApricotPhase4(blockTimestamp)
	rules.IsApricotPhase5 = c.IsApricotPhase5(blockTimestamp)
	rules.IsApricotPhasePre6 = c.IsApricotPhasePre6(blockTimestamp)
	rules.IsApricotPhase6 = c.IsApricotPhase6(blockTimestamp)
	rules.IsApricotPhasePost6 = c.IsApricotPhasePost6(blockTimestamp)
	rules.IsBanff = c.IsBanff(blockTimestamp)
	rules.IsCortina = c.IsCortina(blockTimestamp)

	// Initialize the stateful precompiles that should be enabled at [blockTimestamp].
	rules.Precompiles = make(map[common.Address]precompile.StatefulPrecompiledContract)
	for _, config := range c.enabledStatefulPrecompiles() {
		if utils.IsForked(config.Timestamp(), blockTimestamp) {
			rules.Precompiles[config.Address()] = config.Contract()
		}
	}

	return rules
}

// enabledStatefulPrecompiles returns a list of stateful precompile configs in the order that they are enabled
// by block timestamp.
// Note: the return value does not include the native precompiles [nativeAssetCall] and [nativeAssetBalance].
// These are handled in [evm.precompile] directly.
func (c *ChainConfig) enabledStatefulPrecompiles() []precompile.StatefulPrecompileConfig {
	statefulPrecompileConfigs := make([]precompile.StatefulPrecompileConfig, 0)

	return statefulPrecompileConfigs
}

// CheckConfigurePrecompiles checks if any of the precompiles specified in the chain config are enabled by the block
// transition from [parentTimestamp] to the timestamp set in [blockContext]. If this is the case, it calls [Configure]
// to apply the necessary state transitions for the upgrade.
// This function is called:
// - within genesis setup to configure the starting state for precompiles enabled at genesis,
// - during block processing to update the state before processing the given block.
func (c *ChainConfig) CheckConfigurePrecompiles(parentTimestamp *big.Int, blockContext precompile.BlockContext, statedb precompile.StateDB) {
	// Iterate the enabled stateful precompiles and configure them if needed
	for _, config := range c.enabledStatefulPrecompiles() {
		precompile.CheckConfigure(c, parentTimestamp, blockContext, config, statedb)
	}
}
