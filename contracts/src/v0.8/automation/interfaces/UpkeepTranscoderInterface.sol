// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {UpkeepFormat} from "../UpkeepFormat.sol";

// solhint-disable-next-line interface-starts-with-i
interface UpkeepTranscoderInterface {
  function transcodeUpkeeps(
    UpkeepFormat fromVersion,
    UpkeepFormat toVersion,
    bytes calldata encodedUpkeeps
  ) external view returns (bytes memory);
}
