// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.0;

import {IRebalancer} from "../../interfaces/IRebalancer.sol";

/// @dev this is needed to generate the types to help encode the report offchain
abstract contract RebalancerReportEncoder is IRebalancer {
  /// @dev exposed so that we can encode the report for OCR offchain
  function exposeForEncoding(IRebalancer.LiquidityInstructions memory instructions) public pure {}
}
