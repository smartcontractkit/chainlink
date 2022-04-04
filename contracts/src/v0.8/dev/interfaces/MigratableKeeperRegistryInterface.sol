// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

interface MigratableKeeperRegistryInterface {
  function migrateUpkeeps(uint256[] calldata upkeepIDs, address destination) external;

  function receiveUpkeeps(bytes calldata encodedUpkeeps) external;

  function upkeepTranscoderVersion() external returns (uint8 version);
}
