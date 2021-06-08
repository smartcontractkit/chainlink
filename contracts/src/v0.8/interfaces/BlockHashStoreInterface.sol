// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface BlockHashStoreInterface {
  function getBlockhash(uint256 number) external view returns (bytes32);
}
