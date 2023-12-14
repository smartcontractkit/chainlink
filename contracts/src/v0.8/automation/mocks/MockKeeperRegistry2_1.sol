// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

import "../interfaces/IAutomationRegistryConsumer.sol";

contract MockKeeperRegistry2_1 is IAutomationRegistryConsumer {
  uint96 balance;
  uint96 minBalance;

  constructor() {}

  function getBalance(uint256 id) external view override returns (uint96) {
    return balance;
  }

  function getMinBalance(uint256 id) external view override returns (uint96) {
    return minBalance;
  }

  function cancelUpkeep(uint256 id) external override {}

  function pauseUpkeep(uint256 id) external override {}

  function unpauseUpkeep(uint256 id) external override {}

  function updateCheckData(uint256 id, bytes calldata newCheckData) external {}

  function addFunds(uint256 id, uint96 amount) external override {}

  function withdrawFunds(uint256 id, address to) external override {}
}
