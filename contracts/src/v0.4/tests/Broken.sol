pragma solidity 0.4.24;

// Broken is a contract to aid debugging and testing reverting calls during development.
contract Broken {

  function revertWithMessage(string memory message) public {
    require(false, message);
  }

  function revert() public {
    require(false);
  }
}
