// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../interfaces/AutomationCompatibleInterface.sol";
import "../ConfirmedOwner.sol";
import "./AutoCompoundV3Interface.sol";

contract FurioAutomator is AutomationCompatibleInterface, ConfirmedOwner {
  AutoCompoundV3Interface public s_autoCompounder;
  uint256[] public s_checkers; // can be replaced with a single number if checkers are always starting from 1

  constructor(AutoCompoundV3Interface autoCompounder, uint256[] memory checkers) ConfirmedOwner(msg.sender) {
    s_autoCompounder = autoCompounder;
    s_checkers = checkers;
  }

  function setOptionMarket(AutoCompoundV3Interface autoCompounder) external onlyOwner {
    s_autoCompounder = autoCompounder;
  }

  function setCheckers(uint256[] calldata checkers) external onlyOwner {
    s_checkers = checkers;
  }

  function checkUpkeep(bytes calldata checkData)
    external
    override
    returns (bool upkeepNeeded, bytes memory performData)
  {
    // checks all checkers for each block, return true for the first eligible one
    for (uint256 i = 0; i < s_checkers.length; i++) {
      // their checker func can be simplified to only encode the ID? or it's cheaper to make a static call with selector
      (bool canExec, bytes memory execPayload) = s_autoCompounder.checker(s_checkers[i]);
      if (canExec) {
        return (canExec, execPayload);
      }
    }
    return (false, "0x");

    // sudo-randomly checks a checker per block, maybe need a better source for randomness. but i already see the usage
    // of block.timestamp in Furio contract so maybe this is good enough
    //    uint256 idx = block.number % s_checkers.length;
    //    uint256 checker = s_checkers[idx];
    //    return s_autoCompounder.checker(idx);
  }

  function performUpkeep(bytes calldata performData) external override {
    // I found it difficult and/or super expensive to re-check the conditions here based on AutoCompoundV3 contract.
    // This may cause this performUpkeep to be triggered at unwanted times.
    uint256 id = abi.decode(performData[4:], (uint256));
    s_autoCompounder.compound(id);
  }
}
