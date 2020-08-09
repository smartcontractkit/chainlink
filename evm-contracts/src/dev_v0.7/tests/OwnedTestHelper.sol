pragma solidity ^0.7.0;

import "../Owned.sol";

contract OwnedTestHelper is Owned {

  event Here();

  function modifierOnlyOwner()
    public
    onlyOwner()
  {
    emit Here();
  }

}
