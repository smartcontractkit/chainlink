// SPDX-License-Identifier: MIT
// Adapted from https://github.com/OpenZeppelin/openzeppelin-contracts/blob/97894a140d2a698e5a0f913648a8f56d62277a70/contracts/math/SignedSafeMath.sol

pragma solidity ^0.6.0;

library CheckedMath {

  int256 constant internal INT256_MIN = -2**255;

  /**
   * @dev Subtracts two signed integers, returns false 2nd param on overflow.
   */
  function add(
    int256 a,
    int256 b
  )
    internal
    pure
    returns (int256 result, bool ok)
  {
    int256 c = a + b;
    if ((b >= 0 && c < a) || (b < 0 && c >= a)) return (0, false);

    return (c, true);
  }

  /**
   * @dev Subtracts two signed integers, returns false 2nd param on overflow.
   */
  function sub(
    int256 a,
    int256 b
  )
    internal
    pure
    returns (int256 result, bool ok)
  {
    int256 c = a - b;
    if ((b < 0 && c <= a) || (b >= 0 && c > a)) return (0, false);

    return (c, true);
  }


  /**
   * @dev Multiplies two signed integers, returns false 2nd param on overflow.
   */
  function mul(
    int256 a,
    int256 b
  )
    internal
    pure
    returns (int256 result, bool ok)
  {
    // Gas optimization: this is cheaper than requiring 'a' not being zero, but the
    // benefit is lost if 'b' is also tested.
    // See: https://github.com/OpenZeppelin/openzeppelin-contracts/pull/522
    if (a == 0) return (0, true);
    if (a == -1 && b == INT256_MIN) return (0, false);

    int256 c = a * b;
    if (!(c / a == b)) return (0, false);

    return (c, true);
  }

  /**
   * @dev Divides two signed integers, returns false 2nd param on overflow.
   */
  function div(
    int256 a,
    int256 b
  )
    internal
    pure
    returns (int256 result, bool ok)
  {
    if (b == 0) return (0, false);
    if (b == -1 && a == INT256_MIN) return (0, false);

    int256 c = a / b;

    return (c, true);
  }

}
