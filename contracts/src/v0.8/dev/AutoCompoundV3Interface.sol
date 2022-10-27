// SPDX-License-Identifier: MIT

pragma solidity 0.8.6;

interface AutoCompoundV3Interface {
  function checker(uint256 checker_) external view returns (bool canExec, bytes memory execPayload);

  function compound(uint256 id_) external;
}
