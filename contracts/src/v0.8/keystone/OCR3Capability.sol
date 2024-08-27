// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {OCR2Base} from "./ocr/OCR2Base.sol";

// OCR2Base provides config management compatible with OCR3
contract OCR3Capability is OCR2Base {
  error ReportingUnsupported();

  constructor() OCR2Base() {}

  function typeAndVersion() external pure override returns (string memory) {
    return "Keystone 1.0.0";
  }

  function _beforeSetConfig(uint8 /* _f */, bytes memory /* _onchainConfig */) internal override {
    // no-op
  }

  function transmit(
    // NOTE: If these parameters are changed, expectedMsgDataLength and/or
    // TRANSMIT_MSGDATA_CONSTANT_LENGTH_COMPONENT need to be changed accordingly
    bytes32[3] calldata /* reportContext */,
    bytes calldata /* report */,
    bytes32[] calldata /* rs */,
    bytes32[] calldata /* ss */,
    bytes32 /* rawVs */ // signatures
  ) external override {
    revert ReportingUnsupported();
  }
}
