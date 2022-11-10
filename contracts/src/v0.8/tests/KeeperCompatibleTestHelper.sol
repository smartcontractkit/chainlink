// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../KeeperCompatible.sol";

contract KeeperCompatibleTestHelper is KeeperCompatible {
  function checkUpkeep(bytes calldata) external override returns (bool, bytes memory) {}

  function performUpkeep(bytes calldata) external override {}

  function verifyCannotExecute() public view cannotExecute {}
}
