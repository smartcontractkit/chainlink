// SPDX-License-Identifier: MIT
pragma solidity ^0.6.0;

import "../CheckedMath.sol";

contract CheckedMathTestHelper {
  using CheckedMath for int256;

  function add(int256 a, int256 b)
    external
    pure
    returns (int256 result, bool ok)
  {
    return a.add(b);
  }

  function sub(int256 a, int256 b)
    external
    pure
    returns (int256 result, bool ok)
  {
    return a.sub(b);
  }

  function mul(int256 a, int256 b)
    external
    pure
    returns (int256 result, bool ok)
  {
    return a.mul(b);
  }

  function div(int256 a, int256 b)
    external
    pure
    returns (int256 result, bool ok)
  {
    return a.div(b);
  }

}
