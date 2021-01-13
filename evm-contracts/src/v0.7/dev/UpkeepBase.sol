// SPDX-License-Identifier: MIT
pragma solidity ^0.7.0;

contract UpkeepBase {

  modifier cannotExecute()
  {
    require(tx.origin == address(0), "only for simulated backend");
    _;
  }

}
