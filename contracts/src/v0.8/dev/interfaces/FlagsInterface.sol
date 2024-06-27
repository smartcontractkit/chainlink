// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

interface FlagsInterface {
  function getFlag(address) external view returns (bool);

  function getFlags(address[] calldata) external view returns (bool[] memory);

  function raiseFlag(address) external;

  function raiseFlags(address[] calldata) external;

  function lowerFlag(address) external;

  function lowerFlags(address[] calldata) external;

  function setRaisingAccessController(address) external;

  function setLoweringAccessController(address) external;
}
