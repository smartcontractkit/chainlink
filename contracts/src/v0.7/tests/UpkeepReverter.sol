// SPDX-License-Identifier: MIT
pragma solidity ^0.7.0;

import "../KeeperCompatible.sol";

contract UpkeepReverter is KeeperCompatible {
  function checkUpkeep(bytes calldata data)
    public
    view
    override
    cannotExecute
    returns (bool callable, bytes calldata executedata)
  {
    require(false, "!working");
    return (true, data);
  }

  function performUpkeep(bytes calldata) external pure override {
    require(false, "!working");
  }
}
