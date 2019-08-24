pragma solidity 0.4.24;

import { SignedSafeMath as SignedSafeMath_Chainlink } "../vendor/SignedSafeMath.sol";

contract ConcreteSignedSafeMath {
  using SignedSafeMath_Chainlink for int256;

  function testAdd(int256 _a, int256 _b)
    external
    returns (int256)
  {
    return _a.add(_b);
  }

}
