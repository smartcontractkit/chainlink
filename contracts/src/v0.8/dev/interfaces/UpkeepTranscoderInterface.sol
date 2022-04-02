// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

enum UpkeepTranscoderVersion {
  V1
}

interface UpkeepTranscoderInterface {
  function transcodeUpkeeps(
    UpkeepTranscoderVersion fromVersion,
    UpkeepTranscoderVersion toVersion,
    bytes calldata encodedUpkeeps
  ) external view returns (bytes memory);
}
