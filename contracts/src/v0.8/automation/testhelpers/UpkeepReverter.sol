// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {AutomationCompatible} from "../AutomationCompatible.sol";

contract UpkeepReverter is AutomationCompatible {
  function checkUpkeep(
    bytes calldata data
  ) public view override cannotExecute returns (bool callable, bytes calldata executedata) {
    require(false, "!working");
    return (true, data);
  }

  function performUpkeep(bytes calldata) external pure override {
    require(false, "!working");
  }
}
