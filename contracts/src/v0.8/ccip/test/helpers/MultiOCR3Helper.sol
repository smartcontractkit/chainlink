// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {MultiOCR3Base} from "../../ocr/MultiOCR3Base.sol";

contract MultiOCR3Helper is MultiOCR3Base {
  function transmit(
    uint8 ocrPluginType,
    bytes32[3] calldata reportContext,
    bytes calldata report,
    bytes32[] calldata rs,
    bytes32[] calldata ss,
    bytes32 rawVs
  ) external {
    _transmit(ocrPluginType, reportContext, report, rs, ss, rawVs);
  }

  // TODO: revisit support for different transmit function sigs:
  // /// @dev test transmit function that has extra msg.data
  // function transmit2(
  //   uint256 customArg,
  //   uint8 ocrPluginType,
  //   bytes32[3] calldata reportContext,
  //   bytes calldata report,
  //   bytes32[] calldata rs,
  //   bytes32[] calldata ss,
  //   bytes32 rawVs
  // ) external {
  //   _transmit(ocrPluginType, reportContext, report, rs, ss, rawVs);
  // }

  // /// @dev test transmit function that has reduced args
  // function transmit3(
  //   uint8 ocrPluginType,
  //   bytes32[3] calldata reportContext,
  //   bytes calldata report
  // ) external {
  //   bytes32[] memory emptySigs = new bytes32[](0);
  //   _transmit(ocrPluginType, reportContext, report, emptySigs, emptySigs, bytes32(""));
  // }

  function typeAndVersion() public pure override returns (string memory) {
    return "MultiOCR3BaseHelper 1.0.0";
  }
}
