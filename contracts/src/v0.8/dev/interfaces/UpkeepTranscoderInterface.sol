// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

interface UpkeepTranscoderInterface {
  function transcodeUpkeeps(
    address fromRegistry,
    address toRegistry,
    bytes calldata encodedUpkeeps
  ) external view returns (bytes memory);
}
