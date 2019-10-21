pragma solidity ^0.5.0;

import "../Median.sol";

contract ConcreteMedian {

  function publicGet(int256[] memory _list)
    public
    returns (int256)
  {
    return Median.get(_list);
  }

}
