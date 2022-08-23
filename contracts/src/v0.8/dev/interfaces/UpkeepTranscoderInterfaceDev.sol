// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../UpkeepFormatDev.sol";

interface UpkeepTranscoderInterfaceDev {
  function transcodeUpkeeps(
    UpkeepFormatDev fromVersion,
    UpkeepFormatDev toVersion,
    bytes calldata encodedUpkeeps
  ) external view returns (bytes memory);
}
