pragma solidity 0.4.24;

library SignedSafeMath {

  /**
  * @dev Multiplies two numbers, throws on overflow.
  */
  function mul(int256 _a, int256 _b)
    internal
    pure
    returns (int256 c)
  {
    if (_a == 0) {
      return 0;
    }

    c = _a * _b;
    require(c / _a == _b, "SignedSafeMath: .mul overfow");
    return c;
  }

  /**
  * @dev Integer division of two numbers, truncating the quotient.
  */
  function div(int256 _a, int256 _b)
    internal
    pure
    returns (int256)
  {
    // assert(_b > 0); // Solidity automatically throws when dividing by 0
    // int256 c = _a / _b;
    // assert(_a == _b * c + _a % _b); // There is no case in which this doesn't hold
    return _a / _b;
  }

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
    int256 c = _a + _b;
    require((_b >= 0 && c >= _a) || (_b < 0 && c < _a), "SignedSafeMath: .add overflow");

    return c;
  }
}
