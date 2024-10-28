// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.24;

import {MultiOCR3Base} from "../../ocr/MultiOCR3Base.sol";

contract MultiOCR3Helper is MultiOCR3Base {
  event AfterConfigSet(uint8 ocrPluginType);

  /// @dev OCR plugin type used for transmit.
  ///      Defined in storage since it cannot be passed as calldata due to strict transmit checks
  uint8 internal s_transmitOcrPluginType;

  function setTransmitOcrPluginType(
    uint8 ocrPluginType
  ) external {
    s_transmitOcrPluginType = ocrPluginType;
  }

  /// @dev transmit function with signatures
  function transmitWithSignatures(
    bytes32[3] calldata reportContext,
    bytes calldata report,
    bytes32[] calldata rs,
    bytes32[] calldata ss,
    bytes32 rawVs
  ) external {
    _transmit(s_transmitOcrPluginType, reportContext, report, rs, ss, rawVs);
  }

  /// @dev transmit function with no signatures
  function transmitWithoutSignatures(bytes32[3] calldata reportContext, bytes calldata report) external {
    bytes32[] memory emptySigs = new bytes32[](0);
    _transmit(s_transmitOcrPluginType, reportContext, report, emptySigs, emptySigs, bytes32(""));
  }

  function getOracle(uint8 ocrPluginType, address oracleAddress) external view returns (Oracle memory) {
    return s_oracles[ocrPluginType][oracleAddress];
  }

  function typeAndVersion() public pure override returns (string memory) {
    return "MultiOCR3BaseHelper 1.0.0";
  }

  function _afterOCR3ConfigSet(
    uint8 ocrPluginType
  ) internal virtual override {
    emit AfterConfigSet(ocrPluginType);
  }
}
