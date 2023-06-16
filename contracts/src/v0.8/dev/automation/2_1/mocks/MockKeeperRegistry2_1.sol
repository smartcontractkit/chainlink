// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import '../interfaces/IAutomationRegistryConsumer.sol';

contract MockKeeperRegistry2_1 is IAutomationRegistryConsumer {
  uint256 balance;
  uint256 minBalance;

  constructor(){}

  function getBalance(uint256 id) external override view returns (uint256 balance){
    return balance;
  }

  function getMinBalance(uint256 id) external override view returns (uint96 minBalance){
    return minBalance;
  }

  function cancelUpkeep(uint256 id) override external{}

  function pauseUpkeep(uint256 id) override external{}

  function unpauseUpkeep(uint256 id) override external{}

  function updateCheckData(uint256 id, bytes calldata newCheckData) override external{}

  function addFunds(uint256 id, uint96 amount) override external{}

  function withdrawFunds(uint256 id, address to) override external{}

}