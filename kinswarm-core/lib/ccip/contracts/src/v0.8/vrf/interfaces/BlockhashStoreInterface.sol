// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

// solhint-disable-next-line interface-starts-with-i
interface BlockhashStoreInterface {
  function getBlockhash(uint256 number) external view returns (bytes32);
}
