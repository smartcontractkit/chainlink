pragma solidity ^0.6.0;

import "../dev/SignedSafeMath.sol";

contract ConcreteSignedSafeMath {
  using SignedSafeMath for int256;

  function testAdd(int256 _a, int256 _b)
    external
    pure
    returns (int256)
  {
    return _a.add(_b);
  }

  function testMul(int256 _a, int256 _b)
    external
    pure
    returns (int256)
  {
    return _a.mul(_b);
  }

  function testDiv(int256 _a, int256 _b)
    external
    pure
    returns (int256)
  {
    return _a.div(_b);
  }
}
