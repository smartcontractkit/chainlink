// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../AutomationCompatible.sol";

contract KeeperCompatibleTestHelper is AutomationCompatible {
  function checkUpkeep(bytes calldata) external override returns (bool, bytes memory) {}

  function performUpkeep(bytes calldata) external override {}

  function testCannotExecute() public view cannotExecute {}
}
