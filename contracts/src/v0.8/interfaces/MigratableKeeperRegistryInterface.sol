// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

import "../UpkeepFormat.sol";

interface MigratableKeeperRegistryInterface {
  function migrateUpkeeps(uint256[] calldata upkeepIDs, address destination) external;

  function receiveUpkeeps(bytes calldata encodedUpkeeps) external;

  function upkeepTranscoderVersion() external returns (UpkeepFormat version);
}
