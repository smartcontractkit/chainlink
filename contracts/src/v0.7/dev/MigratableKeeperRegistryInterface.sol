// SPDX-License-Identifier: MIT

pragma solidity ^0.7.0;
pragma abicoder v2;

interface MigratableKeeperRegistryInterface {
  function migrateUpkeeps(uint256[] calldata upkeepIDs, address destination) external;

  function receiveUpkeeps(bytes calldata encodedUpkeeps) external;
}
