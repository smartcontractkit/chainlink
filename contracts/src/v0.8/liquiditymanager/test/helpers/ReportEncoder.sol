// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

import {ILiquidityManager} from "../../interfaces/ILiquidityManager.sol";

/// @dev this is needed to generate the types to help encode the report offchain
abstract contract ReportEncoder is ILiquidityManager {
  /// @dev exposed so that we can encode the report for OCR offchain
  function exposeForEncoding(ILiquidityManager.LiquidityInstructions memory instructions) public pure {}
}
