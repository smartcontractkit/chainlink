pragma solidity ^0.6.0;

import "../Median.sol";

contract MedianTestHelper {

  function publicGet(int256[] memory _list)
    public
    pure
    returns (int256)
  {
    return Median.calculate(_list);
  }

}
