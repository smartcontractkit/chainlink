pragma solidity 0.4.24;

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
    if (_a > 0 && _b > 0) {
      require(c > _a, "SafeMath: addition overflow");
    } else if (_a < 0 && _b < 0) {
      require(c < _a, "SafeMath: addition overflow");
    }
    return c;
  }

}
