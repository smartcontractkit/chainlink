// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

// solhint-disable-next-line interface-starts-with-i
interface UpkeepTranscoderInterfaceV2 {
  function transcodeUpkeeps(
    uint8 fromVersion,
    uint8 toVersion,
    bytes calldata encodedUpkeeps
  ) external view returns (bytes memory);
}
