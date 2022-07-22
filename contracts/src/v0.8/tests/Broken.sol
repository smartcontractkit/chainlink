pragma solidity ^0.8.0;

// Broken is a contract to aid debugging and testing reverting calls during development.
contract Broken {
  error Unauthorized(string reason, int reason2);
  error Empty();

  function revertWithCustomError() pure public {
    revert Unauthorized("param", 121);
  }

  function revertWithMessage(string memory message) pure public {
    require(false, message);
  }

  function revertSilently() pure public {
    require(false);
  }
  
}
