pragma solidity 0.6.2;

import "../dev/Owned.sol";

contract OwnedTestHelper is Owned {

  event Here();

  function modifierOnlyOwner()
    public
    onlyOwner()
  {
    emit Here();
  }

}
