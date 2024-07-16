// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

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

  function _report(
    uint256 /* initialGas */,
    address /* transmitter */,
    uint8 /* signerCount */,
    address[MAX_NUM_ORACLES] memory /* signers */,
    bytes calldata /* report */
  ) internal virtual override {
    revert ReportingUnsupported();
  }
}
