pragma solidity ^0.8.0;

// Broken is a contract to aid debugging and testing reverting calls during development.
contract Broken {
  error Unauthorized(string reason, int256 reason2);

  function revertWithCustomError() public pure {
    revert Unauthorized("param", 121);
  }

  function revertWithMessage(string memory message) public pure {
    require(false, message);
  }

  function revertSilently() public pure {
    require(false);
  }
}
