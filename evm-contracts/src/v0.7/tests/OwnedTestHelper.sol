pragma solidity ^0.7.0;

import "../dev/Owned.sol";

contract OwnedTestHelper is Owned {

  event Here();

  constructor(address owner) Owned(owner) {}

  function modifierOnlyOwner()
    public
    onlyOwner()
  {
    emit Here();
  }

}
