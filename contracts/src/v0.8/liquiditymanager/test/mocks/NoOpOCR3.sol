// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

import {OCR3Base} from "../../ocr/OCR3Base.sol";

// NoOpOCR3 is a mock implementation of the OCR3Base contract that does nothing
// This is so that we can generate gethwrappers for the contract and use the OCR3 ABI in
// Go code.
contract NoOpOCR3 is OCR3Base {
  // solhint-disable-next-line chainlink-solidity/all-caps-constant-storage-variables
  string public constant override typeAndVersion = "NoOpOCR3 1.0.0";

  constructor() OCR3Base() {}

  function _report(bytes calldata, uint64) internal override {
    // do nothing
  }
}
