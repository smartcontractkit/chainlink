pragma solidity 0.4.24;

import "../vendor/SignedSafeMath.sol";

contract ConcreteSignedSafeMath {
  using SignedSafeMath for int256;

  function testAdd(int256 _a, int256 _b)
    external
    returns (int256)
  {
    return _a.add(_b);
  }

}
