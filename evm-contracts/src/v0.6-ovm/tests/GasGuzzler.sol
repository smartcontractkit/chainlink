pragma solidity ^0.6.0;

contract GasGuzzler {
  fallback() external payable {
    while (true) {
    }
  }
}

