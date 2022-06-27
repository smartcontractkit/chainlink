// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface iBlockhashStore {
  function getBlockhash(uint256 number) external view returns (bytes32);
}
