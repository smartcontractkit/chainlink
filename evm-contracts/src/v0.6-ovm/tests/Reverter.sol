pragma solidity ^0.6.0;

contract Reverter {

  fallback() external payable {
    require(false, "Raised by Reverter.sol");
  }

}
