// SPDX-License-Identifier: MIT
pragma solidity 0.6.6;

contract GasGuzzler {
  fallback() external payable {
    while (true) {
    }
  }
}

