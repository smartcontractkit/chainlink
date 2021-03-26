// SPDX-License-Identifier: MIT
pragma solidity ^0.7.0;

interface UpkeepInterface {
  function checkUpkeep(
    bytes calldata data
  )
    external
    returns (
      bool,
      bytes memory
    );

  function performUpkeep(
    bytes calldata data
  ) external;
}
