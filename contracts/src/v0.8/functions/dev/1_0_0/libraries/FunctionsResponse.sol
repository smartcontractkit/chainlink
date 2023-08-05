// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

// @title Library of types that are used during fulfillment of a Functions request
library FunctionsResponse {
  enum FulfillResult {
    USER_SUCCESS, // 0
    USER_ERROR, // 1
    INVALID_REQUEST_ID, // 2
    COST_EXCEEDS_COMMITMENT, // 3
    INSUFFICIENT_GAS_PROVIDED, // 4
    SUBSCRIPTION_BALANCE_INVARIANT_VIOLATION, // 5
    INVALID_COMMITMENT // 6
  }

  struct Commitment {
    uint96 adminFee; // -----------┐
    address coordinator; // -------┘
    address client; // ------------┐
    uint64 subscriptionId; //      |
    uint32 callbackGasLimit; // ---┘
    uint96 estimatedTotalCostJuels;
    uint40 timeoutTimestamp;
    bytes32 requestId;
    uint80 donFee;
    uint40 gasOverheadBeforeCallback;
    uint40 gasOverheadAfterCallback;
  }
}
