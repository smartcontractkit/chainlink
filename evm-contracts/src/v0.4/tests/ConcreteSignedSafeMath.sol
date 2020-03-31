pragma solidity 0.4.24;

import "../vendor/SignedSafeMath.sol";

contract ConcreteSignedSafeMath {
  using SignedSafeMath for int256;

  function testAdd(int256 _a, int256 _b)
    external
    pure
    returns (int256)
  {
    return _a.add(_b);
  }

  function testSub(int256 _a, int256 _b)
    external
    pure
    returns (int256)
  {
    return _a.sub(_b);
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
