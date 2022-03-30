// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

// Compound Finance's oracle interface
interface UniswapAnchoredView {
  function price(string memory symbol) external view returns (uint256);
}
