pragma solidity ^0.6.0;

library SignedSafeMath {

  /**
   * @dev Adds two int256s and makes sure the result doesn't overflow. Signed
   * integers aren't supported by the SafeMath library, thus this method
   * @param _a The first number to be added
   * @param _a The second number to be added
   */
  function add(int256 _a, int256 _b)
    internal
    pure
    returns (int256)
  {
    // solium-disable-next-line zeppelin/no-arithmetic-operations
    int256 c = _a + _b;
    require((_b >= 0 && c >= _a) || (_b < 0 && c < _a), "SignedSafeMath: addition overflow");

    return c;
  }

  /**
   * @notice Computes average of two signed integers, ensuring that the computation
   * doesn't overflow.
   * @dev If the result is not an integer, it is rounded towards zero. For example,
   * avg(-3, -4) = -3
   */
  function avg(int256 _a, int256 _b)
    internal
    pure
    returns (int256)
  {
    int256 remainder = (_a % 2 + _b % 2) / 2;
    return add(add(_a / 2, _b / 2), remainder);
  }
}
