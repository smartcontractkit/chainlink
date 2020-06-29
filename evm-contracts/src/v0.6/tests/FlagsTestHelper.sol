pragma solidity ^0.6.0;

import "../dev/Flags.sol";

contract FlagsTestHelper {
  Flags public flags;

  constructor(
    address flagsContract
  )
    public
  {
    flags = Flags(flagsContract);
  }

  function getFlag(
    address target
  )
    public
    view
    returns(bool)
  {
    return flags.getFlag(target);
  }

}
