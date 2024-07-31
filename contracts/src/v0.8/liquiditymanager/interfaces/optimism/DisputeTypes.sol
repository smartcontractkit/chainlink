// SPDX-License-Identifier: MIT
// Copied from https://github.com/ethereum-optimism/optimism/blob/v1.7.0/packages/contracts-bedrock/src/libraries/DisputeTypes.sol
pragma solidity ^0.8.0;

/// @notice A `GameType` represents the type of game being played.
type GameType is uint32;

/// @notice A `GameId` represents a packed 1 byte game ID, an 11 byte timestamp, and a 20 byte address.
/// @dev The packed layout of this type is as follows:
/// ┌───────────┬───────────┐
/// │   Bits    │   Value   │
/// ├───────────┼───────────┤
/// │ [0, 8)    │ Game Type │
/// │ [8, 96)   │ Timestamp │
/// │ [96, 256) │ Address   │
/// └───────────┴───────────┘
type GameId is bytes32;

/// @notice A dedicated timestamp type.
type Timestamp is uint64;

/// @notice A claim represents an MPT root representing the state of the fault proof program.
type Claim is bytes32;
