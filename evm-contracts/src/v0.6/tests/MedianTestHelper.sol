// SPDX-License-Identifier: MIT
pragma solidity ^0.6.0;

import "../Median.sol";

contract MedianTestHelper {
  function publicGet(int256[] memory list)
    public
    pure
    returns (int256)
  {
    return Median.calculate(list);
  }

  function publicQuickselectTwo(int256[] memory list, uint256 k1, uint256 k2)
    public
    pure
    returns (int256, int256)
  {
    return Median.quickselectTwo(list, 0, list.length - 1, k1, k2);
  }
}
