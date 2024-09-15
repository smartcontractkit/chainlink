// SPDX-License-Identifier: MIT
// Copied from https://github.com/ethereum-optimism/optimism/blob/v1.7.0/packages/contracts-bedrock/src/dispute/DisputeGameFactory.sol
pragma solidity ^0.8.0;

import {GameType, GameId, Timestamp, Claim} from "./DisputeTypes.sol";

interface IOptimismDisputeGameFactory {
  /// @notice Information about a dispute game found in a `findLatestGames` search.
  struct GameSearchResult {
    uint256 index;
    GameId metadata;
    Timestamp timestamp;
    Claim rootClaim;
    bytes extraData;
  }

  /// @notice Finds the `_n` most recent `GameId`'s of type `_gameType` starting at `_start`. If there are less than
  ///         `_n` games of type `_gameType` starting at `_start`, then the returned array will be shorter than `_n`.
  /// @param _gameType The type of game to find.
  /// @param _start The index to start the reverse search from.
  /// @param _n The number of games to find.
  function findLatestGames(
    GameType _gameType,
    uint256 _start,
    uint256 _n
  ) external view returns (GameSearchResult[] memory games_);

  /// @notice The total number of dispute games created by this factory.
  /// @return gameCount_ The total number of dispute games created by this factory.
  function gameCount() external view returns (uint256 gameCount_);
}
