// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;
pragma abicoder v2;

import "../../v0.8/dev/interfaces/UpkeepTranscoderInterface.sol";

interface MigratableKeeperRegistryInterface {
  function migrateUpkeeps(uint256[] calldata upkeepIDs, address destination) external;

  function receiveUpkeeps(bytes calldata encodedUpkeeps) external;

 function upkeepTranscoderVersion() external returns (UpkeepTranscoderVersion version);
}
